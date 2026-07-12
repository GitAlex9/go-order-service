package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migrator struct {
	pool *pgxpool.Pool
}

func NewMigrator(pool *pgxpool.Pool) *Migrator {
	return &Migrator{
		pool: pool,
	}
}

func (m *Migrator) Migrate() error {

	migrations := []struct {
		name  string
		query string
	}{
		{"products", createProductsTable},
		{"customers", createCustomersTable},
		{"orders", createOrdersTable},
		{"order_items", createOrderItemsTable},
	}

	for _, migration := range migrations {

		if _, err := m.pool.Exec(
			context.Background(),
			migration.query,
		); err != nil {

			return fmt.Errorf(
				"creating table %s: %w",
				migration.name,
				err,
			)
		}
	}

	return nil
}

const createProductsTable = `
CREATE TABLE IF NOT EXISTS products (
	id VARCHAR(20) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	price NUMERIC(10,2) NOT NULL,
	stock INTEGER NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);
`

const createCustomersTable = `
CREATE TABLE IF NOT EXISTS customers (
	id VARCHAR(20) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);
`

const createOrdersTable = `
CREATE TABLE IF NOT EXISTS orders (
	id VARCHAR(20) PRIMARY KEY,
	customer_id VARCHAR(20) NOT NULL,
	status VARCHAR(20) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,

	CONSTRAINT fk_orders_customer
		FOREIGN KEY (customer_id)
		REFERENCES customers(id)
);
`

const createOrderItemsTable = `
CREATE TABLE IF NOT EXISTS order_items (
	order_id VARCHAR(20) NOT NULL,
	product_id VARCHAR(20) NOT NULL,
	product_name VARCHAR(255) NOT NULL,
	unit_price NUMERIC(10,2) NOT NULL,
	quantity INTEGER NOT NULL,

	PRIMARY KEY (order_id, product_id),

	CONSTRAINT fk_order_items_order
		FOREIGN KEY (order_id)
		REFERENCES orders(id)
		ON DELETE CASCADE,

	CONSTRAINT fk_order_items_product
		FOREIGN KEY (product_id)
		REFERENCES products(id)
);
`
