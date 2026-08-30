package events

import (
	"github.com/GitAlex9/go-order-service/internal/domain/events"
)

type Registry struct {
	handlers map[string][]events.Handler
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string][]events.Handler)}
}

func (r *Registry) Register(eventName string, handler events.Handler) {
	r.handlers[eventName] = append(r.handlers[eventName], handler)
}

func (r *Registry) HandlersFor(eventName string) []events.Handler {
	return r.handlers[eventName]
}
