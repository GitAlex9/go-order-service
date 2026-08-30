package entities

import (
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestNewOrderItem_Valid(t *testing.T) {
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)

	tests := []struct {
		name         string
		productID    uuid.UUID
		productName  string
		unitPrice    valueobjects.Money
		quantity     int
		wantSubtotal float64
	}{
		{"item normal", productID, "Produto Teste", price, 2, 21.00},
		{"quantidade 1", productID, "Produto Teste", price, 1, 10.50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewOrderItem(tt.productID, tt.productName, tt.unitPrice, tt.quantity)
			if err != nil {
				t.Fatalf("expected creation to succeed, got error: %v", err)
			}
			if item.ProductID() != tt.productID {
				t.Errorf("ProductID() got = %v, want %v", item.ProductID(), tt.productID)
			}
			if item.ProductName() != tt.productName {
				t.Errorf("ProductName() got = %q, want %q", item.ProductName(), tt.productName)
			}
			if !item.UnitPrice().Equals(tt.unitPrice) {
				t.Errorf("UnitPrice() got = %v, want %v", item.UnitPrice().Amount(), tt.unitPrice.Amount())
			}
			if item.Quantity() != tt.quantity {
				t.Errorf("Quantity() got = %d, want %d", item.Quantity(), tt.quantity)
			}
			if item.Subtotal().Amount() != tt.wantSubtotal {
				t.Errorf("Subtotal() got = %v, want %v", item.Subtotal().Amount(), tt.wantSubtotal)
			}
		})
	}
}

func TestNewOrderItem_Invalid(t *testing.T) {
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)

	tests := []struct {
		name        string
		productID   uuid.UUID
		productName string
		unitPrice   valueobjects.Money
		quantity    int
		wantErr     error
	}{
		{"productID vazio", uuid.Nil, "Produto", price, 2, domainerrors.ErrInvalidProductID},
		{"nome vazio", productID, "", price, 2, domainerrors.ErrInvalidProductName},
		{"preço zero", productID, "Produto", valueobjects.Zero(), 2, domainerrors.ErrInvalidProductPrice},
		{"quantidade zero", productID, "Produto", price, 0, domainerrors.ErrInvalidQuantity},
		{"quantidade negativa", productID, "Produto", price, -1, domainerrors.ErrInvalidQuantity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewOrderItem(tt.productID, tt.productName, tt.unitPrice, tt.quantity)
			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if item != nil {
				t.Errorf("expected nil item")
			}
		})
	}
}

func TestOrderItem_Rebuild(t *testing.T) {
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)

	item := RebuildOrderItem(productID, "Produto Teste", price, 3)

	if item.ProductID() != productID {
		t.Errorf("ProductID() got = %v, want %v", item.ProductID(), productID)
	}
	if item.ProductName() != "Produto Teste" {
		t.Errorf("ProductName() got = %q, want %q", item.ProductName(), "Produto Teste")
	}
	if !item.UnitPrice().Equals(price) {
		t.Errorf("UnitPrice() got = %v, want %v", item.UnitPrice().Amount(), price.Amount())
	}
	if item.Quantity() != 3 {
		t.Errorf("Quantity() got = %d, want 3", item.Quantity())
	}
}

func TestOrderItem_Subtotal(t *testing.T) {
	productID := uuid.New()

	tests := []struct {
		name     string
		price    float64
		quantity int
		want     float64
	}{
		{"2 x 10.50 = 21.00", 10.50, 2, 21.00},
		{"3 x 5.25 = 15.75", 5.25, 3, 15.75},
		{"1 x 0.99 = 0.99", 0.99, 1, 0.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := valueobjects.NewMoneyFromFloat(tt.price)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			item, err := NewOrderItem(productID, "Produto", price, tt.quantity)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := item.Subtotal().Amount()
			if got != tt.want {
				t.Errorf("Subtotal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
