package middleware

import (
	"context"

	"github.com/rs/zerolog"
)

func LoggerWithContext(ctx context.Context, log zerolog.Logger) zerolog.Logger {
	return log.With().Str("request_id", RequestIDFromContext(ctx)).Logger()
}
