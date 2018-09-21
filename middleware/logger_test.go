package middleware_test

import (
	"net/http"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerFromContext(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := middleware.NewRequestIDContext(req.Context(), "request-id")

	core, recorded := observer.New(zapcore.InfoLevel)
	log1 := zap.New(core)
	log2 := middleware.LoggerFromContext(ctx, log1)
	log2.Info("Hello world!")

	assert.Equal(t, 1, recorded.Len())

	entry := recorded.All()[0]
	assert.Equal(t, "Hello world!", entry.Message)
	assert.Equal(t, "request-id", entry.ContextMap()["request_id"])
}
