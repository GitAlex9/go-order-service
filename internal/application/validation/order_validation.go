package validation

import (
	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
)

func ValidateCreateOrder(req dto.CreateOrderRequest) (customerID uuid.UUID, verr *ValidationErrors) {
	verr = &ValidationErrors{}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		verr.Add("customer_id", domainerrors.ErrInvalidCustomer.Error())
	}

	if len(req.Items) == 0 {
		verr.Add("items", domainerrors.ErrEmptyOrder.Error())
	}

	for i, item := range req.Items {
		if item.Quantity <= 0 {
			verr.Add("items", domainerrors.ErrInvalidQuantity.Error())
			break // um erro de quantidade já basta pra sinalizar, evita spam repetido
		}
		if _, err := uuid.Parse(item.ProductID); err != nil {
			verr.Add("items", domainerrors.ErrInvalidProductID.Error())
			break
		}
		_ = i
	}

	return customerID, verr
}
