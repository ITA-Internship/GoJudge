# SonarQube Cloud
SonarQube Cloud checks every pull request into `main`. It does not replace tests, the linter or code review

## How it works
`CI` workflow runs the tests, saves the coverage to `coverage.out` and sends the code to SonarQube Cloud.
SonarQube Cloud checks the new code of the pull request and reports the result to GitHub.
If the Quality Gate fails, or the analysis does not finish, the pull request cannot be merged.

The scan settings are in [`sonar-project.properties`](../sonar-project.properties)

## Quality Gate
Only new code is checked:

- no new issues
- all security hotspots reviewed
- coverage of new code at least 80%
- duplication of new code at most 3%

TODO: plan limitations (can the conditions be changed?).

## If the check fails
1. Open the SonarQube check in the pull request and read which condition failed
2. Fix code, add tests, or review the hotspots in SonarQube Cloud, then push again
3. If the analysis itself failed, re-run the check. Merging stays blocked until it passes

## False positive
Do not turn off the check or the rule.
Write in the pull request why you think the finding is wrong and ask the reviewer.
Any suppression needs a written reason
