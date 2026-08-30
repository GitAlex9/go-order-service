package postgres

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid connection (assuming DB running locally)",
			cfg: &Config{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				Database: "go_order_service",
				SSLMode:  "disable",
			},
			wantErr: false,
		},
		{
			name: "invalid connection (wrong host)",
			cfg: &Config{
				Host:     "nonexistent.host",
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				Database: "go_order_service",
				SSLMode:  "disable",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := NewConnection(tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if conn != nil {
					t.Errorf("expected nil connection, got %v", conn)
				}
				return
			}

			if err != nil {
				t.Skipf("Skipping test: database not available: %v", err)
			}
			defer conn.Close()

			if conn.Pool() == nil {
				t.Errorf("Pool() returned nil")
			}
		})
	}
}
