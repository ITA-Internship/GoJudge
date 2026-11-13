<!--
Before opening this pull request:
- Tick one work item type below, and complete the section for that type. Delete the others.
- The GitHub Issue stays open until the Definition of Done for its scope is met. Its label is its work item type (`type:*`), not its status.
- Reference one issue per pull request. Keep the change focused and reasonably small.
- At least one reviewer from a different pair must approve. Additional reviewers are optional.
- The weekly Friday demo is not a prerequisite for merging.
- Do not paste secrets or credentials anywhere in this description.
-->

## Work Item Type

- [ ] Feature
- [ ] Technical task
- [ ] Technical spike
- [ ] Contract

## Related GitHub Issue

Closes #

---

<!-- Feature: complete this section only if the work item is a feature. -->

## Feature Specification

- Specification: `specs/features/F00-short-name.md`
- Feature ID: F00
- Specification status: Approved (implementation pull requests). In Review (specification-only pull requests, which use `docs/`).
- Specification changed in this pull request? Yes / No
  - If yes, link to the renewed review outcome: <!-- A material change requires renewed review before implementation continues. -->

**Acceptance criteria verification** (one row per criterion in the specification)

| AC | Verified by | Result |
|----|-------------|--------|
| AC-01 | <!-- test name, or manual step --> | Pass / Fail / Not applicable |

---

<!-- Technical task: complete this section only if the work item is a technical task. -->

## Acceptance Criteria Verification (technical task)

Criteria are those listed on the issue.

| Criterion | Evidence | Result |
|-----------|----------|--------|
| 1 | <!-- link, command output summary, or observation --> | Met / Not met / Not applicable |

---

<!-- Technical spike: complete this section only if the work item is a technical spike. -->

## Spike Brief and Report (technical spike)

- Spike Brief: on the issue
- Report: `docs/spikes/SPK-XX.md`
- Timebox: <!-- dates and pair-days -->
- Outcome: Go / Conditional go / No-go / Inconclusive
- Spike code: kept on a `spike/` branch and not merged to `main`

**Success criteria** (one row per criterion in the Spike Brief)

| Criterion | Status | Evidence |
|-----------|--------|----------|
| SC-01 | Verified / Partially verified / Not verified | <!-- raw record location --> |

**Failure criteria checked:** <!-- list each failure criterion and whether it was met -->

---

<!-- Contract: complete this section only if the work item is a contract. -->

## Contract and Agreement (contract)

- Contract: `docs/contracts/CON-XX.md`
- Dependent pairs named in the contract: <!-- pairs -->

**Agreement**

| Dependent pair | Agreed? | Recorded on | Date |
|----------------|---------|-------------|------|
| <!-- pair --> | Yes / No | Issue / pull request link | YYYY-MM-DD |

---

## Summary

<!-- What changed, and how it relates to its artefact: the specification, the acceptance criteria, the Spike Brief or the contract. Note any deviation and the reason. -->

-

## Testing and Automated Checks

<!-- For code or configuration changes: list the tests added or changed, and link the CI result. Do not paste secrets. For changes that contain no code (specifications, contracts, spike reports, ADRs), write "Not applicable: evidence is in the section above." -->

- Tests added or changed:
- CI result:

## Documentation Changes

<!-- List documentation updated in this pull request, or write "None" with a reason. -->

- [ ] README updated if setup, structure or status changed
- [ ] Specification updated if requirements changed (material changes need renewed review)
- [ ] Contract updated if an interface changed (changes after agreement need renewed review)
- [ ] ADR added in `docs/architecture/adr/` if an architectural decision was made
- [ ] Other:

## Reviewer Checklist

At least one reviewer from a different pair must approve. Additional reviewers are optional.

- [ ] The work item type is ticked and the matching section is complete.
- [ ] The linked issue and its artefact are correct and current.
- [ ] The change matches its artefact: the approved specification, the acceptance criteria, the Spike Brief or the agreed contract.
- [ ] Every acceptance criterion, success criterion or contract point has its evidence.
- [ ] Tests cover new behaviour and pass in CI, where code changed.
- [ ] No secrets, credentials or environment-specific configuration are committed.
- [ ] No unapproved dependencies or architectural changes are introduced. Architectural decisions have an ADR.
- [ ] Learner-submitted code, if executed, follows the isolation requirements in the specification.
- [ ] Documentation changes are accurate.
- [ ] Any material change to an approved specification or agreed contract has renewed review.
