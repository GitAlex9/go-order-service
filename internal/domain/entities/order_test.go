package entities

import (
	"errors"
	"testing"
	"time"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func newTestOrderItem(t *testing.T, cents int64, quantity int) OrderItem {
	t.Helper()

	price, err := valueobjects.NewMoney(cents)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	item, err := NewOrderItem(uuid.New(), "Produto Teste", price, quantity)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	return *item
}

func timeNow() time.Time { return time.Now() }

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name       string
		customerID uuid.UUID
		items      []OrderItem
		wantErr    error
	}{
		{"pedido válido com um item", uuid.New(), []OrderItem{newTestOrderItem(t, 1000, 1)}, nil},
		{"customer id vazio", uuid.Nil, []OrderItem{newTestOrderItem(t, 1000, 1)}, domainerrors.ErrInvalidCustomer},
		{"sem itens", uuid.New(), []OrderItem{}, domainerrors.ErrEmptyOrder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrder(tt.customerID, tt.items)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewOrder() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil && got != nil {
				t.Errorf("got order = %v, want nil", got)
			}
			if tt.wantErr == nil && got == nil {
				t.Errorf("got nil order, want non-nil")
			}
		})
	}
}

func TestNewOrder_FieldsArePersisted(t *testing.T) {
	customerID := uuid.New()
	item := newTestOrderItem(t, 1000, 2)

	order, err := NewOrder(customerID, []OrderItem{item})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if got := order.ID(); got == uuid.Nil {
		t.Errorf("ID() got = %v, want non-nil uuid", got)
	}

	if got, want := order.CustomerID(), customerID; got != want {
		t.Errorf("CustomerID() got = %v, want %v", got, want)
	}

	if got, want := order.Status(), OrderStatusPending; got != want {
		t.Errorf("Status() got = %v, want %v", got, want)
	}

	if got, want := len(order.Items()), 1; got != want {
		t.Errorf("len(Items()) got = %d, want %d", got, want)
	}

	if got := order.CreatedAt(); got.IsZero() {
		t.Errorf("CreatedAt() got = zero value, want a set timestamp")
	}
}

