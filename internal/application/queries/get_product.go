package queries

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type GetProductHandler struct {
	repo repositories.ProductRepository
}

func NewGetProductHandler(repo repositories.ProductRepository) *GetProductHandler {
	return &GetProductHandler{repo: repo}
}

func (h *GetProductHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
	product, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := mapper.ProductToResponse(product)
	return &response, nil
}
