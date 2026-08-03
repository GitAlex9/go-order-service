package main

import (
	"context"
	"log"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/factory"
	"github.com/GitAlex9/go-order-service/internal/infrastructure/database/postgres"
)

func main() {
	ctx := context.Background()

	cfg := postgres.NewConfig()
	connection, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	migrator := postgres.NewMigrator(connection.Pool())
	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}

	services := factory.NewServiceFactory(connection.Pool())

	admin, err := services.UserService.Create(ctx, dto.CreateUserRequest{
		Email:    "admin@example.com",
		Password: "SenhaForte123!",
		Role:     "admin",
	})
	if err != nil {
		log.Fatal("seed failed:", err)
	}

	log.Println("✓ Admin user created:", admin.Email)
}
