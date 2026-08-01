package validation

import (
	"strings"

	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

func ValidateCreateCustomer(req dto.CreateCustomerRequest) (name string, email valueobjects.Email, cpf valueobjects.CPF, verr *ValidationErrors) {
	verr = &ValidationErrors{}

	name = strings.TrimSpace(req.Name)
	if len(name) < 3 {
		verr.Add("name", domainerrors.ErrEmptyName.Error())
	}

	var err error
	email, err = valueobjects.NewEmail(req.Email)
	if err != nil {
		verr.Add("email", err.Error())
	}

	cpf, err = valueobjects.NewCPF(req.CPF)
	if err != nil {
		verr.Add("cpf", err.Error())
	}

	return name, email, cpf, verr

}
