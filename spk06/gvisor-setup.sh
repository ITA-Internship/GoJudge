#!/usr/bin/env bash
# Adds gVisor (runsc, network=none) to the SPK-06 k3d cluster
# Needs runsc and containerd-shim-runsc-v1 from https://gvisor.dev/docs/user_guide/install/
# Usage: GVISOR_DIR=~/gvisor ./spk06/gvisor-setup.sh
set -euo pipefail
export MSYS_NO_PATHCONV=1 # Git Bash on Windows must not rewrite /usr/... paths

GVISOR_DIR=$(cygpath -m "${GVISOR_DIR:-$HOME/gvisor}") # docker.exe needs a Windows-style path
CONTEXT=k3d-spk06
NODE=k3d-spk06-server-0
CFG=/var/lib/rancher/k3s/agent/etc/containerd

docker exec "$NODE" mkdir -p /usr/local/bin /etc/containerd "$CFG"
docker cp "$GVISOR_DIR/runsc" "$NODE:/usr/local/bin/runsc"
docker cp "$GVISOR_DIR/containerd-shim-runsc-v1" "$NODE:/usr/local/bin/containerd-shim-runsc-v1"
docker cp "$GVISOR_DIR/gvisor-bin" "$NODE:/usr/local/bin/gvisor-bin" # sidecars, runsc looks for them next to itself
docker exec "$NODE" sh -c "chmod -R +x /usr/local/bin"
docker cp spk06/gvisor/runsc-nonet.toml "$NODE:/etc/containerd/runsc-nonet.toml"
docker cp spk06/gvisor/config-v3.toml.tmpl "$NODE:$CFG/config-v3.toml.tmpl"

# k3s regenerates the containerd config on start
docker restart "$NODE"
until kubectl --context "$CONTEXT" get nodes 2>/dev/null | grep -q ' Ready'; do sleep 3; done
kubectl --context "$CONTEXT" apply -f spk06/gvisor/runtimeclass.yaml
docker exec "$NODE" runsc --version
