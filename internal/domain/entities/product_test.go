package entities

import (
	"testing"
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestNewProduct_Valid(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)

	tests := []struct {
		name        string
		description string
		price       valueobjects.Money
		stock       int
	}{
		{"produto normal", "Descrição teste", price, 10},
		{"estoque zero", "Descrição teste", price, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.name, tt.description, tt.price, tt.stock)
			if err != nil {
				t.Fatalf("expected creation to succeed, got error: %v", err)
			}
			if product.ID() == uuid.Nil {
				t.Errorf("expected ID to be set")
			}
			if product.Name() != tt.name {
				t.Errorf("Name() got = %q, want %q", product.Name(), tt.name)
			}
			if product.Description() != tt.description {
				t.Errorf("Description() got = %q, want %q", product.Description(), tt.description)
			}
			if !product.Price().Equals(tt.price) {
				t.Errorf("Price() got = %v, want %v", product.Price().Amount(), tt.price.Amount())
			}
			if product.Stock() != tt.stock {
				t.Errorf("Stock() got = %d, want %d", product.Stock(), tt.stock)
			}
			if !product.IsActive() {
				t.Errorf("expected product to be active")
			}
			if product.CreatedAt().IsZero() || product.UpdatedAt().IsZero() {
				t.Errorf("expected timestamps to be set")
			}
		})
	}
}

func TestNewProduct_Invalid(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)

	tests := []struct {
		name        string
		description string
		price       valueobjects.Money
		stock       int
		wantErr     error
	}{
		{"nome vazio", "Descrição", price, 10, domainerrors.ErrInvalidProductName},
		{"descrição vazia", "Produto", price, 10, domainerrors.ErrInvalidProductDescription},
		{"preço zero", "Produto", valueobjects.Zero(), 10, domainerrors.ErrInvalidProductPrice},
		{"estoque negativo", "Produto", price, -1, domainerrors.ErrInvalidProductStock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.name, tt.description, tt.price, tt.stock)
			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if product != nil {
				t.Errorf("expected nil product")
			}
		})
	}
}

