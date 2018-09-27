package cache

import (
	"io/ioutil"
	"net/http"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// FetchInfo is the container for the information data from a backend.
type FetchInfo struct {
	RemoteURL  string
	StatusCode int
}

// Fetcher is the interface for getting cache data from remote backend based on given key(s).
type Fetcher interface {
	Fetch(key string, r *http.Request) ([]byte, *FetchInfo, error)
	FetchMulti(keys []string, r *http.Request) (map[string][]byte, map[string]*FetchInfo, error)
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

func (p *remoteProvider) Fetch(key string, r *http.Request) ([]byte, *FetchInfo, error) {
	req, err := p.option.Resolver.Resolve(p.Normalize(key), r)
	if err != nil {
		return nil, nil, err
	}

	return p.fetchRequest(req)
}

func (p *remoteProvider) fetchRequest(r *http.Request) ([]byte, *FetchInfo, error) {
	c := p.option.httpClient()

	resp, err := c.Do(r)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &FetchInfo{
			RemoteURL:  r.URL.String(),
			StatusCode: resp.StatusCode,
		}, errors.New("invalid http status: " + resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)

	return b, &FetchInfo{
		RemoteURL:  r.URL.String(),
		StatusCode: resp.StatusCode,
	}, nil
}

func (p *remoteProvider) FetchMulti(keys []string, r *http.Request) (map[string][]byte, map[string]*FetchInfo, error) {
	ks := p.NormalizeMulti(keys)

	mb := make(map[string][]byte)
	mn := make(map[string]*FetchInfo)

	type fetchPayload struct {
		body []byte
		info *FetchInfo
	}

	ec := make(chan error)
	bc := make(chan map[string]fetchPayload)

	for _, k := range ks {
		go func(key string) {
			b, n, err := p.Fetch(key, r)
			if err != nil {
				ec <- errors.Wrap(err, key)
			} else {
				bc <- map[string]fetchPayload{
					key: fetchPayload{
						body: b,
						info: n,
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
				mb[k] = v.body
				mn[k] = v.info
			}
		case err := <-ec:
			mrr = multierror.Append(mrr, err)
		}
	}

	return mb, mn, mrr.ErrorOrNil()
}
