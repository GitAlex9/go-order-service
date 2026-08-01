package queries

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type GetOrderHandler struct {
	repo repositories.OrderRepository
}

func NewGetOrderHandler(repo repositories.OrderRepository) *GetOrderHandler {
	return &GetOrderHandler{repo: repo}
}

func (h *GetOrderHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	order, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := mapper.OrderToResponse(order)
	return &response, nil
}
