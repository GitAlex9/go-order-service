package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	domainevents "github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

type ChangeUserEmailHandler struct {
	repo       repositories.UserRepository
	dispatcher domainevents.Dispatcher
}

func NewChangeUserEmailHandler(repo repositories.UserRepository, dispatcher domainevents.Dispatcher) *ChangeUserEmailHandler {
	return &ChangeUserEmailHandler{repo: repo, dispatcher: dispatcher}
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

	h.dispatcher.Dispatch(ctx, user.Events())
	user.ClearEvents()

	response := mapper.UserToResponse(user)
	return &response, nil
}
