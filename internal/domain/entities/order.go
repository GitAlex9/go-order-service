package entities

import (
	"time"

	"github.com/google/uuid"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

type Order struct {
	id         uuid.UUID
	customerID uuid.UUID
	status     OrderStatus
	items      []OrderItem
	createdAt  time.Time
	updatedAt  time.Time
	events     []interface{}
}

func NewOrder(customerID uuid.UUID, items []OrderItem) (*Order, error) {
	now := time.Now()
	order := &Order{
		id:         uuid.New(),
		customerID: customerID,
		status:     OrderStatusPending,
		items:      items,
		createdAt:  now,
		updatedAt:  now,
		events:     []interface{}{},
	}

	if err := order.validate(); err != nil {
		return nil, err
	}

	return order, nil
}

func RebuildOrder(id, customerID uuid.UUID, status OrderStatus, items []OrderItem, createdAt, updatedAt time.Time) *Order {
	return &Order{
		id:         id,
		customerID: customerID,
		status:     status,
		items:      items,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		events:     []interface{}{},
	}
}

func (o Order) ID() uuid.UUID         { return o.id }
func (o Order) CustomerID() uuid.UUID { return o.customerID }
func (o Order) Status() OrderStatus   { return o.status }
func (o Order) Items() []OrderItem    { return o.items }
func (o Order) CreatedAt() time.Time  { return o.createdAt }
func (o Order) UpdatedAt() time.Time  { return o.updatedAt }

// Event handling methods
func (o *Order) AddEvent(event interface{}) {
	o.events = append(o.events, event)
}

func (o *Order) Events() []interface{} {
	return o.events
}

func (o *Order) ClearEvents() {
	o.events = []interface{}{}
}

func (o *Order) AddItem(item OrderItem) error {
	if o.status != OrderStatusPending {
		return domainerrors.ErrOrderNotEditable
	}

	for i, existing := range o.items {
		if existing.ProductID() == item.ProductID() {
			// já existe: soma a quantidade em vez de duplicar a linha
			merged, err := NewOrderItem(existing.ProductID(), existing.ProductName(), existing.UnitPrice(), existing.Quantity()+item.Quantity())
			if err != nil {
				return err
			}
			o.items[i] = *merged
			o.updatedAt = time.Now()
			return nil
		}
	}

	o.items = append(o.items, item)
	o.updatedAt = time.Now()
	return nil
}

func (o *Order) RemoveItem(productID uuid.UUID) error {
	if o.status != OrderStatusPending {
		return domainerrors.ErrOrderNotEditable
	}

	index := -1
	for i, item := range o.items {
		if item.ProductID() == productID {
			index = i
			break
		}
	}
	if index == -1 {
		return domainerrors.ErrOrderItemNotFound
	}

	o.items = append(o.items[:index], o.items[index+1:]...)
	o.updatedAt = time.Now()

	if len(o.items) == 0 {
		return domainerrors.ErrEmptyOrder
	}
	return nil
}

func (o Order) Total() valueobjects.Money {
	total := valueobjects.Zero()
	for _, item := range o.items {
		total = total.Add(item.Subtotal())
	}
	return total
}

func (o *Order) Pay() error {
	if !o.status.CanTransitionTo(OrderStatusPaid) {
		return domainerrors.ErrInvalidStatusTransition
	}
	o.transitionTo(OrderStatusPaid)
	o.AddEvent(events.OrderPaidEvent{
		OrderID:    o.id,
		CustomerID: o.customerID,
		Total:      o.Total(),
	})
	return nil
}

func (o *Order) Cancel() error {
	if !o.status.CanTransitionTo(OrderStatusCanceled) {
		return domainerrors.ErrInvalidStatusTransition
	}
	o.transitionTo(OrderStatusCanceled)
	o.AddEvent(events.OrderCanceledEvent{OrderID: o.id})
	return nil
}

func (o *Order) transitionTo(status OrderStatus) {
	o.status = status
	o.updatedAt = time.Now()
}

func (o Order) validate() error {
	if o.customerID == uuid.Nil {
		return domainerrors.ErrInvalidCustomer
	}
	if len(o.items) == 0 {
		return domainerrors.ErrEmptyOrder
	}
	return nil
}
