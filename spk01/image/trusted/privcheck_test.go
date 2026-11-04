//go:build linux

package trusted

import (
	"os"
	"strings"
	"syscall"
	"testing"
)

func TestPrivilegeCheck(t *testing.T) {
	// 1. Non-root
	uid := os.Getuid()
	t.Logf("uid=%d gid=%d", uid, os.Getgid())
	if uid == 0 {
		t.Errorf("container runs as root")
	}

	// 2. Escalation attempt: try to become root
	if err := syscall.Setuid(0); err == nil {
		t.Errorf("setuid(0) succeeded: privilege escalation is possible")
	} else {
		t.Logf("setuid(0) blocked: %v", err)
	}

	// 3. No service-account token
	const tokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	if _, err := os.Stat(tokenPath); err == nil {
		t.Errorf("service-account token is mounted at %s", tokenPath)
	} else {
		t.Logf("token absent: %v", err)
	}

	// 4. Kernel flags: NoNewPrivs and effective capabilities
	status, err := os.ReadFile("/proc/self/status")
	if err != nil {
		t.Fatalf("read /proc/self/status: %v", err)
	}
	for _, line := range strings.Split(string(status), "\n") {
		if strings.HasPrefix(line, "NoNewPrivs:") || strings.HasPrefix(line, "CapEff:") {
			t.Log(line)
		}
		if strings.HasPrefix(line, "NoNewPrivs:") && strings.TrimSpace(line[len("NoNewPrivs:"):]) != "1" {
			t.Errorf("no_new_privs is not set: %s", line)
		}
	}
}
