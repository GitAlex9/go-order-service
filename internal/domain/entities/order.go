package entities

import (
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

type Order struct {
	ID         string
	CustomerID string

	status OrderStatus
	Items  []OrderItem

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrder(id string, customerID string, items []OrderItem) (*Order, error) {

	order := &Order{
		ID:         id,
		CustomerID: customerID,
		status:     OrderStatusPending,
		Items:      items,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := order.Validate(); err != nil {
		return nil, err
	}

	return order, nil
}

func RebuildOrder(id string, customerID string, status OrderStatus, items []OrderItem, createdAt time.Time, updatedAt time.Time) *Order {

	return &Order{
		ID:         id,
		CustomerID: customerID,
		status:     status,
		Items:      items,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
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

func (o Order) Status() OrderStatus {
	return o.status
}

func (o *Order) SetStatus(status OrderStatus) {
	o.status = status
	o.UpdatedAt = time.Now()
}

func (o Order) Total() float64 {

	var total float64

	for _, item := range o.Items {
		total += item.Subtotal()
	}

	return total
}

func (o *Order) Pay() error {

	if o.Status() != OrderStatusPending {
		return domainerrors.ErrInvalidStatusTransition
	}

	o.SetStatus(OrderStatusPaid)

	return nil
}

func (o *Order) Cancel() error {

	if o.Status() != OrderStatusPending {
		return domainerrors.ErrInvalidStatusTransition
	}

	o.SetStatus(OrderStatusCanceled)

	return nil
}
