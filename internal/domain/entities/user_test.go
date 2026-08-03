package entities

import (
	"testing"
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestNewUser_Valid(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")

	tests := []struct {
		name     string
		email    valueobjects.Email
		password string
		role     Role
	}{
		{"usuário customer", email, "Teste@123456", RoleCustomer},
		{"usuário admin", email, "Teste@123456", RoleAdmin},
		{"usuário manager", email, "Teste@123456", RoleManager},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.password, tt.role)
			if err != nil {
				t.Fatalf("expected creation to succeed, got error: %v", err)
			}
			if user.ID() == uuid.Nil {
				t.Errorf("expected ID to be set")
			}
			if user.Email().String() != tt.email.String() {
				t.Errorf("Email() got = %q, want %q", user.Email().String(), tt.email.String())
			}
			if user.Role() != tt.role {
				t.Errorf("Role() got = %v, want %v", user.Role(), tt.role)
			}
			if !user.Active() {
				t.Errorf("expected user to be active")
			}
			if user.PasswordHash() == "" {
				t.Errorf("expected password hash to be set")
			}
			if user.CreatedAt().IsZero() || user.UpdatedAt().IsZero() {
				t.Errorf("expected timestamps to be set")
			}
		})
	}
}

func TestNewUser_InvalidRole(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")
	invalidRole := Role("superuser")

	user, err := NewUser(email, "Teste@123456", invalidRole)
	if err != domainerrors.ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
	if user != nil {
		t.Errorf("expected nil user")
	}
}

