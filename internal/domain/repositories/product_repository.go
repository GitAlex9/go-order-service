package repositories

import "github.com/GitAlex9/go-order-service/internal/domain/entities"

type ProductRepository interface {
	Save(product *entities.Product) error
	FindByID(id string) (*entities.Product, error)
	List() ([]*entities.Product, error)
	Exists(id string) (bool, error)
	Delete(id string) error
}
