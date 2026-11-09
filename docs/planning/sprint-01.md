# Sprint 1 Plan: 2–13 November 2026

**Status:** Approved for Sprint 1 execution (2 November 2026), on the requester's instruction, after the SDD documents were made consistent for the four work item types. No GitHub Issues have been created from this plan. The mentor's Day-1 decisions are listed in §8.1 and still need recording.

**Stretch added (9 November 2026):** the mentor approved two stretch items for the second week of Sprint 1, CON-02 (Pair C) and SPEC-F04 (Pair D), totalling 3.0 pair-days. They are recorded in §4.2 and are not part of the committed scope.

**Sources:** [Product Backlog](../product/product-backlog.md) (revision 3, §5 spikes, §9 Sprint 1, §11.3 overhead), [MVP Proposal](../product/mvp-proposal.md), [SDD Workflow](../sdd-workflow.md), [Contributing Guide](../../CONTRIBUTING.md), [Feature Specification Template](../../specs/templates/feature-spec-template.md).

**Relationship to the backlog:** this plan supersedes backlog §9 for Sprint 1 only. The backlog is not modified. Where this plan differs from backlog §9, the differences are listed in §4.

---

## 1. Sprint Goal

By 13 November, the team has a working local platform baseline, a go or no-go answer on running Go tests in isolated Kubernetes Jobs, evidence on Auth0 and RabbitMQ behaviour, and an agreed submission lifecycle and first specification ready for cross-pair review.

The secure execution spike is the Sprint 1 priority. It is the only item on the critical path.

---

## 2. Capacity and Committed Effort

### 2.1 Capacity

| Item | Pair-days |
|------|-----------|
| Gross capacity: 4 pairs × 10 working days | 40 |
| Overhead: training 6, meetings 4, specification review 3, pull request review 1.5, integration 1 (backlog §11.3) | 15.5 |
| **Realistic net capacity** (backlog §11.3) | **24.5** |

### 2.2 Committed effort

| Measure | Pair-days |
|---------|-----------|
| **Committed implementation effort** | **24.5** |
| Stretch effort, mentor-approved 9 November, recorded separately (§4.2) | 3.0 |
| Headroom, from removed double-counting (§3) | 2.0 |

Committed effort equals realistic net. This is intentional, and it has two safeguards:

- Spike timeboxes are fixed. They end on day 10 whatever their state, and they are not extended without MA-12.
- The 2.0 pair-days of headroom in §3 are held per pair (§5.1). They are not banked as committed work.

The backlog's previous Sprint 1 commitment was 28.5 pair-days. That is 4 pair-days over realistic capacity, and this plan removes that over-commitment.

**Stretch adjustment (9 November).** The mentor observed that Pairs C and D completed their originally assigned work earlier than expected, and allocated 3.0 pair-days of additional work for the second week. The adjustment rests on that observed early completion. It does not change the committed scope of 24.5 pair-days, the capacity estimates in §2.1, or the per-pair budgets in §5.1. Committed and stretch effort are reported separately, and the total planned for Sprint 1 is 27.5 pair-days.

---

## 3. Double-Counting Review

Backlog estimates and overhead estimates can count the same effort twice. The overlaps below were checked. Each is resolved as stated.

| # | Overlap | Where it is counted | Resolution |
|---|---------|---------------------|------------|
| 1 | Training (6 pair-days in §11.3) includes Kubernetes, RabbitMQ and Auth0 onboarding. SPK-01, SPK-02 and SPK-03 already include learning those tools. | Overhead, and inside spike timeboxes | Training is set to 4 pair-days: Pair A 1.5, Pair B 0.5, Pair C 0.5, Pair D 1.5. Pairs B and C do not take general training, because their spike learning covers it. **Assumption:** MA-18. |
| 2 | Cluster and broker setup. SPK-01 could be read as including an "interim cluster". INF-03 also builds a cluster. | INF-03 and SPK-01 | Setup is counted once, in INF-03 (Pair A). SPK-01 and SPK-03 start only after the milestones in §6. Spikes do not build environments. |
| 3 | CI debugging. INF-02 (1.5) and integration overhead (1) both cover pipeline fixes. | INF-02 and overhead | CI fixes are counted in INF-02. Integration overhead covers only cluster and environment support outside INF-03. |
| 4 | Spike results shown at the demo. Meetings include demo preparation, and the spike reports are written inside the spike timebox. | Meetings and spike timeboxes | No separate slides. The demo uses the spike report. Demo preparation is counted in meetings only. |
| 5 | Specification review. SPEC-F01 and CON-01 are authored by Pair D, and reviewed by Pairs A and C. | Author's item and review overhead | Authoring is counted in the item (Pair D). Review is counted once, in review overhead (Pairs A and C). Reviewers must not rewrite the specifications. |
| 6 | Pull request review of spike reports. Spike reports are merged through pull requests. | Spike timebox and PR review | Report writing is inside the spike timebox. Review of the report PR is counted once, in PR review overhead. |
| 7 | Pair B's agreement on CON-01. Agreement is a dependent-pair activity. | CON-01 (Pair D) and Pair B headroom | Counted in Pair B's headroom (0.25). It is not added to CON-01. |

