package middleware

import (
	"net/http"
	"runtime/debug"
)

type Notifier interface {
	Notify(err interface{}, stack []byte)
}

type notifyLogger interface {
	Error(msg string, stackTrace string)
}

type Recovery struct {
	Logger notifyLogger
	agent  Notifier
}

func NewRecovery(agent Notifier) *Recovery {
	return &Recovery{agent: agent, Logger: &nopLogger{}}
}

func (v *Recovery) Handler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				v.Logger.Error("ottoman:middleware/recovery", string(debug.Stack()))
				v.agent.Notify(rec, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type nopLogger struct{}

func (n *nopLogger) Error(msg string, stackTrace string) {}
