package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/handlers"
)

func RegisterAuthRoutes(r chi.Router, h *handlers.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
	})
}
