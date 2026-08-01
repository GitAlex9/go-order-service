package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type CancelOrderHandler struct {
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
}

func NewCancelOrderHandler(orderRepo repositories.OrderRepository, productRepo repositories.ProductRepository) *CancelOrderHandler {
	return &CancelOrderHandler{orderRepo: orderRepo, productRepo: productRepo}
}

func (h *CancelOrderHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	order, err := h.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	for _, item := range order.Items() {
		product, err := h.productRepo.FindByID(ctx, item.ProductID())
		if err != nil {
			return nil, err
		}
		if err := product.IncreaseStock(item.Quantity()); err != nil {
			return nil, err
		}
		if err := h.productRepo.Save(ctx, product); err != nil {
			return nil, err
		}
	}

	if err := h.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	response := mapper.OrderToResponse(order)
	return &response, nil
}