func TestProduct_Rename(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	product, _ := NewProduct("Nome Antigo", "Descrição", price, 10)

	tests := []struct {
		name    string
		newName string
		wantErr bool
		want    string
	}{
		{"renomear válido", "Novo Nome", false, "Novo Nome"},
		{"renomear com espaços", "  Outro Nome  ", false, "Outro Nome"},
		{"renomear vazio", "", true, "Nome Antigo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := product.UpdatedAt()
			err := product.Rename(tt.newName)

			if tt.wantErr {
				if err != domainerrors.ErrInvalidProductName {
					t.Errorf("expected ErrInvalidProductName, got %v", err)
				}
				if product.Name() != tt.want {
					t.Errorf("Name() got = %q, want %q", product.Name(), tt.want)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if product.Name() != tt.want {
				t.Errorf("Name() got = %q, want %q", product.Name(), tt.want)
			}
			if product.UpdatedAt() == oldUpdatedAt {
				t.Errorf("expected UpdatedAt to change")
			}
		})
	}
}

func TestProduct_ChangeDescription(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	product, _ := NewProduct("Produto", "Descrição Antiga", price, 10)

	tests := []struct {
		name    string
		newDesc string
		wantErr bool
		want    string
	}{
		{"descrição válida", "Nova Descrição", false, "Nova Descrição"},
		{"descrição com espaços", "  Outra Descrição  ", false, "Outra Descrição"},
		{"descrição vazia", "", true, "Descrição Antiga"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := product.UpdatedAt()
			err := product.ChangeDescription(tt.newDesc)

			if tt.wantErr {
				if err != domainerrors.ErrInvalidProductDescription {
					t.Errorf("expected ErrInvalidProductDescription, got %v", err)
				}
				if product.Description() != tt.want {
					t.Errorf("Description() got = %q, want %q", product.Description(), tt.want)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if product.Description() != tt.want {
				t.Errorf("Description() got = %q, want %q", product.Description(), tt.want)
			}
			if product.UpdatedAt() == oldUpdatedAt {
				t.Errorf("expected UpdatedAt to change")
			}
		})
	}
}

func TestProduct_ChangePrice(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	product, _ := NewProduct("Produto", "Descrição", price, 10)

	newPrice, _ := valueobjects.NewMoneyFromFloat(15.75)
	oldUpdatedAt := product.UpdatedAt()

	err := product.ChangePrice(newPrice)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !product.Price().Equals(newPrice) {
		t.Errorf("Price() got = %v, want %v", product.Price().Amount(), newPrice.Amount())
	}
	if product.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	// Verifica o evento
	productEvents := product.Events()
	lastEvent := productEvents[len(productEvents)-1]
	priceEvent, ok := lastEvent.(events.ProductPriceChangedEvent)
	if !ok {
		t.Errorf("expected ProductPriceChangedEvent, got %T", lastEvent)
	}
	if priceEvent.ProductID != product.ID() {
		t.Errorf("expected ProductID %v, got %v", product.ID(), priceEvent.ProductID)
	}
}

func TestProduct_ChangePrice_Invalid(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	product, _ := NewProduct("Produto", "Descrição", price, 10)

	zeroPrice := valueobjects.Zero()
	err := product.ChangePrice(zeroPrice)
	if err != domainerrors.ErrInvalidProductPrice {
		t.Errorf("expected ErrInvalidProductPrice, got %v", err)
	}
}

func TestProduct_IncreaseStock(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	product, _ := NewProduct("Produto", "Descrição", price, 10)

	tests := []struct {
		name     string
		quantity int
		want     int
		wantErr  bool
	}{
		{"aumentar 5", 5, 15, false},
		{"aumentar 1", 1, 16, false},
		{"quantidade zero", 0, 16, true},
		{"quantidade negativa", -1, 16, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := product.UpdatedAt()
			err := product.IncreaseStock(tt.quantity)

			if tt.wantErr {
				if err != domainerrors.ErrInvalidQuantity {
					t.Errorf("expected ErrInvalidQuantity, got %v", err)
				}
				if product.Stock() == tt.want {
					t.Errorf("Stock() should not change, got %d", product.Stock())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if product.Stock() != tt.want {
				t.Errorf("Stock() got = %d, want %d", product.Stock(), tt.want)
			}
			if product.UpdatedAt() == oldUpdatedAt {
				t.Errorf("expected UpdatedAt to change")
			}

			productEvents := product.Events()
			if len(productEvents) > 0 {
				lastEvent := productEvents[len(productEvents)-1]
				stockEvent, ok := lastEvent.(events.ProductStockIncreasedEvent)
				if !ok {
					t.Errorf("expected ProductStockIncreasedEvent, got %T", lastEvent)
				} else {
					if stockEvent.NewStock != product.Stock() {
						t.Errorf("NewStock got = %d, want %d", stockEvent.NewStock, product.Stock())
					}
				}
			}
		})
	}
}

func TestProduct_DecreaseStock(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	product, _ := NewProduct("Produto", "Descrição", price, 10)

	tests := []struct {
		name     string
		quantity int
		want     int
		wantErr  error
	}{
		{"diminuir 3", 3, 7, nil},
		{"diminuir 1", 1, 6, nil},
		{"quantidade zero", 0, 10, domainerrors.ErrInvalidQuantity},
		{"quantidade negativa", -1, 10, domainerrors.ErrInvalidQuantity},
		{"estoque insuficiente", 20, 10, domainerrors.ErrInsufficientStock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := product.UpdatedAt()
			err := product.DecreaseStock(tt.quantity)

			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if err == nil {
				if product.Stock() != tt.want {
					t.Errorf("Stock() got = %d, want %d", product.Stock(), tt.want)
				}
				if product.UpdatedAt() == oldUpdatedAt {
					t.Errorf("expected UpdatedAt to change")
				}
			}
		})
	}
}

func TestProduct_HasStock(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	product, _ := NewProduct("Produto", "Descrição", price, 10)

	tests := []struct {
		name     string
		quantity int
		want     bool
	}{
		{"tem estoque exato", 10, true},
		{"tem estoque maior", 5, true},
		{"não tem estoque", 15, false},
		{"quantidade zero", 0, true},
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

func TestProduct_ActivateDeactivate(t *testing.T) {
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	product, _ := NewProduct("Produto", "Descrição", price, 10)

	// Desativar
	oldUpdatedAt := product.UpdatedAt()
	product.Deactivate()
	if product.IsActive() {
		t.Errorf("expected product to be inactive")
	}
	if product.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	// Ativar
	oldUpdatedAt = product.UpdatedAt()
	product.Activate()
	if !product.IsActive() {
		t.Errorf("expected product to be active")
	}
	if product.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}
}

func TestProduct_Rebuild(t *testing.T) {
	id := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	now := time.Now()

	product := RebuildProduct(id, "Produto", "Descrição", price, 10, true, now, now)

	if product.ID() != id {
		t.Errorf("ID() got = %v, want %v", product.ID(), id)
	}
	if product.Name() != "Produto" {
		t.Errorf("Name() got = %q, want %q", product.Name(), "Produto")
	}
	if product.Description() != "Descrição" {
		t.Errorf("Description() got = %q, want %q", product.Description(), "Descrição")
	}
	if !product.Price().Equals(price) {
		t.Errorf("Price() got = %v, want %v", product.Price().Amount(), price.Amount())
	}
	if product.Stock() != 10 {
		t.Errorf("Stock() got = %d, want 10", product.Stock())
	}
	if !product.IsActive() {
		t.Errorf("expected product to be active")
	}
	if product.CreatedAt() != now {
		t.Errorf("CreatedAt() got = %v, want %v", product.CreatedAt(), now)
	}
	if product.UpdatedAt() != now {
		t.Errorf("UpdatedAt() got = %v, want %v", product.UpdatedAt(), now)
	}
}
