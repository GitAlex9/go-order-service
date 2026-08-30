package middleware

import (
	"net/http"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/response"
)

// Lembrete da arquitetura, não esquecer: RequireRole só deve vir depois de Authenticate no chain de middlewares,
// já que depende dos claims já estarem no contexto.
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				response.JSONError(w, http.StatusUnauthorized, "authentication required", nil)
				return
			}

			if !allowed[claims.Role] {
				response.JSONError(w, http.StatusForbidden, "insufficient permissions", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
