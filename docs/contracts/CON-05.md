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
- **Unresolved Constraints:** Any field type, length, or `ON DELETE` policy dependent on unfinalized specs is explicitly recorded in Section 4 (`Open Questions`).
- **Open Question IDs:** `C05-OQ-1`..`C05-OQ-19` are reserved for `users`/`roles`/shared conventions; `C05-OQ-20`+ are reserved for `submissions`/`results`.

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
*Draft in progress. States and verdicts must match `CON-01` by name.*

| Field | Type | Nullable | Default | Constraints / Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | — | Primary Key |

### 3.4 `results` *(Owner: Yulia)*
*Draft in progress. Verdicts must match `CON-01` by name. Immutable record (`created_at` only, no `updated_at`).*

| Field | Type | Nullable | Default | Constraints / Description |
|---|---|---|---|---|
| `id` | `UUID` | NO | — | Primary Key |

---

## 4. Open Questions

| ID | Question | Affected Entity / Field | Source / Owner | Blocks Final Approval? |
|---|---|---|---|---|
| C05-OQ-01 | Role source and assignment mechanism (`SPEC-F01` OQ-1) | `users`, `roles` | `SPEC-F01` / Yevhen | Yes |
| C05-OQ-02 | User profile / auth fields scope (`SPEC-F01` OQ-2) | `users` | `SPEC-F01` / Yevhen | Yes |
| C05-OQ-03 | User model details (`SPEC-F01` OQ-9) | `users` | `SPEC-F01` / Yevhen | Yes |
| C05-OQ-04 | UUID generation strategy: DB-side (`gen_random_uuid()`, v4) vs. Go-side (`UUIDv7`) | All entities (`id`) | `INF-05` / Yevhen | No |
| C05-OQ-20 | Duplicate submission handling policy while previous is `QUEUED`/`RUNNING` (`F-5`) | `submissions` | `CON-01` / Yulia | Yes |
| C05-OQ-21 | Exact type and FK constraint for `problem_id` referencing `problems.id` | `submissions.problem_id` | `CON-04` / Yulia | Yes |
