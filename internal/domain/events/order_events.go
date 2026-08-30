package events

import (
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type OrderPaidEvent struct {
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	Total      valueobjects.Money
}

func (OrderPaidEvent) EventName() string { return "order.paid" }

type OrderCanceledEvent struct {
	OrderID uuid.UUID
}

func (OrderCanceledEvent) EventName() string { return "order.canceled" }
