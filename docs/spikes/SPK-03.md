# SPK-03: RabbitMQ and Watermill: Delivery, Redelivery, Idempotency and Dead-Letter Routing

## Metadata
- **Spike ID:** SPK-03
- **Type:** Technical spike
- **Owner:** Pair C
- **Timebox:** 4 pair-days (Days 3–6: Wed 4 Nov – Mon 9 Nov). Report due end of day 6.
- **Status:** Completed within timebox
- **Review follow-up (11 Nov 2026):** SC-E (transactional outbox) was added after the review comments on the report. It is not one of the planned criteria SC-A to SC-D.
- **Dependencies:** S1-01 (INF-03: local RabbitMQ and PostgreSQL), DEC-04 (provisional)
- **Experiment code:** branch `spike/SPK-03-rabbitmq-watermill`, directory `spike/spk-03/RabbitMQ/` (not merged to `main`). All evidence links for SC-A to SC-D point to commit `0b77d50`. The SC-E code and records are in a later commit of the same branch (`98d7603`).
- **Libraries:** `github.com/ThreeDotsLabs/watermill` v1.5.3, `github.com/ThreeDotsLabs/watermill-amqp/v3` v3.1.0 (AMQP client `amqp091-go` v1.11.0), Go 1.27.3. SC-E adds `github.com/ThreeDotsLabs/watermill-sql/v2` v2.0.0 and `github.com/lib/pq` v1.3.0. Broker: RabbitMQ 4.4.0 (Erlang 29.1.2). Database (SC-E only): PostgreSQL 17.11.

## 1. Objective & Hypothesis
- **Objective:** Answer one question: can RabbitMQ with Watermill publish and consume submission messages without loss, redeliver after a consumer failure, keep one final state under duplicate delivery, and route poison messages to a dead-letter route?
- **Hypothesis:** RabbitMQ with Watermill delivers every message, redelivers after a consumer is killed mid-processing, produces one final state under duplicate delivery when an idempotency key is used, and routes poison messages to a dead-letter route.
- **Added after review (SC-E):** does the transactional outbox from the Watermill documentation (the forwarder) keep the database write and the RabbitMQ publish consistent? Hypothesis: a committed business row always produces a message (possibly late), a rolled-back row never produces one, and a crash of the relay can produce a duplicate but not a loss.

## 2. Experiment Setup
- **Broker:** the local RabbitMQ from S1-01, reached on `localhost:5672`.
- **Topology:** Watermill AMQP with `amqp.NewDurableQueueConfig`. A Watermill topic is a durable RabbitMQ queue of the same name. Each criterion uses its own topic (`hello`, `SC-B`, `SC-C`). SC-D uses a separate queue per variant (see the SC-D evidence).
- **Run order:** the queue is declared by the subscriber, not by the publisher, so the consumer is started before a run that publishes messages.
- **Code layout:** one `publisher` and one `consumer` program per criterion in `spike/spk-03/RabbitMQ/SC-A` … `SC-D` (SC-D has three variant folders). SC-E is in `spike/spk-03/RabbitMQ/OUTBOX` and has three programs: `publisher`, `forwarder` and `consumer`. The shared package `rmqenv` builds the AMQP URI (`URI()`) and, for SC-E, the PostgreSQL connection string (`PostgresDSN()`) from the repository-root `.env`.
- **SC-E topology:** PostgreSQL from S1-01, reached on `localhost:5432` through `kubectl port-forward svc/postgres 5432:5432`. Tables: the business table `spk03_submissions` and the outbox table `watermill_outbox_forwarder` (both created by the publisher if missing), and the offsets table `watermill_offsets_outbox_forwarder` (created by the forwarder, which also creates the outbox table if it is missing). The events are delivered to the durable queue `outbox_submissions`, declared by the consumer.
- **Secrets & Environment Policy:** no secrets or credentials are committed. `RABBITMQ_USER` and `RABBITMQ_PASSWORD` (and, for SC-E, `POSTGRES_USER`, `POSTGRES_PASSWORD` and `POSTGRES_DB`) are read from the git-ignored `.env` file, and the raw logs contain no credentials.
- **Raw records:** the console output of each run is stored beside its code as `publisher_log.txt` and `consumer_log.txt` (SC-E: one folder per scenario, with `publisher_log.txt`, `forwarder_log.txt` and `consumer_log.txt`). Two RabbitMQ Management UI screenshots for SC-D (`img.png`) are stored in the same way, on the spike branch. They show no passwords.

