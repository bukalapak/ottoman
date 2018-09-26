package cache

import (
	"io/ioutil"
	"net/http"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// Fetcher is the interface for getting cache data from remote backend based on given key(s).
type Fetcher interface {
	Fetch(key string, r *http.Request) (body []byte, statusCode int, err error)
	FetchMulti(keys []string, r *http.Request) (bodies map[string][]byte, statusCodes map[string]int, err error)
}

// Resolver is the interface for resolving cache key to http request.
type Resolver interface {
	Resolve(key string, r *http.Request) (*http.Request, error)
	ResolveRequest(r *http.Request) (*http.Request, error)
}

// RemoteProvider enhances Provider with remote functionalities.
type RemoteProvider interface {
	Provider
	Fetcher
}

// RemoteOption is the configuration option for the RemoteProvider.
type RemoteOption struct {
	Transport http.RoundTripper
	Timeout   time.Duration
	Resolver  Resolver
}

func (n RemoteOption) httpClient() *http.Client {
	return &http.Client{
		Transport: n.httpTransport(),
		Timeout:   n.httpTimeout(),
	}
}

func (n RemoteOption) httpTransport() http.RoundTripper {
	if n.Transport == nil {
		return http.DefaultTransport
	}

	return n.Transport
}

func (n RemoteOption) httpTimeout() time.Duration {
	if n.Timeout == 0 {
		return 30 * time.Second
	}

	return n.Timeout
}

type remoteProvider struct {
	Provider
	option RemoteOption
}

// NewRemoteProvider returns RemoteProvider from a Provider and RemoteOption.
func NewRemoteProvider(p Provider, opt RemoteOption) RemoteProvider {
	return &remoteProvider{
		Provider: p,
		option:   opt,
	}
}

func (p *remoteProvider) Fetch(key string, r *http.Request) ([]byte, int, error) {
	req, err := p.option.Resolver.Resolve(p.Normalize(key), r)
	if err != nil {
		return nil, 0, err
	}

	return p.fetchRequest(req)
}

func (p *remoteProvider) fetchRequest(r *http.Request) ([]byte, int, error) {
	c := p.option.httpClient()

	resp, err := c.Do(r)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.New("invalid http status: " + resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	return b, resp.StatusCode, err
}

func (p *remoteProvider) FetchMulti(keys []string, r *http.Request) (map[string][]byte, map[string]int, error) {
	ks := p.NormalizeMulti(keys)

	mb := make(map[string][]byte)
	mn := make(map[string]int)

	ec := make(chan error)
	bc := make(chan map[string]map[int][]byte)

	for _, k := range ks {
		go func(key string) {
			z, n, err := p.Fetch(key, r)
			if err != nil {
				ec <- errors.Wrap(err, key)
			} else {
				bc <- map[string]map[int][]byte{
					key: {
						n: z,
					},
				}
			}
		}(k)
	}

	var mrr *multierror.Error

	for i := 0; i < len(ks); i++ {
		select {
		case kb := <-bc:
			for k, v := range kb {
				for n, b := range v {
					mb[k] = b
					mn[k] = n
				}
			}
		case err := <-ec:
			mrr = multierror.Append(mrr, err)
		}
	}

	return mb, mn, mrr.ErrorOrNil()
}
