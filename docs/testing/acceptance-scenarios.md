# TST-05 / ST-02: Cross-Feature Consistency and Acceptance Scenarios

| Field | Value |
|-------|-------|
| Stretch ID / Backlog ID | ST-02 / TST-05 (proposed) |
| GitHub Issue | [#49](https://github.com/ITA-Internship/GoJudge/issues/49) |
| Sprint | Sprint 2 (Stretch) |
| Status | Draft |
| Owner | Pair A |
| Reviewer pair | Pair B, Pair C, or Pair D *(to be confirmed; must differ from Pair A)* |
| Last updated | 2026-11-25 |

> **Scope notice.** This document prepares reviewable, traceable cross-feature acceptance scenarios and records inconsistencies and unresolved decisions across `F03`, `F06`, and `F08`. It does **not** approve any specification or contract, does **not** settle any open product or architectural decision, and does **not** replace or implement the `TST-02` browser-driven end-to-end harness.

---

## 1. Scope and Conventions

### 1.1 Analysed Sources

| Artefact | Location / Status (as of 2026-11-25) | Role in this analysis |
|----------|--------------------------------------|-----------------------|
| `SPEC-F03` *(missing)* | **No specification file exists yet** (outside Sprint 2 committed scope; mentor confirmed on 25 Nov to work from existing descriptions) | Source requirements derived from `docs/product/mvp-proposal.md` (§3, §11), `docs/product/product-backlog.md` (`F03-01`, `F03-02`, `SPEC-F03`), `specs/features/F04-browser-code-editor.md` (`FR-10`–`FR-14`, `OQ-13`), `specs/features/F06-submit-solution.md` (`FR-03`), and `CON-04` (`§3.2`, `P-05`, `V-02`, `V-04`). Gaps recorded explicitly as missing requirements / blockers. |
| `SPEC-F06` | `specs/features/F06-submit-solution.md` (merged on `main`, Status: `In Review`, not formally Approved) | Full inventory (`FR-01`–`FR-11`, `AC-01`–`AC-08`, `OQ-01`–`OQ-04`). |
| `SPEC-F08` | `specs/features/F08-submission-history.md` (branch `docs/spec-f08`, PR `#41`, Status: `In Review`; not yet merged to `main` and not formally Approved) | Full inventory (`FR-01`–`FR-21` including `FR-04a`, `AC-01`–`AC-19`, `OQ-01`–`OQ-10`). |
| `CON-01` | `docs/contracts/CON-01.md` (merged on `main`) | Canonical submission lifecycle states (`QUEUED`, `RUNNING`, `COMPLETED`, `FAILED`, `CANCELLED`), verdicts, invariants, and open `Duplicate Submissions` question. |
| `CON-02` | `docs/contracts/CON-02.md` (merged on `main`; `DEC-04` and `DEC-12` not yet approved) | Publisher contract (`F06-02`: IDs-only payload, `schema_version: 1`, `submissions.dispatch`, `idempotency_key = submission_id`) and crash-recovery / cancellation open questions. |
| `CON-04` | `docs/contracts/CON-04.md` (branch `docs/con-04-draft`, issue `#31`, Status: `Draft`; `DEC-09` Open) | Problem logical content (`§3.2`), package rules (`P-01`–`P-10`), starter template (`§4`), and hidden-test visibility rules (`V-01`–`V-10`). |
| `CON-05` | `docs/contracts/CON-05.md` (branch `docs/con-05-preliminary`, PR `#40`, Status: `Preliminary`; not formally Approved) | Provisional `users`, `roles`, `submissions`, and `results` schema, `submissions_state_verdict_chk`, proposed `submissions_one_active_per_user` index (`C05-OQ-20`), and open questions (`C05-OQ-20`–`C05-OQ-27`). |
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
  2. `F08` history endpoint returns identical `state` and `verdict` values for all five items (`F08 FR-05`, `FR-06`, `AC-05`, `AC-06`). The UI displays `"Pending"` for `QUEUED` and `RUNNING`, `"WRONG_ANSWER"` for `S-WA`, `"INTERNAL_ERROR"` for `S-FAIL`, and `"—"` for `CANCELLED` (`F08 FR-17`, `AC-14`).
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

*Pending `SPEC-F08` submission history (`SCN-HIS-*`) and pagination (`SCN-PAG-*`) scenarios.*

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
| FR-05 | FR | Persist submission in `QUEUED` with `verdict: null` using exact `CON-01` names | `SCN-E2E-01`, `SCN-SUB-01` | Covered | Final PostgreSQL DDL/types owned by `CON-05` (PR `#40`, Preliminary). |
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

### 3.3 `SPEC-F08` – Submission History (`specs/features/F08-submission-history.md`, PR `#41`)

| Req ID | Type | Summary | Mapped Scenarios | Status | Rationale / Blocker |
|--------|------|---------|------------------|--------|---------------------|
| FR-01 | FR | `GET /api/v1/submissions` returns caller's own submissions | `SCN-E2E-01`, `SCN-SEC-02` | Covered | — |
| FR-02 | FR | Derive owner strictly from validated token; ignore query/header/body user ID or role | `SCN-SEC-02` | Covered | — |
| FR-03 | FR | Never return another user's submissions (learners and administrators alike) | `SCN-SEC-02` | Covered | — |
| FR-04 | FR | Return the six `F06` fields plus `problem_title` | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| FR-04a | FR | `problem_title` shows current title (published or unpublished), or `null` if problem missing | Pending (`§2.4`) | Pending | `C04-OQ-08` (delete vs unpublish) open. |
| FR-05 | FR | Exact `CON-01` state and verdict names | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| FR-06 | FR | `verdict` is `null` for `QUEUED`, `RUNNING`, `CANCELLED`; test verdict for `COMPLETED`; `INTERNAL_ERROR` for `FAILED` | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| FR-07 | FR | Read stored verdict only; `null` until `F07-05` populates verdicts | `SCN-E2E-02` | Covered | See `INC-09` (`CON-05` check constraint vs `F07-01` mock). |
| FR-08 | FR | Order by `created_at DESC`, tie-break by `id DESC` | Pending (`§2.4`) | Pending | — |
| FR-09 | FR | Keyset pagination (`page_size` default 20, max 100; `next_page_token`, empty on last page) | Pending (`§2.4`) | Pending | Keyset pagination params/shape open in `F08 OQ-02` (`CON-03`). |
| FR-10 | FR | Reject negative/non-integer `page_size` or invalid `page_token` with `400` | Pending (`§2.4`) | Pending | Pagination rules open in `F08 OQ-02`. |
| FR-11 | FR | `200 OK` with empty `submissions` list when caller has no submissions | Pending (`§2.4`) | Pending | — |
| FR-12 | FR | Never return `source_code`, hidden test data, other users' data, `auth_subject`, email, or tokens | `SCN-E2E-01`, `SCN-SEC-03` | Covered | — |
| FR-13 | FR | `401` without valid access token | `SCN-SEC-01` | Covered | — |
| FR-14 | FR | Endpoint is strictly read-only | Pending (`§2.4`) | Pending | — |
| FR-15 | FR | "My submissions" page reachable only by signed-in users | Pending (`§2.4`) | Pending | — |
| FR-16 | FR | Row displays `problem_title` (or `"Problem unavailable"` when `null`), local time, state, verdict | `SCN-E2E-01`, `SCN-E2E-02` | Pending | Unpublished/missing title fallback pending `§2.4`. |
| FR-17 | FR | UI verdict column shows `"Pending"` (`QUEUED`/`RUNNING`) and `"—"` (`CANCELLED` / other `null`) | `SCN-E2E-02` | Covered | — |
| FR-18 | FR | Empty history message + catalogue link | Pending (`§2.4`) | Pending | — |
| FR-19 | FR | "Load more" control using `next_page_token`, hidden when empty | Pending (`§2.4`) | Pending | Pagination contract open in `F08 OQ-02`. |
| FR-20 | FR | Manual "Refresh" control reloading from page 1 | Pending (`§2.4`) | Pending | Auto-refresh vs manual-only open (`F08 OQ-07`, `MQ-06`). |
| FR-21 | FR | Error message + "Try again" on `400`/`500`/`503`/network error; `401` handled by `F01` sign-in flow | Pending (`§2.4`) | Pending | — |
| AC-01 | AC | Learner sees only own submissions | `SCN-E2E-01`, `SCN-SEC-02` | Covered | — |
| AC-02 | AC | Owner cannot be overridden by query/header/body | `SCN-SEC-02` | Covered | — |
| AC-03 | AC | Administrator sees only own submissions | `SCN-SEC-02` | Covered | — |
| AC-04 | AC | Submission fields match `F06` + `problem_title`, no forbidden fields | `SCN-E2E-01`, `SCN-SEC-03` | Covered | — |
| AC-05 | AC | `QUEUED` and `RUNNING` have `verdict: null` | `SCN-E2E-02` | Covered | — |
| AC-06 | AC | `COMPLETED`, `FAILED`, `CANCELLED` exact states/verdicts | `SCN-E2E-01`, `SCN-E2E-02` | Covered | — |
| AC-07 | AC | Pre-`F07-05` all-`null` verdicts still return `200` | `SCN-E2E-02` | Covered | See `INC-09` regarding `CON-05` check constraint. |
| AC-08 | AC | Ordering by `created_at DESC`, tie-breaker `id DESC` | Pending (`§2.4`) | Pending | — |
| AC-09 | AC | Multi-page traversal (25 rows, `page_size=10` → 10/10/5, no duplicates/gaps) | Pending (`§2.4`) | Pending | Pagination contract open in `F08 OQ-02`. |
| AC-10 | AC | `page_size` bounds and invalid `page_size`/`page_token` (`400 validation_error`) | Pending (`§2.4`) | Pending | Pagination contract open in `F08 OQ-02`. |
| AC-11 | AC | Empty list API + UI empty state | Pending (`§2.4`) | Pending | — |
| AC-12 | AC | Missing/expired token returns `401 unauthorized` | `SCN-SEC-01` | Covered | — |
| AC-13 | AC | Repeated reads do not mutate DB row | Pending (`§2.4`) | Pending | — |
| AC-14 | AC | UI renders title, `"Pending"`, verdict, `"—"`, local time | `SCN-E2E-02` | Covered | — |
| AC-15 | AC | UI "Refresh" reloads updated state/verdict | Pending (`§2.4`) | Pending | Auto-refresh vs manual-only open (`F08 OQ-07`, `MQ-06`). |
| AC-16 | AC | UI "Load more" appends next page and hides when exhausted | Pending (`§2.4`) | Pending | Pagination contract open in `F08 OQ-02`. |
| AC-17 | AC | UI error state + "Try again" on network/`503`/`500` | Pending (`§2.4`) | Pending | — |
| AC-18 | AC | Unpublished problem keeps title; missing problem shows `null` / `"Problem unavailable"` | Pending (`§2.4`) | Pending | `C04-OQ-08` (delete vs unpublish) open. |
| AC-19 | AC | Unknown `state`/`verdict` in storage returns `500 internal_error` and logs submission ID | Pending (`§2.4`) | Pending | `500 internal_error` not yet shared with `F01`/`F06` (`INC-03`, `MQ-07`). |

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
| INC-03 | Error contract inconsistency | `SPEC-F01 §7.1`, `SPEC-F06 §7.1` vs `SPEC-F08 (AC-19, OQ-09)` | `SPEC-F08` introduces `500 internal_error` for unexpected storage/domain states (`AC-19`), whereas `SPEC-F01 §7.1` and `SPEC-F06 §7.1` only list `400`/`401`/`403`/`404`/`409`/`413`/`503` and do not specify a `500` error code or envelope. | `SCN-E2E-02`, `§2.4` |
| INC-04 | Role / access-policy phrasing | `SPEC-F06 (FR-09, §7.1)` vs `SPEC-F08 (FR-03, AC-03)` | Both specs enforce owner-only access for both `learner` and `administrator`, so they are aligned in intent. However, `SPEC-F06 FR-09` does not explicitly verify `administrator` in an AC (only in `§7.1`), whereas `SPEC-F08 AC-03` tests administrator own-only access explicitly; `SPEC-F06` also permits `administrator` on `POST` (`§7.1`) without an AC testing administrator submit. | `SCN-SEC-02` |
| INC-05 | Unconfirmed pagination & missing cursor-expiry / AC rules | `SPEC-F08 (FR-08, FR-09, §6, §8, OQ-02)` | `SPEC-F08 §6` & `§8` define encrypted keyset tokens `(created_at, id)` and describe concurrent-insert stability (new rows do not shift/duplicate subsequent pages and appear after Refresh). However: (1) keyset pagination and its parameters remain **Open** in `F08 OQ-02` (not yet confirmed for `CON-03`); (2) `SPEC-F08` has no `FR` or `AC` testing concurrent inserts; and (3) cursor token TTL / expiry and server encryption key rotation behaviour are unspecified. | `§2.4` (`SCN-PAG-*`) |
| INC-06 | Missing client polling bounds | `SPEC-F06 (FR-11, AC-08, §8)` vs `SPEC-F08 (FR-20, OQ-07)` | `SPEC-F06` specifies regular short polling with bounded interval and retry/timeout handling, but defines no concrete interval, maximum polling duration, or backoff threshold. Separately, `SPEC-F08` uses manual refresh (`OQ-07` open for auto-refresh). | `SCN-SUB-05`, `§2.4` |
| INC-07 | Unmerged / provisional contract dependencies | `CON-04 (#31, Draft)`, `CON-05 (#34, PR #40)`, `CON-02 (DEC-04/DEC-12 open)`, `CON-03 (not started)` | `SPEC-F06` and `SPEC-F08` depend on `CON-04` (`Draft`, `DEC-09` Open), `CON-05` (`Preliminary`, PR `#40`), `CON-02` (`DEC-04` and `DEC-12` unapproved), and `CON-03` (OpenAPI not started). Exact DB constraints, indexes (`F08 OQ-06`), `ON DELETE` rules (`F08 OQ-08` / `C05-OQ-22`, `C04-OQ-08`), and hidden-test feedback detail (`CON-04 §6.4`, `F08 OQ-04`) remain provisional. | `SCN-E2E-01`, `SCN-SUB-06`, `SCN-SEC-03` |
| INC-08 | `CON-01` vs `CON-02` / `CON-05` crash-recovery transition gap | `CON-01 §1` vs `CON-02 §4 (Open Questions)` & `CON-05 (C05-OQ-27)` | `CON-02 §4` and `CON-05 §3.3` re-take a `RUNNING` submission after an Orchestrator crash via atomic `UPDATE ... WHERE state IN ('QUEUED', 'RUNNING')`, whereas `CON-01 §1` only models `RUNNING → QUEUED` for infrastructure retry and does not describe `RUNNING → RUNNING` crash recovery (flagged in `CON-02` Open Questions and `CON-05 C05-OQ-27`). Also, `CON-02` leaves `CANCELLED` propagation to an active Job open (`F07-01`). | `SCN-E2E-02`, `SCN-SUB-06` |
| INC-09 | Pre-`F07-05` `null` verdicts vs `CON-05` DB check constraint | `SPEC-F08 (FR-07, AC-07, §10)` vs `CON-05 (§3.3 submissions_state_verdict_chk)` & `Product Backlog F07-01` | `CON-05 submissions_state_verdict_chk` forbids `verdict IS NULL` when `state IN ('COMPLETED', 'FAILED')`. `SPEC-F08 §10 AC-07` avoids violating this by seeding only `QUEUED`/`RUNNING`/`CANCELLED` rows before `F07-05`; however, if `F07-01`'s Sprint 3 mock executor transitions submissions to `COMPLETED`/`FAILED` before `F07-05` exists, it must write a valid `CON-01` placeholder verdict or violate `CON-05`. | `SCN-E2E-02` |

---

## 5. Unresolved Mentor Questions

| ID | Question & Competing Interpretations | Source Decision / OQ | Affected Scenarios |
|----|--------------------------------------|----------------------|--------------------|
| MQ-01 | **Missing `SPEC-F03` contract:** What are the exact `F03-01` endpoint path, authentication rule (public vs signed-in learners only), draft-problem HTTP status (`404` vs `403`), statement format (`C04-OQ-03`: Markdown vs plain text, embedded vs separate constraints), and whether public examples equal `F05` sample tests (`C04-OQ-04`)? | `SPEC-F03` (missing), `CON-04 (C04-OQ-03, C04-OQ-04)`, `SPEC-F04 OQ-13` | `SCN-E2E-01`, `SCN-SEC-01`, `SCN-SEC-03` |
| MQ-02 | **Duplicate active submission policy & scope:** (1) Return `409 Conflict` without queuing (`CON-01` / `CON-05 C05-OQ-20` default, mentor comment on PR `#39`) vs queue additional submissions (which requires rate limiting to prevent queue starvation)? (2) Does the restriction apply to `QUEUED` + `RUNNING` (`F06` / `CON-05`) or `RUNNING` only (`CON-01`), and is it scoped per `(user_id, problem_id)` (`F06 FR-10`) or globally per `user_id` (`CON-05 C05-OQ-20`)? | `CON-01 §1`, `SPEC-F06 OQ-01`, `CON-05 C05-OQ-20` | `SCN-SUB-04` |
| MQ-03 | **Problem packaging, hidden-test feedback, and compiler diagnostics (`DEC-09`):** (1) Which storage/test/template options (`CON-04 §6.1–§6.3`) are adopted? (2) What post-submission feedback is allowed (`§6.4`: verdict only vs counts vs first failed index vs per-problem)? (3) May filtered `COMPILATION_ERROR` diagnostics be shown (`V-08` / `C04-OQ-10`)? (4) Are problem packages restricted to Go `stdlib` only (`C04-OQ-12`)? | `DEC-09 (MA-14)`, `CON-04 (C04-OQ-09, C04-OQ-10, C04-OQ-12)`, `SPEC-F08 OQ-04` | `SCN-E2E-01`, `SCN-SEC-03` |
| MQ-04 | **Messaging & dispatch consistency (`DEC-04`, `DEC-12`):** Formal approval of RabbitMQ + Watermill (`DEC-04` / `MA-26`) and confirmation of the dispatch consistency mechanism (`DEC-12` / `MA-27`, e.g., Watermill SQL transactional outbox from `SPK-03` SC-E vs another outbox/CDC approach). | `DEC-04 (MA-26)`, `DEC-12 (MA-27)`, `SPEC-F06 (OQ-02, OQ-03)` | `SCN-E2E-01`, `SCN-SUB-01`, `SCN-SUB-06` |
| MQ-05 | **Maximum `source_code` byte limit & oversized HTTP status:** What is the exact byte limit shared across `F04`, `F06`, `CON-04`, and `CON-05`, and should oversized submissions return `413 Content Too Large` (`payload_too_large`) or `400 Bad Request` (`validation_error`)? | `SPEC-F04 OQ-09`, `SPEC-F06 OQ-04`, `CON-04 T-06`, `CON-05 C05-OQ-23` | `SCN-SUB-02` |
| MQ-06 | **Client polling bounds and history refresh model:** What are the allowed polling interval, backoff, and max timeout for `F06-04`, and should `F08-02` remain manual refresh only or auto-refresh active submissions? | `SPEC-F06 FR-11`, `SPEC-F08 OQ-07` | `SCN-SUB-05`, `§2.4` |
| MQ-07 | **Unified error envelope & `500 internal_error` code:** Confirm shared error envelope and codes across `F01`, `F06`, `F08`, and `CON-03`, including whether `500 internal_error` (`F08 AC-19`) is standardised across all API endpoints. | `SPEC-F01 OQ-7`, `SPEC-F08 OQ-09`, `CON-03` | `SCN-E2E-02`, `§2.4` |
