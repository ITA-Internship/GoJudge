# Sprint 1 Plan: 2–13 November 2026

**Status:** Draft for mentor review. Not approved. No GitHub Issues have been created from this plan.

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
| Stretch effort | 0 |
| Headroom, from removed double-counting (§3) | 2.0 |

Committed effort equals realistic net. This is intentional, and it has two safeguards:

- Spike timeboxes are fixed. They end on day 10 whatever their state, and they are not extended without MA-12.
- The 2.0 pair-days of headroom in §3 are held per pair (§5.1). They are not banked as committed work.

The backlog's previous Sprint 1 commitment was 28.5 pair-days. That is 4 pair-days over realistic capacity, and this plan removes that over-commitment.

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

### 4.2 Stretch

**None.** Sprint 1 has no stretch items. Headroom is not used for scope.

### 4.3 Deferred from backlog §9

| Item | Backlog §9 | This plan | Effect |
|------|------------|-----------|--------|
| INF-08 Frontend scaffold | Sprint 1 (1.5) | Sprint 3 | Already moved in backlog revision 3. Not needed until F01-03 in Sprint 4. |
| CON-04 Problem packaging, draft | Sprint 1 (2) | Sprint 2 (draft and final, 4 in total) | DEC-09 product input is needed earlier in Sprint 2. Sprint 2 load rises. |
| SPEC-F02 Problem catalogue specification, draft | Stretch (1.5) | Sprint 2 | SPEC-F02 approval moves to about the end of Sprint 2 week 1. F02-01 may slip into Sprint 2 week 2. |

### 4.4 Ownership changes from backlog §9

Backlog §9 gave Pair A INF-01, INF-02, INF-03, INF-04 and INF-08. This plan gives Pair D INF-01, INF-02 and INF-04, and Pair A only INF-03. Reason: Pair A carries 4 pair-days of INF-03, and a full allocation of the infrastructure items to Pair A would exceed its net capacity (§5.1).

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

- Spike report merged to `docs/spikes/` by end of day 10, with SC-01 to SC-10 each recorded as Verified, Partially verified or Not verified, using §7.1.
- Negative-case programs for the spike, held in a spike branch. Not merged to `main`.

### 5.4 WS-3 Identity and messaging (Pair C)

**Owns:** SPK-02 (days 1–2) and SPK-03 (days 3–6).

**Measurable deliverables:**

- SPK-02 report by end of day 2, with SC-A to SC-C recorded per §7.2.
- SPK-03 report by end of day 6, with SC-A to SC-D recorded per §7.3. Days 7–10 are available for review and for agreeing CON-02 input.

**Reviews by Pair C:** SPEC-F01 (0.5) and CON-01 (1.0). PR review (0.5).

### 5.5 WS-4 Foundations and specifications (Pair D)

**Owns:** INF-01, INF-02, INF-04, CON-01, SPEC-F01.

**Measurable deliverables:**

- **Day 1:** ADR-0001 records the Go version pin (DEC-14) as provisional. INF-01 starts.
- **Day 3:** Repository builds and passes lint locally. INF-02 CI pipeline runs on push and pull request. Green on `main` by day 5.
- **Day 4:** ADR-0002 records the local cluster choice (DEC-06) as provisional. ADR template merged (INF-04).
- **Day 5:** CON-01 draft circulated to Pairs A, B and C for agreement. Agreement recorded by day 8.
- **Day 8:** SPEC-F01 in cross-pair review, reviewed by Pairs A and C.

Pair D's headroom (0.75) is held for review responses and for the mentor's decisions in §8.

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
| MA-17 | **New.** Confirm that spikes are exempt from specification approval, because they produce reports and not features | SDD says no implementation starts before approval, and spikes write code. Without a ruling, SPK-01 is blocked. | Confirm the exemption, with the rule that spike code is not merged to `main` and is not production code. |
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

---

## 9. Proposed GitHub Issues

**These issues are proposed, not created.** Create them only after the mentor accepts this plan. Issue numbers are assigned by GitHub on creation. The S1 numbers here are plan references.

