package logger

import (
	"context"
	"testing"
)

func TestNew_DoesNotPanic(t *testing.T) {
	tests := []string{"development", "production", "test", ""}

	for _, env := range tests {
		t.Run(env, func(t *testing.T) {
			log := New(env)
			if log == nil {
				t.Fatalf("New(%q) returned nil", env)
			}

			log.Debug("debug message", "key", "value")
			log.Info("info message", "key", "value")
			log.Warn("warn message", "key", "value")
			log.Error("error message", "key", "value")
		})
	}
}

func TestLogger_With(t *testing.T) {
	log := New("test")

	scoped := log.With("request_id", "abc-123")

	if scoped == nil {
		t.Fatal("With() returned nil")
	}

	scoped.Info("scoped message")
}

func TestWithContext_And_FromContext(t *testing.T) {
	log := New("test")
	ctx := WithContext(context.Background(), log)

	got := FromContext(ctx)

	if got != log {
		t.Errorf("FromContext() did not return the logger that was stored in the context")
	}
}

func TestFromContext_WithoutStoredLogger_ReturnsDefault(t *testing.T) {
	ctx := context.Background()

	got := FromContext(ctx)

	if got == nil {
		t.Fatal("FromContext() on empty context returned nil, want a default logger")
	}
}
