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
	defaultMaxAttempt   = 3
)

// Option represents configurable configuration for memcache client.
type Option struct {
	Compress     bool
	Timeout      time.Duration
	MaxIdleConns int
	MaxAttempt   int
}

type MemcacheClient interface {
	Set(*memcache.Item) error
	Get(string) (*memcache.Item, error)
	GetMulti([]string) (map[string]*memcache.Item, error)
	Delete(string) error
}

// Memcache is a memcache client. It is safe for unlocked use by multiple concurrent goroutines.
type Memcache struct {
	client MemcacheClient
	option Option
}

// New returns a memcache client using the provided servers and options.
func New(ss []string, option Option) *Memcache {
	option = optionDefaultValue(option)

	c := memcache.New(ss...)
	c.Timeout = option.Timeout
	c.MaxIdleConns = option.MaxIdleConns

	return &Memcache{client: c, option: option}
}

func NewWithClient(mc MemcacheClient, option Option) *Memcache {
	option = optionDefaultValue(option)

	return &Memcache{client: mc, option: option}
}

// Write writes the item for given key.
// It's automatically compress large item.
func (c *Memcache) Write(key string, value []byte, expiration time.Duration) error {
	item := &memcache.Item{
		Key:        key,
		Value:      c.compress(value),
		Expiration: int32(expiration.Seconds()),
	}
	fn := func() error {
		return c.client.Set(item)
	}

	err := c.withRetryOnTimeout(fn)
	return err
}

// Read reads the item for given key.
// It's automatically decode item. Value depending on the client option.
func (c *Memcache) Read(key string) ([]byte, error) {
	var item *memcache.Item
	var err error

	fn := func() error {
		item, err = c.client.Get(key)
		return err
	}

	err = c.withRetryOnTimeout(fn)

	if err != nil {
		return nil, err
	}

	return c.readValue(item.Value)
}

// ReadMulti is a batch version of Read.
// The returned map have exact length as provided keys. For cache miss, an empty byte will be returned.
func (c *Memcache) ReadMulti(keys []string) (map[string][]byte, error) {
	var m map[string]*memcache.Item
	var err error

	fn := func() error {
		m, err = c.client.GetMulti(keys)
		return err
	}

	err = c.withRetryOnTimeout(fn)

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
	return c.option.MaxIdleConns
}

// Delete deletes the item for given key.
func (c *Memcache) Delete(key string) error {
	var err error

	fn := func() error {
		return c.client.Delete(key)
	}

	err = c.withRetryOnTimeout(fn)

	return err
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

func (c *Memcache) withRetryOnTimeout(fn func() error) error {
	var err error

	for i := 0; i < c.option.MaxAttempt; i++ {
		err = fn()
		if !timeoutError(err) {
			return err
		}
	}

	return err
}

func optionDefaultValue(option Option) Option {
	option.Timeout = netTimeout(option.Timeout)
	option.MaxIdleConns = maxIdleConns(option.MaxIdleConns)
	option.MaxAttempt = maxAttempt(option.MaxAttempt)

	return option
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

func maxAttempt(maxAttempt int) int {
	if maxAttempt > 0 {
		return maxAttempt
	}

	return defaultMaxAttempt
}

func timeoutError(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(*memcache.ConnectTimeoutError)
	return ok
}
