package entities

import domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"

type OrderItem struct {
	ProductID   string
	ProductName string
	UnitPrice   float64
	Quantity    int
}

func NewOrderItem(productID string, productName string, unitPrice float64, quantity int) (*OrderItem, error) {

	item := &OrderItem{
		ProductID:   productID,
		ProductName: productName,
		UnitPrice:   unitPrice,
		Quantity:    quantity,
	}

	if err := item.Validate(); err != nil {
		return nil, err
	}

	return item, nil
}

func (oi OrderItem) Validate() error {

	if oi.ProductID == "" {
		return domainerrors.ErrInvalidProductID
	}

	if oi.ProductName == "" {
		return domainerrors.ErrInvalidProductName
	}

	if oi.UnitPrice <= 0 {
		return domainerrors.ErrInvalidProductPrice
	}

	if oi.Quantity <= 0 {
		return domainerrors.ErrInvalidQuantity
	}

	return nil
}

func (oi OrderItem) Subtotal() float64 {
	return oi.UnitPrice * float64(oi.Quantity)
}
