package events

import (
	"testing"

	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

func TestNewDefaultDispatcher_ReturnsNonNilDispatcher(t *testing.T) {
	got := NewDefaultDispatcher(logger.New("test"))

	if got == nil {
		t.Fatal("NewDefaultDispatcher() returned nil")
	}
}
