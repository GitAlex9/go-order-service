package events

import (
	"context"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/domain/events"
)

type fakeHandler struct {
	calls int
}

func (h *fakeHandler) Handle(ctx context.Context, event any) error {
	h.calls++
	return nil
}

func TestRegistry_RegisterAndHandlersFor(t *testing.T) {
	registry := NewRegistry()

	h1 := &fakeHandler{}
	h2 := &fakeHandler{}

	registry.Register("order.paid", h1)
	registry.Register("order.paid", h2)
	registry.Register("order.canceled", h1)

	tests := []struct {
		name      string
		eventName string
		wantCount int
	}{
		{"evento com dois handlers", "order.paid", 2},
		{"evento com um handler", "order.canceled", 1},
		{"evento sem handler registrado", "product.stock_decreased", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := registry.HandlersFor(tt.eventName)
			if len(got) != tt.wantCount {
				t.Errorf("HandlersFor(%q) got %d handlers, want %d", tt.eventName, len(got), tt.wantCount)
			}
		})
	}
}

func TestRegistry_HandlersFor_NeverReturnsNilImplicitly(t *testing.T) {
	registry := NewRegistry()

	got := registry.HandlersFor("nonexistent.event")

	if len(got) != 0 {
		t.Errorf("HandlersFor() for unregistered event got %d handlers, want 0", len(got))
	}
}

var _ events.Handler = (*fakeHandler)(nil)
