# Specification-Driven Development Workflow

This document describes how GoJudge features move from requirement to `done`, and how the weekly demo fits alongside them. It covers the lifecycle, responsibilities, approvals and requirement changes.

It is a companion to the other project documents:

- [CONTRIBUTING.md](../CONTRIBUTING.md): branching, review, Definition of Ready and Definition of Done. Rules are not repeated here.
- [specs/templates/feature-spec-template.md](../specs/templates/feature-spec-template.md): the specification structure. Use it for every feature.

## 1. Principles

- **Specify before you build.** No implementation starts until the feature specification is approved.
- **Contracts first.** API and data contracts are agreed before dependent work starts.
- **Cross-pair review.** Specifications and pull requests are reviewed by a pair other than the author's.
- **Mentor sets priorities.** The mentor owns requirement priorities and approves significant scope changes. Interns clarify, specify and implement.
- **Traceable work.** Every feature links a GitHub Issue, a specification and a pull request.

## 2. Lifecycle

| # | Stage | Owner | Output | Gate to next stage |
|---|-------|-------|--------|--------------------|
| 1 | Requirement | Interns clarify; mentor sets priority | Requirement described in an issue (`spec:draft`), linked to an MVP feature ID (F01–F12) | Issue exists, feature ID is assigned, priority set by the mentor |
| 2 | Feature Specification | Responsible pair | Spec in `specs/features/` based on the template | Spec is complete enough for review |
| 3 | Technical Design and Contracts | Responsible pair, with dependent pairs | Design, contracts and open questions recorded in the spec | Contracts agreed with dependent pairs |
| 4 | Cross-Pair Specification Review | Reviewer pair | Review outcome in the spec's review table (`spec:in-review`) | No open blocking comments |
| 5 | Specification Approval | Reviewer pair | Spec status set to **Approved** (`spec:approved`) | Definition of Ready met (see §5) |
| 6 | Implementation | Responsible pair | Code on a feature branch (`in-progress`), following the approved spec | Tests written and passing locally |
| 7 | Automated Verification | Responsible pair | Tests and CI checks passing; acceptance criteria verified | All configured CI checks pass |
| 8 | Pull Request | Responsible pair | Pull request referencing the issue and spec (`in-review`) | Pull request ready for review |
| 9 | Implementation PR Approval | Reviewer from another pair | Approval on the pull request; comments resolved | At least one approval given |
| 10 | Merge and Done | Responsible pair | Changes integrated into `main`; issue labelled `done` and closed | Definition of Done met (see §5) |
| 11 | Weekly Demo | Responsible pair | Feature shown at the weekly Friday demo | Not a gate. Does not block merge or `done`. |

