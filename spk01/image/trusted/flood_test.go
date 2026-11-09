package trusted

import (
	"os"
	"strings"
	"testing"
)

// TestOutputFlood writes 20 MiB to stdout in 1 KiB lines and then passes.
// Without a limiter the whole flood reaches the container log.
func TestOutputFlood(t *testing.T) {
	line := strings.Repeat("x", 1023) + "\n" // 1 KiB per line
	const total = 20 << 20                   // 20 MiB in total
	for written := 0; written < total; written += len(line) {
		if _, err := os.Stdout.WriteString(line); err != nil {
			t.Fatalf("write: %v", err)
		}
	}
	_, _ = os.Stdout.WriteString("flood-end\n")
}
