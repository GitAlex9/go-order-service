package valueobjects

import (
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestNewEmail_Valid(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"email simples", "cliente@teste.com", "cliente@teste.com"},
		{"email com maiúsculas", "Cliente@Teste.COM", "cliente@teste.com"},
		{"email com espaços", "  cliente@teste.com  ", "cliente@teste.com"},
		{"email com ponto", "cliente.teste@dominio.com", "cliente.teste@dominio.com"},
		{"email com subdomínio", "cliente@mail.dominio.com", "cliente@mail.dominio.com"},
		{"email com números", "cliente123@teste.com", "cliente123@teste.com"},
		{"email com underscore", "cliente_teste@teste.com", "cliente_teste@teste.com"},
		{"email com hífen", "cliente-teste@teste.com", "cliente-teste@teste.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.raw)
			if err != nil {
				t.Fatalf("expected creation to succeed, got error: %v", err)
			}
			if email.String() != tt.want {
				t.Errorf("String() got = %q, want %q", email.String(), tt.want)
			}
		})
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{"string vazia", "", domainerrors.ErrInvalidEmail},
		{"apenas espaços", "   ", domainerrors.ErrInvalidEmail},
		{"sem @", "cliente.teste.com", domainerrors.ErrInvalidEmail},
		{"sem domínio", "cliente@", domainerrors.ErrInvalidEmail},
		{"sem usuário", "@teste.com", domainerrors.ErrInvalidEmail},
		{"domínio inválido (sem ponto)", "cliente@testecom", domainerrors.ErrInvalidEmail},
		{"domínio com caracteres especiais inválidos", "cliente@teste!.com", domainerrors.ErrInvalidEmail},
		{"espaços no meio", "cliente @teste.com", domainerrors.ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.raw)
			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if email != (Email{}) {
				t.Errorf("expected zero value Email, got %v", email)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email, _ := NewEmail("cliente@teste.com")
	got := email.String()
	want := "cliente@teste.com"
	if got != want {
		t.Errorf("String() got = %q, want %q", got, want)
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("cliente@teste.com")
	email2, _ := NewEmail("cliente@teste.com")
	email3, _ := NewEmail("outro@teste.com")

	if !email1.Equals(email2) {
		t.Errorf("expected email1 to equal email2")
	}
	if email1.Equals(email3) {
		t.Errorf("expected email1 to NOT equal email3")
	}
}
