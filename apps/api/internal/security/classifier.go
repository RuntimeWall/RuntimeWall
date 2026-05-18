package security

import (
	"regexp"
	"strings"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
)

var (
	reDestructive = regexp.MustCompile(`(?i)\brm\s+(-[a-zA-Z]*f[a-zA-Z]*\s+.*\/|.*--no-preserve-root|.*\s+\/\s*$)`)
	reExfilPipe   = regexp.MustCompile(`(?i)(curl|wget)\s+[^\s|]+\s*\|\s*(ba)?sh`)
	reExfilURL    = regexp.MustCompile(`(?i)(curl|wget)\s+https?://`)
	reReverseShell = regexp.MustCompile(`(?i)\b(nc|netcat)\s+.*(-e|--exec)\s+`)
	reNetcatListen = regexp.MustCompile(`(?i)\b(nc|netcat)\s+.*-l`)
	rePkgInstall  = regexp.MustCompile(`(?i)\b(apt|apt-get|yum|dnf|apk|pip|pip3|npm|yarn)\s+(.+\s+)?(install|add)\b`)
	reFileModify  = regexp.MustCompile(`(?i)\b(rm|mv|cp|chmod|chown|touch|tee|truncate)\b`)
	reNetTools    = regexp.MustCompile(`(?i)\b(nmap|masscan|nikto|sqlmap|hydra|nc|netcat|tcpdump|tshark)\b`)
)

// Classification is the result of analyzing a shell command.
type Classification struct {
	EventType sandbox.EventType
	Threat    sandbox.ThreatLevel
	Reason    string
}

// Classify inspects a command and returns event type and threat level.
func Classify(command string) Classification {
	cmd := strings.TrimSpace(command)
	if cmd == "" {
		return Classification{EventType: sandbox.EventCommand, Threat: sandbox.ThreatNone}
	}

	switch {
	case reDestructive.MatchString(cmd):
		return Classification{
			EventType: sandbox.EventCommand,
			Threat:    sandbox.ThreatDestructive,
			Reason:    "destructive filesystem operation",
		}
	case reExfilPipe.MatchString(cmd) || reReverseShell.MatchString(cmd):
		return Classification{
			EventType: sandbox.EventCommand,
			Threat:    sandbox.ThreatExfiltration,
			Reason:    "potential remote code execution or reverse shell",
		}
	case reExfilURL.MatchString(cmd) && strings.Contains(cmd, "|"):
		return Classification{
			EventType: sandbox.EventCommand,
			Threat:    sandbox.ThreatExfiltration,
			Reason:    "download piped to shell",
		}
	case reNetcatListen.MatchString(cmd):
		return Classification{
			EventType: sandbox.EventProcessLaunch,
			Threat:    sandbox.ThreatSuspicious,
			Reason:    "network listener detected",
		}
	case rePkgInstall.MatchString(cmd):
		threat := sandbox.ThreatNone
		reason := ""
		if strings.Contains(strings.ToLower(cmd), "nmap") || strings.Contains(strings.ToLower(cmd), "masscan") {
			threat = sandbox.ThreatSuspicious
			reason = "package install of security tool"
		}
		return Classification{
			EventType: sandbox.EventPackageInstall,
			Threat:    threat,
			Reason:    reason,
		}
	case reFileModify.MatchString(cmd):
		threat := sandbox.ThreatSuspicious
		reason := "file modification"
		if strings.Contains(strings.ToLower(cmd), "rm ") {
			threat = sandbox.ThreatDestructive
			reason = "file deletion"
		}
		return Classification{
			EventType: sandbox.EventFileModify,
			Threat:    threat,
			Reason:    reason,
		}
	case reNetTools.MatchString(cmd):
		return Classification{
			EventType: sandbox.EventProcessLaunch,
			Threat:    sandbox.ThreatSuspicious,
			Reason:    "network or security tool execution",
		}
	default:
		return Classification{
			EventType: sandbox.EventCommand,
			Threat:    sandbox.ThreatNone,
		}
	}
}
