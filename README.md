# GoJudge

## Project Overview

GoJudge is a lightweight, LeetCode-style platform where learners solve Go programming challenges, submit their solutions and receive automated evaluation results.

The project is an educational Go backend development internship. Eight interns work in four rotating pairs over 12 weeks to deliver a working MVP while building professional Junior Go Developer skills.

> **Status:** Early planning. The repository currently contains documentation scaffolding only. See [Project Status](#project-status).

## Project Goals

- Deliver a working GoJudge MVP within the 12-week internship.
- Build practical Go backend skills: HTTP APIs, persistence, asynchronous processing and testing.
- Practise Specification-Driven Development (SDD) as a team, from specification through to review and delivery.
- Keep the codebase, documentation and contribution process professional and consistent.

## MVP Functional Scope

| ID | Feature | Status |
|----|---------|--------|
| F01 | Authentication and role-based access | Planned |
| F02 | Problem catalogue | Planned |
| F03 | Problem details | Planned |
| F04 | Browser-based Go code editor | Planned |
| F05 | Run public sample tests | Planned |
| F06 | Submit solutions asynchronously | Planned |
| F07 | Automated evaluation using Go test suites | Planned |
| F08 | Submission history | Planned |
| F09 | Problem administration | Planned |

Optional features F10–F12 are outside the initial delivery commitment.

## Proposed Technology Stack

The following stack is **proposed** and has not yet been adopted in the repository.

| Area | Technology |
|------|------------|
| Backend | Go, `net/http`, optionally `chi` |
| Database | PostgreSQL, `pgx`, `sqlc` |
| Frontend | React, TypeScript, Vite, Monaco Editor |
| Authentication | Auth0 |
| Messaging | RabbitMQ, Watermill |
| Execution orchestration | Kubernetes Jobs, `client-go` |
| Local infrastructure | k3d or Minikube |
| CI/CD | GitHub Actions |
| API documentation | OpenAPI 3.x |

## High-Level Architecture

The proposed submission workflow is:

```text
React client → Go API → PostgreSQL / RabbitMQ → Judge Orchestrator → Kubernetes Job → evaluation results
```

1. The React client submits a solution through the Go API.
2. The API stores the submission in PostgreSQL and publishes a job message to RabbitMQ.
3. The Judge Orchestrator consumes the message and starts a Kubernetes Job for the evaluation.
4. The Kubernetes Job runs the learner's code against the problem's Go test suite.
5. The evaluation result is recorded and returned to the learner, who can review it in their submission history.

> **Security:** Learner code is untrusted. Running it requires validated execution isolation. Execution must not be implemented or enabled until the isolation approach has been specified, reviewed and validated.

## Development Methodology: SDD

GoJudge follows **Specification-Driven Development**:

- Every feature starts with a GitHub Issue.
- Every feature requires a reviewed Markdown specification before implementation begins.
- Specifications are reviewed by a pair other than the author's.
- API and data contracts are agreed before dependent work starts.

See [docs/sdd-workflow.md](docs/sdd-workflow.md) for the full workflow and [specs/templates/feature-spec-template.md](specs/templates/feature-spec-template.md) for the specification template.

## Team Organisation

- **Team:** 8 Go Developer interns.
- **Structure:** 4 rotating pairs. Pairs review each other's specifications and pull requests.
- **Duration:** 12 weeks.

Contribution rules and review responsibilities are described in [CONTRIBUTING.md](CONTRIBUTING.md).

## Repository Structure

```text
.
├── .github/
│   ├── ISSUE_TEMPLATE/
│   │   └── feature.yml          # Feature issue template
│   └── pull_request_template.md # Pull request template
├── docs/
│   └── sdd-workflow.md          # SDD workflow
├── specs/
│   ├── features/                # Approved feature specifications
│   └── templates/
│       └── feature-spec-template.md
├── .gitignore
├── CONTRIBUTING.md
└── README.md
```

## Getting Started

**To be documented.** The repository does not yet contain application code, a Go module, build configuration or local infrastructure definitions. Setup instructions will be added once these exist and have been verified.

## Development Workflow

1. Create a GitHub Issue for the feature or bug.
2. Write a specification from the [feature specification template](specs/templates/feature-spec-template.md) and place it in `specs/features/`.
3. Obtain review of the specification from another pair, and agree on any contracts.
4. Implement on a feature branch and open a pull request that references the issue and the approved specification.
5. Obtain review and passing checks, then merge.

Full rules are in [CONTRIBUTING.md](CONTRIBUTING.md).

## Testing and Quality Standards

- Automated tests are required for new behaviour.
- All configured CI checks must pass before merging.
- Implementation must match the approved specification.
- Test tooling and CI configuration are **to be documented** once they are added to the repository.

## Documentation

| Document | Purpose |
|----------|---------|
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contribution rules, branching, review and definitions of Ready and Done |
| [docs/sdd-workflow.md](docs/sdd-workflow.md) | Specification-Driven Development workflow |
| [specs/templates/feature-spec-template.md](specs/templates/feature-spec-template.md) | Template for feature specifications |
| [specs/features/](specs/features/) | Feature specifications |
| [.github/ISSUE_TEMPLATE/](.github/ISSUE_TEMPLATE/) | GitHub Issue templates |
| [.github/pull_request_template.md](.github/pull_request_template.md) | Pull request template |

## Project Status

- Repository: documentation scaffolding in place.
- Specifications: none approved yet.
- Application code: not yet started.
- Setup and deployment: to be documented.
