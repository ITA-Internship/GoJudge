//go:build linux

package trusted

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestFilesystemWrite(t *testing.T) {
	// 1. The root filesystem must be mounted read-only
	mounts, err := os.ReadFile("/proc/self/mounts")
	if err != nil {
		t.Fatalf("read /proc/self/mounts: %v", err)
	}
	rootRO := false
	for _, line := range strings.Split(string(mounts), "\n") {
		f := strings.Fields(line)
		if len(f) >= 4 && f[1] == "/" {
			t.Logf("root mount: %s", line)
			rootRO = f[3] == "ro" || strings.HasPrefix(f[3], "ro,")
		}
	}
	if !rootRO {
		t.Errorf("root filesystem is not mounted read-only")
	}

	// 2. Writes outside the allowed path must fail
	targets := []string{"/root/hack.txt", "/hack.txt", "/etc/hack.txt", "/workspace/hack.txt", "/usr/local/go/hack.txt"}
	for _, p := range targets {
		err := os.WriteFile(p, []byte("hack"), 0o600)
		switch {
		case err == nil:
			_ = os.Remove(p)
			t.Errorf("write to %s succeeded", p)
		case errors.Is(err, syscall.EROFS):
			t.Logf("write to %s blocked by read-only filesystem: %v", p, err)
		default:
			t.Logf("write to %s blocked by another error: %v", p, err)
		}
	}

	// 3. The allowed path must stay writable
	allowed := filepath.Join(os.TempDir(), "allowed.txt")
	if err := os.WriteFile(allowed, []byte("ok"), 0o600); err != nil {
		t.Errorf("write to allowed path %s failed: %v", allowed, err)
	} else {
		t.Logf("write to allowed path %s succeeded", allowed)
		_ = os.Remove(allowed)
	}
}
