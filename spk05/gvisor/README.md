# SPK-05 gVisor experiments

Standalone experiments with `runsc do`, run before any Kubernetes setup. They answer two questions:
does the SPK-01 Go test binary run under gVisor, and what does `--network=none` do to the egress test.
Not a Kubernetes run: no `runtimeClassName`, containerd, Job or NetworkPolicy.
Trusted tests only.

## Environment
See `results/env.txt`: WSL2 kernel `6.6.87.2-microsoft-standard-WSL2`, Ubuntu 24.04.5 LTS, `runsc release-20261116.0`.

## Setup
1. Install `runsc` from the gVisor apt repository (https://gvisor.dev/docs/user_guide/install/), then `sudo apt-get install -y iptables` (needed by `runsc do` with default network).
2. Build the image and extract the test binary (Windows `cmd`, from the repo root):
   ```
   docker build -t gojudge-spk01:dev spk01\image
   mkdir out
   docker run --rm -v "%cd%\out":/out --entrypoint cp gojudge-spk01:dev /usr/local/bin/trusted.test /out/
   ```
3. Copy `out\trusted.test` into the WSL home folder (`~/spk05/`) and `chmod +x` it. The binary is not committed.

## Commands (run in `~/spk05`)
```
# E1 compatibility
./trusted.test -test.v -test.run '^TestAddSuccess$'                                  # results/native-addsuccess.log
sudo runsc do $PWD/trusted.test -test.v -test.run '^TestAddSuccess$'                 # results/runsc-addsuccess.log
sudo runsc --network=none do $PWD/trusted.test -test.v -test.run '^TestAddSuccess$'  # results/runsc-nonet-addsuccess.log

# E2 egress without network (initial run + 5 runs)
sudo runsc --network=none do $PWD/trusted.test -test.v -test.run '^TestNetworkEgress$'   # results/runsc-nonet-egress.log, runsc-nonet-egress-run1..5.log

# E3 control, default network
sudo runsc do $PWD/trusted.test -test.v -test.run '^TestNetworkEgress$'              # results/runsc-net-egress.log
```
Each command was piped through `2>&1 | tee results/<file>`.

## Results
| ID | Result |
|---|---|
| E1 | `TestAddSuccess` PASS natively and under `runsc` (both network modes) |
| E2 | `TestNetworkEgress` PASS in 6/6 runs: all targets `network is unreachable` (0-3 ms) |
| E3 | `TestNetworkEgress` FAIL (expected): all 3 targets connected (11, 26, 46 ms) |

## Limitations
- `runsc do` only, not the Kubernetes path. Per-handler `network=none` in containerd is not tested.
- The API server target is not probed (its variable exists only inside a Pod).
- One host, no `runc` network control.