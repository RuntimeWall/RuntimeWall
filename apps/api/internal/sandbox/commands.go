package sandbox

import "time"

// CommandSource identifies how a command entered the sandbox.
type CommandSource string

const (
	CommandSourceTerminal CommandSource = "terminal"
)

// CommandRecord is a single executed command inside a sandbox.
type CommandRecord struct {
	ID        string        `json:"id"`
	SandboxID string        `json:"sandbox_id"`
	Command   string        `json:"command"`
	Source    CommandSource `json:"source"`
	Timestamp time.Time     `json:"timestamp"`
}

// CommandTracker stores and streams sandbox command events.
type CommandTracker interface {
	Record(sandboxID string, command string, source CommandSource)
	List(sandboxID string) []CommandRecord
	Subscribe(sandboxID string) (<-chan CommandRecord, func())
	Clear(sandboxID string)
}
