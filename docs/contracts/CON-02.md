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
* **Infrastructure failures** (DB, Kubernetes, ...): retried 3 times with waits of 1s, 2s and 4s
(Watermill `Retry`: initial interval 1s, multiplier 2, max interval 4s)
After that the submission becomes `FAILED` / `INTERNAL_ERROR`

Why: about 7s covers short failures (a dropped DB connection, a Kubernetes API restart),
and longer waits would block the queue for other submissions. Longer outages go to the dead-letter queue (§7)

## 7. Dead-Letter
A message goes to `submissions.dead` when its retries are exhausted, or at once if it cannot be parsed.
The broker does the routing (SPK-03 SC-D, variant B). Messages there are checked by hand.

## 8. Dispatch
How the API keeps the database write and the publish consistent is DEC-12 (open, not approved).
SPK-03 recommends the transactional outbox (SC-E)

## 9. Risks (from SPK-03)
These were not verified by SPK-03. They are risks, not guarantees.
1. The retry counter is in process memory: a crash resets it to zero.
2. Watermill cannot Nack single messages differently, so a broker-side retry counter is not possible through Watermill.
3. `x-delivery-limit` (quorum queues) does not count Nack with requeue, so it is not a retry limit.
4. A consumer is blocked while it retries one message (§6).
5. Without the dead-letter exchange and queue the broker silently drops dead-lettered messages (§2).
6. The idempotency check was tested in memory only, not in PostgreSQL.
7. The "check, then apply" step is not atomic: several consumers need an atomic check (§4).
8. All failures were simulated (`os.Exit(1)`, `Nack`), not real crashes or network failures.
9. The outbox is at-least-once: a forwarder crash can send a message twice (§4).
10. Outbox latency (about 1 s polling) and outbox table clean-up were not measured.
11. An outbox row the forwarder cannot unwrap blocks all rows behind it.
12. The publisher needs `ConfirmDelivery`, and a message to a missing queue is dropped (§2).
13. Broker outage during publish, several forwarders, and message order were not tested. Message order is not guaranteed.

*Open Questions:*
* **Duplicate run:** a message is redelivered while a Job for the same submission still runs. How is a second Job avoided?
(Owner: F07-01)
* **Cancellation:** how does `CANCELLED` reach a running evaluation?
(Owner: F07-01)
