package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/go-chi/chi/v5"
)

// Commands exposes command history (subset of runtime events).
type Commands struct {
	store sandbox.EventStore
}

// NewCommands creates a commands handler.
func NewCommands(store sandbox.EventStore) *Commands {
	return &Commands{store: store}
}

// List handles GET /api/v1/sandboxes/{id}/commands.
func (h *Commands) List(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "event store unavailable")
		return
	}

	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"sandbox_id": id,
		"commands":   commandEvents(h.store.List(id)),
	})
}

// Stream handles GET /api/v1/sandboxes/{id}/commands/stream (SSE).
func (h *Commands) Stream(w http.ResponseWriter, r *http.Request) {
	if h.store == nil {
		writeError(w, http.StatusServiceUnavailable, "event store unavailable")
		return
	}

	id := chi.URLParam(r, "id")
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for _, ev := range commandEvents(h.store.List(id)) {
		if err := writeCommandSSE(w, ev); err != nil {
			return
		}
	}
	flusher.Flush()

	events, unsub := h.store.Subscribe(id)
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
			if ev.Event != sandbox.EventCommand && ev.Event != sandbox.EventPolicyViolation {
				continue
			}
			if err := writeCommandSSE(w, eventToCommand(ev)); err != nil {
				return
			}
			flusher.Flush()
		case <-tick.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}

func commandEvents(events []sandbox.RuntimeEvent) []sandbox.CommandRecord {
	out := make([]sandbox.CommandRecord, 0)
	for _, ev := range events {
		if ev.Event == sandbox.EventCommand || ev.Event == sandbox.EventPolicyViolation ||
			ev.Event == sandbox.EventPackageInstall || ev.Event == sandbox.EventFileModify ||
			ev.Event == sandbox.EventProcessLaunch {
			out = append(out, eventToCommand(ev))
		}
	}
	return out
}

func eventToCommand(ev sandbox.RuntimeEvent) sandbox.CommandRecord {
	return sandbox.CommandRecord{
		ID:        ev.ID,
		SandboxID: ev.SandboxID,
		Command:   ev.Command,
		Source:    ev.Source,
		Timestamp: ev.Timestamp,
	}
}

func writeCommandSSE(w http.ResponseWriter, rec sandbox.CommandRecord) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "event: command\ndata: %s\n\n", data)
	return err
}
