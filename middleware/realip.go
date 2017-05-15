package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

var (
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	contextKeyIP  = ContextKey("RealIP")
)

func RealIP(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := NewIPContext(r.Context(), realIP(r))
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func NewIPContext(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, contextKeyIP, ip)
}

func IPFromContext(ctx context.Context) (string, bool) {
	rip, ok := ctx.Value(contextKeyIP).(string)
	return rip, ok
}

func realIP(r *http.Request) string {
	var ip string

	if xfip := getXForwardedFor(r); xfip != "" {
		ip = xfip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if rrip := remoteIP(r); rrip != "" {
		ip = rrip
	}

	return ip
}

func getXForwardedFor(r *http.Request) string {
	var ip string

	xff := r.Header.Get(xForwardedFor)
	if xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}

	return ip
}

func remoteIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
