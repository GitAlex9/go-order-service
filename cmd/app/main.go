package main

import (
	"fmt"
	"log"

	"github.com/GitAlex9/go-order-service/internal/application/id"
	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/infrastructure/database/postgres"
	repository "github.com/GitAlex9/go-order-service/internal/infrastructure/repositories/postgres"
)

func main() {

	// ==========================
	// Configuração
	// ==========================

	cfg := postgres.NewConfig()

	generator := id.NewCounterGenerator()

	productID := generator.Generate(id.ProductPrefix)
	customerID := generator.Generate(id.CustomerPrefix)
	orderID := generator.Generate(id.OrderPrefix)

	connection, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	log.Println("✓ Connected to PostgreSQL")

	// ==========================
	// Migration
	// ==========================

	migrator := postgres.NewMigrator(connection.Pool())

	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Database migrated")

	// ==========================
	// Repository
	// ==========================

	productRepository := repository.NewProductRepository(connection.Pool())
	customerRepository := repository.NewCustomerRepository(connection.Pool())
	orderRepository := repository.NewOrderRepository(connection.Pool())

	// ==========================
	// Criando Produto / Criando Customer / Criando Order
	// ==========================

	product, err := entities.NewProduct(
		productID,
		"Notebook Gamer",
		"RTX 5070",
		8500,
		10,
	)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	item, err := entities.NewOrderItem(
		product.ID,
		product.Name,
		product.Price,
		1,
	)
	if err != nil {
		log.Fatal(err)
	}

	customer, err := entities.NewCustomer(
		customerID,
		"John Doe",
		"john.doe@example.com",
	)

	if err != nil {
		log.Fatal(err)
	}

	order, err := entities.NewOrder(
		orderID,
		customer.ID,
		[]entities.OrderItem{*item},
	)

	if err != nil {
		log.Fatal(err)
	}

	// ==========================
	// Save
	// ==========================

	if err := productRepository.Save(product); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Product saved")

	if err := customerRepository.Save(customer); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Customer saved")

	if err := orderRepository.Save(order); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Order saved")

	// ==========================
	// Exists
	// ==========================

	exists, err := productRepository.Exists(product.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exists:", exists)

	exists, err = customerRepository.Exists(customer.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exists:", exists)

	exists, err = orderRepository.Exists(order.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exists:", exists)

	// ==========================
	// FindByID
	// ==========================

	found, err := productRepository.FindByID(product.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("========== PRODUCT ==========")
	fmt.Println("ID:", found.ID)
	fmt.Println("Name:", found.Name)
	fmt.Println("Description:", found.Description)
	fmt.Println("Price:", found.Price)
	fmt.Println("Stock:", found.Stock())
	fmt.Println("Active:", found.IsActive())
	fmt.Println()

	foundCustomer, err := customerRepository.FindByID(customer.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== CUSTOMER ==========")
	fmt.Println("ID:", foundCustomer.ID)
	fmt.Println("Name:", foundCustomer.Name)
	fmt.Println("Email:", foundCustomer.Email)
	fmt.Println("Active:", foundCustomer.IsActive())
	fmt.Println()

	foundOrder, err := orderRepository.FindByID(order.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== ORDER ==========")
	fmt.Println("ID:", foundOrder.ID)
	fmt.Println("Customer ID:", foundOrder.CustomerID)
	fmt.Println("Total:", foundOrder.Total())
	fmt.Println("Status:", foundOrder.Status())
	fmt.Println("Created At:", foundOrder.CreatedAt)
	fmt.Println("Updated At:", foundOrder.UpdatedAt)
	fmt.Println()

	// ==========================
	// List
	// ==========================

	products, err := productRepository.List()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== PRODUCTS ==========")

	for _, p := range products {
		fmt.Printf(
			"%s | %s | %.2f | Stock: %d\n",
			p.ID,
			p.Name,
			p.Price,
			p.Stock(),
		)
	}

	customers, err := customerRepository.List()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== CUSTOMERS ==========")

	for _, c := range customers {
		fmt.Printf(
			"%s | %s | %s | Active: %t\n",
			c.ID,
			c.Name,
			c.Email,
			c.IsActive(),
		)
	}

	orders, err := orderRepository.List()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== ORDERS ==========")

	for _, o := range orders {
		fmt.Printf(
			"%s | Customer ID: %s | Total: %.2f | Status: %s\n",
			o.ID,
			o.CustomerID,
			o.Total,
			o.Status(),
		)
	}

	// ==========================
	// Delete
	// ==========================

	if err := orderRepository.Delete(order.ID); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Order deleted")

	if err := customerRepository.Delete(customer.ID); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Customer deleted")

	if err := productRepository.Delete(product.ID); err != nil {
		log.Fatal(err)
	}

	log.Println("✓ Product deleted")

	// ==========================
	// Exists novamente
	// ==========================

	exists, err = productRepository.Exists(product.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exists after delete:", exists)

	exists, err = customerRepository.Exists(customer.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exists after delete:", exists)

	exists, err = orderRepository.Exists(order.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Exists after delete:", exists)
}
