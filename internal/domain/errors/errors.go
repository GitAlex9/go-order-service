package errors

import "errors"

var (
	// ==================== PRODUCT ERRORS ====================
	ErrProductNotFound           = errors.New("product not found")
	ErrInvalidProductID          = errors.New("invalid product id")
	ErrInvalidProductName        = errors.New("invalid product name")
	ErrInvalidProductDescription = errors.New("invalid product description")
	ErrInvalidProductPrice       = errors.New("invalid product price")
	ErrInvalidProductStock       = errors.New("invalid product stock")
	ErrInactiveProduct           = errors.New("inactive product")
	ErrInsufficientStock         = errors.New("insufficient stock")
	ErrProductInUse              = errors.New("product cannot be deleted because it is referenced by existing orders; deactivate it instead")

	// ==================== ORDER ERRORS ====================
	ErrOrderNotFound           = errors.New("order not found")
	ErrEmptyOrder              = errors.New("order must contain at least one item")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidOrderStatus      = errors.New("invalid order status")
	ErrOrderNotEditable        = errors.New("items cannot be modified for paid or cancelled orders")
	ErrOrderItemNotFound       = errors.New("attempted to remove a product that is not in the order")
	ErrOrderNotDeletable       = errors.New("paid orders cannot be deleted, only cancelled orders or pending orders")

	// ==================== CUSTOMER ERRORS ====================
	ErrInvalidCustomer   = errors.New("invalid customer")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrDuplicateCustomer = errors.New("customer with this email or cpf already exists")

	// ==================== USER / AUTH ERRORS ====================
	ErrInvalidID          = errors.New("invalid ID")
	ErrEmptyName          = errors.New("name cannot be empty")
	ErrWeakPassword       = errors.New("the password must be at least 8 characters long")
	ErrInvalidRole        = errors.New("invalid role")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")

	// ==================== GENERIC VALIDATION ====================
	ErrInvalidQuantity = errors.New("invalid quantity")

	// ==================== VALUE OBJECTS ====================
	ErrInsufficientCPFLength = errors.New("CPF must contain exactly 11 digits")
	ErrInvalidCPF            = errors.New("invalid CPF")
	ErrNegativeMoneyAmount   = errors.New("money cannot be negative")

	// ==================== PASSWORD VALIDATION ====================
	ErrPasswordNoNumber         = errors.New("password must contain at least one number")
	ErrPasswordNoUpper          = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower          = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoSpecial        = errors.New("password must contain at least one special character")
	ErrIncorrectCurrentPassword = errors.New("current password is incorrect")

	// ==================== DOMAIN ====================
	ErrNotFound = errors.New("value not found")
)