**Net effect:** the training change releases 2.0 pair-days. That is the headroom in §2.2. It is not used for scope.

---

## 4. Commitments, Stretch and Deferred Work

### 4.1 Committed (24.5 pair-days)

| Item | Pair | Pair-days | Reason for Sprint 1 |
|------|------|-----------|---------------------|
| SPK-01 Secure Go execution with `go test` and Kubernetes Jobs | B | 8 | Priority spike; critical path |
| SPK-02 Auth0 OIDC integration check | C | 2 | Gates DEC-03 and INF-09 |
| SPK-03 RabbitMQ and Watermill check | C | 4 | Gates DEC-04, CON-02 and DEC-12 |
| INF-03 Local infrastructure | A | 4 | Needed by SPK-01 and SPK-03 |
| INF-01 Go module, layout, lint and format | D | 1.5 | Needed by INF-02 and all code |
| INF-02 CI baseline | D | 1.5 | Gates the Definition of Done for every later item |
| INF-04 ADR process and first records | D | 0.5 | Records DEC-14 and DEC-06 |
| CON-01 Submission lifecycle and verdict model | D | 1.5 | Gates SPEC-F06, SPEC-F07 and CON-02 |
| SPEC-F01 Accounts and access specification | D | 1.5 | Gates F01-01 in Sprint 2 |
| **Total** | | **24.5** | |

### 4.2 Stretch (mentor-approved, 9 November 2026)

Stretch items are not committed. They are outside the 24.5 pair-day commitment in §4.1 and §2.2. They are worked in the second week of Sprint 1 (9–13 November). They are not Sprint 1 exit criteria (§11).

| Item | Pair | Pair-days | Dependencies | Deliverable | Acceptance conditions |
|------|------|-----------|--------------|-------------|-----------------------|
| CON-02 Messaging contract (S1-10) | C | 1.5 | CON-01 and SPK-03 (§4.1). Drafting starts from their approved PRs. Final approval waits for both PRs to merge. DEC-04 is provisional. DEC-12 is not decided by this item. | Contract at `docs/contracts/CON-02.md` | Queue names, message schema, idempotency key, acknowledgement, retry and dead-letter behaviour stated. Dependent pairs named. Reviewed by Pair D. Agreement recorded by each dependent pair. Not merged before agreement. |
| SPEC-F04 Browser code editor specification (S1-11) | D | 1.5 | None (backlog §4.1). DEC-13 (session preservation) stays Open until the mechanism is approved. It is recorded in the spec, not decided here. The CON-04 draft is an input to the template format. | Specification at `specs/features/F04-browser-code-editor.md`, status In Review | Written from the template. Acceptance criteria numbered AC-01 onward. Covers Go source editing, starter templates, active-session preservation and edge cases. Reviewed by Pair A. Status In Review (`spec:in-review`). Not approved in Sprint 1. |

**Rules for stretch work:**

- Stretch does not use the headroom in §5.1. It is funded by the early completion described in §2.2.
- Stretch items are authorised to start immediately (mentor, 9 November). Drafting is not blocked by unmerged PRs.
- CON-02 drafting uses the approved SPK-03 and CON-01 PRs. Contract approval waits for both PRs to merge and for the architecture decisions it depends on to be confirmed (§8.4). Drafting is not approval.
- Stretch items do not change the Sprint Goal (§1) and are not exit criteria (§11).
- An item counts as complete only when its deliverable is merged and its acceptance conditions are met. An item not finished by 13 November 15:00 is carried over and is not reported as completed.

**Scheduling.** Backlog §10 places CON-02 in Sprint 2 and SPEC-F04 in Sprint 3. Those rows are unchanged. Whether the two items are removed from Sprints 2 and 3 is decided at Sprint 2 planning.

### 4.3 Deferred from backlog §9

| Item | Backlog §9 | This plan | Effect |
|------|------------|-----------|--------|
| INF-08 Frontend scaffold | Sprint 1 (1.5) | Sprint 3 | Already moved in backlog revision 3. Not needed until F01-03 in Sprint 4. |
| CON-04 Problem packaging, draft | Sprint 1 (2) | Sprint 2 (draft and final, 4 in total) | DEC-09 product input is needed earlier in Sprint 2. Sprint 2 load rises. |
| SPEC-F02 Problem catalogue specification, draft | Stretch (1.5) | Sprint 2 | SPEC-F02 approval moves to about the end of Sprint 2 week 1. F02-01 may slip into Sprint 2 week 2. |

