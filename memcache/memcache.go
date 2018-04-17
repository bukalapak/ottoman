// Package memcache provides a cache.Provider compatible client for the memcached cache server.
package memcache

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bukalapak/ottoman/cache"
	"github.com/bukalapak/ottoman/encoding/json"
)

const (
	defaultTimeout = 100 * time.Millisecond
)

// Option represents configurable configuration for memcache client.
type Option struct {
	Compress bool
	Timeout  time.Duration
	Metric   cache.MetricTracer
}

// Memcache is a memcache client. It is safe for unlocked use by multiple concurrent goroutines.
type Memcache struct {
	client *memcache.Client
	metric cache.MetricTracer
	option Option
}

// New returns a memcache client using the provided servers and options.
func New(ss []string, option Option) *Memcache {
	var m cache.MetricTracer

	if option.Metric != nil {
		m = option.Metric
	} else {
		m = &noopTracer{}
	}

	c := memcache.New(ss...)
	c.Timeout = netTimeout(option.Timeout)

	return &Memcache{client: c, metric: m, option: option}
}

// Read reads the item for given key.
// It's automatically decode item.Value depending on the client option.
func (c *Memcache) Read(key string) ([]byte, error) {
	now := time.Now()
	item, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}

	c.metric.CacheLatency(c.Name(), "Get", time.Since(now))

	return c.readValue(item.Value)
}

// ReadMap reads the item for given key as map[string]interface{}
func (c *Memcache) ReadMap(key string) (map[string]interface{}, error) {
	b, err := c.Read(key)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})

	err = c.decode(bytes.NewReader(b), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// ReadMulti is a batch version of Read.
// The returned map have exact length as provided keys. For cache miss, an empty byte will be returned.
func (c *Memcache) ReadMulti(keys []string) (map[string][]byte, error) {
	now := time.Now()
	m, err := c.client.GetMulti(keys)
	if err != nil {
		return map[string][]byte{}, err
	}

	c.metric.CacheLatency(c.Name(), "GetMulti", time.Since(now))

	z := make(map[string][]byte, len(m))

	for k, v := range m {
		b, _ := c.readValue(v.Value)
		z[k] = b
	}

	return z, nil
}

// Name returns cache storage identifier.
func (c *Memcache) Name() string {
	return "Memcached"
}

func (c *Memcache) decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (c *Memcache) readValue(data []byte) (n []byte, err error) {
	if !c.option.Compress {
		return data, nil
	}

	b := bytes.NewBuffer(data)

	r, err := zlib.NewReader(b)
	if err != nil {
		return data, nil
	}

	return ioutil.ReadAll(r)
}

func netTimeout(timeout time.Duration) time.Duration {
	if timeout != 0 {
		return timeout
	}

	return defaultTimeout
}

type noopTracer struct{}

func (c *noopTracer) CacheLatency(name, action string, n time.Duration) {}
