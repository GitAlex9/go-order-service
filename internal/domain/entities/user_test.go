package entities

import (
	"errors"
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

const validPassword = "SenhaForte123!"

func TestNewUser(t *testing.T) {
	email, err := valueobjects.NewEmail("teste@teste.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		name     string
		password string
		role     Role
		wantErr  error
	}{
		{"dados válidos", validPassword, RoleAdmin, nil},
		{"role inválida", validPassword, Role("desconhecida"), domainerrors.ErrInvalidRole},
		{"senha fraca", "123", RoleCustomer, domainerrors.ErrWeakPassword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(email, tt.password, tt.role)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewUser() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil && got != nil {
				t.Errorf("got user = %v, want nil", got)
			}
			if tt.wantErr == nil && got == nil {
				t.Errorf("got nil user, want non-nil")
			}
		})
	}
}

func TestNewUser_FieldsArePersisted(t *testing.T) {
	email, err := valueobjects.NewEmail("teste@teste.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	user, err := NewUser(email, validPassword, RoleAdmin)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if got := user.ID(); got == uuid.Nil {
		t.Errorf("ID() got = %v, want non-nil uuid", got)
	}
	if got, want := user.Email().String(), "teste@teste.com"; got != want {
		t.Errorf("Email().String() got = %q, want %q", got, want)
	}
	if got, want := user.Role(), RoleAdmin; got != want {
		t.Errorf("Role() got = %v, want %v", got, want)
	}
	if got := user.Active(); !got {
		t.Errorf("Active() got = %v, want true", got)
	}
	if got := user.PasswordHash(); got == "" {
		t.Errorf("PasswordHash() got = empty, want a bcrypt hash")
	}
	if got := user.PasswordHash(); got == validPassword {
		t.Errorf("PasswordHash() should never equal the plain text password")
	}
}

func TestRestoreUser(t *testing.T) {
	id := uuid.New()
	email, err := valueobjects.NewEmail("teste@teste.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	original, err := NewUser(email, validPassword, RoleAdmin)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	restored := RestoreUser(id, email, original.PasswordHash(), RoleAdmin, false, original.CreatedAt(), original.UpdatedAt())

	if got := restored.ID(); got != id {
		t.Errorf("ID() got = %v, want %v", got, id)
	}
	if got := restored.Active(); got {
		t.Errorf("Active() got = %v, want false", got)
	}
	if !restored.CheckPassword(validPassword) {
		t.Errorf("expected restored user to match the original plain text password via its hash")
	}
}

func TestUser_CheckPassword(t *testing.T) {
	email, err := valueobjects.NewEmail("teste@teste.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	user, err := NewUser(email, validPassword, RoleAdmin)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		name  string
		plain string
		want  bool
	}{
		{"senha correta", validPassword, true},
		{"senha incorreta", "SenhaErrada456!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := user.CheckPassword(tt.plain)
			if got != tt.want {
				t.Errorf("CheckPassword(%q) got = %v, want %v", tt.plain, got, tt.want)
			}
		})
	}
}

func TestUser_ChangePassword(t *testing.T) {
	t.Run("senha atual correta permite a troca", func(t *testing.T) {
		email, err := valueobjects.NewEmail("teste@teste.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		user, err := NewUser(email, validPassword, RoleAdmin)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		newPassword := "NovaSenhaForte456!"

		if err := user.ChangePassword(validPassword, newPassword); err != nil {
			t.Fatalf("ChangePassword() error = %v, want nil", err)
		}

		if !user.CheckPassword(newPassword) {
			t.Errorf("expected new password to be active after change")
		}
		if user.CheckPassword(validPassword) {
			t.Errorf("expected old password to no longer work after change")
		}

		userEvents := user.Events()
		if len(userEvents) == 0 {
			t.Fatalf("expected at least one event")
		}
		if _, ok := userEvents[len(userEvents)-1].(events.UserPasswordChangedEvent); !ok {
			t.Errorf("expected UserPasswordChangedEvent, got %T", userEvents[len(userEvents)-1])
		}
	})

	t.Run("senha atual incorreta impede a troca", func(t *testing.T) {
		email, err := valueobjects.NewEmail("teste@teste.com")
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		user, err := NewUser(email, validPassword, RoleAdmin)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = user.ChangePassword("SenhaErradaQualquer1!", "NovaSenhaForte456!")

		if !errors.Is(err, domainerrors.ErrIncorrectCurrentPassword) {
			t.Errorf("ChangePassword() error = %v, want %v", err, domainerrors.ErrIncorrectCurrentPassword)
		}
		if !user.CheckPassword(validPassword) {
			t.Errorf("expected original password to remain unchanged after failed attempt")
		}
	})
}

func TestUser_ChangeEmail(t *testing.T) {
	email, err := valueobjects.NewEmail("old@teste.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	user, err := NewUser(email, validPassword, RoleAdmin)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	newEmail, err := valueobjects.NewEmail("new@teste.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	oldUpdatedAt := user.UpdatedAt()
	user.ChangeEmail(newEmail)

	if got, want := user.Email().String(), "new@teste.com"; got != want {
		t.Errorf("Email().String() got = %q, want %q", got, want)
	}
	if got := user.UpdatedAt(); !got.After(oldUpdatedAt) {
		t.Errorf("UpdatedAt() got = %v, want a timestamp after %v", got, oldUpdatedAt)
	}

	userEvents := user.Events()
	if len(userEvents) == 0 {
		t.Fatalf("expected at least one event")
	}
	evt, ok := userEvents[len(userEvents)-1].(events.UserEmailChangedEvent)
	if !ok {
		t.Fatalf("expected UserEmailChangedEvent, got %T", userEvents[len(userEvents)-1])
	}
	if evt.OldEmail.String() != "old@teste.com" {
		t.Errorf("OldEmail got = %q, want %q", evt.OldEmail.String(), "old@teste.com")
	}
	if evt.NewEmail.String() != "new@teste.com" {
		t.Errorf("NewEmail got = %q, want %q", evt.NewEmail.String(), "new@teste.com")
	}
}

func TestUser_ActivateDeactivate(t *testing.T) {
	email, err := valueobjects.NewEmail("teste@teste.com")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	user, err := NewUser(email, validPassword, RoleAdmin)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	user.Deactivate()
	if got := user.Active(); got {
		t.Errorf("Active() after Deactivate() got = %v, want false", got)
	}
	lastEvent := user.Events()[len(user.Events())-1]
	if _, ok := lastEvent.(events.UserDeactivatedEvent); !ok {
		t.Errorf("expected UserDeactivatedEvent, got %T", lastEvent)
	}

	user.Activate()
	if got := user.Active(); !got {
		t.Errorf("Active() after Activate() got = %v, want true", got)
	}
	lastEvent = user.Events()[len(user.Events())-1]
	if _, ok := lastEvent.(events.UserActivatedEvent); !ok {
		t.Errorf("expected UserActivatedEvent, got %T", lastEvent)
	}
}
