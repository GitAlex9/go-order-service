package repositories

import "github.com/GitAlex9/go-order-service/internal/domain/entities"

type CustomerRepository interface {
	Save(customer *entities.Customer) error
	FindByID(id string) (*entities.Customer, error)
	List() ([]*entities.Customer, error)
	Exists(id string) (bool, error)
	Delete(id string) error
}
