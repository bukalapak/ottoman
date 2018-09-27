package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimestampFieldName = "@timestamp"
}

func JSON() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
}

func Discard() zerolog.Logger {
	return zerolog.Nop()
}
