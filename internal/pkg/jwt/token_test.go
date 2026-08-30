package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenManager_GenerateAndParse(t *testing.T) {
	manager := NewTokenManager("test-secret", time.Hour)
	userID := uuid.New()

	token, err := manager.Generate(userID, "user@teste.com", "admin")
	if err != nil {
		t.Fatalf("Generate() error = %v, want nil", err)
	}
	if token == "" {
		t.Fatalf("Generate() returned empty token")
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	if got := claims.UserID; got != userID {
		t.Errorf("claims.UserID got = %v, want %v", got, userID)
	}
	if got, want := claims.Email, "user@teste.com"; got != want {
		t.Errorf("claims.Email got = %q, want %q", got, want)
	}
	if got, want := claims.Role, "admin"; got != want {
		t.Errorf("claims.Role got = %q, want %q", got, want)
	}
}

func TestTokenManager_Parse_InvalidToken(t *testing.T) {
	manager := NewTokenManager("test-secret", time.Hour)

	tests := []struct {
		name  string
		token string
	}{
		{"token vazio", ""},
		{"token malformado", "isso-não-é-um-jwt"},
		{"token com assinatura corrompida", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.corrompido"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.Parse(tt.token)
			if err == nil {
				t.Errorf("Parse(%q) error = nil, want an error", tt.token)
			}
		})
	}
}

func TestTokenManager_Parse_WrongSecret(t *testing.T) {
	generator := NewTokenManager("secret-a", time.Hour)
	validator := NewTokenManager("secret-b", time.Hour)

	token, err := generator.Generate(uuid.New(), "user@teste.com", "admin")
	if err != nil {
		t.Fatalf("Generate() error = %v, want nil", err)
	}

	_, err = validator.Parse(token)
	if err == nil {
		t.Fatalf("Parse() error = nil, want an error (assinado com secret diferente)")
	}
}

func TestTokenManager_Parse_ExpiredToken(t *testing.T) {
	// TTL negativo força o token a já nascer expirado.
	manager := NewTokenManager("test-secret", -time.Hour)

	token, err := manager.Generate(uuid.New(), "user@teste.com", "admin")
	if err != nil {
		t.Fatalf("Generate() error = %v, want nil", err)
	}

	_, err = manager.Parse(token)
	if err == nil {
		t.Fatalf("Parse() error = nil, want an error (token expirado)")
	}
}
