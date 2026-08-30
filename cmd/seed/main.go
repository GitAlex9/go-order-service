package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/factory"
	"github.com/GitAlex9/go-order-service/internal/infrastructure/database/postgres"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

// cmd/seed — popula dados iniciais (admin + catálogo de produtos) direto no banco,
// sem passar pela API HTTP. Roda uma vez por ambiente:
//
//	go run ./cmd/seed/
//
// É seguro rodar mais de uma vez: tanto o admin quanto os produtos são
// verificados antes de criar, então nada é duplicado.
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

	appLogger := logger.New(getEnv("APP_ENV", "development"))
	services := factory.NewServiceFactory(connection.Pool(), appLogger)

	seedAdmin(ctx, services)
	seedProducts(ctx, services)
}

// seedAdmin cria o usuário administrador padrão, caso ainda não exista.
// A checagem é feita listando usuários existentes (em vez de depender de
// como um eventual erro de duplicidade é propagado internamente) — mais
// robusto e consistente com a estratégia usada em seedProducts.
func seedAdmin(ctx context.Context, services *factory.ServiceFactory) {
	adminEmail := getEnv("SEED_ADMIN_EMAIL", "admin@example.com")
	adminPassword := getEnv("SEED_ADMIN_PASSWORD", "SenhaForte123!")

	existing, err := services.UserService.List(ctx, 0, 1000)
	if err != nil {
		log.Fatal("seed admin: listing existing users failed:", err)
	}

	for _, u := range existing {
		if u.Email == adminEmail {
			log.Printf("⚠ admin user %q already exists, skipping", adminEmail)
			return
		}
	}

	admin, err := services.UserService.Create(ctx, dto.CreateUserRequest{
		Email:    adminEmail,
		Password: adminPassword,
		Role:     "admin",
	})
	if err != nil {
		log.Fatal("seed admin failed:", err)
	}

	log.Println("✓ Admin user created")
	log.Printf("  email: %s\n", admin.Email)
	log.Printf("  role:  %s\n", admin.Role)
	log.Println("  (use SEED_ADMIN_EMAIL / SEED_ADMIN_PASSWORD env vars to customize)")
}

// seedProducts popula um pequeno catálogo de exemplo, útil para testar
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
