package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

func TestCustomerRepository(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewCustomerRepository(db.Pool)

	makeCustomer := func(t *testing.T, name, emailRaw, cpfRaw string) *entities.Customer {
		t.Helper()

		email, err := valueobjects.NewEmail(emailRaw)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		cpf, err := valueobjects.NewCPF(cpfRaw)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		customer, err := entities.NewCustomer(name, email, cpf)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		return customer
	}

	type testCase struct {
		name string
		run  func(t *testing.T)
	}

	tests := []testCase{
		{
			name: "SaveUpdateDuplicate",
			run: func(t *testing.T) {
				customer := makeCustomer(t, "Cliente Teste", "cliente1@teste.com", "52998224725")
				if err := repo.Save(context.Background(), customer); err != nil {
					t.Fatalf("Save() error = %v, want nil", err)
				}

				t.Run("duplicate customer returns ErrDuplicateCustomer", func(t *testing.T) {
					duplicate := makeCustomer(t, "Cliente Duplicado", "cliente1@teste.com", "11144477735")
					if err := repo.Save(context.Background(), duplicate); err != domainerrors.ErrDuplicateCustomer {
						t.Fatalf("Save() error = %v, want %v", err, domainerrors.ErrDuplicateCustomer)
					}
				})

				t.Run("update existing customer", func(t *testing.T) {
					customer2 := makeCustomer(t, "Cliente Original", "cliente2@teste.com", "12345678909")
					if err := repo.Save(context.Background(), customer2); err != nil {
						t.Fatalf("Save() error = %v, want nil", err)
					}

					if err := customer2.Rename("Cliente Atualizado"); err != nil {
						t.Fatalf("Rename() error = %v, want nil", err)
					}
					if err := repo.Save(context.Background(), customer2); err != nil {
						t.Fatalf("Save() update error = %v, want nil", err)
					}

					found, err := repo.FindByID(context.Background(), customer2.ID())
					if err != nil {
						t.Fatalf("FindByID() error = %v, want nil", err)
					}
					if got, want := found.Name(), "Cliente Atualizado"; got != want {
						t.Fatalf("Name() got = %q, want %q", got, want)
					}
				})
			},
		},
		{
			name: "FindByIDFindByEmail",
			run: func(t *testing.T) {
				customer := makeCustomer(t, "Cliente Busca", "busca@teste.com", "86753090070")
				if err := repo.Save(context.Background(), customer); err != nil {
					t.Fatalf("Save() error = %v, want nil", err)
				}

				t.Run("existing id returns customer", func(t *testing.T) {
					found, err := repo.FindByID(context.Background(), customer.ID())
					if err != nil {
						t.Fatalf("FindByID() error = %v, want nil", err)
					}
					if got, want := found.ID(), customer.ID(); got != want {
						t.Fatalf("ID() got = %v, want %v", got, want)
					}
				})

				t.Run("missing id returns ErrNotFound", func(t *testing.T) {
					_, err := repo.FindByID(context.Background(), uuid.New())
					if err != domainerrors.ErrNotFound {
						t.Fatalf("FindByID() error = %v, want %v", err, domainerrors.ErrNotFound)
					}
				})

				t.Run("existing email returns customer", func(t *testing.T) {
					found, err := repo.FindByEmail(context.Background(), customer.Email())
					if err != nil {
						t.Fatalf("FindByEmail() error = %v, want nil", err)
					}
					if got, want := found.Email().String(), customer.Email().String(); got != want {
						t.Fatalf("Email() got = %q, want %q", got, want)
					}
				})

				t.Run("missing email returns ErrNotFound", func(t *testing.T) {
					missingEmail, err := valueobjects.NewEmail("naoexiste@teste.com")
					if err != nil {
						t.Fatalf("setup failed: %v", err)
					}
					_, err = repo.FindByEmail(context.Background(), missingEmail)
					if err != domainerrors.ErrNotFound {
						t.Fatalf("FindByEmail() error = %v, want %v", err, domainerrors.ErrNotFound)
					}
				})
			},
		},
		{
			name: "ListExistsDelete",
			run: func(t *testing.T) {
				cpfs := []string{"11144477735", "12345678909", "52998224725", "00000000191", "98765432100"}
				for i := 0; i < 5; i++ {
					customer := makeCustomer(t, fmt.Sprintf("Cliente %c", 'A'+i), fmt.Sprintf("list%c@teste.com", 'a'+i), cpfs[i])
					if err := repo.Save(context.Background(), customer); err != nil {
						t.Fatalf("Save() error = %v, want nil", err)
					}
				}

				t.Run("pagination returns correct size", func(t *testing.T) {
					list, err := repo.List(context.Background(), 0, 2)
					if err != nil {
						t.Fatalf("List() error = %v, want nil", err)
					}
					if got, want := len(list), 2; got != want {
						t.Fatalf("len(list) got = %d, want %d", got, want)
					}
				})

				t.Run("exists returns true/false", func(t *testing.T) {
					all, err := repo.List(context.Background(), 0, 1)
					if err != nil {
						t.Fatalf("List() error = %v, want nil", err)
					}
					if len(all) != 1 {
						t.Fatalf("expected one customer, got %d", len(all))
					}

					exists, err := repo.Exists(context.Background(), all[0].ID())
					if err != nil {
						t.Fatalf("Exists() error = %v, want nil", err)
					}
					if !exists {
						t.Fatalf("Exists() got = %v, want true", exists)
					}

					notExists, err := repo.Exists(context.Background(), uuid.New())
					if err != nil {
						t.Fatalf("Exists() error = %v, want nil", err)
					}
					if notExists {
						t.Fatalf("Exists() got = %v, want false", notExists)
					}
				})

				t.Run("delete existing and missing returns ErrNotFound", func(t *testing.T) {
					all, err := repo.List(context.Background(), 0, 1)
					if err != nil {
						t.Fatalf("List() error = %v, want nil", err)
					}
					if err := repo.Delete(context.Background(), all[0].ID()); err != nil {
						t.Fatalf("Delete() error = %v, want nil", err)
					}

					err = repo.Delete(context.Background(), uuid.New())
					if err != domainerrors.ErrNotFound {
						t.Fatalf("Delete() error = %v, want %v", err, domainerrors.ErrNotFound)
					}
				})
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
