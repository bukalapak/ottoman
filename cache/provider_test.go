package cache_test

import (
	"fmt"
	"sort"
	"sync"
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
	n := NewCounter()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Counter = n

	b, err := c.Read("foo")
	assert.Nil(t, err)
	assert.Equal(t, []byte(`{"foo":"bar"}`), b)
	assert.Equal(t, 1, n.C)
	assert.Equal(t, 0, n.B)
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
	n := NewCounter()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Counter = n

	keys := []string{
		"fox",
		"api:foo",
	}

	m, err := c.ReadMulti(keys)
	assert.Nil(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, []byte(`{"fox":"baz"}`), m["fox"])
	assert.Equal(t, []byte(`{"foo":"bar"}`), m["foo"])
	assert.Equal(t, 1, n.C)
	assert.Equal(t, 0, n.B)
}

func TestProvider_Fetch(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	n := NewCounter()
	z := NewTracer()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Counter = n
	c.(*cache.Engine).Tracer = z

	b, err := c.Fetch("zoo", q)
	assert.Nil(t, err)
	assert.Equal(t, `{"zoo":"zac"}`, string(b))
	assert.Equal(t, 0, n.C)
	assert.Equal(t, 1, n.B)
	assert.Equal(t, 1, len(z.M))
	assert.Contains(t, z.M[0], "/zoo::200")
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
	n := NewCounter()
	z := NewTracer()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Counter = n
	c.(*cache.Engine).Tracer = z

	keys := []string{
		"api:foo",
		"zoo",
	}

	m, err := c.FetchMulti(keys, q)
	assert.Nil(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, []byte(`{"zoo":"zac"}`), m["zoo"])
	assert.Empty(t, m["foo"])
	assert.Equal(t, 0, n.C)
	assert.Equal(t, 1, n.B)
	assert.Equal(t, 1, len(z.M))
	assert.Contains(t, z.M[0], "/zoo::200")
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
	n := NewCounter()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Counter = n

	for k, v := range m {
		b, err := c.ReadFetch(k, q)
		assert.Nil(t, err)
		assert.Equal(t, v, string(b))
	}

	assert.Equal(t, 1, n.C)
	assert.Equal(t, 1, n.B)
}

func TestProvider_ReadFetchMulti(t *testing.T) {
	h := NewRemoteServer()
	defer h.Close()

	q := NewRequest(h.URL)
	r := NewReader()
	n := NewCounter()
	c := cache.NewProvider(r)
	c.(*cache.Engine).Prefix = "api"
	c.(*cache.Engine).Resolver = NewResolver()
	c.(*cache.Engine).Counter = n

	keys := []string{
		"api:foo",
		"zoo",
	}

	m, err := c.ReadFetchMulti(keys, q)
	assert.Nil(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, []byte(`{"zoo":"zac"}`), m["api:zoo"])
	assert.Equal(t, []byte(`{"foo":"bar"}`), m["api:foo"])
	assert.Equal(t, 1, n.C)
	assert.Equal(t, 1, n.B)
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

type simpleCounter struct {
	mu sync.Mutex
	C  int
	B  int
}

func NewCounter() *simpleCounter {
	return &simpleCounter{}
}

func (c *simpleCounter) IncrCacheCounter() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.C++
}

func (c *simpleCounter) IncrBackendCounter() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.B++
}

type simpleTracer struct {
	mu sync.Mutex
	M  []string
}

func NewTracer() *simpleTracer {
	return &simpleTracer{}
}

func (c *simpleTracer) BackendLatency(route string, code int, n time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.M = append(c.M, fmt.Sprintf("%s::%d::%s", route, code, n.String()))
	sort.Strings(c.M)
}
