#!/usr/bin/env bash
# k3d cluster for the SPK-06 suite (k3d as in ADR-002)
# POD_MAX_PIDS=128 ./spk06/setup.sh sets a kubelet pids limit (our experimental value)
set -euo pipefail

CLUSTER=spk06
CONTEXT=k3d-$CLUSTER

# same k3s version as the project cluster (INF-03)
K3S_IMAGE=$(awk '/^image:/ {print $2}' infra/k3d/k3d-config.yaml)

args=(--image "$K3S_IMAGE" --k3s-arg "--disable=traefik@server:*" --no-lb)
if [[ -n ${POD_MAX_PIDS:-} ]]; then
  args+=(--k3s-arg "--kubelet-arg=pod-max-pids=$POD_MAX_PIDS@server:*")
fi
k3d cluster create "$CLUSTER" "${args[@]}"
# system pods appear a few seconds after the cluster is created
until kubectl --context "$CONTEXT" -n kube-system get pods -o name | grep -q .; do sleep 2; done
kubectl --context "$CONTEXT" -n kube-system wait --for=condition=Ready pod --all --timeout=300s

docker build -t gojudge-spk06:dev spk06/image
k3d image import gojudge-spk06:dev -c "$CLUSTER"

kubectl --context "$CONTEXT" create namespace spk06
kubectl --context "$CONTEXT" apply -f spk06/manifests/deny-egress.yaml
