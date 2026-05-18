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

const version = "0.3.0"

// sweeperRunner is implemented by long-running background sweepers (see
// sandbox/docker.Sweeper). It is wired in if the docker manager provides one.
type sweeperRunner interface {
	Start(ctx context.Context)
	Stop()
}

// Server is the RuntimeWall HTTP API.
type Server struct {
	cfg     config.Config
	http    *http.Server
	docker  sandbox.Manager
	sweeper sweeperRunner
}

// New wires routes and dependencies.
func New(cfg config.Config, sandboxes sandbox.Manager, events sandbox.EventStore) *Server {
	var attacher sandbox.TerminalAttacher
	if a, ok := sandboxes.(sandbox.TerminalAttacher); ok {
		attacher = a
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors)

	health := handler.NewHealth(sandboxes, version)
	sandboxHandler := handler.NewSandboxes(sandboxes)
	terminalHandler := handler.NewTerminal(attacher)
	terminalUI := handler.NewTerminalUI()
	commandsHandler := handler.NewCommands(events)
	eventsHandler := handler.NewEvents(events)
	policiesHandler := handler.NewPolicies(events)

	// Long-lived streams (no request timeout).
	r.Get("/terminal/{id}", terminalUI.Serve)
	r.Get("/api/v1/sandboxes/{id}/attach", terminalHandler.Attach)
	r.Get("/api/v1/sandboxes/{id}/commands/stream", commandsHandler.Stream)
	r.Get("/api/v1/sandboxes/{id}/events/stream", eventsHandler.Stream)
	r.Get("/api/v1/sandboxes/{id}/events/ws", eventsHandler.WebSocket)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(60 * time.Second))

		r.Get("/health", health.ServeHTTP)
		r.Post("/sandbox/create", sandboxHandler.CreateUbuntu)

		r.Route("/api/v1", func(r chi.Router) {
			r.Route("/sandboxes", func(r chi.Router) {
				r.Get("/", sandboxHandler.List)
				r.Post("/", sandboxHandler.Create)
				r.Post("/cleanup", sandboxHandler.Cleanup)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", sandboxHandler.Get)
					r.Delete("/", sandboxHandler.Delete)
					r.Post("/stop", sandboxHandler.Stop)
					r.Get("/commands", commandsHandler.List)
					r.Get("/events", eventsHandler.List)
					r.Get("/policy", policiesHandler.Get)
					r.Put("/policy", policiesHandler.Put)
				})
			})
		})
	})

	return &Server{
		cfg:    cfg,
		docker: sandboxes,
		http: &http.Server{
			Addr:              cfg.Addr,
			Handler:           r,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
		},
	}
}

// AttachSweeper registers a background sweeper that will be started with the
// server and stopped on shutdown. Calling with nil is a no-op.
func (s *Server) AttachSweeper(sw sweeperRunner) {
	if sw != nil {
		s.sweeper = sw
	}
}

// ListenAndServe starts the HTTP server (and background sweeper, if any).
func (s *Server) ListenAndServe() error {
	if s.sweeper != nil {
		s.sweeper.Start(context.Background())
	}
	slog.Info("runtimewall api listening", "addr", s.cfg.Addr, "version", version)
	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the server and closes Docker resources.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.sweeper != nil {
		s.sweeper.Stop()
	}
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
