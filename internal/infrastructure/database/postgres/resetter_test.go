package postgres

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestResetter_Reset_RefusesOutsideDevOrTest(t *testing.T) {
	tests := []struct {
		name    string
		appEnv  string
		wantErr string
	}{
		{"production é bloqueado", "production", `refusing to reset database: APP_ENV is "production", expected 'development' or 'test'`},
		{"ambiente vazio é bloqueado", "", `refusing to reset database: APP_ENV is "", expected 'development' or 'test'`},
		{"staging é bloqueado", "staging", `refusing to reset database: APP_ENV is "staging", expected 'development' or 'test'`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.appEnv)

			mock := &mockExecer{}
			resetter := NewResetter(mock)

			err := resetter.Reset(context.Background())

			if err == nil {
				t.Fatal("Reset() error = nil, want an error")
			}
			if got := err.Error(); got != tt.wantErr {
				t.Errorf("Reset() error = %q, want %q", got, tt.wantErr)
			}

			if got := len(mock.calls); got != 0 {
				t.Errorf("Exec called %d times, want 0 (checagem de ambiente deve vir antes de qualquer DROP)", got)
			}
		})
	}
}

func TestResetter_Reset_AllowsDevelopmentAndTest(t *testing.T) {
	for _, env := range []string{"development", "test"} {
		t.Run(env, func(t *testing.T) {
			t.Setenv("APP_ENV", env)

			mock := &mockExecer{}
			resetter := NewResetter(mock)

			if err := resetter.Reset(context.Background()); err != nil {
				t.Fatalf("Reset() error = %v, want nil", err)
			}

			wantOrder := []string{"order_items", "orders", "customers", "products", "users"}

			if got, want := len(mock.calls), len(wantOrder); got != want {
				t.Fatalf("Exec called %d times, want %d", got, want)
			}

			for i, table := range wantOrder {
				want := "DROP TABLE IF EXISTS " + table
				if got := mock.calls[i]; !strings.Contains(got, want) {
					t.Errorf("call %d got query containing %q, want it to contain %q", i, got, want)
				}
			}
		})
	}
}

func TestResetter_Reset_StopsOnFirstError(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	wantErr := errors.New("connection failed")

	mock := &mockExecer{
		failWhen: func(query string) error {
			if strings.Contains(query, "DROP TABLE IF EXISTS customers") {
				return wantErr
			}
			return nil
		},
	}

	resetter := NewResetter(mock)
	err := resetter.Reset(context.Background())

	if err == nil {
		t.Fatal("Reset() error = nil, want an error")
	}

	wantMsg := "dropping table customers: connection failed"
	if got := err.Error(); got != wantMsg {
		t.Errorf("Reset() error = %q, want %q", got, wantMsg)
	}

	if got, want := len(mock.calls), 3; got != want {
		t.Errorf("Exec called %d times, want %d (deve parar no primeiro erro)", got, want)
	}
}
