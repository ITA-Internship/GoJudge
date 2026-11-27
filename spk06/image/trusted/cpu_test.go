//go:build linux

package trusted

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestCPUBurn keeps every visible CPU busy for 5 s and checks that the
// process got no more CPU time than the cgroup quota allows.
func TestCPUBurn(t *testing.T) {
	limit := cgroupCPULimit(t)

	const wall = 5 * time.Second
	workers := runtime.NumCPU()
	// Go sizes GOMAXPROCS from the cgroup quota; undo it so the quota itself is tested
	runtime.GOMAXPROCS(workers)
	before := cpuTime(t)
	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			x := 0
			for time.Since(start) < wall {
				x++
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	used := cpuTime(t) - before
	cores := used.Seconds() / elapsed.Seconds()
	t.Logf("workers=%d wall=%v cpu=%v used_cores=%.2f limit_cores=%.2f", workers, elapsed.Round(time.Millisecond), used.Round(time.Millisecond), cores, limit)
	logFile(t, "/sys/fs/cgroup/cpu.stat")

	if limit <= 0 {
		t.Logf("no CPU quota visible, behaviour recorded without comparison")
		return
	}
	// 10% margin for scheduler accounting
	if cores > limit*1.1 {
		t.Errorf("used %.2f cores, limit is %.2f", cores, limit)
	}
}

// cgroupCPULimit reads cpu.max ("quota period") and returns cores, or 0 for no limit.
func cgroupCPULimit(t *testing.T) float64 {
	b, err := os.ReadFile("/sys/fs/cgroup/cpu.max")
	if err != nil {
		// no cgroup files under gVisor: keep going, only the measured behaviour is recorded
		t.Logf("cpu.max not readable: %v", err)
		return 0
	}
	t.Logf("cpu.max: %s", strings.TrimSpace(string(b)))
	f := strings.Fields(string(b))
	if len(f) != 2 || f[0] == "max" {
		return 0
	}
	quota, _ := strconv.ParseFloat(f[0], 64)
	period, _ := strconv.ParseFloat(f[1], 64)
	return quota / period
}

func cpuTime(t *testing.T) time.Duration {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		t.Fatalf("getrusage: %v", err)
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

func logFile(t *testing.T, path string) {
	if b, err := os.ReadFile(path); err == nil {
		t.Logf("%s:\n%s", path, b)
	}
}
