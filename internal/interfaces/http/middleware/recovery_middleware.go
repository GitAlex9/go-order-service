package middleware

import (
	"net/http"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/response"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

func Recovery(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.FromContext(r.Context()).Error("panic recovered",
						"error", err,
						"path", r.URL.Path,
					)
					response.JSONError(w, http.StatusInternalServerError, "internal server error", nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
