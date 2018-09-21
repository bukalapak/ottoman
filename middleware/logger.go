package middleware

import (
	"context"

	"github.com/bukalapak/ottoman/logger"
	"go.uber.org/zap"
)

func LoggerFromContext(ctx context.Context, log *zap.Logger) *zap.Logger {
	if log == nil {
		log = logger.Discard()
	}

	return log.With(
		zap.String("request_id", RequestIDFromContext(ctx)),
	)
}
