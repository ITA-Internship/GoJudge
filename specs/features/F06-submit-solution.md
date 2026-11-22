# Feature Spec: F06 – Submit solution

---

## 1. Feature Metadata *(Required)*

| Field            | Value                                                            |
|------------------|------------------------------------------------------------------|
| Feature ID       | F06                                                              |
| Title            | Submit solution                                                  |
| Status           | In Review                                                        |
| Responsible pair | Pair A                                                           |
| Reviewer pair    | Pair B, Pair D                                                   |
| GitHub Issue     | [issue #29](https://github.com/ITA-Internship/GoJudge/issues/29) |
| Last updated     | 2026-11-22                                                       |

---

## 2. Overview and User Story *(Required)*

**Purpose:** Allow an authenticated learner to submit a Go solution for a published problem, persist it in the `QUEUED` state, and dispatch it for asynchronous evaluation without blocking the API request on judge execution.

**User story:**

> As a **learner**,
> I want to submit my Go solution for a problem and receive a submission ID and current processing state,
> so that my code is evaluated asynchronously against hidden tests without keeping my request open.

**Expected value:** Decouples HTTP submission intake from isolated Kubernetes judge execution (`F07`), gives learners immediate confirmation and status polling (`F06-03`, `F06-04`), and persists the submission records that feed submission history (`F08`).

---

## 3. Scope *(Required)*

**In scope:**

- Validating and accepting a single-file Go source submission from an authenticated user for a published problem (`F06-01`).
- Persisting the submission with initial state `QUEUED` and `verdict: null` in accordance with [CON-01](../../docs/contracts/CON-01.md).
- Dispatching an evaluation request to the message broker in accordance with `CON-02` (`F06-02`).
- Providing a submission state lookup endpoint for the owning user (`F06-03`).
- Defining client submit and regular short-polling behavior in the editor workspace (`F06-04`).

**Out of scope:**

- Implementation code for `F06-01` to `F06-04` (starts only after `SPEC-F06` is approved).
- Running public sample tests without creating a persisted submission (`F05`).
- Compiling, isolating, or executing learner-submitted Go code in Kubernetes Jobs (`F07`).
- Listing and filtering a learner's submission history (`F08`).
- Deciding `DEC-12` (dispatch consistency mechanism between database write and message publish), `CON-02` messaging details, `CON-04` problem packaging, or `CON-05` database schema.

---

## 4. Functional Requirements *(Required)*

Each requirement must be numbered and observable, so that a tester can confirm it.

| ID    | Requirement                                                                                                                                                                                                                                                                             |
|-------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| FR-01 | The system shall expose an authenticated endpoint for submitting a Go source solution for a specified problem (`F06-01`).                                                                                                                                                               |
| FR-02 | The system shall identify the submitter (`user_id`) exclusively from the authenticated caller context provided by `F01`, ignoring any user identifier supplied in the request payload.                                                                                                  |
| FR-03 | The system shall verify that the target problem exists and is published before creating a submission, and shall respond `404 Not Found` otherwise.                                                                                                                                      |
| FR-04 | The system shall reject submissions where `source_code` is empty or whitespace-only (resolving the `is_blank` handoff from `F04`), contains the NUL character (`U+0000`), or exceeds the maximum source size in bytes (OQ-04), without persisting a submission or publishing a message. |
| FR-05 | Upon valid input, the system shall persist a submission record in the `QUEUED` state with `verdict` set to `null`, using the exact state and verdict names from `CON-01`.                                                                                                               |
| FR-06 | The system shall dispatch an asynchronous evaluation message for each created submission according to `CON-02` and the dispatch consistency mechanism chosen under `DEC-12` (OQ-03).                                                                                                    |
| FR-07 | The submission creation endpoint shall return immediately after intake without waiting for judge execution, returning HTTP `202 Accepted` with the `id`, `problem_id`, initial state `QUEUED`, `verdict: null`, `created_at`, and `updated_at` timestamps.                              |
| FR-08 | The system shall expose an authenticated endpoint to retrieve the current state and verdict of a single submission by `id` (`F06-03`), returning a `CON-01` verdict when `COMPLETED` or `FAILED`, and `verdict: null` when `QUEUED`, `RUNNING`, or `CANCELLED`. |
| FR-09 | The system shall restrict submission state retrieval to the owning user (`submission.user_id == caller.user_id`) and shall respond `404 Not Found` when a submission is owned by another user.                                                                                          |
| FR-10 | Handling of a new submission while the same user already has an active (`QUEUED` or `RUNNING`) submission for the same problem shall follow the `CON-01` duplicate-submission policy (currently open; recorded as OQ-01).                                                               |
| FR-11 | In the web editor (`F06-04`), triggering Submit shall send the current editor `source_code`, prevent duplicate accidental clicks while the creation request is in flight, and use regular short polling (bounded interval, stopping on terminal state) against the state endpoint until `COMPLETED`, `FAILED`, or `CANCELLED` is reached. |

---

## 5. Acceptance Criteria *(Required)*

Write each criterion as Given / When / Then. Number them AC-01, AC-02, and so on. Every criterion must be testable.

### AC-01 – Successful asynchronous submission

- **Given** an authenticated learner and an existing published problem
- **When** the learner submits valid non-blank Go `source_code` within the allowed byte limit
- **Then** the system persists a submission with state `QUEUED` and `verdict: null`, triggers message dispatch per `CON-02`, and responds immediately with `202 Accepted` containing `id`, `problem_id`, `state: "QUEUED"`, `verdict: null`, `created_at`, and `updated_at`.

### AC-02 – Retrieve own submission state

- **Given** an authenticated learner owns a submission in state `QUEUED`, `RUNNING`, `COMPLETED`, `FAILED`, or `CANCELLED`
- **When** the learner requests the status of that submission by `id`
- **Then** the system responds `200 OK` with the current `CON-01` state, the `CON-01` test or system verdict when `COMPLETED` or `FAILED` respectively, or `verdict: null` when `QUEUED`, `RUNNING`, or `CANCELLED`, and returns no hidden test details.

### AC-03 – Unauthenticated request rejected

- **Given** a caller without a valid Bearer access token
- **When** the caller attempts to create a submission or query a submission's state
- **Then** the system responds `401 Unauthorized` per `F01`, creates no submission record, and dispatches no queue message.

### AC-04 – Invalid source code rejected

- **Given** an authenticated learner and a published problem
- **When** the learner submits empty or whitespace-only `source_code`, `source_code` containing `U+0000`, or `source_code` exceeding the maximum byte limit
- **Then** the system rejects the request with `400 Bad Request` (or `413 Content Too Large` for oversized payloads), creates no submission record, and dispatches no queue message.

### AC-05 – Non-existent or unpublished problem rejected

- **Given** an authenticated learner and a `problem_id` that does not exist or is not published
- **When** the learner submits a solution for that `problem_id`
- **Then** the system responds `404 Not Found`, creates no submission record, and dispatches no queue message.

### AC-06 – Cross-user access to submission state denied

- **Given** learner A owns a submission with a known `id`
- **When** learner B requests the state of learner A's submission by `id`
- **Then** the system responds `404 Not Found` and discloses no submission data.

### AC-07 – Duplicate active submission handling (open question)

- **Given** an authenticated learner already has a submission in `QUEUED` or `RUNNING` state for a problem
- **When** the learner submits another solution for the same problem
- **Then** the system follows the resolved `CON-01` duplicate-submission rule (default in `CON-01`: return conflict error `409 Conflict`; tracked as OQ-01 until Pair D records the final decision).

### AC-08 – Editor submit action and state polling

- **Given** a signed-in learner is editing a non-blank solution in the browser editor (`F04`)
- **When** the learner triggers Submit
- **Then** the UI disables repeat submit clicks until the creation request settles, displays `QUEUED` upon `202 Accepted`, polls the state endpoint until `state` reaches `COMPLETED`, `FAILED`, or `CANCELLED`, and displays the resulting state and verdict.

---

## 6. Technical Design *(Required)*

**Component responsibilities:**

| Component                                            | Responsibility                                                                                                                                                                               |
|------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Web Application (`F04` + `F06-04`)**               | Reads `source_code` and `is_blank` from the browser editor (`F04`), sends the submit request, prevents double-submission in the UI, and polls `GET` submission state until a terminal state. |
| **Go API – Auth & Role Middleware (`F01`)**          | Validates the OIDC Bearer token, resolves the caller's internal `user_id` and role (`learner` or `administrator`), and enforces deny-by-default access control.                              |
| **Go API – Submission Handler (`F06-01`, `F06-03`)** | Validates `problem_id` and `source_code`, checks problem availability, creates the `QUEUED` submission, initiates dispatch, and serves owner-only state lookups.                             |
| **Submission Dispatcher (`F06-02`)**                 | Publishes the evaluation message to RabbitMQ using Watermill according to `CON-02` and the consistency mechanism selected in `DEC-12`.                                                       |
| **PostgreSQL (`CON-05`)**                            | Stores submission records (`submissions`) and, if selected by `DEC-12`, transactional outbox records (`SPK-03` SC-E).                                                                        |
| **Judge Orchestrator (`F07`)**                       | Consumes the `CON-02` message asynchronously and transitions submission state and verdict per `CON-01` (outside `F06`).                                                                      |

**Dependencies:**

- **Accounts and access (`F01`):** Provides `user_id` (`users.id`), authentication middleware, role enforcement, and the standard API error envelope.
- **Problem catalogue / details (`F02` / `F03`):** Provides problem existence and published status checks.
- **Browser code editor (`F04`):** Supplies `source_code` and `is_blank` state (`FR-04`, `AC-17`).
- **Contracts (`CON-01`, `CON-02`, `CON-05`):** `CON-01` defines states and verdicts; `CON-02` (`#23`, Pair B, draft in PR `#38`) defines ID-only messaging (`source_code` stays in PostgreSQL); `CON-05` (`#34`, Pair D, preliminary) defines the `submissions` persistence model.
- **Structured logging (`INF-07`):** Provides request correlation (`X-Request-ID`) and structured logs; must never log `source_code` or credentials.

**Architecture decisions:**

- **DEC-01 (Platform stack – Go, PostgreSQL, RabbitMQ):** Provisional.
- **DEC-04 (Messaging with RabbitMQ and Watermill):** Provisional (`SPK-03` Go, pending MA-26).
- **DEC-12 (Dispatch consistency between DB write and message publish):** Open (`SPK-03` SC-E recommends transactional outbox; decision due before `F06-02`). No new architectural decision is introduced by this specification.

---

## 7. Contracts *(Conditional)*

### 7.1 API Endpoints *(if applicable)*

| Method | Path                                       | Purpose                                                         | Auth required                                |
|--------|--------------------------------------------|-----------------------------------------------------------------|----------------------------------------------|
| `POST` | `/api/v1/problems/{problemId}/submissions` | Create an asynchronous submission (`F06-01`)                    | Yes (`learner`, `administrator`)             |
| `GET`  | `/api/v1/submissions/{submissionId}`       | Retrieve current state and verdict of own submission (`F06-03`) | Yes (`learner`, `administrator`, owner only) |

**Request schema (`POST /api/v1/problems/{problemId}/submissions`):**

- Path parameter: `problemId` (UUID).
- Header: `Authorization: Bearer <access_token>` (caller `user_id` is derived from auth context).
- JSON body:

```json
{
  "source_code": "package main\n\nfunc solve() {}\n"
}
```

**Response schema:**

- `POST` success (`202 Accepted`) and `GET` success (`200 OK`):

```json
{
  "id": "b7e2c1a0-4f8d-4c1b-9e3a-2d1f0a9b8c7d",
  "problem_id": "11111111-2222-3333-4444-555555555555",
  "state": "QUEUED",
  "verdict": null,
  "created_at": "2026-11-17T10:00:00Z",
  "updated_at": "2026-11-17T10:00:00Z"
}
```

- Error envelope (aligned with `F01 §7.1`):

```json
{
  "error": {
    "code": "validation_error",
    "message": "Source code must not be empty."
  }
}
```

| Status | `code`                   | When                                                                                   |
|--------|--------------------------|----------------------------------------------------------------------------------------|
| `400`  | `validation_error`       | Malformed JSON, invalid UUID, blank/whitespace-only `source_code`, or `U+0000` present |
| `401`  | `unauthorized`           | Missing or invalid access token (`F01`)                                                |
| `404`  | `not_found`              | Problem not found / unpublished, or submission not found or not owned by caller        |
| `409`  | `submission_in_progress` | Duplicate active submission if `CON-01` default conflict rule is confirmed (OQ-01)     |
| `413`  | `payload_too_large`      | `source_code` exceeds the maximum allowed byte size (OQ-04)                            |
| `503`  | `service_unavailable`    | Persistence or dispatch dependency unavailable                                         |

### 7.2 Data Models *(if applicable)*

Conceptual `Submission` fields used by `F06` (final PostgreSQL column types, constraints, and indexes are owned by `CON-05` in `#34`):

- `id` (UUID, PK): Generated submission identifier.
- `user_id` (UUID, FK): Local user identifier from `F01` (`users.id`), taken from auth context.
- `problem_id` (UUID, FK): Target published problem identifier.
- `source_code` (text): Untrusted Go source text submitted by the learner (size-bounded per OQ-04).
- `state` (`CON-01` lifecycle state; exact PostgreSQL type — e.g., `ENUM` vs `text` + `CHECK` — owned by `CON-05`): `QUEUED`, `RUNNING`, `COMPLETED`, `FAILED`, `CANCELLED`.
- `verdict` (`CON-01` verdict, nullable; exact PostgreSQL type owned by `CON-05`): must be `null` when `state` is `QUEUED`, `RUNNING`, or `CANCELLED`, and populated only when `state` is `COMPLETED` or `FAILED`.
- `created_at` / `updated_at` (timestamptz): Creation and last state-update timestamps.

**Expected access patterns (input to `CON-05`):** lookup and atomic state update by `id`; checking for an active (`QUEUED` or `RUNNING`) submission by `(user_id, problem_id)` (or by `user_id`, depending on `CON-01`'s final scope for OQ-01); and listing a learner's own submissions ordered by `created_at` (aligned with `F08`).

### 7.3 Messaging Contracts *(if applicable)*

`F06` (`F06-02`) acts exclusively as a publisher under `CON-02` (`#23`, PR `#38`):

- **Publish flow & target:** Publishes one message per submission to the default exchange with routing key `submissions.dispatch`, with no reply message (`CON-02 §1–§2`).
- **Payload & metadata (IDs only):** Sends `{ "schema_version": 1, "submission_id", "problem_id", "user_id", "submitted_at" }` (where `submission_id` maps to `Submission.id`, `user_id` is `F01` `users.id` UUID, and `submitted_at` maps to `created_at`) with metadata `idempotency_key = submission_id`; `source_code` and tests remain in PostgreSQL and are never sent through the broker (`CON-02 §3`).
- **Publisher reliability:** Declares `submissions.dispatch` (with `x-dead-letter-exchange = submissions.dlx`), `submissions.dlx` (fanout exchange), `submissions.dead`, and their binding at startup before the first publish, and enables `ConfirmDelivery` (noting `ConfirmDelivery` confirms broker acceptance only, not queue presence; `CON-02 §2, §9`). DB-to-broker dispatch consistency is governed by `DEC-12` (`CON-02 §8`, OQ-03).

Consumer-side (`F07-01`) idempotency, acknowledgement after PostgreSQL writes, infrastructure retry rules, and dead-letter routing are defined in `CON-02 §4–§7` (with `SPK-03` risks in `CON-02 §9`) and remain subject to `DEC-04` and `DEC-12` approval (OQ-02, OQ-03).

### 7.4 State Transitions *(if applicable)*

All states and verdicts match [CON-01](../../docs/contracts/CON-01.md) verbatim:

- **Lifecycle states:** `QUEUED`, `RUNNING`, `COMPLETED`, `FAILED`, `CANCELLED`.
- **Test verdicts (only when `state == COMPLETED`):** `ACCEPTED`, `WRONG_ANSWER`, `COMPILATION_ERROR`, `RUNTIME_ERROR`, `TIME_LIMIT_EXCEEDED`, `MEMORY_LIMIT_EXCEEDED`, `OUTPUT_LIMIT_EXCEEDED`.
- **System verdict (only when `state == FAILED`):** `INTERNAL_ERROR`.
- **Cancelled state (`state == CANCELLED`):** Terminal state with no verdict (`verdict` is always `null`).
- **F06 boundaries & invariants:**
  - `F06` only creates submissions in `QUEUED` with `verdict: null`.
  - Subsequent transitions (`QUEUED -> RUNNING -> COMPLETED | FAILED`, `QUEUED | RUNNING -> CANCELLED`, and infrastructure retry `RUNNING -> QUEUED`) are driven outside `F06` by the orchestrator (`F07`) and follow `CON-01`'s "first terminal write wins" rule.
  - `F06` (`F06-03`) exposes these states and verdicts as persisted, without inventing new states or mutating terminal verdicts.

---

## 8. Edge Cases and Failure Handling *(Conditional)*

- [x] Invalid or missing input
- [x] Timeouts
- [x] Retries and their limits
- [x] Duplicate processing
- [x] Unavailable dependencies

| Scenario                                                                          | Expected behaviour                                                                                                                                                                                           |
|-----------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Empty or whitespace-only `source_code` (`F04` `is_blank`)**                     | Rejected with `400 Bad Request` (`validation_error`). No submission is persisted and no message is published (resolves `F04` `FR-04` / `AC-17`).                                                             |
| **`source_code` contains NUL byte (`U+0000`) or malformed JSON**                  | Rejected with `400 Bad Request` (`validation_error`). No submission is persisted and no message is published.                                                                                                |
| **`source_code` exceeds `max_source_size` bytes**                                 | Rejected with `413 Content Too Large` (or `400`, OQ-04). Never truncated silently; no submission is persisted or queued.                                                                                     |
| **Non-existent or unpublished (`draft`) `problemId`**                             | Rejected with `404 Not Found` (`not_found`). No submission is persisted or queued.                                                                                                                           |
| **Client sends `user_id`, `state`, or `verdict` in request body**                 | Ignored or rejected by strict schema parsing. `user_id` always comes from auth context (`F01`), and initial `state` is always `QUEUED` with `verdict: null`.                                                 |
| **Duplicate active submission (`QUEUED` or `RUNNING`) for same user and problem** | Behaviour follows the `CON-01` duplicate-submission policy (default: return `409 Conflict` / `submission_in_progress`; tracked as OQ-01 until Pair D confirms).                                              |
| **Rapid double-click on Submit in UI (`F06-04`)**                                 | The editor disables Submit while the `POST` request is in flight (`FR-11`, `AC-08`) to avoid accidental duplicate submissions.                                                                               |
| **PostgreSQL unavailable during submission creation**                             | Returns `503 Service Unavailable` (`service_unavailable`). No partial submission is created and no message is published to RabbitMQ.                                                                         |
| **Message broker unavailable after DB write (or publish timeout)**                | Handled according to `DEC-12` (`CON-02 §8`, OQ-03): e.g., transactional outbox (`SPK-03` SC-E) persists the outbox event in the same DB transaction for asynchronous relay; publisher declares topology at startup and enables `ConfirmDelivery` per `CON-02 §2, §9` (which confirms broker acceptance only). |
| **Broker redelivery / duplicate queue delivery (`F06-02` / `F07`)**               | Governed by `CON-02 §3–§4` (at-least-once delivery, `idempotency_key = submission_id`, skipping terminal states) and `CON-01` ("first terminal write wins") so duplicate deliveries do not corrupt final state or verdict.                                                                                |
| **State polling timeout or transient network error (`F06-04`)**                   | Client stops/backs off polling per bounded retry rules and surfaces a recoverable error without clearing editor content.                                                                                     |

---

## 9. Security Considerations *(Required)*

- [x] **Authentication:** Every `F06` endpoint requires a valid RS256 OIDC Bearer access token validated by `F01` middleware. Missing or invalid tokens return `401 Unauthorized`.
- [x] **Authorisation:** Follows `F01` deny-by-default routing. `user_id` is resolved strictly from the verified `F01` caller context, never from client input. `GET /api/v1/submissions/{submissionId}` is restricted to the owning user (`submission.user_id == caller.user_id`), preventing cross-user data access (IDOR).
- [x] **Input validation:** Learner-submitted `source_code` is **untrusted** input. The API validates `problemId` / `submissionId` format (UUID), rejects blank/whitespace-only code, rejects `U+0000` NUL bytes, and enforces a strict `max_source_size` byte limit (OQ-04) before persistence or dispatch.
- [x] **Sensitive data:** `Authorization` headers, tokens, and raw learner `source_code` are never written to structured logs (`INF-07`). Per `CON-02 §3`, raw `source_code` is also never transmitted through RabbitMQ (only IDs are queued; the judge loads source code from PostgreSQL). Responses from `F06` never return hidden test cases, hidden test outputs, or internal stack traces (`CON-04`).
- [x] **Execution isolation:** **Learner submissions are untrusted code.** `F06` never compiles, executes, or evaluates submitted Go source code in the API process or host container; `F06` only validates, stores, and dispatches the submission for isolated execution by `F07` (subject to `DEC-08` and `SPK-01` / `SPK-05` security boundaries).

---

## 10. Verification Plan *(Required)*

| AC    | Verification type  | Test or scenario                                                                                                                                                                                                        | Automated? |
|-------|--------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|
| AC-01 | Integration        | Authenticated learner submits valid Go code for a published problem; assert `202 Accepted`, DB row in `QUEUED` with `verdict: null`, and message queued per `CON-02`.                                                   | Yes        |
| AC-02 | Integration        | Owner queries `GET /api/v1/submissions/{submissionId}` across `CON-01` states (`QUEUED`, `RUNNING`, `COMPLETED`, `FAILED`, `CANCELLED`); assert `200 OK`, matching state/verdict invariant (`verdict: null` for `QUEUED`, `RUNNING`, `CANCELLED`), and no hidden test leakage. | Yes        |
| AC-03 | Integration        | Send `POST` and `GET` requests without a token or with an invalid token; assert `401 Unauthorized` and zero DB/queue side effects.                                                                                      | Yes        |
| AC-04 | Unit / Integration | Submit empty, whitespace-only, NUL-containing (`U+0000`), and oversized `source_code`; assert `400`/`413` rejection and zero DB/queue side effects.                                                                     | Yes        |
| AC-05 | Integration        | Submit to a non-existent UUID and to an unpublished (`draft`) problem; assert `404 Not Found` and zero DB/queue side effects.                                                                                           | Yes        |
| AC-06 | Integration        | Learner B calls `GET /api/v1/submissions/{submissionId}` for a submission owned by Learner A; assert `404 Not Found` with no submission details disclosed.                                                              | Yes        |
| AC-07 | Integration        | Submit a second solution while an active `QUEUED` or `RUNNING` submission exists for the same user and problem; verify behaviour matches resolved `CON-01` rule (OQ-01).                                                | Yes        |
| AC-08 | End-to-end         | In the browser workspace (`F04`), click Submit on valid code; verify button disables in-flight, state transitions display through polling, and polling stops on terminal state.                                         | Yes        |

---

## 11. Dependencies and Open Questions *(Required)*

**Related features:**

- **F01 – Accounts and access:** Authentication middleware, role enforcement, `users.id` foreign key, and standard error envelope.
- **F02 / F03 – Problem catalogue & details:** Published problem existence and availability check.
- **F04 – Browser code editor:** Supplies `source_code` and `is_blank`; `F06` resolves `F04` OQ-10 by rejecting blank submissions.
- **F07 – Automated judge:** Consumes queued submissions (`CON-02`), executes untrusted code in isolation, and transitions states/verdicts (`CON-01`).
- **F08 – Submission history:** Reads the submission records created by `F06`.

**Open questions:**

| #     | Question                                                                                                                                                                                                   | Owner                         | Status |
|-------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------|--------|
| OQ-01 | **Duplicate active submissions (`CON-01`):** Should submitting while an earlier submission is `QUEUED` or `RUNNING` return `409 Conflict` (CON-01 default) or be queued?                                   | Pair D                        | Open   |
| OQ-02 | **CON-02 finalisation (`#23`):** Final agreement from dependent pairs and formal approval of `CON-02` (`DEC-04` / `MA-26`). | Pair B / Mentor               | Open   |
| OQ-03 | **Dispatch consistency (`DEC-12`):** Confirm whether `F06-02` uses a transactional outbox (`SPK-03` SC-E) or another mechanism to keep PostgreSQL writes and RabbitMQ publishes consistent.                | Engineering / Mentor          | Open   |
| OQ-04 | **Max `source_code` byte size (`FR-04`):** Confirm the exact maximum byte limit (aligned with `F04` OQ-09) and HTTP status (`413` vs `400`).                                                               | Pair A / Pair D               | Open   |

**Blockers:**

- **For Sprint 2 draft completion (`#29`):** None (existing open questions are explicitly recorded above).
- **For formal `Approved` status (Sprint 3):** Resolution of OQ-01 (`CON-01`), dependent-pair agreement on `CON-02` (`#23`), and completion of the `§12` approval checklist.
- **For `F06-01` / `F06-02` implementation (Sprint 3):** Formal approval of `SPEC-F06`, `CON-05` (`#34`), `CON-02` (`#23`), and `DEC-12`.

---

## 12. Review and Approval *(Required)*

| Role     | Name or pair   | Date         | Outcome   |
|----------|----------------|--------------|-----------|
| Author   | Pair A         | 2026-11-19   | Submitted |
| Reviewer | Pair B, Pair D | *YYYY-MM-DD* | *Pending* |

**Approval checklist:**

- [ ] Acceptance criteria are unambiguous and numbered.
- [ ] Every required contract is defined and agreed.
- [ ] Every acceptance criterion has a verification entry.
- [ ] Security considerations are documented.
- [ ] No critical open questions remain.
- [ ] Reviewer is from a different pair.
