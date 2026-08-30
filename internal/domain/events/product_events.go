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

func (ProductStockDecreasedEvent) EventName() string { return "product.stock_decreased" }

type ProductStockIncreasedEvent struct {
	ProductID uuid.UUID
	OldStock  int
	NewStock  int
}

func (ProductStockIncreasedEvent) EventName() string { return "product.stock_increased" }

type ProductPriceChangedEvent struct {
	ProductID uuid.UUID
	OldPrice  valueobjects.Money
	NewPrice  valueobjects.Money
}

func (ProductPriceChangedEvent) EventName() string { return "product.price_change" }
