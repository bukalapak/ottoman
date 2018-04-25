package qtest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/bukalapak/ottoman/cache"
)

type CacheMap struct {
	keymap map[string]string
}

func NewCacheReader(keymap map[string]string) cache.Reader {
	return &CacheMap{keymap: keymap}
}

func (q *CacheMap) Name() string {
	return "qtest/cache:reader"
}

func (q *CacheMap) Read(key string) ([]byte, error) {
	if fname, ok := q.keymap[key]; ok {
		return ioutil.ReadFile(fname)
	}

	return nil, fmt.Errorf("unknown cache key: %s", key)
}

func (q *CacheMap) ReadMulti(keys []string) (map[string][]byte, error) {
	z := make(map[string][]byte, len(keys))

	for _, key := range keys {
		if b, err := q.Read(key); err == nil {
			z[key] = b
		} else {
			z[key] = nil
		}
	}

	return z, nil
}

func NewCacheBackend(reqmap map[string]string) *httptest.Server {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if fixture, ok := reqmap[r.URL.RequestURI()]; ok {
			if b, err := ioutil.ReadFile(fixture); err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write(b)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	}

	return httptest.NewServer(http.HandlerFunc(fn))
}
