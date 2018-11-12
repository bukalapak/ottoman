// Package redis implements a Redis client which compatible with cache.Storage.
package redis

import (
	"errors"
	"time"

	redisc "github.com/go-redis/redis"
)

var errCacheMiss = errors.New("redis: cache miss")

// Option represents configurable configuration for redis client.
type Option struct {
	Addrs    []string
	Password string

	// A database to be selected after connecting to server.
	// Redis Cluster ignores this value.
	DB int

	// Cluster specific flag to enable read-only commands on slave nodes.
	ReadOnly bool
}

type connector interface {
	Set(key string, value interface{}, expiration time.Duration) *redisc.StatusCmd
	Get(key string) *redisc.StringCmd
	MGet(keys ...string) *redisc.SliceCmd
	Incr(key string) *redisc.IntCmd
	Expire(key string, expiration time.Duration) *redisc.BoolCmd
	Del(keys ...string) *redisc.IntCmd
}

// Redis is a Redis client representing a pool of zero or more underlying connections.
// It's saafe for concurrent use by multiple goroutines.
type Redis struct {
	client connector
	name   string
}

// New returns a client to the redis server specified by Option.
func New(opts *Option) *Redis {
	if len(opts.Addrs) == 1 {
		return &Redis{
			name: "Redis",
			client: redisc.NewClient(&redisc.Options{
				Addr:     opts.Addrs[0],
				DB:       opts.DB,
				Password: opts.Password,
			}),
		}
	}

	return &Redis{
		name: "Redis Cluster",
		client: redisc.NewClusterClient(&redisc.ClusterOptions{
			Addrs:    opts.Addrs,
			Password: opts.Password,
			ReadOnly: opts.ReadOnly,
		}),
	}
}

// Write writes the item for given key.
func (c *Redis) Write(key string, value []byte, expiration time.Duration) error {
	cmd := c.client.Set(key, value, expiration)

	return cmd.Err()
}

// Read reads the item for given key.
func (c *Redis) Read(key string) ([]byte, error) {
	cmd := c.client.Get(key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return cmd.Bytes()
}

// ReadMulti is a batch version of Read.
// The returned map have exact length as provided keys. For cache miss, an empty byte will be returned.
func (c *Redis) ReadMulti(keys []string) (map[string][]byte, error) {
	cmd := c.client.MGet(keys...)
	z := make(map[string][]byte)

	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

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
	return c.name
}

// Delete deletes the item for given key.
func (c *Redis) Delete(key string) error {
	cmd := c.client.Del(key)

	n, _ := cmd.Result()
	if n == 0 {
		return errCacheMiss
	}

	return nil
}
