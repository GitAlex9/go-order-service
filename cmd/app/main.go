package main

import (
	"context"
	"fmt"
	"log"

	"github.com/GitAlex9/go-order-service/internal/domain/entities"
	"github.com/GitAlex9/go-order-service/internal/domain/valueobjects"
	"github.com/GitAlex9/go-order-service/internal/infrastructure/database/postgres"
	repository "github.com/GitAlex9/go-order-service/internal/infrastructure/repositories/postgres"
)

func main() {
	ctx := context.Background()

	// ==========================
	// Configuração
	// ==========================

	cfg := postgres.NewConfig()

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
	// Criando Product
	// ==========================

	price, err := valueobjects.NewMoney(350000) // R$ 3.500,00
	if err != nil {
		log.Fatal(err)
	}

	product, err := entities.NewProduct("Notebook", "Notebook gamer", price, 10)
	if err != nil {
		log.Fatal(err)
	}

	// ==========================
	// Criando Customer
	// ==========================

	email, err := valueobjects.NewEmail("john.doe@example.com")
	if err != nil {
		log.Fatal(err)
	}

	cpf, err := valueobjects.NewCPF("111.444.777-35")
	if err != nil {
		log.Fatal(err)
	}

	customer, err := entities.NewCustomer("John Doe", email, cpf)
	if err != nil {
		log.Fatal(err)
	}

	// ==========================
	// Criando Order
	// ==========================

	item, err := entities.NewOrderItem(product.ID(), product.Name(), product.Price(), 2)
	if err != nil {
		log.Fatal(err)
	}

	order, err := entities.NewOrder(customer.ID(), []entities.OrderItem{*item})
	if err != nil {
		log.Fatal(err)
	}

	// ==========================
	// Save
	// ==========================

	if err := productRepository.Save(ctx, product); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Product saved")

	if err := customerRepository.Save(ctx, customer); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Customer saved")

	if err := orderRepository.Save(ctx, order); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Order saved")

	// ==========================
	// Exists
	// ==========================

	exists, err := productRepository.Exists(ctx, product.ID())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Product exists:", exists)

	exists, err = customerRepository.Exists(ctx, customer.ID())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Customer exists:", exists)

	exists, err = orderRepository.Exists(ctx, order.ID())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Order exists:", exists)

	// ==========================
	// FindByID
	// ==========================

	foundProduct, err := productRepository.FindByID(ctx, product.ID())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("========== PRODUCT ==========")
	fmt.Println("ID:", foundProduct.ID())
	fmt.Println("Name:", foundProduct.Name())
	fmt.Println("Description:", foundProduct.Description())
	fmt.Println("Price:", foundProduct.Price())
	fmt.Println("Stock:", foundProduct.Stock())
	fmt.Println("Active:", foundProduct.IsActive())
	fmt.Println()

	foundCustomer, err := customerRepository.FindByID(ctx, customer.ID())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== CUSTOMER ==========")
	fmt.Println("ID:", foundCustomer.ID())
	fmt.Println("Name:", foundCustomer.Name())
	fmt.Println("Email:", foundCustomer.Email())
	fmt.Println("CPF:", foundCustomer.CPF())
	fmt.Println()

	foundOrder, err := orderRepository.FindByID(ctx, order.ID())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== ORDER ==========")
	fmt.Println("ID:", foundOrder.ID())
	fmt.Println("Customer ID:", foundOrder.CustomerID())
	fmt.Println("Total:", foundOrder.Total())
	fmt.Println("Status:", foundOrder.Status())
	fmt.Println("Created At:", foundOrder.CreatedAt())
	fmt.Println("Updated At:", foundOrder.UpdatedAt())
	fmt.Println()

	// ==========================
	// List
	// ==========================

	products, err := productRepository.List(ctx, 0, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== PRODUCTS ==========")
	for _, p := range products {
		fmt.Printf("%s | %s | %s | Stock: %d\n", p.ID(), p.Name(), p.Price(), p.Stock())
	}

	customers, err := customerRepository.List(ctx, 0, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== CUSTOMERS ==========")
	for _, c := range customers {
		fmt.Printf("%s | %s | %s\n", c.ID(), c.Name(), c.Email())
	}

	orders, err := orderRepository.List(ctx, 0, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== ORDERS ==========")
	for _, o := range orders {
		fmt.Printf("%s | Customer ID: %s | Total: %s | Status: %s\n", o.ID(), o.CustomerID(), o.Total(), o.Status())
	}

	// ==========================
	// Delete
	// ==========================

	if err := orderRepository.Delete(ctx, order.ID()); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Order deleted")

	if err := customerRepository.Delete(ctx, customer.ID()); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Customer deleted")

	if err := productRepository.Delete(ctx, product.ID()); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Product deleted")

	// ==========================
	// Exists novamente
	// ==========================

	exists, err = productRepository.Exists(ctx, product.ID())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Product exists after delete:", exists)

	exists, err = customerRepository.Exists(ctx, customer.ID())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Customer exists after delete:", exists)

	exists, err = orderRepository.Exists(ctx, order.ID())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Order exists after delete:", exists)
}
