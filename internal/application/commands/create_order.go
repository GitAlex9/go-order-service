package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/application/validation"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	domainevents "github.com/GitAlex9/go-order-service/internal/domain/events"
)

type CreateOrderHandler struct {
	uow        contracts.UnitOfWork
	dispatcher domainevents.Dispatcher
}

func NewCreateOrderHandler(uow contracts.UnitOfWork, dispatcher domainevents.Dispatcher) *CreateOrderHandler {
	return &CreateOrderHandler{uow: uow, dispatcher: dispatcher}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	customerID, verr := validation.ValidateCreateOrder(req)
	if verr.HasErrors() {
		return nil, verr
	}

	var order *entities.Order
	var productEvents []any

	err := h.uow.Execute(ctx, func(repos contracts.Repositories) error {
		if _, err := repos.Customer.FindByID(ctx, customerID); err != nil {
			return err
		}

		items := make([]entities.OrderItem, 0, len(req.Items))

		for _, itemReq := range req.Items {
			productID, _ := uuid.Parse(itemReq.ProductID)

			product, err := repos.Product.FindByID(ctx, productID)
			if err != nil {
				return err
			}
			if !product.IsActive() {
				return domainerrors.ErrInactiveProduct
			}
			if err := product.DecreaseStock(itemReq.Quantity); err != nil {
				return err
			}

			item, err := entities.NewOrderItem(product.ID(), product.Name(), product.Price(), itemReq.Quantity)
			if err != nil {
				return err
			}
			items = append(items, *item)

			if err := repos.Product.Save(ctx, product); err != nil {
				return err
			}

			productEvents = append(productEvents, product.Events()...)
			product.ClearEvents()
		}

		newOrder, err := entities.NewOrder(customerID, items)
		if err != nil {
			return err
		}

		if err := repos.Order.Save(ctx, newOrder); err != nil {
			return err
		}

		order = newOrder
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
