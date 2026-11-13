# SonarQube Cloud

SonarQube Cloud checks every pull request into `main`. It does not replace
tests, the linter or code review.

## How it works

1. The `Unit tests` job in the `CI` workflow runs the tests and saves the
   coverage to `coverage.out`.
2. The same job runs the SonarQube scan, which sends the code and the coverage
   report to SonarQube Cloud. The scan settings are in
   [`sonar-project.properties`](/sonar-project.properties).
3. SonarQube Cloud checks the new code of the pull request and reports the
   result to GitHub.
4. With `sonar.qualitygate.wait=true` the scan waits for the Quality Gate. The
   job fails when the gate fails, when the analysis fails, or when no result
   arrives within 300 seconds.
5. The ruleset for `main` requires this check, so the pull request cannot be
   merged until it passes.

`main` itself is not analysed separately: the workflow runs only on pull
requests. Coverage and issues are reported for the new code of each pull
request.

### What is analysed

- Go source code, Kubernetes manifests in `infra/` and the workflows in
  `.github/`.
- Test files (`*_test.go`) count as tests, not as source code.
- Excluded: `vendor/`, `testdata/`, generated code (`*.pb.go`, `*_gen.go`,
  `zz_generated*.go`), `.idea/`, and documentation in `docs/` and `specs/`.

## Quality Gate

Only new code is checked (Clean as You Code). Existing code does not block a
pull request.

| Condition on new code | Requirement |
|-----------------------|-------------|
| Issues | No new issues |
| Security Hotspots | 100% reviewed |
| Coverage | At least 80% |
| Duplicated lines | At most 3% |

**Small changes.** When a pull request changes fewer than 20 lines, SonarQube
does not apply the coverage and duplication conditions.

**Plan and gate in use.** The organisation `ita-internship` is on the SonarQube
Cloud **Free plan**. The project uses the built-in **Sonar way** gate, which is
the organisation default. Its conditions on new code:

| Condition | Requirement |
|-----------|-------------|
| Reliability and security | Rated A (no new bugs or vulnerabilities) |
| Technical debt | Maintainability rating is A |
| Security Hotspots Reviewed | 100% |
| Coverage | At least 80.0% |
| Duplicated Lines (%) | At most 3.0% |

The ratings are how Sonar way expresses "no new actionable issues": a new bug
or vulnerability drops the rating below A and fails the gate.

**Free plan limits.** Custom quality gates and quality profiles are not
available, so the built-in Sonar way is used and cannot be edited. The
organisation can have up to 5 members, and only pull requests that target
`main` are analysed (pull requests into other branches get no SonarQube result).

## Reading the results

- Open the SonarQube check on the pull request and follow the link to the
  analysis.
- **Issues** are grouped by software quality (Security, Reliability,
  Maintainability) and severity. Each issue links to its rule, which explains
  why it matters and how to fix it.
- **Security Hotspots** need a human review. They are not always bugs. Each one
  must be marked Safe or Fixed in SonarQube Cloud.
- **Coverage** shows which new lines no test runs.

## If the check fails

| Failed condition | What to do |
|------------------|------------|
| New issues | Fix the code as the rule describes, then push again |
| Security Hotspots not reviewed | Review each hotspot in SonarQube Cloud and mark it Safe or Fixed with a short reason |
| Coverage below 80% | Add tests for the new lines that are not covered |
| Duplication above 3% | Move the repeated code into a shared function |
| Analysis failed or timed out | Open the job log. Re-run the check if SonarQube Cloud was unavailable. Merging stays blocked until it passes |

Do not turn off the check, change the Quality Gate or add exclusions to make a
pull request pass.

## False positive

1. Write in the pull request which rule, file and line, and why you think the
   finding is wrong.
2. Ask a reviewer from another pair to confirm.
3. A project administrator marks the issue as False positive or Accepted in
   SonarQube Cloud, with the reason. No suppression in code (`//NOSONAR`)
   without this review.

## Setup (administrators)

1. Sign in to [SonarQube Cloud](https://sonarcloud.io) with GitHub. Import the
   organisation `ITA-Internship` and the repository `GoJudge`.
2. In Administration → Analysis Method, turn off Automatic Analysis. Coverage
   needs the analysis from CI.
3. Create a token. In GitHub, add it as the Actions secret `SONAR_TOKEN_G`
   (Settings → Secrets and variables → Actions). Never commit it.
4. Copy the project key and organisation key from Project information into
   `sonar-project.properties`.
5. Set the Quality Gate and the new-code definition.
6. In the ruleset for `main`, keep the required checks `Build`, `Lint` and
   `Unit tests` (which now includes the SonarQube scan) and the review from
   another pair.

## Troubleshooting

| Problem | Cause and fix |
|---------|---------------|
| Coverage is 0% | `coverage.out` was not found. Check `sonar.go.coverage.reportPaths` |
| `Project not found` or 401 | Wrong project key or organisation, or `SONAR_TOKEN_G` is missing or expired |
| Analysis runs twice or CI analysis is rejected | Automatic Analysis is still on. Turn it off |
| All code counts as new | The checkout is shallow. The job needs `fetch-depth: 0` |
| Quality Gate timeout | SonarQube Cloud is slow or down. Re-run the check later |
| Pull request from a fork has no result | Secrets are not available to forks. Work from branches in this repository |