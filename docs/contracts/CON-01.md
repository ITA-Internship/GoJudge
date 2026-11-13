# CON-01: Submission lifecycle and verdict model

**Owner:** Pair D
**Dependent Pairs:** Pair A, Pair B, Pair C

## 1. Submission Lifecycle States
Every code submission transitions through an asynchronous lifecycle with the following states:
* **QUEUED**: The submission has been received and stored, waiting for an available evaluation worker.
* **RUNNING**: The evaluation engine is executing the Go test suite against the submission.
* **COMPLETED**: The evaluation finished successfully, and test execution results have been collected. This is a terminal state.
* **FAILED**: The evaluation process aborted unexpectedly before a test verdict could be determined. This is a terminal state.
* **CANCELLED**: The execution was manually aborted by the user before completion. This is a terminal state.

### State Transitions
* Standard flow: `QUEUED` -> `RUNNING` -> `COMPLETED` | `FAILED`
* Cancellation flow: `QUEUED` | `RUNNING` -> `CANCELLED`
* Infrastructure retry flow: `RUNNING` -> `QUEUED` (triggered by the orchestrator).

### Retry Policy and Idempotency
* **Retry Bound:** Infrastructure retries are strictly bounded (e.g., maximum 3 attempts). If retries are exhausted, the state transitions to `FAILED`.
* **Idempotency Rule (First Terminal Write Wins):** To prevent late results from a stalled first attempt conflicting with a subsequent retry, the system enforces a "first terminal write wins" policy. Once a submission transitions to a terminal state (`COMPLETED`, `FAILED`, `CANCELLED`), any subsequent state update requests for that submission are ignored.

### State-Verdict Invariants
* The verdict must remain `null` while the state is `QUEUED` or `RUNNING`.
* The `COMPLETED` state strictly maps to a Test Verdict (e.g., `ACCEPTED`, `WRONG_ANSWER`).
* The `FAILED` state strictly maps to the System Verdict (`INTERNAL_ERROR`).

*Open Questions:*
* **Duplicate Submissions:** If a user submits a solution while another is still `RUNNING`, should the new submission be queued, or should the system return a conflict error? (Owner: Pair D, Deadline: Start of Sprint 2, Default: Return conflict error).

## 2. Evaluation Verdicts
When a submission reaches the `COMPLETED` state, it receives one of the following test verdicts:
* **ACCEPTED**: All tests passed successfully.
* **WRONG_ANSWER**: The code ran, but at least one test failed.
* **COMPILATION_ERROR**: The submitted Go code failed to compile.
* **RUNTIME_ERROR**: The code compiled but crashed during execution (e.g., non-zero exit code or panic).
* **TIME_LIMIT_EXCEEDED**: The execution exceeded the allowed time limit (e.g., Kubernetes Job `activeDeadlineSeconds` reached).
* **MEMORY_LIMIT_EXCEEDED**: The execution exceeded the allowed memory limit.
* **OUTPUT_LIMIT_EXCEEDED**: The execution output is larger than the bounded-output limit.

When a submission reaches the `FAILED` state, it is assigned the system verdict:
* **INTERNAL_ERROR**: An infrastructure or system failure prevented evaluation (e.g., the pod never starts due to an image pull error or pending node, the pod is evicted, the node is lost, or the worker crashes mid-run).

## 3. Infrastructure Failures vs. Test Failures
To ensure clarity, failures are strictly categorized:
* **Test Failures**: Include `WRONG_ANSWER`, `COMPILATION_ERROR`, `RUNTIME_ERROR`, `TIME_LIMIT_EXCEEDED`, `MEMORY_LIMIT_EXCEEDED`, and `OUTPUT_LIMIT_EXCEEDED`. These indicate issues with the learner's code. **Test failures are final and must never trigger automatic execution retries.** Detailed test feedback may be returned according to the visibility rules defined for each problem.
* **Infrastructure Failures**: Manifest as system-level issues not caused by the learner's code. Only infrastructure failures may be retried. The `INTERNAL_ERROR` verdict is the terminal state returned only after all infrastructure retries are exhausted.
* **OOM Conflict Resolution:** If termination is due to memory constraints, a container-level `OOMKilled` (hitting the cgroup limit) is classified as a Test Failure (`MEMORY_LIMIT_EXCEEDED`). A node-level memory eviction (Node Pressure) is classified as an Infrastructure Failure (`INTERNAL_ERROR`).