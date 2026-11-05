package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ITA-Internship/GoJudge/spk01/internal/jobspec"
	"github.com/ITA-Internship/GoJudge/spk01/internal/measure"
	"github.com/ITA-Internship/GoJudge/spk01/internal/report"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// NewClient initializes a Kubernetes clientset using the given or discovered kubeconfig.
func NewClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("KUBECONFIG")
	}
	if kubeconfigPath == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfigPath = filepath.Join(home, ".kube", "config")
		}
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig from %s: %w", kubeconfigPath, err)
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return cs, nil
}

// RunOnce executes a single Job, waits for completion, retrieves logs, and records metrics.
func RunOnce(ctx context.Context, cs *kubernetes.Clientset, opts jobspec.Options, outputDir string) (*measure.RunRecord, error) {
	startTime := time.Now()

	// 1. Build and submit Job manifest
	jobManifest := jobspec.Build(opts)
	createdJob, err := cs.BatchV1().Jobs(opts.Namespace).Create(ctx, jobManifest, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Job: %w", err)
	}
	jobName := createdJob.Name

	if opts.Criterion == "SC-10" {
		return runCancellationTest(ctx, cs, opts, outputDir, startTime, createdJob)
	}

	var (
		podName          string
		podRunningAt     time.Time
		finalCondition   string
		conditionReason  string
		jobFinishedAt    time.Time
		exitCode         int32
		terminatedReason string
		logContent       bytes.Buffer
		logStreamDone    = make(chan struct{})
		streamStarted    bool
		streamCtx, cancelStream = context.WithCancel(ctx)
	)
	defer cancelStream()

	// 2. Poll until Job succeeds or fails, tracking Pod phase transitions and streaming logs live
	pollErr := wait.PollUntilContextTimeout(ctx, 100*time.Millisecond, 90*time.Second, true, func(ctx context.Context) (bool, error) {
		// Discover Pod name
		if podName == "" {
			pods, err := cs.CoreV1().Pods(opts.Namespace).List(ctx, metav1.ListOptions{
				LabelSelector: "batch.kubernetes.io/job-name=" + jobName,
			})
			if err == nil && len(pods.Items) > 0 {
				podName = pods.Items[0].Name
			}
		}

		// Inspect Pod status and launch live log streaming as soon as Pod is Running
		if podName != "" {
			pod, err := cs.CoreV1().Pods(opts.Namespace).Get(ctx, podName, metav1.GetOptions{})
			if err == nil {
				// Record startup transition
				if podRunningAt.IsZero() && (pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed) {
					podRunningAt = time.Now()
				}

				// Capture container exit state while Pod is still alive
				for _, cs := range pod.Status.ContainerStatuses {
					if cs.Name == "judge" && cs.State.Terminated != nil {
						exitCode = cs.State.Terminated.ExitCode
						terminatedReason = cs.State.Terminated.Reason
					}
				}

				// Start live stream
				if !streamStarted && (pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed) {
					streamStarted = true
					go func(targetPod string) {
						defer close(logStreamDone)
						req := cs.CoreV1().Pods(opts.Namespace).GetLogs(targetPod, &corev1.PodLogOptions{
							Container: "judge",
							Follow:    true,
						})
						stream, err := req.Stream(streamCtx)
						if err != nil {
							return
						}
						defer stream.Close()
						_, _ = io.Copy(&logContent, stream)
					}(podName)
				}
			}
		}

		// Check Job completion/failure conditions
		j, err := cs.BatchV1().Jobs(opts.Namespace).Get(ctx, jobName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		for _, c := range j.Status.Conditions {
			if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == corev1.ConditionTrue {
				finalCondition = string(c.Type)
				conditionReason = c.Reason
				jobFinishedAt = time.Now()
				return true, nil
			}
		}

		return false, nil
	})

	if pollErr != nil {
		return nil, fmt.Errorf("job %s did not reach completion within timeout: %w", jobName, pollErr)
	}

	// Wait safely for log stream to complete. Cancel stream context if lingering to avoid data race.
	if streamStarted {
		select {
		case <-logStreamDone:
		case <-time.After(1500 * time.Millisecond):
			cancelStream()
			<-logStreamDone
		}
	} else if podName != "" {
		// Fallback direct read if stream did not start
		req := cs.CoreV1().Pods(opts.Namespace).GetLogs(podName, &corev1.PodLogOptions{
			Container: "judge",
		})
		stream, err := req.Stream(ctx)
		if err == nil {
			defer stream.Close()
			_, _ = io.Copy(&logContent, stream)
		}
	}

	// 5. Save raw log file
	logFileName := fmt.Sprintf("%s-run%02d.log", opts.Criterion, opts.RunIndex)
	logPath := filepath.Join(outputDir, logFileName)
	_ = os.WriteFile(logPath, logContent.Bytes(), 0644)

	// 6. Parse test2json output
	testSummary, err := report.Parse(bytes.NewReader(logContent.Bytes()))
	if err != nil {
		testSummary = &report.Summary{PackageVerdict: "parse_error"}
	}

	// 7. Calculate timings
	startupMs := int64(0)
	if !podRunningAt.IsZero() {
		startupMs = podRunningAt.Sub(startTime).Milliseconds()
	}
	execMs := int64(0)
	if !podRunningAt.IsZero() {
		execMs = jobFinishedAt.Sub(podRunningAt).Milliseconds()
	}
	totalMs := jobFinishedAt.Sub(startTime).Milliseconds()
	var overrunMs int64

	// 8. Determine Criterion pass condition
	isPass := false
	var ttlActualMs int64
	switch opts.Criterion {
	case "SC-01":
		// SC-01 passes if Job Completed, package verdict is "pass", tests were executed (>0), and 0 tests failed
		isPass = (finalCondition == "Complete") && (testSummary.PackageVerdict == "pass") && (testSummary.TestsRun > 0) && (testSummary.TestsFailed == 0)
	case "SC-02":
		// SC-02 passes if Job Failed with DeadlineExceeded within deadline + 5s
		deadlineSeconds := int64(10)
		if opts.ActiveDeadlineSeconds != nil {
			deadlineSeconds = *opts.ActiveDeadlineSeconds
		}
		expectedDeadline := startTime.Add(time.Duration(deadlineSeconds) * time.Second)
		overrunDuration := jobFinishedAt.Sub(expectedDeadline)
		overrunMs = overrunDuration.Milliseconds()
		isPass = (finalCondition == "Failed") && (conditionReason == "DeadlineExceeded") && (overrunDuration <= 5*time.Second)
	case "SC-09":
		// Handled below during TTL observation
	default:
		isPass = (finalCondition == "Complete")
	}

	if opts.Criterion == "SC-09" {
		// In SC-09, we DO NOT delete the Job! We wait for K8s TTL controller to remove it.
		ttlSeconds := int64(10)
		if opts.TTLSecondsAfterFinish != nil {
			ttlSeconds = int64(*opts.TTLSecondsAfterFinish)
		}
		var ttlDeletedAt time.Time
		ttlWaitErr := wait.PollUntilContextTimeout(ctx, 500*time.Millisecond, 80*time.Second, true, func(ctx context.Context) (bool, error) {
			_, jErr := cs.BatchV1().Jobs(opts.Namespace).Get(ctx, jobName, metav1.GetOptions{})
			jobGone := apierrors.IsNotFound(jErr)
			podGone := true
			if podName != "" {
				_, pErr := cs.CoreV1().Pods(opts.Namespace).Get(ctx, podName, metav1.GetOptions{})
				podGone = apierrors.IsNotFound(pErr)
			}
			if jobGone && podGone {
				ttlDeletedAt = time.Now()
				return true, nil
			}
			return false, nil
		})
		if ttlWaitErr == nil {
			ttlDuration := ttlDeletedAt.Sub(jobFinishedAt)
			ttlActualMs = ttlDuration.Milliseconds()
			// Requirement: removed within TTL + 60s
			isPass = (finalCondition == "Complete") && (ttlDuration <= time.Duration(ttlSeconds+60)*time.Second)
		} else {
			isPass = false
		}
	} else {
		// 9. Clean up Job after data extraction for other criteria (Foreground deletion cascades to Pod)
		deletePolicy := metav1.DeletePropagationForeground
		_ = cs.BatchV1().Jobs(opts.Namespace).Delete(ctx, jobName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
	}

	record := &measure.RunRecord{
		Criterion:          opts.Criterion,
		RunIndex:           opts.RunIndex,
		JobName:            jobName,
		PodName:            podName,
		CreatedAt:          startTime,
		PodRunningAt:       podRunningAt,
		JobFinishedAt:      jobFinishedAt,
		StartupLatencyMs:   startupMs,
		ExecutionLatencyMs: execMs,
		TotalDurationMs:    totalMs,
		OverrunMs:          overrunMs,
		TTLActualMs:        ttlActualMs,
		Condition:          finalCondition,
		Reason:             conditionReason,
		TerminatedReason:   terminatedReason,
		ExitCode:           exitCode,
		TestsRun:           testSummary.TestsRun,
		TestsPassed:        testSummary.TestsPassed,
		TestsFailed:        testSummary.TestsFailed,
		PackageVerdict:     testSummary.PackageVerdict,
		Pass:               isPass,
		LogFile:            filepath.Join("results", logFileName),
	}

	return record, nil
}

// runCancellationTest verifies SC-10: manual cancellation terminates a running Job within 30 seconds.
func runCancellationTest(ctx context.Context, cs *kubernetes.Clientset, opts jobspec.Options, outputDir string, startTime time.Time, job *batchv1.Job) (*measure.RunRecord, error) {
	jobName := job.Name
	var (
		podName       string
		podRunningAt  time.Time
		logContent    bytes.Buffer
		logStreamDone = make(chan struct{})
		streamCtx, cancelStream = context.WithCancel(ctx)
	)
	defer cancelStream()

	// 1. Wait until Pod is in Running phase
	err := wait.PollUntilContextTimeout(ctx, 100*time.Millisecond, 30*time.Second, true, func(ctx context.Context) (bool, error) {
		pods, err := cs.CoreV1().Pods(opts.Namespace).List(ctx, metav1.ListOptions{
			LabelSelector: "batch.kubernetes.io/job-name=" + jobName,
		})
		if err == nil && len(pods.Items) > 0 {
			p := pods.Items[0]
			podName = p.Name
			if p.Status.Phase == corev1.PodRunning {
				podRunningAt = time.Now()
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, fmt.Errorf("pod for job %s did not become Running: %w", jobName, err)
	}

	// 2. Start live log stream
	go func(targetPod string) {
		defer close(logStreamDone)
		req := cs.CoreV1().Pods(opts.Namespace).GetLogs(targetPod, &corev1.PodLogOptions{
			Container: "judge",
			Follow:    true,
		})
		stream, err := req.Stream(streamCtx)
		if err != nil {
			return
		}
		defer stream.Close()
		_, _ = io.Copy(&logContent, stream)
	}(podName)

	// Let the workload run actively for 2 seconds before cancellation
	time.Sleep(2 * time.Second)

	// 3. Issue manual cancel request
	cancelRequestedAt := time.Now()
	deletePolicy := metav1.DeletePropagationForeground
	err = cs.BatchV1().Jobs(opts.Namespace).Delete(ctx, jobName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, fmt.Errorf("failed to issue delete for job %s: %w", jobName, err)
	}

	// 4. Poll until BOTH Job and Pod are completely deleted from etcd
	var jobFinishedAt time.Time
	pollErr := wait.PollUntilContextTimeout(ctx, 100*time.Millisecond, 40*time.Second, true, func(ctx context.Context) (bool, error) {
		_, jobErr := cs.BatchV1().Jobs(opts.Namespace).Get(ctx, jobName, metav1.GetOptions{})
		jobGone := apierrors.IsNotFound(jobErr)

		podGone := true
		if podName != "" {
			_, podErr := cs.CoreV1().Pods(opts.Namespace).Get(ctx, podName, metav1.GetOptions{})
			podGone = apierrors.IsNotFound(podErr)
		}

		if jobGone && podGone {
			jobFinishedAt = time.Now()
			return true, nil
		}
		return false, nil
	})
	if pollErr != nil {
		return nil, fmt.Errorf("job %s was not deleted within timeout: %w", jobName, pollErr)
	}

	// 5. Ensure log stream flushes cleanly without race
	select {
	case <-logStreamDone:
	case <-time.After(1500 * time.Millisecond):
		cancelStream()
		<-logStreamDone
	}

	// 6. Save raw logs
	logFileName := fmt.Sprintf("%s-run%02d.log", opts.Criterion, opts.RunIndex)
	logPath := filepath.Join(outputDir, logFileName)
	_ = os.WriteFile(logPath, logContent.Bytes(), 0644)

	cancelDuration := jobFinishedAt.Sub(cancelRequestedAt)
	cancelLatencyMs := cancelDuration.Milliseconds()
	isPass := cancelDuration <= 30*time.Second

	record := &measure.RunRecord{
		Criterion:          opts.Criterion,
		RunIndex:           opts.RunIndex,
		JobName:            jobName,
		PodName:            podName,
		CreatedAt:          startTime,
		PodRunningAt:       podRunningAt,
		JobFinishedAt:      jobFinishedAt,
		StartupLatencyMs:   podRunningAt.Sub(startTime).Milliseconds(),
		ExecutionLatencyMs: cancelDuration.Milliseconds(),
		TotalDurationMs:    jobFinishedAt.Sub(startTime).Milliseconds(),
		CancelLatencyMs:    cancelLatencyMs,
		Condition:          "Cancelled",
		Reason:             "UserCancelRequest",
		Pass:               isPass,
		LogFile:            filepath.Join("results", logFileName),
	}

	return record, nil
}
