package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/go-chi/chi/v5"
)

// Sweeper is implemented by managers that can reap stopped sandboxes.
type Sweeper interface {
	SweepStopped(ctx context.Context, ttl time.Duration) ([]string, error)
}

// Sandboxes exposes REST endpoints for Docker sandbox lifecycle.
type Sandboxes struct {
	manager sandbox.Manager
	sweeper Sweeper
}

// NewSandboxes creates a sandbox handler.
func NewSandboxes(manager sandbox.Manager) *Sandboxes {
	h := &Sandboxes{manager: manager}
	if s, ok := manager.(Sweeper); ok {
		h.sweeper = s
	}
	return h
}

// List handles GET /api/v1/sandboxes.
func (h *Sandboxes) List(w http.ResponseWriter, r *http.Request) {
	if h.manager == nil {
		writeError(w, http.StatusServiceUnavailable, "docker sandbox manager unavailable")
		return
	}

	items, err := h.manager.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if items == nil {
		items = []*sandbox.Sandbox{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"sandboxes": items})
}

// Create handles POST /api/v1/sandboxes.
func (h *Sandboxes) Create(w http.ResponseWriter, r *http.Request) {
	if h.manager == nil {
		writeError(w, http.StatusServiceUnavailable, "docker sandbox manager unavailable")
		return
	}

	var opts sandbox.CreateOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil && r.ContentLength > 0 {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sb, err := h.manager.Create(r.Context(), opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, sb)
}

// Get handles GET /api/v1/sandboxes/{id}.
func (h *Sandboxes) Get(w http.ResponseWriter, r *http.Request) {
	if h.manager == nil {
		writeError(w, http.StatusServiceUnavailable, "docker sandbox manager unavailable")
		return
	}

	id := chi.URLParam(r, "id")
	sb, err := h.manager.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, sandbox.ErrNotFound) {
			writeError(w, http.StatusNotFound, "sandbox not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, sb)
}

// Stop handles POST /api/v1/sandboxes/{id}/stop.
func (h *Sandboxes) Stop(w http.ResponseWriter, r *http.Request) {
	if h.manager == nil {
		writeError(w, http.StatusServiceUnavailable, "docker sandbox manager unavailable")
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.manager.Stop(r.Context(), id); err != nil {
		if errors.Is(err, sandbox.ErrNotFound) {
			writeError(w, http.StatusNotFound, "sandbox not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}

// Cleanup handles POST /api/v1/sandboxes/cleanup. It triggers an immediate
// sweep of stopped sandboxes older than ttl (default 0 = remove all stopped).
func (h *Sandboxes) Cleanup(w http.ResponseWriter, r *http.Request) {
	if h.manager == nil || h.sweeper == nil {
		writeError(w, http.StatusServiceUnavailable, "sandbox sweeper unavailable")
		return
	}

	ttl := time.Duration(0)
	if v := r.URL.Query().Get("ttl"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ttl = d
		}
	}
	// A zero ttl would no-op; treat empty query as "reap everything stopped".
	if ttl == 0 {
		ttl = time.Nanosecond
	}

	removed, err := h.sweeper.SweepStopped(r.Context(), ttl)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if removed == nil {
		removed = []string{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"removed": removed,
		"count":   len(removed),
	})
}

// Delete handles DELETE /api/v1/sandboxes/{id}.
func (h *Sandboxes) Delete(w http.ResponseWriter, r *http.Request) {
	if h.manager == nil {
		writeError(w, http.StatusServiceUnavailable, "docker sandbox manager unavailable")
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.manager.Remove(r.Context(), id); err != nil {
		if errors.Is(err, sandbox.ErrNotFound) {
			writeError(w, http.StatusNotFound, "sandbox not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
