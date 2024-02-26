package logger

import (
	"log/slog"
	"os"
)

type logger struct {
	options  *slog.HandlerOptions
	logLevel *slog.LevelVar
	*slog.Logger
}

func NewLogger() *logger {
	logger := &logger{}
	logLevel := &slog.LevelVar{}
	logger.logLevel = logLevel
	logLevel.Set(slog.LevelDebug)
	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}
	logger.options = &opts
	handler := slog.NewTextHandler(os.Stdout, &opts)
	logger.Logger = slog.New(handler)
	return logger
}

func (l *logger) SetLevel(level slog.Level) {
	l.logLevel.Set(level)
}

func (l *logger) AddSource(add bool) {
	l.options.AddSource = add
}
