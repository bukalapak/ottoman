package logger_test

import (
	"strings"
	"testing"

	"github.com/bukalapak/ottoman/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestJSON(t *testing.T) {
	str := new(strings.Builder)
	log := logger.JSON().WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		str.WriteString(entry.Message)
		return nil
	}))

	log.Info("Hello World!")
	assert.Equal(t, "Hello World!", str.String())
}

func TestDiscard(t *testing.T) {
	str := new(strings.Builder)
	log := logger.Discard().WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
		str.WriteString(entry.Message)
		return nil
	}))

	log.Info("Hello World!")
	assert.Empty(t, str.String())
}
