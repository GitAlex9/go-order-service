package repositories

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error)
	List(ctx context.Context, offset, limit int) ([]*entities.User, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
