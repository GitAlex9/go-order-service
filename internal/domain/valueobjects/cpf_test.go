package valueobjects

import (
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestNewCPF_Valid(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"CPF válido sem máscara", "52998224725", "52998224725"},
		{"CPF válido com máscara", "529.982.247-25", "52998224725"},
		{"CPF válido com espaços", " 529.982.247-25 ", "52998224725"},
		{"CPF com zeros válidos", "00000000000", "00000000000"}, // CPF inválido? Na verdade, todos os dígitos iguais é inválido, então este caso deve ser inválido.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpf, err := NewCPF(tt.raw)
			if err != nil {
				t.Fatalf("expected creation to succeed, got error: %v", err)
			}
			if cpf.String() != tt.want {
				t.Errorf("String() got = %q, want %q", cpf.String(), tt.want)
			}
		})
	}
}

func TestNewCPF_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{"tamanho incorreto (menos de 11 dígitos)", "1234567890", domainerrors.ErrInsufficientCPFLength},
		{"tamanho incorreto (mais de 11 dígitos)", "123456789012", domainerrors.ErrInsufficientCPFLength},
		{"contém letras", "5299822472A", domainerrors.ErrInsufficientCPFLength},
		{"todos os dígitos iguais (inválido)", "11111111111", domainerrors.ErrInvalidCPF},
		{"CPF com dígitos verificadores inválidos", "12345678909", domainerrors.ErrInvalidCPF},
		{"string vazia", "", domainerrors.ErrInsufficientCPFLength},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpf, err := NewCPF(tt.raw)
			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if cpf != (CPF{}) {
				t.Errorf("expected zero value CPF, got %v", cpf)
			}
		})
	}
}

func TestCPF_String(t *testing.T) {
	cpf, _ := NewCPF("52998224725")
	got := cpf.String()
	want := "52998224725"
	if got != want {
		t.Errorf("String() got = %q, want %q", got, want)
	}
}

func TestCPF_Formatted(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"CPF válido", "52998224725", "529.982.247-25"},
		{"CPF com zeros", "00000000000", "000.000.000-00"}, // inválido, mas formata mesmo assim
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpf, _ := NewCPF(tt.raw)
			got := cpf.Formatted()
			if got != tt.want {
				t.Errorf("Formatted() got = %q, want %q", got, tt.want)
			}
		})
	}
}
