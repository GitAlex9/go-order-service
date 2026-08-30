package validation

import (
	"testing"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

func TestValidateCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		req        dto.CreateUserRequest
		wantErrors int
		wantFields []string
	}{
		{
			name:       "dados válidos",
			req:        dto.CreateUserRequest{Email: "admin@teste.com", Password: "SenhaForte123!", Role: "admin"},
			wantErrors: 0,
		},
		{
			name:       "email inválido",
			req:        dto.CreateUserRequest{Email: "email-invalido", Password: "SenhaForte123!", Role: "admin"},
			wantErrors: 1,
			wantFields: []string{"email"},
		},
		{
			name:       "role inválida",
			req:        dto.CreateUserRequest{Email: "admin@teste.com", Password: "SenhaForte123!", Role: "superadmin"},
			wantErrors: 1,
			wantFields: []string{"role"},
		},
		{
			name:       "email e role inválidos ao mesmo tempo",
			req:        dto.CreateUserRequest{Email: "email-invalido", Password: "SenhaForte123!", Role: "superadmin"},
			wantErrors: 2,
			wantFields: []string{"email", "role"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, verr := ValidateCreateUser(tt.req)

			if got := len(verr.Errors); got != tt.wantErrors {
				t.Fatalf("got %d validation errors, want %d (errors: %+v)", got, tt.wantErrors, verr.Errors)
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
