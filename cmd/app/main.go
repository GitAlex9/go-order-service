package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/GitAlex9/go-order-service/internal/application/dto"
	"github.com/GitAlex9/go-order-service/internal/application/factory"
	"github.com/GitAlex9/go-order-service/internal/infrastructure/database/postgres"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ no .env file found, using system environment variables")
	}

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
	// Reset + Migration
	// ==========================

	resetter := postgres.NewResetter(connection.Pool())
	if err := resetter.Reset(ctx); err != nil {
		log.Fatal("reset failed:", err)
	}
	log.Println("✓ Database reset")

	migrator := postgres.NewMigrator(connection.Pool())
	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}
	log.Println("✓ Database migrated")

	// ==========================
	// Factory
	// ==========================

	appLogger := logger.New("development")
	factories := factory.NewServiceFactory(connection.Pool(), appLogger)

	// ======================================================
	// CUSTOMER
	// ======================================================

	fmt.Println()
	fmt.Println("################ CUSTOMER ################")

	createdCustomer, err := factories.CustomerService.Create(ctx, dto.CreateCustomerRequest{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		CPF:   "111.444.777-35",
	})
	if err != nil {
		log.Fatal("create customer failed:", err)
	}
	log.Println("✓ Customer created")
	fmt.Printf("%+v\n", createdCustomer)

	foundCustomer, err := factories.CustomerService.Get(ctx, mustParseUUID(createdCustomer.ID))
	if err != nil {
		log.Fatal("get customer failed:", err)
	}
	fmt.Println("========== FOUND ==========")
	fmt.Printf("%+v\n", foundCustomer)

	updatedCustomer, err := factories.CustomerService.Update(ctx, mustParseUUID(createdCustomer.ID), dto.UpdateCustomerRequest{
		Name: "John Doe Junior",
	})
	if err != nil {
		log.Fatal("update customer failed:", err)
	}
	fmt.Println("========== UPDATED ==========")
	fmt.Printf("%+v\n", updatedCustomer)

	customerList, err := factories.CustomerService.List(ctx, 0, 10)
	if err != nil {
		log.Fatal("list customers failed:", err)
	}
	fmt.Println("========== LIST ==========")
	for _, c := range customerList {
		fmt.Printf("%s | %s | %s\n", c.ID, c.Name, c.Email)
	}

	fmt.Println("========== TESTE VALIDAÇÃO (esperado falhar) ==========")
	_, err = factories.CustomerService.Create(ctx, dto.CreateCustomerRequest{
		Name:  "ab",
		Email: "email-invalido",
		CPF:   "123",
	})
	printExpectedError(err)

	// ======================================================
	// PRODUCT
	// ======================================================

	fmt.Println()
	fmt.Println("################ PRODUCT ################")

	productA, err := factories.ProductService.Create(ctx, dto.CreateProductRequest{
		Name:        "Mouse Gamer",
		Description: "Mouse RGB 16000 DPI",
		Price:       250.00,
		Stock:       5,
	})
	if err != nil {
		log.Fatal("create productA failed:", err)
	}
	log.Println("✓ ProductA created")
	fmt.Printf("%+v\n", productA)

	productB, err := factories.ProductService.Create(ctx, dto.CreateProductRequest{
		Name:        "Teclado Mecânico",
		Description: "Teclado ABNT2 switch blue",
		Price:       450.00,
		Stock:       3,
	})
	if err != nil {
		log.Fatal("create productB failed:", err)
	}
	log.Println("✓ ProductB created")
	fmt.Printf("%+v\n", productB)

	updatedProduct, err := factories.ProductService.Update(ctx, mustParseUUID(productA.ID), dto.UpdateProductRequest{
		Price: 220.00,
	})
	if err != nil {
		log.Fatal("update product failed:", err)
	}
	fmt.Println("========== PRODUCT PRICE UPDATED ==========")
	fmt.Printf("%+v\n", updatedProduct)

	deactivated, err := factories.ProductService.Deactivate(ctx, mustParseUUID(productB.ID))
	if err != nil {
		log.Fatal("deactivate product failed:", err)
	}
	fmt.Println("========== PRODUCT DEACTIVATED ==========")
	fmt.Printf("%+v\n", deactivated)

	reactivated, err := factories.ProductService.Activate(ctx, mustParseUUID(productB.ID))
	if err != nil {
		log.Fatal("activate product failed:", err)
	}
	fmt.Println("========== PRODUCT REACTIVATED ==========")
	fmt.Printf("%+v\n", reactivated)

	fmt.Println("========== TESTE ESTOQUE INSUFICIENTE (esperado falhar) ==========")
	_, err = factories.ProductService.DecreaseStock(ctx, mustParseUUID(productA.ID), dto.AdjustStockRequest{Quantity: 999})
	printExpectedError(err)

	// ======================================================
	// ORDER
	// ======================================================

	fmt.Println()
	fmt.Println("################ ORDER ################")

	fmt.Println("========== ESTOQUE ANTES DO PEDIDO ==========")
	printStock(ctx, factories, productA.ID, productB.ID)

	createdOrder, err := factories.OrderService.Create(ctx, dto.CreateOrderRequest{
		CustomerID: createdCustomer.ID,
		Items: []dto.OrderItemRequest{
			{ProductID: productA.ID, Quantity: 2},
			{ProductID: productB.ID, Quantity: 1},
		},
	})
	if err != nil {
		log.Fatal("create order failed:", err)
	}
	log.Println("✓ Order created")
	fmt.Printf("%+v\n", createdOrder)

	fmt.Println("========== ESTOQUE APÓS PEDIDO ==========")
	printStock(ctx, factories, productA.ID, productB.ID)

	orderID := mustParseUUID(createdOrder.ID)

	paidOrder, err := factories.OrderService.Pay(ctx, orderID)
	if err != nil {
		log.Fatal("pay order failed:", err)
	}
	fmt.Println("========== ORDER PAID ==========")
	fmt.Printf("Status: %s | Total: %.2f\n", paidOrder.Status, paidOrder.Total)

	fmt.Println("========== TESTE: CANCELAR PEDIDO JÁ PAGO (esperado falhar) ==========")
	_, err = factories.OrderService.Cancel(ctx, orderID)
	printExpectedError(err)

	fmt.Println("========== TESTE: DELETAR PEDIDO PAGO (esperado falhar) ==========")
	err = factories.OrderService.Delete(ctx, orderID)
	printExpectedError(err)

	secondOrder, err := factories.OrderService.Create(ctx, dto.CreateOrderRequest{
		CustomerID: createdCustomer.ID,
		Items: []dto.OrderItemRequest{
			{ProductID: productA.ID, Quantity: 1},
		},
	})
	if err != nil {
		log.Fatal("create second order failed:", err)
	}
	log.Println("✓ Second order created")

	fmt.Println("========== ESTOQUE ANTES DO CANCEL ==========")
	printStock(ctx, factories, productA.ID, productB.ID)

	canceledOrder, err := factories.OrderService.Cancel(ctx, mustParseUUID(secondOrder.ID))
	if err != nil {
		log.Fatal("cancel order failed:", err)
	}
	fmt.Println("========== ORDER CANCELED ==========")
	fmt.Printf("Status: %s\n", canceledOrder.Status)

	fmt.Println("========== ESTOQUE APÓS CANCEL (estoque devolvido) ==========")
	printStock(ctx, factories, productA.ID, productB.ID)

	if err := factories.OrderService.Delete(ctx, mustParseUUID(secondOrder.ID)); err != nil {
		log.Fatal("delete canceled order failed:", err)
	}
	log.Println("✓ Canceled order deleted")

	allOrders, err := factories.OrderService.List(ctx, 0, 10)
	if err != nil {
		log.Fatal("list orders failed:", err)
	}
	fmt.Println("========== ALL ORDERS ==========")
	for _, o := range allOrders {
		fmt.Printf("%s | Customer: %s | Status: %s | Total: %.2f\n", o.ID, o.CustomerID, o.Status, o.Total)
	}

	// ======================================================
	// USER + AUTH
	// ======================================================

	fmt.Println()
	fmt.Println("################ USER + AUTH ################")

	createdUser, err := factories.UserService.Create(ctx, dto.CreateUserRequest{
		Email:    "admin@example.com",
		Password: "SenhaForte123!",
		Role:     "admin",
	})
	if err != nil {
		log.Fatal("create user failed:", err)
	}
	log.Println("✓ User created")
	fmt.Printf("%+v\n", createdUser)

	userID := mustParseUUID(createdUser.ID)

	fmt.Println("========== TESTE LOGIN CORRETO ==========")
	loginResp, err := factories.AuthService.Login(ctx, dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "SenhaForte123!",
	})
	if err != nil {
		log.Fatal("login failed:", err)
	}
	fmt.Println("Token gerado (primeiros 40 chars):", loginResp.Token[:40], "...")
	fmt.Printf("Logged in as: %+v\n", loginResp.User)

	fmt.Println("========== TESTE LOGIN SENHA ERRADA (esperado falhar) ==========")
	_, err = factories.AuthService.Login(ctx, dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "senhaerrada",
	})
	printExpectedError(err)

	fmt.Println("========== TESTE LOGIN EMAIL INEXISTENTE (esperado falhar, mesma mensagem genérica) ==========")
	_, err = factories.AuthService.Login(ctx, dto.LoginRequest{
		Email:    "naoexiste@example.com",
		Password: "qualquercoisa",
	})
	printExpectedError(err)

	if err := factories.UserService.ChangePassword(ctx, userID, dto.ChangePasswordRequest{
		CurrentPassword: "SenhaForte123!",
		NewPassword:     "NovaSenhaForte456!",
	}); err != nil {
		log.Fatal("change password failed:", err)
	}
	log.Println("✓ Password changed")

	fmt.Println("========== TESTE LOGIN COM SENHA ANTIGA (esperado falhar) ==========")
	_, err = factories.AuthService.Login(ctx, dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "SenhaForte123!",
	})
	printExpectedError(err)

	fmt.Println("========== TESTE LOGIN COM SENHA NOVA (esperado funcionar) ==========")
	_, err = factories.AuthService.Login(ctx, dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "NovaSenhaForte456!",
	})
	if err != nil {
		log.Fatal("login with new password should have worked:", err)
	}
	log.Println("✓ Login with new password succeeded")

	deactivatedUser, err := factories.UserService.Deactivate(ctx, userID)
	if err != nil {
		log.Fatal("deactivate user failed:", err)
	}
	fmt.Println("========== USER DEACTIVATED ==========")
	fmt.Printf("%+v\n", deactivatedUser)

	fmt.Println("========== TESTE LOGIN COM USUÁRIO INATIVO (esperado falhar) ==========")
	_, err = factories.AuthService.Login(ctx, dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "NovaSenhaForte456!",
	})
	printExpectedError(err)

	reactivatedUser, err := factories.UserService.Activate(ctx, userID)
	if err != nil {
		log.Fatal("reactivate user failed:", err)
	}
	fmt.Println("========== USER REACTIVATED ==========")
	fmt.Printf("%+v\n", reactivatedUser)

	userList, err := factories.UserService.List(ctx, 0, 10)
	if err != nil {
		log.Fatal("list users failed:", err)
	}
	fmt.Println("========== USER LIST ==========")
	for _, u := range userList {
		fmt.Printf("%s | %s | %s | Active: %t\n", u.ID, u.Email, u.Role, u.Active)
	}

	// ======================================================
	// CLEANUP (respeitando ordem de FK)
	// ======================================================

	fmt.Println()
	fmt.Println("################ CLEANUP ################")

	if err := factories.OrderService.Delete(ctx, orderID); err != nil {
		fmt.Println("⚠ delete paid order (expected to fail by rule):", err)
	}

	if err := factories.ProductService.Delete(ctx, mustParseUUID(productA.ID)); err != nil {
		fmt.Println("⚠ delete productA failed (likely still referenced by paid order):", err)
	} else {
		log.Println("✓ ProductA deleted")
	}

	if err := factories.ProductService.Delete(ctx, mustParseUUID(productB.ID)); err != nil {
		fmt.Println("⚠ delete productB failed:", err)
	} else {
		log.Println("✓ ProductB deleted")
	}

	if err := factories.CustomerService.Delete(ctx, mustParseUUID(createdCustomer.ID)); err != nil {
		fmt.Println("⚠ delete customer failed (likely still referenced by paid order):", err)
	} else {
		log.Println("✓ Customer deleted")
	}

	log.Println("✓ Test run finished")
}

// ==========================
// Helpers
// ==========================

func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		log.Fatal("invalid uuid:", err)
	}
	return id
}

func printExpectedError(err error) {
	if err != nil {
		fmt.Println("Erro esperado recebido:", err)
	} else {
		fmt.Println("⚠ Deveria ter falhado, mas não falhou!")
	}
}

func printStock(ctx context.Context, factories *factory.ServiceFactory, productAID, productBID string) {
	pa, err := factories.ProductService.Get(ctx, mustParseUUID(productAID))
	if err != nil {
		log.Fatal("get productA failed:", err)
	}
	pb, err := factories.ProductService.Get(ctx, mustParseUUID(productBID))
	if err != nil {
		log.Fatal("get productB failed:", err)
	}
	fmt.Printf("Mouse: %d | Teclado: %d\n", pa.Stock, pb.Stock)
}
