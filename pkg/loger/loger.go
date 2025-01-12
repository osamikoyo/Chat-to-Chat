package loger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

type Logger struct {
	*zerolog.Logger
}

func New() Logger {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	return Logger{Logger: &logger}
}
