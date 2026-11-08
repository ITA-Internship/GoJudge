# SPK-03: RabbitMQ and Watermill: Delivery, Redelivery, Idempotency and Dead-Letter Routing

## Metadata
- **Spike ID:** SPK-03
- **Type:** Technical spike
- **Owner:** Pair C
- **Timebox:** 4 pair-days (Days 3–6: Wed 4 Nov – Mon 9 Nov). Report due end of day 6.
- **Status:** Completed within timebox
- **Dependencies:** S1-01 (INF-03: local RabbitMQ), DEC-04 (provisional)
- **Experiment code:** branch `spike/SPK-03-rabbitmq-watermill`, directory `spike/spk-03/RabbitMQ/`. All evidence links below point to commit `5e5509d`.
- **Libraries:** `github.com/ThreeDotsLabs/watermill` v1.5.3, `github.com/ThreeDotsLabs/watermill-amqp/v3` v3.1.0 (AMQP client `amqp091-go` v1.11.0), Go 1.27.3.

## 1. Objective & Hypothesis
- **Objective:** Answer one question: can RabbitMQ with Watermill publish and consume submission messages without loss, redeliver after a consumer failure, keep one final state under duplicate delivery, and route poison messages to a dead-letter route?
- **Hypothesis:** RabbitMQ with Watermill delivers every message, redelivers after a consumer is killed mid-processing, produces one final state under duplicate delivery when an idempotency key is used, and routes poison messages to a dead-letter route.

## 2. Experiment Setup
- **Broker:** the local RabbitMQ from S1-01, reached on `localhost:5672`.
- **Topology:** Watermill AMQP with `amqp.NewDurableQueueConfig`. A Watermill topic is a durable RabbitMQ queue of the same name. Each criterion uses its own topic (`hello`, `SC-B`, `SC-C`, `SC-D` and `SC-D.poison`).
- **Run order:** the queue is declared by the subscriber, not by the publisher, so the consumer is started before a run that publishes messages.
- **Code layout:** one `publisher` and one `consumer` program per criterion in `spike/spk-03/RabbitMQ/SC-A` … `SC-D`. The shared package `rmqenv` builds the AMQP URI from the repository-root `.env`.
- **Secrets & Environment Policy:** no secrets or credentials are committed. `RABBITMQ_USER` and `RABBITMQ_PASSWORD` are read from the git-ignored `.env` file, and the raw logs contain no credentials.
- **Raw records:** the console output of each run is stored beside its code as `publisher_log.txt` and `consumer_log.txt`.

## 3. Success Criteria & Recorded Evidence

| Criterion | Status | Evidence Summary |
| :--- | :--- | :--- |
| **SC-A** | **Verified** | 100 of 100 messages published and consumed. The set of 100 message UUIDs is identical on both sides, with no loss and no duplicates. |
| **SC-B** | **Verified** | 5 of 5 consumer kills during processing were followed by redelivery of the same message (same UUID). |
| **SC-C** | **Verified** | 50 of 50 logical operations ended in exactly one applied state, although 100 duplicate deliveries arrived. |
| **SC-D** | **Verified** | 3 of 3 poison messages reached the dead-letter queue after the configured attempts. |

### SC-A Evidence: Publish and Consume (100 of 100)
Raw records: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/5e5509dbac8945b7ec217db6d760d38f76ad67e1/spike/spk-03/RabbitMQ/SC-A/publisher_log.txt) and [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/5e5509dbac8945b7ec217db6d760d38f76ad67e1/spike/spk-03/RabbitMQ/SC-A/consumer_log.txt) (branch `spike/SPK-03-rabbitmq-watermill`, commit `5e5509d`). Run on 7 Nov 2026.
- The publisher sent 100 messages (`message #1` … `message #100`), each with its own UUID. The consumer received 100 messages and acknowledged each one.
- The two sets of UUIDs were compared one by one: 100 unique on each side, identical, none missing and none repeated. Payloads `#1` … `#100` are all present. The order is the order of publication.
- Both logs contain no `ERROR` or `WARN` lines, and both processes closed their connections cleanly.
- The consumer was started after the publisher had finished (consumer connected 13:25:07, publisher finished 13:24:49), so the messages waited in the durable queue until a consumer arrived.

