package middleware

import (
	"context"
	"net/http"
	"regexp"

	uuid "github.com/satori/go.uuid"
)

var (
	contextKeyRequestID = ContextKey("RequestID")
)

func RequestID(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKeyRequestID, reqID(r))
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKeyRequestID).(string)
	return id
}

func reqID(r *http.Request) string {
	if id := r.Header.Get("X-Request-Id"); id != "" {
		return cleanID(id)
	}

	return genID()
}

func genID() string {
	return uuid.NewV4().String()
}

func cleanID(s string) string {
	re := regexp.MustCompile(`[^\w+\-]`)
	return re.ReplaceAllLiteralString(s, "")
}
