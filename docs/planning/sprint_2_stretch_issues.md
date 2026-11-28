# Sprint 2 Stretch Issues

## 1. Purpose and scope

GoJudge's Sprint 2 pairs have completed their assigned work ahead of schedule. This document defines one additional, timeboxed task for each of Pairs A, B and C, aligned with their current workstreams and the Product Backlog. These are proposed planning items; this document does not create GitHub Issues or change existing plans.

The completed work supplied as context is treated as complete: Pair A's SPEC-F03, SPEC-F06 and SPEC-F08; Pair B's SPK-05 and isolation alternatives analysis; and Pair C's INF-09 and F01-01 preparation. The pairs remain responsible for their existing work, review comments and integration. None of the tasks below reopens or duplicates that work.

Each task is estimated at **5 story points / 3.5 pair-days**. Proposed stretch identifiers are separate from established backlog IDs. The tasks have no blocking dependencies on one another.

### Summary

| Pair | Work item ID | Title | Type | Story points | Pair-days |
|---|---|---|---|---:|---:|
| A | **Proposed ST-02**; proposed backlog item **TST-05** | Cross-feature consistency and acceptance scenarios | Technical task | 5 | 3.5 |
| B | **Proposed ST-03**; proposed spike **SPK-06** | Isolation validation and security regression suite | Technical spike | 5 | 3.5 |
| C | **Proposed ST-04**; established backlog item **INF-07 (remaining scope)** | HTTP API foundation and middleware integration | Technical task | 5 | 3.5 |

`ST-02`–`ST-04` are proposed Sprint 2 stretch references, not existing backlog IDs. `TST-05` and `SPK-06` are proposed new backlog IDs. INF-07 is an established backlog item; Pair C's scope is an early implementation of its remaining HTTP/health-endpoint scope, not a second INF-07 item.

## 2. Detailed issue definitions

### Pair A — Cross-feature consistency and acceptance scenarios

**Work item ID:** Proposed stretch reference **ST-02**; proposed backlog item **TST-05**

**Title:** Cross-feature consistency and acceptance scenarios

**Type:** Technical task

**Workstream:** Specifications

**Estimate:** 5 story points / 3.5 pair-days

**Deliverable:** `docs/testing/acceptance-scenarios.md`

#### Backlog reference

This is proposed as TST-05, a new test-and-hardening backlog item, because the backlog's TST-02 is a later browser-driven end-to-end harness and depends on implementation. It prepares, but does not implement, future integration and end-to-end tests. It traces requirements from established backlog features SPEC-F03 (F03-01), SPEC-F06 (F06-01–F06-04) and SPEC-F08 (F08-01–F08-02).

#### Objective

Produce a reviewable, traceable set of cross-feature acceptance scenarios that makes the journey from problem access through solution submission to submission history testable, while exposing inconsistencies and decisions that must not be resolved by assumption.

#### Dependencies and prerequisites

- Use the supplied SPEC-F03, SPEC-F06 and SPEC-F08 work as the requirements baseline, together with CON-01, CON-04, the Product Backlog, and the SDD workflow.
- No other stretch task is a prerequisite. This task does not gate or replace specification approval, contract agreement, or implementation readiness.
- The document is planning/test-design material only. Any change to an approved specification or contract must follow the SDD change and review process.

#### Detailed implementation scope

- Compare the three feature specifications for overlapping terminology, actors, inputs, outputs, state transitions, ownership/access rules, hidden-test visibility, and submission-history behavior.
- Record each inconsistency, missing requirement, or unresolved cross-feature dependency with source references and its effect on scenario expected results.
- Define end-to-end acceptance scenarios for an authenticated learner accessing a published problem, submitting a solution, and seeing that submission in their history, including relevant queued and terminal-state views.
- Include negative and boundary scenarios for access control, hidden-test confidentiality, invalid submissions, repeated submissions, and pagination consistency.
- Create a requirements-to-scenarios traceability matrix. Give each scenario a stable ID and include preconditions, test data/setup, action, observable expected result, source requirement IDs, and any unresolved decision.
- Record mentor questions with the competing interpretations, affected requirements/scenarios, and the decision needed. Mark dependent expected results as unresolved rather than selecting an interpretation.
- Do not modify the existing specifications, contracts, implementation, or automated tests.

#### Subtask breakdown

