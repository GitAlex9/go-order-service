package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/handlers"
	appmiddleware "github.com/GitAlex9/go-order-service/internal/interfaces/http/middleware"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
)

func RegisterUserRoutes(r chi.Router, h *handlers.UserHandler, tokenManager *jwt.TokenManager) {
	r.Route("/users", func(r chi.Router) {
		r.Use(appmiddleware.Authenticate(tokenManager))
		r.Use(appmiddleware.RequireRole("admin"))

		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}/password", h.ChangePassword)
		r.Patch("/{id}/email", h.ChangeEmail)
		r.Patch("/{id}/activate", h.Activate)
		r.Patch("/{id}/deactivate", h.Deactivate)
	})
}
