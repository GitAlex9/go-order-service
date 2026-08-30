package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func setupTestProduct(t *testing.T, repo *fakeProductRepository, stock int) string {
	t.Helper()

	created, err := NewCreateProductHandler(repo).Handle(context.Background(), dto.CreateProductRequest{
		Name: "Notebook", Description: "Notebook gamer", Price: 3500.00, Stock: stock,
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	return created.ID
}

func TestIncreaseStockHandler_Handle(t *testing.T) {
	repo := newFakeProductRepository()
	dispatcher := &fakeDispatcher{}
	productID := setupTestProduct(t, repo, 10)

	handler := NewIncreaseStockHandler(repo, dispatcher)
	id, err := parseTestUUID(productID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	resp, err := handler.Handle(context.Background(), id, dto.AdjustStockRequest{Quantity: 5})
	if err != nil {
		t.Fatalf("Handle() error = %v, want nil", err)
	}
	if resp.Stock != 15 {
		t.Errorf("Stock got = %d, want %d", resp.Stock, 15)
	}
	if len(dispatcher.dispatched) == 0 {
		t.Errorf("expected stock increased event to be dispatched")
	}
}

func TestDecreaseStockHandler_Handle(t *testing.T) {
	t.Run("quantidade válida", func(t *testing.T) {
		repo := newFakeProductRepository()
		dispatcher := &fakeDispatcher{}
		productID := setupTestProduct(t, repo, 10)

		handler := NewDecreaseStockHandler(repo, dispatcher)
		id, err := parseTestUUID(productID)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		resp, err := handler.Handle(context.Background(), id, dto.AdjustStockRequest{Quantity: 4})
		if err != nil {
			t.Fatalf("Handle() error = %v, want nil", err)
		}
		if resp.Stock != 6 {
			t.Errorf("Stock got = %d, want %d", resp.Stock, 6)
		}
	})

	t.Run("estoque insuficiente", func(t *testing.T) {
		repo := newFakeProductRepository()
		dispatcher := &fakeDispatcher{}
		productID := setupTestProduct(t, repo, 3)

		handler := NewDecreaseStockHandler(repo, dispatcher)
		id, err := parseTestUUID(productID)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		_, err = handler.Handle(context.Background(), id, dto.AdjustStockRequest{Quantity: 999})

		if !errors.Is(err, domainerrors.ErrInsufficientStock) {
			t.Errorf("Handle() error = %v, want %v", err, domainerrors.ErrInsufficientStock)
		}
	})
}
