package entities

import (
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

type Customer struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCustomer(id string, name string) (*Customer, error) {

	customer := &Customer{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := customer.Validate(); err != nil {
		return nil, err
	}

	return customer, nil
}

func (c Customer) Validate() error {

	if c.Name == "" {
		return domainerrors.ErrInvalidCustomer
	}

	return nil
}
