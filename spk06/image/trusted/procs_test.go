//go:build linux

package trusted

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// maxSpawn bounds the fixture: it never starts more than this many processes,
// even when no limit is set.
const maxSpawn = 512

// TestProcessLimit starts sleeping child processes until the kernel refuses
// or maxSpawn is reached. It passes only if the pids limit stops it first.
func TestProcessLimit(t *testing.T) {
	b, err := os.ReadFile("/sys/fs/cgroup/pids.max")
	if err != nil {
		t.Fatalf("read pids.max: %v", err)
	}
	pidsMax := strings.TrimSpace(string(b))
	t.Logf("pids.max: %s", pidsMax)

	var children []*exec.Cmd
	defer func() {
		for _, c := range children {
			_ = c.Process.Kill()
			_ = c.Wait()
		}
	}()

	var spawnErr error
	for len(children) < maxSpawn {
		c := exec.Command("sleep", "30")
		if spawnErr = c.Start(); spawnErr != nil {
			break
		}
		children = append(children, c)
	}
	logFile(t, "/sys/fs/cgroup/pids.current")

	if spawnErr == nil {
		t.Errorf("started %d processes, no limit reached (pids.max %s)", len(children), pidsMax)
		return
	}
	t.Logf("stopped after %d processes: %v", len(children), spawnErr)
}