func TestOrder_Total(t *testing.T) {
	tests := []struct {
		name  string
		items []OrderItem
		want  int64
	}{
		{"um item", []OrderItem{newTestOrderItem(t, 1000, 2)}, 2000},
		{"múltiplos itens", []OrderItem{newTestOrderItem(t, 1000, 2), newTestOrderItem(t, 500, 3)}, 3500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(uuid.New(), tt.items)
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			got := order.Total().Cents()
			if got != tt.want {
				t.Errorf("Total().Cents() got = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestOrder_Pay(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus OrderStatus
		wantErr       error
		wantStatus    OrderStatus
	}{
		{"pending pode ser pago", OrderStatusPending, nil, OrderStatusPaid},
		{"paid não pode ser pago de novo", OrderStatusPaid, domainerrors.ErrInvalidStatusTransition, OrderStatusPaid},
		{"canceled não pode ser pago", OrderStatusCanceled, domainerrors.ErrInvalidStatusTransition, OrderStatusCanceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := RebuildOrder(uuid.New(), uuid.New(), tt.initialStatus, []OrderItem{newTestOrderItem(t, 1000, 1)}, timeNow(), timeNow())

			err := order.Pay()

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Pay() error = %v, want %v", err, tt.wantErr)
			}

			if got := order.Status(); got != tt.wantStatus {
				t.Errorf("Status() got = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestOrder_Cancel(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus OrderStatus
		wantErr       error
		wantStatus    OrderStatus
	}{
		{"pending pode ser cancelado", OrderStatusPending, nil, OrderStatusCanceled},
		{"paid não pode ser cancelado", OrderStatusPaid, domainerrors.ErrInvalidStatusTransition, OrderStatusPaid},
		{"canceled não pode ser cancelado de novo", OrderStatusCanceled, domainerrors.ErrInvalidStatusTransition, OrderStatusCanceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := RebuildOrder(uuid.New(), uuid.New(), tt.initialStatus, []OrderItem{newTestOrderItem(t, 1000, 1)}, timeNow(), timeNow())

			err := order.Cancel()

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Cancel() error = %v, want %v", err, tt.wantErr)
			}

			if got := order.Status(); got != tt.wantStatus {
				t.Errorf("Status() got = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestOrder_AddItem(t *testing.T) {
	t.Run("adiciona item novo", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), []OrderItem{newTestOrderItem(t, 1000, 1)})
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		newItem := newTestOrderItem(t, 500, 1)

		if err := order.AddItem(newItem); err != nil {
			t.Fatalf("AddItem() error = %v, want nil", err)
		}

		if got, want := len(order.Items()), 2; got != want {
			t.Errorf("len(Items()) got = %d, want %d", got, want)
		}
	})

	t.Run("mescla quantidade quando produto já existe no pedido", func(t *testing.T) {
		price, err := valueobjects.NewMoney(1000)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		productID := uuid.New()

		firstItem, err := NewOrderItem(productID, "Produto Teste", price, 2)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		order, err := NewOrder(uuid.New(), []OrderItem{*firstItem})
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		duplicateItem, err := NewOrderItem(productID, "Produto Teste", price, 3)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		if err := order.AddItem(*duplicateItem); err != nil {
			t.Fatalf("AddItem() error = %v, want nil", err)
		}

		if got, want := len(order.Items()), 1; got != want {
			t.Fatalf("len(Items()) got = %d, want %d", got, want)
		}

		if got, want := order.Items()[0].Quantity(), 5; got != want {
			t.Errorf("Items()[0].Quantity() got = %d, want %d", got, want)
		}
	})

	t.Run("não permite adicionar item em pedido pago", func(t *testing.T) {
		order := RebuildOrder(uuid.New(), uuid.New(), OrderStatusPaid, []OrderItem{newTestOrderItem(t, 1000, 1)}, timeNow(), timeNow())

		err := order.AddItem(newTestOrderItem(t, 500, 1))

		if !errors.Is(err, domainerrors.ErrOrderNotEditable) {
			t.Errorf("AddItem() error = %v, want %v", err, domainerrors.ErrOrderNotEditable)
		}
	})
}

func TestOrder_RemoveItem(t *testing.T) {
	t.Run("remove item existente", func(t *testing.T) {
		item1 := newTestOrderItem(t, 1000, 1)
		item2 := newTestOrderItem(t, 500, 1)

		order, err := NewOrder(uuid.New(), []OrderItem{item1, item2})
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		if err := order.RemoveItem(item1.ProductID()); err != nil {
			t.Fatalf("RemoveItem() error = %v, want nil", err)
		}

		if got, want := len(order.Items()), 1; got != want {
			t.Fatalf("len(Items()) got = %d, want %d", got, want)
		}

		if got, want := order.Items()[0].ProductID(), item2.ProductID(); got != want {
			t.Errorf("remaining item got = %v, want %v", got, want)
		}
	})

	t.Run("erro ao remover produto que não está no pedido", func(t *testing.T) {
		order, err := NewOrder(uuid.New(), []OrderItem{newTestOrderItem(t, 1000, 1)})
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = order.RemoveItem(uuid.New())

		if !errors.Is(err, domainerrors.ErrOrderItemNotFound) {
			t.Errorf("RemoveItem() error = %v, want %v", err, domainerrors.ErrOrderItemNotFound)
		}
	})

	t.Run("remover último item resulta em pedido vazio", func(t *testing.T) {
		item := newTestOrderItem(t, 1000, 1)

		order, err := NewOrder(uuid.New(), []OrderItem{item})
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		err = order.RemoveItem(item.ProductID())

		if !errors.Is(err, domainerrors.ErrEmptyOrder) {
			t.Errorf("RemoveItem() error = %v, want %v", err, domainerrors.ErrEmptyOrder)
		}
	})

	t.Run("não permite remover item de pedido pago", func(t *testing.T) {
		item := newTestOrderItem(t, 1000, 1)
		order := RebuildOrder(uuid.New(), uuid.New(), OrderStatusPaid, []OrderItem{item}, timeNow(), timeNow())

		err := order.RemoveItem(item.ProductID())

		if !errors.Is(err, domainerrors.ErrOrderNotEditable) {
			t.Errorf("RemoveItem() error = %v, want %v", err, domainerrors.ErrOrderNotEditable)
		}
	})
}
