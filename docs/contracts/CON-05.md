# CON-05: Data Model (Preliminary)

**Status:** Preliminary  
**Owner:** Pair D  
**Dependent Pairs:** Pair A, Pair B, Pair C

> **Disclaimer:** This document is a **Preliminary** draft and is **NOT** the final approved data contract. Formal approval is blocked until `CON-04` is finalized and dependent pairs agree on the schema.

---

## 1. Scope and Exclusions

**In scope (Preliminary Draft):**
- `users`
- `roles`
- `submissions`
- `results`

**Explicitly excluded:**
- `problem` and `test-asset` entities — **outside this draft, waiting for CON-04 final**.
- Database migrations and `INF-05` implementation.

**Blocking Final Approval:**
- Finalization of `CON-04` (required to define `problem` / `test-asset` entities and resolve foreign keys/types referencing `problems`).
- Resolution of open questions from `SPEC-F01` (`OQ-1`, `OQ-2`, `OQ-9`) affecting `users` and `roles`.
- Agreement from dependent pairs (`Pair A`, `Pair B`, `Pair C`).

---

## 2. Database Conventions (PostgreSQL)

- **Naming:** `snake_case`; plural table names (`users`, `roles`, `submissions`, `results`); foreign keys formatted as `<singular_table>_id`.
- **Primary Keys:** `UUID` (`NOT NULL`). Generation strategy (DB-side `gen_random_uuid()` v4 vs. Go-side `UUIDv7`) is deferred to `C05-OQ-04` (`INF-05`).
- **Timestamps:** `TIMESTAMPTZ` stored in UTC. Standard audit columns: `created_at TIMESTAMPTZ NOT NULL DEFAULT now()` and `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`.
  - **`updated_at` policy:** Updated explicitly by the application on `UPDATE` queries (no DB triggers).
  - **Exception (`results`):** `results` is immutable (*first terminal write wins*) and has `created_at` only (no `updated_at`).
- **Enums:** `TEXT NOT NULL` with explicit `CHECK (...)` constraints instead of native PostgreSQL `ENUM` types.
  - **Exception (`submissions.verdict`):** nullable `TEXT`, because `CON-01` requires the verdict to be `null` while the state is `QUEUED` or `RUNNING`. Allowed values are checked together with `state` by `submissions_state_verdict_chk` (section 3.3).
- **Unresolved Constraints:** Any field type, length, or `ON DELETE` policy dependent on unfinalized specs is explicitly recorded in Section 4 (`Open Questions`).
- **Open Question IDs:** two-digit format. `C05-OQ-01`..`C05-OQ-19` are reserved for `users`/`roles`/shared conventions; `C05-OQ-20`+ are reserved for `submissions`/`results`.

---

## 3. Entities

### 3.1 `users` *(Owner: Yevhen)*
*Draft in progress.*

| Field | Type | Nullable | Default | Constraints / Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | — | Primary Key |

### 3.2 `roles` *(Owner: Yevhen)*
*Draft in progress.*

| Field | Type | Nullable | Default | Constraints / Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | — | Primary Key |

### 3.3 `submissions` *(Owner: Yulia)*

One row per submission. Holds the lifecycle **state** (CON-01 §1) and the final **verdict** (CON-01 §2). The pair (`state`, `verdict`) is checked by one DB constraint, so the CON-01 state–verdict invariants cannot be broken by any writer.

| Field | Type | Nullable | Default | Constraints / Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | — (see `C05-OQ-04`) | Primary key. Also the idempotency key for messaging (`CON-02` §4: `idempotency_key = submission_id`), so no separate key column. |
| `user_id` | `UUID` | NO | — | FK → `users.id`. Owner of the submission (`SPEC-F01`: "Submission records reference `users.id` as owner"). `ON DELETE` policy: `C05-OQ-22`. |
| `problem_id` | `UUID` | NO | — | FK → `problems.id` (placeholder). Type and FK: `C05-OQ-21` (blocked by `CON-04`). |
| `source_code` | `TEXT` | NO | — | Submitted Go code. Loaded by the Orchestrator from PostgreSQL, never sent through the broker (`CON-02` §3). Maximum size: `C05-OQ-23`. |
| `state` | `TEXT` | NO | `'QUEUED'` | Names match `CON-01` §1. See `submissions_state_verdict_chk` below. |
| `verdict` | `TEXT` | **YES** | `NULL` | Names match `CON-01` §2. `NULL` until a terminal state with a verdict is written. See `submissions_state_verdict_chk` below. **Exception to section 2** (enums are `NOT NULL`): `CON-01` requires the verdict to be `null` while `QUEUED`/`RUNNING`. |
| `attempt` | `INTEGER` | NO | `0` | Number of times an evaluation was started: incremented on each `QUEUED` → `RUNNING`, and when a `RUNNING` submission is taken again after a redelivered message (crash recovery, `C05-OQ-27`). `CHECK (attempt >= 0)`. Upper bound: `C05-OQ-24`. Fencing value for terminal writes (see *Attempt fencing* below). |
| `created_at` | `TIMESTAMPTZ` | NO | `now()` | Time the submission was received. Maps to `submitted_at` in the `CON-02` message. |
| `updated_at` | `TIMESTAMPTZ` | NO | `now()` | Set by the application on every `UPDATE` (section 2). |

