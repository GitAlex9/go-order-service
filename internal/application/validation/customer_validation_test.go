package validation

import (
	"testing"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

func TestValidateCreateCustomer(t *testing.T) {
	tests := []struct {
		name       string
		req        dto.CreateCustomerRequest
		wantErrors int
		wantFields []string
	}{
		{
			name:       "dados válidos",
			req:        dto.CreateCustomerRequest{Name: "Cliente Teste", Email: "cliente@teste.com", CPF: "52998224725"},
			wantErrors: 0,
		},
		{
			name:       "nome vazio",
			req:        dto.CreateCustomerRequest{Name: "", Email: "cliente@teste.com", CPF: "52998224725"},
			wantErrors: 1,
			wantFields: []string{"name"},
		},
		{
			name:       "email inválido",
			req:        dto.CreateCustomerRequest{Name: "Cliente Teste", Email: "email-invalido", CPF: "52998224725"},
			wantErrors: 1,
			wantFields: []string{"email"},
		},
		{
			name:       "cpf inválido",
			req:        dto.CreateCustomerRequest{Name: "Cliente Teste", Email: "cliente@teste.com", CPF: "123"},
			wantErrors: 1,
			wantFields: []string{"cpf"},
		},
		{
			name:       "todos os campos inválidos ao mesmo tempo",
			req:        dto.CreateCustomerRequest{Name: "ab", Email: "email-invalido", CPF: "123"},
			wantErrors: 3,
			wantFields: []string{"name", "email", "cpf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, verr := ValidateCreateCustomer(tt.req)

			if got := len(verr.Errors); got != tt.wantErrors {
				t.Fatalf("got %d validation errors, want %d (errors: %+v)", got, tt.wantErrors, verr.Errors)
			}

			if tt.wantErrors == 0 {
				if verr.HasErrors() {
					t.Errorf("HasErrors() got = true, want false")
				}
				return
			}

			if !verr.HasErrors() {
				t.Errorf("HasErrors() got = false, want true")
			}

			gotFields := make(map[string]bool)
			for _, fe := range verr.Errors {
				gotFields[fe.Field] = true
			}
			for _, field := range tt.wantFields {
				if !gotFields[field] {
					t.Errorf("expected an error for field %q, got fields: %v", field, verr.Errors)
				}
			}
		})
	}
}
