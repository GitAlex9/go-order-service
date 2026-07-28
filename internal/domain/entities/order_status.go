package entities

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "PENDING"
	OrderStatusPaid     OrderStatus = "PAID"
	OrderStatusCanceled OrderStatus = "CANCELED"
)

var validOrderStatuses = map[OrderStatus]bool{
	OrderStatusPending:  true,
	OrderStatusPaid:     true,
	OrderStatusCanceled: true,
}

func (s OrderStatus) IsValid() bool {
	return validOrderStatuses[s]
}

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	transitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending: {OrderStatusPaid, OrderStatusCanceled},
	}
	for _, allowed := range transitions[s] {
		if allowed == target {
			return true
		}
	}
	return false
}
