package mapper

import (
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
)

func UserToResponse(u *entities.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID().String(),
		Email:     u.Email().String(),
		Role:      u.Role().String(),
		Active:    u.Active(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}
}

func UsersToResponse(users []*entities.User) []dto.UserResponse {
	responses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		responses[i] = UserToResponse(u)
	}
	return responses
}
