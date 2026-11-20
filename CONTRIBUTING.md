# Contributing to GoJudge

This guide applies to every intern on the GoJudge team. It defines how work is specified, branched, reviewed and merged. For project context, see the [README](README.md). The full lifecycle is in [docs/sdd-workflow.md](docs/sdd-workflow.md).

Section numbers are stable. Other documents link to them.

## 1. Contribution Principles

- **Choose the right work item type.** Every issue is one of four types (§2). The type sets what must exist before work starts.
- **Specify before you build.** No feature is implemented without an approved Feature Specification. Spikes, technical tasks and contracts follow the rules in §2 and §3.
- **Keep changes small and focused.** Smaller pull requests are reviewed faster and with better results.
- **Own what you submit.** You are responsible for every line you contribute, including AI-assisted code.
- **Protect secrets and the shared baseline.** Never commit credentials, and never push directly to `main`.

## 2. SDD Rules and Work Item Types

Every GitHub Issue is one of the following types. The type is chosen before the issue is created.

| Type | Use for | Required before work starts | Artefact and location | Branch prefix | Issue form |
|------|---------|-----------------------------|-----------------------|---------------|------------|
| **Feature** | MVP behaviour (F01–F09) | Approved Feature Specification | Specification in [`specs/features/`](specs/features/), based on the [template](specs/templates/feature-spec-template.md) | `feature/` | Feature request (applies `type:feature`) |
| **Technical task** | Infrastructure, tooling, documentation and ADRs | Numbered acceptance criteria in the issue. No Feature Specification. | Pull request. ADRs in [`docs/architecture/adr/`](docs/architecture/adr/) as `ADR-XXX.md`. | `infra/` for infrastructure and tooling. `docs/` for documentation and ADRs. | Technical task (no fixed label; the author selects `type:technical-task`, `type:contract` or `type:adr`) |
| **Technical spike** | Time-boxed investigation that answers one question | Spike Brief in the issue. No Feature Specification. | Report in [`docs/spikes/`](docs/spikes/) as `SPK-XX.md` | `spike/` for experiment code (not merged). `docs/` for the report. | Technical spike (applies `type:spike`) |
| **Contract** | Agreed API, message, data or state contract | Draft contract naming the dependent pairs | Contract in [`docs/contracts/`](docs/contracts/) as `CON-XX.md` | `docs/` | Technical task, with work type "Contract" (the author selects `type:contract`) |

Rules:

- **Features require an approved Feature Specification before implementation.** Implementation means any change to product code that delivers feature behaviour.
- **Technical spikes require a Spike Brief, not a Feature Specification.**
- **Technical tasks require clear acceptance criteria, not a Feature Specification.**
- **Contracts require cross-pair review before dependent implementation begins.** Every dependent pair must agree the contract first.
- **API and data contracts must be agreed before dependent development begins.**
- **Specifications and contracts must be updated when requirements change.** Code, specification and contract must not drift apart.
- **The mentor owns requirement priorities and approves significant scope changes.**
- **Bugfixes** use `fix/`. A bugfix that changes behaviour defined in an approved specification follows the Feature rules for that change. Otherwise it follows the Technical task rules. *(Proposed; to be confirmed by the mentor.)*

## 3. Git Branching Strategy

| Prefix | Convention | Example | Use |
|--------|------------|---------|-----|
| `feature/` | `feature/<feature-id>-<short-description>` | `feature/F06-submit-solution` | Feature implementation, after specification approval |
| `fix/` | `fix/<short-description>` | `fix/submission-status-race` | Bugfixes (see §2) |
| `infra/` | `infra/<short-description>` | `infra/local-cluster` | Infrastructure and tooling technical tasks |
| `spike/` | `spike/<spike-id>-<short-description>` | `spike/SPK-01-job-isolation` | Spike experiment code. Not merged to `main`. |
| `docs/` | `docs/<short-description>` | `docs/spk-01-report` | Documentation, specifications, contracts, ADRs and spike reports |

- Branch from `main`.
- Do not push directly to `main`.
- Keep each branch focused on one work item.
- Specification pull requests use `docs/`. Feature implementation uses `feature/`, and starts only after the specification is approved.

## 4. GitHub Issue Workflow

1. Choose the work item type (§2). Open the matching issue form from [`.github/ISSUE_TEMPLATE/`](.github/ISSUE_TEMPLATE/). Apply exactly one `type:*` label that matches the work item type. The feature form applies `type:feature` and the spike form applies `type:spike`. The technical task form covers technical tasks, contracts and ADR preparation, so it applies no fixed label. The author selects `type:technical-task`, `type:contract` or `type:adr`.
2. Link the issue to its identifier: the feature ID (for example, F06), or the task, spike or contract ID (for example, INF-03, SPK-01, CON-01). The mentor sets the priority of features.
3. Prepare the artefact for the type:
   - **Feature:** write the specification in `specs/features/` and link it from the issue.
   - **Technical spike:** publish the Spike Brief on the issue before the timebox starts.
   - **Technical task:** write the acceptance criteria on the issue.
   - **Contract:** draft the contract in `docs/contracts/` and name the dependent pairs.
