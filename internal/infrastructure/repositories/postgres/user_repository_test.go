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

func TestUserRepository(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewUserRepository(db.Pool)

	makeUser := func(t *testing.T, address string) *entities.User {
		t.Helper()

		email, err := valueobjects.NewEmail(address)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		user, err := entities.NewUser(email, "Senha123!", entities.RoleCustomer)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		return user
	}

	type testCase struct {
		name string
		run  func(t *testing.T)
	}

	tests := []testCase{
		{
			name: "SaveFindByID",
			run: func(t *testing.T) {
				user := makeUser(t, "user1@teste.com")
				if err := repo.Save(context.Background(), user); err != nil {
					t.Fatalf("Save() error = %v, want nil", err)
				}

				found, err := repo.FindByID(context.Background(), user.ID())
				if err != nil {
					t.Fatalf("FindByID() error = %v, want nil", err)
				}
				if got, want := found.Email().String(), user.Email().String(); got != want {
					t.Fatalf("Email() got = %q, want %q", got, want)
				}
			},
		},
		{
			name: "SaveDuplicateReturnsErrDuplicateEmail",
			run: func(t *testing.T) {
				user1 := makeUser(t, "duplicate@teste.com")
				if err := repo.Save(context.Background(), user1); err != nil {
					t.Fatalf("Save() error = %v, want nil", err)
				}

				user2 := makeUser(t, "duplicate@teste.com")
				if err := repo.Save(context.Background(), user2); err != domainerrors.ErrDuplicateEmail {
					t.Fatalf("Save() error = %v, want %v", err, domainerrors.ErrDuplicateEmail)
				}
			},
		},
		{
			name: "FindByEmailMissingReturnsErrNotFound",
			run: func(t *testing.T) {
				email, err := valueobjects.NewEmail("nao@existe.com")
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				_, err = repo.FindByEmail(context.Background(), email)
				if err != domainerrors.ErrNotFound {
					t.Fatalf("FindByEmail() error = %v, want %v", err, domainerrors.ErrNotFound)
				}
			},
		},
		{
			name: "ListExistsDelete",
			run: func(t *testing.T) {
				for i := 0; i < 3; i++ {
					user := makeUser(t, fmt.Sprintf("user%c@teste.com", 'a'+i))
					if err := repo.Save(context.Background(), user); err != nil {
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

				if err := repo.Delete(context.Background(), list[0].ID()); err != nil {
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
