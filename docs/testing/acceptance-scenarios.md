# TST-05 / ST-02: Cross-Feature Consistency and Acceptance Scenarios

| Field | Value |
|-------|-------|
| Stretch ID / Backlog ID | ST-02 / TST-05 (proposed) |
| GitHub Issue | [#49](https://github.com/ITA-Internship/GoJudge/issues/49) |
| Sprint | Sprint 2 (Stretch) |
| Status | Draft |
| Owner | Pair A |
| Reviewer pair | Pair B, Pair C, or Pair D *(to be confirmed; must differ from Pair A)* |
| Last updated | 2026-11-26 |

> **Scope notice.** This document prepares reviewable, traceable cross-feature acceptance scenarios and records inconsistencies and unresolved decisions across `F03`, `F06`, and `F08`. It does **not** approve any specification or contract, does **not** settle any open product or architectural decision, and does **not** replace or implement the `TST-02` browser-driven end-to-end harness.

---

## 1. Scope and Conventions

### 1.1 Analysed Sources

| Artefact | Location / Status (as of 2026-11-26) | Role in this analysis |
|----------|--------------------------------------|-----------------------|
| `SPEC-F03` *(missing)* | **No specification file exists yet** (outside Sprint 2 committed scope; mentor confirmed on 25 Nov to work from existing descriptions) | Source requirements derived from `docs/product/mvp-proposal.md` (§3, §11), `docs/product/product-backlog.md` (`F03-01`, `F03-02`, `SPEC-F03`), `specs/features/F04-browser-code-editor.md` (`FR-10`–`FR-14`, `OQ-13`), `specs/features/F06-submit-solution.md` (`FR-03`), and `CON-04` (`§3.2`, `P-05`, `V-02`, `V-04`). Gaps recorded explicitly as missing requirements / blockers. |
| `SPEC-F06` | `specs/features/F06-submit-solution.md` (merged on `main`, Status: `In Review`, not formally Approved) | Full inventory (`FR-01`–`FR-11`, `AC-01`–`AC-08`, `OQ-01`–`OQ-04`). |
| `SPEC-F08` | `specs/features/F08-submission-history.md` (merged on `main` via PR `#41` on 2026-11-26, Status: `In Review`, not formally Approved) | Full inventory (`FR-01`–`FR-21` including `FR-04a`, `AC-01`–`AC-19`, `OQ-01`–`OQ-10`). |
| `CON-01` | `docs/contracts/CON-01.md` (merged on `main`) | Canonical submission lifecycle states (`QUEUED`, `RUNNING`, `COMPLETED`, `FAILED`, `CANCELLED`), verdicts, invariants, and open `Duplicate Submissions` question. |
| `CON-02` | `docs/contracts/CON-02.md` (merged on `main`; `DEC-04` and `DEC-12` not yet approved) | Publisher contract (`F06-02`: IDs-only payload, `schema_version: 1`, `submissions.dispatch`, `idempotency_key = submission_id`) and crash-recovery / cancellation open questions. |
| `CON-04` | `docs/contracts/CON-04.md` (branch `docs/con-04-draft`, issue `#31`, PR `#47`, Status: `Draft`; `DEC-09` Open) | Problem logical content (`§3.2`), package rules (`P-01`–`P-10`), starter template (`§4`), and hidden-test visibility rules (`V-01`–`V-10`). |
| `CON-05` | `docs/contracts/CON-05.md` (branch `docs/con-05-preliminary`, PR `#46`, which includes the earlier PRs `#40` and `#44`; not yet merged to `main`; Status: `Preliminary`; not formally Approved) | Provisional `users`, `roles`, `submissions`, and `results` schema, `submissions_state_verdict_chk`, proposed `submissions_one_active_per_user` index (`C05-OQ-20`), and open questions (`C05-OQ-20`–`C05-OQ-27`). |
| `SPEC-F01`, `SPEC-F04` | Merged on `main` (Status: `In Review`, not formally Approved) | Supporting identity/auth (`users.id`, `401`, error envelope) and editor (`is_blank`, `U+0000`, `max_source_size`, starter template, fetch timeout) boundaries. |

### 1.2 Scenario ID and Status Conventions

- **Scenario ID prefixes:**
  - `SCN-E2E-XX` — End-to-end cross-feature learner journeys (`F03 → F06 → F08`).
  - `SCN-SUB-XX` — Submission intake, validation, lifecycle, and polling (`F06`, `CON-01`, `CON-02`).
  - `SCN-SEC-XX` — Authentication, authorisation, cross-user isolation, and hidden-test confidentiality (`F01`, `F03`, `F06`, `F08`, `CON-04`, `INF-07`).
  - `SCN-HIS-XX` — Submission history listing, display, and edge cases (`F08`).
  - `SCN-PAG-XX` — Submission history pagination consistency and boundary cases (`F08`).
- **Traceability status values:**
  - `Covered` — Mapped to one or more concrete scenarios with observable expected results.
  - `Partially Blocked` — Core flow is covered, but an unresolved decision (`OQ-*` / `DEC-*` / missing `SPEC-F03`) leaves part of the expected result open.
  - `Blocked` — Cannot produce a final testable expected result until a named missing spec or mentor decision is resolved.
  - `Not applicable` — Outside cross-feature acceptance scope (must state rationale).

---

## 2. Acceptance Scenario Catalogue

### 2.1 End-to-End Learner Journeys

#### `SCN-E2E-01` – Happy-path learner journey: published problem access → submit → queued → accepted → history

- **Source references:** `MVP Proposal §3 (F03, F06, F08)`, `Product Backlog F03-01`, `SPEC-F06 (FR-01, FR-02, FR-03, FR-05, FR-06, FR-07, FR-08, FR-11, AC-01, AC-02, AC-08)`, `SPEC-F08 (FR-01, FR-04, FR-04a, FR-05, FR-06, FR-12, FR-16, AC-01, AC-04, AC-05, AC-06)`, `CON-01 §1–§2`, `CON-02 §1–§3`, `CON-04 (P-05, V-02, V-04)`.
- **Preconditions:**
  - Learner A (`user_id = U-A`) is authenticated with a valid RS256 Bearer token (`F01`).
  - Problem `P-PUB` (`status = published`, title `"Two Sum"`) exists in storage with a statement, public examples, a valid UTF-8 starter template, at least one sample test, and at least one hidden test (publishing precondition is `Proposed` in `CON-04 P-08`, Draft).
  - Learner A has no currently active (`QUEUED` or `RUNNING`) submission for `P-PUB` (or any problem, pending `MQ-02`).
- **Test data:**
  - `problem_id = P-PUB` (UUID).
  - Valid non-blank Go `source_code` implementing the reference solution, free of `U+0000`, within `max_source_size` (`MQ-05`: `F06 OQ-04`).
- **Actions:**
  1. Learner A opens problem `P-PUB` (`F03-01` problem details fetch — *exact endpoint path blocked pending `SPEC-F03`, see `INC-01` / `MQ-01`*).
  2. Learner A submits `source_code` via `POST /api/v1/problems/{P-PUB}/submissions`.
  3. Immediately after creation, Learner A queries `GET /api/v1/submissions/{id}` while the submission is `QUEUED`, and continues regular short polling until a terminal state is reached.
  4. The Judge Orchestrator (`F07`) transitions the submission `QUEUED → RUNNING → COMPLETED` with verdict `ACCEPTED`.
  5. Learner A calls `GET /api/v1/submissions` (`F08-01`) and views the "My submissions" page (`F08-02`).
- **Observable expected results:**
  1. Problem details display title, statement, input/output contract, constraints, public examples, and starter template, and contain **zero** hidden test content (`CON-04 V-02`). *Unresolved parts:* exact `F03-01` JSON schema, statement format (`C04-OQ-03`), public-example vs sample-test equivalence (`C04-OQ-04`) await `MQ-01`, and starter-template shape (`DEC-09 §6.3`) awaits `MQ-03`.
  2. `POST` returns `202 Accepted` immediately without waiting for judging, returning `id` (UUID), `problem_id = P-PUB`, `state = "QUEUED"`, `verdict = null`, `created_at`, and `updated_at`; a `QUEUED` row is persisted in PostgreSQL and one ID-only message (`schema_version: 1`, `idempotency_key = id`) is dispatched per `CON-02`. *Unresolved part:* formal `CON-02` / `DEC-04` approval and `DEC-12` dispatch consistency mechanism await `MQ-04`.
  3. While `QUEUED` or `RUNNING`, `GET /api/v1/submissions/{id}` returns `200 OK` with `verdict: null`.
  4. Once `COMPLETED`, `GET /api/v1/submissions/{id}` returns `200 OK` with `state: "COMPLETED"`, `verdict: "ACCEPTED"`, and no hidden test details; UI polling stops (`F06 FR-11`, `AC-08`). *Unresolved parts:* polling interval/backoff bounds await `MQ-06`, and any post-submission detail beyond the `CON-01` verdict awaits `MQ-03` (`DEC-09 §6.4`).
  5. `GET /api/v1/submissions` includes the submission with the exact six `F06` fields (`id`, `problem_id`, `state`, `verdict`, `created_at`, `updated_at`) plus `problem_title: "Two Sum"`, contains no `source_code`, `user_id`, `email`, `auth_subject`, or token (`F08 FR-12`, `AC-04`), and appears newest first.

#### `SCN-E2E-02` – Queued, running, and terminal states (`COMPLETED`, `FAILED`, `CANCELLED`) consistency across `F06` and `F08`

- **Source references:** `SPEC-F06 (FR-08, FR-11, AC-02, AC-08, §7.2, §7.4)`, `SPEC-F08 (FR-04, FR-05, FR-06, FR-07, FR-17, AC-05, AC-06, AC-07, AC-14)`, `CON-01 §1–§3`, `CON-05 §3.3`.
- **Preconditions:**
  - Learner A owns five submissions for `P-PUB`:
    - `S-QUEUED`: `state = QUEUED`, `verdict = null`.
    - `S-RUN`: `state = RUNNING`, `verdict = null`.
    - `S-WA`: `state = COMPLETED`, `verdict = "WRONG_ANSWER"` (representative of all seven `CON-01` test verdicts: `ACCEPTED`, `WRONG_ANSWER`, `COMPILATION_ERROR`, `RUNTIME_ERROR`, `TIME_LIMIT_EXCEEDED`, `MEMORY_LIMIT_EXCEEDED`, `OUTPUT_LIMIT_EXCEEDED`).
    - `S-FAIL`: reached `FAILED` after infrastructure retries exhausted (`verdict = "INTERNAL_ERROR"`).
    - `S-CANCEL`: reached `CANCELLED` (`verdict = null`).
  - Separate pre-`F07-05` fixture: Learner A owns only non-verdict submissions (`QUEUED`, `RUNNING`, `CANCELLED`, all with `verdict = null`, consistent with `SPEC-F08 §10 AC-07` and `CON-05 submissions_state_verdict_chk`; see `INC-09`).
- **Test data:** Submissions `S-QUEUED`, `S-RUN`, `S-WA`, `S-FAIL`, `S-CANCEL` owned by `U-A` (plus the pre-`F07-05` fixture).
- **Actions:**
  1. Learner A queries `GET /api/v1/submissions/{id}` for each of `S-QUEUED`, `S-RUN`, `S-WA`, `S-FAIL`, `S-CANCEL`.
  2. Learner A calls `GET /api/v1/submissions` and inspects the history UI (`F08-02`).
  3. Repeat `GET /api/v1/submissions` against the pre-`F07-05` fixture where no verdicts are written yet.
- **Observable expected results:**
  1. `F06` state endpoint returns `200 OK` with exact `CON-01` names: `S-QUEUED` → `QUEUED` / `null`; `S-RUN` → `RUNNING` / `null`; `S-WA` → `COMPLETED` / `"WRONG_ANSWER"`; `S-FAIL` → `FAILED` / `"INTERNAL_ERROR"`; `S-CANCEL` → `CANCELLED` / `null`. UI polling stops for `COMPLETED`, `FAILED`, and `CANCELLED` (`F06 FR-11`, `AC-08`).
  2. `F08` history endpoint returns identical `state` and `verdict` values for all five items (`F08 FR-05`, `FR-06`, `AC-05`, `AC-06`). The UI displays the problem title (`"Two Sum"`), the send time (`created_at`) in the user's local time zone, the state, and `"Pending"` for `QUEUED` and `RUNNING`, `"WRONG_ANSWER"` for `S-WA`, `"INTERNAL_ERROR"` for `S-FAIL`, and `"—"` for `CANCELLED` in the verdict column (`F08 FR-16`, `FR-17`, `AC-14`).
  3. On the pre-`F07-05` fixture, `GET /api/v1/submissions` returns `200 OK` with `verdict: null` for every item without inventing verdicts (`F08 FR-07`, `AC-07`).
  4. *Open dependency notes:* (a) How `CANCELLED` propagates to an in-flight execution is open in `CON-02` (Owner: `F07-01`); this scenario verifies only the read contract once `CANCELLED` is persisted. (b) Interaction between pre-`F07-05` `F07-01` mock completion and `CON-05 submissions_state_verdict_chk` is recorded in `INC-09`.

---

### 2.2 Submission Intake, Validation, Lifecycle, and Polling Scenarios

#### `SCN-SUB-01` – Asynchronous intake returns immediately without blocking on judge execution

- **Source references:** `SPEC-F06 (FR-01, FR-02, FR-05, FR-06, FR-07, AC-01, §7.1, §7.3, §8)`, `CON-01 §1`, `CON-02 §1–§3`.
- **Preconditions:**
  - Learner A is authenticated; `P-PUB` is published.
  - Judge Orchestrator (`F07`) is paused or slow (not immediately consuming `submissions.dispatch`).
- **Test data:** Valid Go `source_code`, plus a variant JSON payload containing extraneous client-supplied `"user_id": "<U-B>"`, `"state": "COMPLETED"`, and `"verdict": "ACCEPTED"`.
- **Actions:**
  1. Send `POST /api/v1/problems/{P-PUB}/submissions` with the valid payload.
  2. Repeat with the variant payload containing extra `"user_id"` / `"state"` / `"verdict"` fields (after first submission reaches terminal state or on a clean fixture to avoid `SCN-SUB-04` overlap).
- **Observable expected results:**
  1. Response returns immediately with `202 Accepted` (`state: "QUEUED"`, `verdict: null`), regardless of Judge Orchestrator consumer lag.
  2. Stored `user_id` is strictly `U-A` from the verified `F01` Bearer token; initial `state` is always `QUEUED` and `verdict` is `null` (client-supplied `user_id`/`state`/`verdict` are ignored or rejected by strict schema parsing per `F06 §8`).
  3. Published RabbitMQ message contains only `{ "schema_version": 1, "submission_id", "problem_id", "user_id", "submitted_at" }` and metadata `idempotency_key = submission_id` (`source_code` is absent from the broker payload). *Unresolved part:* `CON-02` remains provisional pending `DEC-04` (`MA-26`) and `DEC-12` (`MA-27`) approval (`MQ-04`).

#### `SCN-SUB-02` – Invalid submission input rejected with zero persistence or queue side effects

- **Source references:** `SPEC-F06 (FR-04, AC-04, §7.1, §8, §9)`, `SPEC-F04 (FR-04, FR-07, FR-08, AC-13, AC-14, AC-17, OQ-09, OQ-11)`, `CON-04 T-06`, `CON-05 C05-OQ-23`.
- **Preconditions:** Learner A is authenticated; `P-PUB` is published.
- **Test data:** Seven invalid request cases:
  - (a) Empty string `""`
  - (b) Whitespace-only `"   \n\t  "` (`F04 is_blank`)
  - (c) Go code containing NUL byte `"package main\x00"` (`U+0000`)
  - (d) Payload exceeding `max_source_size` bytes (*exact byte limit unresolved: `MQ-05`*)
  - (e) Malformed JSON body
  - (f) Non-UUID `{problemId}` path parameter on `POST /api/v1/problems/{problemId}/submissions`
  - (g) Non-UUID `{submissionId}` path parameter on `GET /api/v1/submissions/{submissionId}` (`F06 §9`)
- **Actions:** Send cases (a)–(f) to `POST /api/v1/problems/{problemId}/submissions`, and case (g) to `GET /api/v1/submissions/{submissionId}`.
- **Observable expected results:**
  - Cases (a), (b), (c), (e), (f), (g) return `400 Bad Request` with error envelope `{"error": {"code": "validation_error", "message": "..."}}` (`F01 §7.1`).
  - Case (d) is rejected without silent truncation; **unresolved expected status & threshold:** `413 Content Too Large` (`payload_too_large`) vs `400 Bad Request` (`validation_error`) and the exact byte limit require mentor/pair confirmation (`MQ-05`: `F06 OQ-04`, `F04 OQ-09`, `CON-04 T-06`, `CON-05 C05-OQ-23`).
  - Across all seven cases: zero rows are inserted into PostgreSQL and zero messages are published to RabbitMQ.

#### `SCN-SUB-03` – Non-existent or unpublished (`draft`) problem rejected on submit

- **Source references:** `SPEC-F06 (FR-03, AC-05, §7.1, §8)`, `CON-04 (P-05)`.
- **Preconditions:**
  - Learner A is authenticated.
  - `P-DRAFT` exists with `status = draft`; `P-MISSING` is a valid UUID not present in the database.
- **Test data:** Valid Go `source_code`.
- **Actions:**
  1. Send `POST /api/v1/problems/{P-DRAFT}/submissions`.
  2. Send `POST /api/v1/problems/{P-MISSING}/submissions`.
- **Observable expected results:**
  - Both requests return `404 Not Found` (`{"error": {"code": "not_found", "message": "..."}}`), disclosing no existence of unpublished problems.
  - Zero rows are persisted and zero messages are published.
  - *Cross-feature note:* Access to `P-DRAFT` via `F03-01` (404 vs 403) is separately unresolved pending `SPEC-F03` (`MQ-01`), whereas previously submitted records for a later-unpublished problem remain visible in `F08` history (`F08 FR-04a`, `AC-18`).

#### `SCN-SUB-04` – Repeated / duplicate active submission handling *(Blocked — Unresolved Decision)*

- **Source references:** `SPEC-F06 (FR-10, AC-07, §7.1, §8, OQ-01)`, `CON-01 §1 (Duplicate Submissions Open Question)`, `CON-05 §3.3 & C05-OQ-20`.
- **Preconditions:**
  - Learner A is authenticated and already has submission `S-ACT` for `P-PUB` in state `QUEUED` (sub-case A) or `RUNNING` (sub-case B).
  - Separately (sub-case C): Learner A has `S-ACT` in `RUNNING` (or `QUEUED`) for `P-PUB` and submits to a *different* published problem `P-PUB-2`.
- **Test data:** Valid Go `source_code`.
- **Actions:** Learner A sends a second `POST /api/v1/problems/{problemId}/submissions` while `S-ACT` is still non-terminal.
- **Observable expected results *(Blocked on `MQ-02`: `CON-01` Duplicate Submissions, `F06 OQ-01`, `CON-05 C05-OQ-20`)*:**
  - **Competing interpretation 1 (`CON-01` default, `CON-05 C05-OQ-20` proposed default, mentor preference on PR `#39`):** Reject with `409 Conflict` (`submission_in_progress`), creating no new submission and dispatching no message.
  - **Competing interpretation 2:** Accept and enqueue a second `QUEUED` submission (requires explicit per-user rate-limiting rules to prevent queue starvation; currently not specified in `F06`).
  - **Competing scope interpretations (sub-cases A & C):** `CON-01` mentions `"while another is still RUNNING"` (omitting `QUEUED` and scope), `SPEC-F06 FR-10` scopes the rule to `QUEUED` or `RUNNING` **per `(user_id, problem_id)`** (sub-case C allowed), while `CON-05 C05-OQ-20` proposes `QUEUED` or `RUNNING` **per `user_id` across all problems** via `submissions_one_active_per_user` (sub-case C returns `409`). See `INC-02` and `MQ-02`.

#### `SCN-SUB-05` – Editor submit action, double-click prevention, retry visibility, and polling termination (`F06-04`)

- **Source references:** `SPEC-F06 (FR-11, AC-08, §6, §7.4, §8)`, `SPEC-F04 (FR-04, AC-17)`, `CON-01 §1`.
- **Preconditions:** Learner A is signed in on the problem workspace for `P-PUB` with valid non-blank Go code in the editor.
- **Test data:** Valid Go code; simulated delayed `202 Accepted` response followed by polling responses `QUEUED → RUNNING → QUEUED` (infrastructure retry per `CON-01 §1`) `→ RUNNING → COMPLETED (WRONG_ANSWER)`, plus a variant with a transient network timeout during polling.
- **Actions:**
  1. Learner A double-clicks Submit rapidly.
  2. After `202 Accepted`, the UI performs regular short polling against `GET /api/v1/submissions/{id}` through the retry transition until terminal state.
  3. Variant: a polling request fails or times out transiently.
- **Observable expected results:**
  1. Submit button disables while `POST` is in flight; only one `POST` request is sent.
  2. UI displays `QUEUED`, updates to `RUNNING`, remains polling without treating `RUNNING → QUEUED` (infrastructure retry) as an error or terminal state, and finally displays `COMPLETED` / `WRONG_ANSWER` and stops polling.
  3. On transient polling error/timeout, UI backs off/stops within bounded retry rules and surfaces a recoverable error **without clearing editor content** (`F06 §8`). *Unresolved part:* exact polling interval, backoff schedule, and max timeout/retry count are not specified yet (`INC-06` / `MQ-06`).

#### `SCN-SUB-06` – Persistence and broker unavailability during submission intake

- **Source references:** `SPEC-F06 (§7.1, §8, OQ-03)`, `CON-02 §2, §8, §9`, `SPK-03 SC-E`.
- **Preconditions:** Learner A is authenticated; `P-PUB` is published.
- **Test data:** Valid Go `source_code`.
- **Actions:**
  1. Attempt `POST /api/v1/problems/{P-PUB}/submissions` while PostgreSQL is unreachable.
  2. Attempt `POST` while PostgreSQL is reachable and RabbitMQ is temporarily unreachable.
- **Observable expected results:**
  1. When PostgreSQL is unreachable: API returns `503 Service Unavailable` (`service_unavailable`), no partial record exists, and zero messages are published.
  2. When RabbitMQ is unreachable after DB write: **unresolved expected result — depends on `DEC-12` (`MA-27`, `MQ-04`, `F06 OQ-03`)**:
     - Under the `SPK-03` SC-E recommended transactional outbox, the submission and outbox event commit atomically in PostgreSQL (`202 Accepted`) and the relay publishes once RabbitMQ recovers (zero orphaned submissions).
     - Under direct publish without outbox, the API would either fail with `503` (requiring rollback/compensation) or risk DB/broker divergence.

---

### 2.3 Security, Access Control, and Hidden-Test Confidentiality Scenarios

#### `SCN-SEC-01` – Unauthenticated access denied across submission and history endpoints (`401`)

- **Source references:** `SPEC-F01 (FR-09, FR-10, AC-08, AC-09, AC-14)`, `SPEC-F06 (FR-01, AC-03, §9)`, `SPEC-F08 (FR-13, AC-12, §9)`.
- **Preconditions:** `P-PUB` exists; submission `S-1` owned by Learner A exists.
- **Test data:** (a) No `Authorization` header; (b) non-Bearer scheme (e.g., `Basic ...`); (c) expired, malformed, wrong-issuer, or wrong-audience JWT.
- **Actions:** Send unauthenticated requests to:
  1. `POST /api/v1/problems/{P-PUB}/submissions`
  2. `GET /api/v1/submissions/{S-1}`
  3. `GET /api/v1/submissions`
  4. Problem details endpoint `F03-01` *(auth requirement unspecified — see `INC-01` / `MQ-01`)*.
- **Observable expected results:**
  - Requests 1–3 return `401 Unauthorized` (`{"error": {"code": "unauthorized", "message": "..."}}`), create no DB rows, dispatch no queue messages, disclose no submission data, and log no raw token or `Authorization` header (`F01 AC-14`, `INF-07`).
  - Request 4 (`F03-01`): **Blocked pending `SPEC-F03` (`MQ-01`)** — whether published problem details require authentication (`401`) or are publicly readable (`200`) is not yet defined.

#### `SCN-SEC-02` – Cross-user access denied for both `learner` and `administrator` roles (IDOR protection)

- **Source references:** `SPEC-F06 (FR-02, FR-09, AC-06, §7.1, §9)`, `SPEC-F08 (FR-01, FR-02, FR-03, AC-01, AC-02, AC-03, §9)`.
- **Preconditions:**
  - Learner A (`U-A`, role `learner`) owns submissions `S-A1`, `S-A2`, `S-A3`.
  - Learner B (`U-B`, role `learner`) owns submission `S-B1`.
  - Admin C (`U-C`, role `administrator`) owns submission `S-C1`.
- **Test data:** `S-A1` UUID, plus tampered `?user_id=U-A` query parameter, custom `X-User-Id: U-A` header, and request body attempts.
- **Actions:**
  1. Learner B calls `GET /api/v1/submissions/{S-A1}`.
  2. Admin C calls `GET /api/v1/submissions/{S-A1}`.
  3. Learner A, Learner B, and Admin C each call `GET /api/v1/submissions` (including Learner B and Admin C sending `?user_id=U-A` / header / body overrides).
- **Observable expected results:**
  1. Learner B receives `404 Not Found` (`not_found`) with zero submission data disclosed (`F06 FR-09`, `AC-06`).
  2. Admin C receives `404 Not Found` (`F06 FR-09` restricts retrieval to `submission.user_id == caller.user_id`; see `INC-04`).
  3. On `GET /api/v1/submissions`, Learner A receives only `{S-A1, S-A2, S-A3}`, Learner B receives only `{S-B1}`, and Admin C receives only `{S-C1}` (`F08 FR-02`, `FR-03`, `AC-01`, `AC-02`, `AC-03`), ignoring all client-supplied `user_id` or `role` values.

#### `SCN-SEC-03` – Hidden-test confidentiality and sensitive-field redaction across all learner-facing endpoints and logs

- **Source references:** `CON-04 (V-01, V-02, V-03, V-04, V-05, §5.2)`, `SPEC-F06 (AC-02, §9)`, `SPEC-F08 (FR-12, AC-04, §9)`, `Product Backlog RSK-07`, `INF-07`.
- **Preconditions:**
  - Problem `P-PUB` has public sample tests containing sentinel `SAMPLE-SENTINEL` and hidden tests containing unique sentinels `HIDDEN-INPUT-1`, `HIDDEN-OUTPUT-1`, `HIDDEN-TEST-NAME-1`, and `HIDDEN-CODE-1`.
  - Learner A (`email = learnerA@example.com`, `auth_subject = auth0|sentinel-sub`) submits a solution that prints its received input to `stdout` and panics on hidden test input.
- **Test data:** Hidden test sentinels (`HIDDEN-INPUT-1`, `HIDDEN-OUTPUT-1`, `HIDDEN-TEST-NAME-1`, `HIDDEN-CODE-1`); `Authorization: Bearer TOKEN-SENTINEL`; `source_code` containing `CODE-SENTINEL`.
- **Actions:**
  1. Learner A fetches problem details (`F03-01` — *exact endpoint pending `SPEC-F03`, `MQ-01`*).
  2. Learner A submits the solution (`POST /api/v1/problems/{P-PUB}/submissions`), polls `GET /api/v1/submissions/{id}` through terminal verdict, and calls `GET /api/v1/submissions`.
  3. Inspect all API responses (success and error), RabbitMQ message payloads (`CON-02`), and structured application logs (`INF-07`).
- **Observable expected results:**
  1. None of `HIDDEN-INPUT-1`, `HIDDEN-OUTPUT-1`, `HIDDEN-TEST-NAME-1`, `HIDDEN-CODE-1`, or the learner's hidden-run `stdout`/`stderr`/panic output appear in any learner-facing HTTP response (`F03`, `F06`, `F08`) or error payload (`CON-04 V-02`, `V-05`). Additionally, `GET /api/v1/submissions` never returns `source_code`, `user_id`, `email`, `auth_subject`, or tokens (`F08 FR-12`, `AC-04`).
  2. None of the hidden sentinels, `TOKEN-SENTINEL`, or `CODE-SENTINEL` appear in structured application logs (`CON-04 V-03`, `INF-07`, `F06 §9`).
  3. `CODE-SENTINEL` and hidden sentinels never appear in RabbitMQ messages (`CON-02 §3`).
  4. *Unresolved boundary:* Whether filtered `COMPILATION_ERROR` diagnostics or hidden pass/total counts may be returned depends on `DEC-09 §6.4` and `C04-OQ-10` (`MQ-03`); exact `F03-01` endpoint path awaits `MQ-01`.

#### `SCN-SEC-04` – Untrusted learner code is never compiled or executed by the API intake layer (`F06`)

- **Source references:** `SPEC-F06 (§9 Execution isolation, Input validation)`, `SPEC-F04 (FR-03, AC-16)`, `MVP Proposal §7`.
- **Preconditions:** Learner A is authenticated; `P-PUB` is published.
- **Test data:** Go source code containing an `init()` function or test attempting file writes / network dial / OS command execution (e.g., `os/exec`), plus non-Go text with HTML/script tags `<script>alert(1)</script>`.
- **Actions:**
  1. Submit both payloads via `POST /api/v1/problems/{P-PUB}/submissions` with the Judge Orchestrator disabled (separately, to avoid `SCN-SUB-04` overlap).
  2. Query `GET /api/v1/submissions/{id}` and `GET /api/v1/submissions`.
- **Observable expected results:**
  1. The Go API accepts both payloads as plain UTF-8 text (`202 Accepted`, `QUEUED`) without compiling, evaluating, or executing any code inside the API process or container (`F06 §9`).
  2. No HTML/script payload is rendered as executable markup in any UI view (`F04 FR-03`, `AC-16`). Downstream isolated execution occurs solely in `F07`, gated by `DEC-08` / `MA-03` (outside `F06` scope).

---

### 2.4 Submission History and Pagination Scenarios

> **Status of the source.** `SPEC-F08` is merged on `main` (PR `#41`) with status `In Review` and is not approved. The scenarios below test what `SPEC-F08` currently states. Where the expected result depends on an open decision (`F08 OQ-02`, `OQ-07`, `OQ-09`, `OQ-08`, `OQ-10`, `C04-OQ-08`), the scenario is marked *Partially Blocked* and the open part is written as a competing interpretation, not as a rule (`MQ-07`, `MQ-08`).

#### `SCN-HIS-01` – Empty history for a signed-in learner with no submissions

- **Source references:** `SPEC-F08 (FR-11, FR-18, AC-11)`.
- **Preconditions:**
  - Learner N (`U-N`, role `learner`) is signed in and has never submitted.
  - Learner A (`U-A`) has at least one submission (to show that the empty result is not caused by a missing user).
- **Test data:** Valid access token for `U-N`; valid access token for `U-A`.
- **Actions:**
  1. Learner N calls `GET /api/v1/submissions` (no query parameters).
  2. Learner N opens the "My submissions" page.
  3. Learner N follows the link to the problem catalogue.
- **Observable expected results:**
  1. The API returns `200 OK` with `submissions: []` and `next_page_token: ""`. It does not return `404` and does not return Learner A's items (`F08 FR-11`, `AC-11`).
  2. The page shows the empty-history message and a link to the problem catalogue, with no list rows, no "Load more" control and no error message (`F08 FR-18`, `AC-11`).
  3. The link opens the problem catalogue. *Unresolved part:* the catalogue itself belongs to `F02` and has no specification in this analysis (see `INC-01` for the missing `F03` part); only the presence and target of the link are asserted.

#### `SCN-HIS-02` – Problem title for published, unpublished and unavailable problems *(Partially Blocked — Unresolved Decision)*

- **Source references:** `SPEC-F08 (FR-04, FR-04a, FR-16, AC-04, AC-18, §6, §8)`, `SPEC-F06 (FR-03)`, `CON-04 (§3.2, P-10)`, `CON-05 (§1, C05-OQ-21)`.
- **Preconditions:**
  - Learner A owns three submissions: `S-T1` for problem `P-PUB` (published, title `"Two Sum"`), `S-T2` for problem `P-UNPUB` (title `"Valid Parentheses"`, created as published and **unpublished after** the submission was made), and `S-T3` for problem `P-GONE`, whose problem record cannot be read.
  - `P-GONE` is created with a test double or a fixture without the problem row. Whether a real problem can be deleted is open (`C04-OQ-08`, proposed default: only unpublished), so this scenario does not assume a delete operation exists.
- **Test data:** The three submissions above, with different `created_at` values so that the order is fixed; the title of `P-PUB` is later changed to `"Two Sum II"`.
- **Actions:**
  1. Learner A calls `GET /api/v1/submissions` and opens the "My submissions" page.
  2. Learner A tries `POST /api/v1/problems/{P-UNPUB}/submissions` (a new submission to the unpublished problem).
  3. An administrator renames `P-PUB` to `"Two Sum II"`; Learner A calls `GET /api/v1/submissions` again.
- **Observable expected results:**
  1. All three submissions are listed. `S-T1` has `problem_title: "Two Sum"`, `S-T2` has `problem_title: "Valid Parentheses"` (the title is returned although the problem is not published), and `S-T3` has `problem_title: null` with its `problem_id` still present (`F08 FR-04a`, `AC-18`).
  2. The page shows `"Two Sum"` and `"Valid Parentheses"` for the first two rows, and the fixed text `"Problem unavailable"` with no link for `S-T3` (`F08 FR-16`, `AC-18`).
  3. The new submission to `P-UNPUB` is rejected with `404 Not Found`, while the old `S-T2` stays in the history. This shows that `F06 FR-03` applies only when a submission is created (`F08 §8`, `SCN-SUB-03`).
  4. After the rename, `S-T1` shows `problem_title: "Two Sum II"` (the current title, not a copy taken at submission time).
  5. *Unresolved parts:* (a) whether `S-T3` can occur in production depends on whether problems can be deleted (`C04-OQ-08`) and on the `submissions.problem_id` foreign key in `CON-05` (`C05-OQ-21`); `P-10` in `CON-04` is `Proposed`. (b) The exact `problems` table and title column are not defined (`F08 §7.2`, `CON-05 §1` explicitly excludes `problem` entities, `C05-OQ-21`). Both are recorded in `INC-10` and `MQ-08`.

#### `SCN-HIS-03` – History is read-only; unexpected stored data returns `500 internal_error` *(Partially Blocked — Unresolved Decision)*

- **Source references:** `SPEC-F08 (FR-14, FR-21, AC-13, AC-17, AC-19, §8, §10, OQ-09)`, `SPEC-F06 (§7.1)`, `CON-05 (submissions_state_verdict_chk)`, `CON-01 §1`.
- **Preconditions:**
  - Part A: Learner A owns a `RUNNING` submission `S-RUN` and a `QUEUED` submission `S-Q`.
  - Part B: the data access layer is replaced by a stub that returns one submission with an unknown `state` (for example `"PAUSED"`) and, in a second run, an unknown `verdict` (for example `"MAYBE"`). With the real `CON-05` check constraint such rows cannot be stored (`F08 §8`), so a stub is the only way to test this.
- **Test data:** `S-RUN`, `S-Q`; the two stubbed rows; request identifier header, if present.
- **Actions:**
  1. Part A: Learner A calls `GET /api/v1/submissions` ten times in a row, with and without `?user_id=...` and an unknown query parameter.
  2. Part A: after each call, read `state`, `verdict` and `updated_at` of `S-RUN` and `S-Q` from the database, and count the rows in `submissions` and the messages published to `submissions.dispatch`.
  3. Part B: call `GET /api/v1/submissions` against each stub.
  4. Part B: open the history page while the API returns `500`.
- **Observable expected results:**
  1. Part A: `state`, `verdict` and `updated_at` of both submissions are unchanged, no row is added or removed, and no message is published (`F08 FR-14`, `AC-13`).
  2. Part B: each stubbed call returns `500` with code `internal_error`, a generic message, no `submissions` array and no internal detail (stack trace, table or column names); the server log contains one error entry with the submission ID and the request ID (`F08 AC-19`, `§8`).
  3. Part B: the page shows an error message with "Try again" and no list rows (`F08 FR-21`, `AC-17`).
  4. *Unresolved part:* `500 internal_error` exists only in `SPEC-F08`; `SPEC-F01 §7.1` and `SPEC-F06 §7.1` do not list it (`INC-03`, `MQ-07`, `F08 OQ-09`). The expected code name is therefore the `SPEC-F08` value, not an agreed cross-feature value.

#### `SCN-HIS-04` – History page: authentication gate, Refresh, Load more and Try again *(Partially Blocked — Unresolved Decision)*

- **Source references:** `SPEC-F08 (FR-15, FR-19, FR-20, FR-21, AC-15, AC-16, AC-17, OQ-02, OQ-07, OQ-09)`, `SPEC-F01 (FR-10)`, `SPEC-F06 (FR-11)`.
- **Preconditions:**
  - Learner A is signed in and owns 25 submissions (the dataset of `SCN-PAG-01`), the newest one `S-25` in state `RUNNING`.
  - A mock of the API is available that can delay, fail (`503`, `500`, `400`) or drop a request.
- **Test data:** Dataset of `SCN-PAG-01`; for the Refresh part, `S-25` later changes to `COMPLETED` / `ACCEPTED`.
- **Actions:**
  1. Open the "My submissions" page without a valid session, then with an expired token.
  2. Sign in and open the page. Note the rows shown.
  3. While the page is open, let `S-25` reach `COMPLETED` / `ACCEPTED`, wait, then press "Refresh".
  4. Press "Load more" until the control disappears.
  5. Make the next API call fail (`503`, then `500`, then `400`, then a dropped connection), press "Load more" or "Refresh", then restore the API and press "Try again".
- **Observable expected results:**
  1. A visitor without a valid session cannot see any submission rows. The API call returns `401` with no data (`F08 FR-13`, `AC-12`) and the sign-in flow of `F01` is used. *Unresolved part:* the exact redirect or message is defined by `F01`, not by `F08` (`F08 FR-15` states only that the page is reachable by signed-in users).
  2. After sign-in the page shows 20 rows (default page size), the "Load more" control and the newest row `S-25` with `"Pending"` in the verdict column (`F08 FR-17`, `AC-14`).
  3. After "Refresh", the list is reloaded from the first page and the row shows `COMPLETED` and `ACCEPTED` (`F08 FR-20`, `AC-15`). *Unresolved part:* the scenario asserts nothing about automatic updating of the page; whether the page should refresh itself while some submission is `QUEUED` or `RUNNING` is open (`F08 OQ-07`, `MQ-06`).
  4. "Load more" first adds 5 rows (25 in total) with no duplicates, then the control is hidden because `next_page_token` is empty (`F08 FR-19`, `AC-16`). *Unresolved part:* the token and page-size rules are open (`F08 OQ-02`, `MQ-08`).
  5. For each failure the page shows an error message with "Try again" and does not present partial or stale rows as current data (`F08 FR-21`, `AC-17`). After the API is restored, "Try again" loads the list. *Unresolved part:* `500 internal_error` is not agreed across features (`INC-03`, `MQ-07`).

#### `SCN-PAG-01` – Stable dataset: traversal without duplicates or omissions and tie-break by `id` *(Partially Blocked — Unresolved Decision)*

- **Source references:** `SPEC-F08 (FR-08, FR-09, AC-08, AC-09, AC-16, §6, OQ-02, OQ-06)`.
- **Preconditions:**
  - Learner A (`U-A`) owns exactly the 25 submissions of the dataset below. Learner B (`U-B`) owns 5 other submissions with `created_at` values between those of Learner A (to prove that filtering by owner does not disturb the pages).
  - All submissions are for `P-PUB`. States and verdicts are not important for this scenario and are all `COMPLETED` / `ACCEPTED`.
  - The dataset is loaded by a fixture (not created by calls to `POST`), so that `created_at` and `id` are exactly as shown.
- **Test data (Learner A, all on 2026-11-20, UTC):**
  - The `id` is a UUID of the form `00000000-0000-4000-8000-0000000000NN`, written below as `S-NN`. Ordering by `id` descending is therefore the same as ordering by `NN` descending.
  - Four groups share the same `created_at` and are ordered by `id` descending: `S-23`/`S-24` (10:19:00), `S-14`/`S-15`/`S-16` (10:12:00), `S-10`/`S-11` (10:09:00) and `S-05`/`S-06` (10:04:00). Two of them are split by a page border at `page_size=10`: the group `S-14`–`S-16` between page 1 and page 2, and the group `S-05`/`S-06` between page 2 and page 3.

  | Position (newest first) | `id` | `created_at` | Page at `page_size=10` | Page at default size (20) |
  |---|---|---|---|---|
  | 1 | `S-25` | 10:20:00 | 1 | 1 |
  | 2 | `S-24` | 10:19:00 | 1 | 1 |
  | 3 | `S-23` | 10:19:00 | 1 | 1 |
  | 4 | `S-22` | 10:18:00 | 1 | 1 |
  | 5 | `S-21` | 10:17:00 | 1 | 1 |
  | 6 | `S-20` | 10:16:00 | 1 | 1 |
  | 7 | `S-19` | 10:15:00 | 1 | 1 |
  | 8 | `S-18` | 10:14:00 | 1 | 1 |
  | 9 | `S-17` | 10:13:00 | 1 | 1 |
  | 10 | `S-16` | 10:12:00 | 1 | 1 |
  | 11 | `S-15` | 10:12:00 | 2 | 1 |
  | 12 | `S-14` | 10:12:00 | 2 | 1 |
  | 13 | `S-13` | 10:11:00 | 2 | 1 |
  | 14 | `S-12` | 10:10:00 | 2 | 1 |
  | 15 | `S-11` | 10:09:00 | 2 | 1 |
  | 16 | `S-10` | 10:09:00 | 2 | 1 |
  | 17 | `S-09` | 10:08:00 | 2 | 1 |
  | 18 | `S-08` | 10:07:00 | 2 | 1 |
  | 19 | `S-07` | 10:06:00 | 2 | 1 |
  | 20 | `S-06` | 10:04:00 | 2 | 1 |
  | 21 | `S-05` | 10:04:00 | 3 | 2 |
  | 22 | `S-04` | 10:03:00 | 3 | 2 |
  | 23 | `S-03` | 10:02:00 | 3 | 2 |
  | 24 | `S-02` | 10:01:00 | 3 | 2 |
  | 25 | `S-01` | 10:00:00 | 3 | 2 |

  - Learner B's 5 submissions have `created_at` of 10:03:30, 10:08:30, 10:12:00 (same as the first split tie group `S-14`–`S-16`), 10:16:30 and 10:19:30.
- **Actions:**
  1. Learner A calls `GET /api/v1/submissions?page_size=10`, then calls it again with each returned `next_page_token` as `page_token`, until `next_page_token` is empty. Record every page.
  2. Learner A calls `GET /api/v1/submissions?page_size=25` once, and calls it without `page_size` and then with the returned token.
  3. Repeat step 1 for Learner B with `page_size=2`.
  4. Open the history page as Learner A, press "Load more" once.
- **Observable expected results:**
  1. Step 1 gives three pages with 10, 10 and 5 items. `next_page_token` is non-empty on pages 1 and 2 and `""` on page 3 (`F08 FR-09`, `AC-09`). Page 1 holds `S-25` … `S-16`, page 2 holds `S-15` … `S-06`, page 3 holds `S-05` … `S-01`.
  2. The 25 `id` values of the three pages are all different and are exactly the 25 ids of the dataset (no duplicate, no omission), including the three members of the first split tie group (`S-14`–`S-16`) and both members of the second split tie group (`S-05`/`S-06`) (`F08 AC-09`).
  3. The concatenation of the three pages is identical, item by item, to the single response with `page_size=25`, which has `next_page_token: ""` (stable order, `F08 FR-08`, `AC-08`).
  4. Without `page_size` the first page has 20 items (`S-25` … `S-06`) and the second page has 5 items (`S-05` … `S-01`), with no duplicates (`F08 AC-16`).
  5. No item of Learner B appears in any page of Learner A, and no item of Learner A appears for Learner B; each traversal of Learner B's 5 submissions with `page_size=2` gives 2, 2 and 1 items without duplicates (`F08 FR-03`).
  6. On the page, the first view shows 20 rows and, after "Load more", 25 rows with the control hidden (`F08 AC-16`).
  7. *Unresolved part:* the pagination model is open (`F08 OQ-02`, `MQ-08`). The dataset and the expected pages are valid for any model that orders by `created_at` descending and `id` descending (`F08 FR-08`), but the parameter names, the token and the response shape are the `SPEC-F08` proposal, not an agreed contract.

#### `SCN-PAG-02` – `page_size` and `page_token` rules, and a token that belongs to another user *(Partially Blocked — Unresolved Decision)*

- **Source references:** `SPEC-F08 (FR-02, FR-09, FR-10, AC-10, §8, OQ-02)`.
- **Preconditions:**
  - Learner A owns 130 submissions (fixture; any valid states). Learner B owns 3 submissions.
  - Learner A has obtained a valid `next_page_token` `T-A` from `GET /api/v1/submissions?page_size=10`.
- **Test data:**
  - `page_size` values: *(unset)*, `0`, `1`, `20`, `100`, `101`, `500`, `-1`, `abc`, `1.5`, and an empty value.
  - `page_token` values: a valid token (`T-A`), `T-A` with one character changed, `T-A` cut in half, an arbitrary string (`"abc"`), a base64 string built by the client with a plausible `created_at` and `id`, and an empty value.
- **Actions:**
  1. Learner A calls the endpoint with each `page_size` value.
  2. Learner A calls the endpoint with each invalid or modified `page_token` value.
  3. Learner B calls `GET /api/v1/submissions` with Learner A's valid token `T-A`.
- **Observable expected results:**
  1. *(unset)* and `0` return 20 items; `1` returns 1 item with a non-empty `next_page_token`; `20` returns 20 items; `100` returns 100 items; `101` and `500` return 100 items (lowered, no error) (`F08 FR-09`, `AC-10`, `§8`).
  2. `-1`, `abc` and `1.5` return `400` with code `validation_error` and no items (`F08 FR-10`, `AC-10`).
  3. The changed token, the cut token, `"abc"` and the client-built token return `400` with code `validation_error` and no items (`F08 FR-10`, `AC-10`, `§6`). The valid token returns the next page with no overlap with the previous page.
  4. In step 3, the response never contains an item of Learner A. Learner B receives at most their own 3 submissions, because the query always uses the caller's id (`F08 §8`, `FR-02`, `FR-03`).
  5. *Unresolved parts:* (a) whether the response of step 3 is `200` with Learner B's own list or `400` is not stated: `F08 FR-10` rejects tokens "not issued by the server", while `T-A` was issued by the server to another caller; the scenario asserts only that no data of Learner A is returned (`MQ-08`). (b) An empty `page_size=` or an empty `page_token=` value is not described in `F08`; no expected result is asserted (`MQ-08`). (c) All parameter names and limits are the proposal of `F08 OQ-02`.

#### `SCN-PAG-03` – Concurrent inserts and token lifetime *(Partially Blocked — Unresolved Decision, no final rule asserted)*

- **Source references:** `SPEC-F08 (FR-08, FR-09, §6, §8, OQ-02)`, `INC-05`.
- **Preconditions:** Learner A owns the 25 submissions of `SCN-PAG-01`. Learner A has page 1 with `page_size=10` and its token `T-1`.
- **Test data:** A new submission `S-26` of Learner A with `created_at` 10:21:00 (newer than all others), created by `POST /api/v1/problems/{P-PUB}/submissions` or by the fixture after page 1 was read. Token `T-1` is replayed after (a) a waiting time and (b) a restart of the API and a rotation of the server key. The waiting time is a parameter: **no value is set here, because `F08` defines no token lifetime**.
- **Actions:**
  1. Part A (concurrent insert): read page 1, create `S-26`, then read page 2 and page 3 with the returned tokens. Then call the first page again.
  2. Part B (token lifetime): replay `T-1` after waiting time `T`, then after an API restart, then after a key rotation.
- **Observable expected results:**
  1. Part A, **documented expectation** of `F08 §8` (not an acceptance criterion): `S-26` does not appear in page 2 or page 3, pages 2 and 3 keep the content of `SCN-PAG-01`, and no row is repeated or skipped. The first page called again, or "Refresh", starts with `S-26`. **This is recorded as the current text of `SPEC-F08`, not as a settled requirement:** `SPEC-F08` has no `FR` or `AC` for concurrent inserts and keyset pagination is open (`F08 OQ-02`, `INC-05`).
  2. Part B has **no asserted expected result.** The outcomes under discussion are: the token stays valid (same page returned), or it is rejected with `400 validation_error` as a token "not issued by the server", or a separate expired-token code is introduced. The choice, the lifetime and the behaviour after key rotation are for the mentor (`MQ-08`).
  3. Not asserted either: a row that arrives with a `created_at` older than the position of the cursor (for example a back-dated import or clock skew), a change of `page_size` between two pages, and a row removed between two pages (`INC-10`). They are listed so that the test plan does not silently assume an answer.

---

## 3. Requirements Traceability Matrix

### 3.1 `SPEC-F03` – Problem Details *(Derived from Existing Descriptions — Formal `SPEC-F03` Missing)*

> `SPEC-F03` has no file in `specs/features/` yet (outside Sprint 2 committed scope; gated by `CON-04` final). The rows below inventory all requirements currently stated across `MVP Proposal`, `Product Backlog`, `SPEC-F04`, `SPEC-F06`, and `CON-04 (Draft)`.

| Req ID | Source | Summary | Mapped Scenarios | Status | Rationale / Blocker |
|--------|--------|---------|------------------|--------|---------------------|
| F03-DESC-01 | `MVP Proposal §3 (F03)`, `Backlog F03-01`, `CON-04 §3.2` | Display problem title, statement, input/output contract, and constraints | `SCN-E2E-01` | Partially Blocked | No formal `SPEC-F03` endpoint/JSON schema; title max length (`C04-OQ-01`), statement format (`C04-OQ-03`), and separate vs embedded constraints are open (`MQ-01`). |
| F03-DESC-02 | `MVP Proposal §3 (F03)`, `Backlog F03-01`, `CON-04 §3.2` | Display public examples for a published problem | `SCN-E2E-01` | Partially Blocked | Whether public examples are identical to `F05` sample tests is open (`CON-04 C04-OQ-04`, `MQ-01`). |
| F03-DESC-03 | `CON-04 P-05`, `Backlog F02-01`, `SPEC-F06 FR-03` | Learners access published problems only (draft/unpublished hidden from learners) | `SCN-E2E-01`, `SCN-SUB-03` | Partially Blocked | Draft-problem HTTP status on `F03-01` (`404` vs `403`) and authentication policy on `F03-01` require `SPEC-F03` (`MQ-01`). |
| F03-DESC-04 | `SPEC-F04 FR-10–FR-14, OQ-13`, `Backlog F04-01`, `CON-04 §4` | Problem details deliver the optional starter template to the editor (and define fetch error/timeout boundary) | `SCN-E2E-01` | Partially Blocked | Exact starter template format/shape awaits `DEC-09 §6.3` (`CON-04 T-09`, `SPEC-F04 OQ-06`, `MQ-03`); problem-data retrieval timeout is open (`SPEC-F04 OQ-13`). |
| F03-DESC-05 | `CON-04 V-02, V-04`, `Backlog RSK-07` | Problem details never return hidden test content | `SCN-E2E-01`, `SCN-SEC-03` | Partially Blocked | Hard security invariant is stated in `CON-04 V-02`/`V-04`, but the `F03-01` HTTP endpoint and response schema await `SPEC-F03` (`MQ-01`). |
| F03-DESC-06 | `Backlog F03-02` | React problem details UI page | `SCN-E2E-01` | Blocked | No `SPEC-F03` UI acceptance criteria exist yet (`INC-01`, `MQ-01`). |

### 3.2 `SPEC-F06` – Submit Solution (`specs/features/F06-submit-solution.md`)

| Req ID | Type | Summary | Mapped Scenarios | Status | Rationale / Blocker |
|--------|------|---------|------------------|--------|---------------------|
| FR-01 | FR | Authenticated endpoint to submit a Go solution (`F06-01`) | `SCN-E2E-01`, `SCN-SUB-01`, `SCN-SEC-01` | Covered | — |
| FR-02 | FR | Derive `user_id` exclusively from `F01` auth context; ignore body `user_id` | `SCN-E2E-01`, `SCN-SUB-01`, `SCN-SEC-02` | Covered | — |
| FR-03 | FR | Verify target problem exists and is published; else `404 Not Found` | `SCN-E2E-01`, `SCN-SUB-03` | Covered | — |
| FR-04 | FR | Reject blank/whitespace (`is_blank`), `U+0000`, or oversized `source_code` with no side effects | `SCN-SUB-02` | Partially Blocked | Exact `max_source_size` byte limit and `413` vs `400` status open (`F06 OQ-04`, `F04 OQ-09`, `CON-05 C05-OQ-23`, `MQ-05`). |
| FR-05 | FR | Persist submission in `QUEUED` with `verdict: null` using exact `CON-01` names | `SCN-E2E-01`, `SCN-SUB-01` | Covered | Final PostgreSQL DDL/types owned by `CON-05` (PR `#46`, Preliminary). |
| FR-06 | FR | Dispatch evaluation message per `CON-02` and `DEC-12` | `SCN-E2E-01`, `SCN-SUB-01`, `SCN-SUB-06` | Partially Blocked | `CON-02` approval (`DEC-04` / `MA-26`) and `DEC-12` (`MA-27`) dispatch mechanism are open (`F06 OQ-02`, `OQ-03`, `MQ-04`). |
| FR-07 | FR | Return immediately with `202 Accepted` and submission fields (`QUEUED`, `verdict: null`) | `SCN-E2E-01`, `SCN-SUB-01` | Covered | — |
| FR-08 | FR | Authenticated endpoint to retrieve single submission state and verdict (`verdict: null` for `QUEUED`/`RUNNING`/`CANCELLED`) | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| FR-09 | FR | Restrict state lookup to owning user; return `404 Not Found` for other users | `SCN-SEC-02` | Covered | — |
| FR-10 | FR | Duplicate active submission handling follows `CON-01` policy (`OQ-01`) | `SCN-SUB-04` | Blocked | `CON-01` Duplicate Submissions question (`409` vs queue, plus `QUEUED` and per-problem vs per-user scope conflict with `CON-05 C05-OQ-20`) is open (`INC-02`, `MQ-02`). |
| FR-11 | FR | Editor Submit prevents double-click and uses regular short polling until terminal state | `SCN-E2E-01`, `SCN-E2E-02`, `SCN-SUB-05` | Partially Blocked | Exact polling interval and timeout/backoff bounds not yet specified (`INC-06`, `MQ-06`). |
| AC-01 | AC | Successful asynchronous submission (`202 Accepted`, `QUEUED`, dispatch) | `SCN-E2E-01`, `SCN-SUB-01` | Partially Blocked | API intake is covered; broker dispatch contract remains provisional pending `DEC-04` / `DEC-12` (`MQ-04`). |
| AC-02 | AC | Retrieve own submission state across all 5 `CON-01` states with correct verdict invariant and no hidden tests | `SCN-E2E-01`, `SCN-E2E-02`, `SCN-SEC-03` | Covered | Optional extra feedback beyond verdict remains open in `DEC-09 §6.4` (`MQ-03`). |
| AC-03 | AC | Unauthenticated request rejected (`401 Unauthorized`, no side effects) | `SCN-SEC-01` | Covered | — |
| AC-04 | AC | Invalid `source_code` rejected (`400`/`413`, no side effects) | `SCN-SUB-02` | Partially Blocked | Exact byte limit and `413` vs `400` open (`F06 OQ-04`, `MQ-05`). |
| AC-05 | AC | Non-existent or unpublished problem rejected (`404 Not Found`, no side effects) | `SCN-SUB-03` | Covered | — |
| AC-06 | AC | Cross-user access to submission state denied (`404 Not Found`) | `SCN-SEC-02` | Covered | — |
| AC-07 | AC | Duplicate active submission handling (open question `OQ-01`) | `SCN-SUB-04` | Blocked | Blocked on `CON-01` Duplicate Submissions, `F06 OQ-01`, and `CON-05 C05-OQ-20` (`MQ-02`). |
| AC-08 | AC | Editor submit action and regular short polling until terminal state | `SCN-E2E-01`, `SCN-E2E-02`, `SCN-SUB-05` | Partially Blocked | Exact polling interval/timeout open (`INC-06`, `MQ-06`). |

### 3.3 `SPEC-F08` – Submission History (`specs/features/F08-submission-history.md`)

| Req ID | Type | Summary | Mapped Scenarios | Status | Rationale / Blocker |
|--------|------|---------|------------------|--------|---------------------|
| FR-01 | FR | `GET /api/v1/submissions` returns caller's own submissions | `SCN-E2E-01`, `SCN-SEC-02` | Covered | — |
| FR-02 | FR | Derive owner strictly from validated token; ignore query/header/body user ID or role | `SCN-SEC-02` | Covered | — |
| FR-03 | FR | Never return another user's submissions (learners and administrators alike) | `SCN-SEC-02` | Covered | — |
| FR-04 | FR | Return the six `F06` fields plus `problem_title` | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| FR-04a | FR | `problem_title` shows current title (published or unpublished), or `null` if problem missing | `SCN-HIS-02` | Partially Blocked | The unpublished-problem case is covered. The missing-problem case depends on whether problems can be deleted (`C04-OQ-08`) and on the `submissions.problem_id` foreign key (`C05-OQ-21`); `CON-04 P-10` is `Proposed` (`INC-10`, `MQ-08`). |
| FR-05 | FR | Exact `CON-01` state and verdict names | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| FR-06 | FR | `verdict` is `null` for `QUEUED`, `RUNNING`, `CANCELLED`; test verdict for `COMPLETED`; `INTERNAL_ERROR` for `FAILED` | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| FR-07 | FR | Read stored verdict only; `null` until `F07-05` populates verdicts | `SCN-E2E-02` | Covered | See `INC-09` (`CON-05` check constraint vs `F07-01` mock). |
| FR-08 | FR | Order by `created_at DESC`, tie-break by `id DESC` | `SCN-PAG-01` | Covered | Ordering rule is stated in `F08 FR-08`; the index (`F08 OQ-06`) affects only performance. |
| FR-09 | FR | Keyset pagination (`page_size` default 20, max 100; `next_page_token`, empty on last page) | `SCN-PAG-01`, `SCN-PAG-02`, `SCN-PAG-03` | Partially Blocked | Keyset pagination, parameter names and response shape are open in `F08 OQ-02`; behaviour under concurrent inserts has no `AC` in `F08` (`INC-05`, `MQ-08`). |
| FR-10 | FR | Reject negative/non-integer `page_size` or invalid `page_token` with `400` | `SCN-PAG-02` | Partially Blocked | Rules are open in `F08 OQ-02`; the result of another user's valid token and of empty parameter values is not stated (`MQ-08`). |
| FR-11 | FR | `200 OK` with empty `submissions` list when caller has no submissions | `SCN-HIS-01` | Covered | — |
| FR-12 | FR | Never return `source_code`, hidden test data, other users' data, `auth_subject`, email, or tokens | `SCN-E2E-01`, `SCN-SEC-03` | Covered | — |
| FR-13 | FR | `401` without valid access token | `SCN-SEC-01` | Covered | — |
| FR-14 | FR | Endpoint is strictly read-only | `SCN-HIS-03` | Covered | — |
| FR-15 | FR | "My submissions" page reachable only by signed-in users | `SCN-HIS-04` | Covered | Exact sign-in redirect belongs to `F01`. |
| FR-16 | FR | Row displays `problem_title` (or `"Problem unavailable"` when `null`), local time, state, verdict | `SCN-E2E-01`, `SCN-E2E-02`, `SCN-HIS-02` | Partially Blocked | Display of a title is covered; the `"Problem unavailable"` branch depends on `C04-OQ-08` and `C05-OQ-21` (`INC-10`, `MQ-08`). |
| FR-17 | FR | UI verdict column shows `"Pending"` (`QUEUED`/`RUNNING`) and `"—"` (`CANCELLED` / other `null`) | `SCN-E2E-02` | Covered | — |
| FR-18 | FR | Empty history message + catalogue link | `SCN-HIS-01` | Covered | — |
| FR-19 | FR | "Load more" control using `next_page_token`, hidden when empty | `SCN-PAG-01`, `SCN-HIS-04` | Partially Blocked | Token and page-size rules are open in `F08 OQ-02` (`MQ-08`). |
| FR-20 | FR | Manual "Refresh" control reloading from page 1 | `SCN-HIS-04` | Partially Blocked | Manual refresh is covered; auto-refresh vs manual-only is open (`F08 OQ-07`, `MQ-06`). The scenario asserts nothing about automatic updates. |
| FR-21 | FR | Error message + "Try again" on `400`/`500`/`503`/network error; `401` handled by `F01` sign-in flow | `SCN-HIS-03`, `SCN-HIS-04` | Partially Blocked | `500 internal_error` is not agreed across features (`F08 OQ-09`, `INC-03`, `MQ-07`). |
| AC-01 | AC | Learner sees only own submissions | `SCN-E2E-01`, `SCN-SEC-02` | Covered | — |
| AC-02 | AC | Owner cannot be overridden by query/header/body | `SCN-SEC-02` | Covered | — |
| AC-03 | AC | Administrator sees only own submissions | `SCN-SEC-02` | Covered | — |
| AC-04 | AC | Submission fields match `F06` + `problem_title`, no forbidden fields | `SCN-E2E-01`, `SCN-SEC-03` | Covered | — |
| AC-05 | AC | `QUEUED` and `RUNNING` have `verdict: null` | `SCN-E2E-02` | Covered | — |
| AC-06 | AC | `COMPLETED`, `FAILED`, `CANCELLED` exact states/verdicts | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| AC-07 | AC | Pre-`F07-05` all-`null` verdicts still return `200` | `SCN-E2E-02` | Covered | See `INC-09` regarding `CON-05` check constraint. |
| AC-08 | AC | Ordering by `created_at DESC`, tie-breaker `id DESC` | `SCN-PAG-01` | Covered | Dataset includes four groups with equal `created_at`. |
| AC-09 | AC | Multi-page traversal (25 rows, `page_size=10` → 10/10/5, no duplicates/gaps) | `SCN-PAG-01` | Partially Blocked | Open in `F08 OQ-02` (`MQ-08`); dataset and expected pages are fixed in `SCN-PAG-01`. |
| AC-10 | AC | `page_size` bounds and invalid `page_size`/`page_token` (`400 validation_error`) | `SCN-PAG-02` | Partially Blocked | Limits and parameter names are the proposal of `F08 OQ-02` (`MQ-08`). |
| AC-11 | AC | Empty list API + UI empty state | `SCN-HIS-01` | Covered | — |
| AC-12 | AC | Missing/expired token returns `401 unauthorized` | `SCN-SEC-01` | Covered | — |
| AC-13 | AC | Repeated reads do not mutate DB row | `SCN-HIS-03` | Covered | — |
| AC-14 | AC | UI renders title, `"Pending"`, verdict, `"—"`, local time | `SCN-E2E-02` | Covered | — |
| AC-15 | AC | UI "Refresh" reloads updated state/verdict | `SCN-HIS-04` | Partially Blocked | Manual refresh only; auto-refresh is open (`F08 OQ-07`, `MQ-06`). |
| AC-16 | AC | UI "Load more" appends next page and hides when exhausted | `SCN-PAG-01`, `SCN-HIS-04` | Partially Blocked | Open in `F08 OQ-02` (`MQ-08`). |
| AC-17 | AC | UI error state + "Try again" on network/`503`/`500` | `SCN-HIS-03`, `SCN-HIS-04` | Partially Blocked | The `500` case depends on `F08 OQ-09` (`INC-03`, `MQ-07`). |
| AC-18 | AC | Unpublished problem keeps title; missing problem shows `null` / `"Problem unavailable"` | `SCN-HIS-02` | Partially Blocked | Missing-problem branch depends on `C04-OQ-08` and `C05-OQ-21` (`INC-10`, `MQ-08`). |
| AC-19 | AC | Unknown `state`/`verdict` in storage returns `500 internal_error` and logs submission ID | `SCN-HIS-03` | Partially Blocked | `500 internal_error` not yet shared with `F01`/`F06` (`F08 OQ-09`, `INC-03`, `MQ-07`). |

---

## 4. Identified Inconsistencies and Missing Dependencies

### 4.1 Cross-Feature Consistency Comparison (`F03`, `F06`, `F08`, `CON-01`, `CON-04`)

| Dimension | Alignment across `F03`, `F06`, `F08` and contracts | Gaps / Inconsistencies |
|-----------|----------------------------------------------------|------------------------|
| **Terminology** | `F06` and `F08` share exact JSON field names (`id`, `problem_id`, `state`, `verdict`, `created_at`, `updated_at`), `CON-01` state/verdict names, and `F01` role names (`learner`, `administrator`). `CON-02 §3` uses `submission_id`/`submitted_at` (mapped explicitly in `F06 §7.3`). | `MVP Proposal` (`F03`) says `"public examples"` vs `F05`/`CON-04` `"sample tests"` (`C04-OQ-04`); `F04 §7.2` uses `source_text` vs `F06` `source_code`; `F04 OQ-11` leaves `is_blank` open while `F06 FR-04` rejects whitespace-only; `MVP Proposal §3` omits `OUTPUT_LIMIT_EXCEEDED` and `CANCELLED` added by `CON-01`. |
| **Actors** | `F01`, `F06`, and `F08` recognise `learner`, `administrator`, and unauthenticated callers consistently. | `F03` actor/auth boundary is undefined (`INC-01` / `MQ-01`); `F06` lacks an explicit AC for `administrator` submit/read (`INC-04`). |
| **Inputs** | `F06` (`POST /api/v1/problems/{problemId}/submissions`, `GET /api/v1/submissions/{submissionId}`) and `F08` (`GET /api/v1/submissions`) share `/api/v1/submissions` collection alignment and UUID identifiers. | `F03-01` endpoint path/params missing (`INC-01`); `max_source_size` and `413` vs `400` open (`MQ-05`); `page_size`/`page_token` open for `CON-03` (`F08 OQ-02`). |
| **Outputs** | `F08` reuses the exact six `F06` submission fields and adds nullable `problem_title`. Error envelope `{"error": {"code", "message"}}` matches `F01 §7.1`. | `F08 AC-19` introduces `500 internal_error` missing from `F01 §7.1` & `F06 §7.1` (`INC-03`); `F03-01` response schema missing (`INC-01`); optional judge detail in responses open (`DEC-09 §6.4`). |
| **State transitions** | `F06` and `F08` both enforce `CON-01` states, verdicts, and `verdict: null` for `QUEUED`/`RUNNING`/`CANCELLED`. | Duplicate active submission mismatch across `CON-01`, `F06`, `CON-05` (`INC-02`); `CON-02 §4`/`CON-05 C05-OQ-27` crash recovery `RUNNING → RUNNING` missing from `CON-01` (`INC-08`); `CANCELLED` propagation open (`CON-02`); pre-`F07-05` `null` verdict vs `CON-05` CHECK constraint (`INC-09`). |
| **Ownership** | Both `F06 FR-02` and `F08 FR-02` derive `user_id` strictly from `F01` token context and ignore client-supplied `user_id`/`role`. | Consistent across `F06` and `F08`; `F01 OQ-1` (role source of truth) remains open upstream. |
| **Access rules** | `F06` and `F08` both enforce `401` when unauthenticated and owner-only visibility for both `learner` and `administrator` (`F06` returns `404` on cross-user lookup; `F08` filters by caller `user_id`). | `F03-01` auth policy (public vs signed-in) and draft-problem status (`404` vs `403`) undefined (`INC-01` / `MQ-01`). |
| **Hidden-test visibility** | `F06`, `F08`, and `CON-04 (V-02..V-05)` consistently forbid hidden test content in all learner-facing responses and logs. | Filtered `COMPILATION_ERROR` diagnostics (`V-08` / `C04-OQ-10`), post-submission counts/index (`DEC-09 §6.4`), and judge filesystem isolation (`C04-OQ-09` / `DEC-08`) remain open (`MQ-03`). |
| **Submission-history behaviour** | `F06 FR-03` blocks submitting to an unpublished problem (`404`), while `F08 FR-04a` intentionally keeps past submissions for unpublished problems with current `problem_title` (`CON-04 P-10`, Proposed). | Pagination model (`F08 OQ-02`), cursor TTL/key rotation & missing concurrent-insert AC (`INC-05`), polling/refresh model (`INC-06`), and problem/user deletion policies (`C04-OQ-08`, `F08 OQ-08` / `C05-OQ-22`) remain open. |

### 4.2 Detailed Inconsistencies and Missing Dependencies

| ID | Category | Sources | Description | Affected Scenarios |
|----|----------|---------|-------------|--------------------|
| INC-01 | Missing specification | `Product Backlog SPEC-F03`, `MVP Proposal §3`, `SPEC-F04 OQ-13`, `CON-04 §3.2` | No `specs/features/F03-problem-details.md` file exists yet. As a result, the `F03-01` HTTP path, response JSON schema, authentication requirement (public vs authenticated), draft-problem error code (`404` vs `403`), statement rendering (`C04-OQ-03`), public-example vs sample-test equivalence (`C04-OQ-04`), and retrieval timeout (`F04 OQ-13`) are undefined. | `SCN-E2E-01`, `SCN-SEC-01`, `SCN-SEC-03` |
| INC-02 | Cross-document inconsistency | `CON-01 §1 (Open Questions)` vs `SPEC-F06 (FR-10, AC-07, OQ-01)` vs `CON-05 (C05-OQ-20)` | Three-way mismatch on duplicate active submissions: (1) `CON-01` asks about submitting *"while another is still `RUNNING`"* (omitting `QUEUED` and leaving scope unspecified; Pair D's "Start of Sprint 2" deadline passed unresolved); (2) `SPEC-F06 FR-10` scopes it to `QUEUED` or `RUNNING` **for the same problem** (`(user_id, problem_id)`); (3) `CON-05 C05-OQ-20` proposes one active (`QUEUED`/`RUNNING`) submission **per user across all problems** via `submissions_one_active_per_user`. | `SCN-SUB-04` |
| INC-03 | Error contract inconsistency | `SPEC-F01 §7.1`, `SPEC-F06 §7.1` vs `SPEC-F08 (AC-19, OQ-09)` | `SPEC-F08` introduces `500 internal_error` for unexpected storage/domain states (`AC-19`), whereas `SPEC-F01 §7.1` and `SPEC-F06 §7.1` only list `400`/`401`/`403`/`404`/`409`/`413`/`503` and do not specify a `500` error code or envelope. | `SCN-E2E-02`, `SCN-HIS-03`, `SCN-HIS-04` |
| INC-04 | Role / access-policy phrasing | `SPEC-F06 (FR-09, §7.1)` vs `SPEC-F08 (FR-03, AC-03)` | Both specs enforce owner-only access for both `learner` and `administrator`, so they are aligned in intent. However, `SPEC-F06 FR-09` does not explicitly verify `administrator` in an AC (only in `§7.1`), whereas `SPEC-F08 AC-03` tests administrator own-only access explicitly; `SPEC-F06` also permits `administrator` on `POST` (`§7.1`) without an AC testing administrator submit. | `SCN-SEC-02` |
| INC-05 | Unconfirmed pagination & missing cursor-expiry / AC rules | `SPEC-F08 (FR-08, FR-09, §6, §8, OQ-02)` | `SPEC-F08 §6` & `§8` define encrypted keyset tokens `(created_at, id)` and describe concurrent-insert stability (new rows do not shift/duplicate subsequent pages and appear after Refresh). However: (1) keyset pagination and its parameters remain **Open** in `F08 OQ-02` (not yet confirmed for `CON-03`); (2) `SPEC-F08` has no `FR` or `AC` testing concurrent inserts; and (3) cursor token TTL / expiry and server encryption key rotation behaviour are unspecified. | `SCN-PAG-01`, `SCN-PAG-02`, `SCN-PAG-03` |
| INC-06 | Missing client polling bounds | `SPEC-F06 (FR-11, AC-08, §8)` vs `SPEC-F08 (FR-20, OQ-07)` | `SPEC-F06` specifies regular short polling with bounded interval and retry/timeout handling, but defines no concrete interval, maximum polling duration, or backoff threshold. Separately, `SPEC-F08` uses manual refresh (`OQ-07` open for auto-refresh). | `SCN-SUB-05`, `SCN-HIS-04` |
| INC-07 | Unmerged / provisional contract dependencies | `CON-04 (#31, PR #47, Draft)`, `CON-05 (#34, PR #46)`, `CON-02 (DEC-04/DEC-12 open)`, `CON-03 (not started)` | `SPEC-F06` and `SPEC-F08` depend on `CON-04` (`Draft`, `DEC-09` Open), `CON-05` (`Preliminary`, PR `#46`), `CON-02` (`DEC-04` and `DEC-12` unapproved), and `CON-03` (OpenAPI not started). Exact DB constraints, indexes (`F08 OQ-06`), `ON DELETE` rules (`F08 OQ-08` / `C05-OQ-22`, `C04-OQ-08`), and hidden-test feedback detail (`CON-04 §6.4`, `F08 OQ-04`) remain provisional. | `SCN-E2E-01`, `SCN-SUB-06`, `SCN-SEC-03` |
| INC-08 | `CON-01` vs `CON-02` / `CON-05` crash-recovery transition gap | `CON-01 §1` vs `CON-02 §4 (Open Questions)` & `CON-05 (C05-OQ-27)` | `CON-02 §4` and `CON-05 §3.3` re-take a `RUNNING` submission after an Orchestrator crash via atomic `UPDATE ... WHERE state IN ('QUEUED', 'RUNNING')`, whereas `CON-01 §1` only models `RUNNING → QUEUED` for infrastructure retry and does not describe `RUNNING → RUNNING` crash recovery (flagged in `CON-02` Open Questions and `CON-05 C05-OQ-27`). Also, `CON-02` leaves `CANCELLED` propagation to an active Job open (`F07-01`). | `SCN-E2E-02`, `SCN-SUB-06` |
| INC-09 | Pre-`F07-05` `null` verdicts vs `CON-05` DB check constraint | `SPEC-F08 (FR-07, AC-07, §10)` vs `CON-05 (§3.3 submissions_state_verdict_chk)` & `Product Backlog F07-01` | `CON-05 submissions_state_verdict_chk` forbids `verdict IS NULL` when `state IN ('COMPLETED', 'FAILED')`. `SPEC-F08 §10 AC-07` avoids violating this by seeding only `QUEUED`/`RUNNING`/`CANCELLED` rows before `F07-05`; however, if `F07-01`'s Sprint 3 mock executor transitions submissions to `COMPLETED`/`FAILED` before `F07-05` exists, it must write a valid `CON-01` placeholder verdict or violate `CON-05`. | `SCN-E2E-02` |
| INC-10 | Open deletion and lifecycle policies affecting history | `SPEC-F08 (FR-04a, AC-18, OQ-08, OQ-10)`, `CON-05 (C05-OQ-22)`, `CON-04 (C04-OQ-08, P-10)` | `SPEC-F08` assumes that past submissions stay visible when a problem is unpublished or its record cannot be read (`CON-04 P-10` is `Proposed`). Three related decisions are open and are not described in any specification: (1) whether a problem with submissions can be deleted or only unpublished (`C04-OQ-08`; the proposed default `RESTRICT` would make the `"Problem unavailable"` case a data fault, not a normal state); (2) the `ON DELETE` rule of `submissions.user_id` (`F08 OQ-08` / `C05-OQ-22`; `CASCADE` would remove the whole history of a deleted user); (3) whether a per-problem history (`GET /api/v1/problems/{problemId}/submissions`, the same path as the `F06` create endpoint) is needed (`F08 OQ-10`; default: not in the MVP, so no `F03`, `F06` or `F08` requirement covers it). A row removed between two pages of a traversal is also not described. | `SCN-HIS-02`, `SCN-PAG-03` |

---

## 5. Unresolved Mentor Questions

| ID | Question & Competing Interpretations | Source Decision / OQ | Affected Scenarios |
|----|--------------------------------------|----------------------|--------------------|
| MQ-01 | **Missing `SPEC-F03` contract:** What are the exact `F03-01` endpoint path, authentication rule (public vs signed-in learners only), draft-problem HTTP status (`404` vs `403`), statement format (`C04-OQ-03`: Markdown vs plain text, embedded vs separate constraints), and whether public examples equal `F05` sample tests (`C04-OQ-04`)? | `SPEC-F03` (missing), `CON-04 (C04-OQ-03, C04-OQ-04)`, `SPEC-F04 OQ-13` | `SCN-E2E-01`, `SCN-SEC-01`, `SCN-SEC-03` |
| MQ-02 | **Duplicate active submission policy & scope:** (1) Return `409 Conflict` without queuing (`CON-01` / `CON-05 C05-OQ-20` default, mentor comment on PR `#39`) vs queue additional submissions (which requires rate limiting to prevent queue starvation)? (2) Does the restriction apply to `QUEUED` + `RUNNING` (`F06` / `CON-05`) or `RUNNING` only (`CON-01`), and is it scoped per `(user_id, problem_id)` (`F06 FR-10`) or globally per `user_id` (`CON-05 C05-OQ-20`)? | `CON-01 §1`, `SPEC-F06 OQ-01`, `CON-05 C05-OQ-20` | `SCN-SUB-04` |
| MQ-03 | **Problem packaging, hidden-test feedback, and compiler diagnostics (`DEC-09`):** (1) Which storage/test/template options (`CON-04 §6.1–§6.3`) are adopted? (2) What post-submission feedback is allowed (`§6.4`: verdict only vs counts vs first failed index vs per-problem)? (3) May filtered `COMPILATION_ERROR` diagnostics be shown (`V-08` / `C04-OQ-10`)? (4) Are problem packages restricted to Go `stdlib` only (`C04-OQ-12`)? | `DEC-09 (MA-14)`, `CON-04 (C04-OQ-09, C04-OQ-10, C04-OQ-12)`, `SPEC-F08 OQ-04` | `SCN-E2E-01`, `SCN-SEC-03` |
| MQ-04 | **Messaging & dispatch consistency (`DEC-04`, `DEC-12`):** Formal approval of RabbitMQ + Watermill (`DEC-04` / `MA-26`) and confirmation of the dispatch consistency mechanism (`DEC-12` / `MA-27`, e.g., Watermill SQL transactional outbox from `SPK-03` SC-E vs another outbox/CDC approach). | `DEC-04 (MA-26)`, `DEC-12 (MA-27)`, `SPEC-F06 (OQ-02, OQ-03)` | `SCN-E2E-01`, `SCN-SUB-01`, `SCN-SUB-06` |
| MQ-05 | **Maximum `source_code` byte limit & oversized HTTP status:** What is the exact byte limit shared across `F04`, `F06`, `CON-04`, and `CON-05`, and should oversized submissions return `413 Content Too Large` (`payload_too_large`) or `400 Bad Request` (`validation_error`)? | `SPEC-F04 OQ-09`, `SPEC-F06 OQ-04`, `CON-04 T-06`, `CON-05 C05-OQ-23` | `SCN-SUB-02` |
| MQ-06 | **Client polling bounds and history refresh model:** What are the allowed polling interval, backoff, and max timeout for `F06-04`, and should `F08-02` remain manual refresh only or auto-refresh active submissions? | `SPEC-F06 FR-11`, `SPEC-F08 OQ-07` | `SCN-SUB-05`, `SCN-HIS-04` |
| MQ-07 | **Unified error envelope & `500 internal_error` code:** Confirm shared error envelope and codes across `F01`, `F06`, `F08`, and `CON-03`, including whether `500 internal_error` (`F08 AC-19`) is standardised across all API endpoints. | `SPEC-F01 OQ-7`, `SPEC-F08 OQ-09`, `CON-03` | `SCN-E2E-02`, `SCN-HIS-03`, `SCN-HIS-04` |
| MQ-08 | **Pagination, token lifetime and deletion policies for the history (`F08`):** (1) Is keyset pagination with an opaque token, default 20 and maximum 100, and the response shape `submissions` / `next_page_token` confirmed for `CON-03` (`F08 OQ-02`)? Should `F06` and other list endpoints use the same parameters? (2) Does a `page_token` have a lifetime or expire after a restart or a key rotation? Candidates: it stays valid; it is rejected with `400 validation_error`; a separate expired-token code is added. Is the token bound to the caller, and what is returned when another user sends a valid token (`200` with the caller's own list, or `400`)? Are empty `page_size` and `page_token` values accepted? (3) Is stable behaviour under concurrent inserts a requirement (an `AC` in `F08`) or only a description (`F08 §8`)? (4) Can a problem with submissions be deleted or only unpublished (`C04-OQ-08`), and what is the `ON DELETE` rule of `submissions.user_id` (`F08 OQ-08` / `C05-OQ-22`)? (5) Is a per-problem history needed in the MVP (`F08 OQ-10`)? | `SPEC-F08 (OQ-02, OQ-08, OQ-10)`, `CON-05 (C05-OQ-22)`, `CON-04 (C04-OQ-08)` | `SCN-PAG-01`, `SCN-PAG-02`, `SCN-PAG-03`, `SCN-HIS-02`, `SCN-HIS-04` |