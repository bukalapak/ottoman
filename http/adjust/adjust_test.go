package adjust_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bukalapak/ottoman/http/adjust"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	r := httptest.NewRequest("GET", "https://example.com/foo", nil)
	z := []adjust.Requester{
		sampleHeader,
		func(r *http.Request) error { return sample(r, "https://example.net/bar") },
	}

	err := adjust.Request(r, z...)

	assert.Nil(t, err)
	assert.Equal(t, "world", r.Header.Get("X-Hello"))
	assert.Equal(t, "example.com", r.URL.Host)
	assert.Equal(t, "/bar", r.URL.Path)
}

func TestRequest_error(t *testing.T) {
	r := httptest.NewRequest("GET", "https://example.com/foo", nil)
	z := []adjust.Requester{
		sampleHeader,
		func(r *http.Request) error { return sample(r, "https://exa- -mple.net/bar") },
	}

	err := adjust.Request(r, z...)
	assert.NotNil(t, err)
	assert.Equal(t, "world", r.Header.Get("X-Hello"))
	assert.Equal(t, "example.com", r.URL.Host)
	assert.Equal(t, "/foo", r.URL.Path)
}

func sample(r *http.Request, s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	r.URL.Path = u.Path
	return nil
}

func sampleHeader(r *http.Request) error {
	r.Header.Set("X-Hello", "world")
	return nil
}
