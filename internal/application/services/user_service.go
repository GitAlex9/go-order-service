package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/commands"
	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/queries"
)

var _ contracts.UserService = (*userService)(nil)

type userService struct {
	createHandler         *commands.CreateUserHandler
	changePasswordHandler *commands.ChangePasswordHandler
	changeEmailHandler    *commands.ChangeUserEmailHandler
	activateHandler       *commands.ActivateUserHandler
	deactivateHandler     *commands.DeactivateUserHandler
	getHandler            *queries.GetUserHandler
	listHandler           *queries.ListUsersHandler
}

func NewUserService(
	createHandler *commands.CreateUserHandler,
	changePasswordHandler *commands.ChangePasswordHandler,
	changeEmailHandler *commands.ChangeUserEmailHandler,
	activateHandler *commands.ActivateUserHandler,
	deactivateHandler *commands.DeactivateUserHandler,
	getHandler *queries.GetUserHandler,
	listHandler *queries.ListUsersHandler,
) contracts.UserService {
	return &userService{
		createHandler:         createHandler,
		changePasswordHandler: changePasswordHandler,
		changeEmailHandler:    changeEmailHandler,
		activateHandler:       activateHandler,
		deactivateHandler:     deactivateHandler,
		getHandler:            getHandler,
		listHandler:           listHandler,
	}
}

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	return s.createHandler.Handle(ctx, req)
}
func (s *userService) Get(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	return s.getHandler.Handle(ctx, id)
}
func (s *userService) List(ctx context.Context, offset, limit int) ([]dto.UserResponse, error) {
	return s.listHandler.Handle(ctx, offset, limit)
}
func (s *userService) ChangePassword(ctx context.Context, id uuid.UUID, req dto.ChangePasswordRequest) error {
	return s.changePasswordHandler.Handle(ctx, id, req)
}
func (s *userService) ChangeEmail(ctx context.Context, id uuid.UUID, req dto.ChangeUserEmailRequest) (*dto.UserResponse, error) {
	return s.changeEmailHandler.Handle(ctx, id, req)
}
func (s *userService) Activate(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	return s.activateHandler.Handle(ctx, id)
}
func (s *userService) Deactivate(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	return s.deactivateHandler.Handle(ctx, id)
}
