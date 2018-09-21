package middleware

import (
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

type Notifier interface {
	Notify(err interface{}, stack []byte)
}

type Recovery struct {
	Logger *zap.Logger
	agent  Notifier
}

func NewRecovery(agent Notifier) *Recovery {
	return &Recovery{agent: agent}
}

func (v *Recovery) Handler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log := LoggerFromContext(r.Context(), v.Logger)
				log.Error("ottoman:middleware/recovery",
					zap.String("stack_trace", string(debug.Stack())),
				)

				v.agent.Notify(rec, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
