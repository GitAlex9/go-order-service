package handlers

import (
	"context"
	"fmt"

	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

type LogOrderPaidHandler struct {
	log logger.Logger
}

func NewLogOrderPaidHandler(log logger.Logger) *LogOrderPaidHandler {
	return &LogOrderPaidHandler{log: log}
}

func (h *LogOrderPaidHandler) Handle(ctx context.Context, evt any) error {
	e, ok := evt.(events.OrderPaidEvent)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", evt)
	}

	h.log.Info("order paid",
		"order_id", e.OrderID,
		"customer_id", e.CustomerID,
		"total", e.Total.Amount(),
	)
	return nil
}

type LogOrderCanceledHandler struct {
	log logger.Logger
}

func NewLogOrderCanceledHandler(log logger.Logger) *LogOrderCanceledHandler {
	return &LogOrderCanceledHandler{log: log}
}

func (h *LogOrderCanceledHandler) Handle(ctx context.Context, evt any) error {
	e, ok := evt.(events.OrderCanceledEvent)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", evt)
	}

	h.log.Info("order canceled", "order_id", e.OrderID)
	return nil
}
