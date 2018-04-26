package middleware_test

import (
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
)

func TestContextKey(t *testing.T) {
	key := middleware.ContextKey("Hello")
	assert.Equal(t, "ottoman/middleware context: Hello", key.String())
}
