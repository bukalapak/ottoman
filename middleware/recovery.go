package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
)

type Notifier interface {
	Notify(err interface{}, stack []byte)
}

type Recovery struct {
	Logger zerolog.Logger
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
				log.Error().
					Str("stack_trace", string(debug.Stack())).
					Msg("ottoman:middleware/recovery")

				v.agent.Notify(rec, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
