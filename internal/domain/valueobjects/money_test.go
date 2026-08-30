package valueobjects

import (
	"errors"
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name    string
		cents   int64
		want    int64
		wantErr error
	}{
		{"valor positivo", 1990, 1990, nil},
		{"valor zero", 0, 0, nil},
		{"valor negativo", -100, 0, domainerrors.ErrNegativeMoneyAmount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMoney(tt.cents)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewMoney(%d) error = %v, want %v", tt.cents, err, tt.wantErr)
			}
			if tt.wantErr == nil && got.Cents() != tt.want {
				t.Errorf("got.Cents() = %d, want %d", got.Cents(), tt.want)
			}
		})
	}
}

func TestNewMoneyFromFloat(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		wantCents int64
		wantErr   error
	}{
		{"valor exato", 19.90, 1990, nil},
		{"valor com arredondamento para cima", 19.999, 2000, nil},
		{"valor zero", 0, 0, nil},
		{"valor negativo", -10.00, 0, domainerrors.ErrNegativeMoneyAmount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMoneyFromFloat(tt.amount)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewMoneyFromFloat(%v) error = %v, want %v", tt.amount, err, tt.wantErr)
			}
			if tt.wantErr == nil && got.Cents() != tt.wantCents {
				t.Errorf("got.Cents() = %d, want %d", got.Cents(), tt.wantCents)
			}
		})
	}
}

func TestZero(t *testing.T) {
	got := Zero()
	want := int64(0)

	if got.Cents() != want {
		t.Errorf("Zero().Cents() = %d, want %d", got.Cents(), want)
	}
	if !got.IsZero() {
		t.Errorf("Zero().IsZero() = false, want true")
	}
}

func TestMoney_Amount(t *testing.T) {
	tests := []struct {
		name  string
		cents int64
		want  float64
	}{
		{"valor inteiro", 100000, 1000.00},
		{"valor com centavos", 1990, 19.90},
		{"zero", 0, 0.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoney(tt.cents)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := m.Amount()
			if got != tt.want {
				t.Errorf("Amount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_String(t *testing.T) {
	m, err := NewMoney(350000)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	got := m.String()
	want := "R$ 3500.00"

	if got != want {
		t.Errorf("String() got = %q, want %q", got, want)
	}
}

func TestMoney_Add(t *testing.T) {
	tests := []struct {
		name string
		a    int64
		b    int64
		want int64
	}{
		{"soma dois valores positivos", 1000, 500, 1500},
		{"soma com zero", 1000, 0, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewMoney(tt.a)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}
			b, err := NewMoney(tt.b)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := a.Add(b)
			if got.Cents() != tt.want {
				t.Errorf("Add() got = %d, want %d", got.Cents(), tt.want)
			}
		})
	}
}

func TestMoney_Multiply(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		quantity int
		want     int64
	}{
		{"multiplica por quantidade positiva", 1000, 3, 3000},
		{"multiplica por zero", 1000, 0, 0},
		{"multiplica por um", 1000, 1, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoney(tt.cents)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := m.Multiply(tt.quantity)
			if got.Cents() != tt.want {
				t.Errorf("Multiply(%d) got = %d, want %d", tt.quantity, got.Cents(), tt.want)
			}
		})
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	tests := []struct {
		name string
		a    int64
		b    int64
		want bool
	}{
		{"a maior que b", 2000, 1000, true},
		{"a menor que b", 1000, 2000, false},
		{"a igual a b", 1000, 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewMoney(tt.a)
			b, _ := NewMoney(tt.b)

			got := a.GreaterThan(b)
			if got != tt.want {
				t.Errorf("GreaterThan() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		cents int64
		want  bool
	}{
		{"valor zero", 0, true},
		{"valor positivo", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMoney(tt.cents)

			got := m.IsZero()
			if got != tt.want {
				t.Errorf("IsZero() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    int64
		b    int64
		want bool
	}{
		{"valores iguais", 1000, 1000, true},
		{"valores diferentes", 1000, 2000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewMoney(tt.a)
			b, _ := NewMoney(tt.b)

			got := a.Equals(b)
			if got != tt.want {
				t.Errorf("Equals() got = %v, want %v", got, tt.want)
			}
		})
	}
}
