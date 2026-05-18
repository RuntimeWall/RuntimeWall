package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/go-chi/chi/v5"
)

// Sandboxes exposes REST endpoints for Docker sandbox lifecycle.
type Sandboxes struct {
	manager sandbox.Manager
}

// NewSandboxes creates a sandbox handler.
func NewSandboxes(manager sandbox.Manager) *Sandboxes {
	return &Sandboxes{manager: manager}
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
