package entities

import (
	"strings"
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type Customer struct {
	id        uuid.UUID
	name      string
	email     valueobjects.Email
	cpf       valueobjects.CPF
	userID    *uuid.UUID //Vincular, no postgress, userID para o usuário.
	createdAt time.Time
	updatedAt time.Time
	events    []interface{}
}

func NewCustomer(name string, email valueobjects.Email, cpf valueobjects.CPF) (*Customer, error) {
	name = strings.TrimSpace(name)
	if len(name) < 3 {
		return nil, domainerrors.ErrEmptyName
	}
	now := time.Now()
	return &Customer{
		id:        uuid.New(),
		name:      name,
		email:     email,
		cpf:       cpf,
		createdAt: now,
		updatedAt: now,
		events:    []interface{}{},
	}, nil
}

func RebuildCustomer(id uuid.UUID, name string, email valueobjects.Email, cpf valueobjects.CPF, userID *uuid.UUID, createdAt, updatedAt time.Time) *Customer {
	return &Customer{
		id:        id,
		name:      name,
		email:     email,
		cpf:       cpf,
		userID:    userID,
		createdAt: createdAt,
		updatedAt: updatedAt,
		events:    []interface{}{},
	}
}

func (c *Customer) ID() uuid.UUID             { return c.id }
func (c *Customer) Name() string              { return c.name }
func (c *Customer) Email() valueobjects.Email { return c.email }
func (c *Customer) CPF() valueobjects.CPF     { return c.cpf }
func (c *Customer) UserID() *uuid.UUID        { return c.userID }
func (c *Customer) CreatedAt() time.Time      { return c.createdAt }
func (c *Customer) UpdatedAt() time.Time      { return c.updatedAt }

func (c *Customer) AddEvent(event interface{}) {
	c.events = append(c.events, event)
}

func (c *Customer) Events() []interface{} {
	return c.events
}

func (c *Customer) ClearEvents() {
	c.events = []interface{}{}
}

func (c *Customer) ChangeEmail(newEmail valueobjects.Email) {
	oldEmail := c.email
	c.email = newEmail
	c.updatedAt = time.Now()
	if c.userID != nil {
		c.AddEvent(events.UserEmailChangedEvent{UserID: *c.userID, OldEmail: oldEmail, NewEmail: newEmail})
	}
}

func (c *Customer) Rename(newName string) error {
	oldName := c.name
	newName = strings.TrimSpace(newName)
	if len(newName) < 3 {
		return domainerrors.ErrEmptyName
	}
	c.name = newName
	c.updatedAt = time.Now()
	c.AddEvent(events.CustomerRenamedEvent{CustomerID: c.id, OldName: oldName, NewName: newName})
	return nil
}

func (c *Customer) LinkUser(userID uuid.UUID) {
	c.userID = &userID
	c.updatedAt = time.Now()
	c.AddEvent(events.CustomerLinkedToUserEvent{CustomerID: c.id, UserID: userID})
}
