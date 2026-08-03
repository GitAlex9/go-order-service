package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/handlers"
	appmiddleware "github.com/GitAlex9/go-order-service/internal/interfaces/http/middleware"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
)

func RegisterProductRoutes(r chi.Router, h *handlers.ProductHandler, tokenManager *jwt.TokenManager) {
	r.Route("/products", func(r chi.Router) {
		// leitura é pública
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)

		// escrita exige autenticação + role admin
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.Authenticate(tokenManager))
			r.Use(appmiddleware.RequireRole("admin"))

			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Patch("/{id}/increase-stock", h.IncreaseStock)
			r.Patch("/{id}/decrease-stock", h.DecreaseStock)
			r.Patch("/{id}/activate", h.Activate)
			r.Patch("/{id}/deactivate", h.Deactivate)
		})
	})
}
