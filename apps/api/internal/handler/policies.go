package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/go-chi/chi/v5"
)

// Policies manages per-sandbox security policies.
type Policies struct {
	store sandbox.EventStore
}

// NewPolicies creates a policies handler.
func NewPolicies(store sandbox.EventStore) *Policies {
	return &Policies{store: store}
}

// Get handles GET /api/v1/sandboxes/{id}/policy.
func (h *Policies) Get(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "event store unavailable")
		return
	}
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"sandbox_id": id,
		"policy":     h.store.GetPolicy(id),
	})
}

// Put handles PUT /api/v1/sandboxes/{id}/policy.
func (h *Policies) Put(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "event store unavailable")
		return
	}

	id := chi.URLParam(r, "id")
	var policy sandbox.SecurityPolicy
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		writeError(w, http.StatusBadRequest, "invalid policy body")
		return
	}

	h.store.SetPolicy(id, policy)
	writeJSON(w, http.StatusOK, map[string]any{
		"sandbox_id": id,
		"policy":     policy,
	})
}