**Constraint `submissions_state_verdict_chk`** (CON-01 §1 "State-Verdict Invariants"):

```sql
CONSTRAINT submissions_state_verdict_chk CHECK (
     (state IN ('QUEUED', 'RUNNING', 'CANCELLED') AND verdict IS NULL)
  OR (state = 'COMPLETED' AND verdict IN (
        'ACCEPTED', 'WRONG_ANSWER', 'COMPILATION_ERROR', 'RUNTIME_ERROR',
        'TIME_LIMIT_EXCEEDED', 'MEMORY_LIMIT_EXCEEDED', 'OUTPUT_LIMIT_EXCEEDED'))
  OR (state = 'FAILED' AND verdict = 'INTERNAL_ERROR')
)
```

| `state` | `verdict` |
|---|---|
| `QUEUED`, `RUNNING` | `NULL` |
| `COMPLETED` | one of the 7 test verdicts |
| `FAILED` | `INTERNAL_ERROR` |
| `CANCELLED` | `NULL` (`CON-01` defines no verdict for `CANCELLED`) |

The same constraint also limits `state` and `verdict` to the `CON-01` names, so no separate `CHECK` per column is needed.

**Index `submissions_one_active_per_user`** (proposed default for `C05-OQ-20`):

```sql
CREATE UNIQUE INDEX submissions_one_active_per_user
  ON submissions (user_id)
  WHERE state IN ('QUEUED', 'RUNNING');
```

A user can have at most one active (`QUEUED` or `RUNNING`) submission across all problems. A second `INSERT` fails with `unique_violation` (SQLSTATE `23505`), which the API maps to `409 Conflict` (`CON-01` default). The retry flow `RUNNING → QUEUED` keeps the row active, so it does not conflict with itself.

**Rules for state changes (CON-01 §1):**
- Allowed transitions: `QUEUED → RUNNING`, `RUNNING → COMPLETED | FAILED`, `QUEUED | RUNNING → CANCELLED`, `RUNNING → QUEUED` (infrastructure retry). Taking a `RUNNING` submission again after an Orchestrator crash is not defined in `CON-01`: `C05-OQ-27`. The DB checks valid (`state`, `verdict`) pairs but not transitions; the application enforces transitions with a conditional update.
- **Taking a submission (claim).** The Orchestrator takes a submission with one atomic statement that also counts the attempt (`CON-02` §4):
  ```sql
  UPDATE submissions
  SET state = 'RUNNING', attempt = attempt + 1, updated_at = now()
  WHERE id = $id AND state IN ('QUEUED', 'RUNNING')
  RETURNING attempt;
  ```
  The returned `attempt` is this run's **fencing value**. 0 rows = the submission is already terminal: skip and acknowledge.
- **Attempt fencing.** Every later write of this run (terminal state, `RUNNING → QUEUED`) must match both `state = 'RUNNING'` and the fencing value. A new claim increments `attempt`, so a result from an older, stale attempt affects 0 rows and is ignored:
  ```sql
  UPDATE submissions
  SET state = 'COMPLETED', verdict = $verdict, updated_at = now()
  WHERE id = $id AND state = 'RUNNING' AND attempt = $attempt;
  ```
  `FAILED` / `INTERNAL_ERROR` uses the same condition.
- *First terminal write wins* follows from the same condition: once the state is terminal, `state = 'RUNNING'` is false and every later write affects 0 rows.
- `CANCELLED` is a user action, not part of a run, so it is not fenced: `WHERE id = $id AND state IN ('QUEUED', 'RUNNING')`. A run that finishes after the cancellation affects 0 rows.
- If a `results` row is written (3.4), it is inserted in the same transaction as the terminal update, only when that update affected 1 row, with `results.attempt` = the fencing value.

### 3.4 `results` *(Owner: Yulia)*

Immutable evaluation report of a finished submission: inserted once, never updated or deleted by the application (`created_at` only, no `updated_at`). The verdict is **not** duplicated here; it lives in `submissions.verdict` (3.3).

