package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/domain/repositories"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
)

var _ repositories.OrderRepository = (*OrderRepository)(nil)

type OrderRepository struct {
	db DBTX
}

func NewOrderRepository(db DBTX) *OrderRepository {
	return &OrderRepository{db: db}
}

type orderRow struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type orderItemRow struct {
	ProductID      uuid.UUID
	ProductName    string
	UnitPriceCents int64
	Quantity       int
}

// Save grava o pedido e seus itens. NÃO abre transação própria e espera ser
// chamado sempre dentro de um contexto já transacional, fornecido pelo
// UnitOfWork. É isso que garante a atomicidade entre orders e order_items
// (e, quando aplicável, entre Order e Product também).
func (r *OrderRepository) Save(ctx context.Context, order *entities.Order) error {
	const orderQuery = `
		INSERT INTO orders (id, customer_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, orderQuery,
		order.ID(), order.CustomerID(), order.Status().String(), order.CreatedAt(), order.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	if _, err := r.db.Exec(ctx, `DELETE FROM order_items WHERE order_id = $1`, order.ID()); err != nil {
		return err
	}

	const itemQuery = `
		INSERT INTO order_items (order_id, product_id, product_name, unit_price_cents, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, item := range order.Items() {
		_, err := r.db.Exec(ctx, itemQuery,
			order.ID(), item.ProductID(), item.ProductName(), item.UnitPrice().Cents(), item.Quantity(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	const orderQuery = `
		SELECT id, customer_id, status, created_at, updated_at
		FROM orders WHERE id = $1
	`
	var or orderRow
	err := r.db.QueryRow(ctx, orderQuery, id).
		Scan(&or.ID, &or.CustomerID, &or.Status, &or.CreatedAt, &or.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}

	items, err := r.findItemsByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}

	return orderToDomain(or, items)
}

func (r *OrderRepository) findItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]entities.OrderItem, error) {
	const query = `
		SELECT product_id, product_name, unit_price_cents, quantity
		FROM order_items WHERE order_id = $1
	`
	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entities.OrderItem, 0)
	for rows.Next() {
		var ir orderItemRow
		if err := rows.Scan(&ir.ProductID, &ir.ProductName, &ir.UnitPriceCents, &ir.Quantity); err != nil {
			return nil, err
		}
		price, err := valueobjects.NewMoney(ir.UnitPriceCents)
		if err != nil {
			return nil, err
		}
		items = append(items, *entities.RebuildOrderItem(ir.ProductID, ir.ProductName, price, ir.Quantity))
	}
	return items, rows.Err()
}

func (r *OrderRepository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.Order, error) {
	const query = `
		SELECT id, customer_id, status, created_at, updated_at
		FROM orders WHERE customer_id = $1 ORDER BY created_at DESC
	`
	return r.queryOrders(ctx, query, customerID)
}

func (r *OrderRepository) List(ctx context.Context, offset, limit int) ([]*entities.Order, error) {
	const query = `
		SELECT id, customer_id, status, created_at, updated_at
		FROM orders ORDER BY created_at DESC OFFSET $1 LIMIT $2
	`
	return r.queryOrders(ctx, query, offset, limit)
}

func (r *OrderRepository) queryOrders(ctx context.Context, query string, args ...any) ([]*entities.Order, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderRows []orderRow
	for rows.Next() {
		var or orderRow
		if err := rows.Scan(&or.ID, &or.CustomerID, &or.Status, &or.CreatedAt, &or.UpdatedAt); err != nil {
			return nil, err
		}
		orderRows = append(orderRows, or)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	orders := make([]*entities.Order, 0, len(orderRows))
	for _, or := range orderRows {
		items, err := r.findItemsByOrderID(ctx, or.ID)
		if err != nil {
			return nil, err
		}
		order, err := orderToDomain(or, items)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM orders WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

func orderToDomain(or orderRow, items []entities.OrderItem) (*entities.Order, error) {
	status := entities.OrderStatus(or.Status)
	if !status.IsValid() {
		return nil, domainerrors.ErrInvalidOrderStatus
	}
	return entities.RebuildOrder(or.ID, or.CustomerID, status, items, or.CreatedAt, or.UpdatedAt), nil
}
