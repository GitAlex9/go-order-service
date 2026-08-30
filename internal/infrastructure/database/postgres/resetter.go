package postgres

import (
	"context"
	"fmt"
	"os"
)

type Resetter struct {
	pool Execer
}

func NewResetter(pool Execer) *Resetter {
	return &Resetter{pool: pool}
}

func (r *Resetter) Reset(ctx context.Context) error {
	env := os.Getenv("APP_ENV")
	if env != "development" && env != "test" {
		return fmt.Errorf("refusing to reset database: APP_ENV is %q, expected 'development' or 'test'", env)
	}

	tables := []string{
		"order_items",
		"orders",
		"customers",
		"products",
		"users",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		if _, err := r.pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("dropping table %s: %w", table, err)
		}
	}

	return nil
}
