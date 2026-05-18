package sandbox

import "time"

// EventType categorizes runtime activity inside a sandbox.
type EventType string

const (
	EventCommand         EventType = "command"
	EventPackageInstall  EventType = "package_install"
	EventFileModify      EventType = "file_modify"
	EventProcessLaunch   EventType = "process_launch"
	EventPolicyViolation EventType = "policy_violation"
)

// ThreatLevel describes security severity.
type ThreatLevel string

const (
	ThreatNone         ThreatLevel = "none"
	ThreatSuspicious   ThreatLevel = "suspicious"
	ThreatDestructive  ThreatLevel = "destructive"
	ThreatExfiltration ThreatLevel = "exfiltration"
)

// RuntimeEvent is a single observable sandbox action.
type RuntimeEvent struct {
	ID        string            `json:"id"`
	SandboxID string            `json:"sandbox_id"`
	Event     EventType         `json:"event"`
	Command   string            `json:"command,omitempty"`
	Threat    ThreatLevel       `json:"threat"`
	Blocked   bool              `json:"blocked"`
	Reason    string            `json:"reason,omitempty"`
	Source    CommandSource     `json:"source"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// EventStore records and streams runtime security events.
type EventStore interface {
	RecordCommand(sandboxID, command string, source CommandSource) RuntimeEvent
	List(sandboxID string) []RuntimeEvent
	Subscribe(sandboxID string) (<-chan RuntimeEvent, func())
	Clear(sandboxID string)
	GetPolicy(sandboxID string) SecurityPolicy
	SetPolicy(sandboxID string, policy SecurityPolicy)
}

// SecurityPolicy governs sandbox runtime behavior.
type SecurityPolicy struct {
	BlockNetworkTools        bool `json:"block_network_tools"`
	BlockPackageInstalls     bool `json:"block_package_installs"`
	ReadonlyFilesystem       bool `json:"readonly_filesystem"`
	BlockDestructiveCommands bool `json:"block_destructive_commands"`
	BlockExfiltration        bool `json:"block_exfiltration"`
}

// DefaultPolicy returns secure-by-default runtime policies.
func DefaultPolicy() SecurityPolicy {
	return SecurityPolicy{
		BlockNetworkTools:        true,
		BlockPackageInstalls:     false,
		ReadonlyFilesystem:       true,
		BlockDestructiveCommands: true,
		BlockExfiltration:        true,
	}
}
