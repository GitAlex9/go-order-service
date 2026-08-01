package queries

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type GetCustomerHandler struct {
	repo repositories.CustomerRepository
}

func NewGetCustomerHandler(repo repositories.CustomerRepository) *GetCustomerHandler {
	return &GetCustomerHandler{repo: repo}
}

func (h *GetCustomerHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.CustomerResponse, error) {
	customer, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := mapper.CustomerToResponse(customer)
	return &response, nil
}
