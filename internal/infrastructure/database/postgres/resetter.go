package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Resetter existe apenas para uso em desenvolvimento/testes.
// APAGAR ESTA STRUCT SE FOR USAR EM UM PROJETO REAL, CORRE O RISCO DE APAGAR TODO O BANCO.
type Resetter struct {
	pool *pgxpool.Pool
}

func NewResetter(pool *pgxpool.Pool) *Resetter {
	return &Resetter{pool: pool}
}

// Reset apaga todas as tabelas da aplicação, na ordem inversa da criação
// (respeitando as foreign keys) e recria a partir do Migrator, se fornecido.
//
// Trava de segurança: só executa se a variável de ambiente APP_ENV
// estiver explicitamente definida como "development" ou "test".
func (r *Resetter) Reset(ctx context.Context) error {
	env := os.Getenv("APP_ENV")
	if env != "development" && env != "test" {
		return fmt.Errorf("refusing to reset database: APP_ENV is %q, expected 'development' or 'test'", env)
	}

	tables := []string{
		"order_items", // depende de orders e products
		"orders",      // depende de customers
		"customers",   // depende de users
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
