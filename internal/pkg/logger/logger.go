package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	With(args ...any) Logger
}

type slogLogger struct {
	l *slog.Logger
}

func New(env string) Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &slogLogger{l: slog.New(handler)}
}

func (s *slogLogger) Debug(msg string, args ...any) { s.l.Debug(msg, args...) }
func (s *slogLogger) Info(msg string, args ...any)  { s.l.Info(msg, args...) }
func (s *slogLogger) Warn(msg string, args ...any)  { s.l.Warn(msg, args...) }
func (s *slogLogger) Error(msg string, args ...any) { s.l.Error(msg, args...) }

func (s *slogLogger) With(args ...any) Logger {
	return &slogLogger{l: s.l.With(args...)}
}

type contextKey string

const loggerContextKey contextKey = "logger"

func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerContextKey).(Logger); ok {
		return l
	}
	return New("development")
}
