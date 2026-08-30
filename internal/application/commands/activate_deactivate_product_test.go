package commands

import (
	"context"
	"testing"
)

func TestActivateDeactivateProductHandler_Handle(t *testing.T) {
	repo := newFakeProductRepository()
	productID := setupTestProduct(t, repo, 10)

	id, err := parseTestUUID(productID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	deactivateHandler := NewDeactivateProductHandler(repo)
	resp, err := deactivateHandler.Handle(context.Background(), id)
	if err != nil {
		t.Fatalf("Deactivate Handle() error = %v, want nil", err)
	}
	if resp.Active {
		t.Errorf("Active got = true, want false after deactivate")
	}

	activateHandler := NewActivateProductHandler(repo)
	resp, err = activateHandler.Handle(context.Background(), id)
	if err != nil {
		t.Fatalf("Activate Handle() error = %v, want nil", err)
	}
	if !resp.Active {
		t.Errorf("Active got = false, want true after activate")
	}
}
