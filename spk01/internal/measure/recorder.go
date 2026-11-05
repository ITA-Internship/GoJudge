package measure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RunRecord represents evidence collected from a single Job execution run.
type RunRecord struct {
	Criterion          string    `json:"criterion"`
	RunIndex           int       `json:"run_index"`
	JobName            string    `json:"job_name"`
	PodName            string    `json:"pod_name"`
	CreatedAt          time.Time `json:"created_at"`
	PodRunningAt       time.Time `json:"pod_running_at,omitempty"`
	JobFinishedAt      time.Time `json:"job_finished_at"`
	StartupLatencyMs   int64     `json:"startup_latency_ms"`
	ExecutionLatencyMs int64     `json:"execution_latency_ms"`
	TotalDurationMs    int64     `json:"total_duration_ms"`
	OverrunMs          int64     `json:"overrun_ms,omitempty"`
	CancelLatencyMs    int64     `json:"cancel_latency_ms,omitempty"`
	TTLActualMs        int64     `json:"ttl_actual_ms,omitempty"`
	Condition          string    `json:"condition"` // "Complete" or "Failed"
	Reason             string    `json:"reason,omitempty"`
	TerminatedReason   string    `json:"terminated_reason,omitempty"`
	ExitCode           int32     `json:"exit_code"`
	TestsRun           int       `json:"tests_run"`
	TestsPassed        int       `json:"tests_passed"`
	TestsFailed        int       `json:"tests_failed"`
	PackageVerdict     string    `json:"package_verdict"`
	Pass               bool      `json:"pass"`
	LogFile            string    `json:"log_file"`
}

// AppendRecord writes the RunRecord as a single line of JSON into the target file.
func AppendRecord(filePath string, record RunRecord) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