## 3. Success Criteria & Recorded Evidence

| Criterion | Status | Evidence Summary |
| :--- | :--- | :--- |
| **SC-A** | **Verified** | 100 of 100 messages published and consumed. The set of 100 message UUIDs is identical on both sides, with no loss and no duplicates. |
| **SC-B** | **Verified** | 5 of 5 consumer kills during processing were followed by redelivery of the same message (same UUID). |
| **SC-C** | **Verified** | 50 of 50 logical operations ended in exactly one applied state, although 100 duplicate deliveries arrived. |
| **SC-D** | **Verified** | 3 of 3 poison messages reached a dead-letter queue after the configured attempts, in two variants: application-level `PoisonQueue` (A) and RabbitMQ's native dead-letter exchange with `Retry` (B). A third variant, a quorum queue with `x-delivery-limit` (C), dead-letters only after consumer loss, not after handler errors. |
| **SC-E** (added after review) | **Verified** | Transactional outbox with the Watermill forwarder, PostgreSQL and RabbitMQ. In 6 of 6 outbox scenarios (E1 to E6) `data_without_message=0` and `message_without_data=0`. The 2 baselines without the outbox (B1, B2) showed the two problems (1 row without a message, 1 message without a row). A forwarder crash gave 1 duplicate, no loss. |

### SC-A Evidence: Publish and Consume (100 of 100)
Raw records: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-A/publisher_log.txt) and [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-A/consumer_log.txt) (branch `spike/SPK-03-rabbitmq-watermill`, commit `0b77d50`). Run on 7 Nov 2026.
- The publisher sent 100 messages (`message #1` … `message #100`), each with its own UUID. The consumer received 100 messages and acknowledged each one.
- The two sets of UUIDs were compared one by one: 100 unique on each side, identical, none missing and none repeated. Payloads `#1` … `#100` are all present. The order is the order of publication.
- Both logs contain no `ERROR` or `WARN` lines, and both processes closed their connections cleanly.
- The consumer was started after the publisher had finished (publisher finished 14:00:20, consumer connected 14:00:50 on 2026/11/07), so the messages waited in the durable queue until a consumer arrived.

### SC-B Evidence: Redelivery After a Consumer Is Killed (5 of 5)
Raw records: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-B/publisher_log.txt) and [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-B/consumer_log.txt) (commit `0b77d50`). Run on 7 Nov 2026.
- A new topic (`SC-B`, a new queue) was used, and one message was published: `id=185ef076-68d1-4d9f-aba5-5a26842354cf`.
- The consumer was started six times in a row. Each of the first five runs received the message and was terminated with `os.Exit(1)` before `Ack`. `os.Exit` skips the deferred `Close`, so the connection is dropped like a killed process.
- Every time, the broker delivered the **same message** (same UUID) to the next consumer.
- Each of the 5 kills was followed by a redelivery of the same message: **5 of 5**.
- Method notes: the failure is simulated by `os.Exit(1)` in the consumer code, not by an external kill. The same message was killed and redelivered five times in sequence, rather than five separate messages being published. For the sixth run `os.Exit` was commented out so the message could be acknowledged.

### SC-C Evidence: Duplicates and One Final State (50 of 50)
Raw records: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-C/publisher_log.txt) and [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-C/consumer_log.txt) (commit `0b77d50`). Run on 8 Nov 2026.

