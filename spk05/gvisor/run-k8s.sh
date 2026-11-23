#!/usr/bin/env bash
# usage: run-k8s.sh <label> <manifest> [runs]
label=$1; manifest=$2; runs=${3:-5}
K="sudo k3s kubectl"; out=spk05/gvisor/k8s_results
for i in $(seq 1 "$runs"); do
  job=$($K create -f "$manifest" -o name); name=${job#*/}
  for t in $(seq 1 90); do
    s=$($K get "$job" -n spk01 -o jsonpath='{.status.succeeded}')
    f=$($K get "$job" -n spk01 -o jsonpath='{.status.failed}')
    [ -n "$s$f" ] && break; sleep 1
  done
  pod=$($K get pods -n spk01 -l job-name="$name" -o name)
  $K logs -n spk01 "$pod" > "$out/$label-run$i.log" 2>&1
  code=$($K get -n spk01 "$pod" -o jsonpath='{.status.containerStatuses[0].state.terminated.exitCode}')
  [ "$i" = 1 ] && $K get -n spk01 "$pod" -o yaml > "$out/$label-pod.yaml"
  echo "$label run$i job=$name succeeded=${s:-0} failed=${f:-0} exit=$code" | tee -a "$out/summary.txt"
done
