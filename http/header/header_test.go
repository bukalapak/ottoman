package header_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/bukalapak/ottoman/http/header"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	data := map[string]header.ContentType{
		"*/*":                                header.ContentType{MediaType: "*/*", Quality: 1.0, Parameters: map[string]string{}},
		"*/*; q=0.1":                         header.ContentType{MediaType: "*/*", Quality: 0.1, Parameters: map[string]string{"q": "0.1"}},
		"text/plain":                         header.ContentType{MediaType: "text/plain", Quality: 1.0, Parameters: map[string]string{}},
		"application/json":                   header.ContentType{MediaType: "application/json", Quality: 1.0, Parameters: map[string]string{}},
		"application/msgpack":                header.ContentType{MediaType: "application/msgpack", Quality: 1.0, Parameters: map[string]string{}},
		"application/json; charset=utf-8":    header.ContentType{MediaType: "application/json", Quality: 1.0, Parameters: map[string]string{"charset": "utf-8"}},
		"application/msgpack; charset=utf-8": header.ContentType{MediaType: "application/msgpack", Quality: 1.0, Parameters: map[string]string{"charset": "utf-8"}},
		"application/vnd.example.v4+json":    header.ContentType{MediaType: "application/json", Quality: 1.0, Parameters: map[string]string{"version": "v4", "vendor": "example"}},
		"application/vnd.example.v4+msgpack": header.ContentType{MediaType: "application/msgpack", Quality: 1.0, Parameters: map[string]string{"version": "v4", "vendor": "example"}},
		"application/vnd.example.v4+text":    header.ContentType{MediaType: "text/plain", Quality: 1.0, Parameters: map[string]string{"version": "v4", "vendor": "example"}},
	}

	for s, x := range data {
		c := headerContentTypes(s)
		assert.Equal(t, header.ContentTypes{x}, c)
	}
}

func TestHeader_qualityFactor(t *testing.T) {
	q := "application/json; q=0.8, application/msgpack; q=0.4, */*"
	c := headerContentTypes(q)
	s := make([]string, len(c))

	for i, v := range c {
		s[i] = v.MediaType
	}

	x := []string{"application/msgpack", "application/json", "*/*"}
	assert.Equal(t, x, s)
}

func TestHeader_wildcard(t *testing.T) {
	c := headerContentTypes("application/json, */*; q=0.2")
	w, ok := c.Wildcard()
	assert.True(t, ok)
	assert.Equal(t, "*/*", w.MediaType)
	assert.Equal(t, 0.2, w.Quality)

	c = headerContentTypes("application/json")
	w, ok = c.Wildcard()
	assert.False(t, ok)
	assert.Equal(t, "", w.MediaType)
	assert.Equal(t, float64(0), w.Quality)
}

func TestHeader_contentTypes(t *testing.T) {
	data := map[string]string{
		"application/json":                      "application/json",
		"application/json, */*":                 "application/json",
		"application/json, text/html":           "application/json",
		"application/json, application/msgpack": "application/msgpack",
		"application/msgpack":                   "application/msgpack",
		"application/msgpack, */*":              "application/msgpack",
		"application/msgpack, text/html":        "application/msgpack",
		"application/msgpack, application/json": "application/msgpack",
		"application/vnd.example.v4+json":       "application/json",
		"application/vnd.example.v4+msgpack":    "application/msgpack",
		"application/vnd.example.v4+text":       "text/plain",
		"text/plain, text/html":                 "text/plain",
		"*/*": "*/*",
	}

	for k, x := range data {
		m := headerContentTypes(k)
		c := m.ContentType()

		assert.Equal(t, x, c.MediaType)
		assert.Equal(t, 1.0, c.Quality)

		if len(c.Parameters) != 0 {
			assert.Equal(t, c.MediaType+"; vendor=example; version=v4", c.String())
		}

		if strings.Contains(k, "msgpack") {
			assert.True(t, m.Contains("application/msgpack"))
		}
	}

	v := headerContentTypes("")
	assert.Equal(t, "", v.ContentType().MediaType)
	assert.Equal(t, float64(0), v.ContentType().Quality)
}

func TestContentType_VendorString(t *testing.T) {
	data := map[string]string{
		"*/*":                                "*/*",
		"application/json":                   "application/json",
		"application/json; charset=utf-8":    "application/json; charset=utf-8",
		"application/msgpack":                "application/msgpack",
		"application/msgpack; q=0.2":         "application/msgpack; q=0.2",
		"application/vnd.example.v4+json":    "application/vnd.example.v4+json",
		"application/vnd.example.v4+msgpack": "application/vnd.example.v4+msgpack",
	}

	for k, x := range data {
		c := headerContentType(k)
		assert.Equal(t, x, c.VendorString())
	}
}

func TestContentType_String(t *testing.T) {
	data := map[string]string{
		"*/*":                                "*/*",
		"application/json":                   "application/json",
		"application/json; charset=utf-8":    "application/json; charset=utf-8",
		"application/msgpack":                "application/msgpack",
		"application/msgpack; q=0.2":         "application/msgpack; q=0.2",
		"application/vnd.example.v4+json":    "application/json; vendor=example; version=v4",
		"application/vnd.example.v4+msgpack": "application/msgpack; vendor=example; version=v4",
	}

	for k, x := range data {
		c := headerContentType(k)
		assert.Equal(t, x, c.String())
	}
}

func headerContentTypes(s string) header.ContentTypes {
	h := make(http.Header)
	h.Set("Accept", s)
	return header.ContentHeader(h, "Accept")
}

func headerContentType(s string) header.ContentType {
	m := headerContentTypes(s)
	return m.ContentType()
}
