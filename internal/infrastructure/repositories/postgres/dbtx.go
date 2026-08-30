package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX é satisfeito tanto por *pgxpool.Pool quanto por pgx.Tx — permite que
// cada repository funcione tanto de forma avulsa (fora de uma transação)
// quanto dentro de uma transação aberta pelo UnitOfWork, sem duplicar código.
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
