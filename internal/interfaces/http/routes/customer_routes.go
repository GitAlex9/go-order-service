package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/handlers"
	appmiddleware "github.com/GitAlex9/go-order-service/internal/interfaces/http/middleware"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
)

func RegisterCustomerRoutes(r chi.Router, h *handlers.CustomerHandler, tokenManager *jwt.TokenManager) {
	r.Route("/customers", func(r chi.Router) {
		r.Use(appmiddleware.Authenticate(tokenManager))

		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
