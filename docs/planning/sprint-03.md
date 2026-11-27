# Sprint 3 Plan: 30 November–11 December 2026

**Status:** Mentor-approved Sprint 3 target and pair allocation; execution remains subject to SDD gates. No Sprint 3 issues are created or assigned by this document. Pair letters are workstream placeholders; no individual developers are named. Sprint 2 completion is based on the mentor's authoritative confirmation in [sprint-02.md §16](sprint-02.md#16-sprint-2-closure-addendum-mentor-confirmed).

**Planning sources:** [Product Backlog and decision register](../product/product-backlog.md), [MVP Proposal](../product/mvp-proposal.md), [SDD Workflow](../sdd-workflow.md), [Sprint 1 Plan](sprint-01.md), [Sprint 2 Plan](sprint-02.md), [Sprint 3 Issue Manifest](sprint-03-issues.md), feature specifications F01 and F04, contracts CON-01, CON-02 and CON-05, ADR-003 and ADR-004, and the Sprint 2 closure addendum. Apply the existing CON-04 artifact from approved PR #47 after merge; do not recreate it.

---

## 1. Sprint Goal

Move from preparation to working MVP slices: deliver authenticated caller identity, an editor component driven by approved problem inputs, and a published-problem persistence foundation, while preparing reproducible frontend/database builds.

All implementation remains gated by the SDD workflow. A task may start only when its specification is Approved and its required contracts and decisions are in place. Completion of Sprint 2 artifacts or Friday demonstrations does not satisfy those gates.

## 2. Expected Functional Increment

By Sprint end, the proposed increment is:

- The Go API can validate an authenticated caller and return the caller's own profile, backed by the agreed user data model.
- The browser can display and edit the Go starter template supplied to the editor from problem data. A local fixture or stub page may be used where the problem-details UI is not yet available, as permitted by the backlog's F04-01 soft dependency.
- The agreed schema supports user/profile and published-problem persistence; no problem-catalogue HTTP endpoint is included because F02-02 is not selected.
- API/orchestrator image builds and the frontend scaffold are reproducible through the existing CI path.

This is a set of demonstrable slices, not yet the complete browse-to-judge journey. F03 problem details, F06 submission dispatch, and F07 execution are not included. No learner code is executed in this sprint.

## 3. Capacity, Target and Feasible Alternative

| Measure | Pair-days |
|---------|-----------|
| Calendar gross capacity: 4 pairs × 10 working days | **40** |
| Documented effective task capacity (part-time planning limit) | **20** |
| Target committed workload | **23.5** |
| Target over documented effective capacity | **3.5** |
| Optional stretch workload | **4.0** |
| Target plus optional stretch / over effective capacity | **27.5 / 7.5** |
| Implementation, integration and automated testing | **19.0 (80.9% of target commitment)** |
| Specification/contract completion | **4.5 (19.1% of target commitment)** |

The mentor has approved the 23.5-pair-day target. This is 3.5 pair-days above the documented 20-pair-day effective planning baseline; the approval is specific to this Sprint 3 proposal, not proof of measured availability. The additional 4.0-pair-day stretch remains optional and requires confirmation of capacity.

**Capacity fallback if availability changes:** defer S3-07 and S3-08, reducing the workload to 18.0 pair-days and leaving 2.0 of the documented effective capacity. This is a contingency requiring mentor re-planning, not the approved target. The remaining work is 15.0/18.0 pair-days (83.3%) implementation, integration and testing.

Backlog estimates are retained for existing items (CON-04 2.0, INF-08 1.5, F04-01 4.0, F01-01 4.0, CON-05 1.0, INF-05 4.0). Story points are independent relative estimates, not converted from pair-days. The proposed INF-07 residual is 1.0 pair-day for basic metrics only and is not selected for Sprint 3.

## 4. Selected Backlog Items