| Subtask | Estimate |
|---|---:|
| Inventory and cross-check requirements across SPEC-F03, SPEC-F06 and SPEC-F08 and relevant contracts | 0.75 pair-days |
| Draft end-to-end, negative and boundary scenarios, including stable-dataset pagination checks | 1.25 pair-days |
| Build the traceability matrix and record inconsistencies, dependencies and mentor questions | 0.75 pair-days |
| Cross-pair review, revise the artifact, and prepare evidence for the pull request | 0.75 pair-days |
| **Total** | **3.5 pair-days** |

#### Deliverables

- `docs/testing/acceptance-scenarios.md`, containing the scope, conventions, scenario catalogue, traceability matrix, inconsistencies/dependencies, and unresolved mentor questions.
- A pull request reviewed by a pair other than Pair A. The artifact is not a replacement for, or approval of, any feature specification.

#### Measurable acceptance criteria

1. The traceability matrix lists every functional requirement and acceptance criterion in SPEC-F03, SPEC-F06 and SPEC-F08, or marks it not applicable with a rationale. Every listed item maps to one or more scenario IDs or is explicitly blocked by an unresolved decision.
2. The scenario set includes at least one complete learner journey from published-problem access through solution submission to submission-history visibility, and separately covers queued and terminal submission states where specified.
3. There is at least one explicit scenario for each required risk area: unauthorized or cross-user access, hidden-test confidentiality, invalid submission input, repeated submissions, and pagination consistency.
4. The pagination scenario uses a documented stable dataset and checks that traversing the defined pages yields no duplicate or omitted records. Any behavior under concurrent inserts, or any missing ordering/cursor rule, is recorded as an open question rather than assumed.
5. Every scenario has a unique ID, preconditions, action, observable expected result, and source requirement reference. Unresolved expected results identify the specific mentor decision needed.
6. The document records all identified cross-feature inconsistencies and missing dependencies, or explicitly states that none were found after the full requirement inventory. No open product or architectural decision is presented as settled.
7. A reviewer from another pair has reviewed the artifact, and review comments applicable to its scope are resolved or recorded as explicit open questions.

#### Definition of Done and closing conditions

- All acceptance criteria are evidenced in the pull request.
- The document is reviewed by a pair other than Pair A, updated for review feedback, and merged to `main`.
- No existing specification or contract is changed, and no automated test is implemented.
- Close the proposed work item only when its technical-task Definition of Done is met; the artifact itself does not approve any specification or contract.

#### Explicit exclusions

- Editing or approving SPEC-F03, SPEC-F06, SPEC-F08, CON-01, CON-04, or any other contract.
- Choosing duplicate-submission behavior, authorization policy, hidden-test exposure, pagination model/order, or other open product/API behavior.
- Implementing endpoints, UI, integration tests, browser tests, or any production code.
- Replacing or expanding the separately planned TST-02 end-to-end harness.

#### Open decisions and mentor approval requirements

- Ask the mentor to decide any conflicting requirements or unresolved behavior that affects expected results; include source references and the impact.
- In particular, do not assume the duplicate-submission rule, what result details may be shown without exposing hidden tests, or pagination behavior not fixed by the specifications/contracts.
- Material changes to an approved specification or agreed contract require the existing SDD change process and renewed review; they are outside this task's estimate.

#### Impact on subsequent sprint planning

The scenario catalogue and traceability matrix provide implementation-ready test-design input for F03-01, F06-01–F06-04, F08-01–F08-02 and later integration testing. Mentor decisions and any resulting specification changes remain separate work and must be planned before scenarios that depend on them can be treated as normative.

### Pair B — Isolation validation and security regression suite

**Work item ID:** Proposed stretch reference **ST-03**; proposed technical spike **SPK-06**

**Title:** Isolation validation and security regression suite

**Type:** Technical spike

**Workstream:** Security Architecture

**Estimate:** 5 story points / 3.5 pair-days

**Proposed report:** `docs/spikes/SPK-06.md`

#### Backlog reference

SPK-06 is proposed as a follow-on evidence spike using the established SPK-05 report and the Product Backlog's mandatory isolation controls (Product Backlog §3.1 and MVP Proposal §7) as its baseline. It prepares reusable validation for later F07 isolation work and future mechanism evaluations. It is not a repeat of SPK-05's technology comparison.

#### Objective

Define and demonstrate a repeatable, candidate-neutral validation suite for the isolation requirements, with controlled trusted fixtures and auditable experimental evidence, while preserving the project's E/D/H evidence classifications.

#### Dependencies and prerequisites

- Use SPK-05 findings and the existing isolation requirements as the baseline. The starting report, its known gaps and its evidence classifications are inputs, not items to re-run wholesale.
- Work only in a dedicated isolated experimental environment, using trusted fixtures. Do not change shared or production infrastructure.
- No other stretch task is a prerequisite. Any need to install or alter cluster-level infrastructure requires explicit mentor approval before doing so.

