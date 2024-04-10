package logger

import (
	"github.com/rs/zerolog"
	"os"
)

type Logger struct {
	zerolog.Logger
}

func NewLogger(serviceName string) *Logger {
	lg := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().Str("service", serviceName).
		Timestamp().Logger()

	return &Logger{
		Logger: lg,
	}
}
