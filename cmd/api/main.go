package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/factory"
	domainerrors "github.com/GitAlex9/go-order-service/internal/domain/errors"
	"github.com/GitAlex9/go-order-service/internal/infrastructure/database/postgres"
)

// cmd/seed — popula dados iniciais (hoje: só o admin) direto no banco,
// sem passar pela API HTTP. Roda uma vez por ambiente:
//
//	go run ./cmd/seed/
//
// Como trata-se de teste, vale lembrar que nunca deve ser exposto via cmd/api
// Estava e estou rodando com APP_ENV=development go run ./cmd/api/ e (...)./cmd/app
// em app, se não me engano, na versão anterior eu havia chamado
// o método para dar resete no banco de dados e não tem seed para products.
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ no .env file found, using system environment variables")
	}

	ctx := context.Background()

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

	adminEmail := getEnv("SEED_ADMIN_EMAIL", "admin@example.com")
	adminPassword := getEnv("SEED_ADMIN_PASSWORD", "SenhaForte123!")

	admin, err := services.UserService.Create(ctx, dto.CreateUserRequest{
		Email:    adminEmail,
		Password: adminPassword,
		Role:     "admin",
	})
	if err != nil {
		// Se o admin já existir (rodou o seed mais de uma vez), não dará erro fatal
		// só avisa e segue, em vez de derrubar o processo.
		if errors.Is(err, domainerrors.ErrDuplicateEmail) {
			log.Printf("admin user %q already exists, skipping", adminEmail)
		} else {
			log.Fatal("seed failed:", err)
		}
	} else {
		log.Println("✓ Admin user created")
		log.Printf("  email: %s\n", admin.Email)
		log.Printf("  role:  %s\n", admin.Role)
		log.Println("  (use SEED_ADMIN_EMAIL / SEED_ADMIN_PASSWORD env vars to customize)")
	}

	seedProducts(ctx, services)
}

// popula um pequeno catálogo de exemplo,
// List/Get/Order sem precisar criar produto manualmente antes de cada teste.
// Idempotente por nome: se já existir um produto com aquele nome, pula.
func seedProducts(ctx context.Context, services *factory.ServiceFactory) {
	catalog := []dto.CreateProductRequest{
		{Name: "Mouse Gamer", Description: "Mouse RGB 16000 DPI", Price: 250.00, Stock: 15},
		{Name: "Teclado Mecânico", Description: "Teclado ABNT2 switch blue", Price: 450.00, Stock: 10},
		{Name: "Monitor 27\" 144Hz", Description: "Monitor gamer Full HD IPS", Price: 1800.00, Stock: 6},
		{Name: "Headset Gamer", Description: "Headset com microfone destacável", Price: 320.00, Stock: 20},
		{Name: "Notebook Gamer", Description: "Notebook RTX, 16GB RAM, 1TB SSD", Price: 7500.00, Stock: 3},
	}

	existing, err := services.ProductService.List(ctx, 0, 1000)
	if err != nil {
		log.Fatal("seed products: listing existing failed:", err)
	}

	existingNames := make(map[string]bool, len(existing))
	for _, p := range existing {
		existingNames[p.Name] = true
	}

	created := 0
	for _, item := range catalog {
		if existingNames[item.Name] {
			continue
		}
		if _, err := services.ProductService.Create(ctx, item); err != nil {
			log.Fatal("seed products: creating failed:", err)
		}
		created++
	}

	if created == 0 {
		log.Println("⚠ product catalog already seeded, skipping")
		return
	}

	log.Printf("✓ %d product(s) seeded\n", created)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
