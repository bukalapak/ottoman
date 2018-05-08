package cache

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Engine struct {
	engine    Reader
	Prefix    string
	Resolver  Resolver
	Timeout   time.Duration
	Transport http.RoundTripper
	Counter   MetricCounter
	Tracer    BackendTracer
	Logger    *zap.Logger
}

func NewProvider(r Reader) Provider {
	return &Engine{
		engine:    r,
		Counter:   &noopCounter{},
		Tracer:    &noopTracer{},
		Timeout:   30 * time.Second,
		Transport: http.DefaultTransport,
		Logger:    zap.New(nil),
	}
}

// Name returns cache backend identifier.
func (s *Engine) Name() string {
	return s.engine.Name()
}

// Namespace returns cache Prefix
func (s *Engine) Namespace() string {
	return s.Prefix
}

// Read reads cache data on the cache backend based on key supplied.
func (s *Engine) Read(key string) ([]byte, error) {
	b, err := s.engine.Read(s.Normalize(key))
	if err == nil {
		s.Counter.IncrCacheCounter()
	}

	return b, err
}

// ReadMulti bulk reads multiple cache keys.
func (s *Engine) ReadMulti(keys []string) (map[string][]byte, error) {
	mb, err := s.engine.ReadMulti(s.NormalizeMulti(keys))
	if err == nil {
		s.Counter.IncrCacheCounter()
	}

	return mb, err
}

func (s *Engine) Fetch(key string, r *http.Request) ([]byte, error) {
	req, err := s.Resolver.Resolve(key, r)
	if err != nil {
		return nil, err
	}

	b, err := s.fetchRequest(req)
	if err == nil {
		s.Counter.IncrBackendCounter()
	}

	return b, err
}

func (s *Engine) FetchMulti(keys []string, r *http.Request) (map[string][]byte, error) {
	ks := s.NormalizeMulti(keys)
	mb := make(map[string][]byte, len(ks))
	mx := &sync.Mutex{}

	var wg sync.WaitGroup

	for _, k := range ks {
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

			mx.Lock()
			mb[key] = z
			mx.Unlock()
		}(k)
	}

	wg.Wait()

	return mb, nil
}

func (s *Engine) ReadFetch(key string, r *http.Request) ([]byte, error) {
	if b, err := s.Read(key); err == nil {
		return b, nil
	}

	return s.Fetch(key, r)
}

func (s *Engine) ReadFetchMulti(keys []string, r *http.Request) (map[string][]byte, error) {
	mb, err := s.ReadMulti(keys)
	if err != nil {
		s.Logger.Info("ottoman/cache",
			zap.String("method", "ReadMulti"),
			zap.String("error", err.Error()),
		)
	}

	if len(mb) == 0 {
		mb = make(map[string][]byte)
	}

	cs := s.cachedKeys(mb, keys)
	us := s.uncachedKeys(cs, keys)
	mx, _ := s.FetchMulti(us, r)

	for k, v := range mx {
		mb[k] = v
	}

	return mb, nil
}

func (s *Engine) Normalize(key string) string {
	return Normalize(key, s.Prefix)
}

func (s *Engine) NormalizeMulti(keys []string) []string {
	ks := make([]string, len(keys))

	for i := range keys {
		ks[i] = s.Normalize(keys[i])
	}

	return ks
}

func (s *Engine) fetchRequest(r *http.Request) ([]byte, error) {
	now := time.Now()
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

	s.Tracer.BackendLatency(r.URL.String(), resp.StatusCode, time.Since(now))

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid http status: " + resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *Engine) httpClient() *http.Client {
	return &http.Client{
		Transport: s.Transport,
		Timeout:   s.Timeout,
	}
}

func (s *Engine) uncachedKeys(cs, keys []string) []string {
	us := []string{}

	for _, k := range keys {
		if !sliceContains(cs, k) {
			us = append(us, s.Normalize(k))
		}
	}

	return us
}

func (s *Engine) cachedKeys(mb map[string][]byte, keys []string) []string {
	cs := []string{}

	for k, b := range mb {
		for _, y := range keys {
			if strings.Contains(k, y) && len(b) != 0 {
				cs = append(cs, s.Normalize(y))
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

type noopCounter struct{}

func (c *noopCounter) IncrCacheCounter()   {}
func (c *noopCounter) IncrBackendCounter() {}

type noopTracer struct{}

func (c *noopTracer) BackendLatency(route string, code int, n time.Duration) {}
