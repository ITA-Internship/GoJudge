# SPK-05 gVisor experiments

Two stages. Stage 1 is standalone `runsc do`, run before any Kubernetes setup. Stage 2 runs the same
SPK-01 egress test in k3s with gVisor as a RuntimeClass. They answer: does the SPK-01 Go test binary
run under gVisor, what does `--network=none` do to the egress test, and does it hold in Kubernetes.
Trusted tests only.

## Environment
- Stage 1: `wsl_results/env.txt`: WSL2 kernel `6.6.87.2-microsoft-standard-WSL2`, Ubuntu 24.04.5 LTS, `runsc release-20261116.0`.
- Stage 2: same WSL2 host, k3s `v1.37.2+k3s1` (single node), containerd `2.3.7-k3s1`. SPK-01 used k3s `v1.35.5`, so versions differ. See `k8s_results/env-k3s.txt`, `k8s_results/nodes.txt`.

## Stage 1: standalone `runsc do`

### Setup
1. Install `runsc` from the gVisor apt repository (https://gvisor.dev/docs/user_guide/install/), then `sudo apt-get install -y iptables` (needed by `runsc do` with default network).
2. Build the image and extract the test binary (Windows `cmd`, from the repo root):
   ```
   docker build -t gojudge-spk01:dev spk01\image
   mkdir out
   docker run --rm -v "%cd%\out":/out --entrypoint cp gojudge-spk01:dev /usr/local/bin/trusted.test /out/
   ```
3. Copy `out\trusted.test` into the WSL home folder (`~/spk05/`) and `chmod +x` it. The binary is not committed.

### Commands (run in `~/spk05`)
```
# E1 compatibility
./trusted.test -test.v -test.run '^TestAddSuccess$'                                  # native-addsuccess.log
sudo runsc do $PWD/trusted.test -test.v -test.run '^TestAddSuccess$'                 # runsc-addsuccess.log
sudo runsc --network=none do $PWD/trusted.test -test.v -test.run '^TestAddSuccess$'  # runsc-nonet-addsuccess.log

# E2 egress without network (initial run + 5 runs)
sudo runsc --network=none do $PWD/trusted.test -test.v -test.run '^TestNetworkEgress$'   # runsc-nonet-egress.log, runsc-nonet-egress-run1..5.log

# E3 control, default network
sudo runsc do $PWD/trusted.test -test.v -test.run '^TestNetworkEgress$'              # runsc-net-egress.log
```
Each command was piped through `2>&1 | tee results/<file>`; the logs were then copied to `wsl_results/`.

## Stage 2: k3s with gVisor RuntimeClass (WSL2)

### Setup
1. Install k3s in WSL2 with systemd enabled (https://docs.k3s.io/quick-start). `runsc` and `containerd-shim-runsc-v1` must be in `PATH`.
2. Register the handlers in the k3s containerd template `/var/lib/rancher/k3s/agent/etc/containerd/config-v3.toml.tmpl` (copy: `k8s_results/containerd-config-v3.toml.tmpl`):
   ```
   {{ template "base" . }}

   [plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.'runsc']
     runtime_type = "io.containerd.runsc.v1"

   [plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.'runsc-nonet']
     runtime_type = "io.containerd.runsc.v1"

   [plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.'runsc-nonet'.options]
     TypeUrl = "io.containerd.runsc.v1.options"
     ConfigPath = "/etc/containerd/runsc-nonet.toml"
   ```
   `/etc/containerd/runsc-nonet.toml` (copy: `k8s_results/runsc-nonet.toml`):
   ```
   [runsc_config]
     network = "none"
   ```
   Then `sudo systemctl restart k3s`. k3s does not detect `runsc` and does not create a RuntimeClass by itself.
3. Create the RuntimeClasses (`k8s_results/runtimeclasses.yaml`): `gvisor` (handler `runsc`, default network) and `gvisor-nonet` (handler `runsc-nonet`, `network=none`).
4. Make the SPK-01 image visible to k3s containerd (it has its own image store). From Windows `cmd`: `docker save gojudge-spk01:dev -o gojudge-spk01.tar`; in WSL: `sudo k3s ctr images import gojudge-spk01.tar`. The `.tar` is not committed. The manifest uses `imagePullPolicy: Never`.
5. `sudo k3s kubectl create namespace spk01`.
6. Manifests: `spk01/manifests/sc-04-netdial.yaml` (runc control), `sc-04-deny-egress.yaml` (NetworkPolicy). Here: `sc-04-netdial-gvisor.yaml` and `sc-04-netdial-gvisor-nonet.yaml`, copies of the first with one added line `runtimeClassName: gvisor` / `gvisor-nonet` after `restartPolicy: Never`. If the repo has CRLF line endings, strip them first (`tr -d '\r' < in | sed ... > out`).

### Commands (repo root, WSL)
`run-k8s.sh <label> <manifest> [runs]` creates the Job `runs` times (default 5), waits for it, saves the Pod log, exit code and (for run 1) the Pod YAML to `k8s_results/`, and appends one line per run to `k8s_results/summary.txt`.
```
# E6 smoke test: Pod under gVisor vs plain Pod
sudo k3s kubectl exec nginx-gvisor -- uname -r     # 4.19.0-gvisor
sudo k3s kubectl exec nginx-plain  -- uname -r     # 6.6.87.2-microsoft-standard-WSL2

# E7, E8: no NetworkPolicy
./run-k8s.sh ctrl-runc       spk01/manifests/sc-04-netdial.yaml
./run-k8s.sh gvisor-nopolicy sc-04-netdial-gvisor.yaml

# E9, E10: with NetworkPolicy
sudo k3s kubectl apply -f spk01/manifests/sc-04-deny-egress.yaml
./run-k8s.sh gvisor-policy     sc-04-netdial-gvisor.yaml
./run-k8s.sh ctrl-runc-policy  spk01/manifests/sc-04-netdial.yaml

# E11: is the policy enforced after the Pod is running? (see pol-check.txt)

# E12: network=none handler, NetworkPolicy removed
sudo k3s kubectl delete -f spk01/manifests/sc-04-deny-egress.yaml
./run-k8s.sh gvisor-nonet sc-04-netdial-gvisor-nonet.yaml
```
(Paths of the gVisor manifests and `run-k8s.sh` are under `spk05/gvisor/`.)

## Results
| ID | Result |
|---|---|
| E1 | `TestAddSuccess` PASS natively and under `runsc` (both network modes) |
| E2 | `runsc --network=none do`: `TestNetworkEgress` PASS in 6/6 runs: all targets `network is unreachable` (0-3 ms) |
| E3 | `runsc do`, default network: `TestNetworkEgress` FAIL (expected): all 3 targets connected (11, 26, 46 ms) |
| E6 | Pod with `runtimeClassName: gvisor` is `Running`; `uname -r` is `4.19.0-gvisor`, plain Pod shows the host kernel |
| E7 | runc, no policy: FAIL 5/5 (egress open), so the test sees the open network |
| E8 | gVisor (default network), no policy: FAIL 5/5, gVisor alone does not block egress |
| E9 | gVisor (default network) + NetworkPolicy: FAIL 5/5. Run 1: 1.1.1.1:443 connected in 12 ms, 8.8.8.8:53 in 26 ms, google.com:443 in 35 ms; API server blocked (`connection refused`) |
| E10 | runc + NetworkPolicy: FAIL 5/5, the start-up window of SPK-01 F-03 is reproduced on this k3s |
| E11 | The NetworkPolicy is enforced after start: a running runc Pod labelled `app=spk01` gets `Connection refused` from `wget http://1.1.1.1` (exit 1) |
| E12 | gVisor with `network=none` (`gvisor-nonet`), no policy: PASS 5/5 (`succeeded=1`, exit 0). Run 1: 1.1.1.1:443, 8.8.8.8:53, google.com:443 (DNS server `10.43.0.10:53` unreachable) and API server `10.43.0.1:443` all `network is unreachable`, first probe at 0 s |

Per-run lines: `k8s_results/summary.txt`.

## Known issue (WSL2 + Docker Desktop)
k3s kubelet fails to start with `system validation failed - wrong number of fields (expected 6, got 7)`. Cause: Docker Desktop mounts `C:\Program Files\Docker\Docker\resources` at `/Docker/host`, and the unescaped space in the mount options adds a 7th field to the line in `/proc/mounts`. Check with `awk 'NF!=6 {print NR": "$0}' /proc/mounts`. Workaround: `sudo systemctl stop k3s`, `sudo k3s-killall.sh`, `sudo umount /Docker/host`, then start k3s again. The mount can return after `wsl --shutdown` or a Docker Desktop restart.

## Limitations
- One host (WSL2), one node, one workload. k3s `v1.37.2` differs from SPK-01's `v1.35.5`.
- The test binary has 4 targets and does not write files, so the other criteria (read-only root, memory, output flood) are not tested under gVisor.
- `network=none` is verified only through a dedicated handler (`runsc-nonet`). A Pod that selects a different `runtimeClassName`, or none, gets a normal network. Nothing here enforces the choice.
- With `network=none` the Pod cannot reach DNS or the API server. Results must leave through stdout and the exit code (that is how the Job and `kubectl logs` worked here).
- `pol-check.txt` was captured separately from the loop runs (see the file for the exact commands).