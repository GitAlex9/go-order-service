package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/commands"
	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/queries"
)

var _ contracts.ProductService = (*productService)(nil)

type productService struct {
	createHandler        *commands.CreateProductHandler
	updateHandler        *commands.UpdateProductHandler
	deleteHandler        *commands.DeleteProductHandler
	increaseStockHandler *commands.IncreaseStockHandler
	decreaseStockHandler *commands.DecreaseStockHandler
	activateHandler      *commands.ActivateProductHandler
	deactivateHandler    *commands.DeactivateProductHandler
	getHandler           *queries.GetProductHandler
	listHandler          *queries.ListProductsHandler
}

func NewProductService(
	createHandler *commands.CreateProductHandler,
	updateHandler *commands.UpdateProductHandler,
	deleteHandler *commands.DeleteProductHandler,
	increaseStockHandler *commands.IncreaseStockHandler,
	decreaseStockHandler *commands.DecreaseStockHandler,
	activateHandler *commands.ActivateProductHandler,
	deactivateHandler *commands.DeactivateProductHandler,
	getHandler *queries.GetProductHandler,
	listHandler *queries.ListProductsHandler,
) contracts.ProductService {
	return &productService{
		createHandler:        createHandler,
		updateHandler:        updateHandler,
		deleteHandler:        deleteHandler,
		increaseStockHandler: increaseStockHandler,
		decreaseStockHandler: decreaseStockHandler,
		activateHandler:      activateHandler,
		deactivateHandler:    deactivateHandler,
		getHandler:           getHandler,
		listHandler:          listHandler,
	}
}

func (s *productService) Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	return s.createHandler.Handle(ctx, req)
}

func (s *productService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	return s.updateHandler.Handle(ctx, id, req)
}

func (s *productService) Get(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	return s.getHandler.Handle(ctx, id)
}

func (s *productService) List(ctx context.Context, offset, limit int) ([]dto.ProductResponse, error) {
	return s.listHandler.Handle(ctx, offset, limit)
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteHandler.Handle(ctx, id)
}

func (s *productService) IncreaseStock(ctx context.Context, id uuid.UUID, req dto.AdjustStockRequest) (*dto.ProductResponse, error) {
	return s.increaseStockHandler.Handle(ctx, id, req)
}

func (s *productService) DecreaseStock(ctx context.Context, id uuid.UUID, req dto.AdjustStockRequest) (*dto.ProductResponse, error) {
	return s.decreaseStockHandler.Handle(ctx, id, req)
}

func (s *productService) Activate(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	return s.activateHandler.Handle(ctx, id)
}

func (s *productService) Deactivate(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	return s.deactivateHandler.Handle(ctx, id)
}
