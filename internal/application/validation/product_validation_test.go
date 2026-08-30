package validation

import (
	"testing"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

func TestValidateCreateProduct(t *testing.T) {
	tests := []struct {
		name       string
		req        dto.CreateProductRequest
		wantErrors int
		wantFields []string
	}{
		{
			name:       "dados válidos",
			req:        dto.CreateProductRequest{Name: "Notebook", Description: "Notebook gamer", Price: 3500.00, Stock: 10},
			wantErrors: 0,
		},
		{
			name:       "nome vazio",
			req:        dto.CreateProductRequest{Name: "", Description: "Notebook gamer", Price: 3500.00, Stock: 10},
			wantErrors: 1,
			wantFields: []string{"name"},
		},
		{
			name:       "descrição vazia",
			req:        dto.CreateProductRequest{Name: "Notebook", Description: "", Price: 3500.00, Stock: 10},
			wantErrors: 1,
			wantFields: []string{"description"},
		},
		{
			name:       "preço zero",
			req:        dto.CreateProductRequest{Name: "Notebook", Description: "Notebook gamer", Price: 0, Stock: 10},
			wantErrors: 1,
			wantFields: []string{"price"},
		},
		{
			name:       "preço negativo",
			req:        dto.CreateProductRequest{Name: "Notebook", Description: "Notebook gamer", Price: -10, Stock: 10},
			wantErrors: 1,
			wantFields: []string{"price"},
		},
		{
			name:       "estoque negativo",
			req:        dto.CreateProductRequest{Name: "Notebook", Description: "Notebook gamer", Price: 3500.00, Stock: -1},
			wantErrors: 1,
			wantFields: []string{"stock"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, verr := ValidateCreateProduct(tt.req)

			if got := len(verr.Errors); got != tt.wantErrors {
				t.Fatalf("got %d validation errors, want %d (errors: %+v)", got, tt.wantErrors, verr.Errors)
			}

			if tt.wantErrors == 0 {
				if verr.HasErrors() {
					t.Errorf("HasErrors() got = true, want false")
				}
				return
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
