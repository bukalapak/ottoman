package _qtest

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

type CacheMetric struct {
	cacheLatency *prometheus.HistogramVec
	registry     *prometheus.Registry
}

func NewCacheMetric() *CacheMetric {
	m := &CacheMetric{registry: prometheus.NewRegistry()}

	m.cacheLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "cache_latency_seconds",
		Help: "A histogram of the cache latency in seconds.",
	}, []string{"name", "action"})

	m.registry.MustRegister(m.cacheLatency)

	return m
}

func (m *CacheMetric) Gather(name string) ([]*dto.Metric, error) {
	gf, err := m.registry.Gather()
	if err == nil {
		for _, g := range gf {
			if g.GetName() == name {
				return g.GetMetric(), nil
			}
		}
	}

	return nil, err
}

func (m *CacheMetric) Assert(t *testing.T, name, action string) {
	mc, err := m.Gather("cache_latency_seconds")
	assert.Nil(t, err)
	assert.Len(t, mc, 1)

	labels := map[string]string{
		"name":   name,
		"action": action,
	}

	for _, z := range mc {
		assert.Equal(t, uint64(1), z.GetHistogram().GetSampleCount())
		assert.NotZero(t, z.GetHistogram().GetSampleSum())

		for _, label := range z.GetLabel() {
			assert.Equal(t, labels[label.GetName()], label.GetValue())
		}
	}
}

func (m *CacheMetric) CacheLatency(name, action string, n time.Duration) {
	m.cacheLatency.With(prometheus.Labels{"name": name, "action": action}).Observe(n.Seconds())
}
