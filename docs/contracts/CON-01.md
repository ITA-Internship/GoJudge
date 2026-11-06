## 1. Submission Lifecycle States
Every code submission transitions through an asynchronous lifecycle with the following states:
* **QUEUED**: The submission has been received and stored, waiting for an available evaluation worker.
* **RUNNING**: The evaluation engine is executing the Go test suite against the submission.
* **COMPLETED**: The evaluation finished successfully, and test execution results have been collected.
* **FAILED**: The evaluation process aborted unexpectedly before a test verdict could be determined.

### State Transitions
`QUEUED` -> `RUNNING` -> `COMPLETED` | `FAILED`

## 2. Evaluation Verdicts
When a submission reaches the `COMPLETED` state, it receives one of the following test verdicts:
* **ACCEPTED**: All tests passed successfully.
* **WRONG_ANSWER**: The code executed, but the output did not match the expected result.
* **COMPILATION_ERROR**: The submitted Go code failed to compile.
* **RUNTIME_ERROR**: The code compiled but crashed during execution (e.g., panic).
* **TIME_LIMIT_EXCEEDED**: The execution exceeded the allowed time limit.
* **MEMORY_LIMIT_EXCEEDED**: The execution exceeded the allowed memory limit.

When a submission reaches the `FAILED` state, it is assigned the system verdict:
* **INTERNAL_ERROR**: An infrastructure or system failure prevented evaluation.

## 3. Infrastructure Failures vs. Test Failures
To ensure clarity, failures are strictly categorized:
* **Test Failures**: Include `WRONG_ANSWER`, `COMPILATION_ERROR`, `RUNTIME_ERROR`, `TIME_LIMIT_EXCEEDED`, and `MEMORY_LIMIT_EXCEEDED`. These indicate issues with the learner's code. Detailed test feedback may be returned according to the visibility rules defined for each problem.
* **Infrastructure Failures**: Manifest strictly as an `INTERNAL_ERROR`. This indicates a failure in the GoJudge system, not the learner's code. If an infrastructure failure triggers a retry, the retry must not create duplicate or conflicting final results.