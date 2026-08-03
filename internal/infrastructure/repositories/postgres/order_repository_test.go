package postgres

import (
	"context"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestOrderRepository(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewOrderRepository(db.Pool)
	customerRepo := NewCustomerRepository(db.Pool)
	productRepo := NewProductRepository(db.Pool)

	makeCustomer := func(t *testing.T, emailRaw, cpfRaw string) uuid.UUID {
		t.Helper()

		email, err := valueobjects.NewEmail(emailRaw)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		cpf, err := valueobjects.NewCPF(cpfRaw)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		customer, err := entities.NewCustomer("Cliente Teste", email, cpf)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		if err := customerRepo.Save(context.Background(), customer); err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		return customer.ID()
	}

	makeProduct := func(t *testing.T) uuid.UUID {
		t.Helper()

		price, err := valueobjects.NewMoney(1500)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		product, err := entities.NewProduct("Notebook", "Descrição", price, 10)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		if err := productRepo.Save(context.Background(), product); err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		return product.ID()
	}

	makeOrder := func(t *testing.T, customerID, productID uuid.UUID) *entities.Order {
		t.Helper()

		price, err := valueobjects.NewMoney(1500)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		item, err := entities.NewOrderItem(productID, "Produto Teste", price, 2)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		order, err := entities.NewOrder(customerID, []entities.OrderItem{*item})
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		return order
	}

	type testCase struct {
		name string
		run  func(t *testing.T)
	}

	tests := []testCase{
		{
			name: "SaveFindByID",
			run: func(t *testing.T) {
				order := makeOrder(t, makeCustomer(t, "cliente@teste.com", "52998224725"), makeProduct(t))
				if err := repo.Save(context.Background(), order); err != nil {
					t.Fatalf("Save() error = %v, want nil", err)
				}

				found, err := repo.FindByID(context.Background(), order.ID())
				if err != nil {
					t.Fatalf("FindByID() error = %v, want nil", err)
				}
				if got, want := found.ID(), order.ID(); got != want {
					t.Fatalf("ID() got = %v, want %v", got, want)
				}
				if got, want := found.CustomerID(), order.CustomerID(); got != want {
					t.Fatalf("CustomerID() got = %v, want %v", got, want)
				}
			},
		},
		{
			name: "FindByIDMissing",
			run: func(t *testing.T) {
				_, err := repo.FindByID(context.Background(), uuid.New())
				if err != domainerrors.ErrNotFound {
					t.Fatalf("FindByID() error = %v, want %v", err, domainerrors.ErrNotFound)
				}
			},
		},
		{
			name: "FindByCustomerIDList",
			run: func(t *testing.T) {
				customerID := makeCustomer(t, "cliente@teste.com", "52998224725")
				for i := 0; i < 3; i++ {
					order := makeOrder(t, customerID, makeProduct(t))
					if err := repo.Save(context.Background(), order); err != nil {
						t.Fatalf("Save() error = %v, want nil", err)
					}
				}

				found, err := repo.FindByCustomerID(context.Background(), customerID)
				if err != nil {
					t.Fatalf("FindByCustomerID() error = %v, want nil", err)
				}
				if got, want := len(found), 3; got != want {
					t.Fatalf("len(found) got = %d, want %d", got, want)
				}

				list, err := repo.List(context.Background(), 0, 2)
				if err != nil {
					t.Fatalf("List() error = %v, want nil", err)
				}
				if got, want := len(list), 2; got != want {
					t.Fatalf("len(list) got = %d, want %d", got, want)
				}
			},
		},
		{
			name: "ExistsDelete",
			run: func(t *testing.T) {
				order := makeOrder(t, makeCustomer(t, "cliente@teste.com", "52998224725"), makeProduct(t))
				if err := repo.Save(context.Background(), order); err != nil {
					t.Fatalf("Save() error = %v, want nil", err)
				}

				exists, err := repo.Exists(context.Background(), order.ID())
				if err != nil {
					t.Fatalf("Exists() error = %v, want nil", err)
				}
				if !exists {
					t.Fatalf("Exists() got = %v, want true", exists)
				}

				if err := repo.Delete(context.Background(), order.ID()); err != nil {
					t.Fatalf("Delete() error = %v, want nil", err)
				}

				if err := repo.Delete(context.Background(), uuid.New()); err != domainerrors.ErrNotFound {
					t.Fatalf("Delete() error = %v, want %v", err, domainerrors.ErrNotFound)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db.Reset(t)
			tt.run(t)
		})
	}
}
