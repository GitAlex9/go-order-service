package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

type OrderService interface {
	Create(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	Pay(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	Cancel(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	Get(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error)
	List(ctx context.Context, offset, limit int) ([]dto.OrderResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
