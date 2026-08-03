package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMigrator_Migrate(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	if err != nil {
		t.Skip("Skipping test: testcontainers not available")
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	migrator := NewMigrator(pool)
	if err := migrator.Migrate(); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}

	tables := []struct {
		name        string
		shouldExist bool
	}{
		{"users", true},
		{"products", true},
		{"customers", true},
		{"orders", true},
		{"order_items", true},
		{"nonexistent", false},
	}

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			var exists bool
			query := "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)"
			if err := pool.QueryRow(ctx, query, tt.name).Scan(&exists); err != nil {
				t.Fatalf("failed to check table: %v", err)
			}
			if exists != tt.shouldExist {
				t.Errorf("table %s exists = %v, want %v", tt.name, exists, tt.shouldExist)
			}
		})
	}
}
