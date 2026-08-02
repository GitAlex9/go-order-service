package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

type ChangeUserEmailHandler struct {
	repo repositories.UserRepository
}

func NewChangeUserEmailHandler(repo repositories.UserRepository) *ChangeUserEmailHandler {
	return &ChangeUserEmailHandler{repo: repo}
}

func (h *ChangeUserEmailHandler) Handle(ctx context.Context, id uuid.UUID, req dto.ChangeUserEmailRequest) (*dto.UserResponse, error) {
	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	email, err := valueobjects.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}
	user.ChangeEmail(email)

	if err := h.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	response := mapper.UserToResponse(user)
	return &response, nil
}
