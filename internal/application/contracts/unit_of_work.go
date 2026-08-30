package contracts

import (
	"context"

	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
)

// Repositories agrupa instâncias de repository já vinculadas à mesma
// transação — tudo que for feito através delas dentro de um Execute
// confirma ou desfaz junto.
type Repositories struct {
	Customer repositories.CustomerRepository
	Product  repositories.ProductRepository
	Order    repositories.OrderRepository
	User     repositories.UserRepository
}

// UnitOfWork executa fn dentro de uma única transação. Se fn devolver erro,
// tudo é revertido; se fn concluir sem erro, tudo é confirmado junto.
type UnitOfWork interface {
	Execute(ctx context.Context, fn func(repos Repositories) error) error
}
