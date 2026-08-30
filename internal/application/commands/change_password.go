package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	domainevents "github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type ChangePasswordHandler struct {
	repo       repositories.UserRepository
	dispatcher domainevents.Dispatcher
}

func NewChangePasswordHandler(repo repositories.UserRepository, dispatcher domainevents.Dispatcher) *ChangePasswordHandler {
	return &ChangePasswordHandler{repo: repo, dispatcher: dispatcher}
}

func (h *ChangePasswordHandler) Handle(ctx context.Context, id uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := user.ChangePassword(req.CurrentPassword, req.NewPassword); err != nil {
		return err
	}
	if err := h.repo.Save(ctx, user); err != nil {
		return err
	}

	h.dispatcher.Dispatch(ctx, user.Events())
	user.ClearEvents()

	return nil
}
