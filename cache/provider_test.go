package cache_test

import (
	"testing"
	"time"

	"github.com/bukalapak/ottoman/cache"
	"github.com/stretchr/testify/assert"
)

func TestProvider_Name(t *testing.T) {
	r := NewReader()
	c := cache.NewProvider(r)

	assert.Equal(t, r.Name(), c.Name())
}

func TestProvider_Read(t *testing.T) {
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()

	b, err := c.Read("foo")
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"foo":"bar"}`), b)
}

func TestProvider_Read_namespace(t *testing.T) {
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Prefix = "api"
	c.(*cache.Engine).Resolver = NewResolver()

	b, err := c.Read("foo")
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"foo":"bar"}`), b)

	b, err = c.Read("api:foo")
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"foo":"bar"}`), b)
}

func TestProvider_ReadMulti(t *testing.T) {
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()

	keys := []string{
		"fox",
		"api:foo",
	}

	m, err := c.ReadMulti(keys)
	assert.Nil(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, []byte(`{"fox":"baz"}`), m["fox"])
	assert.Equal(t, []byte(`{"foo":"bar"}`), m["foo"])
}

func TestProvider_Fetch(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()

	b, err := c.Fetch("zoo", q)
	assert.Nil(t, err)
	assert.Equal(t, `{"zoo":"zac"}`, string(b))
}

func TestProvider_Fetch_badKey(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()

	b, err := c.Fetch("err", q)
	assert.Equal(t, "unknown cache", err.Error())
	assert.Nil(t, b)
}

func TestProvider_Fetch_failure(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Transport = &FailureTransport{}

	b, err := c.Fetch("zoo", q)
	assert.NotNil(t, err)
	assert.Nil(t, b)
}

func TestProvider_Fetch_timeout(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Timeout = 1 * time.Microsecond

	b, err := c.Fetch("zoo", q)
	assert.Contains(t, err.Error(), "Client.Timeout exceeded")
	assert.Nil(t, b)
}

func TestProvider_Fetch_backendError(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()

	b, err := c.Fetch("bad", q)
	assert.Contains(t, err.Error(), "invalid http status: 500 Internal Server Error")
	assert.Nil(t, b)
}

func TestProvider_FetchMulti(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()

	keys := []string{
		"api:foo",
		"zoo",
	}

	m, err := c.FetchMulti(keys, q)
	assert.Nil(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, []byte(`{"zoo":"zac"}`), m["zoo"])
	assert.Empty(t, m["foo"])
}

func TestProvider_FetchMulti_namespace(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Prefix = "api"
	c.(*cache.Engine).Resolver = NewResolver()

	keys := []string{
		"api:foo",
		"zoo",
	}

	m, err := c.FetchMulti(keys, q)
	assert.Nil(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, []byte(`{"zoo":"zac"}`), m["api:zoo"])
	assert.Empty(t, m["api:foo"])
}

func TestProvider_FetchMulti_failure(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := &XSample{}
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Transport = &FailureTransport{}

	keys := []string{
		"api:foo",
		"zoo",
	}

	m, err := c.FetchMulti(keys, q)
	assert.Nil(t, err)
	assert.Empty(t, m["foo"])
	assert.Empty(t, m["zoo"])
}

func TestProvider_Namespace(t *testing.T) {
	r := NewReader()
	c := cache.NewProvider(r)
	assert.Equal(t, "", c.Namespace())

	c.(*cache.Engine).Prefix = "api"
	assert.Equal(t, "api", c.Namespace())
}

func TestProvider_ReadFetch(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	m := map[string]string{
		"foo": `{"foo":"bar"}`,
		"zoo": `{"zoo":"zac"}`,
	}

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()

	for k, v := range m {
		b, err := c.ReadFetch(k, q)
		assert.Nil(t, err)
		assert.Equal(t, v, string(b))
	}
}

func TestProvider_ReadFetchMulti(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Prefix = "api"
	c.(*cache.Engine).Resolver = NewResolver()

	keys := []string{
		"api:foo",
		"zoo",
	}

	m, err := c.ReadFetchMulti(keys, q)
	assert.Nil(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, []byte(`{"zoo":"zac"}`), m["api:zoo"])
	assert.Equal(t, []byte(`{"foo":"bar"}`), m["api:foo"])
}

func TestProvider_ReadFetchMulti_failure(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := &XSample{}
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Transport = &FailureTransport{}

	keys := []string{
		"api:foo",
		"zoo",
	}

	m, err := c.ReadFetchMulti(keys, q)
	assert.Nil(t, err)
	assert.Empty(t, m["foo"])
	assert.Empty(t, m["zoo"])
}
