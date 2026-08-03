package response

import (
	"errors"
	"net/http"

	"github.com/GitAlex9/go-order-service/internal/application/validation"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

// HandleError traduz qualquer erro vindo de application/domain para uma resposta HTTP.
// converte para um status code equivalente.
func HandleError(w http.ResponseWriter, err error) {
	var verr *validation.ValidationErrors
	if errors.As(err, &verr) {
		fields := make([]FieldItem, len(verr.Errors))
		for i, fe := range verr.Errors {
			fields[i] = FieldItem{Field: fe.Field, Message: fe.Message}
		}
		JSONError(w, http.StatusUnprocessableEntity, "validation failed", fields)
		return
	}

	status, message := mapDomainError(err)
	JSONError(w, status, message, nil)
}

func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, domainerrors.ErrNotFound),
		errors.Is(err, domainerrors.ErrProductNotFound),
		errors.Is(err, domainerrors.ErrOrderNotFound),
		errors.Is(err, domainerrors.ErrUserNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, domainerrors.ErrDuplicateCustomer),
		errors.Is(err, domainerrors.ErrDuplicateEmail),
		errors.Is(err, domainerrors.ErrProductInUse):
		return http.StatusConflict, err.Error()

	case errors.Is(err, domainerrors.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()

	case errors.Is(err, domainerrors.ErrInsufficientStock),
		errors.Is(err, domainerrors.ErrInactiveProduct),
		errors.Is(err, domainerrors.ErrInvalidStatusTransition),
		errors.Is(err, domainerrors.ErrOrderNotEditable),
		errors.Is(err, domainerrors.ErrOrderNotDeletable),
		errors.Is(err, domainerrors.ErrEmptyOrder),
		errors.Is(err, domainerrors.ErrOrderItemNotFound):
		return http.StatusUnprocessableEntity, err.Error()

	case errors.Is(err, domainerrors.ErrInvalidEmail),
		errors.Is(err, domainerrors.ErrInvalidCPF),
		errors.Is(err, domainerrors.ErrInsufficientCPFLength),
		errors.Is(err, domainerrors.ErrNegativeMoneyAmount),
		errors.Is(err, domainerrors.ErrWeakPassword),
		errors.Is(err, domainerrors.ErrInvalidRole):
		return http.StatusBadRequest, err.Error()

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
