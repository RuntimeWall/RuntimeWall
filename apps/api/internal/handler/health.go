package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
)

// Health reports API and Docker runtime status.
type Health struct {
	sandboxes sandbox.Manager
	version   string
}

// NewHealth creates a health handler.
func NewHealth(sandboxes sandbox.Manager, version string) *Health {
	return &Health{sandboxes: sandboxes, version: version}
}

type healthResponse struct {
	Status  string         `json:"status"`
	Version string         `json:"version"`
	Docker  dockerHealth   `json:"docker"`
}

type dockerHealth struct {
	Enabled   bool   `json:"enabled"`
	Connected bool   `json:"connected"`
	Message   string `json:"message,omitempty"`
}

// ServeHTTP handles GET /health.
func (h *Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{
		Status:  "ok",
		Version: h.version,
	}

	if h.sandboxes == nil {
		resp.Docker = dockerHealth{
			Enabled:   false,
			Connected: false,
			Message:   "docker sandbox manager not configured",
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.sandboxes.Ping(ctx); err != nil {
		resp.Docker = dockerHealth{
			Enabled:   true,
			Connected: false,
			Message:   err.Error(),
		}
		writeJSON(w, http.StatusOK, resp)
		return
	}

	resp.Docker = dockerHealth{
		Enabled:   true,
		Connected: true,
	}
	writeJSON(w, http.StatusOK, resp)
}
