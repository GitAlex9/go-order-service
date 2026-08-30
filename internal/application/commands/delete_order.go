package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

type DeleteOrderHandler struct {
	uow contracts.UnitOfWork
}

func NewDeleteOrderHandler(uow contracts.UnitOfWork) *DeleteOrderHandler {
	return &DeleteOrderHandler{uow: uow}
}

func (h *DeleteOrderHandler) Handle(ctx context.Context, id uuid.UUID) error {
	return h.uow.Execute(ctx, func(repos contracts.Repositories) error {
		order, err := repos.Order.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if order.Status() == entities.OrderStatusPaid {
			return domainerrors.ErrOrderNotDeletable
		}
		return repos.Order.Delete(ctx, id)
	})
}
