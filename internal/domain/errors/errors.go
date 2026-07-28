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
	ErrInvalidOrderStatus      = errors.New("invalid status")
	ErrOrderNotEditable        = errors.New("items cannot be modified for paid or cancelled orders")
	ErrOrderItemNotFound       = errors.New("attempted to remove a product that is not in the order")

	// Customer errors
	ErrInvalidCustomer = errors.New("invalid customer")
	ErrInvalidEmail    = errors.New("invalid email")

	// Generic validation
	ErrInvalidQuantity = errors.New("invalid quantity")

	// User errors
	ErrInvalidID    = errors.New("invalid ID")
	ErrEmptyName    = errors.New("Name cannot be empty")
	ErrWeakPassword = errors.New("The password must be at least 8 characters long")

	// Value object
	ErrInsufficientCPFLength = errors.New("CPF must contain exactly 11 digits")
	ErrInvalidCPF            = errors.New("invalid CPF")
	ErrNegativeMoneyAmount   = errors.New("money cannot be negative")

	// Password
	ErrPasswordNoNumber  = errors.New("password must contain at least one number")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character")

	// Domain
	ErrNotFound = errors.New("value not found")
)
