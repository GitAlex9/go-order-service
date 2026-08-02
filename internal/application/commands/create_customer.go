package commands

import (
	"context"
	"errors"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/application/validation"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type CreateCustomerHandler struct {
	repo repositories.CustomerRepository
}

func NewCreateCustomerHandler(repo repositories.CustomerRepository) *CreateCustomerHandler {
	return &CreateCustomerHandler{repo: repo}
}

func (h *CreateCustomerHandler) Handle(ctx context.Context, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error) {
	name, email, cpf, verr := validation.ValidateCreateCustomer(req)
	if verr.HasErrors() {
		return nil, verr
	}

	existing, err := h.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, domainerrors.ErrDuplicateCustomer
	}

	customer, err := entities.NewCustomer(name, email, cpf)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Save(ctx, customer); err != nil {
		return nil, err
	}

	response := mapper.CustomerToResponse(customer)
	return &response, nil
}
