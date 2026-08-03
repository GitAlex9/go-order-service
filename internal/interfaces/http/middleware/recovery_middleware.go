package middleware

import (
	"log"
	"net/http"

	"github.com/GitAlex9/go-order-service/internal/interfaces/http/response"
)

// Recovery captura panics em qualquer handler e responde 500,
// em vez de deixar o processo Go inteiro cair.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				response.JSONError(w, http.StatusInternalServerError, "internal server error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
