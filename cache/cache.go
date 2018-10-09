// Package cache provides unified access to cache backends.
package cache

import (
	"fmt"
	"strings"
	"time"
)

// Writer is the interface for cache backend implementation for writing cache data.
type Writer interface {
	Write(key string, value []byte, expiration time.Duration) error
}

// Reader is the interface for cache backend implementation for reading cache data.
type Reader interface {
	Name() string
	Read(key string) ([]byte, error)
	ReadMulti(keys []string) (map[string][]byte, error)
}

// Storage is the interface for writing and reading data.
type Storage interface {
	Writer
	Reader
}

// Normalizer is the interface for normalizing cache key
type Normalizer interface {
	Normalize(key string) string
	NormalizeMulti(keys []string) []string
}

// Normalize returns valid cache key. It can automatically detect prefixed/non-prefixed cache key and format the key properly.
func Normalize(key, prefix string) string {
	if n := strings.SplitN(key, ":", 2); len(n) == 2 {
		key = n[1]
	}

	if prefix != "" {
		return fmt.Sprintf("%s:%s", prefix, key)
	}

	return key
}

// NormalizeMulti is batch version of Normalize.
func NormalizeMulti(keys []string, prefix string) []string {
	ks := make([]string, len(keys))

	for i := range keys {
		ks[i] = Normalize(keys[i], prefix)
	}

	return ks
}
