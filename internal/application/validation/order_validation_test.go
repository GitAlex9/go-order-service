package validation

import (
	"testing"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

func TestValidateCreateOrder(t *testing.T) {
	validProductID := uuid.New().String()

	tests := []struct {
		name       string
		req        dto.CreateOrderRequest
		wantErrors int
	}{
		{
			name: "dados válidos",
			req: dto.CreateOrderRequest{
				CustomerID: uuid.New().String(),
				Items:      []dto.OrderItemRequest{{ProductID: validProductID, Quantity: 2}},
			},
			wantErrors: 0,
		},
		{
			name: "customer_id inválido",
			req: dto.CreateOrderRequest{
				CustomerID: "não-é-um-uuid",
				Items:      []dto.OrderItemRequest{{ProductID: validProductID, Quantity: 2}},
			},
			wantErrors: 1,
		},
		{
			name: "sem itens",
			req: dto.CreateOrderRequest{
				CustomerID: uuid.New().String(),
				Items:      []dto.OrderItemRequest{},
			},
			wantErrors: 1,
		},
		{
			name: "quantidade inválida",
			req: dto.CreateOrderRequest{
				CustomerID: uuid.New().String(),
				Items:      []dto.OrderItemRequest{{ProductID: validProductID, Quantity: 0}},
			},
			wantErrors: 1,
		},
		{
			name: "product_id inválido",
			req: dto.CreateOrderRequest{
				CustomerID: uuid.New().String(),
				Items:      []dto.OrderItemRequest{{ProductID: "não-é-um-uuid", Quantity: 1}},
			},
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, verr := ValidateCreateOrder(tt.req)

			if got := len(verr.Errors); got != tt.wantErrors {
				t.Fatalf("got %d validation errors, want %d (errors: %+v)", got, tt.wantErrors, verr.Errors)
			}
		})
	}
}
