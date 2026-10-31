# GoJudge Product Backlog

**Status:** Draft for mentor review. Not approved. No GitHub Issues have been created from this backlog.

**Source of scope:** [MVP Project Proposal](mvp-proposal.md), version 0.1. Feature titles, priorities and expected behaviour are preserved from that document. Where this backlog adds work, it is marked as a proposal.

**Related documents:** [SDD Workflow](../sdd-workflow.md), [Contributing Guide](../../CONTRIBUTING.md), [Feature Specification Template](../../specs/templates/feature-spec-template.md).

---

## 1. Purpose and Conventions

This backlog lists the work needed to deliver the GoJudge MVP between 2 November 2026 and 22 January 2027. It is a planning document. It does not approve architecture, technology or scope.

### 1.1 Scope boundary

- The committed MVP is **F01–F09**, all with priority **Must**.
- **F10–F12** (progress tracking, catalogue filters, leaderboard) are excluded from the committed MVP. They appear nowhere in this backlog. The proposal's scope trade-off is to defer them before weakening security.
- Items marked **Proposed** in the [MVP Proposal §12](mvp-proposal.md#12-decisions-requested-from-technical-expert) are not approved. Their decisions are listed in §8.

### 1.2 Item types

Each backlog item has exactly one type. The type determines what "done" means for it.

| Type | Prefix | Output | Done when |
|------|--------|--------|-----------|
| Specification | `SPEC-` | Feature spec in `specs/features/` | Status **Approved** by a reviewer pair, per [SDD Workflow §4](../sdd-workflow.md#4-approvals) |
| Contract | `CON-` | Agreed API, message, data or state contract | Agreed by every dependent pair and recorded in the relevant spec |
| Implementation | `F01-01` style | Code, tests and documentation on a feature branch | Meets the [Definition of Done](../../CONTRIBUTING.md#11-definition-of-done) |
| Infrastructure | `INF-` | Environment, tooling or deployment | Works as specified, is reproducible, and is documented |
| Test | `TST-` | Test asset, harness or verification run | Runs in the agreed environment with recorded results |
| Hardening | `HRD-` | Fixes or reviews before release | Acceptance evidence recorded for the release checklist |
| Spike | `SPK-` | Time-boxed investigation report | Report merged to `docs/` with a recommendation. Spike code is not production code and does not pass the Definition of Done. |
| Decision | `DEC-` | Recorded decision | Approved by the mentor and recorded in an ADR |
| Mentor approval | `MA-` | Mentor approval point | Mentor's approval is recorded in the relevant issue |
| Risk | `RSK-` | Tracked risk | Reviewed at each sprint boundary |

**Specification work, implementation work and spikes are separate.** A spike may inform a specification, but it never replaces one. No implementation item may start until its specification is approved.

### 1.3 Priority and size

- **Priority:** every committed item is **Must** (required for the MVP). No Should or Could items are committed.
- **Size:** a rough estimate in pair-days. **S** ≈ 1–2, **M** ≈ 3–5, **L** ≈ 6–8. A pair has roughly 10 pair-days in a two-week sprint. Pairs must validate these estimates in Sprint 1.

### 1.4 Sprint calendar

| Sprint | Dates | Proposal milestone (weeks) | Friday demos |
|--------|-------|----------------------------|--------------|
| 1 | 2–13 Nov 2026 | Architecture and vertical slice (1–2) | 6 Nov, 13 Nov |
| 2 | 16–27 Nov 2026 | Core platform (3–4) | 20 Nov, 27 Nov |
| 3 | 30 Nov–11 Dec 2026 | Submission pipeline (5–6) | 4 Dec, 11 Dec |
| 4 | 14–25 Dec 2026 | Judge integration (7–8) | 18 Dec, **25 Dec (public holiday)** |
| 5 | 28 Dec 2026–8 Jan 2027 | End-to-end MVP (9–10) | **1 Jan (public holiday)**, 8 Jan |
| 6 | 11–22 Jan 2027 | Release hardening (11–12) | 15 Jan, 22 Jan (final demo) |

Sprints 4 and 5 span the Christmas and New Year period. Capacity and demo arrangements are a mentor decision (MA-07).

---

## 2. MVP Features

These rows preserve the proposal's requirements and priorities. Implementation items are listed in §3 and sequenced in §6.

| ID | Feature | Priority | Expected behaviour (from proposal) | Specification item | Implementation items |
|----|---------|----------|------------------------------------|--------------------|----------------------|
| F01 | Accounts and access | Must | Register, sign in/out; learner and administrator roles. | SPEC-F01 | F01-01, F01-02, F01-03 |
| F02 | Problem catalogue | Must | List published problems with title, difficulty and solved state. | SPEC-F02 | F02-01, F02-02, F02-03 |
| F03 | Problem details | Must | Display statement, input/output contract, constraints and public examples. | SPEC-F03 | F03-01, F03-02 |
| F04 | Browser code editor | Must | Edit Go source code; load a starter template; preserve current work in the active session. | SPEC-F04 | F04-01, F04-02 |
| F05 | Run sample tests | Must | Execute code against public sample tests and return per-case feedback. | SPEC-F05 | F05-01, F05-02 |
| F06 | Submit solution | Must | Create asynchronous submission; return submission ID and processing state. | SPEC-F06 | F06-01 to F06-04 |
| F07 | Automated judge | Must | Compile and execute submitted Go code in isolation; evaluate hidden tests; produce verdict. | SPEC-F07 | F07-01 to F07-08 |
| F08 | Submission history | Must | Show own submissions, timestamps, states and final verdicts. | SPEC-F08 | F08-01, F08-02 |
| F09 | Problem administration | Must | Admin CRUD, draft/published status, sample and hidden test management. | SPEC-F09 | F09-01 to F09-03 |

---

## 3. Backlog Items

### 3.1 Specifications and contracts

Specifications must be **Approved** before the implementation items they gate. The "Gates" column names the implementation items that cannot start until approval.

| ID | Item | Type | Priority | Size | Depends on | Deliverable | Gates |
|----|------|------|----------|------|------------|-------------|-------|
| SPEC-F01 | Accounts and access specification | Specification | Must | S | CON-01 (draft) | Approved spec in `specs/features/` | F01-01 |
| SPEC-F02 | Problem catalogue specification | Specification | Must | S | CON-05 (draft) | Approved spec | F02-01 |
| SPEC-F03 | Problem details specification | Specification | Must | S | CON-04 | Approved spec | F03-01 |
| SPEC-F04 | Browser code editor specification | Specification | Must | S | DEC-13 | Approved spec | F04-01 |
| SPEC-F05 | Run sample tests specification | Specification | Must | S | DEC-10, SPK-01 | Approved spec | F05-01 |
| SPEC-F06 | Submit solution specification | Specification | Must | S | CON-01, CON-02 | Approved spec | F06-01 |
| SPEC-F07 | Automated judge specification | Specification | Must | M | SPK-01, SPK-05, CON-01, CON-02 | Approved spec with isolation requirements | F07-01, F07-02 |
| SPEC-F08 | Submission history specification | Specification | Must | S | CON-01 | Approved spec | F08-01 |
| SPEC-F09 | Problem administration specification | Specification | Must | S | CON-04, CON-03 | Approved spec | F09-01 |
| CON-01 | Submission lifecycle and verdict model | Contract | Must | S | None | States (QUEUED, RUNNING, COMPLETED, FAILED), transitions, verdicts, infrastructure vs test-failure classification | SPEC-F06, SPEC-F07, F06-01, F07-05 |
| CON-02 | Messaging contract | Contract | Must | S | CON-01, SPK-03 | Queue names, message schema, idempotency key, acknowledgement, retry and dead-letter behaviour | F06-02, F07-01 |
| CON-03 | OpenAPI contract for MVP endpoints | Contract | Must | M | SPEC-F01, SPEC-F02, SPEC-F03, SPEC-F06 | OpenAPI 3.x document, extended as each feature is approved | INF-06, all API items |
| CON-04 | Problem packaging format and test visibility rules | Contract | Must | M | DEC-09 | Packaging format, starter template format, visibility rules | SPEC-F03, F07-04, F09-02 |
| CON-05 | Data model | Contract | Must | M | CON-01, CON-04 | Entities for users, roles, problems, test assets, submissions and results | INF-05, F01-02, F02-01 |

### 3.2 Infrastructure and technical work

| ID | Item | Type | Priority | Size | Depends on | Deliverable |
|----|------|------|----------|------|------------|-------------|
| INF-01 | Go module, repository layout, lint and format configuration | Infrastructure | Must | S | DEC-14 | Module, layout, lint configuration. Layout is proposed and subject to mentor approval. |
| INF-02 | CI baseline on GitHub Actions | Infrastructure | Must | S | INF-01 | Build, lint, unit tests, dependency and security checks |
| INF-03 | Local infrastructure | Infrastructure | Must | M | DEC-06 | Local Kubernetes cluster (k3d or Minikube), PostgreSQL and RabbitMQ for development. Setup steps documented in the README. |
| INF-04 | ADR process and initial records | Infrastructure | Must | S | None | ADR template and records for approved decisions |
| INF-05 | PostgreSQL schema, migrations and data access | Infrastructure | Must | M | CON-05, INF-03, DEC-11 | Migrations, pgx and sqlc setup |
| INF-06 | OpenAPI validation in CI | Infrastructure | Must | S | CON-03, INF-02 | Contract validation step in CI |
| INF-07 | Structured logging, health endpoints and basic metrics | Infrastructure | Must | M | INF-01 | Logging standard, health endpoints, metrics endpoint |
| INF-08 | Frontend scaffold | Infrastructure | Must | S | INF-02 | React, TypeScript and Vite project with Monaco Editor dependency and CI build |
| INF-09 | Auth0 development tenant configuration | Infrastructure | Must | S | DEC-03, SPK-02, MA-05 | Development tenant and application configuration. No secrets in the repository. |
| INF-10 | Judge namespace, RBAC, quotas and Pod security | Infrastructure | Must | L | SPK-01, SPK-05, DEC-08, MA-03 | Namespace isolation, RBAC, resource quotas and Pod security controls, validated |
| INF-11 | Staging environment and deployment | Infrastructure | Must | L | DEC-07, MA-04, INF-02, INF-05 | Repeatable staging deployment |
| INF-12 | Container images for API, orchestrator and judge runner | Infrastructure | Must | M | INF-02, F07-01 | Images with dependency and image security checks in CI |
| INF-13 | Operator README and runbook | Infrastructure | Must | S | INF-11 | Operator documentation and known-limitations list |
| INF-14 | Release checklist and final demo preparation | Infrastructure | Must | S | TST-04 | Checklist and demo script |

### 3.3 Spikes

Spikes are time-boxed. Their output is a report with a recommendation, not production code. Spikes must not run learner-submitted code outside the isolated spike environment.

| ID | Item | Type | Priority | Size | Depends on | Deliverable |
|----|------|------|----------|------|------------|-------------|
| SPK-01 | Secure Go execution feasibility with `go test` and Kubernetes Jobs | Spike | Must | L | INF-03, DEC-01, DEC-02, DEC-06 | Report covering: a `go test` run in a Job created by client-go; CPU, memory and time limits; no network; unprivileged execution; restricted filesystem; log and report collection before cleanup; TTL cleanup; cancellation. Tests use trusted packages only, with negative cases (network access, filesystem writes, infinite loop, memory exhaustion, output flood). Go or no-go recommendation. |
| SPK-02 | Auth0 OIDC integration check | Spike | Must | S | DEC-03 | Sign-in flow and token validation in Go; role claim mapping. Report with recommendations. |
| SPK-03 | RabbitMQ and Watermill check | Spike | Must | M | DEC-04, INF-03 | Publish and consume; acknowledgement, redelivery and dead-letter behaviour; idempotency handling. Report. |
| SPK-04 | Hosting comparison (GKE and VM-hosted k3s) | Spike | Must | M | DEC-07 | Comparison of setup effort, cost, operations and fit with the isolation design. Input to DEC-07. |
| SPK-05 | Isolation mechanism comparison | Spike | Must | M | SPK-01 | Comparison of candidate isolation mechanisms against the proposal's isolation requirements. Input to DEC-08. Candidate mechanisms are listed for evaluation only. |

### 3.4 Feature implementation items

Each item requires its specification to be **Approved** first. Items are grouped by feature.

**F01 Accounts and access**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F01-01 | Sign-in, sign-out and token validation in the Go API | Must | M | SPEC-F01, INF-09, SPK-02 | Authentication middleware; tests with fixture tokens |
| F01-02 | Learner and administrator roles and endpoint authorisation | Must | M | F01-01, CON-05 | Role model; authorisation checks; tests for role boundaries |
| F01-03 | Sign-in and sign-out UI | Must | S | F01-01, INF-08 | React sign-in and sign-out flow |

**F02 Problem catalogue**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F02-01 | Problem persistence and published-status filter | Must | M | SPEC-F02, CON-05, INF-05 | Problem storage; published-only query |
| F02-02 | Catalogue list endpoint with solved state | Must | M | F02-01; solved state completed after F08-01 | List endpoint with title, difficulty and solved state. Solved state is stubbed until F08-01 is available. |
| F02-03 | Catalogue UI | Must | S | F02-02 | React catalogue list |

**F03 Problem details**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F03-01 | Problem details endpoint | Must | M | SPEC-F03, F02-01, CON-04 | Statement, input/output contract, constraints and public examples |
| F03-02 | Problem details UI | Must | S | F03-01 | React problem page |

**F04 Browser code editor**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F04-01 | Monaco Editor integration with Go mode and starter template | Must | M | SPEC-F04, F03-02, INF-08, CON-04 | Editor component; starter template loaded from problem data |
| F04-02 | Active-session preservation of current work | Must | S | F04-01, DEC-13 | Work preserved in the active session, per DEC-13 |

**F05 Run sample tests**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F05-01 | Run-samples endpoint through the isolated execution path | Must | L | SPEC-F05, DEC-10, F07-02, F07-03, F07-05 | Endpoint returning per-case feedback for public sample tests |
| F05-02 | Per-case feedback in the editor | Must | S | F05-01, F04-01 | UI for per-case results |

**F06 Submit solution**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F06-01 | Submission creation endpoint | Must | M | SPEC-F06, F01-02, F03-01, CON-01, CON-05 | Persisted submission; returns submission ID and QUEUED state |
| F06-02 | Dispatch to RabbitMQ through Watermill | Must | M | F06-01, CON-02, SPK-03, DEC-12 | Publisher using the agreed message contract |
| F06-03 | Submission state endpoint | Must | S | F06-01 | Endpoint returning the current state |
| F06-04 | Submit action and state polling in the editor | Must | S | F06-03, F04-01 | UI showing processing state |

**F07 Automated judge**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F07-01 | Orchestrator consumer with state transitions and mock executor | Must | L | SPEC-F07, F06-02, CON-01, CON-02 | Watermill consumer; QUEUED → RUNNING → COMPLETED or FAILED; controlled mock executor |
| F07-02 | Kubernetes Job creation through client-go | Must | L | F07-01, SPK-01, INF-10, DEC-08, MA-03 | One Job per submission with CPU, memory and time limits; unprivileged; no network; no secrets |
| F07-03 | Job monitoring, log and report collection, TTL cleanup and cancellation | Must | M | F07-02 | Monitoring; logs and reports collected before cleanup; cancellation support |
| F07-04 | Hidden test asset handling | Must | M | CON-04, INF-05, F07-02 | Test suites available only to the judge Job; never returned to learners |
| F07-05 | Verdict mapping and failure classification | Must | M | CON-01, F07-03 | Verdicts: ACCEPTED, WRONG_ANSWER, COMPILATION_ERROR, RUNTIME_ERROR, TIME_LIMIT_EXCEEDED, MEMORY_LIMIT_EXCEEDED, INTERNAL_ERROR. Infrastructure failures separated from test failures. |
| F07-06 | Idempotent processing and retry policy | Must | M | F07-01, F07-05 | Redelivery does not create duplicate final results; learner failures are not retried automatically |
| F07-07 | Adversarial security tests for the judge | Must | L | F07-02, F07-03, F07-04 | Tests for network access, filesystem writes, secret access, resource exhaustion, output flood and hidden test reads. Gate for any public exposure. |
| F07-08 | Worker failure and timeout integration tests | Must | M | F07-06, F07-03 | Tests for worker restart mid-run, Job failure and timeout |

**F08 Submission history**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F08-01 | Submission history endpoint | Must | M | SPEC-F08, F06-01, F07-05, F01-02 | Own submissions only, with timestamps, states and verdicts |
| F08-02 | Submission history UI | Must | S | F08-01, F04-01 | React history page |

**F09 Problem administration**

| ID | Item | Priority | Size | Depends on | Deliverable |
|----|------|----------|------|------------|-------------|
| F09-01 | Admin CRUD endpoints with draft and published status | Must | M | SPEC-F09, F01-02, F02-01, CON-03 | Endpoints restricted to administrators |
| F09-02 | Sample and hidden test management | Must | L | F09-01, F07-04, CON-04 | Upload, store and replace test assets |
| F09-03 | Admin UI | Must | M | F09-01, F09-02 | Problem and test management screens |

### 3.5 Test and hardening items

| ID | Item | Type | Priority | Size | Depends on | Deliverable |
|----|------|------|----------|------|------------|-------------|
| TST-01 | Reference challenge fixtures | Test | Must | M | CON-04 | One algorithmic and one HTTP-handler challenge, with reference solutions and negative cases. Required for release acceptance. |
| TST-02 | End-to-end test harness | Test | Must | M | F05-02, F08-02 | Browser-driven journey: browse, edit, run samples, submit, view verdict, view history |
| TST-03 | Concurrency test on staging | Test | Must | M | INF-11, F07-06, DEC-05 | At least two concurrent submissions through Kubernetes Jobs; correct results; no duplicate final evaluations |
| TST-04 | Release acceptance run | Test | Must | M | TST-02, TST-03, HRD-01, HRD-02 | Evidence for each release acceptance criterion |
| HRD-01 | Security review and fixes for the judge | Hardening | Must | M | F07-07, INF-10 | Dedicated security review of the isolation mechanism; findings fixed or accepted by the mentor |
| HRD-02 | Defect fixes for critical and high severity | Hardening | Must | M | TST-02, TST-03 | No unresolved critical or high defects |

---

## 4. Dependency Map

The critical path runs through the judge's isolation work.

```text
SPK-01 → DEC-08 (MA-03) → SPEC-F07 approved → F07-02 → F07-03 → F07-05 → F05-01 → F05-02 → TST-02 → TST-04
                                                          ↘ F07-06 → F07-08
SPK-03 → DEC-12 → CON-02 → F06-02 → F07-01 → F07-06
DEC-03 → SPK-02 → INF-09 → F01-01 → F01-02 → F06-01 / F09-01 / F08-01
CON-01 → CON-05 → INF-05 → F02-01 → F03-01 → F06-01
DEC-08 approval and HRD-01 gate public exposure.
```

---

## 5. Sprint 1 Backlog (2–13 November 2026)

Sprint 1 has four parallel workstreams, one per pair. Pair labels A–D are placeholders. The mentor assigns pairs at the start of the sprint.

Sprint 1 begins with a decision session on day one. Decisions DEC-01, DEC-02, DEC-03, DEC-04, DEC-06 and DEC-14 must be marked **Approved** or changed before the matching spikes and infrastructure work proceed (MA-02).

### 5.1 Workstream overview

| Workstream | Pair | Items | Output by 13 Nov | Demo focus |
|------------|------|-------|------------------|------------|
| WS1 Platform foundations | Pair A | INF-01, INF-02, INF-03, INF-04, INF-08 | CI green on baseline; local cluster and databases running; first ADRs; frontend scaffold | Working baseline |
| WS2 Secure execution spike | Pair B | SPK-01 | Spike report with go or no-go and recommended controls | Negative test results |
| WS3 Identity and messaging spikes | Pair C | SPK-02, SPK-03 | Two spike reports with recommendations | Auth flow and redelivery behaviour |
| WS4 Specifications and contracts | Pair D | SPEC-F01, SPEC-F02, CON-01, CON-04 (draft), DEC-09 proposal | Specs ready for cross-pair review; agreed lifecycle and verdict model | Specs and contracts |

### 5.2 Workstream details

**WS1 Platform foundations (Pair A)**

| ID | Objective | Acceptance criteria | Depends on |
|----|-----------|---------------------|------------|
| DEC-14 | Confirm Go version pin | Version recorded in an ADR | MA-02 |
| INF-01 | Module and layout | `go.mod` present; lint and format configured; layout reviewed | DEC-14 |
| INF-02 | CI baseline | Pipeline runs on push; build, lint and tests run; green on `main` | INF-01 |
| INF-03 | Local infrastructure | Local cluster and PostgreSQL and RabbitMQ run; steps documented | DEC-06 |
| INF-04 | ADR process | Template merged; ADRs for approved DEC items | None |
| INF-08 | Frontend scaffold | Build runs in CI; Monaco dependency installed | INF-02 |

Pair A's first task is INF-03, because Pair B and Pair C need the local cluster early in the sprint.

**WS2 Secure execution spike (Pair B)**

| ID | Objective | Acceptance criteria | Depends on |
|----|-----------|---------------------|------------|
| SPK-01 | Test `go test` in a Kubernetes Job created with client-go | Job runs trusted test packages; limits enforced; network and filesystem restrictions verified; logs collected before cleanup; TTL cleanup verified; report with go or no-go | INF-03 (local cluster), DEC-01, DEC-02 |

No learner-submitted code runs in the spike. Pair B may use an interim local cluster from INF-03 if the formal setup is not ready.

**WS3 Identity and messaging spikes (Pair C)**

| ID | Objective | Acceptance criteria | Depends on |
|----|-----------|---------------------|------------|
| SPK-02 | Auth0 OIDC integration check | Sign-in, token validation and role claim mapping demonstrated; report | DEC-03, MA-05 |
| SPK-03 | RabbitMQ and Watermill check | Publish and consume demonstrated; redelivery and dead-letter behaviour recorded; report | DEC-04, INF-03 |

**WS4 Specifications and contracts (Pair D)**

| ID | Objective | Acceptance criteria | Depends on |
|----|-----------|---------------------|------------|
| CON-01 | Submission lifecycle and verdict model | Lifecycle, transitions and verdicts agreed with Pairs B and C | None |
| CON-04 | Problem packaging (draft) | Draft format and visibility rules reviewed; DEC-09 options listed | DEC-09 |
| SPEC-F01 | Accounts and access specification | Passes review: acceptance criteria numbered, security considerations documented | CON-01 |
| SPEC-F02 | Problem catalogue specification | Passes review: acceptance criteria numbered, dependencies on F08 documented | CON-05 draft |

### 5.3 Sprint 1 exit criteria

- Spike reports for SPK-01, SPK-02 and SPK-03 are merged to `docs/` with recommendations.
- The CI baseline (INF-02) is green on `main`.
- SPEC-F01 and SPEC-F02 are in cross-pair review, with reviews assigned to pairs other than the author's.
- CON-01 is agreed by all dependent pairs.
- Decisions DEC-01 to DEC-04, DEC-06 and DEC-14 are recorded as Approved or changed.
- The 13 November demo presents the spike outcomes and the CI baseline.

Cross-pair reviews of the Sprint 1 specifications consume capacity from all four pairs. The plan allows for this, but it is a known risk (RSK-05).

---

## 6. Proposed Sequence Across Six Sprints

The sequence respects the dependency map in §4. Specification approvals are gates: an implementation item starts only after its specification is approved.

| Sprint | Specifications approved at start | Main items | Exit criteria |
|--------|----------------------------------|------------|---------------|
| **1** (2–13 Nov) | None | INF-01 to INF-04, INF-08; SPK-01 to SPK-03; SPEC-F01, SPEC-F02 (review); CON-01; CON-04 (draft); DEC-01 to DEC-04, DEC-06, DEC-09, DEC-14 | Sprint 1 exit criteria met (§5.3) |
| **2** (16–27 Nov) | SPEC-F01, SPEC-F02 | CON-02, CON-03 (initial), CON-04 (final), CON-05; INF-05, INF-06, INF-09; F01-01, F01-02, F02-01; SPEC-F03, SPEC-F04, SPEC-F06, SPEC-F07 (draft); SPK-04; SPK-05 start; DEC-11 | Auth and catalogue data working; F07 spec in review; SPK-05 results ready |
| **3** (30 Nov–11 Dec) | SPEC-F03, SPEC-F04, SPEC-F06, SPEC-F07 (F07 requirements approved) | F01-03; F02-02 (solved state stubbed until F08-01), F02-03; F03-01, F03-02; F04-01; F06-01 to F06-03; F06-02 (after DEC-12); F07-01 with mock executor; TST-01; INF-07; INF-12 (begin); SPEC-F05, SPEC-F08, SPEC-F09 (draft); DEC-08 decided by sprint end; DEC-12, DEC-13 | Submission reaches QUEUED and is consumed; catalogue and details visible; DEC-08 recorded |
| **4** (14–25 Dec) | SPEC-F05, SPEC-F08, SPEC-F09; DEC-07 and DEC-08 approved; MA-03 and MA-04 recorded | F07-02 to F07-06; F04-02; F06-04; F08-01; F09-01; INF-10; INF-11 (begin, if DEC-07 approved); DEC-10 | Judge runs a Job per submission; verdicts recorded; history API working. Capacity reduced by holidays. |
| **5** (28 Dec–8 Jan) | None | F05-01, F05-02; F07-07, F07-08; F02-02 completed (solved state); F08-02; F09-02, F09-03; INF-11 completed; TST-02; TST-03 | Run-samples works; adversarial tests reported; staging concurrency test run |
| **6** (11–22 Jan) | None | HRD-01; HRD-02; TST-04; INF-12 (complete); INF-13; INF-14; final demo on 22 January | Release acceptance criteria evidenced; security review complete; no open critical or high defects |

---

## 7. Capacity Check

This is a rough estimate. It must be validated by the pairs in Sprint 1.

| Category | Items | Estimated pair-days |
|----------|-------|---------------------|
| Specifications and contracts | SPEC-F01 to SPEC-F09, CON-01 to CON-05 | ~29 |
| Infrastructure | INF-01 to INF-14 | ~44 |
| Spikes | SPK-01 to SPK-05 | ~22 |
| Features | F01-01 to F09-03 | ~116 |
| Test and hardening | TST-01 to TST-04, HRD-01 to HRD-02 | ~22 |
| **Total** | | **~233** |

Capacity:

- Four pairs over six sprints give about 240 pair-days gross.
- Reserving about 25% for cross-pair reviews, demos, onboarding and the holiday period leaves about 180 pair-days net.

**The backlog is over the net capacity by roughly 30%.** The mentor must decide how to handle this (MA-08). Candidate reductions, in the order they should be considered, are:

1. Simplify the admin UI (F09-03) to the minimum needed for the release criteria.
2. Reduce the end-to-end harness (TST-02) to the critical journey only.
3. Use the simplest acceptable option for active-session preservation (F04-02).
4. Limit per-case feedback (F05-02) to pass or fail per case.
5. Reduce the scope of HRD-02 to the defects found in TST-04.

Security items (F07-07, HRD-01, INF-10) and the public-exposure gate must not be reduced.

---

## 8. Decisions

Status values match the [MVP Proposal §12](mvp-proposal.md#12-decisions-requested-from-technical-expert) where a decision appears there. DEC-09 to DEC-14 are added by this backlog and are not in the proposal.

| ID | Decision | Status | Needed by | Source |
|----|----------|--------|-----------|--------|
| DEC-01 | Submission evaluation through the Go testing framework | Proposed | Sprint 1 day 1 | Proposal §12 |
| DEC-02 | Execution orchestration with Kubernetes Jobs and client-go | Proposed | Sprint 1 day 1 | Proposal §12 |
| DEC-03 | Authentication with Auth0 | Proposed | Sprint 1 day 1 | Proposal §12 |
| DEC-04 | Messaging with RabbitMQ and Watermill | Proposed | Sprint 1 day 1 | Proposal §12 |
| DEC-05 | Concurrency target of at least two concurrent submissions | Proposed | Sprint 5 | Proposal §12 |
| DEC-06 | Local infrastructure: k3d or Minikube | Proposed | Sprint 1 day 1 | Proposal §12 |
| DEC-07 | Hosting: GKE or VM-hosted k3s | Open | Sprint 4 | Proposal §12 |
| DEC-08 | Execution isolation mechanism | Open | End of Sprint 3 | Proposal §12 |
| DEC-09 | Problem packaging and test visibility rules | Open | Sprint 1 (draft), Sprint 2 (final) | Proposal §3 says "finalised during architecture design"; not in §12 |
| DEC-10 | Run-samples uses the same isolated execution path as submit | Assumption | Sprint 4 | Not in proposal; assumption to confirm |
| DEC-11 | Database migration tool | Open | Sprint 2 | Not in proposal |
| DEC-12 | Dispatch consistency between database write and message publish | Open | Sprint 3 | Not in proposal; reliability requirement in §6 |
| DEC-13 | Editor session-preservation mechanism | Open | Sprint 3 | Not in proposal; F04 expected behaviour |
| DEC-14 | Go version pin | Open | Sprint 1 day 1 | Proposal §5 says "current supported stable release" |

---

## 9. Risks

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| RSK-01 | Judge isolation takes longer than two sprints | High | High | Start SPK-01 and SPK-05 in Sprint 1 and 2; use the mock executor so the pipeline can progress; mentor decides on release scope if isolation is not validated (proposal §7: remain internal-only or use a vetted external sandbox) |
| RSK-02 | Backlog is over net capacity by roughly 30% | High | High | MA-08; reductions in §7; security items protected |
| RSK-03 | Holiday period reduces capacity in Sprints 4 and 5 | High | Medium | Plan lighter loads in these sprints; MA-07 |
| RSK-04 | Pair rotation every two weeks disrupts features mid-way | Medium | High | Handover rule in the SDD Workflow §3; split features into items no larger than one sprint where possible |
| RSK-05 | Cross-pair specification reviews become a bottleneck | Medium | Medium | Review assignments made at sprint start; reviews counted in capacity |
| RSK-06 | Auth0 tenant access is delayed | Medium | High | MA-05 in Sprint 1; SPK-02 uses a fixture until access arrives |
| RSK-07 | Hidden tests leak through logs, outputs or API responses | Medium | High | F07-04 and F07-07 test this; output bounded; logs reviewed before collection |
| RSK-08 | Interns are new to Kubernetes, RabbitMQ and client-go | High | Medium | Spikes in Sprint 1; pairing; mentor review of spikes |
| RSK-09 | Hosting approval or budget delays staging | Medium | Medium | DEC-07 by Sprint 4; fallback to internal-only deployment |
| RSK-10 | Judge design depends on DEC-08, which is open | Medium | High | Decision needed by end of Sprint 3; SPK-05 prepares options |
| RSK-11 | Demo falls on a public holiday | Certain | Low | MA-07 decides whether to move the 25 December and 1 January demos |

---

## 10. Mentor Approval Points

| ID | Approval | Needed by | Related items |
|----|----------|-----------|---------------|
| MA-01 | Confirm MVP priorities, feature order and the exclusion of F10–F12 | Before Sprint 1 | §2 |
| MA-02 | Approve or change DEC-01, DEC-02, DEC-03, DEC-04, DEC-06 and DEC-14 | Sprint 1 day 1 | §5 |
| MA-03 | Approve the isolation mechanism and its security requirements | Start of Sprint 4 | DEC-08, F07-02, INF-10 |
| MA-04 | Approve hosting and budget | Start of Sprint 4 | DEC-07, INF-11 |
| MA-05 | Confirm Auth0 tenant ownership and admin access | Sprint 2 | INF-09, SPK-02 |
| MA-06 | Approve any significant scope change | As needed | [SDD Workflow §7](../sdd-workflow.md#7-handling-requirement-changes) |
| MA-07 | Confirm holiday capacity, and the handling of the 25 December and 1 January demos | Sprint 3 | §1.4, RSK-03, RSK-11 |
| MA-08 | Accept or re-plan the capacity gap | Sprint 1 | §7, RSK-02 |
| MA-09 | Accept the final security review (HRD-01) before any public exposure | Sprint 6 | HRD-01, F07-07 |

---

## 11. Items Not Assigned and Not Created

- No item is assigned to a named intern.
- No GitHub Issues have been created from this backlog. Issues should be created when items enter Sprint planning, using the feature issue template and the labels in the [SDD Workflow §6](../sdd-workflow.md#6-github-issue-states).
- Specifications for F01–F09 do not yet exist in `specs/features/`. Each is created from the template, and is approved before its implementation items start.
