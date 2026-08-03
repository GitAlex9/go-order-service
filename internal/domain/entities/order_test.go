package entities

import (
	"testing"
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestNewOrder_Valid(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto Teste", price, 2)

	tests := []struct {
		name string
		args []OrderItem
	}{
		{"um item", []OrderItem{*item}},
		{"múltiplos itens", []OrderItem{*item, *item}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(customerID, tt.args)
			if err != nil {
				t.Fatalf("expected creation to succeed, got error: %v", err)
			}
			if order.ID() == uuid.Nil {
				t.Errorf("expected ID to be set")
			}
			if order.CustomerID() != customerID {
				t.Errorf("CustomerID() got = %v, want %v", order.CustomerID(), customerID)
			}
			if order.Status() != OrderStatusPending {
				t.Errorf("Status() got = %v, want %v", order.Status(), OrderStatusPending)
			}
			if len(order.Items()) != len(tt.args) {
				t.Errorf("len(Items()) got = %d, want %d", len(order.Items()), len(tt.args))
			}
			if order.CreatedAt().IsZero() || order.UpdatedAt().IsZero() {
				t.Errorf("expected timestamps to be set")
			}
		})
	}
}

func TestNewOrder_Invalid(t *testing.T) {
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto Teste", price, 2)

	tests := []struct {
		name       string
		customerID uuid.UUID
		items      []OrderItem
		wantErr    error
	}{
		{"customerID vazio", uuid.Nil, []OrderItem{*item}, domainerrors.ErrInvalidCustomer},
		{"sem itens", uuid.New(), []OrderItem{}, domainerrors.ErrEmptyOrder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.customerID, tt.items)
			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if order != nil {
				t.Errorf("expected nil order")
			}
		})
	}
}

func TestOrder_AddItem(t *testing.T) {
	customerID := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)

	item1, _ := NewOrderItem(productID1, "Produto 1", price, 2)
	item2, _ := NewOrderItem(productID2, "Produto 2", price, 3)

	order, _ := NewOrder(customerID, []OrderItem{*item1})

	tests := []struct {
		name      string
		item      OrderItem
		wantCount int
		wantQty   int
	}{
		{"adicionar novo item", *item2, 2, 0},
		{"adicionar item existente (soma)", *item1, 1, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := order.UpdatedAt()
			err := order.AddItem(tt.item)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(order.Items()) != tt.wantCount {
				t.Errorf("len(Items()) got = %d, want %d", len(order.Items()), tt.wantCount)
			}
			if tt.wantQty > 0 {
				for _, item := range order.Items() {
					if item.ProductID() == tt.item.ProductID() && item.Quantity() != tt.wantQty {
						t.Errorf("Quantity() got = %d, want %d", item.Quantity(), tt.wantQty)
					}
				}
			}
			if order.UpdatedAt() == oldUpdatedAt {
				t.Errorf("expected UpdatedAt to change")
			}
		})
	}
}

func TestOrder_AddItem_OrderNotEditable(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto", price, 2)

	order, _ := NewOrder(customerID, []OrderItem{*item})
	order.Pay()

	err := order.AddItem(*item)
	if err != domainerrors.ErrOrderNotEditable {
		t.Errorf("expected ErrOrderNotEditable, got %v", err)
	}
}

func TestOrder_RemoveItem(t *testing.T) {
	customerID := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)

	item1, _ := NewOrderItem(productID1, "Produto 1", price, 2)
	item2, _ := NewOrderItem(productID2, "Produto 2", price, 3)

	order, _ := NewOrder(customerID, []OrderItem{*item1, *item2})

	tests := []struct {
		name      string
		productID uuid.UUID
		wantCount int
		wantErr   error
	}{
		{"remover item existente", productID1, 1, nil},
		{"remover item inexistente", uuid.New(), 1, domainerrors.ErrOrderItemNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := order.UpdatedAt()
			err := order.RemoveItem(tt.productID)

			if err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
			if err == nil {
				if len(order.Items()) != tt.wantCount {
					t.Errorf("len(Items()) got = %d, want %d", len(order.Items()), tt.wantCount)
				}
				if order.UpdatedAt() == oldUpdatedAt {
					t.Errorf("expected UpdatedAt to change")
				}
			}
		})
	}
}

func TestOrder_RemoveItem_OrderNotEditable(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto", price, 2)

	order, _ := NewOrder(customerID, []OrderItem{*item})
	order.Pay()

	err := order.RemoveItem(productID)
	if err != domainerrors.ErrOrderNotEditable {
		t.Errorf("expected ErrOrderNotEditable, got %v", err)
	}
}

