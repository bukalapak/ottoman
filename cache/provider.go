package cache

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bukalapak/ottoman/encoding/json"
	"go.uber.org/zap"
)

type Engine struct {
	name      string
	engine    Reader
	Prefix    string
	Resolver  Resolver
	Timeout   time.Duration
	Transport http.RoundTripper
	Logger    *zap.Logger
}

func NewProvider(r Reader) Provider {
	return &Engine{engine: r, name: r.Name(), Logger: zap.New(nil)}
}

// Name returns cache backend identifier.
func (s *Engine) Name() string {
	return s.name
}

// Namespace returns cache Prefix
func (s *Engine) Namespace() string {
	return s.Prefix
}

// Read reads cache data on the cache backend based on key supplied.
func (s *Engine) Read(key string) ([]byte, error) {
	return s.engine.Read(s.normalize(key))
}

// ReadMap reads cache data as map[string]interface{}.
// It's also expand any cache identifier with actual cache data.
func (s *Engine) ReadMap(key string) (map[string]interface{}, error) {
	return s.engine.ReadMap(s.normalize(key))
}

// ReadMulti bulk reads multiple cache keys.
func (s *Engine) ReadMulti(keys []string) (map[string][]byte, error) {
	ks := make([]string, len(keys))

	for i := range keys {
		ks[i] = s.normalize(keys[i])
	}

	return s.engine.ReadMulti(ks)
}

func (s *Engine) Fetch(key string, r *http.Request) ([]byte, error) {
	if b, err := s.Read(key); err == nil {
		return b, nil
	}

	req, err := s.Resolver.Resolve(key, r)
	if err != nil {
		return nil, err
	}

	return s.fetch(req)
}

func (s *Engine) FetchMap(key string, r *http.Request) (map[string]interface{}, error) {
	b, err := s.Fetch(key, r)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})

	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Engine) FetchMulti(keys []string, r *http.Request) (map[string][]byte, error) {
	mb, err := s.ReadMulti(keys)
	if err != nil {
		s.Logger.Info("ottoman/cache",
			zap.String("method", "ReadMulti"),
			zap.String("error", err.Error()),
		)
	}

	if len(mb) == 0 {
		mb = map[string][]byte{}
	}

	cs := s.cachedKeys(mb, keys)
	us := s.uncachedKeys(cs, keys)

	mutex := &sync.Mutex{}

	var wg sync.WaitGroup

	for _, k := range us {
		wg.Add(1)

		go func(key string) {
			defer wg.Done()

			z, err := s.Fetch(key, r)
			if err != nil {
				s.Logger.Info("ottoman/cache",
					zap.String("method", "Fetch"),
					zap.String("error", err.Error()),
				)
			}

			mutex.Lock()
			mb[key] = z
			mutex.Unlock()
		}(k)
	}

	wg.Wait()

	return mb, nil
}

func (s *Engine) normalize(key string) string {
	return Normalize(key, s.Prefix)
}

func (s *Engine) fetch(r *http.Request) ([]byte, error) {
	c := s.httpClient()

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	s.Logger.Info("ottoman/cache",
		zap.String("request_url", r.URL.String()),
		zap.Int("request_status", resp.StatusCode),
	)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid http status: " + resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *Engine) httpClient() *http.Client {
	return &http.Client{
		Transport: s.httpTransport(),
		Timeout:   s.Timeout,
	}
}

func (s *Engine) httpTransport() http.RoundTripper {
	if s.Transport != nil {
		return s.Transport
	}

	return http.DefaultTransport
}

func (s *Engine) uncachedKeys(cs, keys []string) []string {
	us := []string{}

	for _, k := range keys {
		if !sliceContains(cs, k) {
			us = append(us, s.normalize(k))
		}
	}

	return us
}

func (s *Engine) cachedKeys(mb map[string][]byte, keys []string) []string {
	cs := []string{}

	for k, b := range mb {
		for _, y := range keys {
			if strings.Contains(k, y) && len(b) != 0 {
				cs = append(cs, s.normalize(y))
			}
		}
	}

	return cs
}

func sliceContains(ss []string, k string) bool {
	for _, v := range ss {
		if v == k {
			return true
		}
	}

	return false
}
