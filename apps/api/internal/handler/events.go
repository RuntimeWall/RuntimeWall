package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

// Events exposes runtime security event APIs.
type Events struct {
	store sandbox.EventStore
}

// NewEvents creates an events handler.
func NewEvents(store sandbox.EventStore) *Events {
	return &Events{store: store}
}

// List handles GET /api/v1/sandboxes/{id}/events.
func (h *Events) List(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "event store unavailable")
		return
	}
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"sandbox_id": id,
		"events":     h.store.List(id),
	})
}

// Stream handles GET /api/v1/sandboxes/{id}/events/stream (SSE).
func (h *Events) Stream(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "event store unavailable")
		return
	}
	streamEvents(w, r, h.store, chi.URLParam(r, "id"))
}

var eventsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket handles GET /api/v1/sandboxes/{id}/events (WebSocket).
func (h *Events) WebSocket(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "event store unavailable")
		return
	}

	id := chi.URLParam(r, "id")
	ws, err := eventsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	for _, ev := range h.store.List(id) {
		if err := ws.WriteJSON(ev); err != nil {
			return
		}
	}

	events, unsub := h.store.Subscribe(id)
	defer unsub()

	for {
		select {
		case <-r.Context().Done():
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			if err := ws.WriteJSON(ev); err != nil {
				return
			}
		}
	}
}

func streamEvents(w http.ResponseWriter, r *http.Request, store sandbox.EventStore, id string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for _, ev := range store.List(id) {
		if err := writeEventSSE(w, ev); err != nil {
			return
		}
	}
	flusher.Flush()

	events, unsub := store.Subscribe(id)
	defer unsub()

	tick := time.NewTicker(15 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			if err := writeEventSSE(w, ev); err != nil {
				return
			}
			flusher.Flush()
		case <-tick.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

func writeEventSSE(w http.ResponseWriter, ev sandbox.RuntimeEvent) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "event: runtime\ndata: %s\n\n", data)
	return err
}