**Idempotency approach used:**
- Each logical operation (one submission) gets an **idempotency key** (a UUID) that travels in the message metadata (`idempotency_key`).
- Before applying the side effect, the consumer checks whether the key is already recorded as applied. If it is, the consumer skips the side effect and acknowledges the message. If it is not, it applies the effect and records the key.
- The message UUID is not used as the key, because a producer retry creates a new message UUID for the same operation.

**What was run:**
- The publisher created 50 keys and published each twice: 50 originals, then 50 **producer duplicates** (same key, new message UUID). That is 100 messages with 100 different UUIDs.
- On the first delivery of a key the consumer applied the effect and then called `Nack`, imitating a lost `Ack`. The broker redelivered the same message (same UUID), and the consumer skipped it. The 50 producer duplicates were skipped as well.

**What was observed** (consumer summary line): `deliveries=150 keys=50 applied_total=50 skipped=100 max_applied_per_key=1`.
- 150 deliveries = 50 first deliveries + 50 redeliveries after `Nack` + 50 producer duplicates.
- The 50 applied keys are exactly the 50 published keys. No key was applied more than once.
- All 100 duplicate deliveries (50 of each kind) were skipped, so every one of the 50 operations has one final state: **50 of 50**.

### SC-D Evidence: Poison Messages and the Dead-Letter Route (3 of 3)
Three variants of the dead-letter route were tried. Every run uses the same scenario: nine messages published in one run (normal, poison, normal, three times), so 3 poison and 6 normal messages per run. The poison messages carry the metadata flag `poison=true`, and the handler always fails on them.  Variant B is recorded as the SC-D result because the dead-lettering is done by the broker. A and C are kept for comparison.

**Common setup:**
- **Topology created by hand in RabbitMQ:** a direct exchange `sc-d.dlx`, a queue `sc-d.dead` and a binding with routing key `sc-d.dead`. The code does not create them. In variants B and C the consumer declares the main queue with the arguments `x-dead-letter-exchange` and `x-dead-letter-routing-key`.
- **Queue names:** `SC-D` and `SC-D.poison` (A), `SC-D.2` (B), `SC-D.quorum` (C). Queue arguments cannot be changed after a queue is created, so each variant has its own queue. The dead-letter queue `sc-d.dead` is shared by B and C.
- **Run order:** the consumer is started first, because the subscriber creates the queue and the publisher does not.

| Variant | Who counts the attempts | Who routes to the dead-letter queue | Result |
| :--- | :--- | :--- | :--- |
| A: `PoisonQueue` + `Retry` (all in the application) | application (`Retry`) | application (`PoisonQueue` publishes to `SC-D.poison`) | 3 of 3 |
| **B: `Retry` + native dead-letter queue** | application (`Retry`, in memory) | **broker** (`NoRequeueOnNack`, then `x-dead-letter-exchange`) | **3 of 3** |
| C: quorum queue with `x-delivery-limit` | broker, but only for consumer loss | broker | normal run: never dead-lettered. Crash run: 3 of 3 |

#### Variant A: `PoisonQueue` + `Retry`
Raw records in `SC-D/middleware/`: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/middleware/publisher_log.txt) and [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/middleware/consumer_log.txt) (commit `0b77d50`). Run on 8 Nov 2026.
- The router uses the `PoisonQueue` middleware, which republishes a failed message to the queue `SC-D.poison` once the handler has given up. The `Retry` middleware sits inside it (`MaxRetries = 3`, initial interval 100 ms), so `PoisonQueue` sees the error only after the retries are exhausted. A separate subscriber on `SC-D.poison` reads the dead-lettered messages and prints the failure reason.
- `MaxRetries = 3` means one initial attempt and three retries, which is **4 handler calls** per poison message.

| Poison message | Handler attempts (`FAILED`) | Appeared in `SC-D.poison` | Reason |
| :--- | :--- | :--- | :--- |
| `10e4ff5f…` | 4 | once | `cannot process message` |
| `9279764f…` | 4 | once | `cannot process message` |
| `87d45598…` | 4 | once | `cannot process message` |

