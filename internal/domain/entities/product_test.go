package entities

import (
	"errors"
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

func newTestPrice(t *testing.T, cents int64) valueobjects.Money {
	t.Helper()

	price, err := valueobjects.NewMoney(cents)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	return price
}

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		description string
		price       int64
		stock       int
		wantErr     error
	}{
		{"produto válido", "Notebook", "Notebook gamer", 350000, 10, nil},
		{"nome vazio", "", "Notebook gamer", 350000, 10, domainerrors.ErrInvalidProductName},
		{"descrição vazia", "Notebook", "", 350000, 10, domainerrors.ErrInvalidProductDescription},
		{"preço zero", "Notebook", "Notebook gamer", 0, 10, domainerrors.ErrInvalidProductPrice},
		{"estoque negativo", "Notebook", "Notebook gamer", 350000, -1, domainerrors.ErrInvalidProductStock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price := newTestPrice(t, tt.price)

			got, err := NewProduct(tt.productName, tt.description, price, tt.stock)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewProduct() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil && got != nil {
				t.Errorf("got product = %v, want nil", got)
			}
			if tt.wantErr == nil && got == nil {
				t.Errorf("got nil product, want non-nil")
			}
		})
	}
}

func TestNewProduct_FieldsArePersisted(t *testing.T) {
	price := newTestPrice(t, 350000)

	product, err := NewProduct("Notebook", "Notebook gamer", price, 10)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if got, want := product.Name(), "Notebook"; got != want {
		t.Errorf("Name() got = %q, want %q", got, want)
	}
	if got, want := product.Stock(), 10; got != want {
		t.Errorf("Stock() got = %d, want %d", got, want)
	}
	if got := product.IsActive(); !got {
		t.Errorf("IsActive() got = %v, want true", got)
	}
}

func TestProduct_HasStock(t *testing.T) {
	product, err := NewProduct("Notebook", "Notebook gamer", newTestPrice(t, 350000), 10)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		name     string
		quantity int
		want     bool
	}{
		{"quantidade menor que o estoque", 5, true},
		{"quantidade igual ao estoque", 10, true},
		{"quantidade maior que o estoque", 11, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := product.HasStock(tt.quantity)
			if got != tt.want {
				t.Errorf("HasStock(%d) got = %v, want %v", tt.quantity, got, tt.want)
			}
		})
	}
}

func TestProduct_IncreaseStock(t *testing.T) {
	tests := []struct {
		name      string
		quantity  int
		wantErr   error
		wantStock int
	}{
		{"quantidade positiva", 5, nil, 15},
		{"quantidade zero", 0, domainerrors.ErrInvalidQuantity, 10},
		{"quantidade negativa", -5, domainerrors.ErrInvalidQuantity, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct("Notebook", "Notebook gamer", newTestPrice(t, 350000), 10)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err = product.IncreaseStock(tt.quantity)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("IncreaseStock(%d) error = %v, want %v", tt.quantity, err, tt.wantErr)
			}

			if got := product.Stock(); got != tt.wantStock {
				t.Errorf("Stock() got = %d, want %d", got, tt.wantStock)
			}
		})
	}
}

func TestProduct_DecreaseStock(t *testing.T) {
	tests := []struct {
		name      string
		quantity  int
		wantErr   error
		wantStock int
	}{
		{"quantidade válida", 5, nil, 5},
		{"quantidade igual ao estoque total", 10, nil, 0},
		{"quantidade maior que o estoque", 11, domainerrors.ErrInsufficientStock, 10},
		{"quantidade zero", 0, domainerrors.ErrInvalidQuantity, 10},
		{"quantidade negativa", -5, domainerrors.ErrInvalidQuantity, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct("Notebook", "Notebook gamer", newTestPrice(t, 350000), 10)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err = product.DecreaseStock(tt.quantity)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecreaseStock(%d) error = %v, want %v", tt.quantity, err, tt.wantErr)
			}

			if got := product.Stock(); got != tt.wantStock {
				t.Errorf("Stock() got = %d, want %d", got, tt.wantStock)
			}
		})
	}
}

func TestProduct_ActivateDeactivate(t *testing.T) {
	product, err := NewProduct("Notebook", "Notebook gamer", newTestPrice(t, 350000), 10)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	product.Deactivate()
	if got := product.IsActive(); got {
		t.Errorf("IsActive() after Deactivate() got = %v, want false", got)
	}

	product.Activate()
	if got := product.IsActive(); !got {
		t.Errorf("IsActive() after Activate() got = %v, want true", got)
	}
}

func TestProduct_Rename(t *testing.T) {
	tests := []struct {
		name    string
		newName string
		want    string
		wantErr error
	}{
		{"nome válido", "Notebook Pro", "Notebook Pro", nil},
		{"nome com espaços nas bordas", "  Notebook Pro  ", "Notebook Pro", nil},
		{"nome vazio", "", "", domainerrors.ErrInvalidProductName},
		{"nome só com espaços", "   ", "", domainerrors.ErrInvalidProductName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct("Notebook", "Notebook gamer", newTestPrice(t, 350000), 10)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err = product.Rename(tt.newName)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Rename(%q) error = %v, want %v", tt.newName, err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if got := product.Name(); got != tt.want {
					t.Errorf("Name() got = %q, want %q", got, tt.want)
				}
			}
		})
	}
}

func TestProduct_ChangeDescription(t *testing.T) {
	tests := []struct {
		name           string
		newDescription string
		want           string
		wantErr        error
	}{
		{"descrição válida", "Nova descrição", "Nova descrição", nil},
		{"descrição vazia", "", "", domainerrors.ErrInvalidProductDescription},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct("Notebook", "Notebook gamer", newTestPrice(t, 350000), 10)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			err = product.ChangeDescription(tt.newDescription)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ChangeDescription(%q) error = %v, want %v", tt.newDescription, err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if got := product.Description(); got != tt.want {
					t.Errorf("Description() got = %q, want %q", got, tt.want)
				}
			}
		})
	}
}

func TestProduct_ChangePrice(t *testing.T) {
	tests := []struct {
		name      string
		newCents  int64
		wantErr   error
		wantCents int64
	}{
		{"preço válido", 400000, nil, 400000},
		{"preço zero", 0, domainerrors.ErrInvalidProductPrice, 350000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct("Notebook", "Notebook gamer", newTestPrice(t, 350000), 10)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			newPrice := newTestPrice(t, tt.newCents)

			err = product.ChangePrice(newPrice)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ChangePrice() error = %v, want %v", err, tt.wantErr)
			}

			if got := product.Price().Cents(); got != tt.wantCents {
				t.Errorf("Price().Cents() got = %d, want %d", got, tt.wantCents)
			}
		})
	}
}
