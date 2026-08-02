package factory

import (
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/GitAlex9/go-order-service/internal/application/commands"
	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	"github.com/GitAlex9/go-order-service/internal/application/queries"
	"github.com/GitAlex9/go-order-service/internal/application/services"
	repository "github.com/GitAlex9/go-order-service/internal/infrastructure/repositories/postgres"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
)

type ServiceFactory struct {
	CustomerService contracts.CustomerService
	ProductService  contracts.ProductService
	OrderService    contracts.OrderService
	UserService     contracts.UserService
	AuthService     contracts.AuthService
}

func NewServiceFactory(pool *pgxpool.Pool) *ServiceFactory {
	customerRepo := repository.NewCustomerRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	tokenManager := jwt.NewTokenManager(getJWTSecret(), 24*time.Hour)

	customerService := services.NewCustomerService(
		commands.NewCreateCustomerHandler(customerRepo),
		commands.NewUpdateCustomerHandler(customerRepo),
		commands.NewDeleteCustomerHandler(customerRepo),
		queries.NewGetCustomerHandler(customerRepo),
		queries.NewListCustomersHandler(customerRepo),
	)

	productService := services.NewProductService(
		commands.NewCreateProductHandler(productRepo),
		commands.NewUpdateProductHandler(productRepo),
		commands.NewDeleteProductHandler(productRepo),
		commands.NewIncreaseStockHandler(productRepo),
		commands.NewDecreaseStockHandler(productRepo),
		commands.NewActivateProductHandler(productRepo),
		commands.NewDeactivateProductHandler(productRepo),
		queries.NewGetProductHandler(productRepo),
		queries.NewListProductsHandler(productRepo),
	)

	orderService := services.NewOrderService(
		commands.NewCreateOrderHandler(orderRepo, productRepo, customerRepo),
		commands.NewPayOrderHandler(orderRepo),
		commands.NewCancelOrderHandler(orderRepo, productRepo),
		commands.NewDeleteOrderHandler(orderRepo),
		queries.NewGetOrderHandler(orderRepo),
		queries.NewListOrdersHandler(orderRepo),
	)

	userService := services.NewUserService(
		commands.NewCreateUserHandler(userRepo),
		commands.NewChangePasswordHandler(userRepo),
		commands.NewChangeUserEmailHandler(userRepo),
		commands.NewActivateUserHandler(userRepo),
		commands.NewDeactivateUserHandler(userRepo),
		queries.NewGetUserHandler(userRepo),
		queries.NewListUsersHandler(userRepo),
	)

	authService := services.NewAuthService(
		commands.NewAuthenticateHandler(userRepo, tokenManager),
	)

	return &ServiceFactory{
		CustomerService: customerService,
		ProductService:  productService,
		OrderService:    orderService,
		UserService:     userService,
		AuthService:     authService,
	}
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me" // estudo — em produção, sempre via variável de ambiente
	}
	return secret
}
