package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/GitAlex9/go-order-service/internal/application/factory"
	"github.com/GitAlex9/go-order-service/internal/interfaces/http/routes"
	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
	"github.com/GitAlex9/go-order-service/internal/pkg/logger"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, services *factory.ServiceFactory, tokenManager *jwt.TokenManager, appLogger logger.Logger) *Server {
	router := routes.NewRouter(services, tokenManager, appLogger)

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("✓ Server listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
