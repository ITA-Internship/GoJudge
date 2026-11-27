# SPK-06: isolation check suite

7 isolation checks with trusted tests. The suite does not depend on the isolation mechanism:
the runtime is set with `RUNTIME_CLASS`, the cluster with `CONTEXT`

- main run: gVisor (`RUNTIME_CLASS=gvisor-nonet`), recommended by SPK-05
- baseline: plain k3d, no runtime class (results below)

## Run

From the repo root. Needs Docker, k3d, kubectl, jq, perl

```bash
POD_MAX_PIDS=128 ./spk06/setup.sh # cluster, image, namespace, deny-egress policy
./spk06/run-suite.sh # all checks, or one: ./spk06/run-suite.sh cpu
./spk06/cleanup.sh # delete cluster
```

Other cluster: `CONTEXT=<context> RUNTIME_CLASS=<name> RESULTS=<dir> ./spk06/run-suite.sh`

gVisor run. Needs `runsc`, `containerd-shim-runsc-v1` and the `gvisor-bin` folder from the
[gVisor release archive](https://gvisor.dev/docs/user_guide/install/) in `~/gvisor` (or `GVISOR_DIR`)

```bash
./spk06/gvisor-setup.sh # copies gVisor into the k3d node, handler runsc-nonet (network=none), RuntimeClass gvisor-nonet
RUNTIME_CLASS=gvisor-nonet RESULTS=spk06/results/gVisor ./spk06/run-suite.sh
```

`gvisor-setup.sh` has to be repeated after every new `setup.sh`: the files live inside the node.
On Windows run the scripts in Git Bash. If `kubectl` times out on `host.docker.internal`, use
`kubectl config set-cluster k3d-spk06 --server=https://127.0.0.1:<port>`

## Experimental parameters

Not approved production limits. Taken from SPK-01, except the pids limit

| Parameter | Value | Source |
|---|---|---|
| CPU | limit 1, request 250m | SPK-01 |
| Memory | 512Mi | SPK-01 |
| Time | `activeDeadlineSeconds: 5` | SPK-01 |
| Output | 1 MiB (`outcap`) | SPK-01 |
| `/tmp` | `emptyDir` 256Mi | SPK-01 |
| Processes | 128 per pod (kubelet `pod-max-pids`) | our choice: safe for one node, enough for a Go test |

## Baseline result (k3d)

k3d, k3s v1.35.5. Logs are in `results/k8s (baseline)/`

PASS/FAIL only where the requirement is yes/no. For CPU, output and processes
there is no approved threshold, so only the behaviour is recorded

| Check | Expected | Result |
|---|---|---|
| network | all 4 targets blocked | FAIL 0/5: 1.1.1.1 reached in 38-90 ms (F-03) |
| memory | `OOMKilled` | PASS 3/3: killed after ~500 MiB |
| timeout | stopped within 5 + 5 s | PASS 3/3: 6 s |
| filesystem | only `/tmp` writable | PASS 3/3 |
| cpu | behaviour only | 3/3: used 1.00-1.01 cores of 10, limit 1 |
| output | behaviour only | 3/3: log 1,048,888-1,049,051 bytes, verdict kept |
| process | behaviour only | 3/3: stopped at 119 processes (pod limit 128) |

Without `POD_MAX_PIDS` nothing limits processes: the test starts all 512

## gVisor result (k3d)

Same cluster and k3s v1.35.5, `runsc` release-20261123.1, `RUNTIME_CLASS=gvisor-nonet`
(`network=none`). A Pod under gVisor reports kernel `4.19.0-gvisor`. Logs are in `results/gVisor/`

| Check | Expected | Result |
|---|---|---|
| network | all 4 targets blocked | PASS 5/5: `network is unreachable` from 0 s |
| memory | `OOMKilled` | PASS 3/3, but killed only after 2500 MiB (limit 512Mi) |
| timeout | stopped within 5 + 5 s | PASS 3/3: 6 s |
| filesystem | only `/tmp` writable | PASS 3/3 |
| cpu | behaviour only | 3/3: 1.44-1.51 cores measured inside the sandbox, limit 1 (see below) |
| output | behaviour only | 3/3: log 1,048,935-1,048,955 bytes, verdict kept |
| process | behaviour only | 3/3: stopped at 41-42 processes, `cannot allocate memory` |

### Baseline vs gVisor

| Check | Baseline (runc) | gVisor |
|---|---|---|
| network | FAIL 0/5 (F-03) | PASS 5/5 |
| memory | killed after ~500 MiB | killed after 2500 MiB |
| timeout | 6 s | 6 s |
| filesystem | PASS 3/3 | PASS 3/3 |
| cpu | 1.01 cores (10 CPUs visible) | 1.44-1.51 cores (2 CPUs visible), 1.0 on the host |
| output | 1,048,888-1,049,051 bytes | 1,048,935-1,048,955 bytes |
| process | 119 processes, `resource temporarily unavailable` | 41-42 processes, `cannot allocate memory` |

### Notes

- cpu and process tests read `/sys/fs/cgroup/cpu.max` and `pids.max`. Under gVisor these files do not exist
  (the Pod sees the Sentry kernel, not the host one). The first gVisor run failed at 0 s for that reason.
  The tests now log the missing file and measure the behaviour anyway (`cpu_test.go`, `procs_test.go`).
  Baseline behaviour is the same: the files exist there. `results/gVisor/` holds the second run
- cpu: the measurement inside the sandbox is too high. One probe Pod with the same limits, cgroup of the Pod read
  on the node: `cpu.max` `100000 100000`, ~5.15 s of CPU in ~5 s, throttled in 50 of 63 periods. The limit works on
  the host, the in-sandbox `getrusage` over-reports (1.47 cores in that run). Single run via `kubectl exec`, not a Job.
  Raw data: `./spk06/host-probe.sh cpu`, `results/gVisor/host-probes`
- memory: the test reached 2500 MiB inside the sandbox in 3/3 runs, 5 times the limit. On the host the cgroup of
  the Pod stayed at its limit: `memory.max` 536870912, `memory.peak` 512 MiB, `memory.events` `oom_kill 1`.
  So the limit works, the in-sandbox number does not reflect host memory. Why the test could allocate 2500 MiB
  before the kill is not established. One probe Pod via `kubectl exec`, not a Job. Raw data:
  `./spk06/host-probe.sh memory`, `results/gVisor/host-probes`. The runner counts any `OOMKilled` as
  PASS and does not look at the peak
- process: the pids limit of 128 was not reached. With memory 2Gi (`results/gVisor-process-2Gi/`) the result is the
  same, 41-42 processes, so the stop is not the container memory limit. Whether the pids limit works under gVisor
  is not verified

## Limitations

- one host (Docker Desktop), one node, one k3s version, 3 runs per check (5 for network)
- the runtime is chosen per Pod. Nothing here stops a Pod from using another `runtimeClassName`
- with `network=none` the Pod has no DNS and no API server, results leave only through stdout and the exit code
- Cilium + gVisor is not covered
- thresholds are experimental parameters, not approved requirements
