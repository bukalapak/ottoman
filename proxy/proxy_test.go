package proxy_test

import (
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
	"github.com/stretchr/testify/suite"
)

type ProxySuite struct {
	suite.Suite
	backend *httptest.Server
}

func (suite *ProxySuite) SetupSuite() {
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

	suite.backend = httptest.NewServer(http.HandlerFunc(fn))
}

func (suite *ProxySuite) TearDownSuite() {
	suite.backend.Close()
}

func (suite *ProxySuite) Targeter() proxy.Targeter {
	u, _ := url.Parse(suite.backend.URL)
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

func (suite *ProxySuite) TestProxy() {
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	x := proxy.NewProxy(suite.Targeter())
	x.Forward(rec, req, Transform{})

	assert.Equal(suite.T(), suite.backend.URL, x.Target().String())
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(suite.T(), `{"foo":"bar"}`, strings.TrimSpace(rec.Body.String()))
	assert.Equal(suite.T(), "1", rec.Header().Get("X-Modified"))
}

func (suite *ProxySuite) TestProxy_chunked() {
	req, _ := http.NewRequest("GET", "/stream", nil)
	rec := httptest.NewRecorder()

	x := proxy.NewProxy(suite.Targeter())
	x.FlushInterval = time.Millisecond
	x.Forward(rec, req, Transform{})

	assert.Equal(suite.T(), suite.backend.URL, x.Target().String())
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.True(suite.T(), rec.Flushed)

	for i := 1; i <= 10; i++ {
		assert.Contains(suite.T(), rec.Body.String(), fmt.Sprintf("chunk #%d", i))
	}
}

func TestProxySuite(t *testing.T) {
	suite.Run(t, new(ProxySuite))
}
