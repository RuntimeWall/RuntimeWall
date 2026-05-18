package handler

import (
	"embed"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed static/terminal.html
var terminalHTML embed.FS

// TerminalUI serves the browser terminal page.
type TerminalUI struct{}

// NewTerminalUI creates a terminal UI handler.
func NewTerminalUI() *TerminalUI {
	return &TerminalUI{}
}

// Serve handles GET /terminal/{id}.
func (h *TerminalUI) Serve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "sandbox id required", http.StatusBadRequest)
		return
	}

	f, err := terminalHTML.Open("static/terminal.html")
	if err != nil {
		http.Error(w, "terminal page not found", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.Copy(w, f)
}
