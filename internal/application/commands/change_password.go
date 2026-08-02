package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type ChangePasswordHandler struct {
	repo repositories.UserRepository
}

func NewChangePasswordHandler(repo repositories.UserRepository) *ChangePasswordHandler {
	return &ChangePasswordHandler{repo: repo}
}

func (h *ChangePasswordHandler) Handle(ctx context.Context, id uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := user.ChangePassword(req.CurrentPassword, req.NewPassword); err != nil {
		return err
	}
	return h.repo.Save(ctx, user)
}
