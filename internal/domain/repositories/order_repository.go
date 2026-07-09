package repositories

import "github.com/GitAlex9/go-order-service/internal/domain/entities"

type OrderRepository interface {
	Save(order *entities.Order) error
	FindByID(id string) (*entities.Order, error)
	List() ([]*entities.Order, error)
	Delete(id string) error
}
