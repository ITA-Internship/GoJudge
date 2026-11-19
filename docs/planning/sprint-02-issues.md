# Sprint 2 Issue Manifest (16–27 November 2026)

**Status:** Finalised for mentor review. **GitHub Issues for the eight committed items were created on 13 November 2026**, and #23 was updated (links in §1.1). The stretch item ST-01 has no Issue. Issues have no assignees, and pair composition is finalised during Sprint Planning.

**Approved workload (mentor):** committed task workload **17.5 pair-days**. Maximum committed task workload **20 pair-days**. **Uncommitted 2.5 pair-days.** Meetings, reviews, demos and organisational overhead are excluded from these figures.

**Revision history:** revision 1 proposed S2-07 as F01-01 preparation (3.5). Revision 2 replaced it with INF-07 (4.0) and made INF-09 conditional (1.5). Revision 3 (this file) makes INF-07 the approved committed item, moves INF-09 to **conditional Stretch**, accepts the partial-deliverable model, and records the Sprint 3 follow-up in planning documentation only.

**Source of truth:** [Sprint 2 Plan](sprint-02.md). Backlog IDs, types, full-scope sizes and dependencies come from the [Product Backlog](../product/product-backlog.md) (§4, §5, §10.2). Process rules come from [SDD Workflow](../sdd-workflow.md) and [CONTRIBUTING](../../CONTRIBUTING.md).

**GitHub state checked (read-only, 13 November 2026):** open issues: only #23 (CON-02). No issue, branch or pull request exists for INF-07, CON-04, CON-05, SPEC-F06, SPEC-F08, SPK-05, DEC-11 or INF-09. No logging, request-ID or HTTP middleware code exists on `main` or on any branch.

**Repository state checked for INF-07 (13 November 2026):** `cmd/gojudge/main.go` prints a greeting. Its test checks that greeting. `internal/` and `pkg/` contain only `.keep` files. Module `github.com/ITA-Internship/GoJudge`, Go 1.27.3. CI runs `go build`, `golangci-lint` (`govet`, `errcheck`, `staticcheck`, `ineffassign`, `gosec`, `gofmt`, `goimports`) and `go test -race ./...`.

**Links:** paths are relative to this file, `docs/planning/`.

---

## 1. Summary

### 1.1 Committed issues (17.5 pair-days)

