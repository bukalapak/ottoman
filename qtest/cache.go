package qtest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bukalapak/ottoman/cache"
	"github.com/bukalapak/ottoman/encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type CacheMap struct {
	keymap map[string]string
}

func NewCacheReader(keymap map[string]string) cache.Reader {
	return &CacheMap{keymap: keymap}
}

func (q *CacheMap) Name() string {
	return "qtest/cache:reader"
}

func (q *CacheMap) Read(key string) ([]byte, error) {
	if fname, ok := q.keymap[key]; ok {
		return ioutil.ReadFile(fname)
	}

	return nil, fmt.Errorf("unknown cache key: %s", key)
}

func (q *CacheMap) ReadMap(key string) (map[string]interface{}, error) {
	b, err := q.Read(key)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)

	return m, err
}

func (q *CacheMap) ReadMulti(keys []string) (map[string][]byte, error) {
	z := make(map[string][]byte, len(keys))

	for _, key := range keys {
		if b, err := q.Read(key); err == nil {
			z[key] = b
		} else {
			z[key] = nil
		}
	}

	return z, nil
}

func NewCacheBackend(reqmap map[string]string) *httptest.Server {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if fixture, ok := reqmap[r.URL.RequestURI()]; ok {
			if b, err := ioutil.ReadFile(fixture); err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write(b)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	}

	return httptest.NewServer(http.HandlerFunc(fn))
}

type Metric struct {
	cacheLatency *prometheus.HistogramVec
	registry     *prometheus.Registry
}

func NewMetric() *Metric {
	m := &Metric{registry: prometheus.NewRegistry()}

	m.cacheLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "cache_latency_seconds",
		Help: "A histogram of the cache latency in seconds.",
	}, []string{"name", "action"})

	m.registry.MustRegister(m.cacheLatency)

	return m
}

func (m *Metric) Registry() *prometheus.Registry {
	return m.registry
}

func (m *Metric) Gather(name string) ([]*dto.Metric, error) {
	gf, err := m.Registry().Gather()
	if err == nil {
		for _, g := range gf {
			if g.GetName() == name {
				return g.GetMetric(), nil
			}
		}
	}

	return nil, err
}

func (m *Metric) CacheLatency(name, action string, n time.Duration) {
	m.cacheLatency.With(prometheus.Labels{"name": name, "action": action}).Observe(n.Seconds())
}
