# Sprint 2 Plan: 16–27 November 2026

**Status:** Finalised for mentor review (13 November 2026). GitHub Issues for the eight committed items were created the same day (links in [sprint-02-issues.md §1.1](sprint-02-issues.md)). **Mentor-approved:** sprint dates, the allocation in §3.2 (17.5 committed pair-days), the 20 pair-day maximum, the 2.5 uncommitted pair-days, INF-07 as committed backend work, INF-09 as conditional Stretch, the exclusion of F01-01 from Sprint 2, the partial-deliverable model, the SonarQube identifier INF-15, and the rule that the mentor assigns students to pairs. **Not approved:** any specification or contract, including SPEC-F01, SPEC-F04, SPEC-F06, SPEC-F08, CON-02, CON-04 and CON-05 final. Open technical decisions are listed in §11.4. No GitHub Issues have been created or changed.

**Sources:** [Product Backlog](../product/product-backlog.md) (revision 3; §10.2 Sprint 2 summary), [Sprint 1 Plan](sprint-01.md), [Sprint 2 Issue Manifest](sprint-02-issues.md), [SDD Workflow](../sdd-workflow.md), [Contributing Guide](../../CONTRIBUTING.md), spikes [SPK-01](../spikes/SPK-01.md), [SPK-02](../spikes/SPK-02.md), [SPK-03](../spikes/SPK-03.md), [CON-01](../contracts/CON-01.md), [ADR-001](../architecture/adr/ADR-001.md), [ADR-002](../architecture/adr/ADR-002.md), feature specs [F01](../../specs/features/F01-accounts-and-access.md) and [F04](../../specs/features/F04-browser-code-editor.md). GitHub state was read on 13 November 2026.

---

## 1. Sprint Goal

> Establish the core backend and data contracts while resolving the secure execution architecture blocker, with reduced workload to accommodate the university examination period.

Sprint 2 commits 17.5 pair-days of task work, out of a maximum of 20. Several items are limited by decisions and approvals that are not yet in place, so each one has a stated closing condition and a named blocker (§5). The plan commits only work that can start without an external approval gate.

---

## 2. Sprint 1 Carry-Over Status (verified 13 November 2026)

The mentor has confirmed that completed Sprint 1 work is Done. This table records what the repository shows. Merged PR approvals and closed issues are not specification approvals.

