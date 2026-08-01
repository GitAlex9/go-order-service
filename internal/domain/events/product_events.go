package events

import (
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type ProductStockDecreasedEvent struct {
	ProductID uuid.UUID
	OldStock  int
	NewStock  int
}

type ProductStockIncreasedEvent struct {
	ProductID uuid.UUID
	OldStock  int
	NewStock  int
}

type ProductPriceChangedEvent struct {
	ProductID uuid.UUID
	OldPrice  valueobjects.Money
	NewPrice  valueobjects.Money
}
