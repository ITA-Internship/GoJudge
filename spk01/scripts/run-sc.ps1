param(
    [Parameter(Mandatory)] [string]$Sc,         # e.g. SC-06
    [Parameter(Mandatory)] [string]$Manifest,   # e.g. manifests/sc-06-nonroot.yaml
    [Parameter(Mandatory)] [string]$Expect,     # "pass" or a termination reason, e.g. "OOMKilled"
    [int]$Runs = 5,
    [int]$TimeoutSec = 120
)

$ns = "spk01"
$results = Join-Path $PWD.Path "results"
New-Item -ItemType Directory -Force -Path $results | Out-Null
$jsonl = Join-Path $results ("{0}.jsonl" -f $Sc.ToLower())
$utf8 = New-Object System.Text.UTF8Encoding($false)   # UTF-8 without BOM

for ($i = 1; $i -le $Runs; $i++) {
    $run = "{0:D2}" -f $i
    $createdAt = Get-Date

    $job = kubectl create -f $Manifest -o name
    $jobName = $job.Split("/")[1]
    Write-Host "Run $run -> $jobName"

    # Wait until the Job has either succeeded or failed
    $condition = ""
    $deadline = (Get-Date).AddSeconds($TimeoutSec)
    while (-not $condition -and (Get-Date) -lt $deadline) {
        Start-Sleep -Seconds 1
        if (kubectl get $job -n $ns -o jsonpath='{.status.succeeded}') { $condition = "Complete" }
        elseif (kubectl get $job -n $ns -o jsonpath='{.status.failed}') { $condition = "Failed" }
    }
    $finishedAt = Get-Date
    if (-not $condition) { $condition = "Timeout" }

    $pod = (kubectl get pods -n $ns -l "job-name=$jobName" -o name).Split("/")[1]
    $reason = kubectl get pod $pod -n $ns -o jsonpath='{.status.containerStatuses[0].state.terminated.reason}'
    $code   = kubectl get pod $pod -n $ns -o jsonpath='{.status.containerStatuses[0].state.terminated.exitCode}'

    # Save raw evidence immediately (UTF-8 without BOM)
    $lines = @(kubectl logs $pod -n $ns)
    [IO.File]::WriteAllLines((Join-Path $results "$Sc-run$run.log"), [string[]]$lines, $utf8)
    $podYaml = (kubectl get pod $pod -n $ns -o yaml) -join "`n"
    [IO.File]::WriteAllText((Join-Path $results "$Sc-run$run-pod.yaml"), $podYaml, $utf8)

    # Save init container logs, if the Pod has any
        $inits = kubectl get pod $pod -n $ns -o jsonpath='{.spec.initContainers[*].name}'
        foreach ($ic in ($inits -split ' ' | Where-Object { $_ })) {
            $icLines = @(kubectl logs $pod -n $ns -c $ic)
            [IO.File]::WriteAllLines((Join-Path $results "$Sc-run$run-$ic.log"), [string[]]$icLines, $utf8)
        }

    # Parse go test -json events
    $events = $lines | Where-Object { $_ -like '{*' } | ForEach-Object { $_ | ConvertFrom-Json }
    $testsRun    = @($events | Where-Object { $_.Action -eq 'run'  -and $_.Test }).Count
    $testsPassed = @($events | Where-Object { $_.Action -eq 'pass' -and $_.Test }).Count
    $testsFailed = @($events | Where-Object { $_.Action -eq 'fail' -and $_.Test }).Count
    $verdict = ($events | Where-Object { -not $_.Test -and ($_.Action -eq 'pass' -or $_.Action -eq 'fail') } |
                Select-Object -Last 1).Action

    # A run counts only if the test really started AND the expected outcome happened
    $pass = ($testsRun -ge 1) -and (($reason -eq $Expect) -or ($verdict -eq $Expect))

    $record = [ordered]@{
        criterion         = $Sc
        run_index         = $i
        job_name          = $jobName
        pod_name          = $pod
        created_at        = $createdAt.ToString("o")
        job_finished_at   = $finishedAt.ToString("o")
        total_duration_ms = [int]($finishedAt - $createdAt).TotalMilliseconds
        condition         = $condition
        terminated_reason = $reason
        exit_code         = [int]$code
        tests_passed      = $testsPassed
        tests_failed      = $testsFailed
        package_verdict   = "$verdict"
        pass              = $pass
        log_file          = "results/$Sc-run$run.log"
    }
    [IO.File]::AppendAllText($jsonl, (($record | ConvertTo-Json -Compress) + "`n"), $utf8)

    Write-Host "  condition=$condition reason=$reason exitCode=$code tests=$testsRun verdict=$verdict pass=$pass"
}