| Work item | Backlog reference | Pair | Type | SP | Pair-days | Readiness |
|-----------|------------------|------|------|----|-----------|-----------|
| S3-01 | CON-04 | C | Contract | 3 | 2.0 | Use the existing approved PR #47 artifact after merge; dependent-pair agreement remains |
| S3-02 | INF-08 | A | Technical task | 2 | 1.5 | INF-02 is complete; Pair A owns both frontend scaffold and editor |
| S3-03 | F04-01 | A | Feature | 5 | 4.0 | Component starts after SPEC-F04 Approved; full template integration uses Pair A's INF-08 scaffold and agreed CON-04 terms |
| S3-04 | F01-01 | B | Feature | 5 | 4.0 | Requires SPEC-F01 Approved and INF-09 complete; approved Auth0 and PostgreSQL role-authority decisions; middleware can use repository test doubles |
| S3-05 | CON-05 Sprint 3 slice | D | Contract | 2 | 1.0 | Use merged CON-04 artifact; agree only the Sprint 3 schema slice, with relevant open F01 questions resolved |
| S3-06 | INF-05 | D | Technical task | 5 | 4.0 | Connection/test harness is independent; migrations/data access require agreement on the applicable CON-05 slice |
| S3-07 | SPEC-F02 | C | Feature specification | 3 | 1.5 | Pair C owns specification and persistence; a different pair must approve before S3-08 |
| S3-08 | F02-01 | C | Feature | 5 | 4.0 | Requires SPEC-F02 Approved by another pair, CON-05 slice agreement and INF-05; fixture tests can proceed beforehand |
| S3-09 | INF-12a | B | Technical task | 2 | 1.5 | INF-02 is complete; independent API/orchestrator image work |
| **Target total** | | | | **32** | **23.5** | |

**Optional stretch (not in the target commitment):** S3-10 / TST-01, 5 SP / 4.0 pair-days. Pair assignment remains open until the mentor confirms availability. It requires the merged CON-04 artifact, dependent-pair agreement and confirmed additional capacity.

### Deferred or excluded

- **INF-07 remainder:** only basic metrics may remain. HTTP bootstrap/configuration, liveness/readiness, middleware, timeouts, graceful shutdown, tests and documentation are complete under ST-04. The proposed 1.0 pair-day metrics-only estimate is not selected; confirm acceptance scope before future scheduling.
- **F06-01/F06-02, F07-01 and F08:** specifications, data/dispatch contracts, identity readiness, and downstream dependencies do not yet support a safe functional submission slice. Do not infer approval from completed drafts.
- **F02-02 and F03-01:** F02-02's catalogue endpoint is deferred; applicable API/contract readiness must be confirmed. F03-01 additionally lacks SPEC-F03 and depends on CON-04, CON-05 and INF-05.
- **Judge execution, F05 and any public exposure:** DEC-02 is not approved, SPK-01 reported No-go, and DEC-08 remains open. SPK-06 completion is validation evidence, not architecture approval.
- **SPK-04 and further research-only work:** not selected. Remaining research is limited to the exact decisions that block implementation.

## 5. Pair Allocation

The mentor confirms the post-rotation pair composition. Pair labels are not tied to people. The established SDD rotation policy applies: pairs rotate at sprint boundaries, and the mentor assigns individuals; this plan makes no individual assignments.

| Pair | Work items | Pair-days | Substantive coding task |
|------|------------|-----------|-------------------------|
| A | S3-02 INF-08; S3-03 F04-01 | **5.5** | Frontend scaffold and browser editor |
| B | S3-04 F01-01; S3-09 INF-12a | **5.5** | Authentication/profile API and container images |
| C | S3-01 CON-04; S3-07 SPEC-F02; S3-08 F02-01 | **7.5** | Problem contract, catalogue specification and published-problem persistence |
| D | S3-05 CON-05; S3-06 INF-05 | **5.0** | Schema slice and database foundation |
| **Target total** | | **23.5** | Every pair owns substantive implementation work |

Pair C has the largest target allocation at 7.5 pair-days; the approved total remains within the 40 calendar-gross pair-days. The optional S3-10 assignment is open and is not included in any pair load. If the 18-pair-day fallback is needed, defer S3-07/S3-08 and replan pair loads with the mentor.

## 6. Dependencies and Implementation Sequence

