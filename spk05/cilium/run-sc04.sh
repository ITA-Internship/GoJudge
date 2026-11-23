#!/usr/bin/env bash
# SC-04 egress test on the Cilium cluster
# Same JSONL fields as spk01/scripts/run-sc.ps1
# Usage: ./spk05/cilium/run-sc04.sh [runs] [label]
# Control (no policy): ./spk05/cilium/run-sc04.sh 1 SC-04-control
set -euo pipefail

CONTEXT=k3d-spk05-cilium
NS=spk01
MANIFEST=spk01/manifests/sc-04-netdial.yaml
RUNS=${1:-5}
SC=${2:-SC-04}
TIMEOUT_SEC=120

RESULTS=spk05/cilium/results
mkdir -p "$RESULTS"
JSONL="$RESULTS/$(tr "[:upper:]" "[:lower:]" <<< "$SC").jsonl"

k() { kubectl --context "$CONTEXT" -n "$NS" "$@"; }
now_ms() { perl -MTime::HiRes=time -e 'printf "%d", time*1000'; }
now_iso() { perl -MTime::HiRes=time -MPOSIX=strftime -e '$t=time; printf "%s.%03dZ", strftime("%Y-%m-%dT%H:%M:%S", gmtime($t)), ($t-int($t))*1000'; }

for i in $(seq 1 "$RUNS"); do
  run=$(printf "%02d" "$i")
  created_at=$(now_iso)
  start_ms=$(now_ms)

  job=$(k create -f "$MANIFEST" -o name)
  job_name=${job#*/}
  echo "Run $run -> $job_name"

  condition=""
  deadline=$((SECONDS + TIMEOUT_SEC))
  while [[ -z $condition && $SECONDS -lt $deadline ]]; do
    sleep 1
    if [[ -n $(k get "$job" -o jsonpath='{.status.succeeded}') ]]; then condition=Complete
    elif [[ -n $(k get "$job" -o jsonpath='{.status.failed}') ]]; then condition=Failed
    fi
  done
  finished_at=$(now_iso)
  end_ms=$(now_ms)
  [[ -z $condition ]] && condition=Timeout

  pod=$(k get pods -l "job-name=$job_name" -o name)
  pod=${pod#*/}
  reason=$(k get pod "$pod" -o jsonpath='{.status.containerStatuses[0].state.terminated.reason}')
  code=$(k get pod "$pod" -o jsonpath='{.status.containerStatuses[0].state.terminated.exitCode}')

  log_file="$RESULTS/$SC-run$run.log"
  k logs "$pod" > "$log_file"
  k get pod "$pod" -o yaml > "$RESULTS/$SC-run$run-pod.yaml"
  log_bytes=$(wc -c < "$log_file" | tr -d ' ')

  events=$(grep '^{' "$log_file" || true)
  tests_run=$(jq -s '[.[] | select(.Action=="run" and .Test)] | length' <<< "$events")
  tests_passed=$(jq -s '[.[] | select(.Action=="pass" and .Test)] | length' <<< "$events")
  tests_failed=$(jq -s '[.[] | select(.Action=="fail" and .Test)] | length' <<< "$events")
  verdict=$(jq -rs '[.[] | select((.Test|not) and (.Action=="pass" or .Action=="fail"))] | last | .Action // ""' <<< "$events")

  # pass = test started and all connections blocked
  pass=false
  [[ $tests_run -ge 1 && $verdict == pass ]] && pass=true

  jq -cn \
    --arg criterion "$SC" --argjson run_index "$i" \
    --arg job_name "$job_name" --arg pod_name "$pod" \
    --arg created_at "$created_at" --arg job_finished_at "$finished_at" \
    --argjson total_duration_ms "$((end_ms - start_ms))" \
    --arg condition "$condition" --arg terminated_reason "$reason" \
    --argjson exit_code "${code:-0}" \
    --argjson tests_passed "$tests_passed" --argjson tests_failed "$tests_failed" \
    --arg package_verdict "$verdict" --argjson pass "$pass" \
    --arg log_file "results/$SC-run$run.log" --argjson log_bytes "$log_bytes" \
    '$ARGS.named' >> "$JSONL"

  echo "  condition=$condition reason=$reason exitCode=$code tests=$tests_run verdict=$verdict pass=$pass"
  jq -r 'select(.Output) | .Output | select(test("outbound TCP")) | "    " + (sub("^\\s+netdial_test.go:[0-9]+: "; "") | rtrimstr("\n"))' <<< "$events"
done
