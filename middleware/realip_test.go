package middleware_test

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRealIP(t *testing.T) {
	sampleIP := net.ParseIP("202.212.212.202")

	fn := func(w http.ResponseWriter, r *http.Request) {
		if ip, ok := middleware.IPFromContext(r.Context()); ok {
			io.WriteString(w, ip)
		}
	}

	x := func(req *http.Request, ip string) {
		rec := httptest.NewRecorder()

		middleware.RealIP(http.HandlerFunc(fn)).ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, ip, rec.Body.String())
	}

	t.Run("RealIP", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("X-Forwarded-For", sampleIP.String())
		req.Header.Add("X-Real-IP", "222.222.222.222")

		x(req, sampleIP.String())
	})

	t.Run("X-Forwarded-For", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("X-Forwarded-For", sampleIP.String())

		x(req, sampleIP.String())
	})

	t.Run("X-Real-IP", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("X-Real-IP", sampleIP.String())

		x(req, sampleIP.String())
	})

	t.Run("RemoteAddr", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		host, _, _ := net.SplitHostPort(req.RemoteAddr)

		x(req, host)
	})
}
