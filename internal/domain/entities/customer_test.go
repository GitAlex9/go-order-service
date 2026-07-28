package entities

import (
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestNewCustomer_WithValidData(t *testing.T) {
	// Arrange
	email, err := valueobjects.NewEmail("cliente@teste.com")
	if err != nil {
		t.Fatalf("expected valid email, got error: %v", err)
	}

	cpf, err := valueobjects.NewCPF("52998224725")
	if err != nil {
		t.Fatalf("expected valid cpf, got error: %v", err)
	}

	// Act
	customer, err := NewCustomer("Cliente Teste", email, cpf)

	// Assert
	if err != nil {
		t.Fatalf("expected customer creation to succeed, got error: %v", err)
	}

	if customer.ID() == uuid.Nil {
		t.Fatalf("expected customer id to be set")
	}

	if customer.Name() != "Cliente Teste" {
		t.Fatalf("expected name to be preserved, got %q", customer.Name())
	}

	if customer.Email().String() != "cliente@teste.com" {
		t.Fatalf("expected email to be preserved, got %q", customer.Email().String())
	}

	if customer.CPF().String() != "52998224725" {
		t.Fatalf("expected cpf to be preserved, got %q", customer.CPF().String())
	}

	if customer.UserID() != nil {
		t.Fatalf("expected user id to be nil before linking")
	}

	if customer.CreatedAt().IsZero() || customer.UpdatedAt().IsZero() {
		t.Fatalf("expected creation timestamps to be set")
	}
}

func TestNewCustomer_WithEmptyName(t *testing.T) {
	// Arrange
	email, err := valueobjects.NewEmail("cliente@teste.com")
	if err != nil {
		t.Fatalf("expected valid email, got error: %v", err)
	}

	cpf, err := valueobjects.NewCPF("52998224725")
	if err != nil {
		t.Fatalf("expected valid cpf, got error: %v", err)
	}

	// Act
	customer, err := NewCustomer("   ", email, cpf)

	// Assert
	if err != domainerrors.ErrEmptyName {
		t.Fatalf("expected ErrEmptyName, got %v", err)
	}

	if customer != nil {
		t.Fatalf("expected nil customer when name is empty")
	}
}

func TestCustomer_RenameChangeEmailAndLinkUser(t *testing.T) {
	// Arrange
	email, err := valueobjects.NewEmail("cliente@teste.com")
	if err != nil {
		t.Fatalf("expected valid email, got error: %v", err)
	}

	cpf, err := valueobjects.NewCPF("52998224725")
	if err != nil {
		t.Fatalf("expected valid cpf, got error: %v", err)
	}

	customer, err := NewCustomer("Cliente", email, cpf)
	if err != nil {
		t.Fatalf("expected customer creation to succeed, got error: %v", err)
	}

	// Act
	if err := customer.Rename("Cliente Atualizado"); err != nil {
		t.Fatalf("expected rename to succeed, got error: %v", err)
	}

	newEmail, err := valueobjects.NewEmail("novo@teste.com")
	if err != nil {
		t.Fatalf("expected new email to be valid, got error: %v", err)
	}

	customer.ChangeEmail(newEmail)

	userID := customer.ID()
	customer.LinkUser(userID)

	// Assert
	if customer.Name() != "Cliente Atualizado" {
		t.Fatalf("expected renamed customer name, got %q", customer.Name())
	}

	if customer.Email().String() != "novo@teste.com" {
		t.Fatalf("expected email to be updated, got %q", customer.Email().String())
	}

	if customer.UserID() == nil {
		t.Fatalf("expected user id to be linked")
	}

	if *customer.UserID() != userID {
		t.Fatalf("expected linked user id to match customer id")
	}
}
