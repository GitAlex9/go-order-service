package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type IncreaseStockHandler struct {
	repo repositories.ProductRepository
}

func NewIncreaseStockHandler(repo repositories.ProductRepository) *IncreaseStockHandler {
	return &IncreaseStockHandler{repo: repo}
}

func (h *IncreaseStockHandler) Handle(ctx context.Context, id uuid.UUID, req dto.AdjustStockRequest) (*dto.ProductResponse, error) {
	product, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := product.IncreaseStock(req.Quantity); err != nil {
		return nil, err
	}
	if err := h.repo.Save(ctx, product); err != nil {
		return nil, err
	}
	response := mapper.ProductToResponse(product)
	return &response, nil
}

type DecreaseStockHandler struct {
	repo repositories.ProductRepository
}

func NewDecreaseStockHandler(repo repositories.ProductRepository) *DecreaseStockHandler {
	return &DecreaseStockHandler{repo: repo}
}

func (h *DecreaseStockHandler) Handle(ctx context.Context, id uuid.UUID, req dto.AdjustStockRequest) (*dto.ProductResponse, error) {
	product, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := product.DecreaseStock(req.Quantity); err != nil {
		return nil, err
	}
	if err := h.repo.Save(ctx, product); err != nil {
		return nil, err
	}
	response := mapper.ProductToResponse(product)
	return &response, nil
}

type ActivateProductHandler struct {
	repo repositories.ProductRepository
}

func NewActivateProductHandler(repo repositories.ProductRepository) *ActivateProductHandler {
	return &ActivateProductHandler{repo: repo}
}

func (h *ActivateProductHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	product.Activate()
	if err := h.repo.Save(ctx, product); err != nil {
		return nil, err
	}
	response := mapper.ProductToResponse(product)
	return &response, nil
}

type DeactivateProductHandler struct {
	repo repositories.ProductRepository
}

func NewDeactivateProductHandler(repo repositories.ProductRepository) *DeactivateProductHandler {
	return &DeactivateProductHandler{repo: repo}
}

func (h *DeactivateProductHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	product.Deactivate()
	if err := h.repo.Save(ctx, product); err != nil {
		return nil, err
	}
	response := mapper.ProductToResponse(product)
	return &response, nil
}
