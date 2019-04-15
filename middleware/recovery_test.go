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
	rec1 := httptest.NewRecorder()
	cov1 := middleware.NewRecovery(NewAgent(t))
	cov1.Handler(http.HandlerFunc(fn)).ServeHTTP(rec1, req)

	rec2 := httptest.NewRecorder()
	cov2 := middleware.NewRecovery(NewAgent(t))
	cov2.Logger = NewLogger(t)

	cov2.Handler(http.HandlerFunc(fn)).ServeHTTP(rec2, req)
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

func (l *Logger) Error(err interface{}, stack []byte) {
	assert.Equal(l.t, "!!!", err)
	assert.NotEmpty(l.t, stack)
}
