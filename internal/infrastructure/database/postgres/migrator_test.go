package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

type mockExecer struct {
	calls    []string
	failWhen func(query string) error
}

func (m *mockExecer) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	m.calls = append(m.calls, sql)
	if m.failWhen != nil {
		if err := m.failWhen(sql); err != nil {
			return pgconn.CommandTag{}, err
		}
	}
	return pgconn.CommandTag{}, nil
}

func TestMigrator_Migrate_Success(t *testing.T) {
	mock := &mockExecer{}
	migrator := NewMigrator(mock)

	if err := migrator.Migrate(); err != nil {
		t.Fatalf("Migrate() error = %v, want nil", err)
	}

	wantOrder := []string{
		"CREATE TABLE IF NOT EXISTS users",
		"CREATE TABLE IF NOT EXISTS products",
		"CREATE TABLE IF NOT EXISTS customers",
		"CREATE TABLE IF NOT EXISTS orders",
		"CREATE TABLE IF NOT EXISTS order_items",
	}

	if got, want := len(mock.calls), len(wantOrder); got != want {
		t.Fatalf("Exec called %d times, want %d", got, want)
	}

	for i, want := range wantOrder {
		if got := mock.calls[i]; !strings.Contains(got, want) {
			t.Errorf("call %d got query containing %q, want it to contain %q", i, got, want)
		}
	}
}

func TestMigrator_Migrate_StopsOnFirstError(t *testing.T) {
	wantErr := errors.New("connection failed")

	mock := &mockExecer{
		failWhen: func(query string) error {
			if strings.Contains(query, "CREATE TABLE IF NOT EXISTS customers") {
				return wantErr
			}
			return nil
		},
	}

	migrator := NewMigrator(mock)
	err := migrator.Migrate()

	if err == nil {
		t.Fatal("Migrate() error = nil, want an error")
	}

	wantMsg := "creating table customers: connection failed"
	if got := err.Error(); got != wantMsg {
		t.Errorf("Migrate() error = %q, want %q", got, wantMsg)
	}

	if got, want := len(mock.calls), 3; got != want {
		t.Errorf("Exec called %d times, want %d (deve parar no primeiro erro)", got, want)
	}
}
