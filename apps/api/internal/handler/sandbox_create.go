package handler

import (
	"net/http"
)

// CreateUbuntu handles POST /sandbox/create.
func (h *Sandboxes) CreateUbuntu(w http.ResponseWriter, r *http.Request) {
	if h.manager == nil {
		writeError(w, http.StatusServiceUnavailable, "docker sandbox manager unavailable")
		return
	}

	result, err := h.manager.CreateUbuntu(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, result)
}
