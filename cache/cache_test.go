package cache_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bukalapak/ottoman/cache"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	data := []string{
		"foo/bar",
		"api:foo/bar",
		"bar:foo/bar",
	}

	for i := range data {
		key := cache.Normalize(data[i], "")
		assert.Equal(t, "foo/bar", key)

		key = cache.Normalize(data[i], "api")
		assert.Equal(t, "api:foo/bar", key)

		key = cache.Normalize(data[i], "foo")
		assert.Equal(t, "foo:foo/bar", key)
	}
}

type Sample struct {
	data map[string]string
}

func NewReader() cache.Reader {
	return &Sample{data: map[string]string{
		"foo":     `{"foo":"bar"}`,
		"fox":     `{"fox":"baz"}`,
		"api:foo": `{"foo":"bar"}`,
		"baz":     `x`,
	}}
}

func (m *Sample) Name() string {
	return "cache/reader"
}

func (m *Sample) Read(key string) ([]byte, error) {
	if v, ok := m.data[key]; ok {
		return []byte(v), nil
	}

	return nil, errors.New("unknown cache")
}

func (m *Sample) ReadMap(key string) (map[string]interface{}, error) {
	b, err := m.Read(key)
	if err != nil {
		return nil, err
	}

	z := make(map[string]interface{})
	err = json.Unmarshal(b, &z)

	return z, err
}

func (m *Sample) ReadMulti(keys []string) (map[string][]byte, error) {
	z := make(map[string][]byte, len(keys))

	for _, key := range keys {
		v, _ := m.Read(key)
		z[key] = []byte(v)
	}

	return z, nil
}

type XSample struct{}

func (m *XSample) Read(key string) ([]byte, error) {
	return nil, errors.New("example error from Read")
}

func (m *XSample) ReadMap(key string) (map[string]interface{}, error) {
	return nil, errors.New("example error from ReadMap")
}

func (m *XSample) ReadMulti(keys []string) (map[string][]byte, error) {
	return nil, errors.New("example error from ReadMulti")
}

func (m *XSample) Name() string {
	return "cache/x-reader"
}

func NewRequest(s string) *http.Request {
	r, _ := http.NewRequest("GET", s, nil)
	return r
}

type Match struct{}

func NewResolver() cache.Resolver {
	return &Match{}
}

type FailureTransport struct{}

func (t *FailureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, errors.New("Connection failure")
}

func (m *Match) Resolve(key string, r *http.Request) *http.Request {
	req := new(http.Request)
	req.URL = new(url.URL)
	*req = *r
	*req.URL = *r.URL

	switch key {
	case "zoo", "bad":
		req.URL.Path = "/" + key
	case "api:zoo":
		req.URL.Path = "/zoo"
	}

	return req
}

func NewRemoteServer() *httptest.Server {
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/zoo":
			io.WriteString(w, `{"zoo":"zac"}`)
		case "/zab":
			io.WriteString(w, `remote-x`)
		case "/bad":
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	return httptest.NewServer(http.HandlerFunc(fn))
}
