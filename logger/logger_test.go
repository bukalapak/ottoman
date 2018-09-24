package logger_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/bukalapak/ottoman/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type Entry struct {
	Level     string    `json:"level"`
	Timestamp time.Time `json:"@timestamp"`
	Caller    string    `json:"caller"`
	Message   string    `json:"message"`
}

func TestJSON(t *testing.T) {
	now := time.Now()

	zerolog.TimestampFunc = func() time.Time {
		return now
	}

	buf := new(bytes.Buffer)
	log := logger.JSON().Output(buf)
	log.Info().Msg("Hello world!")

	ent := &Entry{}
	dec := json.NewDecoder(buf)
	err := dec.Decode(ent)

	assert.Nil(t, err)
	assert.Equal(t, "info", ent.Level)
	assert.Equal(t, "Hello world!", ent.Message)
	assert.Contains(t, ent.Caller, "logger_test.go")
	assert.Equal(t, now.Format(time.RFC3339), ent.Timestamp.Format(time.RFC3339))
}

func TestDiscard(t *testing.T) {
	buf := new(bytes.Buffer)
	log := logger.Discard().Output(buf)
	log.Info().Msg("Hello world!")

	assert.Zero(t, buf.Len())
}
