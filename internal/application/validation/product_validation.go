package validation

import (
	"strings"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

func ValidateCreateProduct(req dto.CreateProductRequest) (name, description string, price valueobjects.Money, stock int, verr *ValidationErrors) {
	verr = &ValidationErrors{}

	name = strings.TrimSpace(req.Name)
	if name == "" {
		verr.Add("name", domainerrors.ErrInvalidProductName.Error())
	}

	description = strings.TrimSpace(req.Description)
	if description == "" {
		verr.Add("description", domainerrors.ErrInvalidProductDescription.Error())
	}

	var err error
	price, err = valueobjects.NewMoneyFromFloat(req.Price)
	if err != nil {
		verr.Add("price", err.Error())
	} else if price.IsZero() {
		verr.Add("price", domainerrors.ErrInvalidProductPrice.Error())
	}

	stock = req.Stock
	if stock < 0 {
		verr.Add("stock", domainerrors.ErrInvalidProductStock.Error())
	}

	return name, description, price, stock, verr
}
