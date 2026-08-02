package dto

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CreateOrderRequest struct {
	CustomerID string             `json:"customer_id"`
	Items      []OrderItemRequest `json:"items"`
}

type AddItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
