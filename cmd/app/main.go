package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/GitAlex9/go-order-service/internal/application/factory"
	"github.com/GitAlex9/go-order-service/internal/infrastructure/database/postgres"
	httpserver "github.com/GitAlex9/go-order-service/internal/interfaces/http"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ no .env file found, using system environment variables")
	}

	cfg := postgres.NewConfig()

	connection, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	log.Println("✓ Connected to PostgreSQL")

	migrator := postgres.NewMigrator(connection.Pool())
	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Database migrated")

	services := factory.NewServiceFactory(connection.Pool())
	tokenManager := jwt.NewTokenManager(getJWTSecret(), 24*time.Hour)

	server := httpserver.NewServer(":8080", services, tokenManager)

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed:", err)
		}
	}()

	// Graceful shutdown: espera SIGINT/SIGTERM (Ctrl+C, ou o Docker parando o container)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	log.Println("✓ Server stopped gracefully")
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	return secret
}