- After reaching the poison queue a message was not processed again, and the following messages were handled normally. The 6 normal messages were each received exactly once.
- Result: **3 of 3**. Risks: see §5, item 1.

#### Variant B: `Retry` + native dead-letter queue (SC-D result)
Raw records in `SC-D/built-in_DLQ/native_DLQ+retry_middleware/`: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/built-in_DLQ/native_DLQ%2Bretry_middleware/publisher_log.txt), [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/built-in_DLQ/native_DLQ%2Bretry_middleware/consumer_log.txt) and the screenshot [`img.png`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/built-in_DLQ/native_DLQ%2Bretry_middleware/img.png) (commit `0b77d50`). Run on 11 Nov 2026.
- The consumer sets `NoRequeueOnNack = true` and declares the queue `SC-D.2` with the dead-letter arguments. `Retry` repeats the handler (`MaxRetries = 3`, initial interval 100 ms, multiplier 2, `MaxInterval` 2 s). When the retries are exhausted the router Nacks the message. Because the Nack is without requeue, the **broker** dead-letters it through `sc-d.dlx` to `sc-d.dead`.
- **Who does what:** the routing to the dead-letter queue is done by the broker. The attempt counter is `Retry` in the process, so it is in memory. The variant is therefore not fully broker-side (§5, items 1 and 2).
- **Attempts:** each of the 3 poison messages had 4 handler calls (12 `FAILED` lines in total), which is 1 initial attempt and 3 retries. The 6 normal messages were each received once.
- **Dead-letter queue:** the screenshot `img.png` (Management UI, 11 Nov 2026 11:11:52) shows `sc-d.dead` with 3 Ready messages, and its chart rises from 0 to 3 during the run. It shows the count, not the message ids.
- **Blocking:** while the retries run, the consumer takes no other message (about 2 s for each poison message). The Watermill subscriber processes messages one at a time (§5, item 4).
- Result: **3 of 3**.

