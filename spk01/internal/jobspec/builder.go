package jobspec

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
	TestFilter            string // Regex for -test.run flag (e.g. "^TestAddSuccess$")
	ActiveDeadlineSeconds *int64
	TTLSecondsAfterFinish *int32
	BaselineProfile       bool // If true, disables security hardening (used only for negative controls)
}

// Build creates a batchv1.Job specification based on the provided options.
// By default, it applies the full hardened security profile verified in SPK-01 (SC-05, SC-06).
func Build(opts Options) *batchv1.Job {
	namePrefix := opts.NamePrefix
	if namePrefix == "" {
		namePrefix = "spk01-"
	}

	image := opts.Image
	if image == "" {
		image = "gojudge-spk01:dev"
	}

	// Execution uses precompiled test binary through test2json framing
	cmd := []string{
		"/usr/local/bin/test2json",
		"-t",
		"/usr/local/bin/trusted.test",
		"-test.v=test2json",
	}
	if opts.TestFilter != "" {
		cmd = append(cmd, "-test.run", opts.TestFilter)
	}

	podSpec := corev1.PodSpec{
		RestartPolicy:                 corev1.RestartPolicyNever,
		TerminationGracePeriodSeconds: ptr.To[int64](1),
	}

	container := corev1.Container{
		Name:            "judge",
		Image:           image,
		ImagePullPolicy: corev1.PullNever, // Use pre-imported local image
		Command:         cmd,
	}

	if !opts.BaselineProfile {
		// Hardened Pod-level isolation
		podSpec.AutomountServiceAccountToken = ptr.To(false)
		podSpec.SecurityContext = &corev1.PodSecurityContext{
			RunAsNonRoot: ptr.To(true),
			RunAsUser:    ptr.To[int64](65532),
			RunAsGroup:   ptr.To[int64](65532),
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
		}

		// Hardened Container-level isolation
		container.SecurityContext = &corev1.SecurityContext{
			AllowPrivilegeEscalation: ptr.To(false),
			ReadOnlyRootFilesystem:   ptr.To(true),
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
			},
		}

		// Strict resource boundaries (SC-03)
		container.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("250m"),
				corev1.ResourceMemory: resource.MustParse("512Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("512Mi"),
			},
		}

		// Ephemeral writable tmpfs with size boundary (SC-05)
		container.Env = []corev1.EnvVar{
			{Name: "HOME", Value: "/tmp"},
		}
		container.VolumeMounts = []corev1.VolumeMount{
			{Name: "tmp", MountPath: "/tmp"},
		}
		podSpec.Volumes = []corev1.Volume{
			{
				Name: "tmp",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{
						SizeLimit: ptr.To(resource.MustParse("256Mi")),
					},
				},
			},
		}
	}

	podSpec.Containers = []corev1.Container{container}

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
				Spec: podSpec,
			},
		},
	}
}
