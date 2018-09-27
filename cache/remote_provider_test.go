package cache_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/cache"
	httpclone "github.com/bukalapak/ottoman/http/clone"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRemoteProvider(t *testing.T) {
	h1 := newRemoteServer()
	defer h1.Close()

	c1 := cache.NewProvider(newSample(), "zzz")
	q1 := cache.NewRemoteProvider(c1, cache.RemoteOption{
		Resolver: &resolver{},
	})

	t.Run("Fetch", func(t *testing.T) {
		r, _ := http.NewRequest("GET", h1.URL, nil)
		b, n, err := q1.Fetch("zoo", r)

		assert.Nil(t, err)
		assert.Equal(t, []byte(`{"zoo":"zac"}`), b)
		assert.Equal(t, http.StatusOK, n.StatusCode)
		assert.Contains(t, n.RemoteURL, "/zoo")
	})

	t.Run("Fetch (unknown key)", func(t *testing.T) {
		r, _ := http.NewRequest("GET", h1.URL, nil)
		b, n, err := q1.Fetch("unknown", r)

		assert.NotNil(t, err)
		assert.Nil(t, b)
		assert.Zero(t, n)
	})

	t.Run("Fetch (backend failure)", func(t *testing.T) {
		r, _ := http.NewRequest("GET", h1.URL, nil)
		b, n, err := q1.Fetch("bad", r)

		assert.NotNil(t, err)
		assert.Nil(t, b)
		assert.Equal(t, http.StatusInternalServerError, n.StatusCode)
		assert.Contains(t, n.RemoteURL, "/bad")
	})

	t.Run("Fetch (network failure)", func(t *testing.T) {
		q2 := cache.NewRemoteProvider(c1, cache.RemoteOption{
			Resolver:  &resolver{},
			Transport: &failureTransport{},
			Timeout:   60 * time.Second,
		})

		r, _ := http.NewRequest("GET", h1.URL, nil)
		b, n, err := q2.Fetch("zoo", r)

		assert.NotNil(t, err)
		assert.Zero(t, n)
		assert.Nil(t, b)
	})

	t.Run("FetchMulti", func(t *testing.T) {
		r, _ := http.NewRequest("GET", h1.URL, nil)
		mb, mn, err := q1.FetchMulti([]string{"zoo", "unknown"}, r)

		assert.Contains(t, err.Error(), "zzz:unknown: unknown cache")
		assert.Len(t, mb, 1)
		assert.Equal(t, []byte(`{"zoo":"zac"}`), mb["zzz:zoo"])
		assert.Equal(t, http.StatusOK, mn["zzz:zoo"].StatusCode)
		assert.Contains(t, mn["zzz:zoo"].RemoteURL, "/zoo")
	})
}

type resolver struct{}

func (v *resolver) Resolve(key string, r *http.Request) (*http.Request, error) {
	req, _ := v.ResolveRequest(r)

	keys := map[string]string{
		"zzz:bad": "/bad",
		"zzz:zoo": "/zoo",
	}

	if v, ok := keys[key]; ok {
		req.URL.Path = v
	} else {
		return nil, errors.New("unknown cache")
	}

	return req, nil
}

func (v *resolver) ResolveRequest(r *http.Request) (*http.Request, error) {
	return httpclone.Request(r), nil
}

func newRemoteServer() *httptest.Server {
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/zoo":
			io.WriteString(w, `{"zoo":"zac"}`)
		case "/bad":
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	return httptest.NewServer(http.HandlerFunc(fn))
}

type failureTransport struct{}

func (t *failureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, errors.New("Connection failure")
}
