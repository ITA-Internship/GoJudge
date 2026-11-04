param([int]$Runs = 5)

New-Item -ItemType Directory -Force -Path results | Out-Null

for ($i = 1; $i -le $Runs; $i++) {
    $run = "{0:D2}" -f $i

    # 1. Create the Job; -o name returns "job.batch/spk01-sc-03-xxxxx"
    $job = kubectl create -f manifests/sc-03-memhog.yaml -o name
    $jobName = $job.Split("/")[1]
    Write-Host "Run $run -> $jobName"

    # 2. Wait until the Job fails (OOM or, as a fallback, the 60 s deadline)
    kubectl wait -n spk01 --for=condition=Failed $job --timeout=90s | Out-Null

    # 3. Find the Pod that belongs to this Job
    $pod = (kubectl get pods -n spk01 -l "job-name=$jobName" -o name).Split("/")[1]

    # 4. Read the termination reason and exit code
    $reason = kubectl get pod $pod -n spk01 -o jsonpath='{.status.containerStatuses[0].state.terminated.reason}'
    $code   = kubectl get pod $pod -n spk01 -o jsonpath='{.status.containerStatuses[0].state.terminated.exitCode}'

    # 5. Save the evidence immediately, in UTF-8
    kubectl logs $pod -n spk01 | Set-Content -Encoding utf8 "results/SC-03-run$run.log"
    kubectl get pod $pod -n spk01 -o yaml | Set-Content -Encoding utf8 "results/SC-03-run$run-pod.yaml"
    "SC-03 run$run job=$jobName pod=$pod reason=$reason exitCode=$code" |
        Add-Content -Encoding utf8 "results/SC-03-summary.txt"

    Write-Host "  reason=$reason exitCode=$code"
}