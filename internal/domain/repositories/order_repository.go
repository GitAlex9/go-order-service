package repositories

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entities.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.Order, error)
	List(ctx context.Context, offset, limit int) ([]*entities.Order, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
