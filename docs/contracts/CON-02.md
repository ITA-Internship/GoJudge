# CON-02: Submission messaging contract

**Owner:** Pair B

**Dependent Pairs:** owners of F06-02 (publisher) and F07-01 (consumer);

**Sources:** [CON-01](CON-01.md), [SPK-03](../spikes/SPK-03.md)

> DEC-04 and DEC-12  are not approved. If either changes, this contract changes too and the dependent pairs agree again.

## 1. Flow
The API publishes one message per submission. The Judge Orchestrator consumes it,
runs the evaluation and writes the state and verdict to PostgreSQL. There is no reply message

`API (F06-02) -> submissions.dispatch -> Orchestrator (F07-01) -> PostgreSQL`

## 2. Queues
* **`submissions.dispatch`**: durable queue with new submissions
* **`submissions.dead`**: dead-letter queue, reached through the exchange `submissions.dlx`
* Queues and the dead-letter exchange are declared by code at start up.
Without them the broker silently drops messages (SPK-03 §5 items 5 and 12)

## 3. Message
The message carries IDs only. Source code and tests are loaded from PostgreSQL, never sent through the broker

```json
{
  "submission_id": "uuid",
  "problem_id": "uuid",
  "user_id": "string",
  "submitted_at": "2026-11-17T10:10:10Z"
}
```

Metadata: `idempotency_key = submission_id`

## 4. Idempotency
* Delivery is at least once, so the same submission can arrive twice (SPK-03 SC-C, SC-E)
* Idempotency key is **`submission_id`**, not the message UUID: a producer retry creates a new UUID (SPK-03 SC-C)
* Orchestrator processes a submission only if it has no final state yet (`QUEUED` or `RUNNING`).
If it is already `COMPLETED`, `FAILED` or `CANCELLED`, the message is skipped and acknowledged

## 5. Acknowledgement
Orchestrator acknowledges a message only after the result is written to PostgreSQL.
If it crashes before that, broker redelivers the message (SPK-03 SC-B)

## 6. Retry
* **Test failures** (`WRONG_ANSWER`, `COMPILATION_ERROR`, CON-01 §3): written as `COMPLETED`, acknowledged, never retried
* **Infrastructure failures** (DB, Kubernetes, ...): retried 3 times, 100ms initial interval, doubling, max 2s (SPK-03 SC-D)
After that the submission becomes `FAILED` / `INTERNAL_ERROR`

## 7. Dead-Letter
A message goes to `submissions.dead` when its retries are exhausted, or at once if it cannot be parsed.
The broker does the routing (SPK-03 SC-D, variant B). Messages there are checked by hand.

## 8. Dispatch
How the API keeps the database write and the publish consistent is DEC-12 (open, not approved).
SPK-03 recommends the transactional outbox (SC-E)

## 9. Risks (not verified by SPK-03)
* The retry counter is in memory: a crash resets it (SPK-03 §5, item 1)
* The idempotency check was tested in memory only, not in PostgreSQL (§5, items 6 and 7)
* All failures in SPK-03 were simulated (§5, item 8)

*Open Questions:*
* **Duplicate run:** a message is redelivered while a Job for the same submission still runs. How is a second Job avoided?
(Owner: F07-01)
* **Cancellation:** how does `CANCELLED` reach a running evaluation?
(Owner: F07-01)
