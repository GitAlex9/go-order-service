package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainevents "github.com/GitAlex9/go-order-service/internal/domain/events"
)

type PayOrderHandler struct {
	uow        contracts.UnitOfWork
	dispatcher domainevents.Dispatcher
}

func NewPayOrderHandler(uow contracts.UnitOfWork, dispatcher domainevents.Dispatcher) *PayOrderHandler {
	return &PayOrderHandler{uow: uow, dispatcher: dispatcher}
}

func (h *PayOrderHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	var order *entities.Order

	err := h.uow.Execute(ctx, func(repos contracts.Repositories) error {
		o, err := repos.Order.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := o.Pay(); err != nil {
			return err
		}
		if err := repos.Order.Save(ctx, o); err != nil {
			return err
		}
		order = o
		return nil
	})
	if err != nil {
		return nil, err
	}

	h.dispatcher.Dispatch(ctx, order.Events())
	order.ClearEvents()

	response := mapper.OrderToResponse(order)
	return &response, nil
}
