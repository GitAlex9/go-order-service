package factory

import (
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/GitAlex9/go-order-service/internal/application/commands"
	"github.com/GitAlex9/go-order-service/internal/application/contracts"
	appevents "github.com/GitAlex9/go-order-service/internal/application/events"
	"github.com/GitAlex9/go-order-service/internal/application/queries"
	"github.com/GitAlex9/go-order-service/internal/application/services"
	repository "github.com/GitAlex9/go-order-service/internal/infrastructure/repositories/postgres"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

type ServiceFactory struct {
	CustomerService contracts.CustomerService
	ProductService  contracts.ProductService
	OrderService    contracts.OrderService
	UserService     contracts.UserService
	AuthService     contracts.AuthService
	TokenManager    *jwt.TokenManager
}

func NewServiceFactory(pool *pgxpool.Pool, log logger.Logger) *ServiceFactory {
	customerRepo := repository.NewCustomerRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	uow := repository.NewUnitOfWork(pool, log)
	dispatcher := appevents.NewDefaultDispatcher(log)
	tokenManager := jwt.NewTokenManager(getJWTSecret(), 24*time.Hour)

	customerService := services.NewCustomerService(
		commands.NewCreateCustomerHandler(customerRepo),
		commands.NewUpdateCustomerHandler(customerRepo, dispatcher),
		commands.NewDeleteCustomerHandler(customerRepo),
		queries.NewGetCustomerHandler(customerRepo),
		queries.NewListCustomersHandler(customerRepo),
	)

	productService := services.NewProductService(
		commands.NewCreateProductHandler(productRepo),
		commands.NewUpdateProductHandler(productRepo),
		commands.NewDeleteProductHandler(productRepo),
		commands.NewIncreaseStockHandler(productRepo, dispatcher),
		commands.NewDecreaseStockHandler(productRepo, dispatcher),
		commands.NewActivateProductHandler(productRepo),
		commands.NewDeactivateProductHandler(productRepo),
		queries.NewGetProductHandler(productRepo),
		queries.NewListProductsHandler(productRepo),
	)

	orderService := services.NewOrderService(
		commands.NewCreateOrderHandler(uow, dispatcher),
		commands.NewPayOrderHandler(uow, dispatcher),
		commands.NewCancelOrderHandler(uow, dispatcher),
		commands.NewDeleteOrderHandler(uow),
		queries.NewGetOrderHandler(orderRepo),
		queries.NewListOrdersHandler(orderRepo),
	)

	userService := services.NewUserService(
		commands.NewCreateUserHandler(userRepo),
		commands.NewChangePasswordHandler(userRepo, dispatcher),
		commands.NewChangeUserEmailHandler(userRepo, dispatcher),
		commands.NewActivateUserHandler(userRepo, dispatcher),
		commands.NewDeactivateUserHandler(userRepo, dispatcher),
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
		TokenManager:    tokenManager,
	}
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me" // ⚠️ só pra estudo — em produção, sempre via variável de ambiente
	}
	return secret
}
