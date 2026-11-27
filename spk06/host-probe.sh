#!/usr/bin/env bash
# Host-side cgroup probe: what the node sees while a check runs in a Pod
# Under gVisor the numbers measured inside the sandbox differ from the host ones (README, gVisor result)
# Usage: ./spk06/host-probe.sh cpu|memory   (RUNTIME_CLASS=gvisor-nonet by default, RESULTS=<dir>)
set -euo pipefail
export MSYS_NO_PATHCONV=1 # Git Bash must not rewrite /usr/... paths

AREA=${1:?usage: host-probe.sh cpu|memory}
RUNTIME_CLASS=${RUNTIME_CLASS:-gvisor-nonet}
OUT=${RESULTS:-spk06/results/gVisor/host-probes}
CONTEXT=k3d-spk06
NODE=k3d-spk06-server-0
POD=probe-$AREA

case $AREA in
  cpu)    TEST='^TestCPUBurn$' ;;
  memory) TEST='^TestMemoryExhaustion$' ;;
  *) echo "unknown area: $AREA" >&2; exit 1 ;;
esac

k() { kubectl --context "$CONTEXT" -n spk06 "$@"; }
mkdir -p "$OUT"

# same limits as the suite (cpu 1, memory 512Mi); the Pod only sleeps, the test is started with exec
# so the cgroup of the Pod still exists while it is sampled
k delete pod "$POD" --now --ignore-not-found > /dev/null
k run "$POD" --restart=Never --image=gojudge-spk06:dev --image-pull-policy=Never \
  --overrides="{\"spec\":{\"runtimeClassName\":\"$RUNTIME_CLASS\",\"containers\":[{\"name\":\"$POD\",\"image\":\"gojudge-spk06:dev\",\"imagePullPolicy\":\"Never\",\"command\":[\"sleep\",\"120\"],\"resources\":{\"requests\":{\"cpu\":\"250m\",\"memory\":\"512Mi\"},\"limits\":{\"cpu\":\"1\",\"memory\":\"512Mi\"}}}]}}" > /dev/null
k wait --for=condition=Ready "pod/$POD" --timeout=60s > /dev/null

uid=$(k get pod "$POD" -o jsonpath='{.metadata.uid}')
cg=/sys/fs/cgroup/kubepods/burstable/pod$uid
{
  echo "area=$AREA runtime_class=$RUNTIME_CLASS pod=$POD cgroup=$cg"
  date -u +%Y-%m-%dT%H:%M:%SZ
  docker exec "$NODE" sh -c "echo cpu.max=\$(cat $cg/cpu.max) memory.max=\$(cat $cg/memory.max)"
  echo "--- samples every 0.5 s (cgroup of the Pod on the node)"
} > "$OUT/$AREA-cgroup.txt"

# the cgroup disappears when the sandbox dies, so sample on the node while the test runs
docker exec "$NODE" sh -c "end=\$((\$(date +%s)+14)); while [ \$(date +%s) -lt \$end ]; do
  [ -d $cg ] && echo \"mem_cur=\$(( \$(cat $cg/memory.current)/1048576 ))MiB mem_peak=\$(( \$(cat $cg/memory.peak)/1048576 ))MiB | \$(tr '\n' ' ' < $cg/memory.events) | \$(grep -E 'usage_usec|nr_periods|nr_throttled|throttled_usec' $cg/cpu.stat | tr '\n' ' ')\" 2>/dev/null
  sleep 0.5; done" >> "$OUT/$AREA-cgroup.txt" 2>/dev/null &
sampler=$!
sleep 1

k exec "$POD" -- /usr/local/bin/trusted.test -test.v -test.run "$TEST" > "$OUT/$AREA-test.log" 2>&1 || true
wait "$sampler" || true
k get pod "$POD" -o jsonpath='{.status.containerStatuses[0].state}' > "$OUT/$AREA-pod-state.json"
k delete pod "$POD" --now > /dev/null
echo "saved to $OUT/"
