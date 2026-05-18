package security

import (
	"strings"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
)

// EvaluatePolicy decides if a classified command should be blocked.
func EvaluatePolicy(policy sandbox.SecurityPolicy, cls Classification, command string) (blocked bool, reason string) {
	cmd := strings.ToLower(strings.TrimSpace(command))

	if policy.BlockDestructiveCommands && cls.Threat == sandbox.ThreatDestructive {
		return true, "destructive commands are blocked by policy"
	}
	if policy.BlockExfiltration && cls.Threat == sandbox.ThreatExfiltration {
		return true, "exfiltration patterns are blocked by policy"
	}
	if policy.BlockNetworkTools && cls.Threat == sandbox.ThreatSuspicious &&
		(cls.EventType == sandbox.EventProcessLaunch || strings.Contains(cls.Reason, "network")) {
		return true, "network tools are blocked by policy"
	}
	if policy.BlockPackageInstalls && cls.EventType == sandbox.EventPackageInstall {
		return true, "package installs are blocked by policy"
	}
	if policy.ReadonlyFilesystem && cls.EventType == sandbox.EventFileModify {
		// Allow read-only ops; block mutating file commands.
		if strings.Contains(cmd, "rm ") || strings.Contains(cmd, "mv ") ||
			strings.Contains(cmd, "chmod ") || strings.Contains(cmd, "chown ") ||
			strings.Contains(cmd, "touch ") || strings.Contains(cmd, "tee ") {
			return true, "filesystem is read-only by policy"
		}
	}

	return false, ""
}
