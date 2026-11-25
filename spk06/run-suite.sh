#!/usr/bin/env bash
# SPK-06 isolation suite: runs every check and saves raw evidence
# Usage: ./spk06/run-suite.sh [area ...]   (default: all areas)
# CONTEXT, RUNTIME_CLASS and RESULTS can be set to run it on another candidate
set -euo pipefail

CONTEXT=${CONTEXT:-k3d-spk06}
RUNTIME_CLASS=${RUNTIME_CLASS:-}
RESULTS=${RESULTS:-spk06/results}
NS=spk06
TEMPLATE=spk06/manifests/job.tmpl.yaml
# experimental parameters, taken from SPK-01 (pids limit is our own, set in setup.sh)
OUTPUT_LIMIT=1048576
TIMEOUT_SEC=5

T2J='"/usr/local/bin/test2json", "-t", "/usr/local/bin/trusted.test", "-test.v=test2json", "-test.run"'

# area | runs | cpu | memory | deadline | command
CHECKS="
network|5|1|512Mi|60|[$T2J, \"^TestNetworkEgress\$\"]
cpu|3|1|512Mi|60|[$T2J, \"^TestCPUBurn\$\"]
memory|3|1|512Mi|60|[$T2J, \"^TestMemoryExhaustion\$\"]
timeout|3|1|512Mi|$TIMEOUT_SEC|[$T2J, \"^TestInfiniteLoop\$\"]
filesystem|3|1|512Mi|60|[$T2J, \"^TestFilesystemWrite\$\"]
output|3|1|512Mi|60|[\"/usr/local/bin/outcap\", \"-limit\", \"$OUTPUT_LIMIT\", \"--\", $T2J, \"^TestOutputFlood\$\"]
process|3|1|512Mi|60|[$T2J, \"^TestProcessLimit\$\"]
"

k() { kubectl --context "$CONTEXT" -n "$NS" "$@"; }
now_ms() { perl -MTime::HiRes=time -e 'printf "%d", time*1000'; }

mkdir -p "$RESULTS"
{
  date -u +%Y-%m-%dT%H:%M:%SZ
  echo "context=$CONTEXT runtime_class=${RUNTIME_CLASS:-default}"
  kubectl --context "$CONTEXT" version 2>&1
  kubectl --context "$CONTEXT" get nodes -o wide 2>&1
  docker image inspect gojudge-spk06:dev -f 'image {{.Id}}' 2>&1
} > "$RESULTS/env.txt"

selected=" $* "
while IFS='|' read -r area runs cpu memory deadline command; do
  [[ -z $area ]] && continue
  [[ $# -gt 0 && $selected != *" $area "* ]] && continue
  : > "$RESULTS/$area.jsonl"

  for i in $(seq 1 "$runs"); do
    run=$(printf "%02d" "$i")
    start_ms=$(now_ms)
    job=$(sed -e "s|__AREA__|$area|g" -e "s|__DEADLINE__|$deadline|" \
              -e "s|__CPU__|$cpu|" -e "s|__MEMORY__|$memory|g" \
              -e "s|__COMMAND__|$command|" \
              -e "s|__RUNTIME_CLASS__|$RUNTIME_CLASS|" \
              -e '/runtimeClassName: *$/d' "$TEMPLATE" | k create -f - -o name)
    job_name=${job#*/}
    echo "$area run $run -> $job_name"

    deadline_at=$((SECONDS + deadline + 60))
    state=""
    while [[ -z $state && $SECONDS -lt $deadline_at ]]; do
      sleep 1
      state=$(k get "$job" -o jsonpath='{.status.conditions[?(@.status=="True")].type}' | grep -oE 'Complete|Failed' | head -1 || true)
    done
    [[ -z $state ]] && state=Timeout
    duration_ms=$(( $(now_ms) - start_ms ))

    k get "$job" -o json > "$RESULTS/$area-run$run-job.json"
    job_reason=$(jq -r '[.status.conditions[]? | select(.type=="Failed" and .status=="True") | .reason] | first // ""' "$RESULTS/$area-run$run-job.json")
    # seconds from Job start to its Failed condition (timeout check)
    job_seconds=$(jq -r '[.status.conditions[]? | select(.type=="Failed" and .status=="True") | .lastTransitionTime] | first // empty | fromdateiso8601 as $f | $f - ($ARGS.named.s | fromdateiso8601)' --arg s "$(jq -r '.status.startTime' "$RESULTS/$area-run$run-job.json")" "$RESULTS/$area-run$run-job.json" 2>/dev/null || true)

    pod=$(k get pods -l "job-name=$job_name" -o name | head -1)
    pod=${pod#*/}
    reason="" code=""
    log_file="$RESULTS/$area-run$run.log"
    : > "$log_file"
    if [[ -n $pod ]]; then
      reason=$(k get pod "$pod" -o jsonpath='{.status.containerStatuses[0].state.terminated.reason}')
      code=$(k get pod "$pod" -o jsonpath='{.status.containerStatuses[0].state.terminated.exitCode}')
      k logs "$pod" > "$log_file" 2>&1 || true
      k get pod "$pod" -o yaml > "$RESULTS/$area-run$run-pod.yaml"
    fi
    log_bytes=$(wc -c < "$log_file" | tr -d ' ')
    events=$(grep '^{' "$log_file" || true)
    verdict=$(jq -rs '[.[] | select((.Test|not) and (.Action=="pass" or .Action=="fail"))] | last | .Action // ""' <<< "$events")

    # pass/fail only where the requirement is binary; cpu, output and process
    # have no approved threshold, so only the behaviour is recorded (pass=null)
    pass=false
    case $area in
      memory)  [[ $reason == OOMKilled ]] && pass=true ;;
      timeout) [[ $job_reason == DeadlineExceeded && -n $job_seconds && $job_seconds -le $((deadline + 5)) ]] && pass=true ;;
      cpu|output|process) pass=null ;;
      *)       [[ $verdict == pass ]] && pass=true ;;
    esac

    jq -cn \
      --arg area "$area" --argjson run_index "$i" \
      --arg context "$CONTEXT" --arg runtime_class "${RUNTIME_CLASS:-default}" \
      --arg cpu_limit "$cpu" --arg memory_limit "$memory" --argjson deadline_s "$deadline" \
      --arg job_name "$job_name" --arg pod_name "$pod" \
      --arg state "$state" --arg job_reason "$job_reason" --arg job_seconds "$job_seconds" \
      --arg terminated_reason "$reason" --arg exit_code "$code" \
      --arg verdict "$verdict" --argjson log_bytes "$log_bytes" \
      --argjson duration_ms "$duration_ms" --argjson pass "$pass" \
      --arg log_file "$log_file" \
      '$ARGS.named' >> "$RESULTS/$area.jsonl"

    echo "  state=$state job_reason=$job_reason reason=$reason exit=$code verdict=$verdict log_bytes=$log_bytes pass=$pass"
  done
done <<< "$CHECKS"