func TestNewUser_WeakPassword(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"senha curta", "Teste@1", true},
		{"sem maiúscula", "teste@123456", true},
		{"sem minúscula", "TESTE@123456", true},
		{"sem número", "Teste@abcdef", true},
		{"sem símbolo", "Teste123456", true},
		{"senha forte", "Teste@123456", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(email, tt.password, RoleCustomer)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if user != nil {
					t.Errorf("expected nil user")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if user == nil {
				t.Errorf("expected user to be created")
			}
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")
	user, _ := NewUser(email, "Teste@123456", RoleCustomer)

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"senha correta", "Teste@123456", true},
		{"senha incorreta", "Senha@Errada", false},
		{"senha vazia", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := user.CheckPassword(tt.password)
			if got != tt.want {
				t.Errorf("CheckPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_ChangePassword(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")
	user, _ := NewUser(email, "Current@123456", RoleCustomer)

	tests := []struct {
		name         string
		currentPlain string
		newPlain     string
		wantErr      bool
	}{
		{"senha correta", "Current@123456", "New@123456", false},
		{"senha atual incorreta", "Wrong@123456", "New@123456", true},
		{"nova senha fraca", "Current@123456", "weak", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := user.UpdatedAt()
			err := user.ChangePassword(tt.currentPlain, tt.newPlain)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if user.UpdatedAt() == oldUpdatedAt {
				t.Errorf("expected UpdatedAt to change")
			}
			if !user.CheckPassword(tt.newPlain) {
				t.Errorf("new password does not match")
			}
			if user.CheckPassword(tt.currentPlain) {
				t.Errorf("old password should not match anymore")
			}

			// Verifica evento
			userEvents := user.Events()
			lastEvent := userEvents[len(userEvents)-1]
			passwordEvent, ok := lastEvent.(events.UserPasswordChangedEvent)
			if !ok {
				t.Errorf("expected UserPasswordChangedEvent, got %T", lastEvent)
			}
			if passwordEvent.UserID != user.ID() {
				t.Errorf("expected UserID %v, got %v", user.ID(), passwordEvent.UserID)
			}
		})
	}
}

func TestUser_ChangeEmail(t *testing.T) {
	email, _ := valueobjects.NewEmail("old@teste.com")
	user, _ := NewUser(email, "Teste@123456", RoleCustomer)

	newEmail, _ := valueobjects.NewEmail("new@teste.com")
	oldUpdatedAt := user.UpdatedAt()
	user.ChangeEmail(newEmail)

	if user.Email().String() != "new@teste.com" {
		t.Errorf("Email() got = %q, want %q", user.Email().String(), "new@teste.com")
	}
	if user.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	// Verifica evento
	userEvents := user.Events()
	lastEvent := userEvents[len(userEvents)-1]
	emailEvent, ok := lastEvent.(events.UserEmailChangedEvent)
	if !ok {
		t.Errorf("expected UserEmailChangedEvent, got %T", lastEvent)
	}
	if emailEvent.UserID != user.ID() {
		t.Errorf("expected UserID %v, got %v", user.ID(), emailEvent.UserID)
	}
	if emailEvent.OldEmail.String() != "old@teste.com" {
		t.Errorf("expected OldEmail old@teste.com, got %v", emailEvent.OldEmail)
	}
	if emailEvent.NewEmail.String() != "new@teste.com" {
		t.Errorf("expected NewEmail new@teste.com, got %v", emailEvent.NewEmail)
	}
}

func TestUser_ActivateDeactivate(t *testing.T) {
	email, _ := valueobjects.NewEmail("teste@teste.com")
	user, _ := NewUser(email, "Teste@123456", RoleCustomer)

	// Desativar
	oldUpdatedAt := user.UpdatedAt()
	user.Deactivate()
	if user.Active() {
		t.Errorf("expected user to be inactive")
	}
	if user.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	// Verifica evento de desativação
	userEvents := user.Events()
	lastEvent := userEvents[len(userEvents)-1]
	deactivateEvent, ok := lastEvent.(events.UserDeactivatedEvent)
	if !ok {
		t.Errorf("expected UserDeactivatedEvent, got %T", lastEvent)
	}
	if deactivateEvent.UserID != user.ID() {
		t.Errorf("expected UserID %v, got %v", user.ID(), deactivateEvent.UserID)
	}

	// Ativar
	oldUpdatedAt = user.UpdatedAt()
	user.Activate()
	if !user.Active() {
		t.Errorf("expected user to be active")
	}
	if user.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	// Verifica evento de ativação
	userEvents = user.Events()
	lastEvent = userEvents[len(userEvents)-1]
	activateEvent, ok := lastEvent.(events.UserActivatedEvent)
	if !ok {
		t.Errorf("expected UserActivatedEvent, got %T", lastEvent)
	}
	if activateEvent.UserID != user.ID() {
		t.Errorf("expected UserID %v, got %v", user.ID(), activateEvent.UserID)
	}
}

func TestUser_Restore(t *testing.T) {
	id := uuid.New()
	email, _ := valueobjects.NewEmail("teste@teste.com")
	hash := "$2a$10$N9qo8uLOickgx2ZMRZoMy.Mr/.Zwp6sJ9N7kCjI5M1fBfZ5qVwLqS" // hash fictício
	now := time.Now()

	user := RestoreUser(id, email, hash, RoleAdmin, true, now, now)

	if user.ID() != id {
		t.Errorf("ID() got = %v, want %v", user.ID(), id)
	}
	if user.Email().String() != "teste@teste.com" {
		t.Errorf("Email() got = %q, want %q", user.Email().String(), "teste@teste.com")
	}
	if user.PasswordHash() != hash {
		t.Errorf("PasswordHash() got = %q, want %q", user.PasswordHash(), hash)
	}
	if user.Role() != RoleAdmin {
		t.Errorf("Role() got = %v, want %v", user.Role(), RoleAdmin)
	}
	if !user.Active() {
		t.Errorf("expected user to be active")
	}
	if user.CreatedAt() != now {
		t.Errorf("CreatedAt() got = %v, want %v", user.CreatedAt(), now)
	}
	if user.UpdatedAt() != now {
		t.Errorf("UpdatedAt() got = %v, want %v", user.UpdatedAt(), now)
	}
}
