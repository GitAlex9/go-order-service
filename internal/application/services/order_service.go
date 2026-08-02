package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/commands"
	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/queries"
)

var _ contracts.OrderService = (*orderService)(nil)

type orderService struct {
	createHandler *commands.CreateOrderHandler
	payHandler    *commands.PayOrderHandler
	cancelHandler *commands.CancelOrderHandler
	deleteHandler *commands.DeleteOrderHandler
	getHandler    *queries.GetOrderHandler
	listHandler   *queries.ListOrdersHandler
}

func NewOrderService(
	createHandler *commands.CreateOrderHandler,
	payHandler *commands.PayOrderHandler,
	cancelHandler *commands.CancelOrderHandler,
	deleteHandler *commands.DeleteOrderHandler,
	getHandler *queries.GetOrderHandler,
	listHandler *queries.ListOrdersHandler,
) contracts.OrderService {
	return &orderService{
		createHandler: createHandler,
		payHandler:    payHandler,
		cancelHandler: cancelHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

func (s *orderService) Create(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	return s.createHandler.Handle(ctx, req)
}

func (s *orderService) Pay(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	return s.payHandler.Handle(ctx, id)
}

func (s *orderService) Cancel(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	return s.cancelHandler.Handle(ctx, id)
}

func (s *orderService) Get(ctx context.Context, id uuid.UUID) (*dto.OrderResponse, error) {
	return s.getHandler.Handle(ctx, id)
}

func (s *orderService) List(ctx context.Context, offset, limit int) ([]dto.OrderResponse, error) {
	return s.listHandler.Handle(ctx, offset, limit)
}

func (s *orderService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteHandler.Handle(ctx, id)
}
