# SPK-01: Isolated Go Test Execution in Kubernetes Jobs (Spike Report)

## Metadata
- **Spike ID:** SPK-01
- **Title:** Isolated Go Test Execution in Kubernetes Jobs
- **Type:** Technical spike
- **Owner:** Pair B (Engineer 1 & Engineer 2)
- **Timebox:** 8 pair-days (Sprint 1, Days 1–10: Mon 2 Nov – Fri 13 Nov 2026, 15:00)
- **Status:** Completed within timebox
- **Dependencies:** DEC-01 (Target platform, provisional), DEC-02 (K8s Job per submission, provisional), DEC-08 (Isolation mechanism), INF-02 (Local cluster infrastructure), MA-03
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

---

## 3. Success Criteria & Recorded Evidence

Evidence root: [`spk01/results/`](https://github.com/ITA-Internship/GoJudge/tree/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results) at commit [`17939ad`](https://github.com/ITA-Internship/GoJudge/commit/17939ad95e6fddb064bca8b90057124f6d8daab1).

| Criterion | What to observe | Target | Observed Result | Status | Raw Evidence |
| :---: | :--- | :---: | :--- | :---: | :--- |
| **SC-01** | Trusted `go test` execution via client-go | 10 of 10 | **10/10 PASS**. Job `Complete`, `TestsRun=1`, `TestsFailed=0`, verdict `pass`. Startup ~2.7–3.2s, execution ~1.0s. | **Verified** | [`sc-01.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-01.jsonl), `SC-01-run*.log` |
| **SC-02** | Deadline overrun termination | 10 of 10 | **10/10 PASS**. `activeDeadlineSeconds: 5`, `terminationGracePeriodSeconds: 1`. Killed with `DeadlineExceeded`. Overrun: **+432ms to +697ms** (target $\le 5000\text{ms}$). | **Verified** | [`sc-02.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-02.jsonl), `SC-02-run*.log` |
| **SC-03** | Memory exhaustion limit | 5 of 5 | **5/5 PASS**. Memory limit 512Mi. Container terminated with `reason: OOMKilled`, `exitCode: 137` at ~400MiB. | **Verified** | [`sc-03.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-03.jsonl), `SC-03-run*.log` |
| **SC-04** | Outbound TCP attempt fails | 5 of 5 | **0/5 PASS (FAILED)**. `NetworkPolicy` applied asynchronously. Outbound TCP connection succeeded in first 19ms (`OPEN`). Triggers **F-03**. Gating init-container experiment achieved 5/5 BLOCKED. | **Not verified** | [`sc-04.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-04.jsonl), [`sc-04-gate.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-04-gate.jsonl) |
| **SC-05** | Read-only root FS and restricted write paths | 5 of 5 | **5/5 PASS + 1 control**. `/` mounted `ro` (`EROFS`), writes to `/root`, `/etc`, `/workspace` blocked. Ephemeral `/tmp` (`emptyDir` 256Mi) writable. | **Verified** | [`sc-05.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-05.jsonl), `SC-05-run*.log` |
| **SC-06** | Non-root, no privilege escalation, no SA token | 5 of 5 | **5/5 PASS + 1 control**. `uid=65532`, `setuid(0)` returns `EPERM`, no service-account token mounted, `NoNewPrivs=1`, `CapEff=0`. | **Verified** | [`sc-06.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-06.jsonl), `SC-06-run*.log` |
| **SC-07** | Bounded output flood truncation | 3 of 3 | **3/3 PASS + 1 control**. Flood capped at 1 MiB via `outcap` wrapper while preserving final test verdict NDJSON event. | **Verified** | [`sc-07.jsonl`](https://github.com/ITA-Internship/GoJudge/blob/17939ad95e6fddb064bca8b90057124f6d8daab1/spk01/results/sc-07.jsonl), `SC-07-run*.log` |
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
- **SC-03 (5/5):** `TestMemoryExhaustion` rapidly dirty-allocated 10 MiB buffers. The Linux cgroup killer terminated the container at ~400 MiB under a 512 MiB limit (`terminated.reason: OOMKilled`, `exitCode: 137`).
- **SC-04 (0/5, Failed):** Investigated NetworkPolicy timing. The k3s NetworkPolicy controller programs `iptables` asynchronously upon Pod creation. Fast-starting static binaries execute `net.Dial("tcp", "1.1.1.1:443")` within 19ms, escaping to external networks before `DROP` rules are injected. Established TCP sockets survive subsequent policy application. An init-container network gate (`sc-04-initgate.yaml`) was prototyped, proving that polling until external egress is blocked successfully protects the subsequent application container (5/5 BLOCKED in `sc-04-gate.jsonl`).
- **SC-05 (5/5 + 1 Control):** Verified filesystem isolation. Root overlay mounted `ro`. Writes to `/hack.txt`, `/etc/hack.txt`, `/workspace/hack.txt` failed with `EROFS` (Read-only file system). Write to `/root/` failed with `EACCES`. Ephemeral write to `/tmp/allowed.txt` succeeded under `emptyDir: sizeLimit: 256Mi`. Control run on unhardened spec verified that writes succeed without this configuration.
- **SC-06 (5/5 + 1 Control):** Verified credential and privilege boundaries. Execution under `uid: 65532`, `setuid(0)` blocked (`EPERM`), service-account token absent from `/var/run/secrets/kubernetes.io/serviceaccount/token`, `NoNewPrivs: 1`, and effective capabilities `CapEff: 0000000000000000`. Control run under root confirmed escalation, token access, and full capabilities.
- **SC-07 (3/3 + 1 Control):** Addressed output flooding. Unbounded test output (50+ MiB) threatens node storage and log collection. Implemented `outcap` wrapper (`spk01/image/outcap/main.go`), which caps raw stdout at 1 MiB while preserving structured NDJSON test verdict events and propagating process exit codes. Control run without `outcap` generated 100+ MiB log spam.

---

## 4. Failure Criteria Check

| Failure Criterion | Description | Spike Finding | Status |
| :---: | :--- | :--- | :---: |
| **F-01** | Job cannot run trusted `go test` with result retrieved (SC-01 fails) | SC-01 achieved 10/10 PASS on hardened specification. | **Not met** |
| **F-02** | Job past its deadline is not terminated (SC-02 fails) | SC-02 achieved 10/10 PASS, terminated within deadline + 697ms. | **Not met** |
| **F-03** | Outbound connection from Job succeeds in any run (SC-04 fails) | **F-03 MET.** Initial TCP connection to 1.1.1.1:443 succeeded during the 19ms NetworkPolicy programming window. Reported to mentor on Day 5 (6 Nov). | **MET** |
| **F-04** | Write outside allowed path succeeds in any run (SC-05 fails) | SC-05 achieved 5/5 PASS, all unauthorized writes blocked with `EROFS`/`EACCES`. | **Not met** |
| **F-05** | Logs or reports are lost before cleanup in any run (SC-08 fails) | SC-08 achieved 10/10 PASS across all success, error, and timeout scenarios. | **Not met** |

---

## 5. Outcome & Recommendations

### Spike Outcome: No-go
Per the Spike Brief rules:
> *Failure criteria (any one means No-go): F-03: an outbound connection from the Job succeeds in any run (SC-04 fails).*

Although 9 out of 10 criteria were fully **Verified** (including SC-01 and SC-08 required for Conditional Go), **F-03 was unequivocally triggered**. Standard Kubernetes NetworkPolicy in k3s cannot provide synchronous, race-free network isolation for fast-starting containers. The mentor formally confirmed the **No-go** outcome.

### Risks & Limitations Identified
1. **NetworkPolicy Asynchronous Race Condition:** Standard CNI controller (Flannel + Kube-Router in k3s) programs packet filters asynchronously after Pod network namespace creation. Any precompiled binary executing in <100ms has an open network window to external IP addresses and internal cluster services (`10.43.0.1:443`).
2. **Established Connection Survival:** Conntrack state tracking permits established TCP sessions initiated during the startup window to remain open even after `DROP` rules are programmed.
3. **Container-Level Output Truncation Bypass:** The `outcap` wrapper safely bounds stdout from cooperative binaries. However, an adversarial payload writing directly to `/proc/1/fd/1` or manipulating file descriptors could bypass application-level wrappers unless enforced by container runtime limits.

### Follow-On Decisions & Next Steps
- **DEC-01 (Target Platform: Go, Kubernetes, PostgreSQL, RabbitMQ):** Retained, but Kubernetes networking requires replacement or hardening before multi-tenant execution is permitted.
- **DEC-02 (Kubernetes Job per submission execution model):** **Not approved in its current form.** Must not proceed directly to feature implementation (SPEC-F07 / F07-01) without addressing network gating.
- **DEC-08 (Execution Isolation Mechanism):** The mentor and team must select an isolation architecture from the following options:
  - **Option 1 (Init-Container Network Gate):** Deploy a mandatory gating init-container (`spk01/manifests/sc-04-initgate.yaml`) that verifies external egress is dropped before yielding to the judge container.
  - **Option 2 (eBPF-Based CNI / Cilium):** Replace k3s default CNI with Cilium to enforce socket-level BPF policies synchronously before the container process is executed.
  - **Option 3 (Sandboxed Runtime Class / gVisor):** Enforce network and syscall boundaries at the hypervisor/kernel emulation layer via `runsc` (gVisor).
