package commands

import (
	"context"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/google/uuid"
)

func TestUpdateCustomerHandler_Handle(t *testing.T) {
	repo := newFakeCustomerRepository()
	dispatcher := &fakeDispatcher{}

	created, err := NewCreateCustomerHandler(repo).Handle(context.Background(), dto.CreateCustomerRequest{
		Name: "Nome Antigo", Email: "cliente@teste.com", CPF: "52998224725",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	handler := NewUpdateCustomerHandler(repo, dispatcher)

	id, err := parseTestUUID(created.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	resp, err := handler.Handle(context.Background(), id, dto.UpdateCustomerRequest{Name: "Nome Novo"})
	if err != nil {
		t.Fatalf("Handle() error = %v, want nil", err)
	}
	if resp.Name != "Nome Novo" {
		t.Errorf("Name got = %q, want %q", resp.Name, "Nome Novo")
	}
	if len(dispatcher.dispatched) == 0 {
		t.Errorf("expected at least one event to be dispatched after rename")
	}
}

func parseTestUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
