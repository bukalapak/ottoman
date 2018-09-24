package middleware_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bukalapak/ottoman/logger"
	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
)

func TestLoggerFromContext(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := middleware.NewRequestIDContext(req.Context(), "request-id")

	buf := new(bytes.Buffer)
	log := middleware.LoggerFromContext(ctx, logger.JSON()).Output(buf)
	log.Info().Msg("Hello world!")

	ent := make(map[string]string)
	dec := json.NewDecoder(buf)
	err := dec.Decode(&ent)

	assert.Nil(t, err)
	assert.Equal(t, "request-id", ent["request_id"])
	assert.Equal(t, "info", ent["level"])
	assert.Equal(t, "Hello world!", ent["message"])
}
