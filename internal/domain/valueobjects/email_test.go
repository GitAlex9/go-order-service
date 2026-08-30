package valueobjects

import (
	"errors"
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr error
	}{
		{"email válido", "cliente@teste.com", "cliente@teste.com", nil},
		{"email com maiúsculas é normalizado", "Cliente@Teste.COM", "cliente@teste.com", nil},
		{"email com espaços nas bordas", "  cliente@teste.com  ", "cliente@teste.com", nil},
		{"email vazio", "", "", domainerrors.ErrInvalidEmail},
		{"email só com espaços", "   ", "", domainerrors.ErrInvalidEmail},
		{"email sem @", "clienteteste.com", "", domainerrors.ErrInvalidEmail},
		{"email sem domínio", "cliente@", "", domainerrors.ErrInvalidEmail},
		{"email sem usuário", "@teste.com", "", domainerrors.ErrInvalidEmail},
		{"email sem TLD", "cliente@teste", "", domainerrors.ErrInvalidEmail},
		{"email com espaço no meio", "cliente teste@teste.com", "", domainerrors.ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.raw)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewEmail(%q) error = %v, want %v", tt.raw, err, tt.wantErr)
			}

			if tt.wantErr == nil && got.String() != tt.want {
				t.Errorf("got.String() = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"emails iguais", "cliente@teste.com", "cliente@teste.com", true},
		{"emails diferentes", "cliente@teste.com", "outro@teste.com", false},
		{"emails iguais com case diferente na origem", "Cliente@Teste.com", "cliente@teste.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewEmail(tt.a)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}
			b, err := NewEmail(tt.b)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := a.Equals(b)
			if got != tt.want {
				t.Errorf("Equals() got = %v, want %v", got, tt.want)
			}
		})
	}
}
