# SPK-01: Isolated Go Test Execution in Kubernetes Jobs (Spike Report)

## Metadata
- **Spike ID:** SPK-01
- **Title:** Isolated Go Test Execution in Kubernetes Jobs
- **Type:** Technical spike
- **Owner:** Pair B (Engineer 1 & Engineer 2)
- **Timebox:** 8 pair-days (Sprint 1, Days 1–10: Mon 2 Nov – Fri 13 Nov 2026, 15:00)
- **Status:** Completed within timebox
- **Dependencies:** DEC-01 (Target platform, provisional), DEC-02 (K8s Job per submission, provisional), DEC-08 (Isolation mechanism), INF-03 (Local cluster infrastructure, k3d per ADR-002), MA-03
- **Experiment Code:** Branch `spike/SPK-01-job-lifecycle`, directory `spk01/` (not merged to `main`)
- **Permanent Evidence Commit:** [`17939ad95e6fddb064bca8b90057124f6d8daab1`](https://github.com/ITA-Internship/GoJudge/tree/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results)

---

## 1. Hypothesis & Objectives

**Hypothesis:** A Go test package can run inside a Kubernetes Job created through `client-go`, with CPU, memory and time limits, no network, a restricted filesystem, bounded output, log and report retrieval before cleanup, TTL cleanup and cancellation. Each control is verified by measurement, not assumed.

**Core Question:** Can Kubernetes Jobs orchestrating containers under standard k3s configurations reliably enforce all multi-tenant security and lifecycle boundaries required for untrusted code execution (MVP Proposal §7)?

---

## 2. Architecture & Experiment Setup

The spike implemented a dual-layer investigation:
1. **Orchestrator & Lifecycle (Engineer 1):** A pure Go runner (`spk01/cmd/runner`, `spk01/internal/runner`) utilizing official `k8s.io/client-go` to programmatically build, submit, stream logs, measure latency, verify conditions, and clean up Jobs.
2. **Sandbox Isolation & Adversarial Payloads (Engineer 2):** Container-level controls (`securityContext`, `seccomp`, `capabilities`, `emptyDir`), network policies (`NetworkPolicy`), output bounds (`outcap` wrapper), and attack payloads (`memhog`, `fswrite`, `privcheck`, `netdial`, `flood`).

```
+--------------------------------------------------------------------------------------------------+
|                                    SPK-01 ARCHITECTURE OVERVIEW                                  |
+--------------------------------------------------------------------------------------------------+
|                                                                                                  |
|   1. Runner (client-go)                                                                          |
|      +-------------+        Create Job        +----------------------------------------------+   |
|      |             | -----------------------> | Kubernetes Cluster (k3s v1.35.5-k3s1)        |   |
|      |  JobSpec    |                          |   Namespace: spk01                           |   |
|      |  Builder    |                          |   Pod: runAsUser: 65532, readOnlyRootFS: true|   |
|      +-------------+                          +----------------------------------------------+   |
|             |                                                        |                           |
|             | 2. Live log streaming (Follow: true)                   | Container Running         |
|             |    goroutine + context cancellation                    v                           |
|             v                                         +------------------------------+           |
|      +-------------+                                  | Image: gojudge-spk01:dev     |           |
|      | logContent  | <=============================== | /usr/local/bin/trusted.test  |           |
|      |  (Buffer)   |                                  |   wrapped by test2json/outcap|           |
|      +-------------+                                  +------------------------------+           |
|             |                                                        |                           |
|             | 3. Pod Terminated -> Stop Stream                       v                           |
|             |    Wait <-logStreamDone                 [Complete / Failed / TTL / Cancel]         |
|             v                                                        |                           |
|      +-------------+                                                 |                           |
|      |   Parser    | <--- Validates: TestsRun > 0 && TestsFailed==0  |                           |
|      +-------------+                                                 v                           |
|             |                                         +------------------------------+           |
|             v                                         | Cleanup: Foreground Delete   |           |
|      +-------------+                                  | OR K8s TTL Controller (SC-09)|           |
|      | jsonl + log |                                  +------------------------------+           |
|      +-------------+                                                                             |
+--------------------------------------------------------------------------------------------------+
```

### Key Engineering Findings & Fixes
- **Precompiled Test Binary (49s -> ~1s Execution):** Under non-root (`runAsUser: 65532`) and `readOnlyRootFilesystem: true`, Go cannot access `/root/.cache` and is forced to recompile dependencies into `/tmp/gocache` on every execution, taking ~49 seconds. Precompiling static test binaries during image build (`CGO_ENABLED=0 go test -c -o /usr/local/bin/trusted.test .`) and executing through `/usr/local/bin/test2json` completely eliminated compilation overhead, reducing execution latency to ~1.0s.
- **Protection Against `no tests to run` Trap:** In Go, if `-test.run` contains a typo, `go test` outputs a warning but exits with code 0 (`PASS`). The parser was hardened to track `TestsRun` and require `TestsRun > 0` for any passing verdict.
- **Data-Race-Free Log Capture:** Live log streaming runs under a cancellable `context.WithCancel`. Reading the log buffer is strictly blocked on channel closure `<-logStreamDone`, preventing concurrent read/write races.
- **Foreground Cascade Cleanup:** Calling `DeletePropagationForeground` guarantees that the Job object is retained until the Pod process is fully terminated and cleaned from etcd, preventing zombie pods.

### Job Specification Used for the Runs
- **Hardened specification (all criteria):** generated by `spk01/internal/jobspec/builder.go` and recorded in [`jobspec-hardened.json`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/jobspec-hardened.json) (`go run ./cmd/runner -dump-spec`). Pod level: `automountServiceAccountToken: false`, `runAsNonRoot: true`, `runAsUser/runAsGroup: 65532`, `seccompProfile: RuntimeDefault`, `terminationGracePeriodSeconds: 1`. Container level: `allowPrivilegeEscalation: false`, `readOnlyRootFilesystem: true`, `capabilities.drop: ["ALL"]`, requests `250m / 512Mi`, limits `1 CPU / 512Mi`, `HOME=/tmp`, `/tmp` as `emptyDir` with `sizeLimit: 256Mi`. The dump shows the default test filter; only `-test.run` differs between criteria.
- **Negative controls:** the same builder with `-baseline` (no hardening), recorded in [`baseline-unhardened/jobspec-baseline.json`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/baseline-unhardened/jobspec-baseline.json).
- **SC-01–03, SC-05, SC-06, SC-08–10** were run through the Go runner on this specification.
- **SC-04 and SC-07** were run with `kubectl` manifests that repeat the same security fields (`spk01/manifests/sc-04-*.yaml`, `sc-07-*.yaml`) plus, respectively, the `NetworkPolicy` / init-container gate and the `outcap` wrapper in `command`. Runs and evidence were collected by `spk01/scripts/run-sc.ps1`, which writes the same JSONL fields as `measure.RunRecord` (plus `log_bytes`).
- **First measurement round (archived):** SC-03, SC-05 and SC-06 were first measured with `go test -json` in the container (SC-03 without `securityContext`), with the same outcomes. That round is kept in [`legacy-go-test/`](https://github.com/ITA-Internship/GoJudge/tree/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/legacy-go-test); the results in §3 are from the re-run on the hardened specification above.
- **Namespace policy:** from 6 Nov the `spk01` namespace has the `NetworkPolicy` `spk01-deny-egress` (deny all egress for `app: spk01`) applied, so all later runs executed under it.

---

## 3. Success Criteria & Recorded Evidence

Evidence root: [`spk01/results/`](https://github.com/ITA-Internship/GoJudge/tree/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results) at commit [`17939ad`](https://github.com/ITA-Internship/GoJudge/commit/17939ad95e6fddb064bca8b90057124f6d8daab1).

| Criterion | What to observe | Target | Observed Result | Status | Raw Evidence |
| :---: | :--- | :---: | :--- | :---: | :--- |
| **SC-01** | Trusted `go test` execution via client-go | 10 of 10 | **10/10 PASS**. Job `Complete`, `TestsRun=1`, `TestsFailed=0`, verdict `pass`. Startup ~2.7–3.2s, execution ~1.0s. | **Verified** | [`sc-01.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-01.jsonl), `SC-01-run*.log` |
| **SC-02** | Deadline overrun termination | 10 of 10 | **10/10 PASS**. `activeDeadlineSeconds: 5`, `terminationGracePeriodSeconds: 1`. Killed with `DeadlineExceeded`. Overrun: **+432ms to +697ms** (target $\le 5000\text{ms}$). | **Verified** | [`sc-02.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-02.jsonl), `SC-02-run*.log` |
| **SC-03** | Memory exhaustion limit | 5 of 5 | **5/5 PASS** on the hardened spec. Memory limit 512Mi. Container terminated with `reason: OOMKilled`, `exitCode: 137` after ~400 MiB allocated. | **Verified** | [`sc-03.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-03.jsonl), `SC-03-run*.log` |
| **SC-04** | Outbound TCP attempt fails | 5 of 5 | **Failed in the first run with the policy applied** (1 run): connections to `1.1.1.1:443`, `8.8.8.8:53`, `google.com:443` and the API server `10.43.0.1:443` succeeded (first after 19 ms). This meets **F-03**, so the 5-run series was stopped and the mentor informed. Separate experiment (not counted for SC-04): init-container gate, 5/5 runs with egress blocked from t=1 ms. | **Not verified** | [`sc-04.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-04.jsonl), `SC-04-run01.log`, `SC-04-control-run01.log`, `SC-04-diag-*`, `SC-04-est-control-run01.log`, [`sc-04-gate.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-04-gate.jsonl) |
| **SC-05** | Read-only root FS and restricted write paths | 5 of 5 | **5/5 PASS + 1 control**. `/` mounted `ro`. Writes to `/hack.txt`, `/etc`, `/workspace`, `/usr/local/go` failed with `EROFS`, to `/root` with `EACCES`. Allowed path `/tmp` (`emptyDir` 256Mi) writable. | **Verified** | [`sc-05.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-05.jsonl), `SC-05-run*.log`, `baseline-unhardened/SC-05-run01.log` |
| **SC-06** | Non-root, no privilege escalation, no SA token | 5 of 5 | **5/5 PASS + 1 control**. `uid=65532`, `setuid(0)` fails with `EPERM`, no service-account token mounted, `NoNewPrivs=1`, `CapEff=0`. | **Verified** | [`sc-06.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-06.jsonl), `SC-06-run*.log`, `baseline-unhardened/SC-06-run01.log` |
| **SC-07** | Bounded output flood truncation | 3 of 3 | **3/3 PASS + 1 control**. Payload writes 20 MiB to stdout. With `outcap -limit 1048576` the log is 1,049,886–1,049,902 bytes, truncation marker present, test and package verdicts kept. Control without `outcap`: 22,997,587 bytes. | **Verified** | [`sc-07.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-07.jsonl), `SC-07-run*.log`, `SC-07-control-run01.log.zip` |
| **SC-08** | Log and report retrieval before cleanup | 10 of 10 | **10/10 PASS**. Mixed scenarios: 4 pass (`verdict=pass`), 3 assertion fail (`verdict=fail`, exit 1), 3 deadline timeout (`exit 143`). All logs retrieved before Job deletion. | **Verified** | [`sc-08.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-08.jsonl), `SC-08-run*.log` |
| **SC-09** | TTL controller auto-cleanup | 10 of 10 | **10/10 PASS**. `ttlSecondsAfterFinished: 10`. Runner does not delete Job. K8s TTL controller removed Job & Pod in **9006ms–10007ms** (~10s, threshold $\le 70000\text{ms}$). | **Verified** | [`sc-09.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-09.jsonl), `SC-09-run*.log` |
| **SC-10** | Manual cancellation latency | 5 of 5 | **5/5 PASS**. Running Job cancelled via `DeletePropagationForeground`. Complete deletion of Job and Pod from etcd verified in **1010ms–1110ms** (~1.1s, threshold $\le 30000\text{ms}$). | **Verified** | [`sc-10.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-10.jsonl), `SC-10-run*.log` |

---

### Detailed Evidence per Subsystem

#### A. Orchestration & Job Lifecycle (Engineer 1)
- **SC-01 (10/10):** Executed on fully hardened spec (`runAsUser: 65532`, `readOnlyRootFilesystem: true`, capabilities drop `ALL`). Average startup latency: 2680 ms; average execution latency: 1003 ms. Full NDJSON event stream parsed with `TestsRun=1`, `TestsFailed=0`, and `package_verdict="pass"`.
- **SC-02 (10/10):** Executed infinite loop payload with `activeDeadlineSeconds: 5`. Overrun latency beyond the 5s mark ranged between +432 ms and +697 ms. Container terminated with `exitCode: 143` (SIGTERM by kubelet) and Job transitioned to `Failed (DeadlineExceeded)`.
- **SC-08 (10/10):** Verified across 3 distinct outcomes:
  - Runs 01–04 (`TestAddSuccess`): Job `Complete`, `exit_code: 0`, `package_verdict: "pass"`.
  - Runs 05–07 (`TestFailAssertion`): Job `Failed (BackoffLimitExceeded)`, `exit_code: 1`, `package_verdict: "fail"`. Full failure stack trace captured in logs.
  - Runs 08–10 (`TestInfiniteLoop`): Job `Failed (DeadlineExceeded)`, `exit_code: 143`. Partial execution events captured up to termination.
    In all 10 runs, log file existence and size on disk were verified prior to `Jobs().Delete()`.
- **SC-09 (10/10):** Configured `ttlSecondsAfterFinished: 10`. Monitored via API polling without issuing delete commands. Both Job and Pod returned 404 IsNotFound exactly 9.0s–10.0s after completion.
- **SC-10 (5/5):** Running payload cancelled via `DeletePropagationForeground`. Total time from cancel request to 404 IsNotFound on both Job and Pod was 1010ms–1110ms.

#### B. Sandbox Isolation & Boundary Violations (Engineer 2)
Each criterion was measured as a pair: a negative control without the control under test (the payload must detect the violation and fail), then the target configuration (the payload must pass).

- **SC-03 (5/5):** `TestMemoryExhaustion` allocates 10 MiB chunks and writes every byte (Linux allocates physical memory lazily, so untouched memory would never reach the limit). With a 512Mi limit the container was killed after ~400 MiB had been allocated: `terminated.reason: OOMKilled`, `exitCode: 137` in 5/5 runs. The kill is reported for the whole container even though the test runs as a child of `test2json`. The first round with `go test` and no `securityContext` gave the same result (`legacy-go-test/`).
- **SC-04 (Failed, F-03):**
  - *Control without policy:* all four targets reachable (`SC-04-control-run01.log`).
  - *With `NetworkPolicy` `spk01-deny-egress`:* the policy was accepted (`SC-04-diag-netpol-describe.txt`), the k3s network policy controller was running (`SC-04-diag-k3s-netpol.log`) and its `KUBE-NWPLCY` / `KUBE-POD-FW` chains existed on the node (`SC-04-diag-iptables.txt`). Still, 16 minutes after the policy was applied, the test connected to all four targets, the first after 19 ms (`SC-04-run01.log`).
  - *Cause (measured):* the controller installs the per-Pod filter rules only after the container has started. A diagnostic Pod with the same label, probing once per second, got `OPEN` on the first attempt and `BLOCKED` on the next nine (`SC-04-diag-window.log`). The precompiled test binary connects inside this window; with the earlier ~49 s compilation in the Pod the window was hidden.
  - *Established connections survive:* a connection opened at t=24 ms still carried an HTTP request at t=5.1 s, while new connections were already refused (`SC-04-est-control-run01.log`). This follows the `RELATED,ESTABLISHED -j ACCEPT` rule of the controller, so a connection opened in the window stays usable for the life of the Pod.
  - *Init-container gate experiment (allowed by the mentor, not counted for SC-04):* a `busybox` init container probes `1.1.1.1:443` once per second and lets the judge container start only after the probe is blocked; after 30 attempts it exits 1 and the Pod fails (fail-closed). In 4 of 5 new Pods the first probe was still `OPEN` and the second (~1.04 s later) `BLOCKED`; in 1 Pod it was blocked at once. In 5/5 runs the judge container could not connect at t=1 ms, and all four targets were blocked (`sc-04-gate.jsonl`, `SC-04-gate-run*-net-gate.log`). Extra start-up delay: ~1 s. Limitations are listed in §5.
- **SC-05 (5/5 + 1 Control):** `TestFilesystemWrite` checks `/proc/self/mounts`, attempts writes outside the allowed path and one write inside it. Root overlay mounted `ro`. Writes to `/hack.txt`, `/etc/hack.txt`, `/workspace/hack.txt`, `/usr/local/go/hack.txt` failed with `EROFS`, which shows the filesystem itself blocks them, independent of user permissions. The write to `/root/hack.txt` failed with `EACCES`, because the directory permission check happens before the read-only check; the read-only mount is shown by the mount check. The write to the allowed path `/tmp/allowed.txt` succeeded (`emptyDir`, `sizeLimit: 256Mi`). The unhardened control mounted `/` as `rw`, and all five writes succeeded.
- **SC-06 (5/5 + 1 Control):** `TestPrivilegeCheck` checks the UID, attempts `setuid(0)`, looks for the service-account token and reads `/proc/self/status`. Hardened: `uid=65532`, `setuid(0)` fails with `EPERM`, the token at `/var/run/secrets/kubernetes.io/serviceaccount/token` is absent, `NoNewPrivs: 1`, `CapEff: 0000000000000000`. Unhardened control: `uid=0`, `setuid(0)` succeeds, the token is mounted, `NoNewPrivs: 0`, `CapEff: 00000000a80425fb`.
- **SC-07 (3/3 + 1 Control):** Kubernetes does not truncate container output; kubelet only stores it. `TestOutputFlood` writes 20 MiB to stdout in 1 KiB lines. The `outcap` wrapper (`spk01/image/outcap/main.go`, PID 1) runs `test2json` with stdout and stderr on one pipe and reads it event by event: package-level events (`start`, final `pass`/`fail`) are always forwarded; test-level events are limited by `-event-limit` (64 KiB); test output by `-limit` (1 MiB). It emits one truncation marker and a summary line, drops over-long single lines without buffering them, forwards `SIGTERM`, and exits with the child's exit code. The log is therefore bounded by `limit + event-limit + ~2 KiB`; short lines may still pass after the marker if they fit the remaining budget. Control without `outcap`: 22,997,587 bytes (stored compressed as `SC-07-control-run01.log.zip`).

---

## 4. Failure Criteria Check

| Failure Criterion | Description | Spike Finding | Status |
| :---: | :--- | :--- | :---: |
| **F-01** | Job cannot run trusted `go test` with result retrieved (SC-01 fails) | SC-01 achieved 10/10 PASS on hardened specification. | **Not met** |
| **F-02** | Job past its deadline is not terminated (SC-02 fails) | SC-02 achieved 10/10 PASS, terminated within deadline + 697ms. | **Not met** |
| **F-03** | Outbound connection from Job succeeds in any run (SC-04 fails) | **F-03 MET.** With the deny-egress policy applied, connections to external IPs, DNS-resolved hosts and the API server succeeded before the per-Pod rules were installed (window < ~1 s after start; first connection after 19 ms). Reported to the mentor on Day 5 (6 Nov); the mentor confirmed No-go for SC-04. | **MET** |
| **F-04** | Write outside allowed path succeeds in any run (SC-05 fails) | SC-05 achieved 5/5 PASS (and 5/5 in the first round), all writes outside `/tmp` blocked with `EROFS`/`EACCES`. | **Not met** |
| **F-05** | Logs or reports are lost before cleanup in any run (SC-08 fails) | SC-08 achieved 10/10 PASS across all success, error, and timeout scenarios. | **Not met** |

---

## 5. Outcome & Recommendations

### Spike Outcome: No-go
Per the Spike Brief rules:
> *Failure criteria (any one means No-go): F-03: an outbound connection from the Job succeeds in any run (SC-04 fails).*

Although 9 out of 10 criteria were fully **Verified** (including SC-01 and SC-08 required for Conditional Go), **F-03 was triggered**. In the local cluster from INF-03 (k3s v1.35.5 with its embedded network policy controller), `NetworkPolicy` does not isolate a Pod from its first moment: the per-Pod rules appear after the container has started. The mentor confirmed the **No-go** outcome on 6 Nov; the init-container gate is recorded as a candidate for further architectural research, not as a fix.

### Risks & Limitations Identified
1. **NetworkPolicy start-up window:** the k3s network policy controller installs per-Pod filter rules only after the container has started. In our measurements the window was below ~1 s (first diagnostic probe `OPEN`, the next one blocked; 4 of 5 gate probes `OPEN` on the first attempt). Inside it, external IPs, DNS and the API server (`10.43.0.1:443`) are reachable. Blocked connections were refused (`connection refused`), not silently dropped.
2. **Established connections survive:** the controller's `RELATED,ESTABLISHED -j ACCEPT` rule keeps a connection opened in the window usable after the policy is in place (measured at t=5.1 s).
3. **Init-container gate limitations:** (a) it probes a single target and treats "no answer" as "blocked", so if `1.1.1.1` is unreachable for another reason it opens during the window; probing several targets, or first confirming the target is reachable, would reduce this; (b) the window is measured only with ~1 s resolution (upper bound); (c) it assumes the controller installs all rules of a Pod at once; (d) it adds ~1 s start-up delay and counts towards `activeDeadlineSeconds`; (e) it adds the `busybox` image to the supply chain (could be replaced by a Go gate inside the judge image); (f) if rules take longer than 30 probes, Pods fail to start (safe, but reduces availability).
4. **Output limit bypass (not tested):** `outcap` and the test code run as the same user (`65532`), so a test could open `/proc/1/fd/1` and write to the container's stdout directly, bypassing the wrapper. Possible mitigations: `prctl(PR_SET_DUMPABLE, 0)` in the wrapper, and `PodLogOptions.LimitBytes` when the runner reads logs as a second layer. The spike ran only trusted code, so this is a risk for SPEC-F07.
5. **Writable `/tmp`:** the allowed path is bounded by `sizeLimit: 256Mi` and lives only as long as the Pod; code can fill it up to the limit, after which kubelet evicts the Pod (not immediately).
6. **Scope of measurements:** all runs were on one local k3d cluster (Docker Desktop on Windows, one node). Behaviour of other network plugins (Calico, Cilium) and of a multi-node cluster was not measured.

### Follow-On Decisions & Next Steps
- **DEC-01 (Target Platform: Go, Kubernetes, PostgreSQL, RabbitMQ):** Retained, but Kubernetes networking requires replacement or hardening before multi-tenant execution is permitted.
- **DEC-02 (Kubernetes Job per submission execution model):** **Not approved in its current form.** Must not proceed directly to feature implementation (SPEC-F07 / F07-01) without addressing network gating.
- **DEC-08 (Execution Isolation Mechanism):** The mentor and team must select an isolation architecture. Options identified (only Option 1 was measured in this spike; changing the cluster network plugin requires agreement with Pair A):
  - **Option 1 (Init-Container Network Gate):** a mandatory gating init-container (`spk01/manifests/sc-04-initgate.yaml`) that waits until egress is blocked before the judge container starts. Measured 5/5, with the limitations in §5 item 3.
  - **Option 2 (Different network plugin, e.g. Calico or Cilium):** a plugin that applies policy before the container process starts. Not evaluated.
  - **Option 3 (Sandboxed runtime, e.g. gVisor `runsc`):** stronger syscall isolation; whether it closes the network start-up window was not evaluated.
- **SC-07 output limiting:** if adopted, `outcap` should wrap the command in `jobspec/builder.go`, and the runner should read logs with `LimitBytes` as a second layer (§5 item 4).