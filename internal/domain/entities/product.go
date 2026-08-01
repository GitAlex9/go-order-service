package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

type Product struct {
	id          uuid.UUID
	name        string
	description string
	price       valueobjects.Money
	stock       int
	active      bool
	createdAt   time.Time
	updatedAt   time.Time
	events      []interface{}
}

func NewProduct(name, description string, price valueobjects.Money, stock int) (*Product, error) {
	product := &Product{
		id:          uuid.New(),
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		active:      true,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
		events:      []interface{}{},
	}

	if err := product.validate(); err != nil {
		return nil, err
	}

	return product, nil
}

func RebuildProduct(
	id uuid.UUID,
	name, description string,
	price valueobjects.Money,
	stock int,
	active bool,
	createdAt, updatedAt time.Time,
) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		active:      active,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		events:      []interface{}{},
	}
}

func (p Product) ID() uuid.UUID             { return p.id }
func (p Product) Name() string              { return p.name }
func (p Product) Description() string       { return p.description }
func (p Product) Price() valueobjects.Money { return p.price }
func (p Product) Stock() int                { return p.stock }
func (p Product) IsActive() bool            { return p.active }
func (p Product) CreatedAt() time.Time      { return p.createdAt }
func (p Product) UpdatedAt() time.Time      { return p.updatedAt }

func (p *Product) AddEvent(event interface{}) {
	p.events = append(p.events, event)
}

func (p *Product) Events() []interface{} {
	return p.events
}

func (p *Product) ClearEvents() {
	p.events = []interface{}{}
}

func (p Product) HasStock(quantity int) bool {
	return p.stock >= quantity
}

func (p *Product) Rename(newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return domainerrors.ErrInvalidProductName
	}
	p.name = newName
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangeDescription(newDescription string) error {
	newDescription = strings.TrimSpace(newDescription)
	if newDescription == "" {
		return domainerrors.ErrInvalidProductDescription
	}
	p.description = newDescription
	p.updatedAt = time.Now()
	return nil
}

func (p *Product) ChangePrice(newPrice valueobjects.Money) error {
	if newPrice.IsZero() {
		return domainerrors.ErrInvalidProductPrice
	}
	oldPrice := p.price
	p.price = newPrice
	p.updatedAt = time.Now()
	p.AddEvent(events.ProductPriceChangedEvent{ProductID: p.id, OldPrice: oldPrice, NewPrice: newPrice})
	return nil
}

func (p *Product) IncreaseStock(quantity int) error {
	if quantity <= 0 {
		return domainerrors.ErrInvalidQuantity
	}
	oldStock := p.stock
	p.stock += quantity
	p.updatedAt = time.Now()
	p.AddEvent(events.ProductStockIncreasedEvent{ProductID: p.id, OldStock: oldStock, NewStock: p.stock})
	return nil
}

func (p *Product) DecreaseStock(quantity int) error {
	if quantity <= 0 {
		return domainerrors.ErrInvalidQuantity
	}
	if !p.HasStock(quantity) {
		return domainerrors.ErrInsufficientStock
	}
	oldStock := p.stock
	p.stock -= quantity
	p.updatedAt = time.Now()
	p.AddEvent(events.ProductStockDecreasedEvent{ProductID: p.id, OldStock: oldStock, NewStock: p.stock})
	return nil
}

func (p *Product) Activate() {
	p.active = true
	p.updatedAt = time.Now()
}

func (p *Product) Deactivate() {
	p.active = false
	p.updatedAt = time.Now()
}

func (p Product) validate() error {
	if p.name == "" {
		return domainerrors.ErrInvalidProductName
	}
	if p.description == "" {
		return domainerrors.ErrInvalidProductDescription
	}
	if p.price.IsZero() {
		return domainerrors.ErrInvalidProductPrice
	}
	if p.stock < 0 {
		return domainerrors.ErrInvalidProductStock
	}
	return nil
}
