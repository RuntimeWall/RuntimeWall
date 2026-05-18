package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/config"
	sandboxdocker "github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox/docker"
	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox/tracker"
	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()
	eventStore := tracker.NewStore()

	var sandboxes *sandboxdocker.Manager
	if cfg.EnableDocker {
		mgr, err := sandboxdocker.NewManager(cfg, eventStore)
		if err != nil {
			slog.Warn("docker sandbox manager unavailable; sandbox routes disabled", "error", err)
		} else {
			sandboxes = mgr
			slog.Info("docker sandbox manager ready")
		}
	} else {
		slog.Info("docker sandbox manager disabled via ENABLE_DOCKER=false")
	}

	srv := server.New(cfg, sandboxes, eventStore)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
