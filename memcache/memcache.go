// Package memcache provides a cache.Storage compatible client for the memcached cache server.
package memcache

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

const (
	defaultTimeout      = 100 * time.Millisecond
	defaultMaxIdleConns = 2
)

// Option represents configurable configuration for memcache client.
type Option struct {
	Compress     bool
	Timeout      time.Duration
	MaxIdleConns int
}

// Memcache is a memcache client. It is safe for unlocked use by multiple concurrent goroutines.
type Memcache struct {
	client *memcache.Client
	option Option
}

// New returns a memcache client using the provided servers and options.
func New(ss []string, option Option) *Memcache {
	c := memcache.New(ss...)
	c.Timeout = netTimeout(option.Timeout)
	c.MaxIdleConns = maxIdleConns(option.MaxIdleConns)

	return &Memcache{client: c, option: option}
}

// Write writes the item for given key.
// It's automatically compress large item.
func (c *Memcache) Write(key string, value []byte, expiration time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      c.compress(value),
		Expiration: int32(expiration.Seconds()),
	}

	return c.client.Set(item)
}

// Read reads the item for given key.
// It's automatically decode item. Value depending on the client option.
func (c *Memcache) Read(key string) ([]byte, error) {
	item, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}

	return c.readValue(item.Value)
}

// ReadMulti is a batch version of Read.
// The returned map have exact length as provided keys. For cache miss, an empty byte will be returned.
func (c *Memcache) ReadMulti(keys []string) (map[string][]byte, error) {
	m, err := c.client.GetMulti(keys)
	if err != nil {
		return map[string][]byte{}, err
	}

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

// MaxIdleConns returns client's cache MaxIdleConns option value.
func (c *Memcache) MaxIdleConns() int {
	return c.client.MaxIdleConns
}

// Delete deletes the item for given key.
func (c *Memcache) Delete(key string) error {
	return c.client.Delete(key)
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

func (c *Memcache) compress(data []byte) []byte {
	if !c.option.Compress {
		return data
	}

	if len(data) >= 1024 {
		var b bytes.Buffer

		z := zlib.NewWriter(&b)
		z.Write(data)
		z.Close()

		return b.Bytes()
	}

	return data
}

func netTimeout(timeout time.Duration) time.Duration {
	if timeout != 0 {
		return timeout
	}

	return defaultTimeout
}

func maxIdleConns(maxIdleConns int) int {
	if maxIdleConns > 0 {
		return maxIdleConns
	}

	return defaultMaxIdleConns
}
