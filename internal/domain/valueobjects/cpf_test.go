package valueobjects

import (
	"errors"
	"testing"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func TestNewCPF(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr error
	}{
		{"cpf válido sem máscara", "52998224725", "52998224725", nil},
		{"cpf válido com máscara", "529.982.247-25", "52998224725", nil},
		{"cpf com menos de 11 dígitos", "123456789", "", domainerrors.ErrInsufficientCPFLength},
		{"cpf com mais de 11 dígitos", "529982247256", "", domainerrors.ErrInsufficientCPFLength},
		{"cpf vazio", "", "", domainerrors.ErrInsufficientCPFLength},
		{"cpf com todos os dígitos iguais", "11111111111", "", domainerrors.ErrInvalidCPF},
		{"cpf com dígito verificador inválido", "52998224700", "", domainerrors.ErrInvalidCPF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCPF(tt.raw)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewCPF(%q) error = %v, want %v", tt.raw, err, tt.wantErr)
			}

			if tt.wantErr == nil && got.String() != tt.want {
				t.Errorf("got.String() = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestCPF_Formatted(t *testing.T) {
	cpf, err := NewCPF("52998224725")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	got := cpf.Formatted()
	want := "529.982.247-25"

	if got != want {
		t.Errorf("Formatted() got = %q, want %q", got, want)
	}
}