| Item | Pair | Evidence on `main` / GitHub | Status for Sprint 2 |
|------|------|-----------------------------|---------------------|
| INF-01 Go module, layout, lint | D | PR #12 merged; issue #6 closed | Done |
| INF-02 CI baseline | D | PR #14 merged; issue #7 closed | Done. Green on `main` not checked in this update. |
| INF-03 Local infrastructure | A | PRs #15, #16, #18 merged; issue #2 closed | Done. M1–M3 not re-verified in this update. |
| INF-04 ADR process | D | PR #11 merged; ADR-001 and ADR-002 Provisional; issue #8 closed | Done |
| SPK-01 Secure execution spike | B | PR #25 merged 13 Nov; outcome **No-go** (F-03: egress succeeded before per-Pod rules); issue #3 closed | Done. No-go stands (§5.1). |
| SPK-02 Auth0 OIDC check | C | PR #13 merged; outcome **Go**; issue #4 closed | Done |
| SPK-03 RabbitMQ and Watermill | C | PR #19 merged 13 Nov; outcome **Go** with SC-E added (transactional outbox); issue #5 closed | Done. Input for DEC-04 and DEC-12. |
| CON-01 Submission lifecycle | D | PR #17 merged 10 Nov; `docs/contracts/CON-01.md`; issue #9 closed | Done. Dependent-pair agreement is not recorded as pair-labelled agreement (§8.4 item 3). Duplicate-submission question open (§5.4). |
| SPEC-F01 Accounts and access | D | PR #20 merged 12 Nov; spec status **In Review**; reviewer outcomes in its §12; issue #10 closed | **Not Approved.** Critical OQ-1 and OQ-2 open (§5.2). |
| SPEC-F04 Browser editor (stretch) | D | PR #27 merged 13 Nov; spec status **In Review**; issue #24 closed | **Not Approved.** Not in Sprint 2 scope. |
| CON-02 Messaging contract (stretch) | C | No `docs/contracts/CON-02.md` on `main`; issue #23 **open**; no PR | **Not Done.** Carried into Sprint 2 as S2-05 on #23 (§5.3). |
| SonarQube integration (issue #21) | — | PR #26 merged 13 Nov; issue closed | Done. Identifier **INF-15** confirmed. Issue #21 is now titled "[INF-15] SonarQube…" (resolved). |

---

## 3. Capacity and Committed Effort

### 3.1 Capacity

| Item | Pair-days |
|------|-----------|
| Maximum committed task workload across all four pairs (mentor, approved) | **20** |
| Approved committed task workload (§3.2) | **17.5** |
| Uncommitted capacity (approved) | **2.5** |
| Excluded from the limit: meetings, reviews, planning, demos, handovers and other organisational activities | Accounted for separately (§3.3) |

### 3.2 Committed effort (approved allocation)

| Pair | Committed items | Pair-days |
|------|-----------------|-----------|
| A | SPEC-F06 draft (1.5), SPEC-F08 draft (1.5), CON-04 draft (2.0) | 5.0 |
| B | SPK-05 (3.0), CON-02 completion (1.5) | 4.5 |
| C | INF-07 structured logging and request correlation, Sprint 2 subset (4.0) | 4.0 |
| D | CON-05 preliminary contract (3.0), DEC-11 ADR (1.0) | 4.0 |
| **Total committed** | | **17.5** |
| Uncommitted | | 2.5 |

**Conditional Stretch, not committed:** INF-09 Auth0 development tenant configuration (1.5, Pair C). It is conditional on MA-05, explicit mentor authorisation, and Pair C capacity. It is not in the 17.5 (§4.3, §8.2).

### 3.3 Overhead outside the committed cap

The mentor has confirmed that the 20 pair-day limit covers committed task workload only. Meetings, reviews, planning, demos, handovers and other organisational activities are excluded and accounted for separately. The backlog's Sprint 2 overhead estimate is 15 pair-days (backlog §11.3). It is shown for information only and is not part of the 17.5.

---

## 4. Dependency Validation

### 4.1 Item-by-item check

| Item | Prerequisites | State on 13 Nov | Can start in Sprint 2? | Notes |
|------|---------------|-----------------|------------------------|-------|
| **SPEC-F06 (draft)** | CON-01 (done); CON-02 for approval | Not on `main` | **Yes, draft.** Approval: no. | Draft closes on its Sprint 2 scope (§8.3). Formal approval is Sprint 3. |
| **SPEC-F08 (draft)** | CON-01 (done) | Not on `main` | **Yes, draft.** Approval: no. | Same closure model. Backlog places approval in Sprint 3. |
| **CON-04 (draft)** | DEC-09 (open) | Not on `main` | **Yes, draft**, with DEC-09 recorded as open. Final: **no**. | Sprint 3 remainder 2.0. |
| **SPK-05** | SPK-01 report (merged, No-go) | Report merged | **Yes**, as an evidence spike. | Does not approve DEC-08, DEC-02 or production judge work (§5.1). |
| **CON-02 (completion, #23)** | CON-01 (done), SPK-03 (done); DEC-04 approval (MA-26); dependent-pair agreement | Draft not on `main`. #23 open. | Completion: **yes**. Formal Done: **no** until blockers B1–B3 (§8.4). | Explicit blockers, not weakened criteria. |
| **INF-07 subset (S2-06)** | INF-01 (done), INF-02 (done) | Not implemented. No branch or PR. | **Yes.** No external gate. | Standard library only. No ADR required (mentor-approved). |
| **CON-05 (preliminary)** | CON-01 (done) | Not started | **Yes, preliminary.** Final: **no**, needs CON-04 final. | Sprint 3 remainder 1.0. |
| **DEC-11 (ADR-003)** | Engineering options | Open | **Yes, Provisional ADR.** Approval: MA-23. | Estimate 1.0 accepted provisionally by the mentor. |
| **INF-09 (stretch)** | MA-05; mentor authorisation; Pair C capacity | No record of tenant ownership or admin access | **No, unless all three hold.** | Not in the 17.5. Pair C has 1.0 left against 1.5 (§4.3). |
| **F01-01** | SPEC-F01 Approved; INF-09 | Neither met | **No.** Not in Sprint 2 scope. | Preparation task removed (§5.5). |
| **SPEC-F01 approval** | OQ-1 and OQ-2 closed; checklist; reviewer-pair approval | In Review | **No** for Sprint 2 scope. | Gates F01-01 later. |
| **DEC-09** | Mentor input (MA-14) | Open | Needs decision before CON-04 final | Blocks CON-04 final and SPEC-F03 approval. |
| **DEC-12** | Engineering and mentor decision | Open. SPK-03 recommends the transactional outbox. | Decide by 30 November | Needed for Sprint 3 dispatch. |
| **DEC-08** | SPK-05 evidence | Open | **No** in Sprint 2 | Start of Sprint 4. |

### 4.2 Blockers and closure conditions

1. **CON-02 (S2-05).** Its Sprint 2 completion is committed. Formal Done needs DEC-04 approval (MA-26), agreement from every dependent pair (owners of F06-02 and F07-01 not yet named, MA-28), and review and merge (§8.4). If these are not met by 27 November, #23 stays open with the blocker named.
2. **CON-04 (S2-03).** The draft is committed. The final contract, and its remaining 2.0 pair-days, are Sprint 3 work after DEC-09.
3. **CON-05 (S2-07).** The preliminary contract is committed. Its problem and test-asset entities wait for CON-04 final. Final approval and the remaining 1.0 pair-day are Sprint 3 work.
4. **SPEC-F06 and SPEC-F08 (S2-01, S2-02).** Drafts are committed. Formal approval is Sprint 3 work. SPEC-F06 approval also needs CON-02 agreement.

### 4.3 Stretch and capacity constraint (INF-09)

| Pair | Approved (working figure 5.0) | Remaining | If INF-09 (1.5) is authorised |
|------|-------------------------------|-----------|-------------------------------|
| A | 5.0 | 0.0 | Not affected |
| B | 4.5 | 0.5 | Not affected |
| C | 4.0 | **1.0** | **5.5. Above the working figure of 5.0** |
| D | 4.0 | 1.0 | Not affected |
| **Total** | **17.5** | **2.5** | **19.0** (within the 20 maximum) |

INF-09 may start only if MA-05 is confirmed, the mentor explicitly authorises it, and Pair C has sufficient capacity. At 1.0 remaining, Pair C does not have sufficient capacity for 1.5. This is decision PC-3 (§11.3). The plan does not resolve it by changing the approved allocation.

---

## 5. Blockers and Risks Requiring Action

### 5.1 SPK-01 No-go (secure execution)

SPK-01 reported **No-go**. F-03 was met: an outbound connection from the Job succeeded before per-Pod network rules were installed. The consequences are binding for Sprint 2:

- DEC-02 is not approved.
- No production judge implementation starts in Sprint 2 (F07-02a, F07-02b, F05-01, F07-03 and the rest of the judge chain).
- DEC-08 is not decided in Sprint 2. Its decision is needed at the start of Sprint 4.
- The init-container gate is a candidate for research only, not a fix (SPK-01 §5).

SPK-05 produces evidence and a recommendation only. It classifies each finding as experimental, documentation-based or hypothesis, and records the network-egress timing for each candidate, including the SPK-01 race. A claim that egress is blocked from the first instruction of a container counts as verified only on experimental evidence. Documentation alone is recorded as such, with its residual risk. It approves nothing.

### 5.2 SPEC-F01 is not Approved

SPEC-F01 is In Review. Its §12 records reviewer outcomes from Pair A (10 November) and Pair C (11 November), but its approval checklist is not ticked, and critical OQ-1 (role source) and OQ-2 (DEC-03 approval) are open. A merged pull request and a closed issue do not make it Approved. Approval gates F01-01, which is outside Sprint 2.

### 5.3 CON-02 is not Done

CON-02 is the Sprint 2 completion item on #23. Its formal Done has three blockers: DEC-04 approval, dependent-pair agreement, and review and merge. The acceptance criteria are not weakened to close it.

### 5.4 Open question on duplicate submissions (CON-01)

CON-01 has an open question on duplicate submissions, owned by Pair D, due at the start of Sprint 2. The default in CON-01 is a conflict error. The SPEC-F06 draft records it as an open question unless it is answered first.

### 5.5 INF-07 replaces F01-01 preparation; INF-09 is Stretch

- **F01-01 is not in Sprint 2.** Its preparation task (3.5) is removed. It did not map to a defined subset of F01-01's backlog size, and its deliverables were design notes rather than backend work.
- **INF-07 (Sprint 2 subset, 4.0)** replaces it. Verified: not implemented; no dependency on Auth0, SPEC-F01, CON-02, CON-04 or judge work; standard library only. Scope is limited to structured logging, request-ID middleware, propagation and generation, context availability, unit tests and brief documentation. Health endpoints and metrics are excluded.
- **INF-09 (1.5)** is conditional Stretch (§4.3).

### 5.6 Other risks

| ID | Risk | Mitigation |
|----|------|------------|
| R-S2-1 | Reviews consume more time than planned | Review assignments at sprint start. Reviews are reported separately (§3.3). |
| R-S2-2 | Examinations reduce attendance | Staffing confirmation (MA-11) is not verified in the repository. Absences are recorded at the mid-sprint check. |
| R-S2-3 | Sprint 1 reports merged late (SPK-03 on 13 Nov, CON-01 on 10 Nov) leave less time for Sprint 2 decisions | Decisions are scheduled in §11.1. |
| R-S2-4 | Issues #10 and #24 are closed while their specifications are In Review | Reconcile before Sprint 2 starts (§15, F-2). |

---

## 6. Identifier Checks

### 6.1 SonarQube identifier

The mentor has confirmed **INF-15** for the SonarQube Cloud integration. The identifiers are:

- **INF-05:** PostgreSQL schema, migrations and data access (backlog §4.2).
- **INF-06:** OpenAPI validation in CI (backlog §4.2).
- **INF-15:** SonarQube Cloud integration.

GitHub issue #21 was originally titled "[INF-05] SonarQube Cloud Integration and Mandatory Quality Gate". Its title is now "[INF-15] SonarQube…". This correction is resolved, and this plan did not change the issue.

### 6.2 PostgreSQL migrations

- No migration files exist on `main`: no `migrations/` directory and no `*.sql` files.
- The `spk03_submissions` and outbox tables in SPK-03 SC-E are created by spike code on a spike branch. They are not migrations and must not be reused as the schema.
- DEC-11 is the input for INF-05. Its record is ADR-003, the next free number after ADR-001 and ADR-002. No ADR-003 file exists yet.

---

## 7. Pair Allocation, Rotation and Handover

### 7.1 Workstreams (provisional until the mentor confirms pair composition)

| Pair | Workstream (provisional) | Committed | Pair-days |
|------|--------------------------|-----------|-----------|
| A | Specifications and problem data | SPEC-F06 draft, SPEC-F08 draft, CON-04 draft | 5.0 |
| B | Secure execution and consumer contract | SPK-05, CON-02 completion | 4.5 |
| C | Backend foundation | INF-07 Sprint 2 subset (stretch: INF-09) | 4.0 |
| D | Data architecture | CON-05 preliminary, DEC-11 ADR | 4.0 |
| | | **Total** | **17.5** |

The workstream letters describe work, not membership.

### 7.2 Rotation

The rotation principle is confirmed:

- Students who worked mainly on specifications and contracts in Sprint 1 move towards backend development.
- Students who worked on backend, infrastructure or technical spikes gain specification and architecture experience.

**The mentor assigns students to pairs personally.** This plan does not assign, suggest or infer any individual's pair. Rotation is not a blocker to this plan. Workstreams A–D are provisional, and the review assignments in §9 follow the confirmed composition.

### 7.3 Handover requirement

Before work changes hands, the outgoing and incoming owners exchange a short handover note on the issue or pull request. It states:

- the current status, the artefact location and the branch;
- open questions, with their owners;
- dependencies, with their dates;
- what was not verified.

The incoming owner confirms ownership on the issue (SDD Workflow §3). The handover applies to these Sprint 2 items: SPEC-F01 (reviews and the OQ-1 and OQ-2 owners), CON-02 (Sprint 1 stretch author Pair C to completion owner Pair B), the CON-01 duplicate-submission question, and SPK-05 or CON-05 if their owner changes.

---

## 8. Committed Work, Partial Deliverables, Stretch and Sprint 3 Follow-Up

### 8.1 Committed (approved)

| Item | Pair | Pair-days | Sprint 2 deliverable | Closing condition for the Sprint 2 issue |
|------|------|-----------|----------------------|------------------------------------------|
| SPEC-F06 draft | A | 1.5 | Draft spec at `specs/features/F06-submit-solution.md` | Draft merged, reviewed by a pair other than A |
| SPEC-F08 draft | A | 1.5 | Draft spec at `specs/features/F08-submission-history.md` | Draft merged, reviewed by a pair other than A |
| CON-04 draft | A | 2.0 | Draft at `docs/contracts/CON-04.md`, DEC-09 options listed as open | Draft merged, reviewed by a pair other than A |
| SPK-05 | B | 3.0 | Report `docs/spikes/SPK-05.md` with evidence classes and the network-egress register | Report merged, reviewed by another pair |
| CON-02 completion (#23) | B | 1.5 | Completed contract `docs/contracts/CON-02.md` | **Full Done**: DEC-04 approved, dependent-pair agreement recorded, review and merge (§8.4) |
| INF-07 subset | C | 4.0 | `internal/observability` with tests and `docs/architecture/logging.md` | CI green, tests evidence criteria, merged, reviewed by another pair |
| CON-05 preliminary | D | 3.0 | `docs/contracts/CON-05.md`, status Preliminary | Preliminary draft merged, reviewed by another pair |
| DEC-11 ADR | D | 1.0 | `docs/architecture/adr/ADR-003.md`, Status Provisional | ADR merged as Provisional, reviewed by another pair |
| **Total committed** | | **17.5** | | |

**Estimates.** Backlog sizes: SPEC-F06 1.5, SPEC-F08 1.5, CON-02 1.5, CON-04 4 (split 2 + 2), CON-05 4 (split 3 + 1), SPK-05 3 (timebox), INF-07 4 (full scope, of which the Sprint 2 subset is the maximum timebox), INF-09 1.5, DEC-11 not sized in the backlog (1.0 provisional, mentor-accepted).

### 8.2 Conditional Stretch (not committed)

| Item | Pair | Pair-days | Conditions | Notes |
|------|------|-----------|------------|-------|
| INF-09 Auth0 development tenant | C | 1.5 | MA-05 confirmed; explicit mentor authorisation; Pair C capacity (1.0 remaining, PC-3) | Not in the 17.5. Stays in the backlog if not started. |

Other candidates from earlier drafts (SPK-04, INF-08, INF-12a, a SPEC-F04 review) are not in this plan. They would need a separate mentor decision.

### 8.3 Partial deliverables: what they are and what they are not

Each partial deliverable counts toward its committed pair-days. None is formal approval, and none is completed feature implementation.

| Item | Produces in Sprint 2 | Does not mean |
|------|----------------------|---------------|
| SPEC-F06 / SPEC-F08 draft | Draft specification, reviewed, merged | Approved specification. Approval is Sprint 3. |
| CON-04 draft | Draft contract with DEC-09 open | Final contract. Sprint 3 remainder 2.0. |
| CON-05 preliminary | Preliminary data contract for users, roles, submissions and results | Final approved data contract. Sprint 3 remainder 1.0. |
| CON-02 completion | Complete contract, awaiting approval and agreement | Done, until B1–B3 are met |
| SPK-05 | Evidence matrix and recommendation | DEC-08 approval, or DEC-02, or any production judge work |
| INF-07 subset | Logging, request-ID and access-log library, tests, usage docs | Health endpoints, metrics, HTTP server, F01-01 work |
| DEC-11 ADR | Provisional ADR with options and a recommendation | Approved migration tool. INF-05 is not started. |

**Draft-closure model (mentor-approved).** A draft issue may close when its explicitly defined draft deliverable meets the applicable review and merge requirements. Closing it does **not** mean the full specification or contract is formally Approved. The SDD approval gates and the Definition of Done for final artefacts are unchanged. SDD Workflow §6 now states how a draft-only issue is closed. GitHub labels are type-only (sprint-02-issues.md §2.2). This clarification (PC-1) is made in sdd-workflow.md §6.

### 8.4 CON-02 formal Done: blockers

| # | Blocker | Owner | Status |
|---|---------|-------|--------|
| B1 | DEC-04 approved by the mentor (MA-26) | Mentor | **Open** |
| B2 | Every dependent pair has recorded its agreement (owners of F06-02 and F07-01 are not yet named, MA-28) | Mentor and each dependent pair | **Open** |
| B3 | Pull request approved by a pair other than B, and merged | Pair B and reviewer pair | Not started |

CON-02 is not formally Done until B1–B3 are met. This plan does not assume any of them. PC-2 asks whether the draft-closure model should extend to CON-02. This plan assumes it does not.

### 8.5 Sprint 3 follow-up (planning only)

No successor issue is created by this plan. The work below is recorded here and in sprint-02-issues.md §8.1.

| Item | Sprint 3 work | Remaining estimate | Gate |
|------|---------------|--------------------|------|
| CON-04 final | Completion and approval | **2.0** (approved remainder) | DEC-09 decided; dependent-pair agreement |
| CON-05 final | Completion and approval | **1.0** (approved remainder) | CON-04 final; SPEC-F01 OQ-1 and OQ-2 for users and roles |
| SPEC-F06 approval | Formal approval | Not separately sized (review overhead) | CON-02 agreement; approval checklist |
| SPEC-F08 approval | Formal approval | Not separately sized (review overhead) | Approval checklist |
| CON-02 remainder (#23) | Anything left if B1–B3 are not met by 27 November | Not sized. Mentor to size. | DEC-04; dependent-pair agreement |
| INF-07 remainder | Health endpoints, basic metrics, HTTP server wiring | Not sized. Mentor to size. | Mentor sizing |
| INF-09 (ST-01) | If not started in Sprint 2 | 1.5 (backlog) | MA-05; authorisation; Pair C capacity |
| ADR-003 approval | Mentor approval of the migration tool | No task effort | MA-23 |
| F01-01 | Preparation and implementation | Backlog 4, not reduced | SPEC-F01 Approved; INF-09 done |
| DEC-12 | Dispatch decision | Not a task | By 30 November |
| DEC-08 | Isolation decision | Not a task | Start of Sprint 4 |

---

## 9. Review Responsibilities

CONTRIBUTING §9 requires at least one approval from another pair for every pull request. The author pair is never its own reviewer. The reviewer pairs below are provisional, and the mentor confirms them.

| Item | Author | Reviewer(s) | Dependent pairs |
|------|--------|-------------|-----------------|
| SPEC-F06 draft | A | B and D | C (F06-01 depends on it) |
| SPEC-F08 draft | A | C | D (CON-05 entities) |
| CON-04 draft | A | D | C (F01-01 and F07-04 later) |
| SPK-05 report | B | A or D (not B) | C and D (DEC-08 input) |
| CON-02 | B | D | C (F06-02 owner, Sprint 3) and B (F07-01 owner) |
| INF-07 subset | C | A or D (not C) | — |
| CON-05 preliminary | D | A and C | — |
| DEC-11 ADR | D | C | A (CON-04 and migration use) |
| SPEC-F01 approval (carried) | D (author) | A and C (required) | — |

Reviews are organisational time and are reported outside the 17.5 task figure (§3.3).

---

## 10. Membership

The mentor assigns students to pairs personally (confirmed). This plan does not assign, suggest or infer any individual's pair. The repository holds no pair membership data, so none is used. The mentor's confirmed composition is required before Sprint 2 issues are created, and before reviewers are named. Pair composition and GitHub assignees are finalised during Sprint Planning, and no Issue has an assignee.

---

## 11. Milestones and Decisions

### 11.1 Milestones

| Day | Date | Milestone | Owner |
|-----|------|-----------|-------|
| 1 | Mon 16 Nov | Mentor confirms pair composition (MA-28). Decisions recorded: PC-1, PC-2, PC-3, MA-05, MA-25 (SPEC-F01 OQ-1 and OQ-2 owners). | Mentor |
| 2 | Tue 17 Nov | CON-02 draft completion started. Dependent-pair names requested (B2). DEC-04 decision requested (MA-26). | B, mentor |
| 3 | Wed 18 Nov | SPEC-F06 and SPEC-F08 drafts started. INF-07 design checked against `log/slog`. DEC-11 options recorded. | A, C, D |
| 4 | Thu 19 Nov | SPK-05 matrix started, with evidence classes. CON-04 and CON-05 drafts started. | B, A, D |
| 5 | Fri 20 Nov | **Demo.** Interim: INF-07 tests, SPK-05 matrix, drafts. Mid-sprint check on absences. | All |
| 6 | Mon 23 Nov | CON-02 agreement requested. INF-07 pull request open. | B, C |
| 7 | Tue 24 Nov | SPEC-F06 and SPEC-F08 drafts submitted for review (status In Review in the file). CON-05 preliminary in review. | A, D |
| 8 | Wed 25 Nov | INF-07 reviewed and merged. CON-04 draft in review. | C, A |
| 9 | Thu 26 Nov | SPK-05 report merged. DEC-11 ADR in review. DEC-12 recommendation reviewed. | B, D |
| 10 | Fri 27 Nov | **Sprint 2 demo.** Sprint ends at 15:00. Closing conditions checked per issue. | All |

### 11.2 Exit criteria (27 November, 15:00)

- SPK-05 report merged to `docs/spikes/SPK-05.md`, with criteria SC-A to SC-C recorded, evidence classes stated, and no empty matrix cells.
- INF-07 subset merged, CI green, with tests evidencing criteria 1–11 of sprint-02-issues.md S2-06.
- SPEC-F06 and SPEC-F08 drafts merged and reviewed by another pair. Their status matches the file. Neither is recorded as Approved.
- CON-04 draft and CON-05 preliminary merged and reviewed by another pair. Each is marked Draft or Preliminary.
- DEC-11 ADR-003 merged as Provisional.
- CON-02 (#23): Done if B1–B3 are met. Otherwise open with the blocker named.
- No production judge code merged. No DEC-02 or DEC-08 approval recorded.
- INF-09 (stretch) is not started unless all its conditions are met and the mentor has authorised it.

### 11.3 Mentor decisions

**Confirmed or approved by the mentor:**

| ID | Decision | Status |
|----|----------|--------|
| — | Sprint 2 dates: 16–27 November 2026 (10 working days) | **Approved** |
| MA-19 | Maximum committed task workload: 20 pair-days across all four pairs, excluding meetings, reviews, demos and organisational overhead | **Approved** |
| MA-20 | Committed task workload: 17.5 pair-days (§3.2), with 2.5 uncommitted | **Approved** |
| MA-22 | SonarQube Cloud identifier: **INF-15**. INF-05 and INF-06 keep their backlog meanings. | **Approved** |
| MA-28 | Student rotation is required. The mentor assigns students to pairs personally. | **Approved** (process). No individual assignment recorded here. |
| MA-29 | INF-07 subset replaces F01-01 preparation (S2-06). F01-01 excluded from Sprint 2 | **Approved** |
| MA-31 | Partial-deliverable and draft-closure model (§8.3) | **Approved in principle.** SDD clarification PC-1 is pending. |
| MA-32 | No separate ADR for `log/slog`, provided no new architectural decision is introduced | **Approved** |
| MA-33 | INF-07 subset maximum timebox: 4.0 pair-days | **Approved** |
| — | INF-09 as conditional Stretch (1.5), outside the 17.5 | **Approved** |
| — | DEC-11 provisional estimate of 1.0 pair-day for drafting ADR-003. ADR stays Provisional until MA-23. | **Approved** (estimate only) |
| — | Partial splits: CON-04 2.0 Sprint 2 and 2.0 Sprint 3; CON-05 3.0 Sprint 2 and 1.0 Sprint 3 | **Approved** |

**Open. These are not approved, and this plan does not decide them.**

| ID | Decision | Needed by | Note |
|----|----------|-----------|------|
| MA-05 | Auth0 tenant ownership and admin access | Before any INF-09 work | Gates the INF-09 stretch only |
| PC-3 | Pair C capacity for INF-09: 5.5 against a working figure of 5.0 | Before any INF-09 work | Mentor decides: accept 5.5, reassign, or wait for Sprint 3 |
| PC-1 | SDD clarification: closing draft-only issues, and the status-label rules | Resolved in this step | sdd-workflow.md §6 and CONTRIBUTING.md updated (sprint-02-issues.md §7) |
| Type labels | `type:feature`, `type:contract`, `type:technical-task`, `type:spike` and `type:adr` verified to exist. Mapping in sprint-02-issues.md §2.2. | Resolved | sprint-02-issues.md §2.2 |
| PC-2 | Whether the draft-closure model extends to CON-02 | Day 1 | This plan assumes it does not |
| MA-23 | DEC-11 migration tool approval | Day 3 | ADR-003 stays Provisional until then |
| MA-14 | DEC-09 (problem packaging, test visibility) | Before CON-04 final | Blocks CON-04 final and SPEC-F03 |
| MA-25 | SPEC-F01 OQ-1 (role source) and OQ-2 (DEC-03 approval) | Before SPEC-F01 approval | Not in Sprint 2 scope; gates F01-01 later |
| MA-26 | DEC-04 approval | Day 2 | Blocks CON-02 formal Done (B1) |
| MA-27 | DEC-12 dispatch decision (SPK-03 recommends the outbox) | Day 9 (by 30 November) | Needed for Sprint 3 dispatch |
| MA-12 | Any SPK-05 timebox extension | As needed | Not expected |

### 11.4 Decision status

| ID | Decision | Status |
|----|----------|--------|
| DEC-01 | Evaluation through the Go testing framework | Provisional. SPK-01 is No-go for the execution path. |
| DEC-02 | Execution with Kubernetes Jobs and client-go | **Not approved** (SPK-01 No-go) |
| DEC-03 | Authentication with Auth0 | Provisional. SPK-02 Go. Needs MA-05 and OQ-2. |
| DEC-04 | Messaging with RabbitMQ and Watermill | Provisional. SPK-03 Go. Needs MA-26. |
| DEC-08 | Isolation mechanism | Open. SPK-05 gives evidence only. Needed at start of Sprint 4. |
| DEC-09 | Problem packaging and test visibility | Open |
| DEC-11 | Database migration tool | Open. ADR-003 drafted Provisional. Needs MA-23. |
| DEC-12 | Dispatch consistency | Open. SPK-03 recommends the outbox. Needed by 30 November. |
| DEC-13 | Editor session preservation | Open. Recorded in SPEC-F04. Not Sprint 2. |

No decision in the repository is approved, apart from the mentor's Sprint 2 planning decisions recorded in §11.3. ADR-001 and ADR-002 are Provisional.

---

## 12. Risks

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| RSK-S2-A | CON-02 Done depends on decisions outside Pair B | High | High | B1–B3 named in §8.4. Criteria not weakened. |
| RSK-S2-B | SPEC-F01 cannot be approved, so F01-01 has no approved specification | High | High | OQ-1 and OQ-2 (MA-25). F01-01 is outside Sprint 2. |
| RSK-S2-C | No approved isolation path after the SPK-01 No-go | High | High | SPK-05 evidence. No judge implementation. Option E (internal-only release) remains (backlog §11.5). |
| RSK-S2-D | Pair C cannot take INF-09 within 5.0 | High | Medium | PC-3. INF-09 is not committed. |
| RSK-S2-E | Examinations reduce attendance | High | Medium | Absences tracked at the mid-sprint check |
| RSK-S2-F | Identifier confusion: issue #21 was titled "[INF-05]" for SonarQube | Medium | Low | Resolved. Issue #21 is retitled "[INF-15] SonarQube…" (§15, F-1) |
| RSK-S2-G | Draft issues closed without a clear closure rule, read as approval | Medium | High | PC-1. Closing comment states "draft closed, not approved". |
| RSK-S2-H | INF-07 subset scope creeps into health endpoints or metrics | Medium | Medium | Explicit exclusions in sprint-02-issues.md S2-06 |

---

## 13. Stretch Rules

- The conditional stretch (INF-09) starts only when all its conditions in §4.3 hold, and the mentor has authorised it in writing.
- Stretch items are not exit criteria.
- An item not finished by 27 November, 15:00 is carried over and is not reported as completed.

---

## 14. Files Changed

- `docs/planning/sprint-02.md`: created in the first step. Revised in each later step. Finalised in this step.
- `docs/planning/sprint-02-issues.md`: finalised in this step.
- `docs/product/product-backlog.md`: §10.2 added in this step. §10.1 marked superseded for its figures. No backlog estimate or item changed.

Not changed: GitHub Issues, specifications, contracts, ADRs, application code, the SDD workflow, `sprint-01.md` (beyond the earlier synchronisation), and `.claude/`. Nothing has been committed, pushed or merged.

---

## 15. Follow-Up Actions

| # | Follow-up | Owner | Blocks |
|---|-----------|-------|--------|
| F-1 | **Resolved.** Issue #21 is retitled "[INF-15] SonarQube…". | — | — |
| F-2 | Issues #10 (SPEC-F01) and #24 (SPEC-F04) are closed, but their specifications are In Review. Decide on reopening, or record the state under PC-1. | Mentor, Pairs D and A | Accurate state |
| F-3 | `sprint-01.md` statuses are synchronised for §4.6, §8 and §11. Any remaining references are historical. | — | None |
| F-4 | Pair composition and GitHub assignees. The mentor finalises both during Sprint Planning. No Issue has an assignee. | Mentor | Sprint Planning |
| F-5 | CON-01 duplicate-submission question has no recorded answer. | Pair D, mentor | SPEC-F06 draft completeness |
| F-6 | SPEC-F01 open questions OQ-1 to OQ-9 are open. OQ-1 and OQ-2 are critical. | Pairs C and D, mentor for OQ-2 | SPEC-F01 approval |
| F-7 | **Resolved.** The committed Issues exist (#23, #29 to #35). Links are in sprint-02-issues.md §1.1. | — | — |
