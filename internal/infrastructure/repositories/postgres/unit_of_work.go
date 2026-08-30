package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

var _ contracts.UnitOfWork = (*unitOfWork)(nil)

type unitOfWork struct {
	pool *pgxpool.Pool
	log  logger.Logger
}

func NewUnitOfWork(pool *pgxpool.Pool, log logger.Logger) contracts.UnitOfWork {
	return &unitOfWork{pool: pool, log: log}
}

func (u *unitOfWork) Execute(ctx context.Context, fn func(repos contracts.Repositories) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op se já houve Commit

	repos := contracts.Repositories{
		Customer: NewCustomerRepository(tx),
		Product:  NewProductRepository(tx),
		Order:    NewOrderRepository(tx),
		User:     NewUserRepository(tx),
	}

	if err := fn(repos); err != nil {
		u.log.Warn("unit of work rolled back", "error", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		u.log.Error("unit of work commit failed", "error", err)
		return err
	}

	return nil
}
