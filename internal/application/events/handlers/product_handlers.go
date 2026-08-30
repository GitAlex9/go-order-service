package handlers

import (
	"context"
	"fmt"

	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

type LogStockChangedHandler struct {
	log logger.Logger
}

func NewLogStockChangedHandler(log logger.Logger) *LogStockChangedHandler {
	return &LogStockChangedHandler{log: log}
}

func (h *LogStockChangedHandler) Handle(ctx context.Context, evt any) error {
	switch e := evt.(type) {
	case events.ProductStockDecreasedEvent:
		h.log.Info("product stock decreased",
			"product_id", e.ProductID, "old_stock", e.OldStock, "new_stock", e.NewStock,
		)
	case events.ProductStockIncreasedEvent:
		h.log.Info("product stock increased",
			"product_id", e.ProductID, "old_stock", e.OldStock, "new_stock", e.NewStock,
		)
	default:
		return fmt.Errorf("unexpected event type: %T", evt)
	}
	return nil
}
