package entities

import (
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64

	stock  int
	active bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProduct(id, name, description string, price float64, stock int) (*Product, error) {

	product := &Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		stock:       stock,
		active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	return product, nil
}

func (p Product) Validate() error {

	if p.Name == "" {
		return domainerrors.ErrInvalidProductName
	}

	if p.Description == "" {
		return domainerrors.ErrInvalidProductDescription
	}

	if p.Price <= 0 {
		return domainerrors.ErrInvalidProductPrice
	}

	if p.stock < 0 {
		return domainerrors.ErrInvalidProductStock
	}

	return nil
}

func (p Product) Stock() int {
	return p.stock
}

func (p Product) IsActive() bool {
	return p.active
}

func (p Product) HasStock(quantity int) bool {
	return p.stock >= quantity
}

func (p *Product) IncreaseStock(quantity int) error {

	if quantity <= 0 {
		return domainerrors.ErrInvalidQuantity
	}

	p.stock += quantity
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Product) DecreaseStock(quantity int) error {

	if quantity <= 0 {
		return domainerrors.ErrInvalidQuantity
	}

	if !p.HasStock(quantity) {
		return domainerrors.ErrInsufficientStock
	}

	p.stock -= quantity
	p.UpdatedAt = time.Now()

	return nil
}

func (p *Product) Activate() {
	p.active = true
	p.UpdatedAt = time.Now()
}

func (p *Product) Deactivate() {
	p.active = false
	p.UpdatedAt = time.Now()
}
