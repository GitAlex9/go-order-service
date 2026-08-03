package postgres

import (
	"context"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestProductRepository(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewProductRepository(db.Pool)

	makeProduct := func(t *testing.T) *entities.Product {
		t.Helper()

		price, err := valueobjects.NewMoney(5000)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		product, err := entities.NewProduct("Notebook", "Descrição", price, 10)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		return product
	}

	type testCase struct {
		name string
		run  func(t *testing.T)
	}

	tests := []testCase{
		{
			name: "SaveFindByID",
			run: func(t *testing.T) {
				product := makeProduct(t)
				if err := repo.Save(context.Background(), product); err != nil {
					t.Fatalf("Save() error = %v, want nil", err)
				}

				found, err := repo.FindByID(context.Background(), product.ID())
				if err != nil {
					t.Fatalf("FindByID() error = %v, want nil", err)
				}
				if got, want := found.ID(), product.ID(); got != want {
					t.Fatalf("ID() got = %v, want %v", got, want)
				}
				if got, want := found.Price().Cents(), product.Price().Cents(); got != want {
					t.Fatalf("Price() got = %d, want %d", got, want)
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
			name: "ListExists",
			run: func(t *testing.T) {
				for i := 0; i < 4; i++ {
					product := makeProduct(t)
					if err := repo.Save(context.Background(), product); err != nil {
						t.Fatalf("Save() error = %v, want nil", err)
					}
				}

				list, err := repo.List(context.Background(), 0, 2)
				if err != nil {
					t.Fatalf("List() error = %v, want nil", err)
				}
				if got, want := len(list), 2; got != want {
					t.Fatalf("len(list) got = %d, want %d", got, want)
				}

				exists, err := repo.Exists(context.Background(), list[0].ID())
				if err != nil {
					t.Fatalf("Exists() error = %v, want nil", err)
				}
				if !exists {
					t.Fatalf("Exists() got = %v, want true", exists)
				}
			},
		},
		{
			name: "DeleteMissing",
			run: func(t *testing.T) {
				product := makeProduct(t)
				if err := repo.Save(context.Background(), product); err != nil {
					t.Fatalf("Save() error = %v, want nil", err)
				}

				if err := repo.Delete(context.Background(), product.ID()); err != nil {
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
