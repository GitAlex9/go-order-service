package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestResetter_Reset(t *testing.T) {
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

	tests := []struct {
		name    string
		appEnv  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "development environment allowed",
			appEnv:  "development",
			wantErr: false,
		},
		{
			name:    "test environment allowed",
			appEnv:  "test",
			wantErr: false,
		},
		{
			name:    "production environment blocked",
			appEnv:  "production",
			wantErr: true,
			errMsg:  `refusing to reset database: APP_ENV is "production", expected 'development' or 'test'`,
		},
		{
			name:    "empty environment blocked",
			appEnv:  "",
			wantErr: true,
			errMsg:  `refusing to reset database: APP_ENV is "", expected 'development' or 'test'`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("APP_ENV", tt.appEnv)
			defer os.Unsetenv("APP_ENV")

			resetter := NewResetter(pool)
			err := resetter.Reset(ctx)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if err != nil && err.Error() != tt.errMsg {
					t.Errorf("error = %q, want %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			tables := []string{"order_items", "orders", "customers", "products", "users"}
			for _, table := range tables {
				var exists bool
				query := "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)"
				if err := pool.QueryRow(ctx, query, table).Scan(&exists); err != nil {
					t.Errorf("failed to check table %s: %v", table, err)
				}
				if exists {
					t.Errorf("table %s still exists after reset", table)
				}
			}
		})
	}
}
