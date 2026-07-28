package entities

import (
	"github.com/google/uuid"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

type OrderItem struct {
	productID   uuid.UUID
	productName string
	unitPrice   valueobjects.Money
	quantity    int
}

func NewOrderItem(productID uuid.UUID, productName string, unitPrice valueobjects.Money, quantity int) (*OrderItem, error) {
	item := &OrderItem{
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}

	if err := item.validate(); err != nil {
		return nil, err
	}

	return item, nil
}

func RebuildOrderItem(productID uuid.UUID, productName string, unitPrice valueobjects.Money, quantity int) *OrderItem {
	return &OrderItem{
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}
}

func (oi OrderItem) ProductID() uuid.UUID          { return oi.productID }
func (oi OrderItem) ProductName() string           { return oi.productName }
func (oi OrderItem) UnitPrice() valueobjects.Money { return oi.unitPrice }
func (oi OrderItem) Quantity() int                 { return oi.quantity }

func (oi OrderItem) Subtotal() valueobjects.Money {
	return oi.unitPrice.Multiply(oi.quantity)
}

func (oi OrderItem) validate() error {
	if oi.productID == uuid.Nil {
		return domainerrors.ErrInvalidProductID
	}
	if oi.productName == "" {
		return domainerrors.ErrInvalidProductName
	}
	if oi.unitPrice.IsZero() {
		return domainerrors.ErrInvalidProductPrice
	}
	if oi.quantity <= 0 {
		return domainerrors.ErrInvalidQuantity
	}
	return nil
}
