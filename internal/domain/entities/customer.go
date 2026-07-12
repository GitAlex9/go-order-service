package entities

import (
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

type Customer struct {
	ID    string
	Name  string
	Email string

	active bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCustomer(id, name, email string) (*Customer, error) {

	customer := &Customer{
		ID:        id,
		Name:      name,
		Email:     email,
		active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := customer.Validate(); err != nil {
		return nil, err
	}

	return customer, nil
}

func RebuildCustomer(
	id string,
	name string,
	email string,
	active bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Customer {

	return &Customer{
		ID:        id,
		Name:      name,
		Email:     email,
		active:    active,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (c Customer) Validate() error {

	if c.Name == "" {
		return domainerrors.ErrInvalidCustomer
	}

	if c.Email == "" {
		return domainerrors.ErrInvalidEmail
	}

	return nil
}

func (c Customer) IsActive() bool {
	return c.active
}

func (c *Customer) Activate() {
	c.active = true
	c.UpdatedAt = time.Now()
}

func (c *Customer) Deactivate() {
	c.active = false
	c.UpdatedAt = time.Now()
}
