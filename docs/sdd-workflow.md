# Specification-Driven Development Workflow

This document describes how GoJudge features move from requirement to demo. It covers the lifecycle, responsibilities, and how requirement changes are handled.

It is a companion to the other project documents:

- [CONTRIBUTING.md](../CONTRIBUTING.md): branching, review, Definition of Ready and Definition of Done. Rules are not repeated here.
- [specs/templates/feature-spec-template.md](../specs/templates/feature-spec-template.md): the specification structure. Use it for every feature.

## 1. Principles

- **Specify before you build.** No implementation starts until the feature specification is approved.
- **Contracts first.** API and data contracts are agreed before dependent work starts.
- **Cross-pair review.** Specifications and pull requests are reviewed by a pair other than the author's.
- **Traceable work.** Every feature links a GitHub Issue, a specification and a pull request.

## 2. Lifecycle

| # | Stage | Owner | Output | Gate to next stage |
|---|-------|-------|--------|--------------------|
| 1 | Requirement | Product owner or team | Requirement described in an issue, linked to an MVP feature ID (F01–F12) | Issue exists and feature ID is assigned |
| 2 | Feature Specification | Responsible pair | Spec in `specs/features/` based on the template | Spec is complete enough for review |
| 3 | Technical Design and Contracts | Responsible pair, with dependent pairs | Design, contracts and open questions recorded in the spec | Contracts agreed with dependent pairs |
| 4 | Cross-Pair Specification Review | Reviewer pair | Review outcome in the spec's review table | No open blocking comments |
| 5 | Approval | Reviewer pair | Spec status set to **Approved** | Definition of Ready met (see §5) |
| 6 | Implementation | Responsible pair | Code on a feature branch, following the approved spec | Tests written and passing locally |
| 7 | Automated Verification | Responsible pair | Tests and CI checks passing; acceptance criteria verified | All configured CI checks pass |
| 8 | Pull Request | Responsible pair | Pull request referencing the issue and spec | Pull request ready for review |
| 9 | Code Review | Reviewer pair | Approval from at least one reviewer from another pair | Approval given; comments resolved |
| 10 | Merge | Responsible pair | Changes integrated into `main` | Definition of Done met (see §5) |
| 11 | Demo | Responsible pair | Feature demonstrated to the team | Demo recorded in the issue |

Stages 6–10 follow the branching, pull request and review rules in [CONTRIBUTING.md](../CONTRIBUTING.md#3-git-branching-strategy).

## 3. Responsibilities in Rotating Pairs

The team of eight interns works in four pairs. Pairs rotate, so each developer works with different partners over the 12 weeks.

- **Responsible pair:** owns the feature from specification to demo. Writes the specification, the implementation and the tests. Keeps the issue and spec status current.
- **Reviewer pair:** always a pair other than the responsible pair. Reviews the specification, approves it, and reviews the pull request.
- **Dependent pair:** any pair whose work depends on a contract in the specification. Agrees the contract before dependent development starts.

A pair is never the reviewer of its own work. Reviewer assignment should rotate so that every pair reviews across the cycle.

## 4. Specification Approval

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
- Any change to an approved specification returns it to **In Review** (see §7).

## 5. Definition of Ready and Definition of Done

Both definitions are defined in [CONTRIBUTING.md](../CONTRIBUTING.md#10-definition-of-ready). In short:

- **Ready for development:** specification approved, acceptance criteria defined, contracts agreed, dependencies identified and security considerations documented.
- **Done:** implementation matches the specification, tests and CI checks pass, a reviewer from another pair has approved, documentation is updated, and the feature is in `main`.

Demo is not part of the Definition of Done in CONTRIBUTING, but the team runs a demo for every completed feature (see §6, stage 11).

## 6. GitHub Issue States

GitHub Issues have only *open* and *closed* by default. The team tracks the finer states with labels. The label names below are proposed and must be created in the repository before use.

| State | Label | Issue open? | Set when |
|-------|-------|-------------|----------|
| Requirement captured | `spec:draft` | Yes | Issue is created and the spec is being written (Draft) |
| Awaiting review | `spec:in-review` | Yes | Spec is submitted for cross-pair review (In Review) |
| Ready for development | `spec:approved` | Yes | Spec is approved and the Definition of Ready is met |
| In development | `in-progress` | Yes | Implementation branch is started |
| In code review | `in-review` | Yes | Pull request is open for review |
| Done | `done` (then close) | No | Pull request is merged and the demo is recorded |

Keep exactly one state label on an issue at a time. Update the label when the state changes.

## 7. Handling Requirement Changes

Requirements change. The following rules keep specification and code in step.

**Before implementation starts:**

1. Update the specification.
2. Set its status to **In Review** and set the issue label to `spec:in-review`.
3. Obtain review from the reviewer pair again, before development continues.

**After approval, before merge:**

1. Stop work on the affected part of the implementation.
2. Update the specification and return it to **In Review**. The issue label moves to `spec:in-review`.
3. Obtain re-review. Affected contracts must be agreed again with dependent pairs.
4. Resume implementation only after approval.

**After merge:**

1. Raise a new issue for the change, linked to the original feature ID.
2. Follow the stages in §2 for the change, from stage 2 onward.

Do not change code to match a new requirement without first updating and approving the specification.

## 8. Quick Reference

```text
Requirement (issue, spec:draft)
  → Spec Draft → In Review (spec:in-review) → Approved (spec:approved)
  → Implement on feature branch (in-progress)
  → Tests and CI pass → Pull request (in-review)
  → Cross-pair approval → Merge → Demo → Issue closed (done)
```

Specification changes after approval return the feature to **In Review**.
