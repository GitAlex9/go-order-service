package valueobjects

import (
	"fmt"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

type Money struct {
	cents int64
}

func NewMoney(cents int64) (Money, error) {
	if cents < 0 {
		return Money{}, domainerrors.ErrNegativeMoneyAmount
	}
	return Money{cents: cents}, nil
}

func NewMoneyFromFloat(amount float64) (Money, error) {
	if amount < 0 {
		return Money{}, domainerrors.ErrNegativeMoneyAmount
	}
	cents := int64(amount*100 + 0.5)
	return Money{cents: cents}, nil
}

func Zero() Money {
	return Money{cents: 0}
}

func (m Money) Cents() int64 {
	return m.cents
}

func (m Money) Amount() float64 {
	return float64(m.cents) / 100
}

func (m Money) String() string {
	return fmt.Sprintf("R$ %.2f", m.Amount())
}

func (m Money) Add(other Money) Money {
	return Money{cents: m.cents + other.cents}
}

func (m Money) Multiply(quantity int) Money {
	return Money{cents: m.cents * int64(quantity)}
}

func (m Money) GreaterThan(other Money) bool {
	return m.cents > other.cents
}

func (m Money) IsZero() bool {
	return m.cents == 0
}

func (m Money) Equals(other Money) bool {
	return m.cents == other.cents
}
