package tracker

import (
	"io"
	"testing"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
)

// fakeConn feeds preset payloads to ReadMessage and records writes.
type fakeConn struct {
	reads    [][]byte
	readIdx  int
	writes   [][]byte
	closeErr error
}

func (f *fakeConn) ReadMessage() (int, []byte, error) {
	if f.readIdx >= len(f.reads) {
		return 0, nil, io.EOF
	}
	data := f.reads[f.readIdx]
	f.readIdx++
	return sandbox.TerminalMsgData, data, nil
}

func (f *fakeConn) WriteMessage(_ int, data []byte) error {
	f.writes = append(f.writes, append([]byte(nil), data...))
	return nil
}

func (f *fakeConn) Close() error { return f.closeErr }

// fakeStore captures recorded commands.
type fakeStore struct {
	recorded []sandbox.RuntimeEvent
}

func (s *fakeStore) RecordCommand(sandboxID, command string, src sandbox.CommandSource) sandbox.RuntimeEvent {
	ev := sandbox.RuntimeEvent{
		SandboxID: sandboxID,
		Command:   command,
		Source:    src,
		Event:     sandbox.EventCommand,
		Threat:    sandbox.ThreatNone,
	}
	s.recorded = append(s.recorded, ev)
	return ev
}

func (s *fakeStore) List(string) []sandbox.RuntimeEvent                   { return s.recorded }
func (s *fakeStore) Subscribe(string) (<-chan sandbox.RuntimeEvent, func()) { return nil, func() {} }
func (s *fakeStore) Clear(string)                                          {}
func (s *fakeStore) GetPolicy(string) sandbox.SecurityPolicy                { return sandbox.DefaultPolicy() }
func (s *fakeStore) SetPolicy(string, sandbox.SecurityPolicy)               {}

func TestConn_StripsBracketedPaste(t *testing.T) {
	store := &fakeStore{}
	inner := &fakeConn{
		reads: [][]byte{
			// ESC [ 2 0 0 ~  whoami  ESC [ 2 0 1 ~ \r
			append(append([]byte{0x1b, '[', '2', '0', '0', '~'}, []byte("whoami")...),
				append([]byte{0x1b, '[', '2', '0', '1', '~'}, '\r')...),
		},
	}
	c := NewConn(inner, "sb-1", store)

	_, _, _ = c.ReadMessage()

	if len(store.recorded) != 1 {
		t.Fatalf("recorded = %#v", store.recorded)
	}
	if got := store.recorded[0].Command; got != "whoami" {
		t.Fatalf("command = %q, want %q", got, "whoami")
	}
}

func TestConn_StripsArrowKeys(t *testing.T) {
	store := &fakeStore{}
	inner := &fakeConn{
		reads: [][]byte{
			append(append([]byte("ls"), 0x1b, '[', 'A'), '\r'),
		},
	}
	c := NewConn(inner, "sb-2", store)

	_, _, _ = c.ReadMessage()

	if len(store.recorded) != 1 || store.recorded[0].Command != "ls" {
		t.Fatalf("recorded = %#v", store.recorded)
	}
}
