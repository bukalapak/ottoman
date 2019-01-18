package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bukalapak/ottoman/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		panic("!!!")
	}

	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	cov := middleware.NewRecovery(NewAgent(t))
	cov.Logger = NewLogger(t)

	cov.Handler(http.HandlerFunc(fn)).ServeHTTP(rec, req)
}

type Agent struct {
	t *testing.T
}

func NewAgent(t *testing.T) *Agent {
	return &Agent{t: t}
}

func (a *Agent) Notify(err interface{}, stack []byte) {
	assert.Equal(a.t, "!!!", err)
	assert.NotEmpty(a.t, stack)
}

type Logger struct {
	t *testing.T
}

func NewLogger(t *testing.T) *Logger {
	return &Logger{t: t}
}

func (l *Logger) Error(msg string, stack string) {
	assert.Equal(l.t, "!!!", msg)
	assert.NotEmpty(l.t, stack)
}
