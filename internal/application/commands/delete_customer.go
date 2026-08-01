package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type DeleteCustomerHandler struct {
	repo repositories.CustomerRepository
}

func NewDeleteCustomerHandler(repo repositories.CustomerRepository) *DeleteCustomerHandler {
	return &DeleteCustomerHandler{repo: repo}
}

func (h *DeleteCustomerHandler) Handle(ctx context.Context, id uuid.UUID) error {
	return h.repo.Delete(ctx, id)
}
