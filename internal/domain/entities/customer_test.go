package entities

import (
	"testing"
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestNewCustomer_Valid(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")
	cpf, _ := valueobjects.NewCPF("52998224725")

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"nome normal", "João Silva", "João Silva"},
		{"nome com espaços", "  Maria Souza  ", "Maria Souza"},
		{"nome mínimo", "Ana", "Ana"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := NewCustomer(tt.arg, email, cpf)
			if err != nil {
				t.Fatalf("expected creation to succeed, got error: %v", err)
			}
			if customer.Name() != tt.want {
				t.Errorf("Name() got = %q, want %q", customer.Name(), tt.want)
			}
			if customer.ID() == uuid.Nil {
				t.Errorf("expected ID to be set")
			}
			if customer.UserID() != nil {
				t.Errorf("expected UserID to be nil")
			}
			if customer.CreatedAt().IsZero() || customer.UpdatedAt().IsZero() {
				t.Errorf("expected timestamps to be set")
			}
		})
	}
}

func TestNewCustomer_InvalidName(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")
	cpf, _ := valueobjects.NewCPF("52998224725")

	tests := []struct {
		name string
		arg  string
	}{
		{"nome vazio", ""},
		{"nome com apenas espaços", "   "},
		{"nome curto", "ab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := NewCustomer(tt.arg, email, cpf)
			if err != domainerrors.ErrEmptyName {
				t.Errorf("expected ErrEmptyName, got %v", err)
			}
			if customer != nil {
				t.Errorf("expected nil customer")
			}
		})
	}
}

func TestCustomer_Rename(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")
	cpf, _ := valueobjects.NewCPF("52998224725")

	tests := []struct {
		name    string
		newName string
		wantErr bool
		want    string
	}{
		{"renomear válido", "Novo Nome", false, "Novo Nome"},
		{"renomear com espaços", "  Outro Nome  ", false, "Outro Nome"},
		{"renomear vazio", "", true, "Nome Antigo"},
		{"renomear curto", "ab", true, "Nome Antigo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := NewCustomer("Nome Antigo", email, cpf)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			oldUpdatedAt := customer.UpdatedAt()
			err = customer.Rename(tt.newName)

			if tt.wantErr {
				if err != domainerrors.ErrEmptyName {
					t.Errorf("expected ErrEmptyName, got %v", err)
				}
				if customer.Name() != tt.want {
					t.Errorf("Name() got = %q, want %q", customer.Name(), tt.want)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if customer.Name() != tt.want {
				t.Errorf("Name() got = %q, want %q", customer.Name(), tt.want)
			}
			if customer.UpdatedAt() == oldUpdatedAt {
				t.Errorf("expected UpdatedAt to change")
			}

			customerEvents := customer.Events()
			if len(customerEvents) == 0 {
				t.Errorf("expected at least one event")
			}
			_, ok := customerEvents[len(customerEvents)-1].(events.CustomerRenamedEvent)
			if !ok {
				t.Errorf("expected CustomerRenamedEvent")
			}
		})
	}
}

func TestCustomer_ChangeEmail(t *testing.T) {
	email, _ := valueobjects.NewEmail("old@teste.com")
	cpf, _ := valueobjects.NewCPF("52998224725")
	customer, _ := NewCustomer("Cliente", email, cpf)
	userID := uuid.New()
	customer.LinkUser(userID)

	newEmail, _ := valueobjects.NewEmail("new@teste.com")
	oldUpdatedAt := customer.UpdatedAt()
	customer.ChangeEmail(newEmail)

	if customer.Email().String() != "new@teste.com" {
		t.Errorf("Email() got = %q, want %q", customer.Email().String(), "new@teste.com")
	}
	if customer.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	emailEvents := customer.Events()
	lastEvent := emailEvents[len(emailEvents)-1]
	userEmailEvent, ok := lastEvent.(events.UserEmailChangedEvent)
	if !ok {
		t.Errorf("expected UserEmailChangedEvent, got %T", lastEvent)
	}
	if userEmailEvent.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, userEmailEvent.UserID)
	}
	if userEmailEvent.OldEmail.String() != "old@teste.com" {
		t.Errorf("expected OldEmail old@teste.com, got %v", userEmailEvent.OldEmail)
	}
}

func TestCustomer_LinkUser(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")
	cpf, _ := valueobjects.NewCPF("52998224725")
	customer, _ := NewCustomer("Cliente", email, cpf)

	userID := uuid.New()
	oldUpdatedAt := customer.UpdatedAt()
	customer.LinkUser(userID)

	if customer.UserID() == nil {
		t.Errorf("expected UserID to be set")
	}
	if *customer.UserID() != userID {
		t.Errorf("UserID() got = %v, want %v", *customer.UserID(), userID)
	}
	if customer.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	// Verifica o evento
	customerEvents := customer.Events()
	lastEvent := customerEvents[len(customerEvents)-1]
	linkEvent, ok := lastEvent.(events.CustomerLinkedToUserEvent)
	if !ok {
		t.Errorf("expected CustomerLinkedToUserEvent, got %T", lastEvent)
	}
	if linkEvent.CustomerID != customer.ID() {
		t.Errorf("expected CustomerID %v, got %v", customer.ID(), linkEvent.CustomerID)
	}
	if linkEvent.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, linkEvent.UserID)
	}
}

func TestCustomer_Getters(t *testing.T) {
	id := uuid.New()
	email, _ := valueobjects.NewEmail("teste@teste.com")
	cpf, _ := valueobjects.NewCPF("52998224725")
	userID := uuid.New()
	now := time.Now()

	customer := RebuildCustomer(id, "Teste", email, cpf, &userID, now, now)

	if customer.ID() != id {
		t.Errorf("ID() got = %v, want %v", customer.ID(), id)
	}
	if customer.Name() != "Teste" {
		t.Errorf("Name() got = %q, want %q", customer.Name(), "Teste")
	}
	if customer.Email().String() != "teste@teste.com" {
		t.Errorf("Email() got = %q, want %q", customer.Email().String(), "teste@teste.com")
	}
	if customer.CPF().String() != "52998224725" {
		t.Errorf("CPF() got = %q, want %q", customer.CPF().String(), "52998224725")
	}
	if customer.UserID() == nil || *customer.UserID() != userID {
		t.Errorf("UserID() got = %v, want %v", customer.UserID(), &userID)
	}
	if customer.CreatedAt() != now {
		t.Errorf("CreatedAt() got = %v, want %v", customer.CreatedAt(), now)
	}
	if customer.UpdatedAt() != now {
		t.Errorf("UpdatedAt() got = %v, want %v", customer.UpdatedAt(), now)
	}
}
