package jobspec

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// Options holds parameters to configure a Kubernetes Job for execution.
type Options struct {
	NamePrefix            string
	Namespace             string
	Criterion             string // e.g. "SC-01", "SC-02"
	RunIndex              int
	Image                 string
	TestFilter            string // Regex for -run flag (e.g. "^TestAddSuccess$")
	ActiveDeadlineSeconds *int64
	TTLSecondsAfterFinish *int32
}

// Build creates a batchv1.Job specification based on the provided options.
// This builder is shared between Engineer 1 and Engineer 2.
// Engineer 2 will add securityContext, resource limits, and volumes here.
func Build(opts Options) *batchv1.Job {
	namePrefix := opts.NamePrefix
	if namePrefix == "" {
		namePrefix = "spk01-"
	}

	image := opts.Image
	if image == "" {
		image = "gojudge-spk01:dev"
	}

	cmd := []string{"go", "test", "-json"}
	if opts.TestFilter != "" {
		cmd = append(cmd, "-run", opts.TestFilter)
	}
	cmd = append(cmd, "./...")

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: namePrefix,
			Namespace:    opts.Namespace,
			Labels: map[string]string{
				"app":       "spk01",
				"criterion": opts.Criterion,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            ptr.To[int32](0), // No automatic retries
			ActiveDeadlineSeconds:   opts.ActiveDeadlineSeconds,
			TTLSecondsAfterFinished: opts.TTLSecondsAfterFinish,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "spk01",
						"criterion": opts.Criterion,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy:                 corev1.RestartPolicyNever,
					TerminationGracePeriodSeconds: ptr.To[int64](1),
					Containers: []corev1.Container{
						{
							Name:            "judge",
							Image:           image,
							ImagePullPolicy: corev1.PullNever, // Use pre-imported local image
							Command:         cmd,
						},
					},
				},
			},
		},
	}
}
