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

type CancelOrderHandler struct {
	uow        contracts.UnitOfWork
	dispatcher domainevents.Dispatcher
}

func NewCancelOrderHandler(uow contracts.UnitOfWork, dispatcher domainevents.Dispatcher) *CancelOrderHandler {
	return &CancelOrderHandler{uow: uow, dispatcher: dispatcher}
}

func (h *CancelOrderHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	var order *entities.Order
	var productEvents []any

	err := h.uow.Execute(ctx, func(repos contracts.Repositories) error {
		o, err := repos.Order.FindByID(ctx, id)
		if err != nil {
			return err
		}

		if err := o.Cancel(); err != nil {
			return err
		}

		for _, item := range o.Items() {
			product, err := repos.Product.FindByID(ctx, item.ProductID())
			if err != nil {
				return err
			}
			if err := product.IncreaseStock(item.Quantity()); err != nil {
				return err
			}
			if err := repos.Product.Save(ctx, product); err != nil {
				return err
			}

			productEvents = append(productEvents, product.Events()...)
			product.ClearEvents()
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

	allEvents := append(order.Events(), productEvents...)
	h.dispatcher.Dispatch(ctx, allEvents)
	order.ClearEvents()

	response := mapper.OrderToResponse(order)
	return &response, nil
}
