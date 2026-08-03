package postgres

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		wantHost string
		wantPort string
		wantUser string
		wantPass string
		wantDB   string
		wantSSL  string
	}{
		{
			name:     "default values",
			envVars:  map[string]string{},
			wantHost: "localhost",
			wantPort: "5432",
			wantUser: "postgres",
			wantPass: "postgres",
			wantDB:   "go_order_service",
			wantSSL:  "disable",
		},
		{
			name: "environment variables override defaults",
			envVars: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "5433",
				"DB_USER":     "admin",
				"DB_PASSWORD": "secret",
				"DB_NAME":     "mydb",
				"DB_SSLMODE":  "require",
			},
			wantHost: "db.example.com",
			wantPort: "5433",
			wantUser: "admin",
			wantPass: "secret",
			wantDB:   "mydb",
			wantSSL:  "require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Limpa e seta as variáveis
			os.Clearenv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer os.Clearenv()

			cfg := NewConfig()

			if cfg.Host != tt.wantHost {
				t.Errorf("Host = %q, want %q", cfg.Host, tt.wantHost)
			}
			if cfg.Port != tt.wantPort {
				t.Errorf("Port = %q, want %q", cfg.Port, tt.wantPort)
			}
			if cfg.User != tt.wantUser {
				t.Errorf("User = %q, want %q", cfg.User, tt.wantUser)
			}
			if cfg.Password != tt.wantPass {
				t.Errorf("Password = %q, want %q", cfg.Password, tt.wantPass)
			}
			if cfg.Database != tt.wantDB {
				t.Errorf("Database = %q, want %q", cfg.Database, tt.wantDB)
			}
			if cfg.SSLMode != tt.wantSSL {
				t.Errorf("SSLMode = %q, want %q", cfg.SSLMode, tt.wantSSL)
			}
		})
	}
}

func TestConfig_ConnectionString(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want string
	}{
		{
			name: "standard config",
			cfg: &Config{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				Database: "testdb",
				SSLMode:  "disable",
			},
			want: "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "custom config",
			cfg: &Config{
				Host:     "db.example.com",
				Port:     "5433",
				User:     "admin",
				Password: "secret",
				Database: "mydb",
				SSLMode:  "require",
			},
			want: "postgres://admin:secret@db.example.com:5433/mydb?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.ConnectionString()
			if got != tt.want {
				t.Errorf("ConnectionString() = %q, want %q", got, tt.want)
			}
		})
	}
}
