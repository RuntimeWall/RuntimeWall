package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/config"
	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/handler"
	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const version = "0.1.0"

// Server is the RuntimeWall HTTP API.
type Server struct {
	cfg    config.Config
	http   *http.Server
	docker sandbox.Manager
}

// New wires routes and dependencies.
func New(cfg config.Config, sandboxes sandbox.Manager) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	health := handler.NewHealth(sandboxes, version)
	sandboxHandler := handler.NewSandboxes(sandboxes)

	r.Get("/health", health.ServeHTTP)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/sandboxes", func(r chi.Router) {
			r.Get("/", sandboxHandler.List)
			r.Post("/", sandboxHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", sandboxHandler.Get)
				r.Delete("/", sandboxHandler.Delete)
				r.Post("/stop", sandboxHandler.Stop)
			})
		})
	})

	return &Server{
		cfg:    cfg,
		docker: sandboxes,
		http: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	slog.Info("runtime wall api listening", "addr", s.cfg.Addr, "version", version)
	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the server and closes Docker resources.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.http.Shutdown(ctx); err != nil {
		return err
	}

	if closer, ok := s.docker.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			return fmt.Errorf("close docker manager: %w", err)
		}
	}
	return nil
}