Labels: the six agreed labels are used for specification and implementation items. Spikes have no agreed label. Proposed: `in-progress` when started and `done` when closed. A spike label is an open question for the mentor (§10).

| Ref | Title | Pair | Depends on | Estimate (pair-days) | Acceptance criteria | Labels |
|-----|-------|------|------------|----------------------|---------------------|--------|
| S1-01 | [INF-03] Local infrastructure: cluster, PostgreSQL, RabbitMQ | A | DEC-06 (provisional, day 1) | 4 | (1) Cluster running by end of day 2 (M1). (2) RabbitMQ and PostgreSQL reachable by end of day 3 (M2). (3) Setup steps in the README, reproduced by a second pair (M3, day 4). | `in-progress`, `done` |
| S1-02 | [SPK-01] Secure Go execution spike | B | S1-01 M1; DEC-01 and DEC-02 provisional | 8 | SC-01 to SC-10 recorded with Verified, Partially verified or Not verified per criterion. Outcome per §7.1. Report merged to `docs/spikes/` by 13 Nov 15:00. | `in-progress`, `done` |
| S1-03 | [SPK-02] Auth0 OIDC integration check | C | DEC-03 provisional; MA-05 | 2 | SC-A to SC-C recorded. Failure criteria F-06 and F-07 reported. Report by end of day 2. | `in-progress`, `done` |
| S1-04 | [SPK-03] RabbitMQ and Watermill check | C | S1-01 M2; DEC-04 provisional | 4 | SC-A to SC-D recorded. Failure criteria F-08 to F-10 reported. Report by end of day 6. | `in-progress`, `done` |
| S1-05 | [INF-01] Go module, layout, lint and format | D | DEC-14 (provisional, day 1) | 1.5 | `go.mod` with the Go version from ADR-0001. Lint and format configured. Repository builds and passes lint locally. Layout section in the README. | `in-progress`, `done` |
| S1-06 | [INF-02] CI baseline on GitHub Actions | D | S1-05 | 1.5 | Pipeline runs on push and pull request. Build, lint, unit-test and dependency checks configured. Green on `main` by day 5. No secrets in the repository. | `in-progress`, `done` |
| S1-07 | [INF-04] ADR process and first records | D | None | 0.5 | ADR template merged. ADR-0001 (Go version, DEC-14) and ADR-0002 (local cluster, DEC-06) recorded as Provisional. | `in-progress`, `done` |
| S1-08 | [CON-01] Submission lifecycle and verdict model | D | None | 1.5 | States, transitions and verdicts agreed. Agreement from Pairs A, B and C recorded on the issue by day 8. Infrastructure vs test-failure classification stated. | `in-progress`, `done` |
| S1-09 | [SPEC-F01] Accounts and access specification | D | S1-08 (draft) | 1.5 | Spec written from `specs/templates/feature-spec-template.md`. Acceptance criteria numbered AC-01 onward. Security section complete. Status set to In Review (`spec:in-review`) by day 8. Reviewed by Pairs A and C. | `spec:draft`, then `spec:in-review`, then `spec:approved` at the start of Sprint 2 |

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
| Q1 | Spike label: the six agreed labels do not cover spikes | Mentor decision (§8.1, MA-17 follow-up). Proposed: `in-progress` and `done` only. |
| Q2 | No additional engineering support is assumed | Confirmed by MA-11. If support is added, it is recorded as a backlog change, not absorbed here. |

---

## 11. Sprint 1 Exit and Demo

**Exit criteria (13 November, 15:00):**

- SPK-01, SPK-02 and SPK-03 reports merged to `docs/spikes/`, each with its outcome.
- INF-02 green on `main`.
- INF-03 meets M1, M2 and M3.
- CON-01 agreed by Pairs A, B and C, with agreement on the issue.
- SPEC-F01 in review, with reviews from Pairs A and C.
- ADR-0001 and ADR-0002 recorded as Provisional.

**Demos:**

- **Friday 6 November:** interim SPK-01 results (SC-01 to SC-05).
- **Friday 13 November:** Sprint 1 demo. Spike outcomes, CI baseline and local platform.

Sprint 1 ends at 15:00 on 13 November. Sprint 2 planning follows.