| Field | Type | Nullable | Default | Constraints / Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | — (see `C05-OQ-04`) | Primary key. |
| `submission_id` | `UUID` | NO | — | FK → `submissions.id`. `UNIQUE`: at most one report per submission. `ON DELETE` follows `C05-OQ-22`. |
| `attempt` | `INTEGER` | NO | — | Attempt that produced this report (copy of `submissions.attempt` at write time). `CHECK (attempt >= 1)`. |
| `created_at` | `TIMESTAMPTZ` | NO | `now()` | Time the report was written. |
| *(details)* | — | — | — | Per-test results, output, duration and memory usage: **not defined in this draft**, see `C05-OQ-26`. |

A row exists only for `COMPLETED` and `FAILED` submissions. Whether it is required for every such submission is part of `C05-OQ-26`. A `CHECK` cannot read another table, so this link is enforced by the application, not by the DB.

---

## 4. Open Questions

| ID | Question | Affected Entity / Field | Source / Owner | Blocks Final Approval? |
|---|---|---|---|---|
| C05-OQ-01 | Role source and assignment mechanism (`SPEC-F01` OQ-1) | `users`, `roles` | `SPEC-F01` / Yevhen | Yes |
| C05-OQ-02 | User profile / auth fields scope (`SPEC-F01` OQ-2) | `users` | `SPEC-F01` / Yevhen | Yes |
| C05-OQ-03 | User model details (`SPEC-F01` OQ-9) | `users` | `SPEC-F01` / Yevhen | Yes |
| C05-OQ-04 | UUID generation strategy: DB-side (`gen_random_uuid()`, v4) vs. Go-side (`UUIDv7`) | All entities (`id`) | `INF-05` / Yevhen | No |
| C05-OQ-20 | Duplicate submission while a previous one is `QUEUED`/`RUNNING` (`F-5`). **Proposed default:** one active submission per user across all problems, enforced by the partial unique index `submissions_one_active_per_user`; the API returns `409 Conflict` / `submission_in_progress` (`CON-01` default). `SPEC-F06` (FR-10, OQ-01) currently describes the case per user **and problem**: the scope must be aligned with Pair A, then confirmed by the mentor. | `submissions` | `CON-01`, `SPEC-F06` / Yulia | Yes |
| C05-OQ-21 | Exact type and FK constraint for `problem_id` referencing `problems.id` | `submissions.problem_id` | `CON-04` / Yulia | Yes |
| C05-OQ-22 | `ON DELETE` policy for `submissions.user_id` (and so for `results`): `RESTRICT` (keep history, block user deletion) or `CASCADE`. Depends on whether users can be deleted (`SPEC-F01`). | `submissions.user_id`, `results.submission_id` | `SPEC-F01` / Yulia | Yes |
| C05-OQ-23 | Maximum size of `source_code`, and whether the DB also enforces it with `CHECK (octet_length(source_code) <= N)` or only the API does (`SPEC-F06` FR-04, OQ-04: `413 Content Too Large`). | `submissions.source_code` | `SPEC-F06` / Yulia | No |
| C05-OQ-24 | Maximum number of attempts. **Proposed resolution:** 3 attempts in total, counted in `submissions.attempt`; above 3 the submission becomes `FAILED` / `INTERNAL_ERROR` (`CON-02` §4, §6 in PR #38, aligned with `CON-01`). Closes when `CON-02` is merged. | `submissions.attempt` | `CON-01`, `CON-02` / Yulia | No |
| C05-OQ-25 | ~~Should a terminal write be accepted only from the current attempt?~~ **Resolved (20.11, mentor review on #40):** yes. All writes of a run are fenced by `attempt` (section 3.3, *Attempt fencing*). | `submissions.attempt`, `results.attempt` | `CON-01` / Yulia | — |
| C05-OQ-26 | Which evaluation details `results` stores (per-test results, output, duration, memory), in what form (columns, JSONB or a child table), and whether a row is required for every `COMPLETED`/`FAILED` submission. Depends on test visibility rules. | `results` | `CON-04`, `SPEC-F08` / Yulia | No |
| C05-OQ-27 | Crash recovery: `CON-01` defines retry only as `RUNNING → QUEUED` "triggered by the orchestrator". After an Orchestrator crash nobody makes this transition, and the redelivered message finds the submission still `RUNNING` (`CON-02` §4). Is taking a `RUNNING` submission again an allowed transition (with `attempt + 1`), and does it count toward the attempt bound (`C05-OQ-24`)? Needs a `CON-01` update. | `submissions.state`, `submissions.attempt` | `CON-01`, `CON-02` / Yulia | No |