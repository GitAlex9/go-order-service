package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type PayOrderHandler struct {
	repo repositories.OrderRepository
}

func NewPayOrderHandler(repo repositories.OrderRepository) *PayOrderHandler {
	return &PayOrderHandler{repo: repo}
}

func (h *PayOrderHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	order, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.Pay(); err != nil {
		return nil, err
	}

	if err := h.repo.Save(ctx, order); err != nil {
		return nil, err
	}

	response := mapper.OrderToResponse(order)
	return &response, nil
}
