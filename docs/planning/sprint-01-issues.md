# Sprint 1 GitHub Issues (S1-01 to S1-09)

**Status:** Copy-ready drafts. No GitHub Issues have been created from this file.

**Source of truth:** [Sprint 1 Plan](sprint-01.md), approved for execution. Issue IDs, pairs, estimates, dependencies and acceptance criteria are taken from its §4, §5, §6 and §9. Rules come from [CONTRIBUTING](../../CONTRIBUTING.md) and [docs/sdd-workflow.md](../sdd-workflow.md).

**Links:** paths are relative to this file, `docs/planning/`. They resolve when the file is viewed in the repository.

---

## Conventions

Every issue is one of four work item types. Each type has its own issue form, artefact and branch prefix ([CONTRIBUTING §2](../../CONTRIBUTING.md#2-sdd-rules-and-work-item-types)).

| Type | Issues in Sprint 1 | Required before work | Artefact | Branch | Issue form |
|------|--------------------|----------------------|----------|--------|------------|
| Feature | S1-09 | Approved Feature Specification | `specs/features/` | `docs/` for the specification | Feature request (`spec:draft`) |
| Technical task | S1-01, S1-05, S1-06, S1-07 | Acceptance criteria | Repository changes, README, ADRs in `docs/architecture/adr/` | `infra/` or `docs/` | Technical task (no label) |
| Technical spike | S1-02, S1-03, S1-04 | Spike Brief | Report in `docs/spikes/SPK-XX.md`. Experiment code on `spike/`, not merged. | `spike/` for code, `docs/` for the report | Technical spike (no label) |
| Contract | S1-08 | Draft contract and dependent pairs named | `docs/contracts/CON-01.md` | `docs/` | Technical task, work type "Contract" |

**Labels.** Only the six existing workflow labels are used. No new labels are introduced.

- Technical tasks, spikes and contracts: `in-progress` when started, `in-review` while the pull request is open, `done` when closed.
- S1-09 (feature): `spec:draft`, then `spec:in-review`, then `spec:approved`, then the implementation states.
- Exactly one state label on an issue at a time ([docs/sdd-workflow.md §6](../sdd-workflow.md#6-github-issue-states)).

**Pairs.** Pair labels A–D are placeholders. The mentor assigns named interns to pairs. No named interns appear here.

---

## S1-01: [INF-03] Local infrastructure: cluster, PostgreSQL, RabbitMQ

**Form:** Technical task · **Labels:** `in-progress`, `in-review`, `done` · **Branch:** `infra/local-infrastructure`

```markdown
### Objective

Provide a local Kubernetes cluster, PostgreSQL and RabbitMQ that other Sprint 1 work can use, so that the secure execution and messaging spikes can start.

### Type and ownership

- Type: Technical task (infrastructure)
- Owner: Pair A
- Estimate: 4 pair-days
- Sprint: 1 (2–13 November 2026)

### Dependencies

- DEC-06 (local cluster tool: k3d or Minikube), provisional from Sprint 1 day 1. Confirmed under MA-02.

### Deliverables

- Local cluster running, using the tool chosen under DEC-06.
- PostgreSQL and RabbitMQ running and reachable from the cluster and from a developer machine.
- Setup steps in the README.
- Milestones: M1 cluster ready by end of day 2 (Tue 3 Nov). M2 RabbitMQ and PostgreSQL reachable by end of day 3 (Wed 4 Nov). M3 setup steps in the README, reproduced by a second pair, by end of day 4 (Thu 5 Nov).

### Acceptance criteria

1. The cluster is running by end of day 2. Pair B can create a Job in it.
2. RabbitMQ and PostgreSQL are reachable from the cluster and from a developer machine by end of day 3. Pair C can publish and consume on RabbitMQ.
3. The README contains the setup steps. A second pair reproduces the environment from the README alone.
4. No secrets or credentials are committed.

### Definition of Done

Per [CONTRIBUTING §11](../../CONTRIBUTING.md#11-definition-of-done), technical task:

- Acceptance criteria 1–4 met, with evidence in the pull request.
- Documentation updated (README).
- No secrets or environment-specific configuration committed.
- Pull request approved by a reviewer from another pair.
- Merged to `main`. Issue labelled `done` and closed.

### Documentation

- [Sprint 1 Plan](sprint-01.md), §4.5, §5.2 and §6
- [Product Backlog](../product/product-backlog.md), §4.2 (INF-03)
- [Contributing Guide](../../CONTRIBUTING.md), §6 (no secrets, environment configuration)
```

---

## S1-02: [SPK-01] Secure Go execution spike

**Form:** Technical spike · **Labels:** `in-progress`, `in-review`, `done` · **Branch:** `spike/SPK-01-…` for code, `docs/spk-01-report` for the report

```markdown
### Objective

Answer one question: can `go test` run in an isolated Kubernetes Job, created with client-go, under the controls required by the MVP proposal §7? This is the Sprint 1 priority spike, and it is on the critical path.

### Type and ownership

- Type: Technical spike (Spike Brief below)
- Owner: Pair B
- Estimate: 8 pair-days, timebox days 1–10 (Mon 2 Nov to Fri 13 Nov, report due 15:00 on 13 Nov)

### Dependencies

- S1-01 (INF-03): cluster ready by end of day 2 (M1). Pair B does not build its own cluster.
- DEC-01 and DEC-02, provisional. Tested by this spike, not approved in advance.

### Deliverables

- Spike report at `docs/spikes/SPK-01.md`.
- Each criterion SC-01 to SC-10 recorded as Verified, Partially verified or Not verified, with its evidence.
- Outcome stated: Go, Conditional go, No-go or Inconclusive.
- Experiment code kept on a `spike/` branch, not merged to `main`.

### Acceptance criteria

1. Report merged by 13 November, 15:00.
2. Every criterion SC-01 to SC-10 has a recorded status and evidence reference.
3. Outcome matches the rules in the Spike Brief.
4. Any failure criterion met is reported to the mentor on the day it is met.
5. Interim results shown at the Friday 6 November demo.

### Definition of Done

Per [CONTRIBUTING §11](../../CONTRIBUTING.md#11-definition-of-done), technical spike:

- Timebox respected. Every success criterion recorded with evidence.
- Outcome stated. Failure criteria met during the spike reported to the mentor.
- Report merged to `docs/spikes/SPK-01.md` and reviewed by a pair other than Pair B.
- Experiment code not merged to `main`.
- Follow-on decisions listed (DEC-01, DEC-02, DEC-08).
- Issue labelled `done` and closed, with the outcome in the closing comment.

---

### Spike Brief: SPK-01

**Hypothesis.** A Go test package can run inside a Kubernetes Job created through client-go, with CPU, memory and time limits, no network, a restricted filesystem, bounded output, log and report retrieval before cleanup, TTL cleanup and cancellation. Each control is verified by measurement, not assumed.

**Timebox.** 8 pair-days, days 1–10. Stop at 15:00 on 13 November, whatever the state. No extension without MA-12.

**Experiment.** Use only trusted test packages. No learner-submitted code runs in this spike.

| SC | What to observe |
|----|-----------------|
| SC-01 | A client-go program creates a Job that runs a trusted `go test` package, and the JSON result is retrieved. 10 of 10 runs. |
| SC-02 | A Job that exceeds its deadline is terminated within deadline + 5 s. 10 of 10 runs. |
| SC-03 | A memory-exhaustion program is killed by the limit, with the termination reason captured. 5 of 5 runs. |
| SC-04 | An outbound TCP attempt from the Job fails. 5 of 5 runs. |
| SC-05 | Writes outside the allowed path fail, and the root filesystem is read-only. 5 of 5 runs. |
| SC-06 | The container runs as non-root, with no privilege escalation and no service-account token mounted. Verified from the Job specification and by attempt. 5 of 5. |
| SC-07 | An output flood is truncated at the configured limit. 3 of 3 runs. |
| SC-08 | Logs and reports are retrieved before cleanup. 10 of 10 runs. |
| SC-09 | Finished Jobs are removed within TTL + 60 s. 10 of 10 runs. |
| SC-10 | A running Job is terminated within 30 s of a cancel request. 5 of 5 runs. |

Startup and run latency are recorded for information only. They are not a pass criterion.

**Success criteria.**

- **Go:** SC-01 to SC-10 all Verified.
- **Conditional go:** SC-01 and SC-08 Verified. Each other criterion that is not Verified becomes a named risk for SPEC-F07 and MA-03.

**Failure criteria** (any one means No-go):

- F-01: a Job cannot run a trusted `go test` package with the result retrieved (SC-01 fails).
- F-02: a Job past its deadline is not terminated (SC-02 fails).
- F-03: an outbound connection from the Job succeeds in any run (SC-04 fails).
- F-04: a write outside the allowed path succeeds in any run (SC-05 fails).
- F-05: logs or reports are lost before cleanup in any run (SC-08 fails).

**Inconclusive.** SC-01 Verified, but SC-02 to SC-10 not measured by the timebox end. Reported as Inconclusive, and DEC-02 is not approved.

**Required evidence.**

- For each SC: the observation, the number of runs, and where the raw record is kept (a link or an attached file in the report).
- The Job specification used for the runs.
- The outcome and the list of Not verified or Partially verified criteria with their risks.
- Confirmation that failure criteria F-01 to F-05 were checked.

**If it fails.** DEC-02 is not approved. The mentor chooses between a different execution model and option E in [backlog §11.5](../product/product-backlog.md#115-options-for-the-mentor). No SPEC-F07 or F07 work starts.

**Documentation.** [Sprint 1 Plan](sprint-01.md) §6 and §7.1 · [Product Backlog](../product/product-backlog.md) §5.1 · [MVP Proposal](../product/mvp-proposal.md) §7 (isolation controls).
```

---

## S1-03: [SPK-02] Auth0 OIDC integration check

**Form:** Technical spike · **Labels:** `in-progress`, `in-review`, `done` · **Branch:** `spike/SPK-02-…` for code, `docs/spk-02-report` for the report

```markdown
### Objective

Answer one question: can the Go API validate Auth0 OIDC tokens and map the role claim to learner or administrator, using either the Auth0 development tenant or a documented fixture?

### Type and ownership

- Type: Technical spike (Spike Brief below)
- Owner: Pair C
- Estimate: 2 pair-days, days 1–2 (Mon 2 Nov to Tue 3 Nov, report due end of day 2)

### Dependencies

- DEC-03 (authentication with Auth0), provisional. Tested by this spike, not approved in advance.
- MA-05: tenant access for this spike, or acceptance of the fixture route. Decided on day 1.

### Deliverables

- Spike report at `docs/spikes/SPK-02.md`.
- SC-A, SC-B and SC-C recorded with evidence.
- Outcome stated: Go, Conditional go or No-go.

### Acceptance criteria

1. Report merged by end of day 2 (Tue 3 Nov).
2. SC-A, SC-B and SC-C each recorded as Verified, Partially verified or Not verified.
3. If the fixture route is used, the report says so and records the tenant as not tested.
4. Failure criteria F-06 and F-07 are reported to the mentor on the day they are met.

### Definition of Done

Per [CONTRIBUTING §11](../../CONTRIBUTING.md#11-definition-of-done), technical spike:

- Timebox respected. Every success criterion recorded with evidence.
- Outcome stated. Failure criteria reported to the mentor.
- Report merged to `docs/spikes/SPK-02.md` and reviewed by a pair other than Pair C.
- Experiment code not merged to `main`.
- Follow-on decision listed (DEC-03).
- Issue labelled `done` and closed, with the outcome in the closing comment.

---

### Spike Brief: SPK-02

**Hypothesis.** The Go API can verify an Auth0 OIDC token, and map its role claim to learner or administrator.

**Timebox.** 2 pair-days, days 1–2. Stop at the end of day 2. No extension without MA-12.

**Experiment.** Sign in through the Auth0 development tenant, or through a documented fixture if tenant access is not in place (MA-05). Present the resulting token to a Go middleware.

| SC | What to observe |
|----|-----------------|
| SC-A | Sign-in works through the Auth0 development tenant, or through a documented fixture. Demonstrated once. |
| SC-B | The middleware checks signature, issuer, audience and expiry. It rejects a tampered token, an expired token and a wrong-audience token. 3 of 3. |
| SC-C | The role claim maps to learner or administrator. Correct for two test users. |

**Success criteria.** SC-A, SC-B and SC-C all Verified. A fixture route counts as Partially verified for SC-A, with the tenant listed as a named risk.

**Failure criteria:**

- F-06: the middleware accepts a tampered or expired token (SC-B fails). No-go for F01-01 as designed.
- F-07: neither the development tenant nor a documented fixture is available by day 2. Reported as Inconclusive, and MA-05 is escalated.

**Required evidence.**

- SC-A: a description of the sign-in path used (tenant or fixture) and whether it was tested against the tenant.
- SC-B: the three rejection results, with where the raw record is kept.
- SC-C: the claim mapping for the two users.
- Any secret values are kept out of the report.

**Documentation.** [Sprint 1 Plan](sprint-01.md) §7.2 · [Product Backlog](../product/product-backlog.md) §5.2 · [MVP Proposal](../product/mvp-proposal.md) §5 (authentication).
```

---

## S1-04: [SPK-03] RabbitMQ and Watermill check

**Form:** Technical spike · **Labels:** `in-progress`, `in-review`, `done` · **Branch:** `spike/SPK-03-…` for code, `docs/spk-03-report` for the report

```markdown
### Objective

Answer one question: can RabbitMQ with Watermill publish and consume submission messages without loss, redeliver after a consumer failure, keep one final state under duplicate delivery, and route poison messages to a dead-letter route?

### Type and ownership

- Type: Technical spike (Spike Brief below)
- Owner: Pair C
- Estimate: 4 pair-days, days 3–6 (Wed 4 Nov to Mon 9 Nov, report due end of day 6)

### Dependencies

- S1-01 (INF-03): RabbitMQ and PostgreSQL reachable by end of day 3 (M2).
- DEC-04 (messaging with RabbitMQ and Watermill), provisional. Tested by this spike, not approved in advance.

### Deliverables

- Spike report at `docs/spikes/SPK-03.md`.
- SC-A to SC-D recorded with evidence.
- Outcome stated: Go, Conditional go or No-go.

### Acceptance criteria

1. Report merged by end of day 6 (Mon 9 Nov).
2. SC-A to SC-D each recorded as Verified, Partially verified or Not verified.
3. Failure criteria F-08 to F-10 are reported to the mentor on the day they are met.
4. The outcome informs DEC-04 and DEC-12. Both decisions are taken after this report.

### Definition of Done

Per [CONTRIBUTING §11](../../CONTRIBUTING.md#11-definition-of-done), technical spike:

- Timebox respected. Every success criterion recorded with evidence.
- Outcome stated. Failure criteria reported to the mentor.
- Report merged to `docs/spikes/SPK-03.md` and reviewed by a pair other than Pair C.
- Experiment code not merged to `main`.
- Follow-on decisions listed (DEC-04, DEC-12).
- Issue labelled `done` and closed, with the outcome in the closing comment.

---

### Spike Brief: SPK-03

**Hypothesis.** RabbitMQ with Watermill delivers every message, redelivers after a consumer is killed mid-processing, produces one final state under duplicate delivery when an idempotency key is used, and routes poison messages to a dead-letter route.

**Timebox.** 4 pair-days, days 3–6. Stop at the end of day 6. No extension without MA-12.

**Experiment.** Publish and consume against the local RabbitMQ from S1-01. Simulate consumer failure, duplicate delivery and a poison message.

| SC | What to observe |
|----|-----------------|
| SC-A | Publish and consume with no message loss. 100 of 100 messages. |
| SC-B | A consumer killed during processing causes redelivery. 5 of 5 trials. |
| SC-C | A duplicate delivery of the same message leads to one final state, using an idempotency key. 50 of 50 duplicates. |
| SC-D | A poison message reaches the dead-letter route after the configured attempts. 3 of 3. |

**Success criteria.** SC-A to SC-D all Verified: Go. Conditional go if one or more are Partially verified, with each listed as a named risk for CON-02.

**Failure criteria:**

- F-08: any message is lost in SC-A (not 100 of 100).
- F-09: a redelivered duplicate produces two final states in SC-C (not 50 of 50).
- F-10: a poison message cannot reach the dead-letter route in SC-D (not 3 of 3).

**Required evidence.**

- For each SC: the counts observed, and where the raw record is kept.
- The idempotency approach used in SC-C, described in the report. Not an implementation decision.
- The outcome, and the named risks.

**Documentation.** [Sprint 1 Plan](sprint-01.md) §7.3 · [Product Backlog](../product/product-backlog.md) §5.3 · [MVP Proposal](../product/mvp-proposal.md) §6 (asynchronous processing).
```

---

## S1-05: [INF-01] Go module, layout, lint and format

**Form:** Technical task · **Labels:** `in-progress`, `in-review`, `done` · **Branch:** `infra/go-module-layout`

```markdown
### Objective

Create the Go module and repository layout, and configure lint and format, so that all later code work starts from the same baseline.

### Type and ownership

- Type: Technical task (infrastructure)
- Owner: Pair D
- Estimate: 1.5 pair-days
- Sprint: 1 (2–13 November 2026)

### Dependencies

- DEC-14 (Go version pin), provisional from day 1. Recorded in ADR-001 (S1-07). Confirmed under MA-02.

### Deliverables

- `go.mod` with the Go version from ADR-001.
- Lint and format configuration committed to the repository.
- A layout section in the README describing the directories.
- Layout reviewed by Pairs A to C.

### Acceptance criteria

1. The repository builds and passes lint locally.
2. The Go version in `go.mod` matches ADR-001.
3. The layout is described in the README.
4. No secrets, credentials or environment-specific configuration are committed.
5. Layout reviewed by at least one pair other than Pair D.

### Definition of Done

Per [CONTRIBUTING §11](../../CONTRIBUTING.md#11-definition-of-done), technical task:

- Acceptance criteria 1–5 met, with evidence in the pull request.
- CI green on the pull request once S1-06 is in place. Until then, local checks are recorded.
- README updated.
- Pull request approved by a reviewer from another pair.
- Merged to `main`. Issue labelled `done` and closed.

### Documentation

- [Sprint 1 Plan](sprint-01.md), §4.5 and §5.5
- [Product Backlog](../product/product-backlog.md), §4.2 (INF-01)
- [Contributing Guide](../../CONTRIBUTING.md), §6 (dependencies and configuration)
```

---

## S1-06: [INF-02] CI baseline on GitHub Actions

**Form:** Technical task · **Labels:** `in-progress`, `in-review`, `done` · **Branch:** `infra/ci-baseline`

```markdown
### Objective

Set up a CI pipeline on GitHub Actions that runs on push and pull request, so that every later change is checked automatically.

### Type and ownership

- Type: Technical task (infrastructure)
- Owner: Pair D
- Estimate: 1.5 pair-days
- Sprint: 1 (2–13 November 2026)

### Dependencies

- S1-05 (INF-01): module and lint configuration in place.

### Deliverables

- Pipeline that runs on push and pull request.
- Build, lint, unit-test and dependency checks configured.
- Pipeline green on `main` by end of day 5 (Fri 6 Nov).

### Acceptance criteria

1. The pipeline runs on push and on pull request.
2. Build, lint, unit tests and dependency checks are configured and run.
3. The pipeline is green on `main` by end of day 5.
4. No secrets are stored in the repository. Any credential the pipeline needs is held in GitHub settings, not in files.

### Definition of Done

Per [CONTRIBUTING §11](../../CONTRIBUTING.md#11-definition-of-done), technical task:

- Acceptance criteria 1–4 met, with evidence in the pull request.
- CI green on the pull request and on `main` ([CONTRIBUTING §7](../../CONTRIBUTING.md#7-testing-requirements)).
- README updated to say where the checks are defined.
- Pull request approved by a reviewer from another pair.
- Merged to `main`. Issue labelled `done` and closed.

### Documentation

- [Sprint 1 Plan](sprint-01.md), §4.5 and §5.5
- [Product Backlog](../product/product-backlog.md), §4.2 (INF-02)
- [Contributing Guide](../../CONTRIBUTING.md), §7 (testing requirements)
```

---

## S1-07: [INF-04] ADR process and first records

**Form:** Technical task · **Labels:** `in-progress`, `in-review`, `done` · **Branch:** `docs/adr-process`

```markdown
### Objective

Set up a lightweight Architecture Decision Record (ADR) process, and record the two provisional Sprint 1 decisions, so that decisions have a stable reference.

### Type and ownership

- Type: Technical task (documentation or ADR)
- Owner: Pair D
- Estimate: 0.5 pair-days
- Sprint: 1 (2–13 November 2026)

### Dependencies

- None.

### Deliverables

- ADR template in `docs/architecture/adr/`.
- `docs/architecture/adr/ADR-001.md`: Go version pin (DEC-14), status Provisional.
- `docs/architecture/adr/ADR-002.md`: Local cluster tool (DEC-06), status Provisional.
- An index of ADRs in the same folder.

### Acceptance criteria

1. ADR template merged.
2. ADR-001 and ADR-002 record the provisional choice, the reason, and the decision each depends on.
3. Both records say that the decision is not approved until MA-02 confirms it.
4. Location (`docs/architecture/adr/ADR-XXX.md`) and numbering (three digits) stated in the template.

### Definition of Done

Per [CONTRIBUTING §11](../../CONTRIBUTING.md#11-definition-of-done), technical task (documentation or ADR):

- Acceptance criteria 1–4 met, with evidence in the pull request.
- Documentation updated.
- Pull request approved by a reviewer from another pair.
- Merged to `main`. Issue labelled `done` and closed.

### Documentation

- [Sprint 1 Plan](sprint-01.md), §4.5, §5.5 and §8.1 (MA-02)
- [Product Backlog](../product/product-backlog.md), §6.1 (DEC-14 and DEC-06)
- [SDD Workflow](../sdd-workflow.md), §6 (issue state labels)
```

---

## S1-08: [CON-01] Submission lifecycle and verdict model

**Form:** Technical task (work type "Contract") · **Labels:** `in-progress`, `in-review`, `done` · **Branch:** `docs/con-01-lifecycle`

```markdown
### Objective

Agree one contract for the submission lifecycle, the verdicts, and the split between infrastructure failures and test failures, so that SPEC-F06, SPEC-F07 and CON-02 can be written against it.

### Type and ownership

- Type: Contract
- Owner: Pair D (author). Cross-pair review by at least one pair other than Pair D. Agreement from Pairs A, B and C.
- Estimate: 1.5 pair-days
- Sprint: 1 (2–13 November 2026)

### Dependencies

- None.

### Deliverables

- Contract at `docs/contracts/CON-01.md`: states (QUEUED, RUNNING, COMPLETED, FAILED), transitions, and the verdicts ACCEPTED, WRONG_ANSWER, COMPILATION_ERROR, RUNTIME_ERROR, TIME_LIMIT_EXCEEDED, MEMORY_LIMIT_EXCEEDED and INTERNAL_ERROR.
- A statement of how infrastructure failures are separated from test failures.

### Acceptance criteria

1. The states, transitions and verdicts match [MVP Proposal §3](../product/mvp-proposal.md#verdicts-and-submission-lifecycle).
2. Infrastructure failures and test failures are distinguished, and the distinction is stated.
3. Pairs A, B and C each record agreement on this issue by end of day 8 (Wed 11 Nov).
4. Any disagreement is recorded as an open question, not resolved by silence.

### Definition of Done

Per [CONTRIBUTING §11](../../CONTRIBUTING.md#11-definition-of-done), contract:

- Acceptance criteria 1–4 met.
- Every dependent pair (A, B, C) has recorded its agreement on the issue or the pull request.
- A reviewer from another pair has approved the pull request.
- Contract merged to `docs/contracts/CON-01.md`.
- SPEC-F06 and SPEC-F07 link to it when written (Sprint 2).
- Issue labelled `done` and closed.

### Documentation

- [Sprint 1 Plan](sprint-01.md), §4.5, §5.5 and §3 (item 7)
- [SDD Workflow](../sdd-workflow.md), §1 and §2.4 (contract lifecycle)
- [Contributing Guide](../../CONTRIBUTING.md), §5.2 (contract review)
- [Product Backlog](../product/product-backlog.md), §4.1 (CON-01)
- [MVP Proposal](../product/mvp-proposal.md), §3 (verdicts and submission lifecycle)
```

---

## S1-09: [SPEC-F01] Accounts and access specification

**Form:** Feature request · **Labels:** `spec:draft`, then `spec:in-review` by day 8, then `spec:approved` at the start of Sprint 2 · **Branch:** `docs/spec-f01`

```markdown
### Objective

Write the specification for F01 (accounts and access), so that F01-01 can start in Sprint 2 after approval.

### Type and ownership

- Type: Feature (specification)
- Owner: Pair D (author). Review by Pairs A and C. At least one reviewer is required.
- Estimate: 1.5 pair-days
- Sprint: 1 (2–13 November 2026). Approval at the start of Sprint 2 (16 November 2026).

### Dependencies

- S1-08 (CON-01): lifecycle draft circulated by day 5. Used for any reference to submissions.
- Authentication approach from SPK-02 (S1-03): may be provisional at drafting time. The spec records the open question if SPK-02 is not yet reported.

### Deliverables

- Specification at `specs/features/F01-accounts-and-access.md`, based on [the template](../../specs/templates/feature-spec-template.md).
- Acceptance criteria numbered AC-01 onward.
- Security section complete: authentication, authorisation, input validation, sensitive data.
- Status set to In Review by end of day 8.

### Acceptance criteria

1. Every template section is present. Sections that do not apply say "Not applicable" with a reason.
2. Acceptance criteria are numbered AC-01 onward, written as Given / When / Then.
3. Functional requirements cover F01 as stated in the [MVP Proposal §3](../product/mvp-proposal.md#3-mvp-functional-scope): register, sign in and out, and learner and administrator roles.
4. The security section is complete, and it records any open question from SPK-02.
5. No open question is marked as resolved without a decision reference.
6. Pairs A and C each record a review outcome in the spec's review table.

### Definition of Done

Per [CONTRIBUTING §10 and §11](../../CONTRIBUTING.md#10-definition-of-ready), feature specification:

- Acceptance criteria 1–6 met.
- Status is In Review, with the reviewing pairs recorded.
- Approval is not claimed in Sprint 1. Approval is recorded at the start of Sprint 2 (`spec:approved`).
- Specification merged to `main` through a `docs/` branch.

### Documentation

- [Feature Specification Template](../../specs/templates/feature-spec-template.md)
- [SDD Workflow](../sdd-workflow.md), §2.1, §4 (approvals) and §6 (labels)
- [Contributing Guide](../../CONTRIBUTING.md), §5.1 (specification review)
- [Sprint 1 Plan](sprint-01.md), §4.5 and §5.5
```

---

## Resolved open points

These were open in the previous version and are now resolved by the approved work item rules.

| Point | Resolution |
|-------|------------|
| Spike exemption (MA-17) | Spikes need a Spike Brief, not a Feature Specification ([CONTRIBUTING §2](../../CONTRIBUTING.md#2-sdd-rules-and-work-item-types)). |
| Spike location | `docs/spikes/SPK-XX.md` |
| Contract location and label | `docs/contracts/CON-XX.md`. Contracts use `in-progress`, `in-review` and `done`. |
| Branch naming | `feature/`, `fix/`, `infra/`, `spike/`, `docs/` ([CONTRIBUTING §3](../../CONTRIBUTING.md#3-git-branching-strategy)). |
| Issue creation route | Technical task form (no label) and spike form (no label). Feature form keeps `spec:draft`. |
| ADR location and numbering | `docs/architecture/adr/ADR-XXX.md`. Plan references updated to `ADR-001` and `ADR-002`. |

## Still open (not blocking Sprint 1)

1. **Spike Brief review.** The rules require a Spike Brief, not a review by another pair. [CONTRIBUTING §5.3](../../CONTRIBUTING.md#53-spike-briefs-and-reports) recommends review. Confirm whether it is required.
2. **Bugfix classification.** Proposed in [CONTRIBUTING §2](../../CONTRIBUTING.md#2-sdd-rules-and-work-item-types). Pending mentor confirmation.

---

## Summary of issues

| Ref | Title | Type | Owner | Estimate | Labels | Form |
|-----|-------|------|-------|----------|--------|------|
| S1-01 | [INF-03] Local infrastructure | Technical task | Pair A | 4 | `in-progress`, `in-review`, `done` | Technical task |
| S1-02 | [SPK-01] Secure Go execution spike | Technical spike | Pair B | 8 | `in-progress`, `in-review`, `done` | Technical spike |
| S1-03 | [SPK-02] Auth0 OIDC integration check | Technical spike | Pair C | 2 | `in-progress`, `in-review`, `done` | Technical spike |
| S1-04 | [SPK-03] RabbitMQ and Watermill check | Technical spike | Pair C | 4 | `in-progress`, `in-review`, `done` | Technical spike |
| S1-05 | [INF-01] Go module, layout, lint and format | Technical task | Pair D | 1.5 | `in-progress`, `in-review`, `done` | Technical task |
| S1-06 | [INF-02] CI baseline | Technical task | Pair D | 1.5 | `in-progress`, `in-review`, `done` | Technical task |
| S1-07 | [INF-04] ADR process and first records | Technical task | Pair D | 0.5 | `in-progress`, `in-review`, `done` | Technical task |
| S1-08 | [CON-01] Submission lifecycle and verdict model | Contract | Pair D | 1.5 | `in-progress`, `in-review`, `done` | Technical task (Contract) |
| S1-09 | [SPEC-F01] Accounts and access specification | Feature | Pair D | 1.5 | `spec:draft` → `spec:in-review` → `spec:approved` | Feature request |

Total estimate: 24.5 pair-days, matching the Sprint 1 plan.
