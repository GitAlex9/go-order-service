package mapper

import (
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
)

func OrderToResponse(o *entities.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(o.Items()))
	for i, item := range o.Items() {
		items[i] = dto.OrderItemResponse{
			ProductID:   item.ProductID().String(),
			ProductName: item.ProductName(),
			UnitPrice:   item.UnitPrice().Amount(),
			Quantity:    item.Quantity(),
			Subtotal:    item.Subtotal().Amount(),
		}
	}

	return dto.OrderResponse{
		ID:         o.ID().String(),
		CustomerID: o.CustomerID().String(),
		Status:     o.Status().String(),
		Items:      items,
		Total:      o.Total().Amount(),
		CreatedAt:  o.CreatedAt(),
		UpdatedAt:  o.UpdatedAt(),
	}
}

func OrdersToResponse(orders []*entities.Order) []dto.OrderResponse {
	responses := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = OrderToResponse(o)
	}
	return responses
}
