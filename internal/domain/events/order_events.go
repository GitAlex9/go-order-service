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

type OrderCanceledEvent struct {
	OrderID uuid.UUID
}
