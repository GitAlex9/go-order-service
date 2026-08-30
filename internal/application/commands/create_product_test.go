package commands

import (
	"context"
	"testing"

	"errors"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/validation"
)

func TestCreateProductHandler_Handle(t *testing.T) {
	t.Run("cria produto com dados válidos", func(t *testing.T) {
		repo := newFakeProductRepository()
		handler := NewCreateProductHandler(repo)

		resp, err := handler.Handle(context.Background(), dto.CreateProductRequest{
			Name: "Notebook", Description: "Notebook gamer", Price: 3500.00, Stock: 10,
		})

		if err != nil {
			t.Fatalf("Handle() error = %v, want nil", err)
		}
		if resp.Stock != 10 {
			t.Errorf("Stock got = %d, want %d", resp.Stock, 10)
		}
	})

	t.Run("retorna erro de validação com dados inválidos", func(t *testing.T) {
		repo := newFakeProductRepository()
		handler := NewCreateProductHandler(repo)

		_, err := handler.Handle(context.Background(), dto.CreateProductRequest{
			Name: "", Description: "", Price: -1, Stock: -1,
		})

		var verr *validation.ValidationErrors
		if !errors.As(err, &verr) {
			t.Fatalf("Handle() error = %v, want *validation.ValidationErrors", err)
		}
	})
}
