package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/GitAlex9/go-order-service/internal/infrastructure/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestDatabase struct {
	Pool     *pgxpool.Pool
	Migrator *postgres.Migrator
	Resetter *postgres.Resetter

	ctx       context.Context
	container *tcpostgres.PostgresContainer
}

func SetupTestDB(t *testing.T) *TestDatabase {
	t.Helper()

	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Skipf("docker/testcontainers unavailable: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		terminateContainer(t, ctx, container)
		t.Fatalf("getting postgres connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		terminateContainer(t, ctx, container)
		t.Fatalf("creating postgres pool: %v", err)
	}

	os.Setenv("APP_ENV", "test")

	db := &TestDatabase{
		Pool:      pool,
		Migrator:  postgres.NewMigrator(pool),
		Resetter:  postgres.NewResetter(pool),
		ctx:       ctx,
		container: container,
	}

	db.migrate(t)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func (db *TestDatabase) Reset(t *testing.T) {
	t.Helper()

	if err := db.Resetter.Reset(db.ctx); err != nil {
		t.Fatalf("reset database: %v", err)
	}

	db.migrate(t)
}

func (db *TestDatabase) migrate(t *testing.T) {
	t.Helper()

	if err := db.Migrator.Migrate(); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
}

func (db *TestDatabase) Close() {
	if db == nil {
		return
	}

	if db.Pool != nil {
		db.Pool.Close()
	}

	if db.container != nil {
		_ = db.container.Terminate(db.ctx)
	}

	_ = os.Unsetenv("APP_ENV")
}

func terminateContainer(
	t *testing.T,
	ctx context.Context,
	container *tcpostgres.PostgresContainer,
) {
	t.Helper()

	if container == nil {
		return
	}

	if err := container.Terminate(ctx); err != nil {
		t.Logf("terminating postgres container: %v", err)
	}
}
