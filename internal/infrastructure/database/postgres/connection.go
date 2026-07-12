package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	pool *pgxpool.Pool
}

func NewConnection(cfg *Config) (*Connection, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("creating postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	return &Connection{
		pool: pool,
	}, nil
}

func (c *Connection) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *Connection) Close() {

	if c.pool != nil {
		c.pool.Close()
	}
}
