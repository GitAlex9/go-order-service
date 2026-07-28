package repositories

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	Save(ctx context.Context, customer *entities.Customer) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Customer, error)
	FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.Customer, error)
	List(ctx context.Context, offset, limit int) ([]*entities.Customer, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
