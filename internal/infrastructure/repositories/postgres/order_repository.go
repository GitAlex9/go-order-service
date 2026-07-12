package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repositories.OrderRepository = (*OrderRepositoryPostgres)(nil)
var total float64

type OrderRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepositoryPostgres {
	return &OrderRepositoryPostgres{
		pool: pool,
	}
}

func (r *OrderRepositoryPostgres) Delete(id string) error {
	query := `
		DELETE FROM orders
		WHERE id = $1
	`

	_, err := r.pool.Exec(context.Background(), query, id)
	return err
}

func (r *OrderRepositoryPostgres) List() ([]*entities.Order, error) {
	query := `
		SELECT
			id,
			customer_id,
			status,
			created_at,
			updated_at
		FROM orders
	`

	rows, err := r.pool.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []*entities.Order

	for rows.Next() {
		var (
			id         string
			customerID string
			status     string
			createdAt  time.Time
			updatedAt  time.Time
		)

		err := rows.Scan(
			&id,
			&customerID,
			&status,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		order := entities.RebuildOrder(
			id,
			customerID,
			entities.OrderStatus(status),
			nil,
			createdAt,
			updatedAt,
		)
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepositoryPostgres) Save(order *entities.Order) error {

	query := `
		INSERT INTO orders
		(
			id,
			customer_id,
			status,
			created_at,
			updated_at
		)
		VALUES
		(
			$1,$2,$3,$4,$5
		)
	`

	_, err := r.pool.Exec(
		context.Background(),
		query,
		order.ID,
		order.CustomerID,
		order.Status(),
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

func (r *OrderRepositoryPostgres) FindByID(id string) (*entities.Order, error) {

	query := `
		SELECT
			id,
			customer_id,
			status,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
	`

	row := r.pool.QueryRow(context.Background(), query, id)

	var order entities.Order
	var status string

	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	order.SetStatus(entities.OrderStatus(status))

	return &order, nil
}

func (r *OrderRepositoryPostgres) Exists(id string) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM orders
			WHERE id = $1
		)
	`

	var exists bool

	err := r.pool.QueryRow(context.Background(), query, id).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
