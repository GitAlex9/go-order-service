package valueobjects

import (
	"errors"
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		plain   string
		wantErr error
	}{
		{"senha forte válida", "SenhaForte123!", nil},
		{"senha curta", "Ab1!", domainerrors.ErrWeakPassword},
		{"senha sem maiúscula", "senhaforte123!", domainerrors.ErrPasswordNoUpper},
		{"senha sem minúscula", "SENHAFORTE123!", domainerrors.ErrPasswordNoLower},
		{"senha sem número", "SenhaForte!!!", domainerrors.ErrPasswordNoNumber},
		{"senha sem símbolo", "SenhaForte123", domainerrors.ErrPasswordNoSpecial},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.plain)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewPassword(%q) error = %v, want error containing %v", tt.plain, err, tt.wantErr)
			}
		})
	}
}

func TestPassword_Matches(t *testing.T) {
	password, err := NewPassword("SenhaForte123!")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		name  string
		plain string
		want  bool
	}{
		{"senha correta", "SenhaForte123!", true},
		{"senha incorreta", "SenhaErrada456!", false},
		{"senha vazia", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := password.Matches(tt.plain)
			if got != tt.want {
				t.Errorf("Matches(%q) got = %v, want %v", tt.plain, got, tt.want)
			}
		})
	}
}

func TestNewPasswordFromHash(t *testing.T) {
	original, err := NewPassword("SenhaForte123!")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	restored := NewPasswordFromHash(original.Hash())

	if !restored.Matches("SenhaForte123!") {
		t.Errorf("restored password should match the original plain text")
	}

	if restored.Hash() != original.Hash() {
		t.Errorf("Hash() got = %q, want %q", restored.Hash(), original.Hash())
	}
}
