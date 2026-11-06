# CON-01: Submission lifecycle and verdict model

**Owner:** Pair D
**Dependent Pairs:** Pair A, Pair B, Pair C

## 1. Submission Lifecycle States
Every code submission transitions through an asynchronous lifecycle with the following states:
* **QUEUED**: The submission has been received and stored, waiting for an available evaluation worker.
* **RUNNING**: The evaluation engine is executing the Go test suite against the submission.
* **COMPLETED**: The evaluation finished successfully, and test execution results have been collected. This is a terminal state.
* **FAILED**: The evaluation process aborted unexpectedly before a test verdict could be determined. This is a terminal state.

### State Transitions
* Standard flow: `QUEUED` -> `RUNNING` -> `COMPLETED` | `FAILED`
* Infrastructure retry flow: `RUNNING` -> `QUEUED` (if the system drops the job and explicitly requeues it).

*Open Questions:*
* **Cancellation:** Is submission cancellation (e.g., SPK-01 / SC-10) in scope for the MVP? If so, does it have its own state, or does it transition to `FAILED`?

## 2. Evaluation Verdicts
When a submission reaches the `COMPLETED` state, it receives one of the following test verdicts:
* **ACCEPTED**: All tests passed successfully.
* **WRONG_ANSWER**: The code ran, but at least one test failed.
* **COMPILATION_ERROR**: The submitted Go code failed to compile.
* **RUNTIME_ERROR**: The code compiled but crashed during execution (e.g., non-zero exit code or panic).
* **TIME_LIMIT_EXCEEDED**: The execution exceeded the allowed time limit (e.g., Kubernetes Job `activeDeadlineSeconds` reached).
* **MEMORY_LIMIT_EXCEEDED**: The execution exceeded the allowed memory limit (e.g., Kubernetes Container `OOMKilled`).

*Open Questions:*
* **Bounded Output Limit:** Which verdict should be assigned if the execution output is larger than the bounded-output limit mentioned in MVP Proposal §7?

When a submission reaches the `FAILED` state, it is assigned the system verdict:
* **INTERNAL_ERROR**: An infrastructure or system failure prevented evaluation (e.g., the pod never starts due to an image pull error or pending node, the pod is evicted, the node is lost, or the worker crashes mid-run).

## 3. Infrastructure Failures vs. Test Failures
To ensure clarity, failures are strictly categorized:
* **Test Failures**: Include `WRONG_ANSWER`, `COMPILATION_ERROR`, `RUNTIME_ERROR`, `TIME_LIMIT_EXCEEDED`, and `MEMORY_LIMIT_EXCEEDED`. These indicate issues with the learner's code. **Test failures are final and must never trigger automatic execution retries.** Detailed test feedback may be returned according to the visibility rules defined for each problem.
* **Infrastructure Failures**: Manifest strictly as an `INTERNAL_ERROR`. This indicates a failure in the GoJudge system infrastructure, not the learner's code. Only infrastructure failures may be retried. If an infrastructure failure triggers a retry, the retry mechanism must not create duplicate or conflicting final results.