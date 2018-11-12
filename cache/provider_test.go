package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/cache"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	z1 := newSample()
	c1 := cache.NewProvider(z1, "zzz")
	c2 := cache.NewProvider(newBroken(), "")

	t.Run("Name", func(t *testing.T) {
		assert.Equal(t, z1.Name(), c1.Name())
	})

	t.Run("Namespace", func(t *testing.T) {
		assert.Equal(t, "zzz", c1.Namespace())
		assert.Equal(t, "", c2.Namespace())
	})

	t.Run("Write", func(t *testing.T) {
		err := c1.Write("foo", []byte("bar"), 10*time.Second)
		assert.Nil(t, err)

		v, ok := z1.written["zzz:foo"]
		assert.True(t, ok)
		assert.Equal(t, map[string]string{"10s": "bar"}, v)
	})

	t.Run("Write (failure)", func(t *testing.T) {
		err := c2.Write("foo", []byte("bar"), 10*time.Second)
		assert.NotNil(t, err)
	})

	t.Run("Read", func(t *testing.T) {
		b, err := c1.Read("foo")
		assert.Nil(t, err)
		assert.Equal(t, []byte(`{"zzz":"bar"}`), b)
	})

	t.Run("Read (failure)", func(t *testing.T) {
		_, err := c2.Read("foo")
		assert.NotNil(t, err)
	})

	t.Run("ReadMulti", func(t *testing.T) {
		keys := []string{
			"foo",
			"boo",
		}

		mb, err := c1.ReadMulti(keys)
		assert.Nil(t, err)
		assert.Equal(t, []byte(`{"zzz":"bar"}`), mb["zzz:foo"])
		assert.Equal(t, []byte(`{"zzz":"baz"}`), mb["zzz:boo"])
	})

	t.Run("ReadMulti (failure)", func(t *testing.T) {
		keys := []string{
			"foo",
			"boo",
		}

		_, err := c2.ReadMulti(keys)
		assert.NotNil(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := c1.Delete("foo")
		assert.Nil(t, err)
	})

	t.Run("Delete (failure)", func(t *testing.T) {
		err := c2.Delete("foo")
		assert.NotNil(t, err)
	})
}

type sample struct {
	data    map[string]string
	written map[string]map[string]string
}

func (m *sample) Name() string {
	return "cache/sample"
}

func (m *sample) Write(key string, value []byte, expiration time.Duration) error {
	m.written[key] = map[string]string{
		expiration.String(): string(value),
	}

	return nil
}

func (m *sample) Read(key string) ([]byte, error) {
	if v, ok := m.data[key]; ok {
		return []byte(v), nil
	}

	return nil, errors.New("unknown cache")
}

func (m *sample) ReadMulti(keys []string) (map[string][]byte, error) {
	z := make(map[string][]byte, len(keys))

	for _, key := range keys {
		v, _ := m.Read(key)
		z[key] = []byte(v)
	}

	return z, nil
}

func (m *sample) Delete(key string) error {
	delete(m.data, key)
	return nil
}

func newSample() *sample {
	return &sample{
		written: make(map[string]map[string]string),
		data: map[string]string{
			"foo":     `{"foo":"bar"}`,
			"yyy:foo": `{"yyy":"bar"}`,
			"zzz:foo": `{"zzz":"bar"}`,
			"zzz:boo": `{"zzz":"baz"}`,
		},
	}
}

type broken struct{}

func (*broken) Write(key string, value []byte, expiration time.Duration) error {
	return errors.New("example error from Write")
}

func (*broken) Read(key string) ([]byte, error) {
	return nil, errors.New("example error from Read")
}

func (*broken) ReadMulti(keys []string) (map[string][]byte, error) {
	return nil, errors.New("example error from ReadMulti")
}

func (*broken) Delete(key string) error {
	return errors.New("example error from Delete")
}

func (*broken) Name() string {
	return "cache/broken"
}

func newBroken() *broken {
	return &broken{}
}
