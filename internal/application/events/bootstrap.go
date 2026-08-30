package events

import (
	"github.com/GitAlex9/go-order-service/internal/application/events/handlers"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

func NewDefaultDispatcher(log logger.Logger) events.Dispatcher {
	registry := NewRegistry()

	orderPaid := handlers.NewLogOrderPaidHandler(log)
	orderCanceled := handlers.NewLogOrderCanceledHandler(log)
	stockChanged := handlers.NewLogStockChangedHandler(log)

	registry.Register("order.paid", orderPaid)
	registry.Register("order.canceled", orderCanceled)
	registry.Register("product.stock_decreased", stockChanged)
	registry.Register("product.stock_increased", stockChanged)

	return NewDispatcher(registry, log)
}
