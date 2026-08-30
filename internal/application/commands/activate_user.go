package commands

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/mapper"
	domainevents "github.com/GitAlex9/go-order-service/internal/domain/events"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

type ActivateUserHandler struct {
	repo       repositories.UserRepository
	dispatcher domainevents.Dispatcher
}

func NewActivateUserHandler(repo repositories.UserRepository, dispatcher domainevents.Dispatcher) *ActivateUserHandler {
	return &ActivateUserHandler{repo: repo, dispatcher: dispatcher}
}

func (h *ActivateUserHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Activate()
	if err := h.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	h.dispatcher.Dispatch(ctx, user.Events())
	user.ClearEvents()

	response := mapper.UserToResponse(user)
	return &response, nil
}

type DeactivateUserHandler struct {
	repo       repositories.UserRepository
	dispatcher domainevents.Dispatcher
}

func NewDeactivateUserHandler(repo repositories.UserRepository, dispatcher domainevents.Dispatcher) *DeactivateUserHandler {
	return &DeactivateUserHandler{repo: repo, dispatcher: dispatcher}
}

func (h *DeactivateUserHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Deactivate()
	if err := h.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	h.dispatcher.Dispatch(ctx, user.Events())
	user.ClearEvents()

	response := mapper.UserToResponse(user)
	return &response, nil
}