4. Workflow status is recorded in the artefact (Draft, In Review, Approved, Provisional), not in GitHub labels. The lifecycle for each type is in [docs/sdd-workflow.md §6](docs/sdd-workflow.md#6-github-issue-states):
   - **Feature:** specification Draft → In Review → Approved (recorded in the specification), then implementation and pull request. The Issue is open until the Definition of Done (§11) is met.
   - **Technical task, spike or contract:** the Issue is open while the pull request is in review, and closed when the Definition of Done (§11) is met.
5. Close the issue once the Definition of Done (§11) for its scope is met. A closed draft-only Issue does not mean that the specification, contract or ADR is formally approved. Its closing comment must say that the draft is complete and the document is not formally approved.

Labels classify the work item type only: `type:feature`, `type:contract`, `type:technical-task`, `type:spike` and `type:adr`. Each new Issue has exactly one. Status labels are not used. Historical status labels may remain on older issues, but are not prescribed for new work.

## 5. Review Processes

### 5.1 Specification review (features)

- The author's pair requests review from at least one pair other than its own.
- Reviewers check that the specification covers acceptance criteria, contracts, dependencies and security considerations.
- The specification is approved only when reviewers have no open blocking comments.
- Specification approval is separate from pull request approval. Approving a specification does not approve the code, and approving a pull request does not replace specification approval.
- A material change to an approved specification requires renewed cross-pair review before implementation continues. Material changes include acceptance criteria, functional requirements, contracts, scope, dependencies and security considerations. Editorial changes do not need renewed review.

### 5.2 Contract review

- The author's pair drafts the contract in `docs/contracts/CON-XX.md` and opens a pull request.
- At least one reviewer from another pair reviews the pull request. This is required.
- Every dependent pair records its agreement on the issue or the pull request.
- Dependent implementation starts only after every dependent pair has agreed.
- A change to an agreed contract requires renewed cross-pair review.

### 5.3 Spike Briefs and reports

- The Spike Brief is published on the issue before the timebox starts. It states the hypothesis, timebox, experiment, success and failure criteria, and required evidence.
- The report is written in `docs/spikes/SPK-XX.md`. Each success criterion is recorded as Verified, Partially verified or Not verified, with its evidence.
- The report pull request is reviewed by at least one pair other than the spike's owner.
- A timebox is not extended without mentor approval.

## 6. Implementation Guidelines

- Implement only what the approved specification, the approved contract, or the accepted acceptance criteria describe.
- Do not introduce unapproved dependencies or architectural changes. Raise them as an issue first, and record architectural decisions as ADRs in `docs/architecture/adr/`.
- Do not commit secrets, credentials or environment-specific configuration.
- Treat learner-submitted code as untrusted. Any feature that executes it must follow the isolation requirements in the specification.

## 7. Testing Requirements

- New behaviour must be covered by automated tests.
- Tests must pass locally before a pull request is opened.
- All configured CI checks must pass before merging.
- Test commands and CI configuration will be documented in the README once they exist.
- Items that change no code (specifications, contracts, spike reports, ADRs) provide the evidence listed in the pull request template instead.

## 8. Pull Request Process

- Use the [pull request template](.github/pull_request_template.md) for every work item type. Tick the work item type first.
- Reference the GitHub Issue and the type-specific artefact: the approved specification (feature), the acceptance criteria (technical task), the Spike Brief and report (spike), or the contract and agreement (contract).
- Keep the pull request focused and reasonably small. If it grows, split it.
- Mark the pull request ready for review only when its evidence is complete.

## 9. Code Review Rules

- **At least one reviewer from another pair must approve every pull request**, whatever its work item type. This is required.
- Additional reviewers are optional.
- Reviewers check the change against its artefact: the specification, the contract, the acceptance criteria or the Spike Brief, as applicable. They also check test coverage, readability and security.
- Authors respond to every comment, either by changing the work or by explaining why not.
- Authors must not merge their own pull requests.

## 10. Definition of Ready

A work item is ready to start when the checks for its type are met.

| Type | Ready when |
|------|------------|
| **Feature** | Specification is documented and approved. Acceptance criteria are defined. Relevant API and data contracts are agreed. Dependencies are identified. Security considerations are documented. The mentor's priority is set. |
| **Technical task** | Objective, deliverables and numbered acceptance criteria are on the issue. Owner pair, estimate and dependencies are named. Decisions the task depends on are recorded, marked as provisional where they are not yet approved. |
| **Technical spike** | Spike Brief is on the issue: hypothesis, timebox, experiment, success and failure criteria, and required evidence. Decisions the spike tests are named and marked as provisional. The timebox dates are set. |
| **Contract** | Draft contract is in `docs/contracts/CON-XX.md`, covering its states, fields, messages or errors as relevant. Dependent pairs are named and have been told when review is expected. |

## 11. Definition of Done

A work item is done when the checks for its type are met. Then its issue is closed. No status label is applied.

| Type | Done when |
|------|-----------|
| **Feature** | Implementation matches the approved specification. Tests pass and all configured CI checks pass. A reviewer from another pair has approved. Documentation is updated. The feature is merged into `main`. |
| **Technical task** | Every acceptance criterion is met and its evidence is recorded in the pull request. CI passes where code or configuration changed. A reviewer from another pair has approved. Documentation is updated. No secrets or environment-specific configuration are committed. The change is merged into `main`. |
| **Technical spike** | The timebox is respected. Every success criterion is recorded with its evidence. The outcome is stated (Go, Conditional go, No-go or Inconclusive). Failure criteria met during the spike have been reported to the mentor. The report is merged to `docs/spikes/SPK-XX.md` and reviewed by a pair other than the owner. Spike code is not merged to `main`. Decisions that follow from the outcome are listed. |
| **Contract** | Every dependent pair has recorded its agreement. A reviewer from another pair has approved the pull request. The contract is merged to `docs/contracts/`. Dependent specifications and issues link to it. |

A feature reaches `done` once its implementation is merged and all Definition of Done requirements are satisfied. The weekly Friday demo is not a Definition of Done requirement and is not a prerequisite for merging or for `done`.

## 12. AI-Assisted Development Policy

- AI tools may be used to assist development.
- The responsible developer must understand, test and review all AI-generated code before submitting it.
- AI-generated code is held to the same standards as any other contribution.
- AI output must not be used to bypass specification review, contract agreement, dependency approval or the secrets rules in this guide.
