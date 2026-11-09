# GoJudge Product Backlog

**Status:** Draft for mentor review (revision 3). Not approved. No GitHub Issues have been created from this backlog.

**Source of scope:** [MVP Project Proposal](mvp-proposal.md), version 0.1. Feature titles, priorities, expected behaviour and mandatory security controls are preserved from that document. Where this backlog adds work, it is marked as a proposal.

**Related documents:** [SDD Workflow](../sdd-workflow.md), [Contributing Guide](../../CONTRIBUTING.md), [Feature Specification Template](../../specs/templates/feature-spec-template.md).

---

## 0. Revision Summary

Revision 3 applies the confirmed internship calendar and reassesses capacity against it. It does not remove any mandatory functionality, and it does not weaken any security requirement. No existing backlog item has been removed.

| Area | Change |
|------|--------|
| Calendar | Internship runs 2 November 2026 to 29 January 2027. Six sprints of 10 working days each (60 in total). Non-working days: 24–27 December 2026 and 31 December 2026 – 3 January 2027. Final internship demo on 29 January 2027. See §1.4. |
| Gross capacity | 4 pairs × 60 working days = **240 pair-days**. This is unchanged from revision 2, but it now rests on the confirmed calendar. |
| Realistic net capacity | **About 164 pair-days** (range 140–180), after SDD reviews, pull request review, training, meetings and integration support. Revision 2 used a flat 20% reserve, which understated these costs. See §11.3. |
| Backlog against capacity | **213.5 pair-days** of backlog. This is about **50 pair-days (31%) over** realistic net capacity. The plan does not fit. See §11.4. |
| Sprint 1 | Committed load is 28.5 pair-days against a realistic 24.5. This is 4 pair-days over, and it needs a mentor decision (MA-15). INF-08 moved to Sprint 3. See §9. |
| Plan | Six-sprint plan re-dated to the confirmed calendar, with per-sprint load against realistic net. The release-critical tail is about 51 pair-days and would not finish by 29 January (§10, §11.4). |
| Options | Simplifications, a meeting cadence change, added capacity, scope change, and an internal-only release are set out in §11.5. Scope is not reduced in this revision. |
| Unchanged | F01–F09 requirements, mandatory and security controls (§3.1), all item definitions (§4), spike criteria (§5), decision status (§6), critical path (§7), parallel streams (§8). |

---

## 1. Purpose and Conventions

This backlog lists the work needed to deliver the GoJudge MVP between 2 November 2026 and 29 January 2027. It is a planning document. It does not approve architecture, technology or scope.

### 1.1 Scope boundary

