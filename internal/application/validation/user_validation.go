package validation

import (
	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

func ValidateCreateUser(req dto.CreateUserRequest) (email valueobjects.Email, role entities.Role, verr *ValidationErrors) {
	verr = &ValidationErrors{}

	var err error
	email, err = valueobjects.NewEmail(req.Email)
	if err != nil {
		verr.Add("email", err.Error())
	}

	role = entities.Role(req.Role)
	if !role.IsValid() {
		verr.Add("role", domainerrors.ErrInvalidRole.Error())
	}

	return email, role, verr
}