func TestOrder_RemoveItem_EmptyOrder(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto", price, 2)

	order, _ := NewOrder(customerID, []OrderItem{*item})
	err := order.RemoveItem(productID)

	if err != domainerrors.ErrEmptyOrder {
		t.Errorf("expected ErrEmptyOrder, got %v", err)
	}
	if len(order.Items()) != 0 {
		t.Errorf("len(Items()) got = %d, want 0", len(order.Items()))
	}
}

func TestOrder_Total(t *testing.T) {
	price1, _ := valueobjects.NewMoneyFromFloat(10.50)
	price2, _ := valueobjects.NewMoneyFromFloat(5.25)
	productID1 := uuid.New()
	productID2 := uuid.New()

	item1, _ := NewOrderItem(productID1, "Produto 1", price1, 2)
	item2, _ := NewOrderItem(productID2, "Produto 2", price2, 3)

	customerID := uuid.New()
	order, _ := NewOrder(customerID, []OrderItem{*item1, *item2})

	expected, _ := valueobjects.NewMoneyFromFloat(10.50*2 + 5.25*3)
	got := order.Total()

	if !got.Equals(expected) {
		t.Errorf("Total() got = %v, want %v", got.Amount(), expected.Amount())
	}
}

func TestOrder_Pay(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto", price, 2)

	order, _ := NewOrder(customerID, []OrderItem{*item})
	oldUpdatedAt := order.UpdatedAt()

	err := order.Pay()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.Status() != OrderStatusPaid {
		t.Errorf("Status() got = %v, want %v", order.Status(), OrderStatusPaid)
	}
	if order.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	// Verifica o evento
	orderEvents := order.Events()
	lastEvent := orderEvents[len(orderEvents)-1]
	paidEvent, ok := lastEvent.(events.OrderPaidEvent)
	if !ok {
		t.Errorf("expected OrderPaidEvent, got %T", lastEvent)
	}
	if paidEvent.OrderID != order.ID() {
		t.Errorf("expected OrderID %v, got %v", order.ID(), paidEvent.OrderID)
	}
}

func TestOrder_Cancel(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto", price, 2)

	order, _ := NewOrder(customerID, []OrderItem{*item})
	oldUpdatedAt := order.UpdatedAt()

	err := order.Cancel()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.Status() != OrderStatusCanceled {
		t.Errorf("Status() got = %v, want %v", order.Status(), OrderStatusCanceled)
	}
	if order.UpdatedAt() == oldUpdatedAt {
		t.Errorf("expected UpdatedAt to change")
	}

	// Verifica o evento
	orderEvents1 := order.Events()
	lastEvent := orderEvents1[len(orderEvents1)-1]
	cancelEvent, ok := lastEvent.(events.OrderCanceledEvent)
	if !ok {
		t.Errorf("expected OrderCanceledEvent, got %T", lastEvent)
	}
	if cancelEvent.OrderID != order.ID() {
		t.Errorf("expected OrderID %v, got %v", order.ID(), cancelEvent.OrderID)
	}
}

func TestOrder_InvalidStatusTransition(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto", price, 2)

	order, _ := NewOrder(customerID, []OrderItem{*item})
	order.Pay()

	err := order.Pay()
	if err != domainerrors.ErrInvalidStatusTransition {
		t.Errorf("expected ErrInvalidStatusTransition, got %v", err)
	}

	err = order.Cancel()
	if err != domainerrors.ErrInvalidStatusTransition {
		t.Errorf("expected ErrInvalidStatusTransition, got %v", err)
	}
}

func TestOrder_Rebuild(t *testing.T) {
	id := uuid.New()
	customerID := uuid.New()
	productID := uuid.New()
	price, _ := valueobjects.NewMoneyFromFloat(10.50)
	item, _ := NewOrderItem(productID, "Produto", price, 2)
	now := time.Now()

	order := RebuildOrder(id, customerID, OrderStatusPaid, []OrderItem{*item}, now, now)

	if order.ID() != id {
		t.Errorf("ID() got = %v, want %v", order.ID(), id)
	}
	if order.CustomerID() != customerID {
		t.Errorf("CustomerID() got = %v, want %v", order.CustomerID(), customerID)
	}
	if order.Status() != OrderStatusPaid {
		t.Errorf("Status() got = %v, want %v", order.Status(), OrderStatusPaid)
	}
	if len(order.Items()) != 1 {
		t.Errorf("len(Items()) got = %d, want 1", len(order.Items()))
	}
	if order.CreatedAt() != now {
		t.Errorf("CreatedAt() got = %v, want %v", order.CreatedAt(), now)
	}
	if order.UpdatedAt() != now {
		t.Errorf("UpdatedAt() got = %v, want %v", order.UpdatedAt(), now)
	}
}
