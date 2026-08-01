package queries

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type ListOrdersHandler struct {
	repo repositories.OrderRepository
}

func NewListOrdersHandler(repo repositories.OrderRepository) *ListOrdersHandler {
	return &ListOrdersHandler{repo: repo}
}

func (h *ListOrdersHandler) Handle(ctx context.Context, offset, limit int) ([]dto.OrderResponse, error) {
	orders, err := h.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return mapper.OrdersToResponse(orders), nil
}
