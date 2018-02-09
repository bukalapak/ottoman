// Package cache provides unified access to cache backends.
package cache

import (
	"fmt"
	"net/http"
	"strings"
)

// Reader is the interface for cache backend implementation.
type Reader interface {
	Name() string
	Read(key string) ([]byte, error)
	ReadMap(key string) (map[string]interface{}, error)
	ReadMulti(keys []string) (map[string][]byte, error)
}

// Fetcher is the interface for getting cache key from cache engine as well as to remote backend
type Fetcher interface {
	Fetch(key string, r *http.Request) ([]byte, error)
	FetchRequest(r *http.Request) ([]byte, error)
	FetchMap(key string, r *http.Request) (map[string]interface{}, error)
	FetchMulti(keys []string, r *http.Request) (map[string][]byte, error)
}

// Resolver is the interface for resolving cache key to http request and cache router
type Resolver interface {
	Resolve(key string, r *http.Request) (*http.Request, error)
	ResolveRequest(r *http.Request) (*http.Request, error)
	ResolveLatest(key string) string
}

// Normalizer is the interface for normalizing cache key
type Normalizer interface {
	Normalize(key string) string
	NormalizeMulti(keys []string) []string
}

// Provider wraps several interfaces with additional identifier for getting information about the implementation.
type Provider interface {
	Reader
	Fetcher
	Normalizer
	Namespace() string
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
