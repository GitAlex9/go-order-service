package postgres

import (
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func NewConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Database: getEnv("DB_NAME", "go_order_service"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (c Config) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
