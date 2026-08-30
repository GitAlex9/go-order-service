package events

import (
	"context"
	"fmt"

	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

var _ events.Dispatcher = (*dispatcher)(nil)

type dispatcher struct {
	registry *Registry
	log      logger.Logger
}

func NewDispatcher(registry *Registry, log logger.Logger) events.Dispatcher {
	return &dispatcher{registry: registry, log: log}
}

func (d *dispatcher) Dispatch(ctx context.Context, evts []any) {
	for _, evt := range evts {
		named, ok := evt.(events.Named)
		if !ok {
			d.log.Warn("event does not implement Named, skipping dispatch", "type", eventTypeName(evt))
			continue
		}

		eventName := named.EventName()
		handlers := d.registry.HandlersFor(eventName)

		if len(handlers) == 0 {
			continue
		}

		for _, h := range handlers {
			if err := h.Handle(ctx, evt); err != nil {
				d.log.Error("event handler failed",
					"event", eventName,
					"error", err,
				)
			}
		}
	}
}

func eventTypeName(evt any) string {
	return fmt.Sprintf("%T", evt)
}
