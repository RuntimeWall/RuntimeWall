package sandbox

import "context"

// TerminalAttacher provides interactive shell access to a running sandbox.
type TerminalAttacher interface {
	AttachTerminal(ctx context.Context, sandboxID string, conn TerminalConn) error
}

// TerminalConn is the bidirectional terminal stream (e.g. WebSocket).
type TerminalConn interface {
	ReadMessage() (messageType int, data []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

// Control message types for terminal WebSocket protocol.
const (
	TerminalMsgData   = 0
	TerminalMsgResize = 1
)
