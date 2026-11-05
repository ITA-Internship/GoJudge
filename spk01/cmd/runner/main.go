package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ITA-Internship/GoJudge/spk01/internal/jobspec"
	"github.com/ITA-Internship/GoJudge/spk01/internal/measure"
	"github.com/ITA-Internship/GoJudge/spk01/internal/runner"
)

func main() {
	var (
		scFlag         = flag.String("sc", "SC-01", "Spike criterion identifier to execute (e.g. SC-01)")
		runsFlag       = flag.Int("runs", 10, "Number of consecutive evaluation runs to execute")
		namespaceFlag  = flag.String("namespace", "spk01", "Kubernetes namespace for evaluation Jobs")
		imageFlag      = flag.String("image", "gojudge-spk01:dev", "Container image used for evaluation")
		kubeconfigFlag = flag.String("kubeconfig", "", "Path to kubeconfig file (defaults to ~/.kube/config or $KUBECONFIG)")
		resultsDirFlag = flag.String("results-dir", "results", "Directory where JSONL and log evidence are stored")
		filterFlag     = flag.String("filter", "^TestAddSuccess$", "Test filter regex passed to go test -run")
	)
	flag.Parse()

	sc := strings.ToUpper(*scFlag)
	if err := os.MkdirAll(*resultsDirFlag, 0755); err != nil {
		log.Fatalf("Failed to ensure results directory: %v", err)
	}

	fmt.Println("==================================================================")
	fmt.Printf("SPK-01 Runner | Criterion: %s | Target Runs: %d\n", sc, *runsFlag)
	fmt.Printf("Namespace: %s | Image: %s\n", *namespaceFlag, *imageFlag)
	fmt.Println("==================================================================")

	// Initialize client-go
	cs, err := runner.NewClient(*kubeconfigFlag)
	if err != nil {
		log.Fatalf("Error connecting to Kubernetes: %v", err)
	}

	resultsFile := filepath.Join(*resultsDirFlag, fmt.Sprintf("%s.jsonl", strings.ToLower(sc)))
	successCount := 0
	var totalStartupMs, totalExecMs int64

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for i := 1; i <= *runsFlag; i++ {
		if ctx.Err() != nil {
			fmt.Println("\nExecution interrupted by user signal.")
			break
		}

		fmt.Printf("[%02d/%02d] Starting evaluation run... ", i, *runsFlag)

		opts := jobspec.Options{
			NamePrefix: fmt.Sprintf("spk01-%s-", strings.ToLower(sc)),
			Namespace:  *namespaceFlag,
			Criterion:  sc,
			RunIndex:   i,
			Image:      *imageFlag,
			TestFilter: *filterFlag,
		}

		if sc == "SC-02" {
			// For SC-02: target infinite loop with 5s active deadline
			opts.TestFilter = "^TestInfiniteLoop$"
			deadline := int64(5)
			opts.ActiveDeadlineSeconds = &deadline
		} else if sc == "SC-10" {
			// For SC-10: target infinite loop, runner sends manual cancellation
			opts.TestFilter = "^TestInfiniteLoop$"
		} else if sc == "SC-09" {
			// For SC-09: target normal pass test with 10s TTL controller cleanup
			opts.TestFilter = "^TestAddSuccess$"
			ttl := int32(10)
			opts.TTLSecondsAfterFinish = &ttl
		} else if sc == "SC-08" {
			// For SC-08: mixed scenarios (4 pass, 3 assertion fail, 3 deadline timeout)
			switch {
			case i <= 4:
				opts.TestFilter = "^TestAddSuccess$"
			case i <= 7:
				opts.TestFilter = "^TestFailAssertion$"
			default:
				opts.TestFilter = "^TestInfiniteLoop$"
				deadline := int64(5)
				opts.ActiveDeadlineSeconds = &deadline
			}
		}

		record, runErr := runner.RunOnce(ctx, cs, opts, *resultsDirFlag)
		if runErr != nil {
			fmt.Printf("FAILED: %v\n", runErr)
			continue
		}

		// Persist evidence record
		if err := measure.AppendRecord(resultsFile, *record); err != nil {
			log.Printf("Warning: failed to append record: %v", err)
		}

		if record.Pass {
			successCount++
			totalStartupMs += record.StartupLatencyMs
			totalExecMs += record.ExecutionLatencyMs
			if sc == "SC-02" {
				fmt.Printf("PASS | Job: %s | Condition: %s (%s) | Overrun: +%dms\n",
					record.JobName, record.Condition, record.Reason, record.OverrunMs)
			} else if sc == "SC-10" {
				fmt.Printf("PASS | Job: %s | Condition: %s | CancelLatency: %dms (<=30s)\n",
					record.JobName, record.Condition, record.CancelLatencyMs)
			} else if sc == "SC-09" {
				fmt.Printf("PASS | Job: %s | Condition: %s | TTL Actual: %dms (target <= 70000ms)\n",
					record.JobName, record.Condition, record.TTLActualMs)
			} else if sc == "SC-08" {
				fmt.Printf("PASS | Job: %s | Condition: %s (%s) | Verdict: %s | Log: %s\n",
					record.JobName, record.Condition, record.Reason, record.PackageVerdict, record.LogFile)
			} else {
				fmt.Printf("PASS | Job: %s | Startup: %dms | Exec: %dms | Verdict: %s\n",
					record.JobName, record.StartupLatencyMs, record.ExecutionLatencyMs, record.PackageVerdict)
			}
		} else {
			if sc == "SC-02" {
				fmt.Printf("FAIL | Job: %s | Condition: %s | Reason: %s | Overrun: +%dms\n",
					record.JobName, record.Condition, record.Reason, record.OverrunMs)
			} else if sc == "SC-10" {
				fmt.Printf("FAIL | Job: %s | Condition: %s | CancelLatency: %dms (exceeded 30s)\n",
					record.JobName, record.Condition, record.CancelLatencyMs)
			} else if sc == "SC-09" {
				fmt.Printf("FAIL | Job: %s | Condition: %s | TTL Actual: %dms (exceeded 70000ms)\n",
					record.JobName, record.Condition, record.TTLActualMs)
			} else if sc == "SC-08" {
				fmt.Printf("FAIL | Job: %s | Condition: %s (%s) | Verdict: %s | Log: %s\n",
					record.JobName, record.Condition, record.Reason, record.PackageVerdict, record.LogFile)
			} else {
				fmt.Printf("FAIL | Job: %s | Condition: %s | Reason: %s | Verdict: %s\n",
					record.JobName, record.Condition, record.Reason, record.PackageVerdict)
			}
		}

		// Small interval between runs to ensure clean API state
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("Summary for %s: %d of %d runs passed.\n", sc, successCount, *runsFlag)
	if successCount > 0 {
		fmt.Printf("Average Startup Latency: %d ms | Average Execution: %d ms\n",
			totalStartupMs/int64(successCount), totalExecMs/int64(successCount))
	}

	if successCount == *runsFlag {
		fmt.Printf("STATUS: [VERIFIED] Criterion %s met the threshold (%d/%d).\n", sc, successCount, *runsFlag)
		fmt.Printf("Evidence stored at: %s\n", resultsFile)
	} else {
		fmt.Printf("STATUS: [NOT VERIFIED] (%d/%d passed)\n", successCount, *runsFlag)
	}
	fmt.Println("==================================================================")
}
