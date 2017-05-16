// Package redis implements a Redis client.
package redis

import (
	"time"

	"github.com/bukalapak/ottoman/encoding/json"
	redisc "gopkg.in/redis.v3"
)

// Redis is a Redis client representing a pool of zero or more underlying connections.
// It's saafe for concurrent use by multiple goroutines.
type Redis struct {
	client *redisc.Client
}

// New returns a client to the redis server specified by redisc.Options.
func New(opts *redisc.Options) *Redis {
	return &Redis{
		client: redisc.NewClient(opts),
	}
}

// Read reads the item for given key.
func (c *Redis) Read(key string) ([]byte, error) {
	cmd := c.client.Get(key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return cmd.Bytes()
}

// ReadMap reads the item for given key as map[string]interface{}
func (c *Redis) ReadMap(key string) (map[string]interface{}, error) {
	b, err := c.Read(key)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// ReadMulti is a batch version of Read.
// The returned map have exact length as provided keys. For cache miss, an empty byte will be returned.
func (c *Redis) ReadMulti(keys []string) (map[string][]byte, error) {
	cmd := c.client.MGet(keys...)
	z := make(map[string][]byte)

	for i, k := range keys {
		v, ok := cmd.Val()[i].(string)
		if !ok {
			continue
		}

		z[k] = []byte(v)
	}

	return z, nil
}

// Incr increase counter for given key.
func (c *Redis) Incr(key string) (int64, error) {
	cmd := c.client.Incr(key)
	n, err := cmd.Result()

	return n, err
}

// Expire sets time expiration on a key
func (c *Redis) Expire(key string, expiration time.Duration) (bool, error) {
	cmd := c.client.Expire(key, expiration)
	b, err := cmd.Result()

	return b, err
}

// Name returns cache storage identifier.
func (c *Redis) Name() string {
	return "Redis"
}
