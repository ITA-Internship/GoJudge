# Contributing to GoJudge

This guide applies to every intern on the GoJudge team. It defines how work is specified, branched, reviewed and merged. For project context, see the [README](README.md).

## 1. Contribution Principles

- **Specify before you build.** No feature is implemented without a reviewed specification.
- **Keep changes small and focused.** Smaller pull requests are reviewed faster and with better results.
- **Own what you submit.** You are responsible for every line you contribute, including AI-assisted code.
- **Protect secrets and the shared baseline.** Never commit credentials, and never push directly to `main`.

## 2. Specification-Driven Development Rules

- Every feature starts with a GitHub Issue.
- Every feature requires a Markdown specification based on the [feature specification template](specs/templates/feature-spec-template.md).
- Specifications must be reviewed by another pair before implementation starts.
- API and data contracts must be agreed before any dependent development begins.
- Specifications must be updated when requirements change. Code and specification must not drift apart.
- The mentor owns requirement priorities and approves significant scope changes.

The full workflow is described in [docs/sdd-workflow.md](docs/sdd-workflow.md).

## 3. Git Branching Strategy

| Type | Convention | Example |
|------|------------|---------|
| Feature | `feature/<feature-id>-<short-description>` | `feature/F06-submit-solution` |
| Bugfix | `fix/<short-description>` | `fix/submission-status-race` |

- Branch from `main`.
- Do not push directly to `main`.
- Keep each branch focused on one feature or fix.

## 4. GitHub Issue Workflow

1. Open an issue using the [feature issue template](.github/ISSUE_TEMPLATE/). The template applies the `spec:draft` label.
2. Link the issue to the feature ID (for example, F06) and the MVP scope. The mentor sets the priority.
3. Write the specification in `specs/features/` and link it from the issue. Set the label to `spec:in-review` when you submit it for review.
4. Set the label to `spec:approved` only once the specification is approved and the [Definition of Ready](#10-definition-of-ready) is met.
5. Set the label to `in-progress` when the implementation branch is started, and to `in-review` when the pull request is open.
6. Set the label to `done` and close the issue once the feature is merged and the [Definition of Done](#11-definition-of-done) is met.

Keep exactly one state label on an issue at a time. The full state model is in [docs/sdd-workflow.md](docs/sdd-workflow.md#6-github-issue-states).

## 5. Specification Review Process

- The author's pair requests review from at least one other pair.
- Reviewers check that the specification covers acceptance criteria, contracts, dependencies and security considerations.
- The specification is approved only when reviewers have no open blocking comments.
- Specification approval is separate from pull request approval. Approving a specification does not approve the code, and approving a pull request does not replace specification approval.
- A material change to an approved specification requires renewed cross-pair review before implementation continues. Material changes include acceptance criteria, functional requirements, contracts, scope, dependencies and security considerations. Editorial changes do not need renewed review.

## 6. Implementation Guidelines

- Implement only what the approved specification describes.
- Do not introduce unapproved dependencies or architectural changes. Raise them as an issue first.
- Do not commit secrets, credentials or environment-specific configuration.
- Treat learner-submitted code as untrusted. Any feature that executes it must follow the isolation requirements in the specification.

## 7. Testing Requirements

- New behaviour must be covered by automated tests.
- Tests must pass locally before a pull request is opened.
- All configured CI checks must pass before merging.
- Test commands and CI configuration will be documented in the README once they exist.

## 8. Pull Request Process

- Use the [pull request template](.github/pull_request_template.md).
- Reference the GitHub Issue and the approved specification in the description.
- Keep the pull request focused and reasonably small. If it grows, split it.
- Mark the pull request ready for review only when tests pass and the description is complete.

## 9. Code Review Rules

- At least one reviewer from another pair must approve the pull request. This is implementation approval, not specification approval.
- Reviewers check the implementation against the specification, test coverage, readability and security.
- Authors respond to every comment, either by changing the code or by explaining why not.
- Authors must not merge their own pull requests.

## 10. Definition of Ready

A feature is ready for development when:

- its specification is documented and reviewed;
- its acceptance criteria are defined;
- relevant API and data contracts are agreed;
- its dependencies are identified;
- its security considerations are documented.

## 11. Definition of Done

A feature is complete when:

- the implementation matches the specification;
- tests pass, and all configured CI checks pass;
- code review is approved by a reviewer from another pair;
- documentation is updated;
- the feature is integrated into `main`.

A feature reaches `done` once these requirements are met. The weekly Friday demo is not a Definition of Done requirement and is not a prerequisite for merging.

## 12. AI-Assisted Development Policy

- AI tools may be used to assist development.
- The responsible developer must understand, test and review all AI-generated code before submitting it.
- AI-generated code is held to the same standards as any other contribution.
- AI output must not be used to bypass specification review, dependency approval or the secrets rules in this guide.
