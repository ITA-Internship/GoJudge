# Specification-Driven Development Workflow

This document describes how GoJudge work moves from request to `done`, for each of the four work item types. It covers the lifecycle, responsibilities, approvals, labels and requirement changes.

It is a companion to the other project documents:

- [CONTRIBUTING.md](../CONTRIBUTING.md): branching, review, Definition of Ready and Definition of Done for each type. Rules are not repeated here.
- [specs/templates/feature-spec-template.md](../specs/templates/feature-spec-template.md): the Feature Specification structure. Use it for every feature.

## 1. Principles

- **Specify before you build.** No feature is implemented until its Feature Specification is approved.
- **Right artefact for the work.** Features need a Feature Specification. Technical spikes need a Spike Brief. Technical tasks need acceptance criteria. Contracts need cross-pair review. See §2.
- **Contracts first.** API and data contracts are agreed before dependent work starts.
- **Cross-pair review.** Every work item is reviewed by at least one pair other than the author's.
- **Mentor sets priorities.** The mentor owns requirement priorities and approves significant scope changes. Interns clarify, specify, investigate and implement.
- **Traceable work.** Every item links its GitHub Issue, its artefact and its pull request.

## 2. Work Item Types and Lifecycles

Each GitHub Issue is one of four types. The type is chosen before the issue is created. The table in [CONTRIBUTING §2](../CONTRIBUTING.md#2-sdd-rules-and-work-item-types) sets out the artefacts, locations and branch prefixes.

| Type | Artefact | Approval before work starts | Ends with |
|------|----------|-----------------------------|-----------|
| **Feature** | Feature Specification (`specs/features/`) | Specification approval by a reviewer pair | Implementation merged to `main`, `done` |
| **Technical task** | Acceptance criteria on the issue | None beyond the Definition of Ready | Pull request merged to `main`, `done` |
| **Technical spike** | Spike Brief on the issue; report at `docs/spikes/SPK-XX.md` | None beyond the Spike Brief being published | Report merged, outcome recorded, `done` |
| **Contract** | Contract at `docs/contracts/CON-XX.md` | Cross-pair review, and agreement by every dependent pair | Contract merged, dependent work may start, `done` |

### 2.1 Feature lifecycle

| # | Stage | Owner | Output | Gate to next stage |
|---|-------|-------|--------|--------------------|
| 1 | Requirement | Interns clarify; mentor sets priority | Requirement described in an issue (`spec:draft`), linked to an MVP feature ID (F01–F12) | Issue exists, feature ID is assigned, priority set by the mentor |
| 2 | Feature Specification | Responsible pair | Spec in `specs/features/` based on the template | Spec is complete enough for review |
| 3 | Technical Design and Contracts | Responsible pair, with dependent pairs | Design, contracts and open questions recorded in the spec | Contracts agreed with dependent pairs (CONTRIBUTING §5.2) |
| 4 | Cross-Pair Specification Review | Reviewer pair | Review outcome in the spec's review table (`spec:in-review`) | No open blocking comments |
| 5 | Specification Approval | Reviewer pair | Spec status set to **Approved** (`spec:approved`) | Definition of Ready met (CONTRIBUTING §10) |
| 6 | Implementation | Responsible pair | Code on a `feature/` branch (`in-progress`), following the approved spec | Tests written and passing locally |
| 7 | Automated Verification | Responsible pair | Tests and CI checks passing; acceptance criteria verified | All configured CI checks pass |
| 8 | Pull Request | Responsible pair | Pull request referencing the issue and spec (`in-review`) | Pull request ready for review |
| 9 | Implementation PR Approval | At least one reviewer from another pair | Approval on the pull request; comments resolved | At least one approval given |
| 10 | Merge and Done | Responsible pair | Changes integrated into `main`; issue labelled `done` and closed | Definition of Done met (CONTRIBUTING §11) |
| 11 | Weekly Demo | Responsible pair | Feature shown at the weekly Friday demo | Not a gate. Does not block merge or `done`. |

### 2.2 Technical task lifecycle

1. The issue states the objective, deliverables, numbered acceptance criteria, owner pair, estimate and dependencies.
2. The owner pair works on an `infra/` or `docs/` branch.
3. Each acceptance criterion is evidenced in the pull request.
4. At least one reviewer from another pair approves the pull request.
5. The change is merged, the issue is labelled `done`, and the issue is closed.

### 2.3 Technical spike lifecycle

1. The Spike Brief is published on the issue before the timebox starts: hypothesis, timebox, experiment, success and failure criteria, required evidence.
2. The owner pair runs the experiment on a `spike/` branch. Experiment code is not merged to `main`.
3. Each success criterion is recorded as Verified, Partially verified or Not verified, with its evidence. Failure criteria met during the spike are reported to the mentor on the day they are met.
4. The timebox ends on its last day. The owner pair reports the outcome reached, whatever its state. An extension needs mentor approval.
5. The report is written in `docs/spikes/SPK-XX.md` on a `docs/` branch, and reviewed by at least one pair other than the owner.
6. Decisions that follow from the outcome are recorded separately, and do not change the spike's outcome.

### 2.4 Contract lifecycle

1. The author pair drafts the contract in `docs/contracts/CON-XX.md`, naming the dependent pairs.
2. The author pair opens a pull request on a `docs/` branch.
3. At least one reviewer from another pair reviews it.
4. Every dependent pair records its agreement on the issue or the pull request.
5. The contract is merged. Dependent implementation may now start.
6. A later change to the contract needs renewed cross-pair review and renewed agreement.

## 3. Responsibilities in Rotating Pairs

The team of eight interns works in four development pairs. Pairs rotate every two weeks, at sprint boundaries.

- **Responsible pair:** owns the work item from its artefact to merge. Writes the specification, brief, criteria or contract, the implementation and the tests, as the type requires. Keeps the issue and state label current. Presents features at the weekly demo.
- **Reviewer pair:** always a pair other than the responsible pair. Reviews and approves the artefact, and reviews the pull request. At least one reviewer is required. Additional reviewers are optional.
- **Dependent pair:** any pair whose work depends on a contract. Agrees the contract before dependent development starts.
- **Mentor:** owns requirement priorities and approves significant scope changes (see §7). The mentor does not replace cross-pair review.

A pair is never the reviewer of its own work. Reviewer assignment should rotate so that every pair reviews across the cycle.

**Work that spans a rotation:** a work item may not finish within one two-week sprint. Before the rotation, the outgoing responsible pair updates the issue with the current label, artefact status and open questions. The incoming responsible pair confirms ownership on the issue.

## 4. Approvals

Each type has its own approval. The approvals are not interchangeable.

| Approval | Type | Given by | Gates |
|----------|------|----------|-------|
| **Specification approval** | Feature | Reviewer pair | Start of implementation. Status **Approved** and label `spec:approved`. |
| **Contract agreement** | Contract | Every dependent pair, after cross-pair review | Start of dependent implementation |
| **Implementation PR approval** | All types | At least one reviewer from another pair | Merge |

Specification approval does not approve the code. Pull request approval does not replace specification approval. A Spike Brief is published, not approved; it does not need a reviewer pair before the spike starts.

The Feature Specification status is one of the values defined in the [template](../specs/templates/feature-spec-template.md#1-feature-metadata-required):

| Status | Meaning |
|--------|---------|
| Draft | Being written by the responsible pair. Not ready for review. |
| In Review | Submitted to the reviewer pair. Changes may be requested. |
| Approved | Reviewed, with no open blocking comments. Development may start. |

Approval rules:

- Only the reviewer pair sets a specification to **Approved**.
- The author must not approve their own specification.
- A specification must not be marked **Approved** while acceptance criteria are ambiguous, a required contract is undefined, or a critical question is open.
- A materially changed approved specification returns to **In Review** and needs renewed cross-pair review (see §7).

## 5. Definition of Ready and Definition of Done

Both definitions are set for each type in [CONTRIBUTING §10](../CONTRIBUTING.md#10-definition-of-ready) and [CONTRIBUTING §11](../CONTRIBUTING.md#11-definition-of-done). In short:

- **Feature, ready:** specification approved, acceptance criteria defined, contracts agreed, dependencies identified, security considerations documented.
- **Feature, done:** implementation matches the specification, tests and CI checks pass, a reviewer from another pair has approved, documentation is updated, and the feature is merged into `main`.
- **Technical task, ready:** objective, deliverables, numbered acceptance criteria, owner pair, estimate and dependencies are on the issue.
- **Technical task, done:** every acceptance criterion is evidenced, CI passes where code changed, a reviewer from another pair has approved, and the change is merged.
- **Technical spike, ready:** Spike Brief is on the issue, with the timebox dates set.
- **Technical spike, done:** the timebox is respected, every criterion is recorded, the outcome is stated, and the report is merged and reviewed.
- **Contract, ready:** draft contract is in `docs/contracts/`, with the dependent pairs named.
- **Contract, done:** every dependent pair has agreed, a reviewer from another pair has approved, and the contract is merged.

A work item reaches `done` once its Definition of Done is satisfied. The weekly Friday demo is separate. It is not part of the Definition of Done, and no item is held back from `done` to wait for it.

## 6. GitHub Issue States

GitHub Issues have only *open* and *closed* by default. The team tracks workflow states with the six agreed labels. No other labels are used for state. Exactly one state label is on an issue at a time.

| Label | Feature | Technical task | Technical spike | Contract | Meaning |
|-------|:-------:|:--------------:|:---------------:|:--------:|---------|
| `spec:draft` | Yes | No | No | No | Requirement captured; specification being written (Draft). Applied by the feature form. |
| `spec:in-review` | Yes | No | No | No | Specification submitted for cross-pair review (In Review) |
| `spec:approved` | Yes | No | No | No | Specification approved and Definition of Ready met |
| `in-progress` | Yes | Yes | Yes | Yes | Work has started on its branch |
| `in-review` | Yes | Yes | Yes | Yes | Pull request is open for review |
| `done` | Yes | Yes | Yes | Yes | Definition of Done met. The issue is closed. |

The issue forms for technical tasks and spikes apply no label. The team applies `in-progress` when work starts.

## 7. Handling Requirement Changes

Requirements change. The rules below keep each artefact in step with the code.

**Who decides:** the mentor approves changes to requirement priorities and significant scope changes. Interns may clarify and update a specification without mentor approval, unless the change alters priority or scope significantly.

**What needs renewed review:**

- **Feature specification:** a material change to an approved specification. Material changes include changes to acceptance criteria, functional requirements, contracts, scope, dependencies or security considerations. Editorial changes, such as typo fixes or formatting, do not need renewed review.
- **Contract:** any change to an agreed contract needs renewed cross-pair review and renewed agreement.
- **Technical task or spike:** a change to its acceptance criteria or success criteria is recorded on the issue. A change that alters the scope of a spike's timebox needs mentor approval (MA-12).

**Feature, before implementation starts:**

1. Update the specification.
2. If the change is significant, obtain mentor approval.
3. Set its status to **In Review** and set the issue label to `spec:in-review`.
4. Obtain review from the reviewer pair again, before development continues.

**Feature, after approval, before merge:**

1. Stop work on the affected part of the implementation.
2. If the change is significant, obtain mentor approval.
3. Update the specification and return it to **In Review**. The issue label moves to `spec:in-review`.
4. Obtain re-review. Affected contracts must be agreed again with dependent pairs.
5. Resume implementation only after approval.

**After merge:**

1. Raise a new issue for the change, linked to the original feature ID.
2. Follow the stages in §2.1 for the change, from stage 2 onward.

Do not change code to match a new requirement without first updating and approving the specification.

## 8. Quick Reference

```text
Feature
  Requirement (spec:draft; priority set by mentor)
    → Spec Draft → In Review (spec:in-review) → Approved (spec:approved)   [specification approval]
    → Implement on feature/ branch (in-progress)
    → Tests and CI pass → Pull request (in-review)
    → Implementation PR approval (another pair) → Merge → Issue closed (done)
    → Weekly Friday demo (not a gate)

Technical task
  Acceptance criteria on issue (in-progress) → infra/ or docs/ branch
    → Pull request (in-review) → PR approval (another pair) → Merge → Issue closed (done)

Technical spike
  Spike Brief on issue (in-progress) → spike/ branch (experiment, not merged)
    → Report at docs/spikes/SPK-XX.md on docs/ branch (in-review) → Review (another pair)
    → Merge → Issue closed (done)

Contract
  Draft docs/contracts/CON-XX.md (in-progress) → docs/ branch
    → Pull request (in-review) → Cross-pair review → Dependent pairs agree
    → Merge → Issue closed (done)
```

Material changes to an approved specification or agreed contract return it to review.
