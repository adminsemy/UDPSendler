package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	options  *slog.HandlerOptions
	logLevel *slog.LevelVar
	*slog.Logger
}

func NewLogger() *Logger {
	logger := &Logger{}
	logLevel := &slog.LevelVar{}
	logger.logLevel = logLevel
	logLevel.Set(slog.LevelDebug)
	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("02.01.2006 15:04:05"))
			}
			return a
		},
	}
	logger.options = &opts
	handler := slog.NewTextHandler(os.Stdout, &opts)
	logger.Logger = slog.New(handler)

	return logger
}

func (l *Logger) SetLevel(level slog.Level) {
	l.logLevel.Set(level)
}