#### Spike brief

**Question:** Can the team define and repeatably execute a reusable set of controlled checks for the stated isolation requirements, while recording evidence and limitations consistently enough to support future candidate evaluation?

**Timebox:** 3.5 pair-days. Stop at the agreed timebox end; any extension requires mentor approval under the SDD workflow.

**Experiment:** Define a reusable matrix and controlled trusted fixtures for network isolation from the first instruction, CPU limits, memory limits, execution timeout, filesystem restrictions, output limits, and process limits. Exercise the suite only in an isolated experimental environment. Where an approved test environment or a numeric requirement is unavailable, record the gap as unverified/open rather than creating a new approval or threshold.

**Success criteria:** All seven control areas have a documented validation case and evidence status; the repeatable mechanism and evidence format are documented; controlled runs and their limitations are reported. Not verified or partially verified results are acceptable spike outcomes if reported accurately.

**Failure/stop criteria:** Stop and notify the mentor if a trusted fixture can affect resources outside the isolated environment, if a proposed check would require running learner-submitted code, or if the experiment requires unapproved cluster-wide changes. Record the event and residual risk in the report; do not broaden the spike to resolve it.

#### Detailed implementation scope

- Turn SPK-05 findings and applicable isolation requirements into a reusable validation matrix covering all seven areas above.
- Specify controlled, bounded, trusted fixtures and a repeatable runner or documented execution mechanism. Keep fixtures and experiment code on the existing `spike/` branch convention; do not merge experiment code to `main`.
- Define how each run records the requirement/control, environment and versions, fixture identifier, command or invocation, repetitions, observed result, E/D/H classification, evidence reference, limitation, and residual risk.
- Preserve the existing classifications: **E** (experimental), **D** (documentation-based), and **H** (hypothesis/untested). A D or H result must not be described as empirically verified.
- For network isolation, make the first test action the egress probe and record whether it succeeds or is blocked. Use five repeated runs when the isolated environment supports the test; record any unavailable or inconclusive run.
- For CPU, memory, timeout, filesystem, output, and process controls, perform three repeated runs per supported case. Record configured limits and observed behavior; use existing approved thresholds if defined, and mark missing thresholds as open rather than inventing them.
- Include reproducibility instructions, cleanup, environmental assumptions, coverage limitations, and residual risks.
- Do not compare technologies, recommend or select a production mechanism, or claim a decision has been approved.

#### Subtask breakdown

| Subtask | Estimate |
|---|---:|
| Map SPK-05 findings and isolation requirements into a reusable validation matrix; identify threshold/environment gaps | 0.75 pair-days |
| Prepare bounded trusted fixtures and the repeatable execution/evidence-collection mechanism | 1.25 pair-days |
| Run supported cases in the isolated environment and classify results E/D/H | 1.0 pair-days |
| Write reproducibility/limitations report, obtain other-pair review, and revise | 0.5 pair-days |
| **Total** | **3.5 pair-days** |

#### Deliverables

- Proposed report `docs/spikes/SPK-06.md` with the validation matrix, per-case outcomes, E/D/H evidence, evidence references, reproducibility steps, limitations, and residual risks.
- Trusted test fixtures and experiment mechanism on the spike branch only; no experiment code merged to `main`.
- A reviewed report with the spike outcome stated as Go, Conditional go, No-go, or Inconclusive, as supported by the evidence.

#### Measurable acceptance criteria

1. The validation matrix has a separate row for network, CPU, memory, timeout, filesystem, output, and process limits. Every row states the requirement source, fixture/check, expected observation or unresolved threshold, evidence class, evidence location, and limitations.
2. The fixture set contains a bounded trusted test for each control area. The report states that no learner-submitted code was executed and identifies the isolated environment used.
3. The repeatable mechanism documents the invocation, prerequisites, per-run output/evidence, and cleanup. Another pair can follow the written steps in the same available isolated environment without relying on undocumented setup.
4. The network check initiates its first egress probe at fixture start and records five runs when supported. Each of the other supported control checks records three runs. Deviations are reported with a reason; no missing run is represented as a pass.
5. Each result is classified E, D, or H consistently. Experimental claims include the environment/version, fixture/check, observed result, run count and evidence reference. Documentation claims include their source. Untested behavior remains H.
6. The report includes reproducibility instructions, explicit scope limitations, unresolved thresholds or prerequisites, and residual risks for every partially verified or unverified control.
7. The report does not repeat the SPK-05 candidate comparison, approve DEC-08 or DEC-02, select a production mechanism, execute learner-submitted code, or implement the production judge.
8. A pair other than Pair B reviews the report; applicable comments are resolved or recorded as limitations/open questions.