### 4.4 Ownership changes from backlog §9

Backlog §9 gave Pair A INF-01, INF-02, INF-03, INF-04 and INF-08. This plan gives Pair D INF-01, INF-02 and INF-04, and Pair A only INF-03. Reason: Pair A carries 4 pair-days of INF-03, and a full allocation of the infrastructure items to Pair A would exceed its net capacity (§5.1).

---

### 4.5 Work item types

Each Sprint 1 issue is one of the four work item types in [CONTRIBUTING §2](../../CONTRIBUTING.md#2-sdd-rules-and-work-item-types). Estimates, owners and dependencies are unchanged from §4.1 and §4.2.

| Work item | Type | Required before work | Artefact | Branch prefix |
|-----------|------|----------------------|----------|---------------|
| SPEC-F01 (S1-09) | Feature | Approved Feature Specification | `specs/features/F01-accounts-and-access.md` | `docs/` for the specification |
| CON-01 (S1-08) | Contract | Cross-pair review and dependent-pair agreement | `docs/contracts/CON-01.md` | `docs/` |
| INF-04 (S1-07) | Technical task (documentation or ADR) | Acceptance criteria | `docs/architecture/adr/ADR-001.md`, `ADR-002.md` and the template | `docs/` |
| INF-01, INF-02, INF-03 (S1-05, S1-06, S1-01) | Technical task (infrastructure) | Acceptance criteria | Repository changes and README | `infra/` |
| SPK-01, SPK-02, SPK-03 (S1-02, S1-03, S1-04) | Technical spike | Spike Brief on the issue | Report at `docs/spikes/SPK-01.md`, `SPK-02.md`, `SPK-03.md`. Experiment code on `spike/`, not merged. | `spike/` for code. `docs/` for reports. |
| CON-02 (S1-10, stretch) | Contract | Draft contract, cross-pair review and dependent-pair agreement | `docs/contracts/CON-02.md` | `docs/` |
| SPEC-F04 (S1-11, stretch) | Feature | Approved Feature Specification (approval is not a Sprint 1 outcome) | `specs/features/F04-browser-code-editor.md` | `docs/` for the specification |

### 4.6 Status of committed and stretch work (9 November 2026)

Three states are kept separate for every item:

- **Implementation or documentation completed:** the artefact is written and the owner pair has finished its work.
- **Pull request state:** open, approved or merged. CONTRIBUTING §9 requires at least one approval from a reviewer from another pair before merge. Owner approval alone does not meet this rule.
- **Formally done:** only after the merge, when the Definition of Done in CONTRIBUTING §11 is met. The issue is then labelled `done` and closed.

Merged state is read from `main`. Completion and PR states for SPK-03, CON-01 and SPEC-F01 are as reported by the mentor on 9 November. The repository does not show those PRs, and this plan has not verified them.

| Item | Category | Pair-days | Implementation or documentation | Pull request | Formally done |
|------|----------|-----------|---------------------------------|--------------|---------------|
| INF-01 | Committed | 1.5 | Completed | Merged (PR #12) | Merged. Issue closure not checked here. |
| INF-02 | Committed | 1.5 | Completed | Merged (PR #14) | Merged. Issue closure and green `main` not checked here. |
| INF-03 | Committed | 4 | Completed | Merged (PRs #15, #16, #18) | Merged. Issue closure not checked here. M1 to M3 not re-verified here. |
| INF-04 | Committed | 0.5 | Completed. ADR-001 and ADR-002 present, both Provisional. | Merged (PR #11) | Merged. Issue closure not checked here. |
| SPK-02 | Committed | 2 | Completed. Report at `docs/spikes/SPK-02.md`. | Merged (PR #13) | Merged. Issue closure not checked here. |
| SPK-03 | Committed | 4 | Completed by Pair C, as reported by the mentor. Report not on `main`. | Open. Owner approval reported. Cross-pair approval outstanding. | Not done. Awaits cross-pair approval and merge. |
| CON-01 | Committed | 1.5 | Completed by Pair D, as reported by the mentor. Contract not on `main`. | Open. Owner approval reported. Cross-pair approval and agreement from Pairs A, B and C outstanding. | Not done. Awaits approval, agreement and merge. |
| SPEC-F01 | Committed | 1.5 | Completed by Pair D, as reported by the mentor. Specification not on `main`. | Open. Owner approval reported. Reviews from Pairs A and C outstanding. | Not done. Awaits review, approval and merge. Spec approval (`spec:approved`) is a Sprint 2 start event. |
| SPK-01 | Committed | 8 | Not confirmed. No report in the repository, and the mentor's confirmation does not cover it. | None | Not done. |
| CON-02 | Stretch | 1.5 | Drafting authorised 9 November. No draft in the repository. | None | Not done. |
| SPEC-F04 | Stretch | 1.5 | Drafting authorised 9 November. No draft in the repository. | None | Not done. |

**Note.** Stretch drafting may start now (§4.2). Drafting is not approval. CON-02 is not approved until the SPK-03 and CON-01 PRs are merged and the decisions in §8.4 are confirmed.

---

## 5. Four Parallel Workstreams

Pair labels A–D are placeholders. The mentor assigns named interns to pairs at the start of the sprint. Named interns are not assigned in this plan.

### 5.1 Pair budgets

Each pair has 10 pair-days. Overhead is split by role.

| Pair | Workstream | Committed | Baseline overhead (training, meetings, integration) | Review overhead | Used | Headroom |
|------|------------|-----------|------------------------------------------------------|-----------------|------|----------|
| A | WS-1 Local platform | 4.0 (INF-03) | 2.75 | 2.5 (SPEC-F01 review 1.0, CON-01 review 0.5, PR review 1.0) | 9.25 | **0.75** |
| B | WS-2 Secure execution | 8.0 (SPK-01) | 1.75 | 0 | 9.75 | **0.25** |
| C | WS-3 Identity and messaging | 6.0 (SPK-02, SPK-03) | 1.75 | 2.0 (SPEC-F01 review 0.5, CON-01 review 1.0, PR review 0.5) | 9.75 | **0.25** |
| D | WS-4 Foundations and specifications | 6.5 (INF-01, INF-02, INF-04, CON-01, SPEC-F01) | 2.75 | 0 | 9.25 | **0.75** |
| **Total** | | **24.5** | **9.0** | **4.5** | **38.0** | **2.0** |

Pair B has no review load, so its 0.25 headroom is the least. Pair B's agreement on CON-01 (§3, item 7) uses that headroom.

**Stretch load is not in this table.** CON-02 (Pair C, 1.5) and SPEC-F04 (Pair D, 1.5) are stretch under §4.2. Each exceeds the baseline headroom of its pair (0.25 for Pair C, 0.75 for Pair D). They rely on the time freed by the early completion in §2.2, and this table is not changed to show them. Their review load is also not budgeted here: CON-02 review (Pair D) and SPEC-F04 review (Pair A) are not in the review overhead column.

### 5.2 WS-1 Local platform (Pair A)

**Owns:** INF-03.

**Measurable deliverables:**

- **M1, end of day 2:** Local cluster running, using the tool chosen under DEC-06. `kubectl`-level access works from the developer machine. Pair B can create a Job in it.
- **M2, end of day 3:** RabbitMQ and PostgreSQL running and reachable from the local cluster and from the developer machine. Pair C can publish and consume on RabbitMQ.
- **M3, end of day 4:** Setup steps in the README. Another pair can reproduce the environment from the README alone.

**Reviews by Pair A:** SPEC-F01 (1.0 pair-day) and CON-01 (0.5). PR review of spike reports (1.0).

### 5.3 WS-2 Secure execution (Pair B)

**Owns:** SPK-01, the priority spike. See §6 for its daily plan.

**Measurable deliverables:**

- Spike report merged to `docs/spikes/SPK-01.md` by end of day 10, with SC-01 to SC-10 each recorded as Verified, Partially verified or Not verified, using §7.1.
- Negative-case programs for the spike, held in a spike branch. Not merged to `main`.

### 5.4 WS-3 Identity and messaging (Pair C)

**Owns:** SPK-02 (days 1–2) and SPK-03 (days 3–6).

**Measurable deliverables:**

- SPK-02 report by end of day 2, with SC-A to SC-C recorded per §7.2.
- SPK-03 report by end of day 6, with SC-A to SC-D recorded per §7.3. Days 7–10 are available for review and for agreeing CON-02 input.

**Reviews by Pair C:** SPEC-F01 (0.5) and CON-01 (1.0). PR review (0.5).

**Stretch (§4.2):** CON-02 authoring (1.5 pair-days). Drafting starts now, from the approved SPK-03 and CON-01 PRs. Final approval waits for both to merge.

### 5.5 WS-4 Foundations and specifications (Pair D)

**Owns:** INF-01, INF-02, INF-04, CON-01, SPEC-F01.

**Measurable deliverables:**

- **Day 1:** ADR-001 (`docs/architecture/adr/ADR-001.md`) records the Go version pin (DEC-14) as provisional. INF-01 starts.
- **Day 3:** Repository builds and passes lint locally. INF-02 CI pipeline runs on push and pull request. Green on `main` by day 5.
- **Day 4:** ADR-002 (`docs/architecture/adr/ADR-002.md`) records the local cluster choice (DEC-06) as provisional. ADR template merged (INF-04).
- **Day 5:** CON-01 draft (`docs/contracts/CON-01.md`) circulated to Pairs A, B and C for agreement. Agreement recorded by day 8.
- **Day 8:** SPEC-F01 in cross-pair review, reviewed by Pairs A and C.

Pair D's headroom (0.75) is held for review responses and for the mentor's decisions in §8.

**Stretch (§4.2):** SPEC-F04 authoring (1.5 pair-days), drafting from 9 November, reviewed by Pair A. Pair D also reviews CON-02, Pair C's stretch item.

---

## 6. Secure Execution Spike: Daily Plan

The spike is the Sprint 1 priority. It has the most dependencies and the longest critical path.

| Day | Date | Pair B (SPK-01) | Pair A (INF-03) | Gate |
|-----|------|-----------------|-----------------|------|
| 1 | Mon 2 Nov | Read spike criteria (§7.1) and the proposal §7 controls. Write the trusted `go test` package and a client-go program skeleton. No cluster yet. | Start INF-03. Choose tool under DEC-06 (provisional). | Day-1 decisions (§8) recorded |
| 2 | Tue 3 Nov | Write negative-case programs (network attempt, filesystem write, infinite loop, memory exhaustion, output flood). Still no cluster. | Cluster running (M1). | **M1: cluster ready** |
| 3 | Wed 4 Nov | Start SC-01 and SC-06 in the cluster: create a Job that runs the trusted package. | RabbitMQ and PostgreSQL (M2). | M2 |
| 4 | Thu 5 Nov | SC-02 (time limit), SC-03 (memory limit). | README steps (M3). | — |
| 5 | Fri 6 Nov | SC-04 (network), SC-05 (filesystem). **Demo: interim results.** | Support Pair B if the cluster needs it. | Friday demo |
| 6 | Mon 9 Nov | SC-07 (output bounded), SC-08 (logs before cleanup). | Review spike PRs. | — |
| 7 | Tue 10 Nov | SC-09 (TTL cleanup), SC-10 (cancellation). | Review spike PRs. | — |
| 8 | Wed 11 Nov | Re-run any criterion that is not yet 10/10 or 5/5. Write the report. | Respond to review comments. | — |
| 9 | Thu 12 Nov | Finish the report. Record outcome per §7.1. | Respond to review comments. | — |
| 10 | Fri 13 Nov | **Timebox ends at 15:00.** Merge the report. Sprint demo. | Sprint demo. | Sprint demo |

Pair B has 8 pair-days of spike work across days 1–10, with 1.75 pair-days of overhead and 0.25 headroom (§5.1).

**If the cluster is late:** M1 is a hard dependency of SC-01. If M1 slips past day 2, Pair B runs the negative-case programs against a single-node local cluster stood up by Pair A, and records the slip in the report. Pair B does not build its own cluster. That would double-count INF-03 (§3, item 2).

---

## 7. Spike Timeboxes, Success and Failure Criteria

Every spike has a timebox, success criteria and failure criteria. A criterion is **Verified** only if the measure in the backlog (§5) is met. A criterion that is not measured is **Not verified**. Nothing is assumed.

General rules:

- The timebox ends at 15:00 on the last day. The pair reports the outcome reached at that point.
- A spike is not extended without mentor approval (MA-12).
- A failure criterion met on any day is reported at once to the mentor. It is not held until the report.

### 7.1 SPK-01: Secure Go execution (8 pair-days, days 1–10)

Success criteria: SC-01 to SC-10 from [backlog §5.1](../product/product-backlog.md#51-spk-01-secure-go-execution-with-go-test-and-kubernetes-jobs-8-pair-days). Run only trusted packages.

**Outcome categories** (from the backlog's outcome rule, restated):

| Outcome | Condition |
|---------|-----------|
| **Go** | SC-01 to SC-10 all Verified |
| **Conditional go** | SC-01 and SC-08 Verified. Every other criterion Not verified or Partially verified is a named risk for SPEC-F07 and MA-03. |
| **No-go** | Any failure criterion below is met |
| **Inconclusive** | Timebox ends with SC-01 Verified but SC-02 to SC-10 not measured. Reported as Inconclusive; no DEC-02 approval. |

**Failure criteria** (any one is a no-go):

- **F-01:** A Job cannot run a trusted `go test` package with a retrieved result (SC-01 fails).
- **F-02:** A Job that exceeds its deadline is not terminated (SC-02 fails).
- **F-03:** An outbound network connection from the Job succeeds in any run (SC-04 fails).
- **F-04:** A write outside the allowed path succeeds in any run (SC-05 fails).
- **F-05:** Logs or reports are lost before cleanup in any run (SC-08 fails).

Failure consequence: DEC-02 is not approved. The mentor is asked to choose between a different execution model, and option E in backlog §11.5 (internal-only release). No F07 work starts.

### 7.2 SPK-02: Auth0 OIDC integration (2 pair-days, days 1–2)

Success criteria: SC-A to SC-C from [backlog §5.2](../product/product-backlog.md#52-spk-02-auth0-oidc-integration-2-pair-days).

**Failure criteria:**

- **F-06:** The Go middleware accepts a tampered or expired token (SC-B fails). This is a no-go for F01-01 as designed.
- **F-07:** Neither the Auth0 development tenant nor a documented fixture is available by day 2. The spike is Inconclusive, and MA-05 is escalated.

**Outcome:** Go, Conditional go (SC-A on a fixture only), or No-go. DEC-03 is recorded as approved only after MA-05 (§8).

### 7.3 SPK-03: RabbitMQ and Watermill (4 pair-days, days 3–6)

Success criteria: SC-A to SC-D from [backlog §5.3](../product/product-backlog.md#53-spk-03-rabbitmq-and-watermill-4-pair-days).

**Failure criteria:**

- **F-08:** Any message is lost in SC-A (not 100 of 100).
- **F-09:** A redelivered duplicate produces two final states in SC-C (not 50 of 50).
- **F-10:** A poison message cannot reach the dead-letter route in SC-D (not 3 of 3).

**Outcome:** Go, Conditional go, or No-go. DEC-04 is recorded as approved only after the report. DEC-12 (dispatch consistency) is decided after this report, and no later than the start of Sprint 3.

---

## 8. Decisions

### 8.1 Required on Day 1 (mentor)

| ID | Decision | Why now | Recommended outcome |
|----|----------|---------|---------------------|
| MA-02 | Confirm the provisional choices for DEC-14 (Go version pin) and DEC-06 (k3d or Minikube) | INF-01 and INF-03 cannot start without them | Confirm both as provisional. Reversible. |
| MA-13 | Confirm DEC-01 to DEC-04 are the right hypotheses to test | SPK-01 to SPK-03 test these | Confirm. They stay Provisional until the reports. |
| MA-15 | Accept the Sprint 1 plan at 24.5 pair-days, replacing the backlog's 28.5 | Sprint 1 cannot start at 28.5 | Accept §2 and §4 of this plan. |
| MA-05 | Confirm Auth0 tenant access for SPK-02, or accept a documented fixture | SPK-02 starts on day 1 | If access is not in place, accept the fixture route with F-07 as the fallback. |
| MA-12 | Confirm the spike timebox rule (no extension without approval) | Sets how each spike ends | Confirm. |
| MA-17 | **Resolved.** Spikes need a Spike Brief, not a Feature Specification (approved work item rules, CONTRIBUTING §2). | Previously, SDD had no spike rule and SPK-01 was blocked. | Spike experiment code stays on a `spike/` branch and is not merged to `main`. |
| MA-11 | Confirm all four pairs are present for all 10 working days | The plan assumes it | Confirm, or give the absence and the plan will be re-cut. |
| MA-18 | **New.** Confirm pairs B and C need no general training beyond their spikes (§3, item 1) | Basis for the 2.0 headroom | Confirm. If not, training returns to 6 and headroom falls to 0. |

### 8.2 Decisions that follow spike results

| ID | Decision | Input | Needed by |
|----|----------|-------|-----------|
| DEC-01 | Submission evaluation through the Go testing framework | SPK-01 report | Sprint 1 demo (13 Nov) or the first day of Sprint 2 |
| DEC-02 | Execution with Kubernetes Jobs and client-go | SPK-01 report | Same |
| DEC-03 | Authentication with Auth0 | SPK-02 report and MA-05 | Start of Sprint 2 (INF-09) |
| DEC-04 | Messaging with RabbitMQ and Watermill | SPK-03 report | Start of Sprint 2 (CON-02) |
| DEC-12 | Dispatch consistency between database write and publish | SPK-03 report | Start of Sprint 3 (30 Nov) |
| DEC-08 | Isolation mechanism | SPK-01 and SPK-05 | Start of Sprint 4. Not a Sprint 1 decision. |

### 8.3 Decisions that are not spike-dependent

| ID | Decision | Needed by | Input |
|----|----------|-----------|-------|
| DEC-09 | Problem packaging and test visibility | Start of Sprint 2 (MA-14) | Product input from the mentor |
| DEC-11 | Database migration tool | Sprint 2 | Engineering ADR |
| DEC-13 | Editor session preservation | Sprint 3 | Recorded in SPEC-F04 |

### 8.4 Mentor decisions on stretch work

Recorded for the 9 November stretch authorisation. None of these changes the committed scope or the approved 24.5 pair-day Sprint 1 Plan (MA-15).

| # | Decision | Status | Needed by | Why |
|---|----------|--------|-----------|-----|
| 1 | Stretch authorised: CON-02 (Pair C) and SPEC-F04 (Pair D), 3.0 pair-days, outside the committed 24.5. Drafting may start immediately. | Resolved (mentor, 9 Nov) | — | §2.2, §4.2 |
| 2 | Reviewers: CON-02 reviewed by Pair D, with agreement from every dependent pair. SPEC-F04 reviewed by Pair A. | Resolved for the reviewers (mentor, 9 Nov). Dependent pairs open (item 6). | — | §4.2 |
| 3 | Cross-pair approval on the SPK-03, CON-01 and SPEC-F01 PRs. The mentor confirmed owner approval only. CONTRIBUTING §9 requires an approval from another pair before merge. | Open | Before each PR is merged | CONTRIBUTING §5.2, §9 |
| 4 | Merge of the SPK-03 and CON-01 PRs. Needed before CON-02 is approved. Drafting is not blocked. | Open | Before CON-02 approval. No calendar deadline. | CONTRIBUTING §2.4 (SDD workflow) and §11 |
| 5 | Confirm the architecture decisions CON-02 depends on: DEC-04 (provisional, tested by SPK-03), and whether DEC-12 affects dispatch in CON-02. | Open | Before CON-02 approval | §4.2 |
| 6 | Name the dependent pairs for CON-02. The owners of F06-02 (dispatch) and F07-01 (consumer) are not assigned in Sprint 1. Their agreement is required. | Open | Before CON-02 agreement is requested | [SDD Workflow](../sdd-workflow.md) §2.4 |
| 7 | Stretch stays outside the Sprint Goal (§1) and the exit criteria (§11). | Resolved (mentor, 9 Nov) | — | §11 |
| 8 | DEC-13 (editor session preservation) stays Open until the mechanism is approved. | Open | Sprint 3 (backlog §6.3) | §4.2 |
| 9 | Whether CON-02 and SPEC-F04 are removed from the Sprint 2 and Sprint 3 rows of backlog §10. | Open | Sprint 2 planning | Avoids counting the items twice |

---

## 9. Proposed GitHub Issues

**These issues are proposed, not created.** Create them only after the mentor accepts this plan. Issue numbers are assigned by GitHub on creation. The S1 numbers here are plan references.

Labels: only the six agreed workflow labels are used. No new labels are introduced. Feature items use the `spec:` states, then `in-progress`, `in-review` and `done`. Technical tasks, spikes and contracts use `in-progress`, `in-review` (while the pull request is open) and `done`. See [CONTRIBUTING §4](../../CONTRIBUTING.md#4-github-issue-workflow).

| Ref | Title | Pair | Depends on | Estimate (pair-days) | Acceptance criteria | Labels |
|-----|-------|------|------------|----------------------|---------------------|--------|
| S1-01 | [INF-03] Local infrastructure: cluster, PostgreSQL, RabbitMQ | A | DEC-06 (provisional, day 1) | 4 | (1) Cluster running by end of day 2 (M1). (2) RabbitMQ and PostgreSQL reachable by end of day 3 (M2). (3) Setup steps in the README, reproduced by a second pair (M3, day 4). | `in-progress`, `in-review`, `done` |
| S1-02 | [SPK-01] Secure Go execution spike | B | S1-01 M1; DEC-01 and DEC-02 provisional | 8 | SC-01 to SC-10 recorded with Verified, Partially verified or Not verified per criterion. Outcome per §7.1. Report merged to `docs/spikes/SPK-01.md` by 13 Nov 15:00. | `in-progress`, `in-review`, `done` |
| S1-03 | [SPK-02] Auth0 OIDC integration check | C | DEC-03 provisional; MA-05 | 2 | SC-A to SC-C recorded. Failure criteria F-06 and F-07 reported. Report by end of day 2. | `in-progress`, `in-review`, `done` |
| S1-04 | [SPK-03] RabbitMQ and Watermill check | C | S1-01 M2; DEC-04 provisional | 4 | SC-A to SC-D recorded. Failure criteria F-08 to F-10 reported. Report by end of day 6. | `in-progress`, `in-review`, `done` |
| S1-05 | [INF-01] Go module, layout, lint and format | D | DEC-14 (provisional, day 1) | 1.5 | `go.mod` with the Go version from ADR-001. Lint and format configured. Repository builds and passes lint locally. Layout section in the README. | `in-progress`, `in-review`, `done` |
| S1-06 | [INF-02] CI baseline on GitHub Actions | D | S1-05 | 1.5 | Pipeline runs on push and pull request. Build, lint, unit-test and dependency checks configured. Green on `main` by day 5. No secrets in the repository. | `in-progress`, `in-review`, `done` |
| S1-07 | [INF-04] ADR process and first records | D | None | 0.5 | ADR template merged in `docs/architecture/adr/`. ADR-001 (Go version, DEC-14) and ADR-002 (local cluster, DEC-06) recorded as Provisional. | `in-progress`, `in-review`, `done` |
| S1-08 | [CON-01] Submission lifecycle and verdict model | D | None | 1.5 | States, transitions and verdicts agreed. Agreement from Pairs A, B and C recorded on the issue by day 8. Infrastructure vs test-failure classification stated. Contract in `docs/contracts/CON-01.md`. | `in-progress`, `in-review`, `done` |
| S1-09 | [SPEC-F01] Accounts and access specification | D | S1-08 (draft) | 1.5 | Spec written from `specs/templates/feature-spec-template.md`. Acceptance criteria numbered AC-01 onward. Security section complete. Status set to In Review (`spec:in-review`) by day 8. Reviewed by Pairs A and C. | `spec:draft`, then `spec:in-review`, then `spec:approved` at the start of Sprint 2 |
| S1-10 (stretch) | [CON-02] Messaging contract | C | S1-08 (CON-01) and S1-04 (SPK-03) PRs, approved, merge pending. Drafting starts now. | 1.5 | Contract in `docs/contracts/CON-02.md` covering queue names, message schema, idempotency key, acknowledgement, retry and dead-letter behaviour. Dependent pairs named. Reviewed by Pair D. Agreement recorded by each dependent pair. Not merged before agreement. | `in-progress`, `in-review`, `done` |
| S1-11 (stretch) | [SPEC-F04] Browser code editor specification | D | None (backlog §4.1); DEC-13 recorded in the spec | 1.5 | Spec at `specs/features/F04-browser-code-editor.md`, from the template. Acceptance criteria numbered AC-01 onward. Covers Go source editing, starter templates, active-session preservation and edge cases. Status In Review (`spec:in-review`). Not approved in Sprint 1. | `spec:draft`, then `spec:in-review`; `spec:approved` is not a Sprint 1 outcome |

**Not proposed as Sprint 1 issues:** SPEC-F02, CON-04 (Sprint 2), INF-08 (Sprint 3), and all other backlog items.

---

## 10. Risks and Open Questions for Sprint 1

| ID | Risk or question | Mitigation |
|----|------------------|------------|
| R1 | Zero slack in committed effort. Any unplanned work displaces committed items. | Headroom is held per pair (§5.1). Spikes are timeboxed. Any new work goes to the mentor before it is accepted. |
| R2 | M1 (cluster by day 2) slips | Pair A prioritises INF-03. Pair B's days 1–2 do not depend on the cluster. Fallback in §6. |
| R3 | Pair B's negative-case tests take longer than the timebox | Criteria are ordered by security value (SC-04, SC-05 before SC-09, SC-10). Unmeasured criteria are reported, not extended. |
| R4 | The training assumption (MA-18) is wrong, and pairs B and C need general training | Headroom falls to 0. Re-cut at the Day 1 decision session. |
| R5 | Cross-pair review of SPEC-F01 and CON-01 takes longer than the review overhead | Review is allocated to Pairs A and C from day 5. Pair D does not chase reviews outside its headroom. |
| R6 | Spike reports are merged late, so the Sprint 2 decisions (DEC-01 to DEC-04) slip | Reports are merged on the day their timebox ends. |
| R7 | SPEC-F02 and CON-04 move to Sprint 2, raising its load | Sprint 2 load to be re-checked at the Sprint 2 planning. |
| R8 | Stretch work relies on PRs not yet merged (SPK-03, CON-01, SPEC-F01). Stretch could displace committed review work in week 2, when Pair D's CON-01 agreement and SPEC-F01 review are due. | Drafting only before merge. Approval waits for merge, cross-pair approval and confirmed decisions (§8.4, items 3–5). Committed items come first. Any stretch item not finished by 13 November is carried over. |
| Q1 | Spike label | Resolved. Spikes use `in-progress`, `in-review` and `done`, within the six agreed labels. |
| Q2 | No additional engineering support is assumed | Confirmed by MA-11. If support is added, it is recorded as a backlog change, not absorbed here. |

---

## 11. Sprint 1 Exit and Demo

**Exit criteria (13 November, 15:00):**

- SPK-01, SPK-02 and SPK-03 reports merged to `docs/spikes/SPK-01.md`, `SPK-02.md` and `SPK-03.md`, each with its outcome.
- INF-02 green on `main`.
- INF-03 meets M1, M2 and M3.
- CON-01 agreed by Pairs A, B and C, with agreement on the issue.
- SPEC-F01 in review, with reviews from Pairs A and C.
- ADR-001 and ADR-002 recorded as Provisional in `docs/architecture/adr/`.

**Demos:**

- **Friday 6 November:** interim SPK-01 results (SC-01 to SC-05).
- **Friday 13 November:** Sprint 1 demo. Spike outcomes, CI baseline and local platform.

Stretch items (§4.2) are not exit criteria. If complete by the same deadline, they are shown at the 13 November demo. If not, they are carried over.

Sprint 1 ends at 15:00 on 13 November. Sprint 2 planning follows.
