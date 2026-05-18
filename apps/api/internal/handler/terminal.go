package handler

import (
	"net/http"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var terminalUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// wsTerminalConn adapts gorilla WebSocket to sandbox.TerminalConn.
type wsTerminalConn struct {
	*websocket.Conn
}

func (w *wsTerminalConn) ReadMessage() (int, []byte, error) {
	msgType, data, err := w.Conn.ReadMessage()
	if err != nil {
		return 0, nil, err
	}
	switch msgType {
	case websocket.BinaryMessage:
		return sandbox.TerminalMsgData, data, nil
	case websocket.TextMessage:
		// JSON resize or raw input from browser
		if len(data) > 0 && data[0] == '{' {
			return sandbox.TerminalMsgResize, data, nil
		}
		return sandbox.TerminalMsgData, data, nil
	default:
		return sandbox.TerminalMsgData, data, nil
	}
}

func (w *wsTerminalConn) WriteMessage(messageType int, data []byte) error {
	return w.Conn.WriteMessage(websocket.BinaryMessage, data)
}

// Terminal handles browser and CLI terminal attach over WebSocket.
type Terminal struct {
	attacher sandbox.TerminalAttacher
}

// NewTerminal creates a terminal handler.
func NewTerminal(attacher sandbox.TerminalAttacher) *Terminal {
	return &Terminal{attacher: attacher}
}

// Attach handles GET /api/v1/sandboxes/{id}/attach (WebSocket upgrade).
func (h *Terminal) Attach(w http.ResponseWriter, r *http.Request) {
	if h.attacher == nil {
		writeError(w, http.StatusServiceUnavailable, "docker sandbox manager unavailable")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "sandbox id required")
		return
	}

	ws, err := terminalUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	conn := &wsTerminalConn{Conn: ws}
	if err := h.attacher.AttachTerminal(r.Context(), id, conn); err != nil {
		_ = ws.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
	}
}
