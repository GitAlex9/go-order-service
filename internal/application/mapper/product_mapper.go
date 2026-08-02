package mapper

import (
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
)

func ProductToResponse(p *entities.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          p.ID().String(),
		Name:        p.Name(),
		Description: p.Description(),
		Price:       p.Price().Amount(),
		Stock:       p.Stock(),
		Active:      p.IsActive(),
		CreatedAt:   p.CreatedAt(),
		UpdatedAt:   p.UpdatedAt(),
	}
}

func ProductsToResponse(products []*entities.Product) []dto.ProductResponse {
	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ProductToResponse(p)
	}
	return responses
}