#### Variant C: quorum queue with `x-delivery-limit`
Raw records in `SC-D/built-in_DLQ/quorum_queue/` (commit `0b77d50`): normal run [`publisher_log_normal.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/built-in_DLQ/quorum_queue/publisher_log_normal.txt) and [`normal_consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/built-in_DLQ/quorum_queue/normal_consumer_log.txt), crash run [`publisher_log_for_crash.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/built-in_DLQ/quorum_queue/publisher_log_for_crash.txt) and [`crash_consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/built-in_DLQ/quorum_queue/crash_consumer_log.txt), and the screenshot [`img.png`](https://github.com/ITA-Internship/GoJudge/blob/0b77d500b8abbffea359d9211f58021628cfe54f/spike/spk-03/RabbitMQ/SC-D/built-in_DLQ/quorum_queue/img.png). Run on 11 Nov 2026.
- **Queue settings:** `SC-D.quorum` with `x-queue-type=quorum`, `x-delivery-limit=3`, `x-dead-letter-exchange=sc-d.dlx` and `x-dead-letter-routing-key=sc-d.dead`. The handler Nacks a poison message and Acks any other.
- **Numeric headers:** a quorum queue adds numeric headers to each delivery (`x-acquired-count`, which grows with every delivery, and `x-delivery-count`, which appears only after a consumer is lost). The Watermill `DefaultMarshaler` expects strings, so the consumer uses `PreprocessDelivery` to convert them. Without it a redelivery breaks the unmarshalling.
- **Normal run** (`normal_consumer_log.txt`): Nack with requeue. The first poison message came back again and again: 51 `FAILED` lines in the log until the run was stopped with Ctrl+C. The delivery limit never dead-lettered it, and the messages behind it were not processed (8 messages stayed in `SC-D.quorum`). In RabbitMQ 4.4 the delivery limit did not count Nack with requeue.
- **Crash run** (`crash_consumer_log.txt`): the consumer exits with `os.Exit(1)` on a poison message before the Nack, and is started again by hand. Each poison message arrived **4 times** (the first delivery and 3 redeliveries, so a limit of 3 gives 4 deliveries) and then the broker dead-lettered it with the reason `delivery_limit` (read in the Management UI; not captured in the screenshots). 3 of 3 poison messages, and 6 of 6 normal messages were received. The crash version differs from the normal one only by the `os.Exit(1)` line (commented out in the committed code). To reproduce it, uncomment the line and restart the consumer manually, 4 starts per poison message.
- **Contents of `sc-d.dead`:** the screenshot `img.png` of this variant (Management UI, 13:11:12) shows `SC-D.quorum` with 8 Ready (the blocked messages), `sc-d.dead` with 10 Ready and `SC-D.2` empty. All 10 messages in `sc-d.dead` came from `SC-D.quorum`: 4 with the reason `delivery_limit` (3 were expected) and 6 with the reason `rejected`. The `rejected` ones are most likely from an earlier quorum attempt with `NoRequeueOnNack = true`, which is not logged (a Nack without requeue dead-letters at once, so the limit is not used). Why there are 4 `delivery_limit` messages and not 3 was not established. The broker state accumulated over several runs and the screenshots were not taken after each run, so this report relies on the logs.
- Result: the delivery limit works for the loss of a consumer, not for handler errors (§5, item 3).

*The screenshot was taken once, after all SC-D runs, so the counts in it are cumulative. `sc-d.dead` holds messages from several runs and variants (6 `rejected` and 4 `delivery_limit`), and `SC-D.quorum` holds the 8 messages left behind by the normal run. The screenshot does not show the state after each variant on its own.*

### SC-E Evidence: Transactional Outbox, Database and Broker Stay Consistent (added after review)
**Result: Verified.** All 8 scenarios gave the expected values in one run of `runall` on 11 Nov 2026. SC-E is not one of the planned criteria (SC-A to SC-D) and has no failure criterion in the spike plan.

**Question.** The review asked whether the outbox pattern from the Watermill documentation ([forwarder](https://watermill.io/advanced/forwarder/)) avoids two situations: data is persisted but no message is sent, and the reverse, a message is sent but the data is not persisted. SC-A to SC-D do not use a database, so this is the DEC-12 question (§6).

**How the outbox works in the experiment (the forwarder pattern from the documentation):**
- **Application side (`publisher`).** For each submission one PostgreSQL transaction writes the business row (`spk03_submissions`) and the event. The event is written by a Watermill SQL publisher that is bound to the transaction (`watermillSQL.NewPublisher(tx, …)`) and decorated with `forwarder.NewPublisher`. The decorator wraps the message in an envelope that carries the destination topic and writes it to the outbox table. Commit or rollback decides both writes together.
- **Relay (`forwarder`).** The Watermill forwarder (`forwarder.NewForwarder`) reads the outbox table with the SQL subscriber, unwraps the envelope, publishes the original message to the queue `outbox_submissions` and only then acknowledges the row (an offset in `watermill_offsets_outbox_forwarder`). The AMQP publisher uses `ConfirmDelivery = true`, so `Publish` returns after the broker has confirmed the message.
- **Checker (`consumer`).** It receives the messages and, after Ctrl+C, compares the delivered `submission_id` values with the rows in `spk03_submissions`. `data_without_message` counts rows with no message. `message_without_data` counts messages with no row. Both must be 0. Repeated deliveries are counted as `duplicates` and are allowed, because the outbox is at-least-once.
- **Fault injection.** In the failure scenarios submission #3 of 5 fails. A fault is `tx.Rollback()` (an application error) or `os.Exit(1)` (a killed process, as in SC-B). PostgreSQL rolls back an open transaction when the connection drops. The forwarder can be killed after it published a message and before it acknowledged the row (`-crash-after-publish`).
- **Baseline without the outbox (B1, B2).** The same application, but the event is published straight to RabbitMQ, after the commit (B1) or before it (B2). They show that the two situations from the review really happen without the outbox, and that the checker detects them.

**Scenarios** (`publisher -mode=…`; `Expected` is what the design predicts, `Observed` is the consumer's `SUMMARY` line from the run):

| ID | Scenario | Fault | Expected | Observed |
| :--- | :--- | :--- | :--- | :--- |
| E1 | `outbox`, forwarder running | none | 5 committed, 5 delivered, `data_without_message=0`, `message_without_data=0` | `deliveries=5 unique=5 duplicates=0 db_rows=5 data_without_message=0 message_without_data=0` |
| E2 | `outbox`, forwarder started after the publisher has finished | forwarder is down while the data is committed | after the publisher: `submissions_in_db=5`, `outbox_rows_added=5`, nothing received yet. After the forwarder starts: 5 delivered, 0 and 0 | `deliveries=5 unique=5 duplicates=0 db_rows=5 data_without_message=0 message_without_data=0` |
| E3 | `outbox-rollback` | #3: error after both writes, rollback | `submissions_in_db=4`, `outbox_rows_added=4`, 4 delivered, 0 and 0 | `deliveries=4 unique=4 duplicates=0 db_rows=4 data_without_message=0 message_without_data=0` |
| E4 | `outbox-crash-before-commit` | #3: process killed before `Commit` | `db_rows=2`, 2 delivered, 0 and 0 | `deliveries=2 unique=2 duplicates=0 db_rows=2 data_without_message=0 message_without_data=0` |
| E5 | `outbox-crash-after-commit` | #3: process killed right after `Commit` | `db_rows=3`, 3 delivered (#3 too, although the publisher is dead), 0 and 0 | `deliveries=3 unique=3 duplicates=0 db_rows=3 data_without_message=0 message_without_data=0` |
| E6 | `outbox`, forwarder with `-crash-after-publish`, then restarted | forwarder killed after the first publish, before the Ack | `deliveries=6`, `unique=5`, `duplicates=1` (the same submission delivered twice), 0 and 0 | `deliveries=6 unique=5 duplicates=1 db_rows=5 data_without_message=0 message_without_data=0` |
| B1 | `direct-commit-first` (no outbox) | #3: killed between `Commit` and `Publish` | `db_rows=3`, 2 delivered, `data_without_message=1` | `deliveries=2 unique=2 duplicates=0 db_rows=3 data_without_message=1 message_without_data=0` |
| B2 | `direct-publish-first` (no outbox) | #3: DB rolled back after the publish | `db_rows=4`, 5 delivered, `message_without_data=1` | `deliveries=5 unique=5 duplicates=0 db_rows=4 data_without_message=0 message_without_data=1` |

**How to run** (RabbitMQ and PostgreSQL must be reachable on `localhost`, and `.env` must hold the credentials). One command runs all eight scenarios one after another and writes all logs in one go:

`go run ./spike/spk-03/RabbitMQ/OUTBOX/runall`

For each scenario the runner purges the queue, starts the consumer, starts the forwarder when the scenario needs it (before the publisher; after it in E2; killed after the first publish and restarted in E6; not at all in B1 and B2), runs the publisher in the scenario's mode, and waits until the consumer has had no message for 10 s. It then compares the consumer's `SUMMARY` line with the expected values and prints PASS or FAIL per scenario. The publisher empties `spk03_submissions` at the start of each run. The programs can also be started by hand (consumer, forwarder, publisher `-mode=<mode>`); each has a `-log <file>` flag.

Raw records: every program writes its own log. For each scenario they are in `spike/spk-03/RabbitMQ/OUTBOX/<scenario>/` as `publisher_log.txt`, `forwarder_log.txt` (E6 also `forwarder_restart_log.txt`) and `consumer_log.txt`. The result table of the whole run is `spike/spk-03/RabbitMQ/OUTBOX/runall_log.txt` (commit `98d7603`).

**What was observed** (run of `runall` on 11 Nov 2026, all 8 scenarios as expected, see `runall_log.txt`):
- **Data persisted, no message (the first situation from the review):** it did not occur with the outbox. E5 (the publisher killed right after `Commit`) and E2 (the forwarder not running during the commits) both ended with `data_without_message=0`: the message was sent later by the forwarder. Without the outbox (B1) the same crash left 1 row without a message.
- **Message sent, no data (the reverse):** it did not occur with the outbox. E3 (rollback) and E4 (killed before `Commit`) ended with `message_without_data=0`, and the rows of #3 are in neither table. Without the outbox (B2) the rollback left 1 message without a row.
- **Relay crash:** in E6 the forwarder was killed after its first publish and before the Ack, and the restart sent that row again: 6 deliveries for 5 submissions (`duplicates=1`), nothing lost. The outbox is at-least-once (§5, item 9).
- In every scenario the publisher's own summary and the database agree with the consumer's `db_rows`.

**Method notes.**
- The business table `spk03_submissions` is a stand-in. There is no real submission logic, and the consumer only counts.
- Crashes are simulated by `os.Exit(1)` and a rollback by `tx.Rollback()` (§5, item 8).
- A first run of `runall` showed 8 FAIL (1 or 2 of 5 messages received). The cause was most likely a consumer and a forwarder from an earlier manual run that were still running and took messages and rows. `runall` now stops if it finds another consumer on the queue or another forwarder on the outbox. The results above are from the next run, with nothing else running.
- A first attempt to build SC-E wrote the events to the outbox table without `forwarder.NewPublisher`. The forwarder could not unwrap them (`unknown destination topic`). That attempt is not part of the evidence; it is the reason for §5, item 11.

## 4. Failure Criteria Check
- **F-08 (any message lost in SC-A, not 100 of 100):** **Not met.** 100 of 100 messages were delivered.
- **F-09 (redelivered duplicate produces two final states in SC-C, not 50 of 50):** **Not met.** 50 of 50 operations have exactly one final state.
- **F-10 (poison message cannot reach the dead-letter route in SC-D, not 3 of 3):** **Not met.** 3 of 3 reached the dead-letter queue.

No failure criterion of the plan (F-08 to F-10) was met.

SC-E has no failure criterion in the plan. Its own check is that `data_without_message` and `message_without_data` are both 0 in E1 to E6. Result: **met in all six scenarios** (not a failure).

## 5. Risks and Limitations
These do not change a criterion's status, but they bound what the evidence proves. They are inputs for CON-02 and DEC-12.
1. **Attempts are counted in the application.** For handler errors, the retry counter is `middleware.Retry` in process memory, in variants A and B. If a message kills the process, the counter restarts from zero. In variant A the message may then never reach the poison queue, and the application also makes an extra publish. In variant B the routing to the dead-letter queue is done by the broker, so B is not fully broker-side.
2. **Watermill cannot Nack individual messages differently.** `NoRequeueOnNack` applies to the whole queue (`nackMsg` in `watermill-amqp` v3.1.0, `pkg/amqp/subscriber.go`). A purely broker-side retry counter for handler errors is therefore not possible through Watermill. It would need the plain `amqp091` client (Nack with requeue up to a threshold, then without requeue). This was not tested.
3. **`x-delivery-limit` protects only against consumer loss.** In RabbitMQ 4.4 it did not count Nack with requeue (variant C, normal run). No separate limit on `x-acquired-count` was found. Retry combined with a quorum queue and `x-delivery-limit` as a safeguard against crash loops was not tested.
4. **A consumer is blocked while retries run.** The Watermill subscriber processes messages one at a time, so during the retries of one poison message (about 2 s with the settings used) nothing else is taken from that queue. Several consumers on one queue would remove this. That was not tested. A goroutine inside the handler is not a fix: the router acknowledges the message when the handler returns, so the `Ack` would be sent before the work is done.
5. **The dead-letter topology is created outside the code.** The exchange `sc-d.dlx`, the queue `sc-d.dead` and the binding were created by hand. Without them the broker silently drops dead-lettered messages. 
6. **The final state in SC-C is in memory.** The map of applied keys lives in the consumer process and is lost on restart. A real design needs a durable store (for example PostgreSQL). SC-E uses PostgreSQL only for the business row and the outbox. The idempotent consumer is still in memory, and a durable idempotency check combined with the database write was not tested.
7. **The idempotency check is safe only for one consumer handling one message at a time.** The "check, then apply" step is not atomic. With several consumers or concurrent handlers it needs an atomic operation.
8. **Failures are simulated.** A consumer kill is `os.Exit(1)` in code (SC-B and SC-D variant C), a lost `Ack` is a `Nack` after the effect (SC-C), and a poison message is a metadata flag that always fails (SC-D). In SC-E a crash is `os.Exit(1)` before or after `Commit` and in the forwarder after the publish, and an application error is `tx.Rollback()`. No process was killed from outside, and no network or broker failure was injected.
9. **The outbox is at-least-once, not exactly-once.** If the forwarder dies after the publish and before the row is acknowledged, the row is sent again (E6: a second delivery of the same submission). Consumers need the idempotency key from SC-C (`submission_id` in SC-E), and items 6 and 7 apply to it.
10. **Latency, load and clean-up are not measured.** The forwarder polls the outbox table (`PollInterval` and `ResendInterval` default to 1 s in `watermill-sql` v2.0.0), so a message can wait about a second. The SQL subscriber only moves an offset when a row is acknowledged and does not delete rows, so the outbox table grows and needs a clean-up job.
11. **A row the forwarder cannot unwrap blocks the outbox.** By default the forwarder Nacks it (`AckWhenCannotUnwrap` is false in `forwarder.go`), the SQL subscriber sends it again every second, and the rows behind it are not forwarded. This was seen once during development (see the SC-E method notes). Only messages written through `forwarder.NewPublisher` may be put in the outbox table.
12. **The AMQP publisher needs `ConfirmDelivery`, and the queue must exist.** `amqp.NewDurableQueueConfig` leaves `ConfirmDelivery` off (`watermill-amqp` v3.1.0, `config.go`), so without it the forwarder could acknowledge a row that the broker never stored. SC-E sets it to true. A message for a queue that does not exist is dropped by the broker (default exchange, `Mandatory` is not set), which is why the consumer is started first. Neither point was tested.
13. **Not tested:** RabbitMQ unavailable while the forwarder runs (a failed publish should leave the row unacknowledged and retried, but the reconnect behaviour of `watermill-amqp` was not checked), several forwarder instances on one outbox, and the order of events across transactions. The SQL subscriber reads only rows of transactions older than the oldest open one (`pg_snapshot_xmin` in `DefaultPostgreSQLSchema`, read from the source), so a long-running transaction holds back all later rows. This was not tested.

## 6. Outcome & Follow-On Decisions
- **Outcome:** **Go.** SC-A, SC-B, SC-C and SC-D are all Verified, and no failure criterion was met. The Go stands with the risks in §5. In particular, a poison message reaches the dead-letter queue, but for handler errors the retry count is kept in the application and not in the broker (§5, items 1 to 4). These risks are named for CON-02. SC-E was added after the review and is input for DEC-12. It does not change SC-A to SC-D.
- **DEC-04 (messaging with RabbitMQ and Watermill):** recommendation: approve RabbitMQ with Watermill for the asynchronous submission pipeline, with the limits in §5 accepted as constraints on the design, not as reasons to change the choice.
- **DEC-12 (dispatch consistency between the database write and the publish):** tested by SC-E with the Watermill transactional outbox (forwarder), PostgreSQL and RabbitMQ. Result: in E1 to E6 `data_without_message=0` and `message_without_data=0`, while the baselines B1 and B2 without the outbox showed both problems. Recommendation: use the transactional outbox (Watermill forwarder) for the dispatch, with the constraints that delivery is at-least-once and consumers need an idempotency key (§5, items 9 to 13). The failures were simulated (§5, item 8).
