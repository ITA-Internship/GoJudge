# Feature Spec: F08 – Submission History

---

## 1. Feature Metadata *(Required)*

| Field | Value |
|-------|-------|
| Feature ID | F08 |
| Title | Submission History |
| Status | In Review |
| Responsible pair | Pair A (yanferens, xhanjo) | 2026-11-19 | Submitted |
| Reviewer pair | Pair C (provisional, to be confirmed by the mentor) |
| GitHub Issue | [#30](https://github.com/ITA-Internship/GoJudge/issues/30) |
| Last updated | 2026-11-19 |

---

## 2. Overview and User Story *(Required)*

**Purpose:** A learner needs to see what happened to the solutions they sent. This feature gives each learner a list of their own submissions, with the time each one was sent, its current state and, when it is available, its final verdict.

**User story:**

> As a *learner*,
> I want *to see the history of my own submissions with their time, state and final verdict*,
> so that *I can follow the progress of my work and know whether my solutions were accepted*.

**Expected value:** The learner does not need to remember or re-send a solution to learn its result. The history is the single place where the outcome of every past submission can be found. It also shows the learner which submissions are still waiting or running. The same data is later used by F02-04 to mark problems as solved in the catalogue.

---

## 3. Scope *(Required)*

**In scope:**

- A paginated list of the signed-in user's own submissions (backend endpoint, backlog item F08-01).
- For each submission: its ID, the problem it belongs to, the time it was sent, the time of its last state change, its state and its final verdict.
- States and verdicts exactly as named in [CON-01](../../docs/contracts/CON-01.md). States: `QUEUED`, `RUNNING`, `COMPLETED`, `FAILED`, `CANCELLED`. Verdicts: `ACCEPTED`, `WRONG_ANSWER`, `COMPILATION_ERROR`, `RUNTIME_ERROR`, `TIME_LIMIT_EXCEEDED`, `MEMORY_LIMIT_EXCEEDED`, `OUTPUT_LIMIT_EXCEEDED`, and `INTERNAL_ERROR` for `FAILED`.
- An empty (`null`) verdict while a submission is `QUEUED` or `RUNNING`.
- A history page in the web application (backlog item F08-02).

**Out of scope:**

- Submissions of other users, for any role.
- Creating, cancelling or re-sending a submission (F06).
- Creating a submission (`POST /api/v1/problems/{problemId}/submissions`) and reading one submission (`GET /api/v1/submissions/{submissionId}`), with its polling in the editor (F06-01, F06-03, F06-04).
- Producing verdicts (the evaluation itself, F07-05).
- Viewing the submitted source code.
- Per-test results, program output, run time and memory usage (pending CON-04 visibility rules and CON-05 `C05-OQ-26`, see OQ-04).
- Filtering, searching and sorting options beyond the default order (newest first).
- Administrator or instructor views of all submissions.
- Notifications or live push updates when a verdict becomes available.

---

## 4. Functional Requirements *(Required)*

Each requirement must be numbered and observable, so that a tester can confirm it.

**Backend (F08-01)**

| ID | Requirement |
|----|-------------|
| FR-01 | The system shall provide `GET /api/v1/submissions`, which returns the submissions owned by the caller. The path is the collection of the F06 endpoint `GET /api/v1/submissions/{submissionId}`. |
| FR-02 | The system shall determine the owner only from the validated access token, through the user ID that F01 resolves. Any user ID, owner or role value sent in the query string, headers or body shall be ignored. |
| FR-03 | The system shall never return a submission owned by another user. This applies to learners and administrators alike. |
| FR-04 | For each submission, the system shall return `id`, `problem_id`, `state`, `verdict`, `created_at` and `updated_at`, with the same names, types and meaning as the submission object of SPEC-F06 (field definitions in §7.1). |
| FR-05 | The system shall return `state` and `verdict` values exactly as named in CON-01 and shall not rename, translate or merge them. |
| FR-06 | The system shall return `verdict` as `null` while the state is `QUEUED` or `RUNNING`, and for `CANCELLED`. For `COMPLETED` it shall return one of the seven test verdicts, and for `FAILED` it shall return `INTERNAL_ERROR`. |
| FR-07 | The system shall only read the verdict stored for the submission and shall never compute or guess one. Until F07-05 writes verdicts, the field is `null`, and F08-01 is still complete. |
| FR-08 | The system shall order the list by `created_at`, newest first, and by `id` descending when two submissions have the same `created_at`. |
| FR-09 | The system shall return results in pages. The page size is set by `page_size`: unset or `0` means the default of 20, values above 100 are lowered to 100. A `next_page_token` is returned when more results exist, and an empty string on the last page. |
| FR-10 | The system shall reject a negative or non-integer `page_size`, or a `page_token` that was not issued by the server, with `400` and shall not return any submissions in that case. |
| FR-11 | The system shall return `200` with an empty `submissions` list when the caller has no submissions. |
| FR-12 | The system shall not include source code, hidden test data, other users' data, the identity provider subject, the email or any token in the response. |
| FR-13 | The system shall respond `401` to a request without a valid access token (F01 FR-10), with no submission data. |
| FR-14 | The endpoint shall be read-only. A call shall not change any submission, state or verdict. |

**Web application (F08-02)**

| ID | Requirement |
|----|-------------|
| FR-15 | The web application shall provide a "My submissions" page, reachable only by a signed-in user. |
| FR-16 | The page shall show one row per submission with the problem, the time it was sent in the user's local time zone, the state and the verdict. |
| FR-17 | The page shall show "Pending" in the verdict column for `QUEUED` and `RUNNING`, and "—" for `CANCELLED` and for any other `null` verdict. |
| FR-18 | The page shall show a clear message and a link to the problem catalogue when the user has no submissions. |
| FR-19 | The page shall provide a "Load more" control that loads the next page with `next_page_token`, and hide it when `next_page_token` is empty. |
| FR-20 | The page shall provide a "Refresh" control that reloads the list from the first page, so that the user can see state changes. |
| FR-21 | The page shall show an error message with a "Try again" control when the request fails, and shall not show partial or stale data as if it were current. |

---

## 5. Acceptance Criteria *(Required)*

### AC-01 – Learner sees only their own submissions

- **Given** learner A has 3 submissions and learner B has 1 submission
- **When** learner A calls `GET /api/v1/submissions`
- **Then** the response is `200` with exactly A's 3 submissions, and B's submission is not in it

### AC-02 – Owner cannot be changed by the request

- **Given** learner A and learner B each have submissions
- **When** learner A calls `GET /api/v1/submissions?user_id=<B's ID>`, or sends B's ID in a header or body
- **Then** the response contains only A's submissions, exactly as without the extra value

### AC-03 – Administrator sees only their own submissions

- **Given** a user with role `administrator` has 1 submission and learner B has 2 submissions
- **When** the administrator calls `GET /api/v1/submissions`
- **Then** the response contains only the administrator's own submission

### AC-04 – Fields of a submission

- **Given** a learner has one submission
- **When** they call `GET /api/v1/submissions`
- **Then** the item contains `id`, `problem_id`, `state`, `verdict`, `created_at` and `updated_at`, the same fields as `GET /api/v1/submissions/{submissionId}` (F06) returns for that submission, and contains no `source_code`, `user_id`, email, identity provider subject or token

### AC-05 – Verdict is null while the submission is not finished

- **Given** a learner has one submission in `QUEUED` and one in `RUNNING`
- **When** they call `GET /api/v1/submissions`
- **Then** both items have `verdict: null` and `state` equal to `QUEUED` and `RUNNING`

### AC-06 – Verdicts of finished submissions

- **Given** a learner has a `COMPLETED` submission with verdict `WRONG_ANSWER`, a `FAILED` submission with verdict `INTERNAL_ERROR` and a `CANCELLED` submission
- **When** they call `GET /api/v1/submissions`
- **Then** the items show `WRONG_ANSWER`, `INTERNAL_ERROR` and `null` respectively, with names exactly as in CON-01

### AC-07 – No verdicts before F07-05

- **Given** verdicts are not yet written by the judge (before F07-05), so every stored verdict is `null`
- **When** a learner calls `GET /api/v1/submissions`
- **Then** the response is `200`, every item has `verdict: null`, and no verdict is invented

### AC-08 – Newest first

- **Given** a learner sent submissions at 10:00, 10:05 and 10:10
- **When** they call `GET /api/v1/submissions`
- **Then** the items are in the order 10:10, 10:05, 10:00

### AC-09 – Pages

- **Given** a learner has 25 submissions
- **When** they call `GET /api/v1/submissions?page_size=10`, then call it again with each returned `next_page_token` as `page_token`
- **Then** they receive pages of 10, 10 and 5 items, the last page has an empty `next_page_token`, and every submission appears exactly once

### AC-10 – Page size and page token rules

- **Given** a learner has 130 submissions
- **When** they call the endpoint without `page_size`, with `page_size=0`, with `page_size=500`, and then with `page_size=-1`, `page_size=abc` or a `page_token` that was not issued by the server
- **Then** the first two calls return 20 items, the third returns 100 items, and each of the last three calls returns `400` with code `validation_error` and no items

### AC-11 – Empty history

- **Given** a signed-in learner who has never submitted
- **When** they call `GET /api/v1/submissions` and open the "My submissions" page
- **Then** the API returns `200` with `submissions: []` and `next_page_token: ""`, and the page shows the empty-history message with a link to the catalogue

### AC-12 – Not signed in

- **Given** a request with no access token or with an invalid one
- **When** it calls `GET /api/v1/submissions`
- **Then** the response is `401` with code `unauthorized` and contains no submission data

### AC-13 – History does not change data

- **Given** a learner has a `RUNNING` submission
- **When** they call `GET /api/v1/submissions` several times
- **Then** the submission's `state`, `verdict` and `updated_at` in the database are unchanged by these calls

### AC-14 – History page shows states and verdicts

- **Given** a signed-in learner with submissions in `QUEUED`, `COMPLETED` (`ACCEPTED`) and `CANCELLED`
- **When** they open the "My submissions" page
- **Then** each row shows the problem, the local send time and the state, and the verdict column shows "Pending", "ACCEPTED" and "—" respectively

### AC-15 – Refresh shows a new state

- **Given** the history page shows a submission as `RUNNING`, and the submission then reaches `COMPLETED` with `ACCEPTED`
- **When** the learner presses "Refresh"
- **Then** the row shows `COMPLETED` and `ACCEPTED`

### AC-16 – Load more

- **Given** a learner has 25 submissions
- **When** they open the history page and press "Load more"
- **Then** the page first shows 20 rows, then 25 rows with no duplicates, and the "Load more" control is hidden

### AC-17 – Error on the page

- **Given** the API is unavailable
- **When** the learner opens the history page
- **Then** the page shows an error message with "Try again" and no list rows

---

## 6. Technical Design *(Required)*

**Component responsibilities:**

| Component | Responsibility |
|-----------|----------------|
| Web application – "My submissions" page (F08-02) | Calls `GET /api/v1/submissions` with the access token, shows the rows, "Pending" and "—" placeholders, local time, "Load more", "Refresh", empty and error states. Never sends a user ID. |
| Go API – authentication and role middleware (F01) | Validates the token, resolves the caller's `users.id` into the request context. The route is registered with the "authenticated" policy (any role). |
| Go API – history handler (F08-01) | Validates `page_size` and `page_token`, reads the caller's ID only from the request context, calls the data access layer and builds the response. |
| Data access (INF-05) | One read-only query on `submissions`, filtered by `user_id`, ordered by `created_at DESC, id DESC`, using keyset pagination (sketch below). |
| PostgreSQL – `submissions` table (CON-05 §3.3) | Stores state, verdict, `created_at` and `updated_at`. Its `submissions_state_verdict_chk` constraint keeps state and verdict consistent with CON-01. |

**Query sketch** (for illustration, the final query belongs to F08-01):

```sql
SELECT id, problem_id, state, verdict, created_at, updated_at
FROM submissions
WHERE user_id = $user_id
  AND ($token_created_at IS NULL OR (created_at, id) < ($token_created_at, $token_id))
ORDER BY created_at DESC, id DESC
LIMIT $page_size + 1;
```

One extra row is read to know whether `next_page_token` is needed. The page token is an opaque, URL-safe string that encodes the last item's `created_at` and `id` and is encrypted with a server key, so that the client cannot read, change or build it. Keyset pagination is chosen over page numbers so that a new submission sent while the learner is paging does not shift rows or show a row twice. An index on `submissions (user_id, created_at DESC, id DESC)` is proposed to CON-05 (OQ-06).

**Dependencies:**

- F01 (F01-02): authentication, `users.id` of the caller, `401` behaviour and the error envelope.
- F06 (F06-01, F06-03): creates the submissions that F08 lists, and defines the submission object and the paths `POST /api/v1/problems/{problemId}/submissions` and `GET /api/v1/submissions/{submissionId}` that F08 follows (SPEC-F06 draft, PR #39).
- F07 (F07-05): writes the verdicts that F08 shows. F08-01 does not wait for it (backlog F08-01).
- CON-01: state and verdict names and the state–verdict invariants.
- CON-05 (preliminary, PR #40): `submissions` table. Field names in this spec follow its current draft.
- CON-03: the OpenAPI document must include the endpoint in §7.1.
- INF-05: database access layer (pgx, sqlc).

**Architecture decisions:**

- None. No new architectural decision is introduced. Paths and field names follow SPEC-F06. Parameter names for pagination (`page_size`, `page_token`, `next_page_token`) follow Google AIP-158. Keyset pagination is a local design choice of this endpoint and is recorded as OQ-02 for agreement.

---

## 7. Contracts *(Conditional)*

### 7.1 API Endpoints *(if applicable)*

| Method | Path | Purpose | Auth required |
|--------|------|---------|---------------|
| GET | `/api/v1/submissions` | List the caller's own submissions, newest first | Yes (any role) |

All submission endpoints, so that the paths stay consistent with SPEC-F06:

| Method | Path | Purpose | Owner |
|--------|------|---------|-------|
| POST | `/api/v1/problems/{problemId}/submissions` | Create a submission for a problem | F06-01 |
| GET | `/api/v1/submissions/{submissionId}` | Read one submission and its current state | F06-03 |
| GET | `/api/v1/submissions` | List the caller's own submissions | F08-01 |

**Request schema:**

- Header `Authorization: Bearer <access token>` (F01).
- Query parameters (names as in Google AIP-158):

| Parameter | Type | Required | Rules |
|-----------|------|----------|-------|
| `page_size` | integer | No | Unset or `0`: default 20. Above 100: lowered to 100. Negative or not an integer: `400`. |
| `page_token` | string | No | Value of `next_page_token` from an earlier response. Any other value gives `400`. |

- No body. Any other query parameter (for example `user_id` or `problem_id`) is ignored.

**Response schema:**

Success, `200`. Each element of `submissions` is the same submission object as in SPEC-F06:

```json
{
  "submissions": [
    {
      "id": "0b6f4c7e-1f2a-4b8e-9d3c-5a6b7c8d9e0f",
      "problem_id": "3c2b1a09-8f7e-4d6c-b5a4-0f1e2d3c4b5a",
      "state": "COMPLETED",
      "verdict": "ACCEPTED",
      "created_at": "2026-11-18T10:10:00Z",
      "updated_at": "2026-11-18T10:10:07Z"
    },
    {
      "id": "7d8e9f00-1a2b-4c3d-8e4f-5a6b7c8d9e0a",
      "problem_id": "3c2b1a09-8f7e-4d6c-b5a4-0f1e2d3c4b5a",
      "state": "QUEUED",
      "verdict": null,
      "created_at": "2026-11-18T10:05:00Z",
      "updated_at": "2026-11-18T10:05:00Z"
    }
  ],
  "next_page_token": ""
}
```

| Field | Type | Source (CON-05) | Notes |
|-------|------|-----------------|-------|
| `submissions` | array | — | Empty array when there are no submissions. |
| `id` | UUID string | `submissions.id` | Submission ID, the same as returned by F06 and used in `GET /api/v1/submissions/{submissionId}`. |
| `problem_id` | UUID string | `submissions.problem_id` | Type depends on CON-04 (`C05-OQ-21`). Problem title: OQ-05. |
| `state` | string | `submissions.state` | One of the CON-01 states. |
| `verdict` | string or `null` | `submissions.verdict` | One of the CON-01 verdicts, or `null` (FR-06). |
| `created_at` | RFC 3339 UTC string | `submissions.created_at` | Time the submission was received. |
| `updated_at` | RFC 3339 UTC string | `submissions.updated_at` | Time of the last state change. |
| `next_page_token` | string | — | Opaque. Empty string on the last page. |

Errors use the F01 envelope (F01 §7.1, F01 OQ-7):

| Status | `code` | When |
|--------|--------|------|
| `400` | `validation_error` | Negative or non-integer `page_size`, unknown `page_token` (AC-10) |
| `401` | `unauthorized` | Missing or invalid token (AC-12) |
| `503` | `service_unavailable` | Database unavailable |

```json
{ "error": { "code": "validation_error", "message": "Invalid query parameter: page_size." } }
```

### 7.2 Data Models *(if applicable)*

F08 adds no table and no column. It reads `submissions` as defined in CON-05 §3.3 (preliminary): `id`, `user_id`, `problem_id`, `state`, `verdict`, `created_at`, `updated_at`. It does not read `source_code` or `attempt`, and it does not read `results` in the MVP (OQ-04). One index is proposed to CON-05 (OQ-06).

### 7.3 Messaging Contracts *(if applicable)*

Not applicable. F08 neither publishes nor consumes messages. It only reads the state that other components store.

### 7.4 State Transitions *(if applicable)*

Not applicable. F08 changes no state. It shows the states and transitions defined in CON-01 §1 as they are stored. Because of the infrastructure retry flow (`RUNNING` → `QUEUED`), a learner may see a submission go back from `RUNNING` to `QUEUED` between two refreshes. This is expected and is not an error.

---

## 8. Edge Cases and Failure Handling *(Conditional)*

- [x] Invalid or missing input
- [x] Timeouts
- [x] Retries and their limits
- [x] Duplicate processing
- [x] Unavailable dependencies

| Scenario | Expected behaviour |
|----------|--------------------|
| No token or invalid token | `401`, no data (AC-12). |
| `page_size` unset or `0` | Default 20. |
| `page_size` above 100 | Lowered to 100 (AC-10). |
| `page_size` negative or not an integer | `400 validation_error` (AC-10). |
| `page_token` not issued by the server, damaged or in a wrong format | `400 validation_error`. No data. |
| `page_token` issued to another user | Cannot reveal data: the query always uses the caller's ID. |
| Unknown query parameter, including `user_id` | Ignored (AC-02). |
| The learner has no submissions | `200`, `submissions: []`, `next_page_token: ""` (AC-11). |
| A new submission is created while the learner is paging | It does not appear in the following pages and no row is repeated. It appears after "Refresh" (keyset pagination, §6). |
| A submission changes state between two pages | Each page shows the state at the time it was read. "Refresh" shows the current state. |
| A submission goes from `RUNNING` back to `QUEUED` (infrastructure retry) | Shown as stored (§7.4). |
| A submission stays `RUNNING` after an Orchestrator crash until it is taken again (CON-05 `C05-OQ-27`) | Shown as `RUNNING` with `verdict: null`. F08 does not detect or repair stuck submissions. |
| The problem of a past submission was later unpublished | The submission is still listed. History shows the learner's own data, and the published check of SPEC-F06 FR-03 applies only when a submission is created. |
| Two submissions have the same `created_at` | Order by `id` descending, so the order is stable between calls (FR-08). |
| Verdicts not yet written by F07-05 | `verdict: null` for every item (AC-07). |
| A stored state or verdict is not a CON-01 name | Cannot happen while the CON-05 check constraint is in place. If it still occurs, the API returns `500` and logs an error with the submission ID, and does not pass the unknown value to the client. |
| Database unavailable or the query times out | `503 service_unavailable`, no partial list. The page shows the error with "Try again" (AC-17). |
| Repeated calls | Safe, the endpoint is read-only (AC-13). Retries are left to the user ("Try again"). The web application does not retry automatically. |
| Very long history | Only one page (at most 100 items) is read per request. |

---

## 9. Security Considerations *(Required)*

- [x] **Authentication:** The caller is identified by the F01 access token, validated by the F01 middleware on every request. The route uses the "authenticated" policy. Without a valid token the response is `401` (AC-12).
- [x] **Authorisation:** A user can read only their own submissions. The owner filter `user_id = <caller's users.id>` is always applied in the query and is taken only from the request context set by F01. No request field can change it (FR-02, AC-02). Administrators get no wider access through this endpoint (FR-03, AC-03). Any future view of other users' submissions requires a separate, administrator-only endpoint and specification. The page token only encodes a position in the caller's own list and grants no access, and the owner filter is applied to every page, so a forged or copied token cannot reveal other users' data. For the F06 endpoint `GET /api/v1/submissions/{submissionId}`, another user's submission must give the same response as a missing one (OQ-03).
- [x] **Input validation:** `page_size` is an integer, not negative, lowered to 100 when larger. `page_token` is opaque, encrypted by the server and checked before use. The query is parameterised (sqlc), and no request value is placed in SQL text. Other parameters are ignored.
- [x] **Sensitive data:** The response contains no source code, test data, hidden test details, email, identity provider subject or token (FR-12, AC-04). Logs for this endpoint contain the user ID and request ID (INF-07), not tokens or submission content.
- [ ] **Execution isolation:** Not applicable. F08 does not run learner-submitted code.

---

## 10. Verification Plan *(Required)*

| AC | Verification type | Test or scenario | Automated? |
|----|-------------------|------------------|------------|
| AC-01 | Integration | Seed 3 submissions for A and 1 for B. A's token returns exactly A's 3 IDs. | Yes |
| AC-02 | Integration | A's token with `user_id=<B>` in the query, a header and a body: same result as without it. | Yes |
| AC-03 | Integration | Administrator token returns only the administrator's own submission. | Yes |
| AC-04 | Integration | Check the item keys: the six expected fields are present and no forbidden field is present. Compare one item with the response of `GET /api/v1/submissions/{submissionId}`. | Yes |
| AC-05 | Integration | Seed `QUEUED` and `RUNNING` submissions, assert `verdict` is `null`. | Yes |
| AC-06 | Integration | Seed `COMPLETED`/`WRONG_ANSWER`, `FAILED`/`INTERNAL_ERROR` and `CANCELLED`, assert the values and exact names. | Yes |
| AC-07 | Integration | Seed only rows with `verdict` `null` (states `QUEUED`, `RUNNING`, `CANCELLED`), assert `200` and `null` everywhere. | Yes |
| AC-08 | Integration | Seed three rows with different `created_at`, assert newest-first order. Add two rows with equal `created_at` and assert ordering by `id`. | Yes |
| AC-09 | Integration | Seed 25 rows, page with `page_size=10`, assert 10/10/5, last `next_page_token` empty, no duplicates or gaps. | Yes |
| AC-10 | Unit and integration | Table test: unset and `0` give 20, `500` gives 100, `-1`, `abc` and a forged token give `400 validation_error`. | Yes |
| AC-11 | Integration and end-to-end | API returns empty list for a new user. The page shows the empty message and catalogue link. | Yes (API), No (page) |
| AC-12 | Integration | No header, and an expired token: `401 unauthorized`, no data. | Yes |
| AC-13 | Integration | Read a `RUNNING` row, call the endpoint 3 times, assert the DB row is unchanged. | Yes |
| AC-14 | Component (UI) | Render the page with mocked API data, assert "Pending", "ACCEPTED", "—" and local time. | Yes |
| AC-15 | Component (UI) | Mock first response `RUNNING`, second `COMPLETED`/`ACCEPTED`, press "Refresh", assert the new row. | Yes |
| AC-16 | Component (UI) | Mock 20 + 5 items with a page token, press "Load more", assert 25 rows and hidden control. | Yes |
| AC-17 | Component (UI) | Mock a failed request, assert the error message, "Try again" and no rows. | Yes |

The full journey (submit, then see the submission in the history) is covered by TST-02.

---

## 11. Dependencies and Open Questions *(Required)*

**Related features:**

- **F01 – Accounts and access.** Provides authentication, the caller's `users.id`, `401` and the error envelope. F08 relies on F01-02 and does not re-implement authentication.
- **F06 – Submit solution.** Creates the submissions F08 lists (`POST /api/v1/problems/{problemId}/submissions`), owns `GET /api/v1/submissions/{submissionId}` (F06-03) and the polling in the editor (F06-04), and defines the submission object that F08 reuses.
- **F07 – Automated judge.** F07-05 writes the verdicts F08 shows. F08-01 does not wait for it.
- **F02 – Problem catalogue.** F02-04 uses F08-01 to mark solved problems.
- **CON-01 – Submission lifecycle and verdict model.** Names of states and verdicts, invariants.
- **CON-04 – Problem packaging and test visibility.** Decides the `problem_id` type and whether any test details may be shown to learners.
- **CON-05 – Data model (preliminary).** `submissions` and `results` tables.
- **CON-03 – OpenAPI contract.** Must include the endpoint in §7.1.

**Open questions:**

| # | Question | Owner | Status |
|---|----------|-------|--------|
| OQ-01 | `CANCELLED` has no verdict in CON-01. What does the history show in the verdict column? | Pair A, Pair D | Resolved: `verdict` is `null` for `CANCELLED`, as in SPEC-F06 FR-08 and §7.4 (PR #39) and CON-05 `submissions_state_verdict_chk` (PR #40). The page shows "—" (FR-17). |
| OQ-02 | Pagination: confirm keyset pagination with an opaque page token, default 20 and maximum 100, and the response shape (`submissions`, `next_page_token`) for CON-03. Should SPEC-F06 and other list endpoints use the same parameters? | Pair A, CON-03 owner | Open |
| OQ-03 | SPEC-F06: `GET /api/v1/submissions/{submissionId}` for another user's submission should return `404`, the same as for a missing one, so that its existence is not revealed. | Pair A (SPEC-F06) | Resolved: SPEC-F06 FR-09 and AC-06 (PR #39) return `404 not_found` for another user's submission. |
| OQ-04 | Should the history show run time, memory or per-test details from `results`? Depends on CON-04 visibility rules and CON-05 `C05-OQ-26`. Default for the MVP: no. | Pair A (CON-04), Pair D | Open |
| OQ-05 | Should each item include the problem title (a join with `problems`, or a separate call by the web application)? Blocked by CON-04 and the `problems` entity in CON-05. Default: `problem_id` only, the web application shows the title from the catalogue. | Pair A, Pair D | Open |
| OQ-06 | Add the index `submissions (user_id, created_at DESC, id DESC)` to CON-05? | Pair D | Open |
| OQ-07 | Is a manual "Refresh" enough in the MVP, or should the page refresh automatically while some submission is `QUEUED` or `RUNNING`? Default: manual. | Pair A, mentor | Open |
| OQ-08 | `ON DELETE` policy for `submissions.user_id` (`C05-OQ-22`). If users can be deleted with `CASCADE`, their history disappears. Confirm this is acceptable. | Pair D | Open |
| OQ-09 | Error envelope and code names to be agreed for all endpoints (F01 OQ-7). SPEC-F06 (PR #39) and this spec already use the same envelope and codes (`validation_error`, `unauthorized`, `not_found`, `service_unavailable`). | Pair D, Pair A | Open |
| OQ-10 | Is a per-problem history (`GET /api/v1/problems/{problemId}/submissions`, the same path as the F06 create endpoint) needed, for example on the problem page or for F02-04? Default: not in the MVP. | Pair A | Open |

**Blockers:**

- None for this draft. CON-05 (preliminary), CON-04 and CON-03 must be agreed before F08-01 implementation starts.

---

## 12. Review and Approval *(Required)*

| Role | Name or pair | Date | Outcome |
|------|--------------|------|---------|
| Author | Pair A | *YYYY-MM-DD* | *Submitted* |
| Reviewer | Pair C | *YYYY-MM-DD* | *Approved / Changes requested* |

**Approval checklist:**

- [ ] Acceptance criteria are unambiguous and numbered.
- [ ] Every required contract is defined and agreed.
- [ ] Every acceptance criterion has a verification entry.
- [ ] Security considerations are documented.
- [ ] No critical open questions remain.
- [ ] Reviewer is from a different pair.