1. **Approval preflight:** reviewer pairs independently decide whether SPEC-F01, SPEC-F02 and SPEC-F04 are Approved. DEC-03, DEC-09, the PostgreSQL role-authority choice and DEC-11/MA-23 are approved and are not pending gates. SPEC-F01 OQ-2 is recorded as resolved by ADR-004; remaining questions and reviewer approval are separate.
2. **CON-04 artifact:** PR #47 is approved and awaiting merge, per mentor confirmation. Do not recreate the contract. S3-01 uses the existing artifact after merge, records the approved DEC-09 choices, and obtains dependent-pair agreement; CON-04 does not become Approved automatically.
3. **Independent technical starts:** after INF-02, Pair A builds INF-08 and then integrates F04-01; Pair B prepares INF-12a. Pair D may begin the schema-independent PostgreSQL connection configuration and integration-test harness within INF-05, but must not create migrations, tables or generated models before agreement on the applicable CON-05 slice.
4. **Parallel feature components:** after SPEC-F04 is Approved, Pair A may build the editor against approved inputs and fixtures, then integrate with its scaffold. Pair B may implement Auth0 token validation and middleware against fixtures and a repository test double, independently of PostgreSQL. Caller provisioning and `GET /api/v1/me` persistence integrate only after CON-05/INF-05 are agreed and ready.
5. **Problem-data path:** Pair C drafts SPEC-F02; a reviewer pair other than C must approve it before Pair C implements F02-01. Pair C can build tests around approved repository interfaces and fixtures, but problem persistence waits for the approved specification, agreement on the relevant CON-05 slice and the INF-05 adapter.
6. **Contract and schema sequence:** Pair C updates the existing CON-04 artifact after PR #47 merges, using the approved options and obtaining dependent-pair agreement. Pair D agrees only the CON-05 Sprint 3 slice after relevant SPEC-F01 schema questions are resolved; unrelated contract questions remain open. Schema-dependent INF-05 migrations use approved ADR-003 and follow the agreed slice.
7. **Integration points:** first freeze approved feature/contract interfaces; then connect auth and problem repositories to the agreed database adapter; then integrate the editor with the scaffold and CON-04 input. Use fixtures/mocks before these points, but do not claim end-to-end acceptance until real adapters and interfaces are integrated. One owner at a time for shared migration, router and CI files; every PR gets cross-pair review.

## 7. Approved Decisions and Remaining SDD Gates

### Mentor-approved decisions

| Decision | Approved choice | Evidence / implementation boundary |
|----------|----------------|--------------------------------|
| DEC-03 | Auth0; `github.com/auth0/go-jwt-middleware/v3` is the API authentication foundation. | Successful SPK-02 validation outcome; recorded in ADR-004. Auth0 validates identity/JWTs; it does not grant application roles. |
| SPEC-F01 OQ-1 / CON-05 C05-OQ-01 | PostgreSQL is authoritative for application roles and permissions. | Recorded in ADR-004 and the decision register. Do not grant privileges solely from JWT role claims. |
| DEC-11 / MA-23 | `pressly/goose`. | ADR-003 is Approved; its rationale and alternatives are retained. |
| DEC-09 §6.1 | Option A: PostgreSQL rows. | PostgreSQL remains the problem/test-asset storage system. |
| DEC-09 §6.2 | Option B: Go test files per problem. | This is the approved test organization. |
| DEC-09 §6.3 | Option B: function, type or handler in a fixed package. | Supports algorithmic functions and HTTP-handler challenges. |
| DEC-09 §6.4 | Option B: verdict and permitted test counts. | Never disclose hidden-test content, names, inputs, expected outputs or unfiltered diagnostics. DEC-09 does not resolve DEC-08 or other isolation decisions. |

### Required to start selected work

The decisions above are approved mentor decisions, not pending gates. Feature-spec approvals remain with reviewer pairs; contract agreement remains with dependent pairs.

| SDD approval gate | Required approver | Before | Blocks |
|-------------------|------------------|--------|--------|
| SPEC-F01 Approved, with approved decisions recorded and remaining open questions handled | Reviewer pair | F01-01 | S3-04 |
| SPEC-F02 Approved | Reviewer pair other than Pair C | F02-01 | S3-07 → S3-08 |
| SPEC-F04 Approved | Reviewer pair other than the responsible pair (confirm after rotation) | F04-01 | S3-03 |
| Merged CON-04 artifact and agreement on the applicable terms | All dependent pairs | CON-04-dependent implementation | S3-01, S3-03, S3-05, S3-08, optional S3-10 |
| Agreement on the CON-05 Sprint 3 schema slice | Dependent pairs, after relevant F01 schema questions are resolved | Schema-dependent INF-05 work | S3-05, S3-06, S3-08 |

SPEC-F02 does not yet have a file in the documented inventory. Drafting it is not approval; a reviewer pair other than C must approve it before S3-08 implementation. SPEC-F01 OQ-2 is resolved by approved DEC-03 / ADR-004; OQ-1 is resolved by the PostgreSQL role-authority decision. OQ-9 and other applicable questions remain for SPEC-F01/CON-05 review. This does not approve SPEC-F01.

DEC-13 (editor session preservation) is not a dependency of the scoped F04-01 backlog item. It need not block the editor component if reviewers approve SPEC-F04 with preservation separated as F04-02; if the reviewer considers it a blocker to SPEC-F04 approval, resolve it before approval. No F04-02 implementation is planned.

