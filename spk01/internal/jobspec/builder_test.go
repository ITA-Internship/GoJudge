package jobspec

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestBuildHardenedByDefault(t *testing.T) {
	deadline := int64(5)
	ttl := int32(10)
	job := Build(Options{
		NamePrefix:            "spk01-sc-01-",
		Namespace:             "spk01",
		Criterion:             "SC-01",
		Image:                 "gojudge-spk01:dev",
		TestFilter:            "^TestAddSuccess$",
		ActiveDeadlineSeconds: &deadline,
		TTLSecondsAfterFinish: &ttl,
	})

	if job.Spec.BackoffLimit == nil || *job.Spec.BackoffLimit != 0 {
		t.Fatalf("expected BackoffLimit=0, got %v", job.Spec.BackoffLimit)
	}
	if job.Spec.ActiveDeadlineSeconds == nil || *job.Spec.ActiveDeadlineSeconds != 5 {
		t.Fatalf("expected ActiveDeadlineSeconds=5")
	}
	if job.Spec.TTLSecondsAfterFinished == nil || *job.Spec.TTLSecondsAfterFinished != 10 {
		t.Fatalf("expected TTLSecondsAfterFinished=10")
	}

	podSpec := job.Spec.Template.Spec
	if podSpec.AutomountServiceAccountToken == nil || *podSpec.AutomountServiceAccountToken != false {
		t.Errorf("expected AutomountServiceAccountToken=false")
	}
	if podSpec.TerminationGracePeriodSeconds == nil || *podSpec.TerminationGracePeriodSeconds != 1 {
		t.Errorf("expected TerminationGracePeriodSeconds=1")
	}
	if podSpec.SecurityContext == nil || podSpec.SecurityContext.RunAsNonRoot == nil || !*podSpec.SecurityContext.RunAsNonRoot {
		t.Errorf("expected Pod RunAsNonRoot=true")
	}
	if podSpec.SecurityContext.RunAsUser == nil || *podSpec.SecurityContext.RunAsUser != 65532 {
		t.Errorf("expected RunAsUser=65532")
	}
	if podSpec.SecurityContext.SeccompProfile == nil || podSpec.SecurityContext.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
		t.Errorf("expected SeccompProfile RuntimeDefault")
	}

	if len(podSpec.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(podSpec.Containers))
	}
	c := podSpec.Containers[0]
	if c.SecurityContext == nil {
		t.Fatalf("expected container SecurityContext")
	}
	if c.SecurityContext.AllowPrivilegeEscalation == nil || *c.SecurityContext.AllowPrivilegeEscalation != false {
		t.Errorf("expected AllowPrivilegeEscalation=false")
	}
	if c.SecurityContext.ReadOnlyRootFilesystem == nil || !*c.SecurityContext.ReadOnlyRootFilesystem {
		t.Errorf("expected ReadOnlyRootFilesystem=true")
	}
	if c.SecurityContext.Capabilities == nil || len(c.SecurityContext.Capabilities.Drop) != 1 || c.SecurityContext.Capabilities.Drop[0] != "ALL" {
		t.Errorf("expected Capabilities Drop=[ALL]")
	}
}

func TestBuildBaselineProfile(t *testing.T) {
	job := Build(Options{
		NamePrefix:      "spk01-baseline-",
		Namespace:       "spk01",
		Criterion:       "SC-06-control",
		BaselineProfile: true,
	})

	podSpec := job.Spec.Template.Spec
	if podSpec.SecurityContext != nil {
		t.Errorf("expected nil PodSecurityContext in baseline profile")
	}
	if podSpec.Containers[0].SecurityContext != nil {
		t.Errorf("expected nil Container SecurityContext in baseline profile")
	}
}
