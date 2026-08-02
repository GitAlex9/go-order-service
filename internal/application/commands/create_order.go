package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/application/validation"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type CreateOrderHandler struct {
	orderRepo    repositories.OrderRepository
	productRepo  repositories.ProductRepository
	customerRepo repositories.CustomerRepository
}

func NewCreateOrderHandler(
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	customerRepo repositories.CustomerRepository,
) *CreateOrderHandler {
	return &CreateOrderHandler{orderRepo: orderRepo, productRepo: productRepo, customerRepo: customerRepo}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	customerID, verr := validation.ValidateCreateOrder(req)
	if verr.HasErrors() {
		return nil, verr
	}

	if _, err := h.customerRepo.FindByID(ctx, customerID); err != nil {
		return nil, err
	}

	items := make([]entities.OrderItem, 0, len(req.Items))
	reservedProducts := make([]struct {
		product  *entities.Product
		quantity int
	}, 0, len(req.Items))

	for _, itemReq := range req.Items {
		productID, _ := uuid.Parse(itemReq.ProductID)

		product, err := h.productRepo.FindByID(ctx, productID)
		if err != nil {
			return nil, err
		}
		if !product.IsActive() {
			return nil, domainerrors.ErrInactiveProduct
		}
		if err := product.DecreaseStock(itemReq.Quantity); err != nil {
			return nil, err
		}

		item, err := entities.NewOrderItem(product.ID(), product.Name(), product.Price(), itemReq.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
		reservedProducts = append(reservedProducts, struct {
			product  *entities.Product
			quantity int
		}{product, itemReq.Quantity})
	}

	order, err := entities.NewOrder(customerID, items)
	if err != nil {
		return nil, err
	}

	if err := h.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	//Será resolvida com UnitOfWork
	for _, rp := range reservedProducts {
		if err := h.productRepo.Save(ctx, rp.product); err != nil {
			return nil, err
		}
	}

	response := mapper.OrderToResponse(order)
	return &response, nil
}
