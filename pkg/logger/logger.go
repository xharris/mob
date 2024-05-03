package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	handlerOptions *slog.HandlerOptions
	logger         *slog.Logger
}

func NewLogger(opts ...LoggerOption) Logger {
	handler := &slog.HandlerOptions{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handler))
	l := Logger{
		handlerOptions: handler,
		logger:         logger,
	}
	for _, opt := range opts {
		opt(&l)
	}
	return l
}

type LoggerOption func(*Logger)

func WithDebug() LoggerOption {
	return func(l *Logger) {
		l.handlerOptions.Level = slog.LevelDebug
	}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
