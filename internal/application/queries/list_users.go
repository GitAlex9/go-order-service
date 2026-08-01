package queries

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type ListUsersHandler struct {
	repo repositories.UserRepository
}

func NewListUsersHandler(repo repositories.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{repo: repo}
}

func (h *ListUsersHandler) Handle(ctx context.Context, offset, limit int) ([]dto.UserResponse, error) {
	users, err := h.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return mapper.UsersToResponse(users), nil
}
