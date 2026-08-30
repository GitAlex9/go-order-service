package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestDeleteCustomerHandler_Handle(t *testing.T) {
	repo := newFakeCustomerRepository()

	created, err := NewCreateCustomerHandler(repo).Handle(context.Background(), dto.CreateCustomerRequest{
		Name: "Cliente Teste", Email: "cliente@teste.com", CPF: "52998224725",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	id, err := parseTestUUID(created.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	handler := NewDeleteCustomerHandler(repo)

	if err := handler.Handle(context.Background(), id); err != nil {
		t.Fatalf("Handle() error = %v, want nil", err)
	}

	_, err = repo.FindByID(context.Background(), id)
	if !errors.Is(err, domainerrors.ErrNotFound) {
		t.Errorf("expected customer to be deleted, FindByID error = %v, want %v", err, domainerrors.ErrNotFound)
	}
}
