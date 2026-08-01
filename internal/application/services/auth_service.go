package services

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/application/commands"
	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/dto"
)

var _ contracts.AuthService = (*authService)(nil)

type authService struct {
	authenticateHandler *commands.AuthenticateHandler
}

func NewAuthService(authenticateHandler *commands.AuthenticateHandler) contracts.AuthService {
	return &authService{authenticateHandler: authenticateHandler}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	return s.authenticateHandler.Handle(ctx, req)
}
