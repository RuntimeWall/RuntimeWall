package tracker

import (
	"fmt"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
)

// Conn wraps a terminal connection with runtime monitoring and policy enforcement.
type Conn struct {
	inner     sandbox.TerminalConn
	sandboxID string
	store     sandbox.EventStore
	lineBuf   []byte
	escState  escState
}

type escState int

const (
	escNone escState = iota
	escSeenESC
	escInCSI
)

// NewConn wraps conn for monitored terminal I/O.
func NewConn(inner sandbox.TerminalConn, sandboxID string, store sandbox.EventStore) *Conn {
	return &Conn{inner: inner, sandboxID: sandboxID, store: store}
}

func (c *Conn) ReadMessage() (int, []byte, error) {
	msgType, data, err := c.inner.ReadMessage()
	if err != nil {
		return msgType, data, err
	}

	if msgType != sandbox.TerminalMsgData || c.store == nil {
		return msgType, data, nil
	}

	var forward []byte
	for _, b := range data {
		// Discard ANSI / CSI escape sequences (arrow keys, bracketed paste,
		// terminal mode toggles, etc.) so they don't pollute logged commands.
		// Bytes are still forwarded to the container untouched.
		switch c.escState {
		case escSeenESC:
			if b == '[' {
				c.escState = escInCSI
			} else {
				c.escState = escNone
			}
			forward = append(forward, b)
			continue
		case escInCSI:
			if b >= 0x40 && b <= 0x7e {
				c.escState = escNone
			}
			forward = append(forward, b)
			continue
		}

		if b == 0x1b {
			c.escState = escSeenESC
			forward = append(forward, b)
			continue
		}

		if b == '\r' || b == '\n' {
			cmd := cleanCommand(string(c.lineBuf))
			c.lineBuf = c.lineBuf[:0]
			if cmd == "" {
				forward = append(forward, b)
				continue
			}

			ev := c.store.RecordCommand(c.sandboxID, cmd, sandbox.CommandSourceTerminal)
			if ev.Blocked {
				msg := fmt.Sprintf("\r\n\x1b[31m[RuntimeWall] BLOCKED\x1b[0m (%s): %s\r\n", ev.Threat, ev.Reason)
				_ = c.inner.WriteMessage(sandbox.TerminalMsgData, []byte(msg))
				continue
			}
			forward = append(forward, b)
			continue
		}

		if b == 127 || b == 8 {
			if len(c.lineBuf) > 0 {
				c.lineBuf = c.lineBuf[:len(c.lineBuf)-1]
			}
			forward = append(forward, b)
			continue
		}
		if b == 3 || b == 21 {
			c.lineBuf = c.lineBuf[:0]
			forward = append(forward, b)
			continue
		}
		if b >= 32 && b < 127 {
			c.lineBuf = append(c.lineBuf, b)
		}
		forward = append(forward, b)
	}

	if len(forward) == 0 {
		return msgType, nil, nil
	}
	return msgType, forward, nil
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	return c.inner.WriteMessage(messageType, data)
}

func (c *Conn) Close() error {
	return c.inner.Close()
}
