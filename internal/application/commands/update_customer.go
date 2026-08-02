package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

type UpdateCustomerHandler struct {
	repo repositories.CustomerRepository
}

func NewUpdateCustomerHandler(repo repositories.CustomerRepository) *UpdateCustomerHandler {
	return &UpdateCustomerHandler{repo: repo}
}

func (h *UpdateCustomerHandler) Handle(ctx context.Context, id uuid.UUID, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error) {
	customer, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		if err := customer.Rename(req.Name); err != nil {
			return nil, err
		}
	}

	if req.Email != "" {
		email, err := valueobjects.NewEmail(req.Email)
		if err != nil {
			return nil, err
		}
		customer.ChangeEmail(email)
	}

	if err := h.repo.Save(ctx, customer); err != nil {
		return nil, err
	}

	response := mapper.CustomerToResponse(customer)
	return &response, nil
}