Stages 6–10 follow the branching, pull request and review rules in [CONTRIBUTING.md](../CONTRIBUTING.md#3-git-branching-strategy).

## 3. Responsibilities in Rotating Pairs

The team of eight interns works in four development pairs. Pairs rotate every two weeks, at sprint boundaries.

- **Responsible pair:** owns the feature from specification to merge. Writes the specification, the implementation and the tests. Keeps the issue and spec status current. Presents the feature at the weekly demo.
- **Reviewer pair:** always a pair other than the responsible pair. Reviews and approves the specification, and reviews the pull request.
- **Dependent pair:** any pair whose work depends on a contract in the specification. Agrees the contract before dependent development starts.
- **Mentor:** owns requirement priorities and approves significant scope changes (see §7). The mentor does not replace cross-pair review.

A pair is never the reviewer of its own work. Reviewer assignment should rotate so that every pair reviews across the cycle.

**Work that spans a rotation:** a feature may not finish within one two-week sprint. Before the rotation, the outgoing responsible pair updates the issue with the current label, spec status and open questions. The incoming responsible pair confirms ownership on the issue.

## 4. Approvals

Two separate approvals apply to every feature. They are not interchangeable.

| Approval | Given by | Gates |
|----------|----------|-------|
| **Specification approval** | Reviewer pair | Start of implementation. Status **Approved** and label `spec:approved`. |
| **Implementation PR approval** | At least one reviewer from another pair | Merge. Approval of the pull request. |

Specification approval does not approve the code. Pull request approval does not replace specification approval.

The specification status is one of the values defined in the [template](../specs/templates/feature-spec-template.md#1-feature-metadata-required):

| Status | Meaning |
|--------|---------|
| Draft | Being written by the responsible pair. Not ready for review. |
| In Review | Submitted to the reviewer pair. Changes may be requested. |
| Approved | Reviewed, with no open blocking comments. Development may start. |

Approval rules:

- Only the reviewer pair sets the status to **Approved**.
- The author must not approve their own specification.
- A specification must not be marked **Approved** while acceptance criteria are ambiguous, a required contract is undefined, or a critical question is open.
- A materially changed approved specification returns to **In Review** and needs renewed cross-pair review (see §7).

## 5. Definition of Ready and Definition of Done

Both definitions are defined in [CONTRIBUTING.md](../CONTRIBUTING.md#10-definition-of-ready). In short:

- **Ready for development:** specification approved, acceptance criteria defined, contracts agreed, dependencies identified and security considerations documented.
- **Done:** implementation matches the specification, tests and CI checks pass, a reviewer from another pair has approved, documentation is updated, and the feature is merged into `main`.

A feature reaches `done` once its implementation is merged and all Definition of Done requirements are satisfied. The weekly Friday demo is separate. It is not part of the Definition of Done, and a feature is not held back from `done` to wait for it.

## 6. GitHub Issue States

GitHub Issues have only *open* and *closed* by default. The team tracks the workflow states with the six agreed labels below. The issue form applies `spec:draft` on creation. The other labels must exist in the repository.

| State | Label | Issue open? | Set when |
|-------|-------|-------------|----------|
| Requirement captured | `spec:draft` | Yes | Issue is created and the spec is being written (Draft) |
| Awaiting review | `spec:in-review` | Yes | Spec is submitted for cross-pair review (In Review) |
| Ready for development | `spec:approved` | Yes | Spec is approved and the Definition of Ready is met |
| In development | `in-progress` | Yes | Implementation branch is started |
| In code review | `in-review` | Yes | Pull request is open for review |
| Done | `done` | No | Implementation is merged and all Definition of Done requirements are met. Issue is closed. |

Keep exactly one state label on an issue at a time. Update the label when the state changes.

## 7. Handling Requirement Changes

Requirements change. The following rules keep specification and code in step.

**Who decides:** the mentor approves changes to requirement priorities and significant scope changes. Interns may clarify and update a specification without mentor approval, unless the change alters priority or scope significantly.

**What needs renewed review:** a material change to an approved specification. Material changes include changes to acceptance criteria, functional requirements, contracts, scope, dependencies or security considerations. Editorial changes, such as typo fixes or formatting, do not need renewed review.

**Before implementation starts:**

1. Update the specification.
2. If the change is significant, obtain mentor approval.
3. Set its status to **In Review** and set the issue label to `spec:in-review`.
4. Obtain review from the reviewer pair again, before development continues.

**After approval, before merge:**

1. Stop work on the affected part of the implementation.
2. If the change is significant, obtain mentor approval.
3. Update the specification and return it to **In Review**. The issue label moves to `spec:in-review`.
4. Obtain re-review. Affected contracts must be agreed again with dependent pairs.
5. Resume implementation only after approval.

**After merge:**

1. Raise a new issue for the change, linked to the original feature ID.
2. Follow the stages in §2 for the change, from stage 2 onward.

Do not change code to match a new requirement without first updating and approving the specification.

## 8. Quick Reference

```text
Requirement (issue, spec:draft; priority set by mentor)
  → Spec Draft → In Review (spec:in-review) → Approved (spec:approved)   [specification approval]
  → Implement on feature branch (in-progress)
  → Tests and CI pass → Pull request (in-review)
  → Implementation PR approval (another pair) → Merge → Issue closed (done)
  → Weekly Friday demo (not a gate)
```

Material specification changes after approval return the feature to **In Review**.
