package routes

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/GitAlex9/go-order-service/internal/application/factory"
	"github.com/GitAlex9/go-order-service/internal/interfaces/http/handlers"
	appmiddleware "github.com/GitAlex9/go-order-service/internal/interfaces/http/middleware"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

func NewRouter(services *factory.ServiceFactory, tokenManager *jwt.TokenManager, log logger.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(appmiddleware.Logging(log))
	r.Use(appmiddleware.Recovery(log))

	authHandler := handlers.NewAuthHandler(services.AuthService)
	customerHandler := handlers.NewCustomerHandler(services.CustomerService)
	productHandler := handlers.NewProductHandler(services.ProductService)
	orderHandler := handlers.NewOrderHandler(services.OrderService)
	userHandler := handlers.NewUserHandler(services.UserService)

	r.Route("/api/v1", func(r chi.Router) {
		RegisterAuthRoutes(r, authHandler)
		RegisterCustomerRoutes(r, customerHandler, tokenManager)
		RegisterProductRoutes(r, productHandler, tokenManager)
		RegisterOrderRoutes(r, orderHandler, tokenManager)
		RegisterUserRoutes(r, userHandler, tokenManager)
	})

	return r
}
