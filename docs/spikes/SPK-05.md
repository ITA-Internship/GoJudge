# SPK-05: Isolation Mechanism Comparison (Spike Report)

## Metadata
- **Spike ID:** SPK-05 (issue #32)
- **Type:** Technical spike, evidence and recommendation only
- **Owner:** Pair B
- **Timebox:** 3 pair-days (Sprint 2)
- **Starting point:** SPK-01 No-go (F-03: egress succeeded before per-Pod NetworkPolicy rules were installed)
- **Out of scope:** approving DEC-08 / DEC-02, production judge code, running learner code

## Evidence classes
- **E**: run in the spike environment (command, image version, result, log reference recorded)
- **D**: cited documentation or the SPK-01 report (version and source recorded, not re-run)
- **H**: untested and undocumented; a risk, never a finding

## Candidates
- **A. Baseline:** Kubernetes Job + `NetworkPolicy` on k3s (SPK-01 as-is). Reference column.
- **B. Init-container gate:** baseline + init container that waits until egress is blocked (SPK-01 Option 1).
- **C. Cilium:** alternative network plugin (SPK-01 Option 2). Calico is not covered in this report.
- **D. gVisor (`runsc`):** sandboxed runtime (SPK-01 Option 3).

## SC-A: Requirement-by-candidate matrix

Each cell: `class` + short result + source key (see Sources). Documentation was read on 18 Nov 2026. Cilium docs: "stable" = 1.20.5. gVisor docs show no version number. Where a docs page is silent, the cell is **H**, not D.

| # | Requirement (proposal §7) | A. Baseline | B. Init gate | C. Cilium | D. gVisor |
|---|---|---|---|---|---|
| R1 | Unprivileged execution (non-root, no privilege escalation) | **E**: uid 65532, `setuid(0)` EPERM, NoNewPrivs=1 (SPK-01 SC-06) | **E**: same Pod spec (SPK-01 gate manifest) | **D**: not a CNI feature, `securityContext` is set on the Pod [K1]. No Cilium page says it changes it | **H**: docs do not state that Pod `securityContext` is honored under `runsc`; compatibility page warns of "known feature gaps" in Kubernetes [G8]. Inside the sandbox a user can only manipulate virtualized resources [G1] |
| R2 | No network access (egress blocked) | **E**: **Not verified**, F-03 met (SPK-01 SC-04). **D**: NetworkPolicy needs a plugin that enforces it [K2] | **E**: blocked at t=1 ms in 5/5 runs, limitations SPK-01 §5.3 (`sc-04-gate.jsonl`) | **D**: in default mode endpoints "have unrestricted network access until selected by policy" [C1]; `always` mode enforces "on all endpoints even if no rules select" them [C1], resulting verdict not spelled out. Policy is programmed in eBPF when the endpoint is regenerated [C2]. **H**: timing relative to container start is not documented [C2] | **E**: standalone `runsc --network=none do`, 6/6 runs: all 3 targets `network is unreachable` within 0-3 ms, DNS lookup fails too (E2). Control with default network: all 3 targets connected, test FAIL (E3). **D**: `--network=none` cuts external networking, only loopback remains [G4]. **H**: in Kubernetes, per-handler `network=none` in containerd [G6], interplay with NetworkPolicy and log retrieval are not tested [G3, G5] |
| R3 | Read-only / minimal filesystem | **E**: `/` is `ro`, writes fail EROFS/EACCES, `/tmp` emptyDir 256Mi (SC-05) | **E**: same spec | **D**: not a CNI concern, same as A (no Cilium page covers it) | **D**: without overlay "you will have a read-only rootfs"; overlay is a writable tmpfs layer, memory-backed overlay "can bloat container memory usage" [G7]. **H**: emptyDir and `readOnlyRootFilesystem` interplay is not described [G7] |
| R4 | CPU, memory and time limits | **E**: time (SC-02), memory OOMKilled (SC-03). CPU limit: **H** | **E**: same as A | **D**: limits are enforced by kubelet [K3], not by the CNI | **D**: "All Sentry threads and Sentry memory are subject to a container cgroup" [G9]; gVisor "relies on the host resource mechanisms (cgroups)" [G1]. **H**: how Kubernetes limits and `activeDeadlineSeconds` map to the sandbox is not covered [G9] |
| R5 | Bounded output | **E**: `outcap` truncates 20 MiB to ~1 MiB (SC-07). `/proc/1/fd/1` bypass: **H** | **E**: same as A | **D**: output handling is not a network concern, same as A | **H**: no gVisor page covers container stdout/log size limits |
| R6 | Log retrieval, TTL cleanup, cancellation | **E**: SC-08, SC-09 (~10 s), SC-10 (~1.1 s) | **E**: same as A (gate adds ~1 s) | **D**: same Kubernetes Job mechanisms as A; Cilium does not touch them (no page covers it) | **H**: Job TTL and cancellation are Kubernetes-side, not documented for `runsc`. **D**: sandbox is deleted before CNI `DEL` on cleanup [G10] |
| R7 | No app secrets / no Kubernetes API credentials | **E**: no SA token mounted (SC-06). API server `10.43.0.1:443` reachable in the window (SC-04) | **E**: token absent, API server blocked after gate | **D**: egress to the API server can be restricted with L3-L7 `CiliumNetworkPolicy` [C3]. **H**: window timing as in R2 | **H**: docs do not cover Pod service-account token handling under `runsc` |
| R8 | Protection of hidden tests | **H**: depends on problem packaging (DEC-09), not on the mechanism | **H**: same | **H**: same | **H**: same |
| R9 | RBAC, namespace isolation, quotas, Pod security | **H**: only namespace `spk01` + NetworkPolicy used; RBAC and quotas not measured | **H**: same | **D**: namespace-scoped `CiliumNetworkPolicy` and `CiliumClusterwideNetworkPolicy` exist [C3]. RBAC, quotas, Pod security stay Kubernetes-side: **H** | **D**: selected per Pod with `runtimeClassName` [G11]. RBAC, quotas, Pod security: **H** |
| R10 | Stronger boundary than a plain container ("Docker alone is not sufficient") | **D**: host kernel is shared; a container workload is "only one system call away from host compromise" [G2] | **D**: same as A | **D**: Cilium is a network layer, not a sandbox. Threat model: after container escape, traffic "will no longer be subject to" Cilium network policies [C4] | **D**: "never passes through any system call to the host"; escape needs both a Sentry and a host-kernel exploit [G2]. Limits: "not a substitute for a secure architecture", no protection against hardware side channels [G1] |

## SC-C: Isolation-gap register (network egress timing)

Reference finding: SPK-01 F-03. With a deny-egress NetworkPolicy, the first connection succeeded after 19 ms, because the k3s controller installs per-Pod rules only after the container starts (window under ~1 s, ~1 s resolution).

| Candidate | Egress blocked from first instruction? | Evidence | Window / limitation | Named risk |
|---|---|---|---|---|
| A. Baseline | No | E (Not verified) | Under ~1 s; established connections survive (SPK-01 §5.1-2) | Egress in the start-up window |
| B. Init gate | Yes for the judge container, from t=1 ms | E (5/5) | SPK-01 §5.3 (a-f): single probe target, ~1 s resolution, +1 s start-up, `busybox` in the supply chain | Residual limits listed in SPK-01 §5.3 |
| C. Cilium | Not documented | D for default semantics, **H** for timing | In default mode the endpoint is unrestricted until selected by policy [C1]; start-up ordering is not described [C2]. `always` mode might remove the gap but its start-up effect is not documented | Same shape as F-03 unless `always` mode or equivalent is verified. Not verified (H) |
| D. gVisor | `network=none` in standalone `runsc do`: yes, no network is configured, so there is no start-up window by construction (6/6 runs). In Kubernetes: not verified. With NetworkPolicy: not documented | E for standalone `runsc do` (E2, E3), D for the mode itself [G4], **H** for Kubernetes, NetworkPolicy and timing | Not run through containerd/Kubernetes. `network=none` per runtime handler is not confirmed [G6]. API server target not probed (variable exists only inside a Pod). 3 targets, 1 host (WSL2) | Not verified in Kubernetes (H). Unknown whether Job lifecycle and log retrieval work with no Pod network |

A "blocked from the first instruction" claim counts as verified only with E evidence. D-only is recorded as "Documented only".

## Experiments run (E evidence)

Environment: Windows host, WSL2 (kernel `6.6.87.2-microsoft-standard-WSL2`), Ubuntu 24.04.5 LTS, `runsc release-20261116.0` (spec 1.2.1). Workload: precompiled `trusted.test` from the SPK-01 image (`spk01/image`), trusted tests only. Raw logs: `spk05/gvisor/wsl_results/` on branch `spike/SPK-05-isolation`, evidence commit `248148568618091270b65a33c4b45f5fdbeec398`. Commands are listed in `spk05/gvisor/README.md`.

| ID | Experiment | Result |
|---|---|---|
| E1 | Compatibility: `TestAddSuccess` natively, and under `runsc do` (default network and `--network=none`) | PASS in all three: the Go test binary runs under gVisor |
| E2 | `TestNetworkEgress` under `runsc --network=none do`, 1 initial run + 5 loop runs | 6/6 PASS: 1.1.1.1:443, 8.8.8.8:53 and google.com:443 all `network is unreachable` (0-3 ms), DNS lookup also unreachable |
| E3 | Control: `TestNetworkEgress` under `runsc do` with default network | FAIL (expected): all 3 targets connected in 11, 26 and 46 ms, so the test detects an open network |

Limitations: `runsc do` is a quick test command, not the Kubernetes path (no `runtimeClassName`, containerd, Job, NetworkPolicy). One host. The Kubernetes API server target is not probed because the variable it comes from exists only inside a Pod. No `runc` network control was run on the same host.

## SC-B: DEC-08 recommendation

TODO: recommended option, rejected options with reasons. Nothing is adopted without MA-03.

## Sources
- **K1** Kubernetes, Security Context (page last modified 2026-09-27): https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
- **K2** Kubernetes, Network Policies (last modified 2026-04-14): https://kubernetes.io/docs/concepts/services-networking/network-policies/
- **K3** Kubernetes, Resource Management for Pods and Containers: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
- **C1** Cilium 1.20.5, Policy Enforcement Modes: https://docs.cilium.io/en/stable/security/policy/intro/
- **C2** Cilium 1.20.5, Endpoint Lifecycle: https://docs.cilium.io/en/stable/security/policy/lifecycle/
- **C3** Cilium 1.20.5, Kubernetes Network Policy: https://docs.cilium.io/en/stable/network/kubernetes/policy/
- **C4** Cilium 1.20.5, Threat Model: https://docs.cilium.io/en/stable/security/threat-model/
- **G1** gVisor, Security Model: https://gvisor.dev/docs/architecture_guide/security/
- **G2** gVisor, Introduction to gVisor security: https://gvisor.dev/docs/architecture_guide/intro/
- **G3** gVisor, Networking Guide: https://gvisor.dev/docs/architecture_guide/networking/
- **G4** gVisor, Networking (user guide): https://gvisor.dev/docs/user_guide/networking/
- **G5** gVisor, Production guide: https://gvisor.dev/docs/user_guide/production/
- **G6** gVisor, Containerd Advanced Configuration: https://gvisor.dev/docs/user_guide/containerd/configuration/
- **G7** gVisor, Filesystem: https://gvisor.dev/docs/user_guide/filesystem/
- **G8** gVisor, Application compatibility: https://gvisor.dev/docs/user_guide/compatibility/
- **G9** gVisor, Resource Model: https://gvisor.dev/docs/architecture_guide/resources/
- **G10** gVisor, Using CNI: https://gvisor.dev/docs/tutorials/cni/
- **G11** gVisor, Kubernetes Quick Start: https://gvisor.dev/docs/user_guide/quick_start/kubernetes/
- gVisor experiment logs: `spk05/gvisor/wsl_results/` (branch `spike/SPK-05-isolation`)
- SPK-01 report: `docs/spikes/SPK-01.md`, evidence commit `17939ad95e6fddb064bca8b90057124f6d8daab1`
- MVP Proposal §7: `docs/product/mvp-proposal.md`