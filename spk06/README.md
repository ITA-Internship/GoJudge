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

k3d, k3s v1.35.5. Logs are in `results/`

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
