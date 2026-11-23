# SPK-05: Isolation Mechanism Comparison (Spike Report)

## Metadata
- **Spike ID:** SPK-05 (issue #32)
- **Type:** Technical spike, evidence and recommendation only
- **Owner:** Pair B
- **Timebox:** 3 pair-days (Sprint 2)
- **Starting point:** SPK-01 No-go (F-03: egress succeeded before per-Pod NetworkPolicy rules were installed)
- **Outcome:** evidence collected for all four candidates; Cilium and gVisor with `network=none` block egress from the first attempt (E); recommendation in SC-B; nothing is approved
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
| R1 | Unprivileged execution (non-root, no privilege escalation) | **E**: uid 65532, `setuid(0)` EPERM, NoNewPrivs=1 (SPK-01 SC-06) | **E**: same Pod spec (SPK-01 gate manifest) | **D**: not a CNI feature, `securityContext` is set on the Pod [K1]. No Cilium page says it changes it | **E**: a Job with the SPK-01 `securityContext` (runAsNonRoot uid 65532, drop ALL, seccomp RuntimeDefault, `readOnlyRootFilesystem`) is accepted and runs to completion under gVisor (E12). **H**: enforcement (e.g. `setuid(0)` EPERM) is not tested under `runsc`; compatibility page warns of "known feature gaps" in Kubernetes [G8]. Inside the sandbox a user can only manipulate virtualized resources [G1] |
| R2 | No network access (egress blocked) | **E**: **Not verified**, F-03 met (SPK-01 SC-04). Reproduced on k3s v1.37.2 in WSL2: runc + NetworkPolicy, `TestNetworkEgress` FAIL 5/5 (E10). The policy does block once the Pod is running: `Connection refused` from a running Pod (E11). **D**: NetworkPolicy needs a plugin that enforces it [K2] | **E**: blocked at t=1 ms in 5/5 runs, limitations SPK-01 §5.3 (`sc-04-gate.jsonl`) | **E**: SC-04 on k3d + Cilium 1.20.5, 5/5 runs: all 4 targets blocked from the first attempt (`i/o timeout` after the test's 3s dial timeout, packets are dropped). Control without the policy: all 4 targets connected (E5). **D**: in default mode endpoints "have unrestricted network access until selected by policy" [C1]; policy is programmed in eBPF when the endpoint is regenerated [C2] | **E**: in k3s, RuntimeClass `gvisor-nonet` (handler `runsc-nonet` with `network=none`), no NetworkPolicy: `TestNetworkEgress` PASS 5/5, all 4 targets `network is unreachable`, including the DNS server and the API server (E12). Standalone `runsc --network=none do`: 6/6 PASS, control with default network FAIL (E2, E3). The plain `gvisor` handler (default network) does **not** block: FAIL 5/5 without policy and 5/5 with NetworkPolicy (E8, E9). **D**: `--network=none` cuts external networking, only loopback remains [G4]. **H**: a Pod that gets another `runtimeClassName`, or none, has a normal network; nothing in this spike enforces the choice. Single host and k3s version |
| R3 | Read-only / minimal filesystem | **E**: `/` is `ro`, writes fail EROFS/EACCES, `/tmp` emptyDir 256Mi (SC-05) | **E**: same spec | **D**: not a CNI concern, same as A (no Cilium page covers it) | **D**: without overlay "you will have a read-only rootfs"; overlay is a writable tmpfs layer, memory-backed overlay "can bloat container memory usage" [G7]. A Job with `readOnlyRootFilesystem` and an emptyDir `/tmp` starts under gVisor (E12). **H**: writes were not tested; emptyDir and `readOnlyRootFilesystem` interplay is not described [G7] |
| R4 | CPU, memory and time limits | **E**: time (SC-02), memory OOMKilled (SC-03). CPU limit: **H** | **E**: same as A | **D**: limits are enforced by kubelet [K3], not by the CNI | **D**: "All Sentry threads and Sentry memory are subject to a container cgroup" [G9]; gVisor "relies on the host resource mechanisms (cgroups)" [G1]. **H**: how Kubernetes limits and `activeDeadlineSeconds` map to the sandbox is not covered [G9] |
| R5 | Bounded output | **E**: `outcap` truncates 20 MiB to ~1 MiB (SC-07). `/proc/1/fd/1` bypass: **H** | **E**: same as A | **D**: output handling is not a network concern, same as A | **H**: no gVisor page covers container stdout/log size limits |
| R6 | Log retrieval, TTL cleanup, cancellation | **E**: SC-08, SC-09 (~10 s), SC-10 (~1.1 s) | **E**: same as A (gate adds ~1 s) | **D**: same Kubernetes Job mechanisms as A; Cilium does not touch them (no page covers it) | **E**: the Job completed and `kubectl logs` returned the test output with no Pod network (E12). **H**: Job TTL and cancellation are not tested under `runsc`. **D**: sandbox is deleted before CNI `DEL` on cleanup [G10] |
| R7 | No app secrets / no Kubernetes API credentials | **E**: no SA token mounted (SC-06). API server `10.43.0.1:443` reachable in the window (SC-04) | **E**: token absent, API server blocked after gate | **E**: API server `10.43.0.1:443` blocked in 5/5 runs, reached in the control (E4, E5). **D**: egress to the API server can be restricted with L3-L7 `CiliumNetworkPolicy` [C3]. SA token: same Pod spec as A | **E**: with `network=none` the API server `10.43.0.1:443` is unreachable (E12). **H**: ServiceAccount token handling under `runsc` is not tested (the manifest sets `automountServiceAccountToken: false`) |
| R8 | Protection of hidden tests | **H**: depends on problem packaging (DEC-09), not on the mechanism | **H**: same | **H**: same | **H**: same |
| R9 | RBAC, namespace isolation, quotas, Pod security | **H**: only namespace `spk01` + NetworkPolicy used; RBAC and quotas not measured | **H**: same | **D**: namespace-scoped `CiliumNetworkPolicy` and `CiliumClusterwideNetworkPolicy` exist [C3]. RBAC, quotas, Pod security stay Kubernetes-side: **H** | **D**: selected per Pod with `runtimeClassName` [G11]. **E**: a Pod with `runtimeClassName: gvisor` runs under gVisor (`uname -r` = `4.19.0-gvisor`; a plain Pod shows the host kernel) (E6). RBAC, quotas, Pod security: **H** |
| R10 | Stronger boundary than a plain container ("Docker alone is not sufficient") | **D**: host kernel is shared; a container workload is "only one system call away from host compromise" [G2] | **D**: same as A | **D**: Cilium is a network layer, not a sandbox. Threat model: after container escape, traffic "will no longer be subject to" Cilium network policies [C4] | **D**: "never passes through any system call to the host"; escape needs both a Sentry and a host-kernel exploit [G2]. **E**: the Pod sees the Sentry kernel, not the host kernel (E6). Limits: "not a substitute for a secure architecture", no protection against hardware side channels [G1] |

## SC-C: Isolation-gap register (network egress timing)

Reference finding: SPK-01 F-03. With a deny-egress NetworkPolicy, the first connection succeeded after 19 ms, because the k3s controller installs per-Pod rules only after the container starts (window under ~1 s, ~1 s resolution).

| Candidate | Egress blocked from first instruction? | Evidence | Window / limitation | Named risk |
|---|---|---|---|---|
| A. Baseline | No | E (Not verified); reproduced on k3s v1.37.2 in WSL2 (E10) | Under ~1 s; established connections survive (SPK-01 §5.1-2). Policy is enforced once the Pod runs (E11) | Egress in the start-up window |
| B. Init gate | Yes for the judge container, from t=1 ms | E (5/5) | SPK-01 §5.3 (a-f): single probe target, ~1 s resolution, +1 s start-up, `busybox` in the supply chain | Residual limits listed in SPK-01 §5.3 |
| C. Cilium | Yes, no connection succeeded in any run (default mode, policy created before the Job) | E (5/5 + control, E4, E5) | One node, macOS Docker Desktop, 5 runs. The test's first dial starts right after the container starts; no finer timing. Drop reason (`Policy denied`) not captured. Policy created after the Pod and `always` mode not tested | Start-up ordering is not documented [C2], so the result can depend on load or version (H) |
| D. gVisor | With the `network=none` handler: yes, the Pod has no network, so there is no start-up window by construction. First probe at 0 s, 5/5 runs in k3s (E12), 6/6 standalone (E2). With the default-network handler: no. Connections succeeded 12-35 ms after start with a deny-egress NetworkPolicy (E9), same as runc (E10) | E for both handlers (E2, E9, E10, E12), D for the mode itself [G4] | `network=none` needs its own runtime handler and RuntimeClass. The Pod has no DNS and no API server access, so code and tests must be pushed in from outside and results must leave through stdout and the exit code. One host (WSL2), 4 targets, k3s v1.37.2 (SPK-01 used v1.35.5) | A Pod that selects another `runtimeClassName`, or none, runs with a normal network; no check in this spike rejects it (H). Behavior of the other criteria under `runsc` is mostly untested (see SC-A) |

A "blocked from the first instruction" claim counts as verified only with E evidence. D-only is recorded as "Documented only".

## Experiments run (E evidence)

### gVisor

**Stage 1, standalone `runsc do`.** Environment: Windows host, WSL2 (kernel `6.6.87.2-microsoft-standard-WSL2`), Ubuntu 24.04.5 LTS, `runsc release-20261116.0` (spec 1.2.1). Workload: precompiled `trusted.test` from the SPK-01 image (`spk01/image`), trusted tests only. Raw logs: `spk05/gvisor/wsl_results/` on branch `spike/SPK-05-isolation`, evidence commit `248148568618091270b65a33c4b45f5fdbeec398`. Commands are listed in `spk05/gvisor/README.md`.

| ID | Experiment | Result |
|---|---|---|
| E1 | Compatibility: `TestAddSuccess` natively, and under `runsc do` (default network and `--network=none`) | PASS in all three: the Go test binary runs under gVisor |
| E2 | `TestNetworkEgress` under `runsc --network=none do`, 1 initial run + 5 loop runs | 6/6 PASS: 1.1.1.1:443, 8.8.8.8:53 and google.com:443 all `network is unreachable` (0-3 ms), DNS lookup also unreachable |
| E3 | Control: `TestNetworkEgress` under `runsc do` with default network | FAIL (expected): all 3 targets connected in 11, 26 and 46 ms, so the test detects an open network |

Limitations (stage 1): `runsc do` is a quick test command, not the Kubernetes path (no `runtimeClassName`, containerd, Job, NetworkPolicy). One host. The Kubernetes API server target is not probed because the variable it comes from exists only inside a Pod.

**Stage 2, k3s with gVisor as a RuntimeClass.** Same WSL2 host, k3s `v1.37.2+k3s1`, single node. Workload: the SPK-01 `sc-04-netdial.yaml` Job (namespace `spk01`, `runAsNonRoot` uid 65532, `readOnlyRootFilesystem`, drop ALL, `automountServiceAccountToken: false`), with `runtimeClassName` added for the gVisor variants. 5 runs per row; PASS = all 4 targets blocked (1.1.1.1:443, 8.8.8.8:53, google.com:443, API server `10.43.0.1:443`). Raw logs: `spk05/gvisor/k8s_results/`, evidence commit `784eecd723cc049c3169b88a71cc32b6aaa5ae93`.

| ID | Experiment | Result |
|---|---|---|
| E6 | Smoke test: Pod `runtimeClassName: gvisor` vs plain Pod, `uname -r` | `Running`; `4.19.0-gvisor` vs host kernel `6.6.87.2-microsoft-standard-WSL2` |
| E7 | Control: runc, no NetworkPolicy | FAIL 5/5: egress open, the test detects it |
| E8 | gVisor (default network), no NetworkPolicy | FAIL 5/5: gVisor alone does not block egress |
| E9 | gVisor (default network) + deny-egress NetworkPolicy | FAIL 5/5. Run 1: 1.1.1.1:443 connected in 12 ms, 8.8.8.8:53 in 26 ms, google.com:443 in 35 ms; API server `connection refused` |
| E10 | Control: runc + deny-egress NetworkPolicy | FAIL 5/5: the F-03 start-up window is reproduced on this k3s |
| E11 | Is the NetworkPolicy enforced after start? Running runc Pod labelled `app=spk01`, 25 s after creation, `wget -T 5 http://1.1.1.1` | `Connection refused`, exit 1: the policy works once the Pod runs (`pol-check.txt`) |
| E12 | gVisor with `network=none` (RuntimeClass `gvisor-nonet`, handler `runsc-nonet`), no NetworkPolicy | **PASS 5/5** (`succeeded=1`, exit 0). Run 1: all 4 targets `network is unreachable` (0-3 ms), including DNS server `10.43.0.10:53` and API server `10.43.0.1:443` |

Limitations (stage 2):
- One host (WSL2), one node, one workload. k3s `v1.37.2` differs from SPK-01's `v1.35.5` and from the Cilium environment.
- No learner code was delivered into the Pod: the workload is a precompiled test binary baked into the image.
- Only the egress test ran under gVisor in Kubernetes. Read-only root writes, memory limit, output flood and non-root enforcement were not exercised under `runsc`.
- `network=none` is applied through a dedicated handler. Nothing here stops a Pod from using another RuntimeClass.
- Known WSL2 issue: Docker Desktop's `/Docker/host` mount breaks the k3s kubelet; workaround in `spk05/gvisor/README.md`.

### Cilium

Environment: macOS 27.0.1, Docker Desktop (Engine 29.5.3), k3d v5.9.0, k3s v1.36.4-k3s1, Cilium 1.20.5
k3s flannel and network policy controller are disabled, kube-proxy is kept, so only the NetworkPolicy engine differs from SPK-01.
Workload: SPK-01 image, `spk01/manifests/sc-04-netdial.yaml` and `sc-04-deny-egress.yaml` unchanged.
Logs: `spk05/cilium/results/` on branch `spike/SPK-05-isolation`. Commands are listed in `spk05/cilium/README.md`

| ID | Experiment | Result |
|---|---|---|
| E4 | SC-04 with the deny-egress NetworkPolicy, 5 runs | 5/5 PASS: 1.1.1.1:443, 8.8.8.8:53, google.com:443 and API server 10.43.0.1:443 all `i/o timeout` after 3 s, DNS lookup also timed out |
| E5 | Control: same Job without the NetworkPolicy, 1 run | FAIL (expected): all 4 targets connected (0-235 ms), so the network works and the policy does the blocking |

Limitations: one node on Docker Desktop. 5 runs. The policy exists before the Job is created. Default enforcement mode only

## SC-B: DEC-08 recommendation

**Status:** recommendation and evidence only. This report does not approve DEC-08, DEC-02 or any production judge work. **Nothing is adopted without MA-03.**

### What the decision turns on

DEC-08 asks which mechanism isolates the judge's execution of untrusted code. Of the ten requirements in SC-A, only two separate the candidates:

| Requirement | What the evidence shows |
|---|---|
| R2: egress blocked from the first attempt | A fails (E8). B, C and D with `network=none` block it on E evidence (SPK-01, E4/E5, E12) |
| R10: boundary stronger than a plain container | Only D provides it, and only on **D-class** evidence [G2]. Cilium is a network layer, not a sandbox [C4] |

R1, R3-R6 and R9 are handled by Kubernetes (Pod `securityContext`, limits, Jobs, namespace controls) and do not distinguish A, B and C. They are mostly **H** under `runsc` (see SC-A). R8 depends on DEC-09 (problem packaging), not on the isolation mechanism.

### Options and verdicts

| Option | Verdict | Reasons |
|---|---|---|
| **A. Baseline** (k3s NetworkPolicy) | Rejected | The start-up window of SPK-01 F-03 is reproduced on k3s v1.37.2: 5/5 runs connected with a deny-egress policy in place (E8). The policy does block once the Pod is running (E9), so the flaw is timing, not enforcement |
| **B. Init-container gate** | Rejected as the primary mechanism (research candidate only, as the sprint plan states) | Blocked from t=1 ms in 5/5 runs (SPK-01), but with the SPK-01 §5.3 limits: single probe target, ~1 s resolution, +1 s start-up, `busybox` in the supply chain, assumes the controller installs all rules at once. It does not change the boundary (R10) |
| **C. Cilium** | Fallback | 5/5 blocked with the policy created before the Job, control without the policy connected (E4, E5). Costs: it replaces the network plugin of the whole cluster (SPK-01 notes this needs agreement with Pair A), needs eBPF with `bpffs` and `cgroup2` mounted, start-up ordering is not documented [C2] (H). The boundary stays that of a plain container (R10) [C4] |
| **D. gVisor with `network=none`** | Proposed primary candidate, on the conditions below | 5/5 blocked in k3s with the first probe at 0 s, including DNS and the API server (E12); 6/6 standalone (E2). Without `network=none` gVisor does **not** close the window (E8, E9). Stronger boundary per documentation only (R10) [G2]. Costs: `runsc` on every node, a dedicated handler and RuntimeClass, the Pod has no network, so learner code and tests must be pushed in from outside as volumes and results leave only through stdout and the exit code |

### Recommendation

Take **D (gVisor with `network=none`)** to MA-03 as the primary candidate: it is the only option that combines closing the network window on E evidence (R2) with a stronger execution boundary (R10), and it changes only the Pods of the judge, not the whole cluster. Keep **C (Cilium)** as the fallback if the conditions below cannot be met. Reject **A** and **B** as the primary mechanism for the reasons above.

C and D work at different layers (network and runtime), so they are not exclusive. A combination was not tested (H). A namespace deny-egress NetworkPolicy can stay as a second layer; its interplay with `network=none` was not tested (H).

**Delivering code into the Pod.** With `network=none` the Pod cannot fetch anything itself. If D is chosen, the judge worker (outside the Pod) should push the learner code and the tests in as read-only volumes, for example a ConfigMap, which Kubernetes mounts from the node side without any Pod network, instead of the Pod pulling them. This fits the draft rule that the judge side loads source code and tests from PostgreSQL (CON-02 §3 draft, CON-04 P-07). The concrete mechanism (size limit and cleanup of the ConfigMap, RBAC for the worker, keeping hidden tests out of reach of learner code per CON-04 V-06) belongs to F07 and DEC-09. It was not tested here, because no learner code was run (H).

If neither C nor D can be demonstrated safe, the MVP proposal §7 already states the fallback: release stays internal-only or uses a vetted external sandbox service.

### What the recommendation rests on

- **E:** R2 for C and D (E4, E5, E12); the F-03 window reproduced on A (E8); a Pod under `runsc` sees the Sentry kernel, not the host kernel (E6).
- **D:** R10 for gVisor [G2]; `--network=none` semantics [G4]; Cilium is not a sandbox [C4].
- **H:** everything listed as a condition below.

### Conditions before MA-03 can adopt D (each is a named risk, class H)

1. **Other requirements under `runsc` are untested.** Non-root enforcement (`setuid(0)` EPERM), read-only root and `/tmp` writes, memory limit and OOM mapping, output cap, TTL cleanup and cancellation were not run under gVisor. Compiling learner code (`go test` with the learner file) was not run under gVisor either: the experiments used a precompiled test binary, so compile time and overhead under `runsc` are unknown. The compatibility page warns of "known feature gaps" in Kubernetes [G8].
2. **Environment coverage.** `runsc` was run on one host (WSL2, k3s v1.37.2). It was not run inside k3d, the local cluster tool in ADR-002 (Provisional), nor on the hosting that DEC-07 will choose. Start-up cost and overhead were not measured.
3. **The choice of RuntimeClass is per Pod.** A Pod with another `runtimeClassName`, or none, has a normal network. The Job builder (F07-02a, a pure Go function that can be unit-tested without a cluster) can set and test it; cluster-level enforcement (INF-10) was not tested.
4. **No network means input must be pushed in and results leave only through stdout and the exit code.** Log retrieval worked in the experiment (E12). Delivering learner code and tests as volumes was not tested (see above). The output cap and log limits from SPK-01 §5 item 4 still apply to F07-02b.
5. **Fallback C is not fully characterised.** Policy created after the Pod, `always` mode, behavior under load, multi-node, the drop reason (not captured) and eBPF support on the target hosts were not tested.
6. **Sample size.** One host per candidate, 4 probe targets, 5 runs, different k3s versions across candidates (SPK-01 v1.35.5, Cilium v1.36.4, gVisor v1.37.2).
7. **Security review.** Only trusted code was run. The MVP proposal §7 requires a dedicated security review and adversarial tests before public deployment.

### Outcome

Evidence was collected for all four candidates. Two of them (C, and D with `network=none`) block egress from the first attempt on E evidence. The recommendation above goes to MA-03. It is a recommendation: it adopts nothing and approves neither DEC-08 nor DEC-02.

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
- gVisor experiment logs: `spk05/gvisor/wsl_results/` (stage 1) and `spk05/gvisor/k8s_results/` (stage 2), branch `spike/SPK-05-isolation`
- Cilium experiment logs: `spk05/cilium/results/` (branch `spike/SPK-05-isolation`)
- SPK-01 report: `docs/spikes/SPK-01.md`, evidence commit `17939ad95e6fddb064bca8b90057124f6d8daab1`
- MVP Proposal §7: `docs/product/mvp-proposal.md`