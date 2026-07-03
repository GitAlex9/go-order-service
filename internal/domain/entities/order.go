package entities

import (
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

type Order struct {
	ID         string
	CustomerID string

	Status OrderStatus
	Items  []OrderItem

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrder(id string, customerID string, items []OrderItem) (*Order, error) {

	order := &Order{
		ID:         id,
		CustomerID: customerID,
		Status:     OrderStatusPending,
		Items:      items,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}

func (o Order) Validate() error {

	if o.CustomerID == "" {
		return domainerrors.ErrInvalidCustomer
	}

	if len(o.Items) == 0 {
		return domainerrors.ErrEmptyOrder
	}

	return nil
}

func (o Order) Total() float64 {

	var total float64

	for _, item := range o.Items {
		total += item.Subtotal()
	}

	return total
}

func (o *Order) Pay() error {

	if o.Status != OrderStatusPending {
		return domainerrors.ErrInvalidStatusTransition
	}

	o.Status = OrderStatusPaid
	o.UpdatedAt = time.Now()

	return nil
}

func (o *Order) Cancel() error {

	if o.Status != OrderStatusPending {
		return domainerrors.ErrInvalidStatusTransition
	}

	o.Status = OrderStatusCanceled
	o.UpdatedAt = time.Now()

	return nil
}
