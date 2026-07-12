package errors

import "errors"

var (
	// Product errors
	ErrProductNotFound = errors.New("product not found")

	ErrInvalidProductID          = errors.New("invalid product id")
	ErrInvalidProductName        = errors.New("invalid product name")
	ErrInvalidProductDescription = errors.New("invalid product description")
	ErrInvalidProductPrice       = errors.New("invalid product price")
	ErrInvalidProductStock       = errors.New("invalid product stock")

	ErrInactiveProduct   = errors.New("inactive product")
	ErrInsufficientStock = errors.New("insufficient stock")

	// Order errors
	ErrOrderNotFound           = errors.New("order not found")
	ErrEmptyOrder              = errors.New("order must contain at least one item")
	ErrInvalidStatusTransition = errors.New("invalid status transition")

	// Customer errors
	ErrInvalidCustomer = errors.New("invalid customer")
	ErrInvalidEmail    = errors.New("invalid email")

	// Generic validation
	ErrInvalidQuantity = errors.New("invalid quantity")
)
