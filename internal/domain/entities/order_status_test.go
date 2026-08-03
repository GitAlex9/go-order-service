package entities

import "testing"

func TestOrderStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{"pending é válido", OrderStatusPending, true},
		{"paid é válido", OrderStatusPaid, true},
		{"canceled é válido", OrderStatusCanceled, true},
		{"string vazia é inválida", OrderStatus(""), false},
		{"status desconhecido é inválido", OrderStatus("SHIPPED"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   string
	}{
		{"pending", OrderStatusPending, "PENDING"},
		{"paid", OrderStatusPaid, "PAID"},
		{"canceled", OrderStatusCanceled, "CANCELED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.String()
			if got != tt.want {
				t.Errorf("String() got = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOrderStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from OrderStatus
		to   OrderStatus
		want bool
	}{
		{"pending para paid é permitido", OrderStatusPending, OrderStatusPaid, true},
		{"pending para canceled é permitido", OrderStatusPending, OrderStatusCanceled, true},
		{"pending para pending não é permitido", OrderStatusPending, OrderStatusPending, false},
		{"paid para canceled não é permitido", OrderStatusPaid, OrderStatusCanceled, false},
		{"paid para pending não é permitido", OrderStatusPaid, OrderStatusPending, false},
		{"canceled para paid não é permitido", OrderStatusCanceled, OrderStatusPaid, false},
		{"canceled para pending não é permitido", OrderStatusCanceled, OrderStatusPending, false},
		{"status desconhecido não permite nenhuma transição", OrderStatus("SHIPPED"), OrderStatusPaid, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.want {
				t.Errorf("CanTransitionTo(%s -> %s) got = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