### Unresolved but not blocking this selection

DEC-04/CON-02 agreement and DEC-12 dispatch consistency block future F06-02, not the selected work. DEC-02/DEC-08 block production judge work; DEC-07 is a Sprint 4 hosting decision. No work in this plan assumes those decisions.

## 8. Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| The approved CON-04 artifact is not merged from PR #47 when needed | Contract-dependent implementation cannot use the agreed interface | S3-01 uses the existing artifact after merge; do not recreate it. Keep dependent integration gated until the artifact is available and pairs agree. |
| Target exceeds the documented 20 effective pair-day baseline | Mentor has approved the 23.5-pair-day Sprint 3 target; optional work requires additional capacity confirmation | Track workload against the approved target; do not infer availability for optional TST-01 from Sprint 2 estimates. |
| SPEC-F01 or SPEC-F04 remains In Review, or SPEC-F02 is not approved | Corresponding feature implementation cannot start | Reviewer pairs complete the existing SDD review and record status. No Sprint 2 completion or Friday demo is treated as approval. |
| CON-05 Sprint 3 slice cannot be agreed until relevant F01 schema questions close | Database and identity integration are delayed | PostgreSQL role authority is approved, but relevant schema representation questions and dependent-pair agreement remain. Pair B can complete token validation against a repository test double; schema-dependent integration waits. |
| Sprint 3 schema slice is unresolved | Schema implementation cannot start | Pair D limits pre-agreement work to connection/test-harness setup with no tables or generated models; migrations use approved ADR-003 only after the applicable slice is agreed. |
| INF-07 metrics acceptance scope remains unclear | Residual task could duplicate ST-04 or introduce an unapproved metrics design | Leave the proposed 1.0 metrics-only estimate outside Sprint 3; do not schedule HTTP/bootstrap/health/middleware work or select a metrics technology by assumption. |
| Pair C's 7.5 target pair-days reduce integration margin | Pair C owns the specification and corresponding persistence work | Keep specification review with a different pair and sequence F02-01 only after approval; use fixtures while the agreed data adapter is in progress. |
| Isolation approach is still unapproved | No production judge or public execution | No submitted code is executed. Carry DEC-08 to its Sprint 4 decision point and preserve the security gate. |
| Remaining effort still exceeds calendar capacity | MVP date at risk | Mentor decides scope/capacity/release target; do not reduce mandatory security controls. |

## 9. Integration and Testing Strategy

- Each feature implementation follows only its Approved specification. Unit and component tests cover the acceptance criteria in that spec.
- INF-05 may prepare a PostgreSQL connection/test harness without schema assumptions. Migration tests and data-access integration use only the agreed CON-05 Sprint 3 slice and mentor-approved ADR-003 option; unrelated unresolved CON-05 sections remain open.
- F01-01 tests token validation, missing/invalid tokens, claims and safe logging with fixtures. Middleware and handler unit tests use a repository test double; database provisioning/profile integration is verified only after the agreed adapter is ready.
- F02-01 uses fixtures/test doubles for published/unpublished behavior and repository errors; persistence integration is tested against the agreed CON-05 schema only. It does not implement the catalogue endpoint F02-02.
- F04-01 tests editing and starter-template behavior with approved inputs/fixtures, including the spec's blank/unusable-template, NUL, size-limit and untrusted-text cases. It does not test preservation, sample execution or submission.
- INF-08 and INF-12a builds/image checks use existing CI conventions. No new framework, router, metrics backend or infrastructure is introduced.
- Pair D owns the shared migration directory; shared router/CI files have one owner at a time. API and editor work can proceed on component/fixture paths before database and frontend integration, then integrate serially against approved interfaces.
- Require another pair's review for every issue. Friday demos are progress checkpoints and do not replace spec approval, contract agreement, PR review, or release acceptance.

## 10. Friday Demonstrations

| Date | Progress checkpoint |
|------|---------------------|
| Fri 4 Dec | Show the frontend scaffold, independent auth/editor components and fixture tests where their specs are approved, CON-04/CON-05 decision status, and schema-independent database harness progress. |
| Fri 11 Dec | Demonstrate the authenticated caller-profile path, editor component with approved starter-template fixture, and published-problem persistence if SPEC-F02/CON-05/INF-05 gates have cleared. Show container build checks, automated test evidence and unresolved gates. |

These are progress demonstrations only, not formal approvals, merge gates, or release gates.

## 11. Sprint Exit Criteria