### SC-B Evidence: Redelivery After a Consumer Is Killed (5 of 5)
Raw records: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/5e5509dbac8945b7ec217db6d760d38f76ad67e1/spike/spk-03/RabbitMQ/SC-B/publisher_log.txt) and [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/5e5509dbac8945b7ec217db6d760d38f76ad67e1/spike/spk-03/RabbitMQ/SC-B/consumer_log.txt) (commit `5e5509d`). Run on 7 Nov 2026.
- A new topic (`SC-B`, a new queue) was used, and one message was published: `id=185ef076-68d1-4d9f-aba5-5a26842354cf`.
- The consumer was started six times in a row. Each of the first five runs received the message and was terminated with `os.Exit(1)` before `Ack`. `os.Exit` skips the deferred `Close`, so the connection is dropped like a killed process.
- Every time, the broker delivered the **same message** (same UUID) to the next consumer.
- Each of the 5 kills was followed by a redelivery of the same message: **5 of 5**.
- Method notes: the failure is simulated by `os.Exit(1)` in the consumer code, not by an external kill. The same message was killed and redelivered five times in sequence, rather than five separate messages being published. For the sixth run `os.Exit` was commented out so the message could be acknowledged.

### SC-C Evidence: Duplicates and One Final State (50 of 50)
Raw records: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/5e5509dbac8945b7ec217db6d760d38f76ad67e1/spike/spk-03/RabbitMQ/SC-C/publisher_log.txt) and [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/5e5509dbac8945b7ec217db6d760d38f76ad67e1/spike/spk-03/RabbitMQ/SC-C/consumer_log.txt) (commit `5e5509d`). Run on 8 Nov 2026.

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
Raw records: [`publisher_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/5e5509dbac8945b7ec217db6d760d38f76ad67e1/spike/spk-03/RabbitMQ/SC-D/publisher_log.txt) and [`consumer_log.txt`](https://github.com/ITA-Internship/GoJudge/blob/5e5509dbac8945b7ec217db6d760d38f76ad67e1/spike/spk-03/RabbitMQ/SC-D/consumer_log.txt) (commit `5e5509d`). Run on 8 Nov 2026.
- **Dead-letter route:** the Watermill router uses the `PoisonQueue` middleware, which republishes a failed message to the queue `SC-D.poison` once the handler has given up. The `Retry` middleware sits inside it (`MaxRetries = 3`, initial interval 100 ms), so `PoisonQueue` sees the error only after the retries are exhausted. A separate subscriber on `SC-D.poison` reads the dead-lettered messages and prints the failure reason.
- **Poison messages:** nine messages were published in three rounds (normal, poison, normal). The three poison messages carry the metadata flag `poison=true`, and the handler always returns an error for them.
- **Configured attempts:** `MaxRetries = 3` means one initial attempt and three retries, which is **4 handler calls** per poison message.

| Poison message | Handler attempts (`FAILED`) | Appeared in `SC-D.poison` | Reason |
| :--- | :--- | :--- | :--- |
| `eaf62665…` | 4 | once | `cannot process message` |
| `348a40b3…` | 4 | once | `cannot process message` |
| `094db0ee…` | 4 | once | `cannot process message` |

- After reaching the poison queue a message was not processed again, and the following messages were handled normally, so the main queue was not blocked.
- The 6 normal messages were each received exactly once.
- Result: **3 of 3** poison messages reached the dead-letter route.

## 4. Failure Criteria Check
- **F-08 (any message lost in SC-A, not 100 of 100):** **Not met.** 100 of 100 messages were delivered.
- **F-09 (redelivered duplicate produces two final states in SC-C, not 50 of 50):** **Not met.** 50 of 50 operations have exactly one final state.
- **F-10 (poison message cannot reach the dead-letter route in SC-D, not 3 of 3):** **Not met.** 3 of 3 reached the dead-letter queue.

No failure criterion was met.

## 5. Risks and Limitations
These do not change a criterion's status, but they bound what the evidence proves. They are inputs for CON-02 and DEC-12.
1. **Dead-letter route is an application-level mechanism.** It is Watermill's `PoisonQueue` middleware publishing to a separate queue. RabbitMQ's native dead-letter exchange (`x-dead-letter-exchange`) was not tested.
2. **The final state in SC-C is in memory.** The map of applied keys lives in the consumer process and is lost on restart. A real design needs a durable store (for example PostgreSQL).
3. **The idempotency check is safe only for one consumer handling one message at a time.** The "check, then apply" step is not atomic. With several consumers or concurrent handlers it needs an atomic operation.
4. **Failures are simulated.** A consumer kill is `os.Exit(1)` in code (SC-B), a lost `Ack` is a `Nack` after the effect (SC-C), and a poison message is a metadata flag that always fails (SC-D).

## 6. Outcome & Follow-On Decisions
- **Outcome:** **Go.** SC-A, SC-B, SC-C and SC-D are all Verified, and no failure criterion was met. The risks in §5 are named for CON-02.
- **DEC-04 (messaging with RabbitMQ and Watermill):** recommendation: approve RabbitMQ with Watermill for the asynchronous submission pipeline. 
- **DEC-12 (dispatch consistency between the database write and the publish):** **not tested by this spike.** This report only confirms that publishing and acknowledgement work.
