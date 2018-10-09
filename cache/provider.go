package cache

import (
	"time"
)

// Provider wraps Storage interface with additional functionalities.
type Provider interface {
	Storage
	Normalizer
	Namespace() string
}

// NewProvider returns Provider from a Storage and prefix.
func NewProvider(z Storage, prefix string) Provider {
	return &provider{
		engine: z,
		prefix: prefix,
	}
}

type provider struct {
	engine Storage
	prefix string
}

// Name returns cache backend identifier.
func (p *provider) Name() string {
	return p.engine.Name()
}

// Namespace returns cache Prefix
func (p *provider) Namespace() string {
	return p.prefix
}

func (p *provider) Normalize(key string) string {
	return Normalize(key, p.prefix)
}

func (p *provider) NormalizeMulti(keys []string) []string {
	return NormalizeMulti(keys, p.prefix)
}

// Write writes cache data to the cache backend based on key supplied.
func (p *provider) Write(key string, value []byte, expiration time.Duration) error {
	return p.engine.Write(p.Normalize(key), value, expiration)
}

// Read reads cache data on the cache backend based on key supplied.
func (p *provider) Read(key string) ([]byte, error) {
	return p.engine.Read(p.Normalize(key))
}

// ReadMulti bulk reads multiple cache keys.
func (p *provider) ReadMulti(keys []string) (map[string][]byte, error) {
	return p.engine.ReadMulti(p.NormalizeMulti(keys))
}
