package commands

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/application/validation"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type CreateProductHandler struct {
	repo repositories.ProductRepository
}

func NewCreateProductHandler(repo repositories.ProductRepository) *CreateProductHandler {
	return &CreateProductHandler{repo: repo}
}

func (h *CreateProductHandler) Handle(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	name, description, price, stock, verr := validation.ValidateCreateProduct(req)
	if verr.HasErrors() {
		return nil, verr
	}

	product, err := entities.NewProduct(name, description, price, stock)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Save(ctx, product); err != nil {
		return nil, err
	}

	response := mapper.ProductToResponse(product)
	return &response, nil
}
