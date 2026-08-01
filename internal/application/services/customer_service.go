package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/commands"
	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/queries"
)

var _ contracts.CustomerService = (*customerService)(nil)

type customerService struct {
	createHandler *commands.CreateCustomerHandler
	updateHandler *commands.UpdateCustomerHandler
	deleteHandler *commands.DeleteCustomerHandler
	getHandler    *queries.GetCustomerHandler
	listHandler   *queries.ListCustomersHandler
}

func NewCustomerService(
	createHandler *commands.CreateCustomerHandler,
	updateHandler *commands.UpdateCustomerHandler,
	deleteHandler *commands.DeleteCustomerHandler,
	getHandler *queries.GetCustomerHandler,
	listHandler *queries.ListCustomersHandler,
) contracts.CustomerService {
	return &customerService{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

func (s *customerService) Create(ctx context.Context, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error) {
	return s.createHandler.Handle(ctx, req)
}

func (s *customerService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error) {
	return s.updateHandler.Handle(ctx, id, req)
}

func (s *customerService) Get(ctx context.Context, id uuid.UUID) (*dto.CustomerResponse, error) {
	return s.getHandler.Handle(ctx, id)
}

func (s *customerService) List(ctx context.Context, offset, limit int) ([]dto.CustomerResponse, error) {
	return s.listHandler.Handle(ctx, offset, limit)
}

func (s *customerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteHandler.Handle(ctx, id)
}