| Ref | Backlog ID | Proposed title | Type | Status | Existing issue | Pair (provisional) | Pair-days |
|-----|------------|----------------|------|--------|----------------|--------------------|-----------|
| S2-01 ([#29](https://github.com/ITA-Internship/GoJudge/issues/29)) | SPEC-F06 (draft) | [SPEC-F06] Submit solution specification (draft) | Feature | NEW | — | A | 1.5 |
| S2-02 ([#30](https://github.com/ITA-Internship/GoJudge/issues/30)) | SPEC-F08 (draft) | [SPEC-F08] Submission history specification (draft) | Feature | NEW | — | A | 1.5 |
| S2-03 ([#31](https://github.com/ITA-Internship/GoJudge/issues/31)) | CON-04 (draft) | [CON-04] Problem packaging and test visibility contract (draft) | Contract | NEW | — | A | 2.0 |
| S2-04 ([#32](https://github.com/ITA-Internship/GoJudge/issues/32)) | SPK-05 | [SPK-05] Isolation mechanism comparison | Spike | NEW | — | B | 3.0 |
| S2-05 ([#23](https://github.com/ITA-Internship/GoJudge/issues/23)) | CON-02 (completion) | [CON-02] Messaging contract (stretch) | Contract | **EXISTING #23** | #23 | B | 1.5 |
| S2-06 ([#33](https://github.com/ITA-Internship/GoJudge/issues/33)) | INF-07 (Sprint 2 subset) | [INF-07] Structured logging and request correlation | Technical task | NEW | — | C | 4.0 |
| S2-07 ([#34](https://github.com/ITA-Internship/GoJudge/issues/34)) | CON-05 (preliminary) | [CON-05] Data model: preliminary contract (users, roles, submissions, results) | Contract | NEW | — | D | 3.0 |
| S2-08 ([#35](https://github.com/ITA-Internship/GoJudge/issues/35)) | DEC-11 | [DEC-11] Database migration tool: ADR-003 (Provisional) | Technical task (ADR) | NEW | — | D | 1.0 |
| | | **Committed total** | | | | | **17.5** |

### 1.2 Conditional Stretch (not committed, not in the 17.5)

| Ref | Backlog ID | Proposed title | Type | Status | Existing issue | Pair | Pair-days |
|-----|------------|----------------|------|--------|----------------|------|-----------|
| ST-01 | INF-09 | [INF-09] Auth0 development tenant configuration | Technical task | **CONDITIONAL STRETCH** (MA-05, explicit mentor authorisation, Pair C capacity) | — | C | 1.5 |

With ST-01 authorised, the total is 19.0, within the 20 maximum. Pair C would reach 5.5, above its 5.0 working figure. See §5 and decision PC-3.

### 1.3 Reconciliation by pair

| Pair | Committed issues | Sum | Approved allocation | Match |
|------|------------------|-----|---------------------|-------|
| A | S2-01 (1.5), S2-02 (1.5), S2-03 (2.0) | 5.0 | 5.0 | Yes |
| B | S2-04 (3.0), S2-05 (1.5) | 4.5 | 4.5 | Yes |
| C | S2-06 (4.0) | 4.0 | 4.0 | Yes. ST-01 (1.5) is not counted. |
| D | S2-07 (3.0), S2-08 (1.0) | 4.0 | 4.0 | Yes |
| **Total** | eight committed issues | **17.5** | 17.5 | Yes. Uncommitted **2.5**. |

No issue is counted twice. Issue #23 (S2-05) is counted once, as Sprint 2 completion work.

### 1.4 Gates and existing issues with no pair-days

| Gate or issue | Status | Action |
|---------------|--------|--------|
| SPEC-F01 formal approval (#10, closed) | In Review. Critical OQ-1 and OQ-2 open. Not in Sprint 2 scope. | Mentor decision on reopening (§6). Gates F01-01 later. |
| SPEC-F04 approval (#24, closed) | In Review. Not in Sprint 2 scope. | Mentor decision on reopening (§6). |
| SonarQube (INF-15; issue #21, now titled "[INF-15] SonarQube…") | INF-15 is the confirmed identifier. The administrative title correction is **resolved**. Not Sprint 2 work. | None. Issue #21 was not modified in this step. |

---

## 2. Conventions

Each issue is one of the four SDD work item types ([CONTRIBUTING §2](../../CONTRIBUTING.md#2-sdd-rules-and-work-item-types)).

| Type | Issues | Required before work | Artefact | Branch | Issue form |
|------|--------|----------------------|----------|--------|------------|
| Feature (specification draft) | S2-01, S2-02 | Specification written from the template. Approval by the reviewer pair, which is not a Sprint 2 closing condition. | `specs/features/` | `docs/` | Feature request (applies `type:feature`) |
| Contract (draft or preliminary) | S2-03, S2-05, S2-07 | Draft naming the dependent pairs | `docs/contracts/CON-XX.md` | `docs/` | Technical task, work type "Contract" |
| Technical task | S2-06, S2-08, ST-01 | Numbered acceptance criteria on the issue | Code and tests under `internal/`, documentation, or ADR in `docs/architecture/adr/` | `infra/` or `docs/` | Technical task form (select `type:technical-task` or `type:adr`) |
| Technical spike | S2-04 | Spike Brief on the issue, before the timebox | Report `docs/spikes/SPK-05.md`. Experiment code on `spike/`, not merged. | `spike/`, `docs/` | Technical spike (applies `type:spike`) |

**Labels: type only (mentor policy).** GitHub labels identify the type of Issue, not its workflow status. Status-oriented labels (`spec:draft`, `spec:in-review`, `spec:approved`, `in-progress`, `in-review`, `done`) are not used for Sprint 2. The label policy is in §2.2.

- Status is recorded in the artefact (Draft, In Review, Approved, Provisional), not in GitHub labels.
- Completion of an Issue is its Open or Closed state, with PR review and merge evidence.
- The type-label policy is in §2.2. All five `type:*` labels were verified to exist.

### 2.2 GitHub label policy (mentor-approved; labels verified)

The five approved `type:*` labels were verified read-only on 13 November 2026. The two older type-oriented labels below are listed for reference only and are not used for Sprint 2:

| Label | Description (from repository) |
|-------|-------------------------------|
| `documentation` | Improvements or additions to documentation |
| `enhancement` | New feature or request |

No label was created, renamed or deleted by this manifest.

| Work type | Issues | Proposed label | Basis | Status |
|-----------|--------|----------------|-------|--------|
| Feature specification (draft) | S2-01, S2-02 | `type:feature` | Approved type label, verified to exist. | Verified |
| Contract (draft, preliminary or completion) | S2-03, S2-05, S2-07 | `type:contract` | Approved type label, verified to exist. | Verified |
| ADR | S2-08 | `type:adr` | Approved type label, verified to exist. | Verified |
| Technical task (code and tests) | S2-06, ST-01 | `type:technical-task` | Approved type label, verified to exist. | Verified |
| Technical spike | S2-04 | `type:spike` | Approved type label, verified to exist. | Verified |

**Assignees.** No Issue has an assignee. Each Issue names its provisional pair in its body (Pair A, B, C or D). The mentor assigns students after rotation.

**Resolved gaps (G-1 to G-4):** the five labels now exist, so the proxies are no longer used.

- **G-1:** resolved. `type:feature` exists.
- **G-2:** resolved. `type:contract` and `type:adr` exist.
- **G-3:** resolved. `type:technical-task` exists.
- **G-4:** resolved. `type:spike` exists.
- Workflow documents and forms were updated to the type-only policy in this step. The conflicts are listed in §7 (G-5 to G-8), now resolved.

**Pairs.** Pair letters are provisional workstreams. The mentor assigns individual students. **No student is named in this manifest.** Pair composition and GitHub assignees are finalised during Sprint Planning. No Issue has an assignee.

### 2.1 Completion semantics (mentor-approved model)

Each issue states four things separately:

- **Sprint 2 deliverable:** what is merged or produced in Sprint 2.
- **Closing conditions:** the checks that let the Sprint 2 issue close.
- **Formal approval:** specification or contract approval, where it applies. It is never claimed by this manifest, and it is not a closing condition for a draft.
- **Sprint 3 follow-up:** work outside Sprint 2. It is recorded in §8 for planning only. No successor issue is created.

**Approved model for draft issues (S2-01, S2-02, S2-03, S2-07).**

- A draft issue may close when its explicitly defined draft deliverable meets the applicable review and merge requirements.
- Closing a draft issue does **not** mean the full specification or contract is formally Approved.
- Final approval, and the remaining effort, belong to Sprint 3 (§8).
- The SDD approval gates and the Definition of Done for final specifications and contracts are unchanged.
- Closing a draft-only issue records that the Sprint 2 draft deliverable is complete. It does not change the artefact's status or approval. Its closing comment states that the draft is complete and the document is not formally approved. SDD Workflow §6 now states this rule (PC-1, resolved in this step).

---

## 3. Issues

### S2-01: [SPEC-F06] Submit solution specification (draft)

**GitHub Issue:** [#29](https://github.com/ITA-Internship/GoJudge/issues/29)

**Backlog ID:** SPEC-F06 · **Form:** Feature request · **Type label:** `type:feature` · **Branch:** `docs/spec-f06` · **Status:** NEW · **Pair:** A (provisional) · **Pair-days:** 1.5

**Objective and scope.** Draft the feature specification for submitting a solution: asynchronous submission, the QUEUED state returned to the learner, and links to the lifecycle and messaging contracts.

**Sprint 2 draft deliverable.**
- `specs/features/F06-submit-solution.md` from `specs/templates/feature-spec-template.md`.
- Status `Draft`, moved to `In Review` when the author pair submits it.

**Testable acceptance criteria (draft).**
1. The file exists at `specs/features/F06-submit-solution.md` and uses the template's required sections.
2. Acceptance criteria are numbered AC-01 onward, each with a verification entry.
3. Every state and verdict matches CON-01 by name. The spec does not redefine them.
4. Message behaviour is cited from CON-02 and not restated. Values CON-02 does not fix are listed as open questions.
5. The duplicate-submission behaviour is recorded as an explicit open question unless CON-01 answers it (sprint-02.md §5.4).
6. The security section states that learner submissions are untrusted.

**Closing conditions (Sprint 2).** Draft merged by pull request, approved by a pair other than A. Status matches the file. Criteria 1–6 evidenced in the pull request.

**Formal approval (not a Sprint 2 closing condition).** Requires CON-02 agreement and a complete approval checklist. Recorded in §8 as a Sprint 3 follow-up.

**Dependencies.** CON-01 (merged, done). CON-02 (S2-05) for approval only; the draft may use its draft.

**Estimate.** 1.5 pair-days (backlog §4.1), covering the draft. Review is overhead outside the task limit.

---

### S2-02: [SPEC-F08] Submission history specification (draft)

**GitHub Issue:** [#30](https://github.com/ITA-Internship/GoJudge/issues/30)

**Backlog ID:** SPEC-F08 · **Form:** Feature request · **Type label:** `type:feature` · **Branch:** `docs/spec-f08` · **Status:** NEW · **Pair:** A (provisional) · **Pair-days:** 1.5

**Objective and scope.** Draft the feature specification for submission history: a learner sees their own submissions with timestamps, states and final verdicts.

**Sprint 2 draft deliverable.** `specs/features/F08-submission-history.md` from the template, status `Draft` or `In Review`.

**Testable acceptance criteria (draft).**
1. The file exists and uses the template's required sections.
2. Acceptance criteria are numbered AC-01 onward, each with a verification entry.
3. Lifecycle states and verdicts match CON-01 by name. Only the learner's own submissions are shown.
4. The verdict field is described as populated once F07-05 exists, consistent with backlog F08-01.
5. The security section covers access control for other learners' submissions.

**Closing conditions (Sprint 2).** Draft merged by pull request, approved by a pair other than A. Status matches the file. Criteria 1–5 evidenced.

**Formal approval (not a Sprint 2 closing condition).** Sprint 3 follow-up (§8). The backlog places SPEC-F08 approval in Sprint 3.

**Dependencies.** CON-01 (merged, done). No other blocker.

**Estimate.** 1.5 pair-days (backlog §4.1), draft only.

---

### S2-03: [CON-04] Problem packaging and test visibility contract (draft)

**GitHub Issue:** [#31](https://github.com/ITA-Internship/GoJudge/issues/31)

**Backlog ID:** CON-04 · **Form:** Technical task (work type "Contract") · **Type label:** `type:contract` · **Branch:** `docs/con-04-draft` · **Status:** NEW (draft) · **Pair:** A (provisional) · **Pair-days:** 2.0

**Objective and scope.** Draft the contract for problem packaging, the starter template format and the visibility rules for sample and hidden tests.

**Sprint 2 draft deliverable (approved split: 2.0).**
- `docs/contracts/CON-04.md`, marked **Draft**.
- Packaging format, starter template format and visibility rules, each stated or listed as an open question.
- Each DEC-09 option listed with its trade-offs. No option selected.
- Dependent pairs named: F03-01, F07-04 (judge), F09-02 (admin), SPEC-F03.

**Testable acceptance criteria (draft).**
1. Each of the three areas has a stated value or a listed open question. None is assumed.
2. Hidden tests are stated as never returned to a learner, in any response or log.
3. DEC-09 is marked Open. The draft does not treat it as decided.
4. Dependent pairs are named, and the draft lists which decisions block final approval.
5. The draft is merged to `docs/contracts/CON-04.md` with status Draft.

**Closing conditions (Sprint 2).** Criteria 1–5 evidenced. Draft merged by pull request, reviewed by a pair other than A.

**Formal approval (not a Sprint 2 closing condition).** Final contract and agreement from every dependent pair. Needs DEC-09. Sprint 3 remainder: **2.0 pair-days** (§8).

**Dependencies.** DEC-09 (MA-14), open. It blocks the final contract, not the draft.

**Estimate.** 2.0 pair-days (approved split of backlog CON-04, size 4).

---

### S2-04: [SPK-05] Isolation mechanism comparison

**GitHub Issue:** [#32](https://github.com/ITA-Internship/GoJudge/issues/32)

**Backlog ID:** SPK-05 · **Form:** Technical spike · **Type label:** `type:spike` · **Report branch:** `docs/spk-05-report` · **Experiment branch (not merged):** `spike/SPK-05-*` · **Status:** NEW · **Pair:** B (provisional) · **Pair-days:** 3.0 (timebox)

**Objective and scope.** Compare candidate isolation mechanisms for judge execution against the proposal §7 isolation requirements. The starting evidence is the SPK-01 No-go, including the network-egress race. **The spike produces evidence and a recommendation only.** It does not approve DEC-08, DEC-02 or production judge implementation. Adoption needs MA-03.

**Evidence classes.** Every matrix cell and finding carries exactly one class:

- **E, Experimental:** run in the isolated spike environment. The command, image version, result and log reference are recorded.
- **D, Documentation-based:** taken from cited upstream documentation or the SPK-01 report, with version and source. Not re-run here.
- **H, Hypothesis:** not tested and not documented. Listed as a risk, never as a finding.

No cell may be described as secure, approved or verified without an E or D source. "Not verified" is an acceptable and expected outcome.

**Sprint 2 deliverable.**
- Spike Brief on the issue before the timebox starts: hypothesis, timebox, experiment, success and failure criteria, required evidence.
- Report `docs/spikes/SPK-05.md` containing:
  - the requirement-by-candidate matrix (SC-A);
  - the DEC-08 recommendation, with rejected options and reasons (SC-B);
  - an isolation-gap register (SC-C, below), with limitations for each candidate.

**Testable acceptance criteria.**
1. **SC-A (backlog §5.5):** every proposal §7 requirement is mapped to every candidate, each cell with an evidence class and source. No cell is empty.
2. **SC-B (backlog §5.5):** a written DEC-08 recommendation, with rejected options and reasons. It states that nothing is adopted without MA-03.
3. **SC-C (spike-brief addition):** an isolation-gap register. For each candidate it states **network-egress timing**: whether egress is blocked from the first instruction of the container, or only after a delay. It cites the SPK-01 finding: outbound egress succeeded in the window before per-Pod NetworkPolicy rules were installed (F-03; first connection after 19 ms; window under about 1 s on the k3s embedded controller). Each candidate's egress control is recorded as Verified (E), Documented only (D) or Not verified (H), with its limitations. Each Not verified item is a named risk.

**Empirical standard for egress (mentor clarification).** A claim that egress is blocked from the first instruction of a container is treated as **verified only on E evidence**. A D-only claim is recorded as "Documented only", and it is not equivalent to empirical verification of this security property. Where no E evidence exists, the residual risk is recorded in the register: the window before the per-Pod rules, its measured or cited length, and which candidates it affects. The spike does not extend its timebox to close this gap without MA-12.
4. The SPK-01 init-container gate is listed as a candidate for research only, not as a fix. Its known limitations (SPK-01 §5) are carried over.
5. The report runs no learner-submitted code. Any experiment uses trusted packages only, in the isolated spike environment.
6. The report is reviewed by a pair other than B (sprint-02.md §9).
7. The report contains no DEC-08, DEC-02 or production judge approval and no production judge code.

**Closing conditions (Sprint 2).** Timebox respected. Criteria 1–7 evidenced. Outcome stated. Report merged to `docs/spikes/SPK-05.md` and reviewed by another pair. No spike code on `main`.

**Dependencies.** SPK-01 report (merged, outcome No-go). DEC-02 is not approved by the No-go. DEC-08 is Open.

**Estimate.** 3.0 pair-days, the timebox. An extension needs MA-12.

**Sprint 3 follow-up.** DEC-08 decision (start of Sprint 4). Any production judge work waits for it.

---

### S2-05: [CON-02] Messaging contract (completion) — EXISTING #23

**GitHub Issue:** [#23](https://github.com/ITA-Internship/GoJudge/issues/23), reused and updated. Title is now "[CON-02] Messaging contract (Sprint 2 completion)".

**Backlog ID:** CON-02 · **Existing issue:** **#23** "[CON-02] Messaging contract (stretch)", open · **Form:** Technical task (work type "Contract") · **Type label:** `type:contract` · **Branch:** `docs/con-02-messaging` · **Status:** **EXISTING, update required** · **Pair:** B (provisional; Pair C was the Sprint 1 author) · **Pair-days:** 1.5

**Objective and scope.** Complete the messaging contract for submission dispatch and processing: queue names, message schema, idempotency key, acknowledgement, retry and dead-letter behaviour. It is agreed by the dependent pairs and merged.

**Sprint 2 deliverable.**
- `docs/contracts/CON-02.md` completed from the Sprint 1 draft, with the SPK-03 risks listed.
- Pull request on `docs/con-02-messaging`, reviewed by a pair other than B.
- Dependent pairs named, and their agreement recorded on #23 or the pull request.

**Testable acceptance criteria.**
1. The contract states queue names, message schema, idempotency key, acknowledgement, retry behaviour and dead-letter behaviour. None is left to the implementation.
2. Every attempt count or limit is stated as a value or listed as an open question.
3. SPK-03 risks are listed. Criteria SPK-03 did not verify are shown as risks, not guarantees.
4. DEC-04 and DEC-12 are marked not approved unless the mentor records their approval (MA-26, MA-27).
5. Infrastructure failures and test failures are distinguished, consistent with CON-01.
6. Dependent pairs are named. Each records its agreement before merge.
7. #23 is updated to Pair B ownership and the Sprint 2 scope (update note below).

**Formal Done: explicit blockers.** CON-02 cannot reach formal Done until all of the following hold. None is assumed, and the criteria are not weakened:

| # | Blocker | Owner | Status |
|---|---------|-------|--------|
| B1 | **DEC-04 approved by the mentor** (MA-26). CON-02 is the contract that relies on it. | Mentor | **Open** |
| B2 | **Every dependent pair has recorded its agreement.** Dependent pairs are the owners of F06-02 and F07-01, which are not yet named. | Mentor (names, MA-28); each dependent pair | **Open** |
| B3 | **Review and merge.** Pull request approved by a pair other than B, then merged. | Pair B; reviewer pair | Not started |

**Does this issue close as a draft?** The approved draft-closure model covers S2-01, S2-02, S2-03 and S2-07 only. CON-02 is a contract completion, and CONTRIBUTING §11 applies to it in full. The draft-closure model does not apply to CON-02 unless the mentor extends it (PC-2). Until then, if B1–B3 are not met by 27 November, 15:00, #23 stays open with the missing blocker named. It is not reported as done.

DEC-12 does not block Done, provided the contract names dispatch as dependent on DEC-12 and states that a DEC-12 decision that changes dispatch needs renewed agreement.

**Update note for existing issue #23, for mentor approval, not applied.**
- Owner: Pair C (Sprint 1) to Pair B (Sprint 2). Pair C remains a dependent pair.
- Sprint: "Sprint 1, second week, stretch" to "Sprint 2 (16–27 November 2026), committed".
- Estimate: unchanged at 1.5.
- The Sprint 1 "not an exit criterion" text is superseded.

**Dependencies.** CON-01 (merged, done). SPK-03 (merged, done). B1 to B3 above. DEC-12 (MA-27) needed by 30 November for Sprint 3 dispatch.

**Estimate.** 1.5 pair-days (backlog §4.1).

**Sprint 3 follow-up.** If #23 is not Done by 27 November, its remaining work and blockers carry over. No estimate is invented for the residual.

---

### S2-06: [INF-07] Structured logging and request correlation (Sprint 2 subset)

**GitHub Issue:** [#33](https://github.com/ITA-Internship/GoJudge/issues/33)

**Backlog ID:** INF-07 (Sprint 2 subset) · **Form:** Technical task · **Type label:** `type:technical-task` · **Branch:** `infra/INF-07-structured-logging` · **Status:** NEW · **Pair:** C (provisional) · **Pair-days:** 4.0 (maximum timebox, approved)

**Objective and scope.** Build reusable backend logging and request correlation for future API handlers, as a Go library under `internal/`.

**Approved Sprint 2 scope (only these):**
- Structured JSON logging using the standard-library `log/slog`.
- HTTP request ID middleware.
- Propagation of a valid inbound request ID, or generation of a new one when absent or invalid.
- Request ID availability through `context.Context`.
- Unit tests for request ID generation, propagation and logging.
- Brief usage documentation.

**Excluded from this issue:** health endpoints, basic metrics, authentication, judge execution, an HTTP server or router choice, any error response format, and unrelated observability features.

**Why no separate ADR (mentor-approved).** `log/slog` and the standard `net/http` and `crypto/rand` packages introduce no new architectural decision and no new dependency. The middleware is a plain `http.Handler` wrapper and does not choose a router. No ADR is required, provided the implementation does not introduce a new architectural decision. If it does, the work stops and the mentor decides.

**Verified suitability.**
- Not implemented on `main` or any branch.
- Depends only on INF-01 and INF-02, both done. No dependency on Auth0, SPEC-F01, CON-02, CON-04, DEC-08 or judge code.
- Fits the README Repository Structure: private application code under `internal/`.

**Sprint 2 deliverable.**
- Package `internal/observability` (proposed name):
  - `NewLogger(w io.Writer, level slog.Level) *slog.Logger`: JSON handler.
  - `RequestID` middleware: keeps a valid inbound `X-Request-ID`, otherwise generates one. Sets the response header and stores the ID in the request context.
  - `RequestIDFromContext(ctx) string`: returns the ID, or `""` when none is set.
  - `AccessLog` middleware: one line per request with `method`, `path` (no query string), `status`, `duration_ms` and `request_id`.
- Unit tests for each function, using `httptest`, run under `go test -race`.
- Brief usage documentation `docs/architecture/logging.md`: allowed log fields, request ID rule, and how a handler uses the logger and context. Documentation, not an ADR.

**Security requirements (mandatory).**
- No credentials, `Authorization` headers, cookies, tokens, passwords or secrets in logs. Headers are not logged.
- No learner-submitted source code. Request bodies are not logged.
- No query string in logs, because it can carry tokens.
- The middleware writes no error body, so no internal detail reaches responses.
- No personal data beyond the request ID, method, path, status and duration.

**Testable acceptance criteria.**
1. `NewLogger` writes one JSON object per line with `time`, `level` and `msg` plus structured fields. A test decodes each line.
2. A valid inbound `X-Request-ID` (1–64 characters, `[A-Za-z0-9_-]`) is kept unchanged. It appears in the response header, in `RequestIDFromContext` and in the access-log line. A test checks all three.
3. A missing, empty, over-long (more than 64) or invalid inbound ID is replaced by a generated 32-character lowercase hex ID from `crypto/rand`. A test checks each invalid case, and that two generated IDs differ.
4. `RequestIDFromContext` returns `""` when no middleware ran. A test checks this.
5. The access-log line contains `method`, `path`, `status`, `duration_ms` and `request_id`, and no query string. A test requests `/x?token=SENTINEL-1` and asserts that `SENTINEL-1` is absent from the log.
6. The access log never contains the `Authorization` header, a cookie or the request body. A test sends `Authorization: Bearer SENTINEL-2` and a body containing `SENTINEL-3`, and asserts neither appears in the log.
7. The middleware writes no body of its own. A test checks that a handler's body passes through unchanged, and that the middleware adds nothing on a 500.
8. `go.mod` gains no `require` line, and `go mod tidy` leaves no diff.
9. CI passes: `go build ./...`, `golangci-lint` (including `gosec`, `errcheck`, `gofmt`, `goimports`), and `go test -race ./...`.
10. `docs/architecture/logging.md` exists and lists the allowed log fields and the request ID rule.
11. No token, key or secret appears in any file in the pull request.

**Closing conditions (Sprint 2).** Criteria 1–11 evidenced in the pull request. CI green. Approved by a pair other than C. Merged to `main`. Closed as Done when the conditions above are met. No status label is applied.

**Dependencies.** INF-01 (done), INF-02 (done). No external blocker.

**Estimate.** 4.0 pair-days, the maximum timebox. If the work finishes early, the unused time is not reallocated without the mentor's approval.

**Sprint 3 follow-up (not an issue; sized by the mentor).** Health endpoints, basic metrics, and wiring the middleware into an HTTP server. These are the rest of backlog INF-07 and are unsized (§8).

---

### S2-07: [CON-05] Data model: preliminary contract (users, roles, submissions, results)

**GitHub Issue:** [#34](https://github.com/ITA-Internship/GoJudge/issues/34)

**Backlog ID:** CON-05 (preliminary) · **Form:** Technical task (work type "Contract") · **Type label:** `type:contract` · **Branch:** `docs/con-05-data-model` · **Status:** NEW (preliminary; **not an approved final contract**) · **Pair:** D (provisional) · **Pair-days:** 3.0 (approved split of backlog CON-05, size 4)

**Objective and scope.** Draft the preliminary data model for entities that do not depend on CON-04: users, roles, submissions and results. Problem and test-asset entities are excluded until CON-04 is final.

**Sprint 2 preliminary deliverable (approved).**
- `docs/contracts/CON-05.md`, status **Preliminary**, listing the entities for users, roles, submissions and results.
- For each entity: fields, keys and relations, consistent with CON-01 states and verdicts, and with SPEC-F01's `users` model (OQ-9 open).
- Open questions listed. None is resolved by assumption.

**Testable acceptance criteria (preliminary).**
1. The document is marked **Preliminary** and states that it is not the final approved data contract.
2. It covers users, roles, submissions and results. Problem and test-asset entities are listed as "outside this draft, waiting for CON-04 final".
3. Submission states and verdicts match CON-01 by name.
4. Every field type or constraint is stated or listed as an open question.
5. OQ-9 remains open unless the mentor records a decision. `users` fields that depend on OQ-1 (role source) are marked open.
6. Dependent pairs are named, and the draft lists which parts block final approval.

**Closing conditions (Sprint 2).** Criteria 1–6 evidenced. Preliminary draft merged by pull request, reviewed by a pair other than D.

**Formal approval (not a Sprint 2 closing condition).** Needs CON-04 final and dependent-pair agreement. Sprint 3 remainder: **1.0 pair-day** (§8).

**Dependencies.** CON-01 (merged, done). CON-04 final (Sprint 3) for problem and test-asset entities and final approval. SPEC-F01 OQ-1 and OQ-2 (open) affect the users and roles entities. INF-05 is not in Sprint 2 and depends on DEC-11 and this contract.

---

### S2-08: [DEC-11] Database migration tool: ADR-003 (Provisional)

**GitHub Issue:** [#35](https://github.com/ITA-Internship/GoJudge/issues/35)

**Backlog ID:** DEC-11 · **Form:** Technical task (ADR) · **Type label:** `type:adr` · **Branch:** `docs/adr-003-migration-tool` · **Status:** NEW · **Pair:** D (provisional) · **Pair-days:** 1.0 (provisional estimate, mentor-accepted)

**Objective and scope.** Record the database migration tool decision as an ADR, following the ADR process (INF-04, [ADR README](../architecture/adr/README.md), [template](../architecture/adr/template.md)). The ADR records options and a recommendation.

**Sprint 2 deliverable.**
- `docs/architecture/adr/ADR-003.md` from the template, Status **Provisional**.
- At least two options with trade-offs, including compatibility with the pgx and sqlc approach named in INF-05.
- Recommendation, and the approval condition: "NOT approved until confirmed by MA-23".

**Testable acceptance criteria.**
1. `ADR-003.md` exists, uses the ADR template, and has Status Provisional.
2. Two or more options are listed with trade-offs. No option is presented as approved.
3. The approval condition names MA-23.
4. No migration is implemented and INF-05 is not started.
5. ADR-001 and ADR-002 are not modified.

**Closing conditions (Sprint 2).** Criteria 1–5 evidenced. ADR merged as Provisional, reviewed by a pair other than D.

**Formal approval (not a Sprint 2 closing condition).** Mentor approval of the migration tool decision (MA-23). The merge does not imply approval.

**Dependencies.** MA-23 for approval only. The repository has no migration files, so there is no existing convention to follow.

**Estimate.** 1.0 pair-day, the provisional estimate the mentor has accepted. It is not in the backlog and is recorded as such.

---

## 4. Stretch Candidate (not committed)

### ST-01: [INF-09] Auth0 development tenant configuration (conditional Stretch)

**Backlog ID:** INF-09 · **Form:** Technical task · **Type label (if authorised):** `type:technical-task`. Not created as an Issue (conditional Stretch, §4). · **Branch:** `infra/INF-09-auth0-tenant` · **Status:** **CONDITIONAL STRETCH, not committed** · **Pair:** C (provisional) · **Pair-days:** 1.5 (not in the 17.5)

**Objective and scope.** Configure the Auth0 development tenant and application for the Go API, with no secrets in the repository. This is configuration, not authentication code.

**May start only when all of these hold:**
1. **MA-05 confirms** Auth0 tenant ownership and admin access. Not confirmed in the repository.
2. **The mentor explicitly authorises** this stretch task.
3. **Pair C has sufficient remaining capacity.** Pair C is at 4.0 for S2-06. Its working figure is 5.0, so 1.0 remains, which is **less than 1.5**. The condition is not met unless the mentor decides how to resolve it (PC-3).

**Deliverables (if started).**
- Development tenant and application configured, documented in the README or `infra/` setup notes.
- A record of tenant ownership and admin access, by role name only.
- Secret values supplied by environment, not committed.

**Acceptance criteria (if started).**
1. The tenant and application exist. The configuration is reproducible from the documented steps.
2. No secret, token, client secret or tenant credential is in the repository or in CI logs (checked in review).
3. The role-claim mapping is **not** implemented here. OQ-1 is open in SPEC-F01.
4. The configuration matches DEC-03 as Provisional. The issue does not claim DEC-03 approval.

**Closing conditions (if started).** Criteria 1–4 evidenced. Approved by a pair other than C. Merged.

**Estimate.** 1.5 pair-days, from backlog §4.2. Not counted in the committed total.

**If not started in Sprint 2.** It stays in the backlog, unchanged, for Sprint 3 planning (§8).

---

## 5. Estimate Checks and Inconsistencies

| # | Item | Finding | Handling |
|---|------|---------|----------|
| 1 | S2-06 INF-07 (4.0) | The backlog INF-07 full size (4) also covers health endpoints and basic metrics, which are excluded. | 4.0 is the approved maximum timebox for the subset. The remainder is unsized (§8). Not re-estimated. |
| 2 | Backlog INF-07 name | The backlog name is "Structured logging, health endpoints and basic metrics". The Sprint 2 title is narrower. | Title marked "(Sprint 2 subset)". Backlog keeps the full name. |
| 3 | S2-08 DEC-11 (1.0) | Not sized in the backlog. | Provisional estimate, accepted by the mentor. |
| 4 | S2-05 CON-02 (1.5) | Some Sprint 1 stretch drafting may exist, but no draft is on `main`. | 1.5 preserved as the approved estimate. Confirm whether Sprint 1 drafting exists. |
| 5 | S2-03 CON-04 (2.0) and S2-07 CON-05 (3.0) | Approved splits of backlog sizes 4 and 4. Remainders 2.0 and 1.0 are Sprint 3. | Recorded per issue and in §8. |
| 6 | S2-01 and S2-02 (1.5 each) | Backlog size covers the specification. Formal approval is not sized separately. | Approval effort is overhead outside the task limit. No new estimate is invented. |
| 7 | ST-01 INF-09 (1.5) | Pair C would be at 5.5 if authorised. The working cap is 5.0. | PC-3. Not resolved here. |
| 8 | S2-04 SPK-05 (3.0) | Timebox equals backlog size. | Consistent. |
| 9 | Backlog §10 Sprint 2 row (38.5) | Historical. Not the Sprint 2 planning figure. | Preserved as historical context in the backlog (§10.2). |

---

## 6. Existing Issues and Duplicate Check

| Existing issue | Title | State | Relation to Sprint 2 | Action |
|----------------|-------|-------|----------------------|--------|
| #23 | [CON-02] Messaging contract (stretch) | Open | S2-05 | **Reuse. Update owner, sprint and scope. Do not duplicate.** |
| #10 | [SPEC-F01] Accounts and access specification | Closed | Not Sprint 2 scope. Gates F01-01 later. Spec is In Review. | Mentor decides on reopening (§7, F-2). Do not create a second SPEC-F01 issue. |
| #24 | [SPEC-F04] Browser code editor specification (stretch) | Closed | Not Sprint 2 scope. Spec is In Review. | Mentor decides on reopening (§7, F-2). |
| #9 | [CON-01] Submission lifecycle | Closed | Dependency of S2-01, S2-02, S2-05, S2-07 | None. |
| #21 | "[INF-15] SonarQube…" (retitled by the mentor) | Closed | Identifier INF-15. Not Sprint 2 work. | **Resolved.** Not modified in this step. |

**Created Sprint 2 Issues:** #29 to #35 for S2-01 to S2-08, excluding S2-05. #23 is updated for S2-05. Links are in §1.1.

**No duplicate risk for the other committed items.** None of S2-01 to S2-08 or ST-01 has an existing issue, branch or PR.

**Not Sprint 2 issues (outside the approved scope):** F01-01 (implementation and preparation), SPEC-F02, SPEC-F03, SPEC-F07, CON-03, INF-05, INF-08, INF-12a, F02-01, F07-*, and all Sprint 3+ items.

---

## 7. Process Clarifications and Follow-Up (no action taken)

| ID | Item | Owner |
|----|------|-------|
| PC-1 | **Resolved in this step.** SDD Workflow §6 now states the closure rule for draft-only Issues. The label policy is type-only (§2.2). | — |
| G-5 to G-8 | **Resolved in this step.** The conflicting rules now follow the type-only policy: `docs/sdd-workflow.md` §6, `CONTRIBUTING.md` §2, §4 and §11, and `.github/ISSUE_TEMPLATE/feature.yml`. The quotes below are the superseded text. | — |

**G-5 (SDD Workflow §6, [docs/sdd-workflow.md](../sdd-workflow.md#6-github-issue-states)).** "The team tracks workflow states with the six agreed labels. No other labels are used for state. Exactly one state label is on an issue at a time." The mapping table on the same page lists `spec:draft`, `spec:in-review`, `spec:approved`, `in-progress`, `in-review` and `done` as state labels.

**G-6 (CONTRIBUTING §4, [CONTRIBUTING.md](../../CONTRIBUTING.md#4-github-issue-workflow)).** Step 1: "Only the feature form applies `spec:draft`." Step 4: "Keep exactly one state label on an issue at a time." Step 5: "Set `done` and close the issue once the Definition of Done (§11) for its type is met." The closing line of the same section lists the six labels as the only labels.

**G-7 (issue form, `.github/ISSUE_TEMPLATE/feature.yml`, line 4).** `labels: ["spec:draft"]`. Applies a status label automatically when a feature issue is created.

**G-8 (CONTRIBUTING §11).** "Then it is labelled `done` and its issue is closed." Requires a status label at completion.
| PC-2 | Does the draft-closure model extend to CON-02's completion, or does CON-02 keep full contract DoD (B1–B3)? This manifest assumes full DoD. | Mentor |
| PC-3 | **Pair C capacity.** INF-09 would take Pair C to 5.5 against a 5.0 working figure. Options: the mentor accepts 5.5 for Pair C; INF-09 is assigned to another pair; or INF-09 waits for Sprint 3. Not decided here. | Mentor |
| F-1 | **Resolved.** Issue #21 is retitled "[INF-15] SonarQube…". | — |
| F-2 | Issues #10 (SPEC-F01) and #24 (SPEC-F04) are closed, but their specifications are In Review. Decide on reopening, or record the closed-but-In-Review state under PC-1. | Mentor, Pairs D and A |
| F-3 | **Resolved.** Issues #23 and #29 to #35 exist for the eight committed items (§1.1). Pair composition and assignees remain open (F-4). | — |
| F-4 | Pair composition. The mentor assigns students personally after rotation. No student is named here. | Mentor |

---

## 8. Validation

| Check | Result |
|-------|--------|
| Committed issues are exactly the approved items (S2-01 to S2-08) | Yes. INF-09 is listed only as ST-01, conditional Stretch. |
| Individual committed estimates sum to 17.5 | 1.5 + 1.5 + 2.0 + 3.0 + 1.5 + 4.0 + 3.0 + 1.0 = **17.5** |
| Pair allocations match the Sprint 2 plan | A 5.0, B 4.5, C 4.0, D 4.0 = 17.5 |
| INF-09 is not counted in the committed total | Yes. It is ST-01, 1.5, conditional. |
| No F01-01 work remains in Sprint 2 | Yes. F01-01 is excluded. Its former preparation task is replaced by S2-06 (INF-07). |
| Draft completion is distinguished from formal approval | Yes. Each draft issue states "Formal approval (not a Sprint 2 closing condition)". |
| Issue #23 is reused, not duplicated | Yes. S2-05 is EXISTING #23. |
| Type labels | Yes. Each committed Issue has one `type:*` label, and all five were verified to exist. No status label is prescribed. |
| No GitHub assignee on any Issue | Yes. Provisional pair named in each body. |
| INF-09 not created as an Issue | Yes. ST-01 is conditional and not scheduled. |
| Uncommitted capacity | 20 − 17.5 = **2.5** |
| Conditional total if ST-01 is authorised | 17.5 + 1.5 = 19.0, within the 20 maximum. Pair C at 5.5 (PC-3). |

### 8.1 Sprint 3 follow-up (planning only; no successor issues created)

| Item | Sprint 3 work | Estimate | Gate |
|------|---------------|----------|------|
| CON-04 final | Completion and formal approval of the contract | **2.0** (approved remainder) | DEC-09 decided (MA-14); dependent-pair agreement |
| CON-05 final | Completion and formal approval of the data contract | **1.0** (approved remainder) | CON-04 final; SPEC-F01 OQ-1 and OQ-2 for the users and roles entities |
| SPEC-F06 approval | Formal specification approval | Not separately sized. Review overhead (sprint-02.md §3.3). | CON-02 agreement; approval checklist; no critical open question |
| SPEC-F08 approval | Formal specification approval | Not separately sized. Review overhead. | Approval checklist; no critical open question |
| CON-02 (#23) | Remaining B1–B3 if not Done by 27 November | **Not sized.** Mentor to size any residual. | DEC-04 approved; dependent-pair agreement |
| INF-07 remainder | Health endpoints, basic metrics, HTTP server wiring | **Not sized.** The backlog total is 4 and the Sprint 2 subset takes 4. Mentor to size. | Mentor sizing (MA-33 follow-up) |
| INF-09 (ST-01) | If not started in Sprint 2: the stretch, unchanged | 1.5 (backlog) | MA-05; mentor authorisation; Pair C capacity (PC-3) |
| DEC-11 / ADR-003 | Mentor approval of the migration tool decision | Not sized. No task effort. | MA-23 |
| F01-01 | Preparation and implementation, after SPEC-F01 approval and INF-09 | Backlog 4 (not reduced) | SPEC-F01 Approved; INF-09 done |
| DEC-12 | Dispatch decision | Not a task | Needed by 30 November |
| DEC-08 | Isolation decision | Not a task | Start of Sprint 4 |

---

## 9. Mentor Follow-Up Before Issue Creation

1. Approve the committed list (S2-01 to S2-08) and the stretch (ST-01), and record PC-1 to PC-3.
2. Confirm pair composition (F-4) before any reviewer or owner is named on an issue. Assignees stay empty.
3. Done: the eight committed Issues exist (§1.1). ST-01 is created only if authorised. Pair composition and GitHub assignees are finalised during Sprint Planning.
4. The issues have **no** GitHub milestone or project assignment.
