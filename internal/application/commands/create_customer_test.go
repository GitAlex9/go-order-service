package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/validation"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestCreateCustomerHandler_Handle(t *testing.T) {
	t.Run("cria cliente com dados válidos", func(t *testing.T) {
		repo := newFakeCustomerRepository()
		handler := NewCreateCustomerHandler(repo)

		resp, err := handler.Handle(context.Background(), dto.CreateCustomerRequest{
			Name: "Cliente Teste", Email: "cliente@teste.com", CPF: "52998224725",
		})

		if err != nil {
			t.Fatalf("Handle() error = %v, want nil", err)
		}
		if resp == nil {
			t.Fatal("Handle() returned nil response")
		}
		if resp.Name != "Cliente Teste" {
			t.Errorf("Name got = %q, want %q", resp.Name, "Cliente Teste")
		}
	})

	t.Run("retorna erro de validação agregado com dados inválidos", func(t *testing.T) {
		repo := newFakeCustomerRepository()
		handler := NewCreateCustomerHandler(repo)

		_, err := handler.Handle(context.Background(), dto.CreateCustomerRequest{
			Name: "ab", Email: "email-invalido", CPF: "123",
		})

		var verr *validation.ValidationErrors
		if !errors.As(err, &verr) {
			t.Fatalf("Handle() error = %v, want *validation.ValidationErrors", err)
		}
		if len(verr.Errors) != 3 {
			t.Errorf("got %d validation errors, want 3", len(verr.Errors))
		}
	})

	t.Run("retorna erro de duplicidade quando email já existe", func(t *testing.T) {
		repo := newFakeCustomerRepository()
		handler := NewCreateCustomerHandler(repo)

		first := dto.CreateCustomerRequest{Name: "Cliente Um", Email: "duplicado@teste.com", CPF: "52998224725"}
		if _, err := handler.Handle(context.Background(), first); err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		second := dto.CreateCustomerRequest{Name: "Cliente Dois", Email: "duplicado@teste.com", CPF: "11144477735"}
		_, err := handler.Handle(context.Background(), second)

		if !errors.Is(err, domainerrors.ErrDuplicateCustomer) {
			t.Errorf("Handle() error = %v, want %v", err, domainerrors.ErrDuplicateCustomer)
		}
	})
}
