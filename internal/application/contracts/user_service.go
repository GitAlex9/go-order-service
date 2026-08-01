package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

type UserService interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	Get(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	List(ctx context.Context, offset, limit int) ([]dto.UserResponse, error)
	ChangePassword(ctx context.Context, id uuid.UUID, req dto.ChangePasswordRequest) error
	ChangeEmail(ctx context.Context, id uuid.UUID, req dto.ChangeUserEmailRequest) (*dto.UserResponse, error)
	Activate(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	Deactivate(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
}
