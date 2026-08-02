package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

type CustomerService interface {
	Create(ctx context.Context, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error)
	Get(ctx context.Context, id uuid.UUID) (*dto.CustomerResponse, error)
	List(ctx context.Context, offset, limit int) ([]dto.CustomerResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
