package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

type ProductService interface {
	Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	Get(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error)
	List(ctx context.Context, offset, limit int) ([]dto.ProductResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	IncreaseStock(ctx context.Context, id uuid.UUID, req dto.AdjustStockRequest) (*dto.ProductResponse, error)
	DecreaseStock(ctx context.Context, id uuid.UUID, req dto.AdjustStockRequest) (*dto.ProductResponse, error)
	Activate(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error)
	Deactivate(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error)
}
