package proxy_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/proxy"
	"github.com/stretchr/testify/assert"
)

func TestProxy(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/stream" {
			if flusher, ok := w.(http.Flusher); ok {
				for i := 1; i <= 10; i++ {
					fmt.Fprintf(w, "chunk #%d\n", i)
					flusher.Flush()
					time.Sleep(10 * time.Millisecond)
				}

				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"foo":"bar"}`)
	}

	backend := httptest.NewServer(http.HandlerFunc(fn))
	defer backend.Close()

	t.Run("Proxy", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		x := proxy.NewProxy(targeter(backend))
		x.Forward(rec, req, Transform{})

		assert.Equal(t, backend.URL, x.Target().String())
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.Equal(t, `{"foo":"bar"}`, strings.TrimSpace(rec.Body.String()))
		assert.Equal(t, "1", rec.Header().Get("X-Modified"))
	})

	t.Run("Proxy-Chunked", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/stream", nil)
		rec := httptest.NewRecorder()

		x := proxy.NewProxy(targeter(backend))
		x.FlushInterval = time.Millisecond
		x.Forward(rec, req, Transform{})

		assert.Equal(t, backend.URL, x.Target().String())
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, rec.Flushed)

		for i := 1; i <= 10; i++ {
			assert.Contains(t, rec.Body.String(), fmt.Sprintf("chunk #%d", i))
		}
	})

	t.Run("Proxy-Bad-Gateway", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		x := proxy.NewProxy(targeter(backend))
		x.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
			rw.Header().Set("X-Fail", "1")
			rw.WriteHeader(http.StatusBadGateway)
		}

		x.Forward(rec, req, FailingTransform{})

		assert.Equal(t, backend.URL, x.Target().String())
		assert.Equal(t, http.StatusBadGateway, rec.Code)
		assert.Equal(t, "1", rec.Header().Get("X-Fail"))
	})
}

func targeter(backend *httptest.Server) proxy.Targeter {
	u, _ := url.Parse(backend.URL)
	t := proxy.NewTarget(u)

	return t
}

type Transform struct{}

func (s Transform) Director(t proxy.Targeter) func(r *http.Request) {
	return func(r *http.Request) {
		u := t.Target()
		r.URL.Scheme = u.Scheme
		r.URL.Host = u.Host
	}
}

func (s Transform) RoundTrip(r *http.Request) (*http.Response, error) {
	return http.DefaultTransport.RoundTrip(r)
}

func (s Transform) ModifyResponse(resp *http.Response) error {
	resp.Header.Set("X-Modified", "1")
	return nil
}

type FailingTransform struct{}

func (s FailingTransform) Director(t proxy.Targeter) func(r *http.Request) {
	return func(r *http.Request) {
		u := t.Target()
		r.URL.Scheme = u.Scheme
		r.URL.Host = u.Host
	}
}

func (s FailingTransform) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("some error")
}

func (s FailingTransform) ModifyResponse(resp *http.Response) error {
	return nil
}
