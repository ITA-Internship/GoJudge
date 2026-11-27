# Sprint 3 Issue Manifest (30 November–11 December 2026)

**Status:** Mentor-approved Sprint 3 target and pair allocation. This manifest is planning documentation only; no GitHub Issues are created or updated. Pair assignments use A–D only; do not assign individual developers. Sprint 2 completion does not imply approval of unrelated specifications or contracts.

**Capacity:** calendar gross capacity is 40 pair-days; documented effective task capacity is 20. The mentor approved the 23.5-pair-day target (32 SP); optional stretch is an additional 4.0 pair-days (5 SP), subject to availability confirmation. If availability changes, defer S3-07/S3-08 for an 18.0-pair-day fallback and replan pair loads with the mentor. Story points are not converted to pair-days.

**Type labels:** `type:contract`, `type:technical-task`, and `type:feature`, following the [SDD Workflow §6](../sdd-workflow.md#6-github-issue-states). Every proposed issue has exactly one type label.

---

## 1. Summary

| Work item | Backlog reference | Proposed issue title | Type / label | Pair | SP | Pair-days |
|-----------|-------------------|----------------------|--------------|------|----|-----------|
| S3-01 | CON-04 | [CON-04] Finalize problem packaging and test-visibility contract | Contract / `type:contract` | C | 3 | 2.0 |
| S3-02 | INF-08 | [INF-08] Create the frontend scaffold and CI build | Technical task / `type:technical-task` | A | 2 | 1.5 |
| S3-03 | F04-01 | [F04-01] Integrate the browser editor and problem starter template | Feature / `type:feature` | A | 5 | 4.0 |
| S3-04 | F01-01 | [F01-01] Implement API authentication and caller profile | Feature / `type:feature` | B | 5 | 4.0 |
| S3-05 | CON-05 Sprint 3 slice | [CON-05] Agree the Sprint 3 data-model slice | Contract / `type:contract` | D | 2 | 1.0 |
| S3-06 | INF-05 | [INF-05] Implement PostgreSQL migrations and data access | Technical task / `type:technical-task` | D | 5 | 4.0 |
| S3-07 | SPEC-F02 | [SPEC-F02] Specify the problem catalogue | Feature / `type:feature` | C | 3 | 1.5 |
| S3-08 | F02-01 | [F02-01] Implement published-problem persistence | Feature / `type:feature` | C | 5 | 4.0 |
| S3-09 | INF-12a | [INF-12a] Build API and orchestrator container images | Technical task / `type:technical-task` | B | 2 | 1.5 |
| **Target total** | | | | | **32** | **23.5** |
| S3-10 (optional) | TST-01 | [TST-01] Create reference challenge fixtures | Technical task / `type:technical-task` | Open | 5 | 4.0 |

### Pair reconciliation

| Pair | Issues | Pair-days | Capacity remaining |
|------|--------|-----------|--------------------|
| A | S3-02, S3-03 | 5.5 | — |
| B | S3-04, S3-09 | 5.5 | — |
| C | S3-01, S3-07, S3-08 | 7.5 | — |
| D | S3-05, S3-06 | 5.0 | — |
| **Target total** | | **23.5** | **3.5 above the 20 effective pair-day limit** |

**Optional stretch:** S3-10 / TST-01, 5 SP and 4.0 pair-days; pair assignment remains open until the mentor confirms availability. The total including stretch would be 27.5 pair-days. It is not part of the target commitment and requires the merged CON-04 artifact and dependent-pair agreement.

**Capacity fallback:** defer S3-07/S3-08, leaving 18.0 pair-days and 2.0 pair-days against the documented effective baseline. Its implementation/integration/testing share is 15.0/18.0 (83.3%). Replan pair loads with the mentor if this fallback is needed; do not infer new scope.

The mentor confirms pair composition after the established sprint-boundary rotation. Pair labels are not individual assignments.

## 2. Shared Readiness and SDD Rules

1. Feature implementation S3-03, S3-04 and S3-08 may start only when their corresponding feature specifications are **Approved** by reviewer pairs other than their authors. Pair C owns both SPEC-F02 and F02-01; a different pair must approve SPEC-F02 before S3-08 implementation.
2. Contract work follows cross-pair review and agreement by dependent pairs. A completed draft is not final agreement.
3. DEC-03, PostgreSQL role authority, DEC-09 and DEC-11/MA-23 are approved and recorded in the decision register; ADR-003/ADR-004 record the applicable architecture decisions. SPEC-F01 remains subject to reviewer-pair approval and remaining open questions; CON-05 still requires agreement on the applicable Sprint 3 schema slice.
4. INF-09 is treated as completed under the mentor's blanket Sprint 2 closure confirmation; its details are not reconstructed.
5. Every PR is reviewed by at least one pair other than the owner. Friday demonstrations are not approval or release gates.
6. S3-04 token validation may use fixtures and a repository test double before the persistent adapter is ready. S3-03 may build the editor component using approved inputs and fixtures before full frontend integration. S3-06 may prepare only schema-independent connection/test-harness work before the CON-05 slice is agreed. If an approved interface or contract is unavailable, stop that dependent portion rather than inventing it.
7. No task introduces a new technology, infrastructure component, router, or architectural decision. If the approved specification or contract leaves a blocking implementation choice open, stop that portion and request the relevant approval.

## 3. Proposed Issues

### S3-01: [CON-04] Finalize problem packaging and test-visibility contract

**Backlog ID:** CON-04 (remaining 3 SP / 2.0 pair-days; broader approval remains subject to SDD gates)
**Type / label:** Contract / `type:contract`
**Assigned pair:** C
**Estimate:** 3 SP; **2.0 pair-days**
**Objective:** Update the existing CON-04 draft from the approved artifact after PR #47 is merged, recording the approved DEC-09 choices and obtaining dependent-pair agreement without automatically marking the contract Approved.

**Dependencies and prerequisites**

- Use the existing artifact from approved PR #47 after merge; do not recreate CON-04.
- DEC-09 is approved. Record its choices: §6.1 Option A (PostgreSQL rows); §6.2 Option B (Go test files per problem); §6.3 Option B (function, type or handler in a fixed package); §6.4 Option B (verdict and permitted counts).
- Cross-pair review and agreement from every dependent pair named in the contract.

**References**

- Product Backlog §4.1 CON-04, §4.4 F03-01/F04-01/F07-04/F09-02.
- [MVP Proposal §3 and §7](../product/mvp-proposal.md).
- [SDD Workflow §2.4 and §4](../sdd-workflow.md).
- Dependent artifacts: SPEC-F03, SPEC-F04, CON-05, and the future judge/admin work that consumes the contract.
- DEC-09 approved choices in the Product Backlog decision register and Sprint 3 plan.

**Detailed scope and subtasks**

1. Review the merged artifact and apply the approved DEC-09 choices (0.25 pair-day).
2. Record PostgreSQL row storage and Go test files per problem; support algorithmic-function and HTTP-handler challenge solutions in a fixed package (0.5 pair-day).
3. Record learner-visible verdict and permitted test counts, excluding hidden-test content/names/inputs/expected outputs and unfiltered execution diagnostics (0.5 pair-day).
4. Record downstream data dependencies without adding unrequested fields or formats (0.25 pair-day).
5. Obtain cross-pair review and dependent-pair agreement; preserve all remaining open questions and the correct contract status (0.5 pair-day).

**Deliverables**

- Updated existing `docs/contracts/CON-04.md` after PR #47 merges; no artifact recreation.
- Approved DEC-09 choices and named dependent-pair agreement recorded in the artifact.
- Remaining blocking questions and contract status accurately recorded; no automatic Approved status.

**Acceptance criteria**

1. The PR #47 artifact is merged and reviewed. The four approved DEC-09 options are recorded accurately.
2. The format supports algorithmic functions and HTTP-handler challenges. Learners receive only the verdict and permitted test counts; hidden content, names, inputs, expected outputs and unfiltered diagnostics are never exposed.
3. Remaining open questions are recorded as blocking where applicable; DEC-09 approval is not presented as resolving DEC-08 or other isolation decisions.
4. CON-05 and the dependent specifications can reference the final terms without inventing fields or formats.
5. Every dependent pair agrees after cross-pair review; unresolved disagreement prevents the contract being marked agreed.

**Testing/review:** Contract consistency review against the MVP proposal, SPEC-F04, CON-05, and all dependents. At least one reviewer pair other than C.

**Definition of Done:** Existing contract artifact updated after PR #47 merge, reviewed, and agreed by dependent pairs; remaining blockers are explicit. Do not mark CON-04 Approved unless all SDD approval conditions are met.

**Explicit exclusions:** implementation, migrations, hidden-test storage code, judge execution, choosing an isolation mechanism, F07 implementation, and approval of DEC-02/DEC-08.

**Mentor approval:** DEC-09 is already approved. Pair agreement and formal contract approval remain separate SDD gates.

---

### S3-02: [INF-08] Create the frontend scaffold and CI build

**Backlog ID:** INF-08
**Type / label:** Technical task / `type:technical-task`
**Assigned pair:** A
**Estimate:** 2 SP; **1.5 pair-days**
**Objective:** Provide the repository's planned React, TypeScript, and Vite frontend scaffold with the Monaco Editor dependency and a reproducible CI build.

**Dependencies and prerequisites**

- INF-02 CI baseline (completed Sprint 1 scope).
- No feature specification, API contract, auth provider, or editor-preservation decision is needed for scaffold-only work. INF-08 may proceed independently of S3-03's editor component.

**References**

- Product Backlog §4.2 INF-08, §9.4.
- [MVP Proposal §5](../product/mvp-proposal.md) (React/TypeScript/Vite and Monaco proposal).
- No ADR authorizes additional tools or infrastructure.

**Detailed scope and subtasks**

1. Add the scaffold needed to build the planned frontend (0.5 pair-day).
2. Include Monaco as a dependency only; do not implement editor behavior under this issue (0.25 pair-day).
3. Connect the frontend build to the existing CI baseline using the repository's established workflow (0.5 pair-day).
4. Document the local build command in existing project documentation where needed (0.25 pair-day).

**Deliverables**

- Reproducible frontend scaffold and CI build.
- Minimal local build instructions.

**Acceptance criteria**

1. The frontend scaffold uses the backlog's React, TypeScript, and Vite choices and includes the Monaco dependency.
2. The documented local build completes successfully.
3. The existing CI runs and verifies the frontend build.
4. No API, authentication, browser-editor behavior, or persistence feature is implemented.
5. No unrelated infrastructure or dependency is added.

**Testing requirements:** Local frontend build and CI build; review dependency changes and build output.

**Definition of Done:** Acceptance criteria evidenced in the PR, CI green, cross-pair PR approval, and changes merged.

**Explicit exclusions:** F04-01 editor behavior, F01 sign-in, API routes, DEC-13 preservation, additional UI frameworks, or hosting infrastructure.

**Mentor approval:** Not required unless implementation would introduce an unapproved technology or architecture decision.

---

### S3-03: [F04-01] Integrate the browser editor and problem starter template

**Backlog ID:** F04-01
**Type / label:** Feature / `type:feature`
**Assigned pair:** A
**Estimate:** 5 SP; **4.0 pair-days**
**Objective:** Implement the specified Go editor component and load the starter template supplied by problem data.

**Dependencies and prerequisites**

- SPEC-F04 is Approved by its reviewer pair before implementation.
- The editor component may be built independently with fixtures conforming to approved specification/contract interfaces; full scaffold and problem-template integration follows Pair A's INF-08 scaffold and the agreed CON-04 template terms.
- F03-02 is a soft dependency: use the backlog-permitted stub/fixture host if the problem-details UI is not ready.
- DEC-13 is not a dependency of the scoped F04-01 item. Any open question that reviewers identify as blocking SPEC-F04 approval must be resolved before approval.

**References**

- Product Backlog §4.4 F04-01 (size 4; F03-02 soft dependency), F04-02 (excluded).
- [SPEC-F04](../../specs/features/F04-browser-code-editor.md), especially FR-01–FR-14 and relevant approved acceptance criteria.
- Agreed CON-04 template terms; [MVP Proposal §3](../product/mvp-proposal.md).
- SDD Workflow §2.1 and §4.

**Detailed scope and subtasks**

1. Implement the editor component with Go editing support against the inputs in Approved SPEC-F04; use an approved-interface fixture where needed (1.0 pair-day).
2. Load the starter template from the problem-data interface/fixture defined by Approved SPEC-F04 and agreed CON-04 terms; do not invent a format or embed a default challenge template (0.75 pair-day).
3. Implement the approved behavior for blank/unusable templates, plain-text rendering, NUL rejection and the agreed size limit (1.0 pair-day).
4. Keep template/source content as text; do not execute or render it as HTML (0.25 pair-day).
5. Add focused component tests for the approved acceptance criteria (1.0 pair-day).

**Deliverables**

- Editor component hosted in the scaffold.
- Template loading path from problem data or the allowed local test fixture.
- Automated tests for editing and template behavior.

**Acceptance criteria**

1. A learner can edit Go source text; the editor offers no language selector.
2. A valid starter template is loaded from the problem-data input and appears unchanged.
3. Missing/unusable template handling matches the Approved spec and agreed CON-04 terms.
4. Pasted template/source content is plain text; NUL and over-limit behavior follows the Approved spec's resolved limits and does not silently truncate.
5. Typing/pasting does not execute, compile, or submit source code.
6. Unit/component tests cover the implemented acceptance cases, and CI passes.

**Testing requirements:** Component tests for normal editing, template loading/fallback, text-only rendering, NUL, and size boundaries; no live judge or submission path.

**Definition of Done:** SPEC-F04 was Approved before implementation; feature acceptance criteria are evidenced; cross-pair PR approval and CI pass; implementation merged.

**Explicit exclusions:** F04-02 session preservation (including DEC-13), F03-02 full problem-details UI, F05 sample execution, F06 submission, multi-file editing, other languages, and new frontend technologies.

**Mentor approval:** DEC-13 is not required for this bounded item unless a reviewer records it as a blocker to SPEC-F04 approval. DEC-09 choices are approved; the mentor does not approve the feature spec in place of its reviewer pair.

---

### S3-04: [F01-01] Implement API authentication and caller profile

**Backlog ID:** F01-01
**Type / label:** Feature / `type:feature`
**Assigned pair:** B
**Estimate:** 5 SP; **4.0 pair-days**
**Objective:** Implement the approved F01 authentication path in the Go API, including caller identity provisioning and own-profile retrieval.

**Dependencies and prerequisites**

- SPEC-F01 is Approved by its reviewer pair before implementation.
- DEC-03 (Auth0) and PostgreSQL as the authoritative role source are approved; SPEC-F01 remains subject to reviewer-pair approval and its remaining open questions. OQ-2's DEC-03 approval is documented in ADR-004.
- INF-09 is treated as complete based on Sprint 2's authoritative closure; verify the approved setup is available without exposing its secrets.
- Token validation and middleware unit tests do not depend on S3-05/S3-06; use fixture tokens and a repository test double. User provisioning and persisted profile integration wait for the agreed CON-05/INF-05 interface.
- Existing HTTP foundation and middleware integration from ST-04; no replacement router/server is selected here.

**References**

- Product Backlog §4.4 F01-01 (size 4), F01-02 (separate and excluded).
- [SPEC-F01](../../specs/features/F01-accounts-and-access.md), especially FR-01–FR-10, FR-13, FR-16–FR-17 and approved API/error contract.
- [ADR-004](../architecture/adr/ADR-004.md), agreed CON-05 Sprint 3 schema slice, INF-07 logging/request correlation, and SDD Workflow §2.1/§4. SPK-02 is validation evidence only.

**Detailed scope and subtasks**

1. Implement Auth0 token-validation middleware using `github.com/auth0/go-jwt-middleware/v3` and the Approved SPEC-F01 rules; add unit tests using fixture tokens and a repository test double (1.5 pair-days).
2. Implement provisioning only against the agreed CON-05 Sprint 3 slice and INF-05 interface; preserve idempotency for concurrent first requests (0.75 pair-day).
3. Implement the own-profile endpoint and approved success/error responses (1.0 pair-day).
4. Log authentication failures using the existing safe structured-logging conventions; never log tokens or credentials (0.25 pair-day).
5. Add token-fixture, handler and persistence tests for valid, invalid, unavailable-dependency and repeated-provisioning cases (0.5 pair-day).

**Deliverables**

- API authentication middleware and caller identity context.
- Own-profile endpoint.
- Automated unit and integration tests.

**Acceptance criteria**

1. Requests with valid credentials follow the Approved specification and resolve the local caller ID.
2. Missing/invalid token cases return the specified unauthorized response and do not provision a user.
3. Concurrent/repeated first requests result in one local user record and stable user ID.
4. The caller can read only their own profile; no token, password, or subject identifier is returned or logged.
5. PostgreSQL is authoritative for application roles and permissions; JWT role claims alone never grant privileges.
6. Database/identity-provider failures follow the Approved specification without success-shaped fallback.
7. Tests cover token validation, provisioning, endpoint behavior, and error mapping; CI passes.

**Testing requirements:** Offline token fixtures; database integration tests using the approved model; secret/log-content assertions; race-enabled tests per the project's existing CI convention.

**Definition of Done:** SPEC-F01 Approved before implementation; INF-09 available; agreed CON-05/INF-05 integration; all acceptance evidence, cross-pair PR approval, CI pass, and merge.

**Explicit exclusions:** F01-02 role authorization/admin designation, F01-03 sign-in UI, F06 submissions, role decisions not in the Approved spec, custom identity providers/protocols, the experimental administrator typo mapping, and Resource Owner Password Flow.

**Mentor approval:** No additional DEC-03 or role-authority decision is pending. The reviewer pair alone approves SPEC-F01; remaining specification open questions must be handled under the SDD workflow.

---

### S3-05: [CON-05] Agree the Sprint 3 data-model slice

**Backlog ID:** CON-05 (Sprint 3 schema slice; broader contract remains Preliminary)
**Type / label:** Contract / `type:contract`
**Assigned pair:** D
**Estimate:** 2 SP; **1.0 pair-day**
**Objective:** Agree only the CON-05 entities and fields required by Sprint 3's user/profile and published-problem work; do not imply that unrelated submission/result questions or the entire contract are approved.

**Dependencies and prerequisites**

- S3-01 agreement on the applicable CON-04 problem terms, using the merged PR #47 artifact.
- Apply approved SPEC-F01 OQ-1 / CON-05 C05-OQ-01 (PostgreSQL role authority) and DEC-03 / OQ-2 (Auth0); resolve the remaining OQ-9 schema representation questions needed by the Sprint 3 slice.
- Dependent pairs are named, cross-pair review is complete, and agreement is obtained for the Sprint 3 slice.

**References**

- Product Backlog §4.1 CON-05 (size 4; 3.0 completed preliminary in Sprint 2, 1.0 remaining).
- [CON-05 preliminary](../contracts/CON-05.md), [CON-01](../contracts/CON-01.md), agreed CON-04 terms, SPEC-F01 and SPEC-F08.
- SDD Workflow §2.4 and §4.

**Detailed scope and subtasks**

1. Define only the `users`/`roles` fields and relationships required by approved F01 behavior, respecting CON-01 and the resolved F01 role/IdP questions (0.25 pair-day).
2. Define only the problem fields/status needed by F02-01 from agreed CON-04 terms and approved SPEC-F02 (0.25 pair-day).
3. Keep test-asset, submission and result questions outside this slice explicitly open unless needed by a selected Sprint 3 acceptance criterion; do not infer their schema (0.25 pair-day).
4. Record field constraints and dependent-pair agreement for this slice; retain CON-05's overall Preliminary status if unrelated blocking questions remain (0.25 pair-day).

**Deliverables**

- `docs/contracts/CON-05.md` updated with an explicitly agreed Sprint 3 schema slice.
- Recorded open-question dispositions and pair agreement; the full contract remains Preliminary unless all of its approval conditions are met.

**Acceptance criteria**

1. The Sprint 3 `users`/`roles` and problem slice matches approved SPEC-F01/SPEC-F02, CON-01 and agreed CON-04 terms.
2. Every field/type/constraint needed by S3-04/S3-06/S3-08 is agreed or explicitly blocks that dependent subtask; no assumption is presented as a decision.
3. Unrelated test-asset/submission/result questions remain visible and do not receive inferred values.
4. Every dependent pair agrees to this slice after a non-owner pair review; no full-CON-05 approval is claimed if its remaining conditions are open.

**Testing/review:** Cross-document consistency review; schema implementers and dependent feature pairs record agreement.

**Definition of Done:** Sprint 3 schema slice is reviewed, merged and agreed by its dependent pairs. CON-05 remains Preliminary unless all contract-wide approval requirements are met; INF-05 uses only the agreed slice.

**Explicit exclusions:** database migrations, application code, unapproved schema assumptions, and changing SPEC-F01 or CON-04.

**Mentor approval:** Mentor resolves requirement/decision questions that block finalization; cross-pair agreement remains mandatory and is not replaced by mentor approval.

---

### S3-06: [INF-05] Implement PostgreSQL migrations and data access

**Backlog ID:** INF-05
**Type / label:** Technical task / `type:technical-task`
**Assigned pair:** D
**Estimate:** 5 SP; **4.0 pair-days**
**Objective:** Implement the agreed Sprint 3 schema/migration and PostgreSQL data-access foundation using approved contract slices and the formally approved migration-tool decision.

**Dependencies and prerequisites**

- INF-03 local PostgreSQL environment and INF-02 CI baseline are available.
- The connection configuration and database integration-test harness may be prepared independently, without creating tables, migrations or generated models.
- Migration and schema-dependent access work waits for S3-05 agreement on the applicable Sprint 3 schema slice. ADR-003 / DEC-11 is Approved; preserve all other unresolved schema questions.
- No migration or generated model may encode unresolved CON-05 questions.

**References**

- Product Backlog §4.2 INF-05 (size 4).
- Agreed CON-05 Sprint 3 slice, CON-04, CON-01; [ADR-003](../architecture/adr/ADR-003.md) (Approved).
- S3-04 F01-01 integration; later F02/F06/F08 work consumes the data-access layer.

**Detailed scope and subtasks**

1. Prepare pgx connection configuration and a PostgreSQL integration-test harness that assumes no application tables or schema (1.0 pair-day).
2. Configure only the `pressly/goose` migration tool approved under DEC-11 (0.5 pair-day).
3. Implement and serialize migrations from the agreed Sprint 3 CON-05 slice under one Pair D owner (0.75 pair-day).
4. Provide data-access operations required by the Sprint 3 F01-01 profile and F02-01 published-problem paths, using only agreed interfaces (1.0 pair-day).
5. Add migration/data-access tests and document the local workflow (0.75 pair-day).

**Deliverables**

- Version-controlled migrations derived from agreed contracts.
- pgx/sqlc setup and data-access layer for the Sprint 3 user/profile path.
- Automated migration/data-access tests and concise usage documentation.

**Acceptance criteria**

1. The schema-independent harness connects to the documented local PostgreSQL environment without adding an application schema.
2. Migration tooling matches the formally approved ADR-003 decision; no alternative is introduced.
3. A clean local PostgreSQL instance can apply the migration set after agreement on the applicable CON-05 slice.
4. User/profile and published-problem operations match the agreed Sprint 3 contract slices; provisioning is idempotent as required by SPEC-F01.
5. No schema fields, constraints, or deletion/role policies are invented where the contract remains open.
6. Tests verify persistence and explicit database-failure propagation; CI passes.
7. Migration ownership and execution instructions are documented.

**Testing requirements:** Test connection/configuration independently. After gates pass, apply migrations to a clean database; run user persistence and repeat/concurrent provisioning tests plus published-status repository tests; verify failures propagate explicitly; CI and race-enabled tests pass where configured.

**Definition of Done:** The applicable CON-05 slice is agreed; migration and access behavior meet acceptance criteria; cross-pair PR approval, CI pass, documentation updated, and merge. If slice gates do not clear, only independently completed harness criteria may be evidenced; the INF-05 backlog item remains open.

**Explicit exclusions:** Unapproved schema or migration decisions, catalogue HTTP endpoint F02-02, sample/submission messaging, judge worker, and any modification to contracts/ADR to suit implementation.

**Mentor approval:** DEC-11 / MA-23 is Approved. Contract agreement and PR review remain separate gates.

---

### S3-07: [SPEC-F02] Specify the problem catalogue

**Backlog ID:** SPEC-F02
**Type / label:** Feature / `type:feature`
**Assigned pair:** C
**Estimate:** 3 SP; **1.5 pair-days**
**Objective:** Create a reviewable specification for the Must-priority problem catalogue so its published-problem persistence work can proceed without inferred behavior.

**Dependencies and prerequisites**

- Obtain the mentor-confirmed CON-04 draft required by the backlog; the contract file was not in the planning inventory.
- Use only the Product Backlog F02 requirements and MVP Proposal behavior as scope.
- A reviewer pair other than the responsible Sprint 3 pair must approve the specification before S3-08 starts.

**References**

- Product Backlog §2 F02, §4.1 SPEC-F02 and §4.4 F02-01/F02-02.
- MVP Proposal F02 problem-catalogue behavior; SDD Workflow §§2.1 and 4.
- CON-04 draft and preliminary CON-05 as inputs only; neither is treated as final here.

**Detailed scope and subtasks**

1. Draft `specs/features/F02-problem-catalogue.md` from the existing feature-spec template (0.25 pair-day).
2. Specify only the published problem listing and title/difficulty/solved-state behavior already stated in the backlog; record solved state as unknown until F02-04 (0.5 pair-day).
3. Identify required data/API contract inputs and record unresolved questions instead of selecting fields, sort order, pagination or visibility by assumption (0.5 pair-day).
4. Submit for cross-pair review and address blocking comments; the reviewer pair controls Approved status (0.25 pair-day).

**Deliverables**

- F02 feature specification in the existing `specs/features/` location.
- Review outcome and open questions recorded in the specification.

**Acceptance criteria**

1. All documented F02 Must behavior and the F02-01/F02-02 boundaries are traceable to backlog items.
2. The specification distinguishes published-only listing from future solved-state work.
3. API/data behaviors not specified by the backlog or an agreed contract remain explicit open questions.
4. A reviewer pair other than the author pair records approval before S3-08 can start.

**Testing requirements:** Specification checklist and cross-document consistency review against the MVP Proposal, Product Backlog and available contracts.

**Definition of Done:** Specification has cross-pair review, no blocking open questions, and reviewer-pair Approved status. A draft or merged document alone is not approval.

**Explicit exclusions:** F02 implementation, selecting undocumented API semantics, problem details F03, user-specific solved-state implementation F02-04, and changing CON-04/CON-05.

**Mentor approval:** Mentor resolves requirement questions and confirms any scope change; the reviewer pair alone approves the specification.

---

### S3-08: [F02-01] Implement published-problem persistence

**Backlog ID:** F02-01
**Type / label:** Feature / `type:feature`
**Assigned pair:** C
**Estimate:** 5 SP; **4.0 pair-days**
**Objective:** Persist problem records and provide the published-status filter defined by the Approved F02 specification and agreed Sprint 3 data-model slice.

**Dependencies and prerequisites**

- S3-07 SPEC-F02 is Approved before feature implementation.
- The CON-05 Sprint 3 slice is agreed, including problem fields/status semantics derived from CON-04.
- INF-05 supplies the approved migration/data-access path. Pair C may prepare interface-level tests with fixtures/test doubles, but must not define schema fields or migrations.

**References**

- Product Backlog §4.4 F02-01 (size 4); §4.1 SPEC-F02 and CON-05.
- Approved SPEC-F02; agreed CON-04 problem terms and CON-05 Sprint 3 slice; S3-06 INF-05.
- SDD Workflow §§2.1 and 4.

**Detailed scope and subtasks**

1. Define behavior only from Approved SPEC-F02 and agreed repository/data interfaces (0.5 pair-day).
2. Implement the F02-01 problem repository operations, including the published-status filter, using Pair D's agreed data-access layer (1.5 pair-days).
3. Use fixtures/test doubles for unit-level behavior while the schema adapter is in progress; integrate only after CON-05/INF-05 interfaces are agreed (0.75 pair-day).
4. Add focused tests for published/unpublished records, empty results and repository/database errors (1.25 pair-days).

**Deliverables**

- Problem persistence and published-only repository query.
- Automated unit and database integration tests against the agreed schema.

**Acceptance criteria**

1. Published records are returned and unpublished records are excluded according to Approved SPEC-F02.
2. Repository fields, constraints and query behavior match the agreed CON-05 Sprint 3 slice; no schema is inferred from preliminary drafts.
3. Empty results and database failures follow the approved specification without a success-shaped fallback.
4. Tests cover published/unpublished behavior, errors and the agreed database integration path.
5. CI passes and another pair reviews the implementation.

**Testing requirements:** Unit tests with fixtures/test doubles; integration tests against the approved PostgreSQL schema after S3-06 is ready.

**Definition of Done:** SPEC-F02 Approved before implementation; agreed CON-05/INF-05 interface used; acceptance evidence, cross-pair review and CI are complete; changes are merged.

**Explicit exclusions:** Catalogue HTTP endpoint F02-02, catalogue UI F02-03, solved-state F02-04, problem details F03, and any unapproved data/API shape.

**Mentor approval:** No separate decision is presumed. Mentor resolves any requirement or scope question; formal specification and contract gates remain mandatory.

---

### S3-09: [INF-12a] Build API and orchestrator container images

**Backlog ID:** INF-12a
**Type / label:** Technical task / `type:technical-task`
**Assigned pair:** B
**Estimate:** 2 SP; **1.5 pair-days**
**Objective:** Add reproducible API and orchestrator container-image build/checks specified by the backlog, using the existing CI baseline.

**Dependencies and prerequisites**

- INF-02 CI baseline is available.
- This task does not depend on database schema, identity-provider decisions or judge-execution approval.

**References**

- Product Backlog §4.2 INF-12a (size 1.5; API and orchestrator images with dependency/image security checks in CI).
- Existing CI and repository contribution conventions; SDD Workflow §2.2.

**Detailed scope and subtasks**

1. Define the API and orchestrator image build inputs required by INF-12a (0.5 pair-day).
2. Connect image builds and existing dependency/image security checks to CI (0.75 pair-day).
3. Document the existing CI build/checks; do not add a deployment or runtime platform (0.25 pair-day).

**Deliverables**

- Reproducible API and orchestrator image builds.
- CI image/dependency security checks and concise build instructions.

**Acceptance criteria**

1. Both images build reproducibly using the existing project build conventions.
2. Existing dependency and image-security checks run in CI and report failures explicitly.
3. No credentials/secrets are embedded in image layers or build output.
4. No new hosting, registry, orchestrator behavior or infrastructure component is introduced.

**Testing requirements:** Build both images and run the existing configured image/dependency checks; verify CI evidence.

**Definition of Done:** Acceptance criteria are evidenced, CI passes, another pair reviews the PR, and changes are merged.

**Explicit exclusions:** INF-12b judge-runner image, staging/deployment, Kubernetes changes, and changes to application behavior.

**Mentor approval:** Not required unless the implementation introduces a new technology, hosting target or architecture decision.

---

### S3-10: [TST-01] Create reference challenge fixtures (optional stretch)

**Backlog ID:** TST-01
**Type / label:** Technical task / `type:technical-task`
**Assigned pair:** Open (mentor to confirm after availability review)
**Estimate:** 5 SP; **4.0 pair-days**
**Objective:** Create the two trusted challenge fixtures in the existing backlog so later feature and test work has repeatable positive and negative examples.

**Dependencies and prerequisites**

- Optional only; requires confirmed additional effective capacity beyond the target commitment.
- Use the merged CON-04 artifact and approved DEC-09 choices; dependent-pair agreement on applicable terms is required before writing fixtures.
- Fixtures are test assets only; no learner-submitted code is executed.

**References**

- Product Backlog §4.5 TST-01 (size 4; one algorithmic and one HTTP-handler challenge, reference solutions and negative cases).
- Agreed CON-04 terms and SDD Workflow §2.2.

**Detailed scope and subtasks**

1. Create one algorithmic challenge fixture using only the agreed CON-04 packaging format (1.25 pair-days).
2. Create one HTTP-handler challenge fixture in the same agreed format (1.25 pair-days).
3. Include reference solutions and negative cases for both fixtures as required by the backlog (1.0 pair-day).
4. Validate the fixture files against the agreed format and document how later automated tests consume them (0.5 pair-day).

**Deliverables**

- Two version-controlled reference challenge fixtures with negative cases.
- Validation/tests and concise fixture-use instructions.

**Acceptance criteria**

1. Both fixtures conform to the agreed CON-04 terms and are usable by later automated tests.
2. Each has a reference solution and documented negative cases.
3. No hidden test is exposed through learner-facing output; fixture content is not used to approve DEC-02/DEC-08.
4. No learner code is executed and no execution infrastructure is added.
5. Cross-pair review and CI validation pass.

**Testing requirements:** Validate fixture structure and run only trusted fixture-validation tests; no judge execution.

**Definition of Done:** Merged CON-04 artifact, agreement on applicable terms, capacity confirmation, acceptance evidence, cross-pair review and CI pass are recorded. Otherwise keep this optional issue unstarted.

**Explicit exclusions:** Judge execution, hidden-test service, F07 implementation, infrastructure, and approval of isolation decisions.

**Mentor approval:** Confirm optional capacity and pair assignment; dependent-pair contract agreement remains required.

---

## 4. Non-Selected INF-07 Remainder

The Product Backlog carries a proposed 1.0 pair-day estimate for basic metrics only. ST-04 completed HTTP bootstrap/configuration, liveness/readiness endpoints, request-ID/access-log middleware, timeouts, graceful shutdown, tests and documentation; none is allocated or estimated again. Basic metrics remain unselected because the existing backlog does not define the exact acceptance surface. Confirm scope before scheduling metrics later; do not assume a metrics technology.

## 5. Capacity and Dependency Validation

| Check | Result |
|-------|--------|
| Calendar gross capacity | **40 pair-days** (4 pairs × 10 working days) |
| Documented effective task capacity | **20 pair-days** |
| Target committed workload / overage | **23.5 / 3.5 pair-days** |
| Optional stretch / total target plus stretch | **4.0 / 27.5 pair-days** (7.5 over effective capacity) |
| Capacity-safe alternative / reserve | **18.0 / 2.0 pair-days** |
| Target pair workload | A 5.5, B 5.5, C 7.5, D 5.0 pair-days |
| Target coding/integration/testing | **19.0 pair-days (80.9%)**, above the 80% target |
| Capacity-safe alternative coding/integration/testing | **15.0 pair-days (83.3%)** |
| Coding task for every pair | Yes for the approved target: editor, API identity, problem persistence and database work. If the fallback defers S3-08, replan Pair C's coding work with the mentor. |
| Completed Sprint 1/Sprint 2 scope duplicated | No; ST-02–ST-04, ST-01/INF-09 and completed ST-04 server/health work are not rescheduled |
| INF-07 remaining scope | Basic metrics only, proposed 1.0 pair-day; not committed. All listed ST-04 activities are excluded |
| Approval gates | Explicitly listed; implementation starts only after feature/contract/ADR gates |
| Parallel work | Auth middleware uses repository test doubles; editor component uses approved fixtures; DB connection/test harness starts without schema assumptions |
| Remaining blocking dependencies | CON-04 PR #47 merge and dependent-pair agreement; CON-05 Sprint 3 slice agreement; SPEC-F01/F02/F04 reviewer approval; SPEC-F02 approval by a pair other than C |
| End-to-end user journey | Not claimed; catalogue endpoint F02-02, submissions and judge execution are outside this sprint |

## 6. Final Approval Note

This manifest is a proposal for mentor review. It does not create GitHub Issues, approve specifications/contracts/ADRs, authorize architectural decisions, or authorize implementation. Wait for mentor review and the specific SDD gates before creating issues or beginning work.
