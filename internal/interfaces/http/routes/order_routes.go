package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/handlers"
	appmiddleware "github.com/GitAlex9/go-order-service/internal/interfaces/http/middleware"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
)

func RegisterOrderRoutes(r chi.Router, h *handlers.OrderHandler, tokenManager *jwt.TokenManager) {
	r.Route("/orders", func(r chi.Router) {
		r.Use(appmiddleware.Authenticate(tokenManager))

		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}/pay", h.Pay)
		r.Patch("/{id}/cancel", h.Cancel)
		r.Delete("/{id}", h.Delete)
	})
}
