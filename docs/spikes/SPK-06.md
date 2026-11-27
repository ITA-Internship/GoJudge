# SPK-06: Isolation Validation and Security Regression Suite (Spike Report)

## Metadata
- **Spike ID:** SPK-06 (issue #50, stretch ST-03)
- **Type:** Technical spike, evidence only
- **Owner:** Pair B
- **Timebox:** 3.5 pair-days (Sprint 2, stretch).
- **Starting point:** SPK-05 Conditional go: gVisor with `network=none` proposed as primary candidate, with most other controls untested under `runsc` (SPK-05, conditions 1 and 2)
- **Outcome:** Conditional go. The suite is defined and was run repeatably on two runtimes. Six of seven control areas are verified under gVisor (E). The process limit under gVisor is not verified (H). Reproduction by another pair is pending
- **Experiment code:** branch `spike/SPK-06-security-suite`, directory `spk06/` (not merged to `main`). Evidence commit: [`83429ee28b2cdee3dae5d4b2a6a47bdd236dbeba`](https://github.com/ITA-Internship/GoJudge/tree/83429ee28b2cdee3dae5d4b2a6a47bdd236dbeba/spk06)
- **Out of scope:** repeating the SPK-05 comparison, a new DEC-08 recommendation, approving DEC-08 or DEC-02, selecting a production mechanism, the production judge, learner-submitted code

## Evidence classes
Same as SPK-05.
- **E**: run in the spike environment (command, versions, result, evidence location recorded)
- **D**: cited documentation or an earlier report (source recorded, not re-run)
- **H**: untested or undocumented; a risk, never a finding

## Environment

Isolated experimental environment: a dedicated local k3d cluster `spk06`, one server node, created for this spike. No shared or production infrastructure was used or changed. Only trusted fixtures were run: no learner-submitted code was executed.

| | Baseline run (runc) | gVisor run |
|---|---|---|
| Dates | 25 Nov 2026 | 26-27 Nov 2026 |
| Cluster | k3s v1.35.5, containerd 2.2.3-k3s1 (k3d, image from `infra/k3d/k3d-config.yaml`) | same. k3d v5.9.0, Docker Engine 29.2.1 recorded for this machine only |
| Node kernel | `6.12.76-linuxkit` | `6.6.87.2-microsoft-standard-WSL2` |
| Runtime | default (`runc`), no RuntimeClass | RuntimeClass `gvisor-nonet`, handler `runsc-nonet` (`network=none`), `runsc` release-20261123.1 |
| Fixture image | `gojudge-spk06:dev`, built from `spk06/image` | same source |
| Network policy | deny-egress `NetworkPolicy` in namespace `spk06` | same, `network=none` also set by the handler |

The two runs were made on different machines, so the node kernels differ (see Limitations). The baseline cpu and process checks were run with the first version of those two tests; the gVisor run used the changed version (F-5).

### Experimental parameters
Not approved production limits. Taken from SPK-01, except the pids limit, which is our own choice.

| Parameter | Value | Source |
|---|---|---|
| CPU | limit 1, request 250m | SPK-01 |
| Memory | 512Mi | SPK-01 |
| Time | `activeDeadlineSeconds: 5` | SPK-01 |
| Output | 1 MiB (`outcap`) | SPK-01 |
| `/tmp` | `emptyDir` 256Mi | SPK-01 |
| Processes | 128 per pod (kubelet `pod-max-pids`) | our choice: safe for one node, enough for a Go test |

Final limits are not defined and are agreed separately (see Open questions). Where no threshold exists, the check records the behaviour and draws no conclusion about compliance.

## Success criteria status

| ID | Criterion (from the brief) | Status | Evidence |
|---|---|---|---|
| S1 | All seven control areas have documented validation cases and evidence statuses | Verified | SC-A matrix below |
| S2 | The repeatable execution mechanism and the evidence format are documented | Verified | `spk06/README.md`, `run-suite.sh`, `*.jsonl` per area |
| S3 | Supported experimental runs and their limitations are reported | Partially verified | Runs in E1-E17. Not verified: the pids limit under gVisor, a second pair's reproduction, Cilium + gVisor (not run). See Limitations |

## SC-A: Validation matrix

One row per control area. Requirement source: MVP Proposal §7 and Product Backlog §3.1; the SPK-05 requirement is given as R-number. Evidence IDs are in the next section.

| Area | Requirement source | Fixture / check | Expected observation or open threshold | Class | Evidence | Limitations |
|---|---|---|---|---|---|---|
| Network | Proposal §7 "no network access"; SPK-05 R2 | `TestNetworkEgress`: the first and only action is a TCP dial to `1.1.1.1:443`, `8.8.8.8:53`, `google.com:443` and the API server, 3 s dial timeout | All 4 targets blocked. No open threshold | E | E1 (baseline FAIL 0/5), E8 (gVisor PASS 5/5) | 4 targets, one host. The RuntimeClass is chosen per Pod; nothing here enforces it |
| CPU | Proposal §7 "CPU limits"; SPK-05 R4 | `TestCPUBurn`: busy loop on all visible CPUs for 5 s, reports CPU time used | Behaviour only. The limit of 1 core is an experimental parameter, not an approved threshold | E | E2, E9, E16 | Under gVisor the in-sandbox figure is too high; the host figure is from one probe Pod (F-3) |
| Memory | Proposal §7 "memory limits"; SPK-05 R4; SPK-01 SC-03 | `TestMemoryExhaustion`: allocates and touches 10 MiB chunks until killed | Pod `OOMKilled`. The runner does not check how much was allocated | E | E3, E10, E17 | The fixture has no own cap, it relies on the memory limit. Allocation reached 2500 MiB inside the sandbox before the kill (F-2) |
| Timeout | Proposal §7 "time limits"; SPK-05 R4; SPK-01 SC-02 | `TestInfiniteLoop` with `activeDeadlineSeconds: 5` | Job `DeadlineExceeded` within 5 + 5 s | E | E4, E11 | The Pod is deleted when the deadline passes, so the log is empty. The evidence is the Job status |
| Filesystem | Proposal §7 "read-only/minimal filesystem"; SPK-05 R3; SPK-01 SC-05 | `TestFilesystemWrite`: `/` read-only, writes to 5 paths fail, `/tmp` stays writable | Only `/tmp` writable. No open threshold | E | E5, E12 | `readOnlyRootFilesystem` plus `emptyDir`; not tested with learner code |
| Output | Proposal §7 "bounded output"; SPK-05 R5; SPK-01 SC-07 | `TestOutputFlood` writes 20 MiB, wrapped by `outcap -limit 1048576` | Behaviour only: log capped near 1 MiB, verdict kept | E | E6, E13 | `outcap` is part of the image. A bypass through `/proc/1/fd/1` is H (SPK-05) |
| Process | Proposal §7 "process limits"; SPK-05 R4 | `TestProcessLimit`: starts `sleep` processes up to 512 until the kernel refuses | Behaviour only. The limit of 128 (kubelet `pod-max-pids`) is our experimental choice | E on baseline, H under gVisor | E7, E14, E15 | Under gVisor the stop (41-42) is not caused by the pids limit; the limit itself is not verified (F-4) |

## SC-B: Per-case outcomes

Runs: network 5, other areas 3. All runs completed; no run is missing. Evidence directory names are relative to `spk06/results/` in the evidence commit.

| ID | Area | Runtime | Runs | Result | Evidence |
|---|---|---|---|---|---|
| E1 | Network | runc | 5 | FAIL 0/5: all 4 targets reached in run 1, `1.1.1.1:443` in 38-90 ms over the 5 runs (SPK-01 F-03) | `k8s (baseline)/network-*` |
| E2 | CPU | runc | 3 | 1.00-1.01 cores used (10 CPUs visible), limit 1 | `k8s (baseline)/cpu-*` |
| E3 | Memory | runc | 3 | `OOMKilled` (137), killed after 500 MiB | `k8s (baseline)/memory-*` |
| E4 | Timeout | runc | 3 | `DeadlineExceeded`, Job failed after 6 s | `k8s (baseline)/timeout-*` |
| E5 | Filesystem | runc | 3 | PASS 3/3 | `k8s (baseline)/filesystem-*` |
| E6 | Output | runc | 3 | log 1,048,888-1,049,051 bytes, verdict kept | `k8s (baseline)/output-*` |
| E7 | Process | runc | 3 | stopped at 119 processes, `resource temporarily unavailable` (pod limit 128) | `k8s (baseline)/process-*` |
| E8 | Network | gVisor | 5 | PASS 5/5: all 4 targets `network is unreachable` from 0 s, including DNS and the API server | `gVisor/network-*` |
| E9 | CPU | gVisor | 3 | 1.44, 1.51, 1.46 cores used inside the sandbox (2 CPUs visible), limit 1 | `gVisor/cpu-*` |
| E10 | Memory | gVisor | 3 | `OOMKilled` (137) in 3/3, after 2500 MiB allocated inside the sandbox | `gVisor/memory-*` |
| E11 | Timeout | gVisor | 3 | `DeadlineExceeded`, Job failed after 6 s | `gVisor/timeout-*` |
| E12 | Filesystem | gVisor | 3 | PASS 3/3 | `gVisor/filesystem-*` |
| E13 | Output | gVisor | 3 | log 1,048,935-1,048,955 bytes, verdict kept | `gVisor/output-*` |
| E14 | Process | gVisor | 3 | stopped at 42 processes, `cannot allocate memory`. `pids.max` file not available | `gVisor/process-*` |
| E15 | Process | gVisor, memory 2Gi | 3 | stopped at 42, 42, 41 processes, `cannot allocate memory` (same result as E14) | `gVisor-process-2Gi/` |
| E16 | CPU, host side | gVisor | 1 | Pod cgroup on the node: `cpu.max` `100000 100000`, ~5.15 s of CPU in ~5 s, throttled in 50 of 65 periods. In-sandbox figure in the same run: 1.47 cores | `gVisor/host-probes/cpu-*` |
| E17 | Memory, host side | gVisor | 1 | Pod cgroup on the node: `memory.max` 536870912, `memory.peak` 512 MiB, `oom_kill 1`. In-sandbox figure in the same run: 2500 MiB | `gVisor/host-probes/memory-*` |

E16 and E17 are single runs of a probe Pod with the suite limits, started with `kubectl exec` and sampled on the node every 0.5 s (`spk06/host-probe.sh`). They are not Jobs.

The end of `memory-test.log` (`waiting on pid 3 ... EOF`, `command terminated with exit code 128`) is expected in E17: the sandbox is killed by the OOM kill, so the `kubectl exec` session loses it. The Pod state in `memory-pod-state.json` is `OOMKilled`, exit code 137.

## Findings

| ID | Finding | Class |
|---|---|---|
| F-1 | Network: runc fails 0/5 (the start-up window of SPK-01 F-03 is reproduced). gVisor with `network=none` blocks all 4 targets in 5/5 runs, from 0 s. This repeats SPK-05 E12 in k3d with the same fixture | E (E1, E8) |
| F-2 | Memory: the limit works on the host, the Pod cgroup stayed at 512 MiB and one OOM kill was recorded (E17). The fixture reported 2500 MiB inside the sandbox before the kill (E10), against 500 MiB on runc (E3). Why the allocation went beyond the limit is not established | E for both figures, H for the cause |
| F-3 | CPU: the limit works on the host (~1.0 core, throttled, E16). The in-sandbox `getrusage` figure is 1.44-1.51 cores, too high compared with the host. The cause is not established; Sentry may account CPU time itself | E for both figures, H for the cause |
| F-4 | Process: under gVisor the fixture stops at 41-42 processes with `cannot allocate memory`, not at the pids limit of 128. Raising the Pod memory to 2Gi did not change it (E15), so it is not the container memory limit. Whether a pids limit applies to processes inside the sandbox is not verified | E for the stop, H for the pids limit |
| F-5 | Under gVisor the files `/sys/fs/cgroup/cpu.max` and `pids.max` do not exist (the Pod sees the Sentry kernel, `4.19.0-gvisor`). The first gVisor run of cpu and process failed at 0 s because the two tests read them first. The tests now log the missing file and measure the behaviour anyway (`cpu_test.go`, `procs_test.go`). Behaviour on runc is unchanged because the files exist there | E |
| F-6 | The runner marks memory as PASS for any `OOMKilled` and does not look at the peak, so the memory PASS in E10 does not say how much was allocated | E |
| F-7 | Timeout, filesystem and output behave the same on runc and gVisor | E (E4-E6, E11-E13) |

**Under `runsc`, in-sandbox metrics (CPU time, allocated memory, cgroup files) cannot be used as evidence that a limit is respected (F-2, F-3, F-5).** Enforcement is verified from the host cgroup of the Pod (E16, E17) or from the Pod and Job status (`OOMKilled`, `DeadlineExceeded`).

Documentation-based reading (D): "All Sentry threads and Sentry memory are subject to a container cgroup" (gVisor Resource Model, SPK-05 source G9). This agrees with F-2 and F-3 (host cgroup enforces, the sandbox reports its own figures) but is not proof of their cause.

## Execution mechanism and reproducibility

Full instructions are in `spk06/README.md` on the spike branch. Summary:

- **Prerequisites:** Docker, k3d, kubectl, jq, perl. For gVisor: `runsc`, `containerd-shim-runsc-v1` and the `gvisor-bin` folder from the gVisor release archive. On Windows, Git Bash
- **Invocation:**
  ```bash
  POD_MAX_PIDS=128 ./spk06/setup.sh          # cluster, image, namespace, deny-egress policy
  ./spk06/gvisor-setup.sh                    # gVisor into the node, RuntimeClass gvisor-nonet
  RUNTIME_CLASS=gvisor-nonet RESULTS=spk06/results/gVisor ./spk06/run-suite.sh
  ./spk06/host-probe.sh cpu                  # host-side cgroup probe, also: memory
  ```
- **Per-run evidence:** for every run, `<area>-runNN.log`, `-job.json`, `-pod.yaml` and one line in `<area>.jsonl` (area, run index, runtime class, limits, Job and container state, exit code, verdict, log bytes, duration). `env.txt` records the context, versions and image. `env.txt` is rewritten on every `run-suite.sh` invocation, so it describes the last invocation in that folder
- **Cleanup:** `./spk06/cleanup.sh` deletes the cluster. Results stay in `spk06/results/`. `gvisor-setup.sh` has to be repeated after every new `setup.sh`, because its files live inside the node
- **Reproduction by another pair:** not done yet. It is required by the acceptance criteria and is open

## Failure and stop criteria

None of the three stop criteria from the brief was met.
- No fixture affected resources outside the experimental cluster. The only outbound traffic is the network fixture's TCP dials, which are the check itself
- No check required learner-submitted code
- No cluster-wide or shared infrastructure change: the kubelet `pod-max-pids` argument and gVisor are set only on the node of the `spk06` cluster

## Limitations and residual risks

- One node, one k3s version, 3 runs per check (5 for network)
- The baseline and gVisor runs were made on different machines (node kernels `6.12.76-linuxkit` and `6.6.87.2-microsoft-standard-WSL2`), so the comparison is indicative. The baseline cpu and process checks used the first version of the two tests
- Process limit under gVisor: not verified (F-4, H)
- CPU and memory under gVisor are verified on the host from one probe Pod each (E16, E17), not from a Job. The cause of the in-sandbox figures is not established (H)
- The suite is not tied to a mechanism, but only runc and gVisor were run. Cilium + gVisor was not run: no agreed environment
- The runtime is chosen per Pod; nothing here stops a Pod from using another `runtimeClassName` (SPK-05)
- With `network=none` the Pod has no DNS and no API server, so results leave only through stdout and the exit code
- Only trusted fixtures were run. Compile time and overhead of learner code under `runsc`, TTL cleanup, cancellation and non-root enforcement under `runsc` were not tested (H, SPK-05 condition 1)
- A second pair has not yet reproduced the runs

## Open questions
- Final limits for CPU, memory, processes and output. The values used here are experimental parameters (SC-A), not requirements. They are agreed separately, with these results as input
- Whether the pids limit applies to processes inside a gVisor sandbox.
- Cause of the 2500 MiB allocation and of the 41-42 process stop under gVisor

## Outcome

**Conditional go.** A reusable suite with seven bounded trusted fixtures, a runner and an evidence format exists and was executed repeatably on runc and on gVisor with `network=none`. Network, timeout, filesystem and output are verified on both runtimes. CPU and memory limits are verified on the host under gVisor, while the in-sandbox figures are not reliable. The process limit under gVisor is not verified.

Conditions for a full go: reproduction by a second pair, a check of the process limit under gVisor, and a decision on final limits.

This report selects and approves no production isolation mechanism and does not approve DEC-08 or DEC-02.

### Follow-on decisions (listed, not made)
- Agree the final CPU, memory, process and output limits, using E1-E17 as input
- Decide whether host-side cgroup metrics become the standard way to verify CPU and memory under gVisor
- Decide whether the runner's memory check should also record the peak
- Decide on a Cilium + gVisor run once an environment is agreed

## Sources
- SPK-05 report: `docs/spikes/SPK-05.md` (R2-R5, conditions 1 and 2, gVisor source G9)
- SPK-01 report: `docs/spikes/SPK-01.md` (experimental limits, SC-02, SC-03, SC-05, SC-07, F-03)
- MVP Proposal §7: `docs/product/mvp-proposal.md`
- Product Backlog §3.1: `docs/product/product-backlog.md`
- Experiment code and logs: branch `spike/SPK-06-security-suite`, `spk06/`, evidence commit `83429ee28b2cdee3dae5d4b2a6a47bdd236dbeba`
