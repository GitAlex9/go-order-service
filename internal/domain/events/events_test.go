package events

import "testing"

func TestEventNames(t *testing.T) {
	tests := []struct {
		name  string
		event Named
		want  string
	}{
		{"order paid", OrderPaidEvent{}, "order.paid"},
		{"order canceled", OrderCanceledEvent{}, "order.canceled"},
		{"product stock decreased", ProductStockDecreasedEvent{}, "product.stock_decreased"},
		{"product stock increased", ProductStockIncreasedEvent{}, "product.stock_increased"},
		{"customer linked to user", CustomerLinkedToUserEvent{}, "customer.linked_to_user"},
		{"customer renamed", CustomerRenamedEvent{}, "customer.renamed"},
		{"user password changed", UserPasswordChangedEvent{}, "user.password_changed"},
		{"user email changed", UserEmailChangedEvent{}, "user.email_changed"},
		{"user deactivated", UserDeactivatedEvent{}, "user.deactivated"},
		{"user activated", UserActivatedEvent{}, "user.activated"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.event.EventName()
			if got != tt.want {
				t.Errorf("EventName() got = %q, want %q", got, tt.want)
			}
		})
	}
}
