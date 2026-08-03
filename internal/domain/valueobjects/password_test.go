package valueobjects

import (
	"errors"
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestNewPassword_Valid(t *testing.T) {
	tests := []struct {
		name    string
		plain   string
		wantErr error
	}{
		{"senha forte com símbolos", "Teste@123456", nil},
		{"senha forte com números", "Teste123456!", nil},
		{"senha forte com minúsculas e maiúsculas", "Teste@123", nil}, // 8 caracteres, válido
		{"senha com caracteres especiais variados", "Teste!@#$123", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := NewPassword(tt.plain)
			if err != tt.wantErr {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && password.Hash() == "" {
				t.Errorf("expected hash to be set")
			}
		})
	}
}

func TestNewPassword_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		plain   string
		wantErr error
	}{
		{"senha curta", "Test@12", domainerrors.ErrWeakPassword},
		{"senha sem maiúscula", "teste@123456", domainerrors.ErrPasswordNoUpper},
		{"senha sem minúscula", "TESTE@123456", domainerrors.ErrPasswordNoLower},
		{"senha sem número", "Teste@abcdef", domainerrors.ErrPasswordNoNumber},
		{"senha sem símbolo", "Teste123456", domainerrors.ErrPasswordNoSpecial},
		{"senha vazia", "", domainerrors.ErrWeakPassword},
		{"senha com apenas maiúsculas", "TESTE@123", domainerrors.ErrPasswordNoLower},
		{"senha com apenas minúsculas", "teste@123", domainerrors.ErrPasswordNoUpper},
		{"senha com apenas números", "12345678", errors.Join(domainerrors.ErrPasswordNoUpper, domainerrors.ErrPasswordNoLower, domainerrors.ErrPasswordNoSpecial)},
		{"senha com apenas símbolos", "@#$%^&*()", errors.Join(domainerrors.ErrPasswordNoUpper, domainerrors.ErrPasswordNoLower, domainerrors.ErrPasswordNoNumber)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := NewPassword(tt.plain)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if password != (Password{}) {
				t.Errorf("expected zero value Password, got %v", password)
			}
		})
	}
}

func TestPassword_Matches(t *testing.T) {
	plain := "Teste@123456"
	password, _ := NewPassword(plain)

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"senha correta", plain, true},
		{"senha incorreta", "SenhaErrada", false},
		{"senha vazia", "", false},
		{"senha com diferença de caractere", "Teste@12345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := password.Matches(tt.input)
			if got != tt.want {
				t.Errorf("Matches(%q) got = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPassword_Hash(t *testing.T) {
	plain := "Teste@123456"
	password, _ := NewPassword(plain)

	hash := password.Hash()
	if hash == "" {
		t.Errorf("Hash() returned empty string")
	}
	// Verifica se o hash tem o formato do bcrypt (começa com "$2a$")
	if len(hash) < 4 || hash[:4] != "$2a$" {
		t.Errorf("Hash() does not start with bcrypt prefix: %s", hash)
	}
}

func TestNewPasswordFromHash(t *testing.T) {
	plain := "Teste@123456"
	original, _ := NewPassword(plain)
	hash := original.Hash()

	restored := NewPasswordFromHash(hash)

	if restored.Hash() != hash {
		t.Errorf("Hash() got = %q, want %q", restored.Hash(), hash)
	}
	if !restored.Matches(plain) {
		t.Errorf("restored password does not match original")
	}
}

func TestValidateStrength_ErrorCombination(t *testing.T) {
	plain := "abc"
	err := validateStrength(plain)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
