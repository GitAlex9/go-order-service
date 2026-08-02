package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

type UpdateProductHandler struct {
	repo repositories.ProductRepository
}

func NewUpdateProductHandler(repo repositories.ProductRepository) *UpdateProductHandler {
	return &UpdateProductHandler{repo: repo}
}

func (h *UpdateProductHandler) Handle(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		if err := product.Rename(req.Name); err != nil {
			return nil, err
		}
	}

	if req.Description != "" {
		if err := product.ChangeDescription(req.Description); err != nil {
			return nil, err
		}
	}

	if req.Price != 0 {
		price, err := valueobjects.NewMoneyFromFloat(req.Price)
		if err != nil {
			return nil, err
		}
		if err := product.ChangePrice(price); err != nil {
			return nil, err
		}
	}

	if err := h.repo.Save(ctx, product); err != nil {
		return nil, err
	}

	response := mapper.ProductToResponse(product)
	return &response, nil
}
