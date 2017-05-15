package middleware

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	contextKeyLogger = ContextKey("Logger")
)

type Logger struct {
	log *zap.Logger
}

func NewLogger(log *zap.Logger) *Logger {
	return &Logger{log: log}
}

func (l *Logger) Handler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := NewLoggerContext(r.Context(), l.log)
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func NewLoggerContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, log)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	var logger *zap.Logger

	if log, ok := ctx.Value(contextKeyLogger).(*zap.Logger); ok {
		logger = log
	} else {
		logger = nullLogger()
	}

	return logger.With(
		zap.String("request_id", RequestIDFromContext(ctx)),
	)
}

func JSONLogger() *zap.Logger {
	log, _ := logConfig().Build()
	return log
}

func nullLogger() *zap.Logger {
	return zap.New(nil)
}

func logConfig() zap.Config {
	n := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zap.Config{
		Level:       zap.NewAtomicLevel(),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    n,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}
