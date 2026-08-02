package mapper

import (
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
)

func CustomerToResponse(c *entities.Customer) dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:        c.ID().String(),
		Name:      c.Name(),
		Email:     c.Email().String(),
		CPF:       c.CPF().Formatted(),
		CreatedAt: c.CreatedAt(),
		UpdatedAt: c.UpdatedAt(),
	}
}

func CustomersToResponse(customers []*entities.Customer) []dto.CustomerResponse {
	responses := make([]dto.CustomerResponse, len(customers))
	for i, c := range customers {
		responses[i] = CustomerToResponse(c)
	}
	return responses
}
