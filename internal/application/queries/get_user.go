package queries

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type GetUserHandler struct {
	repo repositories.UserRepository
}

func NewGetUserHandler(repo repositories.UserRepository) *GetUserHandler {
	return &GetUserHandler{repo: repo}
}

func (h *GetUserHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := mapper.UserToResponse(user)
	return &response, nil
}
