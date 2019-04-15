package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/bukalapak/ottoman/notify"
)

type RecoveryLogger interface {
	Error(err interface{}, stackTrace []byte)
}

type Recovery struct {
	Logger RecoveryLogger
	agent  notify.Notifier
}

func NewRecovery(agent notify.Notifier) *Recovery {
	return &Recovery{agent: agent, Logger: &nopLogger{}}
}

func (v *Recovery) Handler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				v.Logger.Error(rec, debug.Stack())
				v.agent.Notify(rec, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type nopLogger struct{}

func (n *nopLogger) Error(err interface{}, stackTrace []byte) {}