- The committed MVP is **F01–F09**, all with priority **Must**.
- **F10–F12** (progress tracking, catalogue filters, leaderboard) are excluded from the committed MVP. They appear nowhere in this backlog. The proposal's scope trade-off is to defer them before weakening security.
- Items marked **Proposed** in the [MVP Proposal §12](mvp-proposal.md#12-decisions-requested-from-technical-expert) are not approved. Their status is tracked in §6.

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
| Spike | `SPK-` | Time-boxed investigation report | Report merged to `docs/` with a recommendation and success criteria results. Spike code is not production code and does not pass the Definition of Done. |
| Decision | `DEC-` | Recorded decision | Approved by the mentor and recorded in an ADR |
| Mentor approval | `MA-` | Mentor approval point | Mentor's approval is recorded in the relevant issue |
| Risk | `RSK-` | Tracked risk | Reviewed at each sprint boundary |

**Specification work, implementation work and spikes are separate.** A spike may inform a specification, but it never replaces one. No implementation item may start until its specification is approved.

### 1.3 Priority, size and capacity basis

- **Priority:** every committed item is **Must** (required for the MVP). No Should or Could items are committed.
- **Size:** a planning estimate of implementation effort in **pair-days**. Sizes exclude SDD reviews, meetings and training, which are counted separately in §11.3. Pairs must validate these estimates in Sprint 1 and report actuals at each sprint boundary.
- **Capacity basis:** each sprint has 10 working days. A pair therefore has 10 pair-days per sprint and 60 across the internship. Four pairs give 240 gross pair-days. This assumes all four pairs are fully staffed for every working day (MA-11).

### 1.4 Sprint calendar

Non-working days are 24–27 December 2026 and 31 December 2026 – 3 January 2027. Each sprint below contains exactly 10 working days.

| Sprint | Dates (working days) | Proposal milestone (weeks) | Friday demos |
|--------|----------------------|----------------------------|--------------|
| 1 | 2–13 Nov 2026 (10) | Architecture and vertical slice (1–2) | 6 Nov, 13 Nov |
| 2 | 16–27 Nov 2026 (10) | Core platform (3–4) | 20 Nov, 27 Nov |
| 3 | 30 Nov – 11 Dec 2026 (10) | Submission pipeline (5–6) | 4 Dec, 11 Dec |
| 4 | 14–29 Dec 2026 (10; 24–27 Dec excluded) | Judge integration (7–8) | 18 Dec; the holiday week has no Friday demo. Proposed: 29 Dec (MA-07). |
| 5 | 30 Dec 2026 – 14 Jan 2027 (10; 31 Dec – 3 Jan excluded) | End-to-end MVP (9–10) | 8 Jan |
| 6 | 15–28 Jan 2027 (10) | Release hardening (11–12) | 15 Jan, 22 Jan |
| Final | 29 Jan 2027 (Friday) | — | **Internship final demo** |

Sprint 5 starts on 30 December, which is the only working day before the 4 January restart. Its effective start is 4 January, which is a risk (RSK-17). No demo falls on 25 December or 1 January, because both are non-working days.

---

## 2. MVP Features

These rows preserve the proposal's requirements and priorities. Implementation items are listed in §4.4 and sequenced in §10.

| ID | Feature | Priority | Expected behaviour (from proposal) | Specification item | Implementation items |
|----|---------|----------|------------------------------------|--------------------|----------------------|
| F01 | Accounts and access | Must | Register, sign in/out; learner and administrator roles. | SPEC-F01 | F01-01, F01-02, F01-03 |
| F02 | Problem catalogue | Must | List published problems with title, difficulty and solved state. | SPEC-F02 | F02-01 to F02-04 |
| F03 | Problem details | Must | Display statement, input/output contract, constraints and public examples. | SPEC-F03 | F03-01, F03-02 |
| F04 | Browser code editor | Must | Edit Go source code; load a starter template; preserve current work in the active session. | SPEC-F04 | F04-01, F04-02 |
| F05 | Run sample tests | Must | Execute code against public sample tests and return per-case feedback. | SPEC-F05 | F05-01, F05-02 |
| F06 | Submit solution | Must | Create asynchronous submission; return submission ID and processing state. | SPEC-F06 | F06-01 to F06-04 |
| F07 | Automated judge | Must | Compile and execute submitted Go code in isolation; evaluate hidden tests; produce verdict. | SPEC-F07 | F07-01 to F07-08 |
| F08 | Submission history | Must | Show own submissions, timestamps, states and final verdicts. | SPEC-F08 | F08-01, F08-02 |
| F09 | Problem administration | Must | Admin CRUD, draft/published status, sample and hidden test management. | SPEC-F09 | F09-01 to F09-03 |

---

## 3. Mandatory Work and Optional Complexity

This section separates what the MVP must contain from implementation choices that can be simplified without removing functionality.

### 3.1 Mandatory (not reduced in any option)

- **F01–F09 behaviour** as listed in §2.
- **Verdict model** (ACCEPTED, WRONG_ANSWER, COMPILATION_ERROR, RUNTIME_ERROR, TIME_LIMIT_EXCEEDED, MEMORY_LIMIT_EXCEEDED, INTERNAL_ERROR) and the separation of infrastructure failures from test failures.
- **Asynchronous lifecycle** QUEUED → RUNNING → COMPLETED or FAILED, with idempotent processing and no duplicate final results (proposal §6, §8).
- **Execution isolation controls** from proposal §7: dedicated Job per submission; CPU, memory and time limits; unprivileged; no network; restricted filesystem; bounded output; no access to secrets or Kubernetes API credentials; RBAC, namespace isolation, resource quotas and Pod security; Job TTL and cleanup; logs collected before cleanup; cancellation; hidden tests protected.
- **Failed learner solutions are not retried automatically** (proposal §7).
- **Adversarial tests and a dedicated security review** before public deployment (proposal §7, §11). Public exposure remains blocked until both pass.
- **Concurrency** of at least two independent submissions in staging (proposal §8, §11).
- **Reference fixtures** for one algorithmic and one HTTP-handler challenge (proposal §11).
- **Worker failure and timeout scenarios** tested (proposal §8).
- **Staging deployment** that is reproducible, CI green, and an operator README with known limitations (proposal §8).
- **Role-based access** and protection of hidden tests.

### 3.2 Optional complexity (may be simplified, not removed)

Each simplification below keeps the functionality listed in §3.1. The saving is an estimate and is used in §11.5.

| Item | Full option | Simplified option | Estimated saving |
|------|-------------|-------------------|------------------|
| F09-03 Admin UI | Full problem and test management screens | Forms covering the same CRUD, draft/published and test upload actions, with no extra polish | 1.5 pair-days |
| TST-02 End-to-end harness | Several journeys | One critical journey: browse, edit, run samples, submit, view verdict, view history | 1 pair-day |
| F05-02 Per-case feedback | Formatted diff output | Pass or fail per case, with bounded output text | 0.5 pair-days |
| F04-02 Session preservation | Designed mechanism | Simplest mechanism agreed in SPEC-F04 (DEC-13) | 0.5 pair-days |
| INF-07 Observability | Full metrics set | Structured logs, health endpoints and basic counters only (proposal wording) | 1 pair-day |
| INF-11 Staging | Full scripted environment | Single scripted environment with the minimum configuration needed for TST-03 | 1 pair-day |
| F07-07 Adversarial tests | Separate test harness | Reuse the SPK-01 negative-case harness, extended to production code paths. Coverage is not reduced. | 1 pair-day |

**Total simplification saving: about 6.5 pair-days.** These savings are offered to the mentor as options (MA-10). They are not applied in the plan in §10, so the plan still shows the unreduced load.

---

## 4. Backlog Items

### 4.1 Specifications and contracts

Specifications must be **Approved** before the implementation items they gate.

| ID | Item | Type | Size | Depends on | Deliverable | Gates |
|----|------|------|------|------------|-------------|-------|
| SPEC-F01 | Accounts and access specification | Specification | 1.5 | CON-01 (draft) | Approved spec in `specs/features/` | F01-01 |
| SPEC-F02 | Problem catalogue specification | Specification | 1.5 | CON-04 (draft) | Approved spec | F02-01 |
| SPEC-F03 | Problem details specification | Specification | 1.5 | CON-04 (final) | Approved spec | F03-01 |
| SPEC-F04 | Browser code editor specification | Specification | 1.5 | None | Approved spec; session mechanism recorded (DEC-13) | F04-01 |
| SPEC-F05 | Run sample tests specification | Specification | 1.5 | SPK-01, DEC-10 | Approved spec | F05-01 |
| SPEC-F06 | Submit solution specification | Specification | 1.5 | CON-01, CON-02 (draft) | Approved spec | F06-01 |
| SPEC-F07 | Automated judge specification | Specification | 4 | SPK-01, CON-01, CON-02 | Approved spec with isolation requirements from proposal §7 | F07-01, F07-02a |
| SPEC-F08 | Submission history specification | Specification | 1.5 | CON-01 | Approved spec | F08-01 |
| SPEC-F09 | Problem administration specification | Specification | 1.5 | CON-04, CON-03 (draft) | Approved spec | F09-01 |
| CON-01 | Submission lifecycle and verdict model | Contract | 1.5 | None | States, transitions, verdicts, infrastructure vs test-failure classification | SPEC-F06, SPEC-F07, F06-01, F07-05 |
| CON-02 | Messaging contract | Contract | 1.5 | CON-01, SPK-03 | Queue names, message schema, idempotency key, acknowledgement, retry and dead-letter behaviour | F06-02, F07-01 |
| CON-03 | OpenAPI contract for MVP endpoints | Contract | 4 | Approved specs for the endpoints in each extension | OpenAPI 3.x document. Initial part in Sprint 2, extended in Sprint 3. | INF-06, all API items |
| CON-04 | Problem packaging format and test visibility rules | Contract | 4 | DEC-09 (draft in Sprint 1, final in Sprint 2) | Packaging format, starter template format, visibility rules | SPEC-F03, F07-04, F09-02 |
| CON-05 | Data model | Contract | 4 | CON-01, CON-04 | Entities for users, roles, problems, test assets, submissions and results | INF-05, F01-02, F02-01 |

### 4.2 Infrastructure and technical work

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| INF-01 | Go module, repository layout, lint and format configuration | 1.5 | DEC-14 | Module, layout, lint configuration. Layout is proposed and subject to mentor approval. |
| INF-02 | CI baseline on GitHub Actions | 1.5 | INF-01 | Build, lint, unit tests, dependency and security checks |
| INF-03 | Local infrastructure | 4 | DEC-06 | Local Kubernetes cluster (k3d or Minikube), PostgreSQL and RabbitMQ for development. Setup steps documented in the README. |
| INF-04 | ADR process and initial records | 1 | None | ADR template and records for approved decisions |
| INF-05 | PostgreSQL schema, migrations and data access | 4 | CON-05, INF-03, DEC-11 | Migrations, pgx and sqlc setup |
| INF-06 | OpenAPI validation in CI | 1.5 | CON-03, INF-02 | Contract validation step in CI |
| INF-07 | Structured logging, health endpoints and basic metrics | 4 | INF-01 | Logging standard, health endpoints, basic metrics (see §3.2) |
| INF-08 | Frontend scaffold | 1.5 | INF-02 | React, TypeScript and Vite project with Monaco Editor dependency and CI build. Planned for Sprint 3 (see §9.4). |
| INF-09 | Auth0 development tenant configuration | 1.5 | DEC-03, MA-05 | Development tenant and application configuration. No secrets in the repository. |
| INF-10 | Judge namespace, RBAC, resource quotas and namespace network policy | 4 | DEC-08, MA-03, INF-03 | Cluster-level controls for the judge namespace, validated. Job-level security settings are in F07-02a. |
| INF-11 | Staging environment and deployment | 7 | DEC-07, MA-04, INF-02, INF-05 | Repeatable staging deployment (see §3.2) |
| INF-12a | Container images for API and orchestrator | 1.5 | INF-02 | Images with dependency and image security checks in CI |
| INF-12b | Container image for the judge runner | 1.5 | INF-12a, F07-02b | Image with image security checks in CI |
| INF-13 | Operator README and runbook | 1.5 | INF-11 | Operator documentation and known-limitations list |

INF-14 from the first draft is removed. Its content, the release checklist, is now part of TST-04.

### 4.3 Spikes

Spikes are time-boxed. Their output is a report with a recommendation, not production code. Spikes must not run learner-submitted code outside the isolated spike environment. Timeboxes and success criteria are in §5.

| ID | Item | Timebox | Depends on | Deliverable |
|----|------|---------|------------|-------------|
| SPK-01 | Secure Go execution feasibility with `go test` and Kubernetes Jobs | 8 | INF-03 (interim cluster allowed), DEC-01, DEC-02 (provisional) | Report with a go, conditional go or no-go outcome against SC-01 to SC-10 (§5.1) |
| SPK-02 | Auth0 OIDC integration check | 2 | DEC-03 (provisional) | Report against SC-A to SC-C (§5.2) |
| SPK-03 | RabbitMQ and Watermill check | 4 | DEC-04 (provisional), INF-03 | Report against SC-A to SC-D (§5.3) |
| SPK-04 | Hosting comparison (GKE and VM-hosted k3s) | 3 | None (feeds DEC-07) | Report against SC-A to SC-C (§5.4) |
| SPK-05 | Isolation mechanism comparison | 3 | SPK-01 | Report against SC-A and SC-B (§5.5) |

### 4.4 Feature implementation items

Each item requires its specification to be **Approved** first. Items are grouped by feature.

**F01 Accounts and access**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F01-01 | Sign-in, sign-out and token validation in the Go API | 4 | SPEC-F01, INF-09 | Authentication middleware; tests with fixture tokens |
| F01-02 | Learner and administrator roles and endpoint authorisation | 4 | F01-01, CON-05 | Role model; authorisation checks; tests for role boundaries |
| F01-03 | Sign-in and sign-out UI | 1.5 | F01-01, INF-08 | React sign-in and sign-out flow |

**F02 Problem catalogue**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F02-01 | Problem persistence and published-status filter | 4 | SPEC-F02, CON-05, INF-05 | Problem storage; published-only query |
| F02-02 | Catalogue list endpoint (title, difficulty) | 3 | F02-01 | List endpoint with title and difficulty. Solved state is returned as "unknown" until F02-04. |
| F02-03 | Catalogue UI | 1.5 | F02-02 | React catalogue list |
| F02-04 | Solved state in the catalogue | 1.5 | F02-02, F08-01 | Solved state populated from the learner's submissions |

**F03 Problem details**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F03-01 | Problem details endpoint | 4 | SPEC-F03, F02-01, CON-04 | Statement, input/output contract, constraints and public examples |
| F03-02 | Problem details UI | 1.5 | F03-01 | React problem page |

**F04 Browser code editor**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F04-01 | Monaco Editor integration with Go mode and starter template | 4 | SPEC-F04, INF-08, CON-04; F03-02 is a soft dependency (a stub page can host the editor) | Editor component; starter template loaded from problem data |
| F04-02 | Active-session preservation of current work | 1.5 | F04-01, DEC-13 | Work preserved in the active session, using the mechanism in DEC-13 |

**F05 Run sample tests**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F05-01 | Run-samples endpoint through the isolated execution path | 7 | SPEC-F05, DEC-10, F07-02b, F07-03, F07-05 | Endpoint returning per-case feedback for public sample tests |
| F05-02 | Per-case feedback in the editor | 1.5 | F05-01, F04-01 | UI for per-case results (see §3.2) |

**F06 Submit solution**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F06-01 | Submission creation endpoint | 4 | SPEC-F06, F01-02, CON-01, CON-05 | Persisted submission; returns submission ID and QUEUED state |
| F06-02 | Dispatch to RabbitMQ through Watermill | 4 | F06-01, CON-02, SPK-03, DEC-12 | Publisher using the agreed message contract |
| F06-03 | Submission state endpoint | 1.5 | F06-01 | Endpoint returning the current state |
| F06-04 | Submit action and state polling in the editor | 1.5 | F06-03, F04-01 | UI showing processing state |

**F07 Automated judge**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F07-01 | Orchestrator consumer with state transitions and controlled mock executor | 4 | SPEC-F07, F06-02, CON-01, CON-02 | Watermill consumer; QUEUED → RUNNING → COMPLETED or FAILED; mock executor for pipeline tests |
| F07-02a | Job specification builder with mandated security settings | 4 | SPEC-F07, DEC-08 | Pure Go builder for one Job per submission: CPU, memory and time limits; unprivileged; read-only root filesystem; no network; no service-account token; no secrets. Unit-tested without a cluster. |
| F07-02b | Job lifecycle through client-go | 4 | F07-02a, INF-10, SPK-01, MA-03 | Create, watch and cancel Jobs in the cluster |
| F07-03 | Job monitoring, log and report collection, TTL cleanup | 4 | F07-02b | Monitoring; logs and reports collected before cleanup; TTL cleanup |
| F07-04 | Hidden test asset handling | 4 | CON-04, INF-05, F07-02b | Test suites available only to the judge Job; never returned to learners |
| F07-05 | Verdict mapping and failure classification | 3 | CON-01, F07-03 | Verdicts from CON-01; infrastructure failures separated from test failures |
| F07-06 | Idempotent processing and retry policy | 4 | F07-01, F07-05 | Redelivery does not create duplicate final results; learner failures are not retried automatically |
| F07-07 | Adversarial security tests for the judge | 6 | F07-02b, F07-03, F07-04, SPK-01 harness | Tests for network access, filesystem writes, secret access, resource exhaustion, output flood and hidden test reads. Gate for any public exposure. |
| F07-08 | Worker failure and timeout integration tests | 3 | F07-06, F07-03 | Tests for worker restart mid-run, Job failure and timeout |

F07-02 from the first draft (size 8) is split into F07-02a and F07-02b for granularity. The split does not reduce effort by itself.

**F08 Submission history**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F08-01 | Submission history endpoint | 3 | SPEC-F08, F06-01, CON-01, F01-02 | Own submissions only, with timestamps and states. Verdict field is populated once F07-05 is in place, so F08-01 no longer waits for it. |
| F08-02 | Submission history UI | 1.5 | F08-01, F04-01 | React history page |

**F09 Problem administration**

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| F09-01 | Admin CRUD endpoints with draft and published status | 4 | SPEC-F09, F01-02, F02-01, CON-03 | Endpoints restricted to administrators |
| F09-02 | Sample and hidden test management | 7 | F09-01, F07-04, CON-04 | Upload, store and replace test assets |
| F09-03 | Admin UI | 4 | F09-01, F09-02 | Problem and test management screens (see §3.2) |

### 4.5 Test and hardening items

| ID | Item | Size | Depends on | Deliverable |
|----|------|------|------------|-------------|
| TST-01 | Reference challenge fixtures | 4 | CON-04 | One algorithmic and one HTTP-handler challenge, with reference solutions and negative cases |
| TST-02 | End-to-end test harness | 3 | F05-02, F08-02 | Browser-driven critical journey (see §3.2) |
| TST-03 | Concurrency test on staging | 3 | INF-11, F07-06, DEC-05 | At least two concurrent submissions through Kubernetes Jobs; correct results; no duplicate final evaluations |
| TST-04 | Release acceptance run and release checklist | 4 | TST-02, TST-03, HRD-01, HRD-02 | Evidence for each release acceptance criterion; checklist and demo script (absorbs INF-14) |
| HRD-01 | Security review and fixes for the judge | 4 | F07-07, INF-10 | Dedicated security review of the isolation mechanism; findings fixed or accepted by the mentor (MA-09) |
| HRD-02 | Defect fixes for critical and high severity | 3 | TST-02, TST-03 | No unresolved critical or high defects |

---

## 5. Technical Spikes: Timeboxes and Success Criteria

Each spike has a fixed timebox. When the timebox ends, the pair writes the report with the outcome that has been reached. A spike is not extended without mentor approval (MA-12). Each criterion is marked **Verified**, **Partially verified** or **Not verified**. A criterion that is not verified is listed as an open risk. It is never assumed.

### 5.1 SPK-01: Secure Go execution with `go test` and Kubernetes Jobs (8 pair-days)

Runs only trusted test packages. No learner-submitted code runs in the spike.

| SC | Criterion | Measure (pass threshold) |
|----|-----------|--------------------------|
| SC-01 | A client-go program creates a Job that runs a trusted `go test` package, and the JSON result is retrieved | 10 of 10 runs |
| SC-02 | Time limit: a Job that exceeds its deadline is terminated | Terminated within deadline + 5 s, 10 of 10 runs |
| SC-03 | Memory limit: a memory-exhaustion test program is killed by the limit | 5 of 5 runs, with termination reason captured |
| SC-04 | No network: an outbound TCP attempt from the Job fails | 5 of 5 runs |
| SC-05 | Restricted filesystem: writes outside the allowed path fail, and the root filesystem is read-only | 5 of 5 runs |
| SC-06 | Unprivileged: the container runs as non-root, with no privilege escalation and no service-account token mounted | Verified from the Job spec and by attempt, 5 of 5 |
| SC-07 | Output bounded: an output flood is truncated at the configured limit | 3 of 3 runs |
| SC-08 | Logs and reports are retrieved before cleanup | 10 of 10 runs |
| SC-09 | TTL cleanup removes finished Jobs | Removed within TTL + 60 s, 10 of 10 runs |
| SC-10 | Cancellation terminates a running Job | Terminated within 30 s, 5 of 5 runs |

Startup and run latency are recorded for information (median and range). They are not a pass criterion.

**Outcome rule:** *Go* requires SC-01 to SC-10 verified. *Conditional go* requires SC-01 and SC-08 verified, with every other unverified criterion listed as a named risk for SPEC-F07 and MA-03. Otherwise *no-go*.

### 5.2 SPK-02: Auth0 OIDC integration (2 pair-days)

| SC | Criterion | Measure |
|----|-----------|---------|
| SC-A | Sign-in works through the Auth0 development tenant, or through a documented fixture if tenant access is not yet available (RSK-06) | Demonstrated once |
| SC-B | The Go middleware checks signature, issuer, audience and expiry | Rejects a tampered token, an expired token and a wrong-audience token, 3 of 3 |
| SC-C | The role claim maps to learner or administrator | Correct for two test users |

### 5.3 SPK-03: RabbitMQ and Watermill (4 pair-days)

| SC | Criterion | Measure |
|----|-----------|---------|
| SC-A | Publish and consume with no message loss | 100 of 100 messages |
| SC-B | A consumer killed during processing causes redelivery | 5 of 5 trials |
| SC-C | A duplicate delivery of the same message leads to one final state, using an idempotency key | 50 of 50 duplicates |
| SC-D | A poison message reaches the dead-letter route after the configured attempts | 3 of 3 |

### 5.4 SPK-04: Hosting comparison, GKE and VM-hosted k3s (3 pair-days)

| SC | Criterion | Measure |
|----|-----------|---------|
| SC-A | Comparison table covering setup steps, operations burden, fit with the judge isolation controls, learning curve and budget | Complete table with a source for each entry |
| SC-B | Cost estimate from published pricing, with the pricing source and date recorded | Estimate for the staging footprint |
| SC-C | No cloud resources are created | Nothing provisioned without MA-04 |

### 5.5 SPK-05: Isolation mechanism comparison (3 pair-days)

| SC | Criterion | Measure |
|----|-----------|---------|
| SC-A | Each isolation requirement from proposal §7 is mapped to each candidate, with its evidence source (SPK-01 result or cited documentation) | Complete matrix; no empty cells |
| SC-B | Recommendation for DEC-08, with rejected options and reasons | Written recommendation. Nothing is adopted without MA-03. |

---

## 6. Decisions

Status values are: **Proposed** (from the proposal, not approved), **Provisional** (adopted only for spike work, not approved), **Open** (not yet decided), **Assumption** (to be confirmed). **Approved: none.** No decision in this backlog is approved.

### 6.1 Needed on Sprint 1 day 1 (to start work)

These are needed to start work on day 1. They are reversible and low-cost, so provisional choices are adopted for Sprint 1 only.

| ID | Decision | Status | Provisional choice for Sprint 1 | Why day 1 |
|----|----------|--------|---------------------------------|-----------|
| DEC-14 | Go version pin | Open | Latest stable release at sprint start, pinned in an ADR | INF-01 needs a version |
| DEC-06 | Local infrastructure: k3d or Minikube | Proposed | One of the two, chosen by the team on day 1 | INF-03 and SPK-01 need a cluster |

### 6.2 Hypotheses for the Sprint 1 spikes (not approved on day 1)

The spikes test these decisions. They stay **Provisional** until the spike report is accepted. They are not approved on day 1.

| ID | Decision | Status | Approved when | Evidence |
|----|----------|--------|---------------|----------|
| DEC-01 | Submission evaluation through the Go testing framework | Provisional | After SPK-01 report | SPK-01 |
| DEC-02 | Execution orchestration with Kubernetes Jobs and client-go | Provisional | After SPK-01 report | SPK-01 |
| DEC-03 | Authentication with Auth0 | Provisional | After SPK-02 report and MA-05 | SPK-02 |
| DEC-04 | Messaging with RabbitMQ and Watermill | Provisional | After SPK-03 report | SPK-03 |

### 6.3 Need spike evidence or later input

| ID | Decision | Status | Needed by | Input required |
|----|----------|--------|-----------|----------------|
| DEC-05 | Concurrency target of at least two concurrent submissions | Proposed | Sprint 6 | TST-03 evidence in staging |
| DEC-07 | Hosting: GKE or VM-hosted k3s | Open | Start of Sprint 4 (14 Dec 2026) | SPK-04 report, MA-04 |
| DEC-08 | Execution isolation mechanism | Open | Start of Sprint 4 (14 Dec 2026, MA-03) | SPK-01 and SPK-05 reports |
| DEC-09 | Problem packaging and test visibility rules | Open | Draft in Sprint 1, final in Sprint 2 | Product input from the mentor; the proposal says "finalised during architecture design" and it is not in its §12 table |
| DEC-10 | Run-samples uses the same isolated execution path as submit | Assumption | Sprint 4 | SPK-01 result and SPEC-F05. Must be confirmed; security-sensitive. |
| DEC-11 | Database migration tool | Open | Sprint 2 | Engineering ADR. No spike needed. Not in the proposal. |
| DEC-12 | Dispatch consistency between database write and message publish | Open | **Start of Sprint 3** (30 Nov 2026; F06-02 depends on it) | SPK-03 report. Reliability requirement in proposal §6. Not in the proposal's §12 table. |
| DEC-13 | Editor session-preservation mechanism | Open | Sprint 3 (in SPEC-F04) | Product input, recorded in the spec. Not in the proposal's §12 table. |

### 6.4 Notes

- **DEC-12 moves earlier.** The first draft placed it in Sprint 3 without a day. F06-02 starts on the first day of Sprint 3, so DEC-12 must be decided by then.
- **DEC-14 and DEC-06** are the only decisions required on day 1 for work to start.

---

## 7. Dependency Map and Critical Path

### 7.1 Critical path

The critical path runs through the judge. It has **no schedule slack**. Any slip in Sprint 4 pushes the end-to-end work and the release.

```text
INF-03 → SPK-01 (8) → SPK-05 (3) → DEC-08 (MA-03, start of Sprint 4)
      → F07-02a (4) → F07-02b (4) → F07-03 (4) → F07-05 (3)
      → F05-01 (7) → F05-02 (1.5) → TST-02 (3) → TST-04 (4) → release
                  ↘ F07-07 (6) → HRD-01 (4) → TST-04
```

Pair-day figures are the sizes in §4. The chain runs from the first day of Sprint 1 to the last sprint.

### 7.2 Second long chain: staging

```text
DEC-07 (start of Sprint 4) → INF-11 (7) → TST-03 (3) → TST-04 → release
```

### 7.3 Other dependencies

```text
DEC-03 (MA-05) → INF-09 → F01-01 → F01-02 → F06-01 → F06-02 → F07-01 → F07-06 → F07-08
SPK-03 → DEC-12 (start of Sprint 3) → F06-02
CON-01 → CON-05 → INF-05 → F02-01 → F03-01 → F06-01
F07-05 → F05-01 (verdicts) ; F08-01 → F02-04 (solved state, not on the critical path)
```

---

## 8. Parallel Work and Integration Risks

These streams can run in parallel without integration risk, provided the listed integration points are respected.

| Stream | Runs in parallel with | Integration point | Control |
|--------|-----------------------|-------------------|---------|
| Platform (INF-01, 02, 03, 04; INF-08 in Sprint 3) | Specs (SPEC-F01, F02) and spikes | Local cluster used by SPK-01 and SPK-03 | Pair A owns INF-03. Spikes may use an interim cluster until then. |
| Backend core (F01-01, F02-01, F06-01) | Frontend scaffold (INF-08) and UI against mocked CON-03 responses | Endpoints must match the OpenAPI document | UI wiring starts only after the endpoint's CON-03 section is approved. INF-06 checks contract conformance in CI. |
| Judge: consumer (F07-01) | Job builder (F07-02a) and hidden test handling (F07-04) | Job creation at F07-02b | CON-02 frozen before F07-01 starts. One pair owns the Job specification. |
| Admin (F09-01) | Judge chain in Sprint 4 | Shared database migrations (INF-05) | Migrations are serialised, with one owner per sprint. CON-05 frozen before F09-01 starts. |
| Learner UI (F02-03, F03-02, F04-01) | Backend endpoints in Sprints 3 and 4 | Shared OpenAPI document | Same rule as backend core. |
| Staging (INF-11) | Judge development | Integration at TST-03 | Staging is built against the INF-05 schema only, not against judge internals. |

**Shared-file risks** apply to the migrations folder, the OpenAPI file, the Job specification and the HTTP router. Each shared file has one owner per sprint. Changes to a shared file are merged in order, and the owner reviews every change to it.

---

## 9. Sprint 1 Backlog (2–13 November 2026)

Sprint 1 has four parallel workstreams, one per pair. Pair labels A–D are placeholders. The mentor assigns pairs at the start of the sprint.

Day one is a decision session (MA-13): the team confirms the provisional choices in §6.1 and the spike hypotheses in §6.2. Nothing in §6.2 is approved on day one.

### 9.1 Committed and stretch

| Workstream | Pair | Committed | Stretch | Estimated load |
|------------|------|-----------|---------|----------------|
| WS1 Platform foundations | Pair A | INF-01, INF-02, INF-03, INF-04 | None | 8 |
| WS2 Secure execution spike | Pair B | SPK-01 (timebox 8) | Start SPK-05 reading list on day 10 | 8 |
| WS3 Identity and messaging spikes | Pair C | SPK-02 (timebox 2), SPK-03 (timebox 4) | Review SPEC-F01 for another pair; CON-02 (1.5, mentor-approved 9 Nov) | 6 |
| WS4 Specifications and contracts | Pair D | CON-01, CON-04 (draft), SPEC-F01 (in review), SPEC-F02 (draft) | SPEC-F04 (1.5, mentor-approved 9 Nov) | 6.5 |

Committed implementation load is **28.5 pair-days**. Realistic Sprint 1 net capacity is **24.5 pair-days** (§11.3), so the earlier committed set was **4 pair-days over**. The approved Sprint 1 Plan ([sprint-01.md](../planning/sprint-01.md) §2 and §4.1, MA-15) commits 24.5 pair-days and is authoritative. The 28.5 figure is a historical backlog estimate and is not reopened here. The options in §9.4 are not needed for Sprint 1. The gap is not hidden in a stretch item. The two mentor-approved stretch items, CON-02 (Pair C) and SPEC-F04 (Pair D), total 3.0 pair-days. They are recorded in [sprint-01.md §4.2](../planning/sprint-01.md#42-stretch-mentor-approved-9-november-2026) and are not counted in the committed load above. Their rows in §10 (Sprints 2 and 3) are unchanged until Sprint 2 planning.

### 9.2 Workstream detail

**WS1 Platform foundations (Pair A).** Pair A's first task is INF-03, because Pair B and Pair C need a cluster early in the sprint. INF-08 has moved to Sprint 3 (§9.4).

| ID | Objective | Acceptance criteria |
|----|-----------|---------------------|
| INF-03 | Local infrastructure | Local cluster, PostgreSQL and RabbitMQ run; setup steps in the README; ready by day 3 |
| INF-01 | Module and layout | `go.mod` present, with the Go version from DEC-14; lint and format configured; layout reviewed |
| INF-02 | CI baseline | Pipeline runs on push; build, lint and tests run; green on `main` |
| INF-04 | ADR process | Template merged; ADRs for the provisional choices in §6.1 |

**WS2 Secure execution spike (Pair B).** Uses the cluster from INF-03 (or an interim cluster on the same version). Timebox starts on day 3 and ends on day 10, which is 8 pair-days. Results are reported against SC-01 to SC-10 in §5.1 at the sprint demo.

**WS3 Identity and messaging spikes (Pair C).** SPK-02 uses an Auth0 tenant if MA-05 is in place, otherwise a documented fixture (RSK-06). SPK-03 uses the local RabbitMQ from INF-03.

**WS4 Specifications and contracts (Pair D).** CON-01 is agreed with Pairs B and C before it is marked draft-complete. CON-04 is a draft, and DEC-09 options are listed for the mentor. SPEC-F01 and SPEC-F02 are submitted for cross-pair review.

### 9.3 Sprint 1 exit criteria

- The spike reports for SPK-01, SPK-02 and SPK-03 are merged to `docs/` with their outcomes against the success criteria in §5.
- The CI baseline (INF-02) is green on `main`.
- SPEC-F01 and SPEC-F02 are in cross-pair review, with the reviewing pair different from the author's pair.
- CON-01 is agreed by all dependent pairs.
- DEC-14 and DEC-06 are recorded. The §6.2 hypotheses are recorded as Provisional, with the evidence that will decide them.
- The 13 November demo presents the spike outcomes and the CI baseline.

### 9.4 Sprint 1 deferred work and options

Deferred from Sprint 1:

- INF-08 (frontend scaffold, 1.5) moves to Sprint 3. Its first consumer is F01-03 in Sprint 4.
- SPK-05 starts in Sprint 2 (3 pair-days). It needs the SPK-01 report.
- SPEC-F02 approval moves to the start of Sprint 2 if its review is not complete by 13 November.

Options if MA-15 does not accept the 4 pair-day over-commit:

- **Option 1.** Move the CON-04 draft (2) and SPEC-F02 draft (1.5) to Sprint 2. Sprint 1 falls to 25, which is still 0.5 over. Sprint 2 rises by 3.5.
- **Option 2.** Keep the committed set and accept Sprint 1 running at about 116% of realistic net. Recovery depends on Sprint 2 absorbing the carry-over.

Cross-pair reviews consume capacity from all four pairs. This is a known risk (RSK-05).

---

## 10. Six-Sprint Plan with Load

Each sprint's planned load is the sum of the items allocated to it. Gross capacity per sprint is **40 pair-days** (4 pairs × 10 working days). Realistic net per sprint is taken from §11.3. Loads above net are flagged. Loads above gross are not feasible.

| Sprint | Dates | Items (planned load) | Load | Realistic net | Over (+) / spare (−) | Notes |
|--------|-------|----------------------|------|---------------|----------------------|-------|
| **1** | 2–13 Nov | INF-01, INF-02, INF-03, INF-04; SPK-01, SPK-02, SPK-03; CON-01; CON-04 (draft, 2); SPEC-F01, SPEC-F02 (draft) | **28.5** | 24.5 | **+4** | Approved Sprint 1 Plan is 24.5 (MA-15, sprint-01.md §2); 28.5 is the earlier estimate (§9.1). Stretch CON-02 and SPEC-F04 (3.0) are outside this load (§9.1). |
| **2** | 16–27 Nov | CON-04 (final, 2); CON-02 (1.5); CON-03 (initial, 2); CON-05 (4); INF-05 (4); INF-07 (4); INF-09 (1.5); INF-12a (1.5); F01-01 (4); F02-01 (4); SPEC-F03, SPEC-F06 (1.5 each); SPEC-F07 (4); SPK-05 (3) | **38.5** | 25 | **+13.5** | SPEC-F01 and SPEC-F02 approved at the start. SPEC-F07 approval at the start of Sprint 3. |
| **3** | 30 Nov – 11 Dec | F01-02 (4); F03-01 (4); F06-01 (4); F06-02 (4); F07-01 (4); SPK-04 (3); CON-03 extension (2); SPEC-F04, SPEC-F05, SPEC-F08, SPEC-F09 (1.5 each); TST-01 (4); INF-08 (1.5) | **36.5** | 25 | **+11.5** | DEC-12 decided on the first day. SPEC-F07 approved on the first day, gating F07-02a. |
| **4** | 14–29 Dec | F04-01 (4); F07-02a (4); F07-02b (4); F07-03 (4); F07-04 (4); F07-05 (3); F07-06 (4); F08-01 (3); F09-01 (4); INF-10 (4) | **38** | 27.5 | **+10.5** | DEC-07 and DEC-08 approved at the start (14 Dec), gating INF-11 and F07-02b. MA-03 and MA-04 recorded. Holiday break mid-sprint. |
| **5** | 30 Dec – 14 Jan | F05-01 (7); F07-07 (6); INF-11 (7); F02-02 (3); F02-03, F02-04, F03-02, F04-02, F06-03, F06-04, F08-02, F01-03, F05-02, INF-06, INF-12b (1.5 each, 11 items) | **39.5** | 28.5 | **+11** | Effective start 4 Jan (RSK-17). Critical path. Near gross. |
| **6** | 15–28 Jan | F07-08 (3); F09-02 (7); F09-03 (4); TST-02 (3); TST-03 (3); TST-04 (4); HRD-01 (4); HRD-02 (3); INF-13 (1.5) | **32.5** | 33 | **−0.5** | Demos 15 and 22 Jan. Final internship demo 29 Jan. |
| **Total** | 2 Nov – 29 Jan | | **213.5** | **163.5** | **+50** | About 31% over realistic net (§11.4) |

**Reading this plan.** Gross capacity is not the binding constraint in most sprints: every sprint's load is under 40. Realistic net is the constraint. The shortfall builds in Sprints 1–5 and reaches about 50.5 pair-days by the end of Sprint 5. Sprint 6 absorbs 0.5 of it and does not recover the rest. Work that cannot be moved earlier, because of dependencies, is what slips. §11.4 lists it.

---

## 11. Capacity Reassessment

### 11.1 Backlog estimate

| Category | Items | Pair-days (revision 3) | Revision 2 | First draft |
|----------|-------|------------------------|------------|-------------|
| Specifications and contracts | SPEC-F01 to SPEC-F09, CON-01 to CON-05 | 31 | 31 | 29 |
| Infrastructure | INF-01 to INF-13 (INF-14 removed) | 36 | 36 | 44 |
| Spikes | SPK-01 to SPK-05 | 20 | 20 | 22 |
| Features | F01-01 to F09-03 | 105.5 | 105.5 | 116 |
| Test and hardening | TST-01 to TST-04, HRD-01 to HRD-02 | 21 | 21 | 22 |
| **Total** | | **213.5** | **213.5** | **233** |

The backlog itself is unchanged by this revision. What changes is the capacity it is measured against. The first draft's 233 is superseded: it mixed sizes for the same kinds of item, and it counted INF-14 separately from TST-04. The figure of 213.5 is the corrected estimate.

### 11.2 Gross capacity under the confirmed calendar

- Six sprints × 10 working days = **60 working days**.
- Four pairs × 60 working days = **240 gross pair-days**.
- Holidays (24–27 December and 31 December – 3 January) are already excluded from the 60 working days.
- The earlier plan assumed 40 gross pair-days in every sprint. Under the confirmed calendar, each sprint has 10 working days, so 40 per sprint holds. The earlier flat 5% holiday contingency is no longer needed for the calendar itself. Leave, sickness and staffing changes are not in the calendar, and are covered by MA-11.

### 11.3 Realistic net capacity

Net capacity is gross capacity less the work that a pair must do on the backlog's behalf but that is not in the backlog's sizes. The categories below are planning estimates.

| Overhead category | S1 | S2 | S3 | S4 | S5 | S6 | Total | Basis |
|-------------------|----|----|----|----|----|----|-------|-------|
| Specification and contract review (author and reviewer) | 3 | 6 | 5 | 1 | 1 | 0 | 16 | About 1 pair-day per spec or contract reviewed, plus re-review |
| Pull request review | 1.5 | 2 | 2.5 | 3 | 2.5 | 1 | 12.5 | About 50 pull requests at about 0.25 pair-day each |
| Training and onboarding (beyond spike learning) | 6 | 1 | 0.5 | 0.5 | 0 | 0 | 8 | Go, Kubernetes, RabbitMQ and Auth0 onboarding for four pairs |
| Meetings and ceremonies | 4 | 4 | 4 | 4 | 4 | 4 | 24 | Planning, review, retrospective, demo preparation and mentor sessions, about 1 pair-day per pair per sprint |
| Integration, CI and environment support | 1 | 2 | 3 | 4 | 4 | 2 | 16 | About 10% of implementation load, weighted to the judge and staging |
| **Total overhead** | **15.5** | **15** | **15** | **12.5** | **11.5** | **7** | **76.5** | |
| Gross capacity | 40 | 40 | 40 | 40 | 40 | 40 | 240 | §11.2 |
| **Realistic net** | **24.5** | **25** | **25** | **27.5** | **28.5** | **33** | **163.5** | |

**Sensitivity range for realistic net:**

- **Best case, about 180.** Meetings reduced to three per sprint (saves 6), pull request review at 0.15 pair-day each (saves 5), and integration at 6% (saves 6). Overhead falls to about 59.5.
- **Worst case, about 140.** Pull request review at 0.4 each (20), integration at 15% (24), training 12, specification review 20, meetings 24. Overhead rises to 100.

These overhead figures are unvalidated (RSK-16). They are the largest uncertainty in this assessment.

### 11.4 Fit test

The backlog, at 213.5, is compared with realistic net, 163.5.

| Sprint | Load | Realistic net | Over (+) / spare (−) | Cumulative shortfall |
|--------|------|---------------|----------------------|----------------------|
| 1 | 28.5 | 24.5 | +4 | 4 |
| 2 | 38.5 | 25 | +13.5 | 17.5 |
| 3 | 36.5 | 25 | +11.5 | 29 |
| 4 | 38 | 27.5 | +10.5 | 39.5 |
| 5 | 39.5 | 28.5 | +11 | 50.5 |
| 6 | 32.5 | 33 | −0.5 | **50** |

**Result: the backlog does not fit.** Under the base case, it is about **50 pair-days (31%) over** realistic net capacity. Under the best case of about 180 it is still about 33 over (19%). Under the worst case of about 140 it is about 74 over (52%).

**Work that would not finish by 29 January at base-case capacity.** Taking the last items in dependency order, the release-critical tail adds up to about 51 pair-days:

| Item | Size | Why it is at risk |
|------|------|-------------------|
| F05-01 Run-samples endpoint | 7 | On the critical path |
| F07-07 Adversarial security tests | 6 | Gate for public exposure and for HRD-01 |
| INF-11 Staging environment | 7 | Needed for TST-03 and the release |
| F09-02 Sample and hidden test management | 7 | Needed for the admin flow in the release criteria |
| F09-03 Admin UI | 4 | Needed for the admin flow in the release criteria |
| F07-08 Worker failure and timeout tests | 3 | Release criterion (proposal §11) |
| TST-02, TST-03 | 3 each | Release evidence |
| TST-04 Release acceptance run | 4 | Release evidence |
| HRD-01 Security review and fixes | 4 | Required before public deployment |
| HRD-02 Defect fixes | 3 | No unresolved critical or high defects |

These items are mandatory under the proposal. None may be dropped to make the estimate fit (§3.1).

### 11.5 Options for the mentor

No option below reduces mandatory functionality or security controls, except option D, which is listed only to show what would be required.

| Option | Effect on the 50 pair-day shortfall | Trade-off |
|--------|-------------------------------------|-----------|
| **A.** Apply the §3.2 simplifications | Reduces it by 6.5 to about 43.5 | Admin UI, test harness and observability are minimal. Functionality is unchanged. |
| **B.** Reduce meetings from four to three per sprint | Reduces it by 6 (about 37.5 after A) | A team decision. Must not remove the Friday demo or the sprint review. |
| **C.** Add capacity. Most efficiently, a support engineer or mentor time on INF-11 and judge integration | About 37.5 pair-days after A and B. Equal to about one pair for 38 of the 60 working days. | Needs a decision on who and when (MA-16). Support must be within team rules. |
| **D.** Reduce mandatory scope | Not applied. Would reduce the shortfall by the size of the removed items. | Contrary to this revision's instruction to preserve requirements. Requires explicit approval outside this backlog. |
| **E.** Release internal-only, if judge isolation is not validated by the start of Sprint 4 (proposal §7) | Does not reduce effort. It changes the release target. | Requires MA-09 and a documented limitation. Mandatory security controls still apply. |
| **F.** Change the internship end date | Not applied. Outside this backlog. | The confirmed end date is 29 January 2027. |

**Combined path that fits, on paper:** A plus B plus C closes the base-case gap, with C at about 37.5 pair-days. A plus B plus C at the worst case would still be short. Because the best and worst cases differ by about 40 pair-days, the mentor should not rely on any single estimate.

**Security items (F07-07, HRD-01, INF-10, and the judge's Job-level controls in F07-02a) must not be reduced under any option.**

---

## 12. Risks

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|------------|--------|------------|
| RSK-01 | Judge isolation takes longer than planned. The critical path has no slack. | High | High | SPK-01 and SPK-05 in Sprints 1–2; the mock executor lets the pipeline progress; option E in §11.5 |
| RSK-02 | Backlog exceeds realistic net capacity by about 50 pair-days (31%) | High | High | MA-08 and MA-16; options in §11.5; security items protected |
| RSK-03 | Leave, sickness or pair changes reduce the 60 working days, most of all in the holiday-adjacent Sprints 4 and 5 | High | High | MA-11 confirms staffing. Absence is tracked at each sprint boundary. |
| RSK-04 | Pair rotation every two weeks disrupts features mid-way | Medium | High | Handover rule in the SDD Workflow §3; items sized to fit within one sprint where possible |
| RSK-05 | Cross-pair specification reviews become a bottleneck | Medium | Medium | Review assignments made at sprint start; reviews counted in §11.3 |
| RSK-06 | Auth0 tenant access is delayed | Medium | High | MA-05 in Sprint 1; SPK-02 uses a documented fixture until access arrives |
| RSK-07 | Hidden tests leak through logs, outputs or API responses | Medium | High | F07-04 and F07-07 test this; output bounded; logs reviewed before collection |
| RSK-08 | Interns are new to Kubernetes, RabbitMQ and client-go | High | Medium | Spikes in Sprint 1; pairing; mentor review of spike reports |
| RSK-09 | Hosting approval or budget delays staging (INF-11, TST-03) | Medium | Medium | DEC-07 at the start of Sprint 4 (14 Dec); option E as fallback |
| RSK-10 | Judge design depends on DEC-08, which is open | Medium | High | Decision needed at the start of Sprint 4 (14 Dec). SPK-05 prepares options during Sprint 2. |
| RSK-11 | No Friday demo falls on 25 December or 1 January, and the holiday week has no ordinary demo | Certain | Low | Proposed demo on 29 December (MA-07). Confirm before Sprint 4. |
| RSK-12 | Mentor approval latency blocks the gates in §6 and §10 | Medium | High | Approval points listed in §13; approvals requested at sprint boundaries |
| RSK-13 | Integration conflicts on shared files (migrations, OpenAPI, Job specification, router) | Medium | Medium | One owner per shared file per sprint; §8 rules |
| RSK-14 | Spike success criteria are not met, leaving risks unresolved for Sprint 2 | Medium | High | Outcome rule in §5.1 defines the conditional path; MA-12 decides on extension |
| RSK-15 | Sprint 1 committed load (28.5) exceeds realistic net (24.5) | High | Medium | MA-15 and the options in §9.4 |
| RSK-16 | Overhead estimates in §11.3 are wrong. The range is 140–180. | High | High | Validate in Sprint 1 with timesheet-style actuals; re-plan at each sprint boundary |
| RSK-17 | Sprint 5 begins with one working day (30 Dec), then a 4-day break. The critical path's staging and run-samples work starts with little margin. | High | High | Start INF-11 planning in Sprint 4 (DEC-07 by 14 Dec); protect the 4 Jan restart |
| RSK-18 | The release-critical tail (about 51 pair-days, §11.4) is not complete by 29 January | High | High | MA-08 and MA-16; options in §11.5; release target option E |

---

## 13. Mentor Approval Points

| ID | Approval | Needed by | Related items |
|----|----------|-----------|---------------|
| MA-01 | Confirm MVP priorities, feature order and the exclusion of F10–F12 | Before Sprint 1 | §2 |
| MA-02 | Confirm the provisional choices in §6.1 (DEC-14 and DEC-06) | Sprint 1 day 1 | §6.1 |
| MA-03 | Approve the isolation mechanism (DEC-08) and its security requirements | Start of Sprint 4 (14 Dec) | DEC-08, F07-02a, F07-02b, INF-10 |
| MA-04 | Approve hosting and budget (DEC-07), and any cloud resource creation | Start of Sprint 4 (14 Dec) | DEC-07, SPK-04, INF-11 |
| MA-05 | Confirm Auth0 tenant ownership and admin access | Sprint 1 (for SPK-02), Sprint 2 (for INF-09) | DEC-03, INF-09, SPK-02 |
| MA-06 | Approve any significant scope change | As needed | [SDD Workflow §7](../sdd-workflow.md#7-handling-requirement-changes) |
| MA-07 | Confirm the demo arrangements for the holiday weeks, including the proposed 29 December demo | Sprint 3 | §1.4, RSK-11 |
| MA-08 | Choose one or more options in §11.5 to close the shortfall | Sprint 1 | §11, RSK-02, RSK-18 |
| MA-09 | Accept the final security review (HRD-01) before any public exposure, and approve option E if needed | Sprint 6 | HRD-01, F07-07 |
| MA-10 | Approve or decline each simplification in §3.2 | Sprint 1 | §3.2, §11.5 |
| MA-11 | Confirm the overhead assumptions in §11.3 and the staffing of all four pairs for the 60 working days | Sprint 1 | §1.3, §11.2, §11.3, RSK-03 |
| MA-12 | Decide on any spike extension beyond its timebox | As needed | §5 |
| MA-13 | Confirm the spike hypotheses in §6.2 are the right ones to test | Sprint 1 day 1 | §6.2, §9 |
| MA-14 | Decide DEC-09 (packaging and test visibility rules) | Sprint 2 | DEC-09, CON-04 |
| MA-15 | Decide how to handle the 4 pair-day over-commit in Sprint 1 (§9.4) | Sprint 1 day 1 | §9, RSK-15 |
| MA-16 | Decide on added capacity (option C in §11.5) and on the meeting cadence (option B) | Sprint 2 | §11.5, RSK-02, RSK-18 |

---

## 14. Items Not Assigned and Not Created

- No item is assigned to a named intern.
- No GitHub Issues have been created from this backlog. Issues should be created when items enter Sprint planning, using the feature issue template and the labels in the [SDD Workflow §6](../sdd-workflow.md#6-github-issue-states).
- Specifications for F01–F09 do not yet exist in `specs/features/`. Each is created from the template, and is approved before its implementation items start.