#### Definition of Done and closing conditions

- Respect the timebox and record evidence/status for every success criterion, including unverified results.
- State the outcome, report any stop/failure criterion to the mentor when met, and list follow-on decisions without making them.
- Merge the reviewed report to `docs/spikes/SPK-06.md`; keep experimental code off `main`.
- Close the proposed spike only when the spike Definition of Done is met. No decision or production approval is implied.

#### Explicit exclusions

- Repeating SPK-05's technology comparison or issuing another DEC-08 recommendation.
- Approving DEC-08 or DEC-02, selecting a production isolation mechanism, or changing the approved architecture.
- Running learner-submitted code, production judge implementation, or production-like untrusted workloads.
- Installing or changing shared/production cluster plugins, runtimes, policies, or other infrastructure without mentor approval.
- Claiming coverage for a control that was not actually measured.

#### Open decisions and mentor approval requirements

- Obtain mentor approval before any cluster-wide change or new experimental infrastructure is introduced.
- Ask the mentor to resolve any missing numeric isolation threshold needed to interpret a result; until then, report the observation without a pass/fail claim against an invented limit.
- Any mechanism-selection or architecture decision remains outside SPK-06 and follows the existing DEC-08/DEC-02 approval gates.

#### Impact on subsequent sprint planning

SPK-06 gives future F07/INF-10 work a repeatable security regression baseline and exposes evidence gaps for mentor prioritization. The report can inform later isolation evaluation but does not remove DEC-08/DEC-02 gates, establish production readiness, or reduce implementation estimates by assumption.

### Pair C — HTTP API foundation and middleware integration

**Work item ID:** Proposed stretch reference **ST-04**; established backlog item **INF-07 (remaining scope)**

**Title:** HTTP API foundation and middleware integration

**Type:** Technical task

**Workstream:** Backend

**Estimate:** 5 story points / 3.5 pair-days

**Deliverables:** HTTP server bootstrap/configuration, liveness and readiness endpoints, tests, and server usage documentation

#### Backlog reference

This is an early implementation of the remaining INF-07 scope identified in the Sprint 2 planning references: health endpoints and HTTP server wiring. It integrates the existing INF-07 structured logging and request-ID middleware. Track the work against the existing INF-07 backlog item; do not create a duplicate INF-07 item or count this estimate in addition to the same remaining INF-07 work without mentor reconciliation.

#### Objective

Provide a minimal, well-tested HTTP server foundation with safe lifecycle behavior, health endpoints, and the existing observability middleware, ready to host future authentication and functional endpoints without implementing them now.

#### Dependencies and prerequisites

- INF-01 (Go module/repository foundation) and the existing INF-07 logging and request-ID middleware/documented middleware contract.
- Follow repository conventions and approved architectural decisions. Use existing dependencies only; do not add a third-party dependency.
- No other stretch task is a prerequisite. F01-01 work is explicitly not duplicated or required as a dependency for this foundation.

#### Detailed implementation scope

- Implement server bootstrap and typed configuration for the listen address, HTTP server timeouts, and graceful-shutdown deadline, with documented defaults and validation.
- Register liveness and readiness endpoints. Liveness reports that the process is running. Readiness reports whether the server is accepting requests and becomes not-ready during shutdown; it must not claim database, authentication, or judge dependencies are healthy when those checks are not part of this task.
- Integrate the existing request-ID and access-log middleware using the documented order: RequestID outside AccessLog. Health requests receive normal request IDs and access-log records without logging query strings, credentials, or request bodies.
- Implement graceful shutdown and set non-zero, documented HTTP server timeouts appropriate to the repository's existing conventions.
- Prepare an extensible server/handler registration point for future authentication and functional endpoints, without adding them.
- Add focused unit and integration tests using existing Go testing conventions, and document configuration, startup, shutdown, endpoints, and middleware integration.

#### Subtask breakdown

| Subtask | Estimate |
|---|---:|
| Define server configuration, bootstrap and lifecycle behavior | 0.75 pair-days |
| Add liveness/readiness handlers and integrate existing middleware | 1.0 pair-days |
| Add unit and integration tests for health behavior, request correlation and graceful shutdown | 1.0 pair-days |
| Document configuration/startup/shutdown, obtain cross-pair review, and address feedback | 0.75 pair-days |
| **Total** | **3.5 pair-days** |

