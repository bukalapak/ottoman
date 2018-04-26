package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := middleware.LoggerFromContext(r.Context())

		assert.False(t, log.Core().Enabled(zapcore.DebugLevel))
		assert.True(t, log.Core().Enabled(zapcore.InfoLevel))
	}

	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	log := middleware.NewLogger(middleware.JSONLogger())
	log.Handler(http.HandlerFunc(fn)).ServeHTTP(rec, req)
}

func TestLoggerFromContext(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	log := middleware.LoggerFromContext(req.Context())
	assert.Equal(t, zap.New(nil), log)
}
