package repositories

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Save(ctx context.Context, product *entities.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	List(ctx context.Context, offset, limit int) ([]*entities.Product, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
