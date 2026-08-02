package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type DeleteOrderHandler struct {
	repo repositories.OrderRepository
}

func NewDeleteOrderHandler(repo repositories.OrderRepository) *DeleteOrderHandler {
	return &DeleteOrderHandler{repo: repo}
}

func (h *DeleteOrderHandler) Handle(ctx context.Context, id uuid.UUID) error {
	order, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status() == entities.OrderStatusPaid {
		return domainerrors.ErrOrderNotDeletable
	}

	return h.repo.Delete(ctx, id)
}