#### Deliverables

- Minimal HTTP server bootstrap with validated configuration, bounded timeouts, and graceful shutdown.
- Liveness and readiness endpoints, integrated with existing request-ID and access-log middleware.
- Unit/integration tests covering the acceptance criteria.
- Documentation describing configuration, startup, shutdown, endpoint semantics, and middleware order.

#### Measurable acceptance criteria

1. Server configuration covers the listen address, all configured HTTP server timeouts, and shutdown deadline. Invalid configuration is surfaced as an error; documented defaults are applied only where repository convention allows them.
2. Tests start the server on a test listener and verify that the liveness endpoint returns success while running, and the readiness endpoint reports ready while accepting requests and not-ready during shutdown.
3. The configured server timeouts and shutdown deadline are non-zero and their effective values are asserted by tests or by a focused configuration test.
4. Integration tests verify that health requests pass through the documented middleware order, return an `X-Request-ID`, and produce an access-log record containing the request ID and status. Tests also verify that a query-string sentinel is absent from the log.
5. A shutdown test triggers graceful shutdown, confirms it completes within the configured deadline, and confirms the server no longer reports ready. Shutdown errors are surfaced rather than silently ignored.
6. The server exposes a registration point for future handlers, but no authentication middleware or functional/business endpoint is added.
7. No database migration or new third-party dependency is introduced. Tests run using the repository's documented test/build conventions.
8. Documentation describes the configuration values/defaults, how to start and stop the server, endpoint semantics, and middleware integration/order.
9. A reviewer from another pair reviews the change; applicable comments are resolved or recorded.

#### Definition of Done and closing conditions

- All acceptance criteria are evidenced in the pull request; tests and applicable CI checks pass.
- Documentation is updated, the change is reviewed by another pair, and merged to `main`.
- No secret, environment-specific configuration, new authentication logic, business endpoint, migration, or unapproved dependency is included.
- Close the proposed stretch work only when the technical-task Definition of Done is met. Reconcile the remaining INF-07 estimate with the mentor so the same scope is not counted twice.

#### Explicit exclusions

- Repeating INF-09 tenant configuration or F01-01 authentication preparation/implementation.
- Authentication, authorization, token validation, or role logic.
- Business/functional API endpoints, database access, migrations, or judge integration.
- Basic metrics or other observability features beyond the existing logging/request-ID middleware and requested health endpoints.
- New third-party dependencies or an unapproved router/server architecture.

#### Open decisions and mentor approval requirements

- Follow any existing repository convention for configuration source and endpoint naming. If none exists, document the minimal choice in the pull request; seek mentor approval before making a decision that changes an agreed API contract or introduces a new architectural dependency.
- Keep readiness local to server lifecycle in this task. Any future database or external-dependency readiness checks require their own agreed scope and must not be implied by these endpoints.
- The mentor should reconcile this early INF-07 delivery with the INF-07 remaining-scope estimate during subsequent planning.

#### Impact on subsequent sprint planning

This work advances part of INF-07 earlier than planned and provides a place to integrate future authentication and functional handlers. Subsequent plans should mark the delivered health/server-wiring scope as consumed, size any remaining INF-07 work (including metrics, if still required), and avoid counting this 3.5 pair-day task twice.

## 3. Total estimated effort

| Pair | Story points | Pair-days |
|---|---:|---:|
| A | 5 | 3.5 |
| B | 5 | 3.5 |
| C | 5 | 3.5 |
| **Total** | **15** | **10.5** |

These are proposed stretch estimates and are additional to existing Sprint 2 estimates. The project mentor must confirm capacity and the Sprint 2 timebox before authorizing any item.

## 4. Cross-pair dependency summary

| Stretch item | Required inputs | Dependency on other stretch items |
|---|---|---|
| ST-02 / TST-05 (Pair A) | SPEC-F03, SPEC-F06, SPEC-F08, relevant contracts and SDD workflow | None |
| ST-03 / SPK-06 (Pair B) | SPK-05 findings and existing isolation requirements; approved isolated test environment | None |
| ST-04 / INF-07 remaining scope (Pair C) | INF-01 and existing INF-07 logging/request-ID middleware | None |

There are **no blocking dependencies between ST-02, ST-03 and ST-04**. Cross-pair review remains required by the SDD workflow, but review is not an implementation dependency. Any unresolved mentor decision or prerequisite is recorded within the relevant item and does not silently transfer to another pair's stretch task.
