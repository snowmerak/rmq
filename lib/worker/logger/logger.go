package logger

import (
	"github.com/rs/zerolog"
	"io"
)

//go:bean
type Logger struct {
	zerolog.Logger
}

func New(w io.Writer) *Logger {
	return &Logger{
		Logger: zerolog.New(w).With().Timestamp().Logger(),
	}
}
