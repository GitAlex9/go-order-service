package middleware

import (
	"net/http"
	"strings"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/response"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
)

func Authenticate(tokenManager *jwt.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.JSONError(w, http.StatusUnauthorized, "missing authorization header", nil)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.JSONError(w, http.StatusUnauthorized, "authorization header must be 'Bearer <token>'", nil)
				return
			}

			claims, err := tokenManager.Parse(parts[1])
			if err != nil {
				response.JSONError(w, http.StatusUnauthorized, "invalid or expired token", nil)
				return
			}

			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