- All committed work meets its issue acceptance criteria, is reviewed by another pair, and is merged under the SDD workflow.
- SPEC-F01/SPEC-F02/SPEC-F04 are marked Approved only by their reviewer pairs; if a selected feature spec remains unapproved, its implementation does not count as complete.
- The merged CON-04 artifact reflects approved DEC-09 choices and applicable dependent-pair agreement; DEC-09 approval does not automatically make CON-04 Approved.
- The applicable CON-05 Sprint 3 schema slice is agreed and implemented by dependent work; do not require or imply finalization of the entire product data model.
- INF-05 uses `pressly/goose` under approved ADR-003 and the agreed CON-05 Sprint 3 slice.
- Tests and CI pass for changed code; no credentials, tokens, learner source, or hidden tests are added to logs or documentation.
- No production judge execution or public exposure is enabled by Sprint 3.
- The mentor-approved 23.5-pair-day target is completed or explicitly reported as incomplete. Optional TST-01 proceeds only with separate capacity confirmation; its pair assignment remains open until then.

## 12. Roadmap and Remaining Capacity Forecast

| Point in plan | Estimated backlog remaining | Capacity/forecast |
|---------------|-----------------------------|-------------------|
| Sprint 2 closure | **170.0 pair-days** | Includes proposed 1.0 pair-day for basic metrics only; unsized approval overhead excluded |
| After target Sprint 3 commitment | **146.5 pair-days** | 23.5 mentor-approved; 3.5 above documented 20 effective baseline |
| After optional TST-01 | **142.5 pair-days** | Only if the 4.0 pair-day optional test-fixture scope is authorized and completed |
| After capacity-safe Sprint 3 alternative | **152.0 pair-days** | 18.0 committed of 20 effective; 2.0 reserve |
| Sprint 4–6 | Existing net forecast: **89 pair-days** | Historical base-case forecast from Product Backlog §11.3 |
| Conservative capacity through Sprint 6 | **109 pair-days** | 20 effective baseline in S3 + 89 net in S4–S6; closure shortfall **61** |
| Observed-throughput scenario | **112 pair-days** | 28 estimated pair-days × 4 remaining sprints; closure shortfall **58**, not a capacity guarantee |
| Target-plan comparison | **112.5 pair-days** | 23.5 mentor-approved S3 target + 89 S4–S6; closure shortfall **57.5** |

The estimated-work throughput scenario is only a comparison indicator: it uses the completed Sprint 2 estimate of 28 pair-days as a planning rate, not observed working hours or guaranteed future availability. At the target S3 commitment, 146.5 pair-days remain against 89 base-case pair-days in S4–S6 (a 57.5 gap). The capacity-safe alternative leaves a 63-pair-day gap after S3. Even if the target plus optional TST-01 were delivered, the S4–S6 base case leaves a 53.5-pair-day gap. Sprint 4 remains the proposed judge-integration milestone, but DEC-08 must be approved before any production execution path starts. Sprint 5–6 testing and release hardening remain dependent on APIs, staging, isolation and concurrency validation. Do not reduce mandatory MVP or security scope to close the forecast gap.

## 13. Validation Summary

- Sprint 1 and Sprint 2 completed scope is not reallocated; Sprint 2's original commitment remains 17.5 pair-days. Completed estimated Sprint 2 workload is 28 pair-days, not measured actual time.
- ST-02–ST-04 completed estimated effort is recorded once as an aggregate 15 SP / 10.5 pair-days and is not added to Sprint 3 capacity.
- ST-04's HTTP bootstrap/configuration, liveness/readiness, middleware, timeouts, graceful shutdown, tests and documentation are excluded from future scope. Only basic metrics may remain, with a proposed 1.0-pair-day metrics-only estimate; it is not selected.
- CON-04/CON-05 remainders are 2.0/1.0 pair-days; both are included once.
- The target commitment is 23.5 pair-days (19.0 implementation/integration/testing, **80.9%**) plus 4.0 optional testing stretch; it exceeds documented 20 effective pair-days by 3.5 before stretch.
- The fallback is 18.0 pair-days with 2.0 reserve and **83.3%** implementation/integration/testing. Optional stretch requires additional confirmed capacity.
- Pair target loads are A 5.5, B 5.5, C 7.5 and D 5.0; every pair owns substantive code. The optional stretch is unassigned.
- F01 authentication, editor and problem-persistence components use fixtures/test doubles to enable parallel work; integration waits for approved specs, contracts and schema. No provisional architecture is treated as approved.
- The increment is not a complete browse-to-judge journey; catalogue endpoint, submission dispatch, execution, staging and security release gates remain outside Sprint 3.
