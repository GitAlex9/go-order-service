package postgres

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    Config
	}{
		{
			name:    "usa valores padrão quando nenhuma env var está definida",
			envVars: map[string]string{},
			want: Config{
				Host: "localhost", Port: "5432", User: "postgres",
				Password: "postgres", Database: "go_order_service", SSLMode: "disable",
			},
		},
		{
			name: "usa valores das env vars quando definidas",
			envVars: map[string]string{
				"DB_HOST": "db.example.com", "DB_PORT": "6543", "DB_USER": "admin",
				"DB_PASSWORD": "secret", "DB_NAME": "prod_db", "DB_SSLMODE": "require",
			},
			want: Config{
				Host: "db.example.com", Port: "6543", User: "admin",
				Password: "secret", Database: "prod_db", SSLMode: "require",
			},
		},
	}

	keys := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, k := range keys {
				os.Unsetenv(k)
			}
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got := *NewConfig()
			if got != tt.want {
				t.Errorf("NewConfig() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestConfig_ConnectionString(t *testing.T) {
	cfg := Config{
		Host: "localhost", Port: "5432", User: "postgres",
		Password: "secret", Database: "mydb", SSLMode: "disable",
	}

	got := cfg.ConnectionString()
	want := "postgres://postgres:secret@localhost:5432/mydb?sslmode=disable"

	if got != want {
		t.Errorf("ConnectionString() got = %q, want %q", got, want)
	}
}
