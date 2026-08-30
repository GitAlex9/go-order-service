package postgres

import (
	"context"
	"fmt"
)

type Migrator struct {
	pool Execer
}

func NewMigrator(pool Execer) *Migrator {
	return &Migrator{pool: pool}
}

func (m *Migrator) Migrate() error {

	migrations := []struct {
		name  string
		query string
	}{
		{"users", createUsersTable},
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

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash VARCHAR(255) NOT NULL,
	role VARCHAR(20) NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`

const createProductsTable = `
CREATE TABLE IF NOT EXISTS products (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
	stock INTEGER NOT NULL CHECK (stock >= 0),
	active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`

const createCustomersTable = `
CREATE TABLE IF NOT EXISTS customers (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,
	cpf VARCHAR(11) NOT NULL UNIQUE,
	user_id UUID REFERENCES users(id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`

const createOrdersTable = `
CREATE TABLE IF NOT EXISTS orders (
	id UUID PRIMARY KEY,
	customer_id UUID NOT NULL,
	status VARCHAR(20) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

	CONSTRAINT fk_orders_customer
		FOREIGN KEY (customer_id)
		REFERENCES customers(id)
);
`

const createOrderItemsTable = `
CREATE TABLE IF NOT EXISTS order_items (
	order_id UUID NOT NULL,
	product_id UUID NOT NULL,
	product_name VARCHAR(255) NOT NULL,
	unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
	quantity INTEGER NOT NULL CHECK (quantity > 0),

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
