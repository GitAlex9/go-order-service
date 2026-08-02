package dto

import "time"

type OrderItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	CustomerID string              `json:"customer_id"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
	Total      float64             `json:"total"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}
