#!/usr/bin/env bash
# k3d cluster with Cilium instead of k3s NetworkPolicy
set -euo pipefail

NODE=k3d-spk05-cilium-server-0
CONTEXT=k3d-spk05-cilium
CILIUM_VERSION=1.20.5

# no flannel, no k3s network policy
# no metrics-server: it hangs namespace deletion before Cilium is up
k3d cluster create spk05-cilium \
  --k3s-arg "--flannel-backend=none@server:*" \
  --k3s-arg "--disable-network-policy@server:*" \
  --k3s-arg "--disable=traefik@server:*" \
  --k3s-arg "--disable=metrics-server@server:*" \
  --no-lb

# Cilium needs bpffs and cgroup2
# Run again after a node restart
docker exec "$NODE" sh -c '
  mount bpffs /sys/fs/bpf -t bpf &&
  mount --make-shared /sys/fs/bpf &&
  mkdir -p /run/cilium/cgroupv2 &&
  mount -t cgroup2 none /run/cilium/cgroupv2 &&
  mount --make-shared /run/cilium/cgroupv2
'

# 0.0.0.0:<port> from kubeconfig does not work inside the cluster
NODE_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$NODE")

# keep kube-proxy: only the NetworkPolicy engine differs from SPK-01
cilium install --context "$CONTEXT" --version "$CILIUM_VERSION" \
  --set kubeProxyReplacement=false \
  --set k8sServiceHost="$NODE_IP" \
  --set k8sServicePort=6443

cilium status --context "$CONTEXT" --wait
kubectl --context "$CONTEXT" -n kube-system wait --for=condition=Ready pod --all --timeout=300s
kubectl --context "$CONTEXT" get nodes

# SPK-01 image (imagePullPolicy: Never)
docker build -t gojudge-spk01:dev spk01/image
k3d image import gojudge-spk01:dev -c spk05-cilium

# same namespace and deny-egress policy as SPK-01
kubectl --context "$CONTEXT" create namespace spk01
kubectl --context "$CONTEXT" apply -f spk01/manifests/sc-04-deny-egress.yaml
