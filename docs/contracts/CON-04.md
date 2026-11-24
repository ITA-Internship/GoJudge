# CON-04: Problem packaging, starter template and test visibility

| Field | Value |
|-------|-------|
| Contract ID | CON-04 |
| Status | Draft |
| Owner | Pair A |
| Reviewer pair | Pair D (provisional, `sprint-02.md` §9) |
| Dependent pairs | Pair C (F07-04 and later consumers), Pair D (CON-05 entities, SPEC-F04), Pair B (CON-02, F07-01) |
| GitHub Issue | [#31](https://github.com/ITA-Internship/GoJudge/issues/31) |
| Blocking decision | DEC-09 (MA-14), **Open** |
| Last updated | 2026-11-24 |

> **Draft.** This document is the Sprint 2 draft of CON-04. It is **not** an approved contract. DEC-09 is open, and no DEC-09 option is selected here (section 6). Values marked **Stated** come from an approved or merged source that is named next to them. Values marked **Proposed** are suggestions for review and are not decided. Everything else is an open question (section 8). Final approval is planned for Sprint 3, after DEC-09 is decided and every dependent pair has agreed (section 9).

---

## 1. Purpose and Scope

This contract defines three things that several features share:

1. **Problem package:** what a problem consists of, and where its parts are stored (section 3).
2. **Starter template:** the code a learner sees in the editor when they open a problem (section 4).
3. **Test visibility:** which tests and test results a learner may see, and which must stay hidden (section 5).

**In scope:** the logical content of a problem, the starter template format, the visibility rules for sample and hidden tests, and the DEC-09 options.

**Out of scope:** database tables and columns (CON-05), HTTP endpoints and their JSON (SPEC-F03, SPEC-F09, CON-03), message formats (CON-02), Kubernetes Job settings (F07-02a, DEC-08), and the verdict model (CON-01). This contract states what these documents must respect, not how they implement it.

---

## 2. Terms

| Term | Meaning in this contract |
|------|--------------------------|
| Problem | One coding challenge that learners can solve. It has a draft or published status (F09-01). |
| Statement | The text that describes the task to the learner. |
| Submission contract | What the learner's code must provide: a full program, a function with a given signature, or an HTTP handler. |
| Starter template | Go source text that the editor shows when a learner opens a problem with no preserved work (SPEC-F04 FR-10). |
| Test case | One check of the learner's code. |
| Sample test | A test case that is public. Its content may be shown to learners and is used by "Run samples" (F05). |
| Hidden test | A test case that is private. Only the judge (F07-04) and administrators (F09-02) may access its content. |
| Test asset | Any file or data needed to run test cases: test code, inputs, expected outputs, fixtures. |
| Hidden content | The content of hidden tests: test code, inputs, expected outputs, fixtures, test names and the output produced while they run. |
| Learner-facing | Any API response, page, message or log line that a learner can see or request. |

---

## 3. Problem Package

### 3.1 What is already fixed

| # | Value | Status | Source |
|---|-------|--------|--------|
| P-01 | Only Go solutions are supported in the MVP. | Stated | MVP proposal §3, "Problem Format and Evaluation Model" |
| P-02 | Each problem defines a submission contract, starter code, a predefined test suite and the expected evaluation behaviour. | Stated | MVP proposal §3, "Problem Format and Evaluation Model" |
| P-03 | Submitted code is combined with trusted test assets and evaluated with `go test` in an isolated environment. | Stated | MVP proposal §3 and §6; CON-01 §1 (`RUNNING`: "executing the Go test suite") |
| P-04 | The release must include at least one algorithmic challenge and one HTTP-handler challenge. | Stated | MVP proposal §11 |
| P-05 | A problem has a draft or published status. Learners only see published problems. | Stated | Product backlog F02-01, F09-01; SPEC-F06 FR-03 |
| P-06 | Each problem has sample tests and hidden tests. Administrators upload, store and replace them. | Stated | MVP proposal §2 and §3 (F05, F09); product backlog F09-02 |
| P-07 | The submission message carries IDs only. The judge side loads source code and tests from PostgreSQL. | Stated (pending DEC-04 / DEC-12) | CON-02 §3 (`docs/contracts/CON-02.md`, merged; DEC-04 and DEC-12 not yet approved). Any DEC-09 option that stores tests outside PostgreSQL needs a CON-02 change (section 6.1). |

### 3.2 Logical content of a problem

This table lists the parts of a problem and who uses them. It defines content, not storage. Storage is DEC-09 (section 6.1), and the table and column names belong to CON-05.

| Part | Required | Used by | Status |
|------|----------|---------|--------|
| ID | Yes | F02, F03, F06 (`problem_id`), F08, CON-05 | **Proposed:** UUID, so that it matches `submissions.problem_id` in CON-05 (preliminary, PR #40) and the examples in SPEC-F06 and SPEC-F08. Closes CON-05 `C05-OQ-21` (PR #40) when agreed. |
| Title | Yes, non-empty | F02-02, F03-01, SPEC-F08 (`problem_title`) | Stated (F02-02 lists the title). Maximum length: C04-OQ-01. |
| Difficulty | Yes | F02-02, F11 | Stated (F02-02). Allowed values: C04-OQ-02. |
| Statement | Yes | F03-01 | Stated (F03-01). Text format (plain text or Markdown) and rendering rules: C04-OQ-03. |
| Input/output contract and constraints | Yes | F03-01 | Stated (F03-01). Whether they are separate fields or part of the statement: C04-OQ-03. |
| Public examples | Yes, at least one | F03-01 | Stated (F03-01). Whether they are the same as the sample tests: C04-OQ-04. |
| Submission contract | Yes | F04, F06, F07 | Stated (P-02). Shape: DEC-09 (section 6.3). |
| Starter template | Optional | F04-01 | Section 4. A problem without a template is valid (SPEC-F04 FR-12). |
| Time limit and memory limit | Yes | F07-02a | Stated: F07-02a sets CPU, memory and time limits per Job. Per problem or global: C04-OQ-05. |
| Sample tests | Yes, at least one | F05, F03-01 | Stated (P-06). Format: DEC-09 (section 6.2). |
| Hidden tests | Yes, at least one | F07-04 | Stated (P-06). Format: DEC-09 (section 6.2). |
| Reference solution | Proposed: yes, admin-only | F09-02, MVP proposal §8 ("deterministic verdicts for agreed reference solutions") | **Proposed.** Used to check a problem before it is published. Never learner-facing. C04-OQ-06. |
| Status (draft or published) | Yes | F02-01, F06, F09-01 | Stated (P-05). |
| Topic tags | No | F11 (Should) | Not needed in the MVP Must scope. C04-OQ-02. |

### 3.3 Rules for every package

| # | Rule | Status |
|---|------|--------|
| P-08 | A problem may be published only if it has a title, a statement, at least one sample test and at least one hidden test. | **Proposed** (needed by F05 and F07-04; enforced by F09-01) |
| P-09 | Changing the tests of a published problem does not change the verdicts of submissions that were already judged. | **Proposed.** Whether a new version of the tests is kept as a separate version: C04-OQ-07. |
| P-10 | Unpublishing a problem does not delete it. Past submissions keep their `problem_id`, and SPEC-F08 keeps showing the title (SPEC-F08 FR-04a). | **Proposed**, consistent with SPEC-F08 (merged). Deleting a problem with submissions: C04-OQ-08. |

---

## 4. Starter Template

### 4.1 What is already fixed

| # | Value | Status | Source |
|---|-------|--------|--------|
| T-01 | The template comes only from problem data. The editor has no built-in template text. | Stated | SPEC-F04 FR-11 |
| T-02 | The template is untrusted text and is shown as plain text only. | Stated | SPEC-F04 FR-14 |
| T-03 | A template that is missing, empty, whitespace-only, contains NUL (U+0000) or is longer than the maximum source size is treated as "no usable template": the editor opens empty with a notice. | Stated | SPEC-F04 FR-12, AC-03 |
| T-04 | The template is delivered to the editor as part of the problem details (F03-01). | Stated | SPEC-F04 §6, product backlog F04-01 |

### 4.2 Format

| # | Value | Status |
|---|-------|--------|
| T-05 | One Go source file, as UTF-8 text without NUL (U+0000), with `\n` line endings. | **Proposed** (matches SPEC-F06 single-file and NUL rules; line-ending handling aligns with SPEC-F04 OQ-14) |
| T-06 | Maximum size: the same limit as `source_code` (SPEC-F04 OQ-09, SPEC-F06 OQ-04, CON-05 `C05-OQ-23` in PR #40), so that a template can always be submitted unchanged. | **Proposed** |
| T-07 | When provided, the template must compile together with the problem's test suite, so that submitting it unchanged gives a test verdict (normally `WRONG_ANSWER`), not `COMPILATION_ERROR`. | **Proposed.** Checked by F09-02 before publishing (C04-OQ-06); optionality per section 3.2 stands. |
| T-08 | The template contains no part of any hidden test, and no hint derived from hidden content. | **Stated as a visibility rule** (V-02). |
| T-09 | The shape of the template (full program, function or HTTP handler, and the package name) follows the submission contract. | Open: DEC-09 (section 6.3). |
| T-10 | The format closes SPEC-F04 OQ-06 once DEC-09 is decided. | Dependency note |

---

## 5. Test Visibility

### 5.1 Rules

| # | Rule | Status |
|---|------|--------|
| V-01 | Every test case is either a sample test or a hidden test, never both. | Stated (P-06) |
| V-02 | **Hidden content is never returned to a learner, in any response or log.** This covers every learner-facing API response (problem details, run samples, submission, submission state, history, error responses), every page, every message a learner can read, and every log line a learner can see. It covers test code, inputs, expected outputs, fixtures, test names and the output produced while hidden tests run. | **Stated** (issue #31, product backlog F07-04, RSK-07) |
| V-03 | Hidden content is written to no application log at any level. Logs refer to tests by ID or index only. | **Stated** as part of V-02; INF-07 logging rules apply. |
| V-04 | Only the judge (F07-04) and administrator-only endpoints (F09-02) may read hidden test assets. No learner-facing endpoint reads them. | **Stated** (product backlog F07-04: "available only to the judge Job"; F09-01: admin endpoints restricted to administrators) |
| V-05 | Output from a hidden-test run (stdout, stderr, `go test` output, panic messages) is not returned to the learner (except filtered compilation diagnostics if permitted under V-08 / C04-OQ-10), because the learner's code can print the hidden inputs it receives. Only the CON-01 verdict and the details allowed by section 6.4 are returned. | **Stated** as part of V-02 |
| V-06 | The learner's code must not be able to read hidden test assets during execution beyond what the test run itself passes to it (for example, by reading test files from the filesystem and printing them). | **Stated as a requirement for F07-04 and DEC-08.** How it is enforced: C04-OQ-09. Tested by F07-07 ("hidden test reads"). |
| V-07 | Sample tests and their results may be shown to the learner, including inputs, expected outputs and the program output, bounded by the output limit (`OUTPUT_LIMIT_EXCEEDED`, CON-01 §2). | **Stated** (MVP proposal F05: "per-case feedback") |
| V-08 | A compilation error message may be shown to the learner only if it cannot contain hidden test code. | **Proposed.** Depends on how the code is combined with the tests (DEC-09, section 6.2). C04-OQ-10. |
| V-09 | An administrator may view and download hidden test assets only through administrator-only endpoints. | **Stated** (F09-01, F09-02) |
| V-10 | The data in the `results` table (CON-05 §3.4, preliminary PR #40) that a learner may see follows section 6.4. Data that may contain hidden content is never placed in a learner-facing response. | Stated rule; storage is CON-05 `C05-OQ-26` (PR #40). |

### 5.2 How the rules are checked

| Rule | Check | Owner |
|------|-------|-------|
| V-02, V-04 | API tests: every learner-facing endpoint is called for a problem with hidden tests, and the response is searched for the hidden inputs, expected outputs and test names. | F03, F05, F06, F08 implementation pull requests |
| V-03 | Log review: logs from a full judge run are searched for hidden inputs and expected outputs (RSK-07). | F07-04, INF-07 |
| V-05 | Judge test: a solution that prints its input fails a hidden test; the learner-facing result contains only the verdict and the allowed details. | F07-04, F07-07 |
| V-06 | Adversarial test: a solution tries to open and print test files. | F07-07 |
| V-07 | F05 test: sample input, expected output and program output are returned for sample tests only. | F05 |

---

## 6. DEC-09 Options

**Status: Open.** DEC-09 is decided by the mentor (MA-14). This section lists the options and their trade-offs. **No option is selected in this draft.** DEC-09 has four connected parts.

### 6.1 Where problems and test assets are stored

| Option | Description | For | Against |
|--------|-------------|-----|---------|
| A. Rows in PostgreSQL | Problem fields and each test case are stored as rows (text columns). | Matches CON-02 §3 (tests loaded from PostgreSQL). One storage, one backup, transactions with publishing. Simple admin forms (F09-03). | Large test data makes the table big. Test code for HTTP-handler problems does not fit "input and output" rows well. |
| B. Package archive in PostgreSQL | The admin uploads one archive (for example, `.zip` or `.tar.gz` with `problem.yaml`, `statement.md`, `starter.go`, `tests/`). It is stored as a binary column and unpacked by the judge. | Keeps CON-02 §3. One upload contains everything; a package can be versioned as a whole. Easy to prepare and test locally. | Needs archive validation (size, paths, file types) against malicious uploads. Problem fields must still be copied into columns for the catalogue. |
| C. Files in object storage | Same package as B, stored in object storage (for example, a bucket); PostgreSQL keeps metadata and a reference. | Good for large test data. | New infrastructure and credentials, conflicts with the MVP stack (MVP proposal §5). Needs a CON-02 change, because tests are no longer loaded from PostgreSQL. |
| D. Files in the Git repository | Problems live in a folder of the repository and are loaded on deployment. | Review of problems through pull requests; no upload code. | Hidden tests become visible to everyone with repository access. Does not meet F09-02 (admin upload). |

### 6.2 Test format

| Option | Description | For | Against |
|--------|-------------|-----|---------|
| A. Input/output pairs | Each test case is an input text and an expected output text. A fixed harness runs the program and compares the output. | Simple for algorithmic problems; easy admin forms; per-case feedback for F05 is natural. | Does not cover HTTP-handler challenges (P-04). The proposal asks for "Go testing framework instead of stdin/stdout-only evaluation" (MVP proposal §12, Proposed). |
| B. Go test files per problem | The admin writes `_test.go` files that call the learner's code. The judge runs `go test`. | Covers functions, components and HTTP handlers (P-02, P-04) with the standard library (`net/http/httptest`). Matches P-03. | Test files are compiled alongside the learner's code (same package or external `<pkg>_test` package), so V-06 and V-08 still require isolating test sources from runtime access and filtering compiler diagnostics. Admins must write Go tests. Mapping individual test results to per-case feedback needs `go test -json`. |
| C. Hybrid | Test cases are data (input and expected output, or request and expected response); one Go test harness per problem type reads them and runs them with `go test`. | Covers both challenge types; admins mostly edit data; per-case results are easy to report. | More design work; harnesses for each challenge type must be written and reviewed. |

### 6.3 Submission contract and template shape

| Option | Description | For | Against |
|--------|-------------|-----|---------|
| A. Full program | The learner writes `package main` with `func main()`; tests run the program. | Familiar from competitive programming; works with 6.2 A. | Does not work for HTTP handlers or functions; tests cannot call the code directly. |
| B. Function, type or handler in a fixed package | The learner implements a given signature (for example, `func Solve(nums []int, target int) []int` or an `http.Handler`) in a package named by the problem. | Works with 6.2 B and C; natural for HTTP-handler challenges. | The package name and signature become part of each problem and must be shown in the statement and template. |
| C. Per problem | Each problem chooses A or B. | Flexible. | Two judge paths to build and test. |

### 6.4 Feedback after a submission

| Option | What the learner sees after hidden tests run | For | Against |
|--------|----------------------------------------------|-----|---------|
| A. Verdict only | The CON-01 verdict. | Simplest; no risk of leaking hidden content. | Learner gets little guidance. |
| B. Verdict and counts | The verdict and the number of passed and total hidden tests. | Some guidance, no hidden content. | Reveals the number of hidden tests. |
| C. Verdict, counts and first failed test index | Adds the index (for example, "test 4") of the first failed hidden test, without its content. | More guidance. | Allows probing a single hidden test by repeated submissions. |
| D. Per problem | Each problem sets one of A–C. The MVP proposal §3 says "Detailed test feedback may be returned according to the visibility rules defined for each problem". | Follows the proposal text. | Adds a field to every problem and a branch to every learner-facing response. |

None of A–D allows hidden content to be returned (V-02).

---

## 7. Dependent Work

| Dependent | What it needs from CON-04 | Pair (per `sprint-02.md` §9, provisional) |
|-----------|---------------------------|--------------------------------------------|
| SPEC-F03, F03-01 (problem details) | Fields of section 3.2; public examples; V-02 | To be confirmed by the mentor |
| SPEC-F04, F04-01 (editor) | Section 4; closes SPEC-F04 OQ-06 | Pair D (SPEC-F04 author) |
| F05 (run samples) | Sample tests, V-07, 6.2 | To be confirmed by the mentor |
| F07-04 (hidden test assets, judge) | 6.1, 6.2, V-02 to V-06 | Pair C (`sprint-02.md` §9: "F07-04 later") |
| SPEC-F09, F09-02 (admin problem and test management) | 3.2, 6.1, 6.2, P-08, P-09, V-09 | To be confirmed by the mentor |
| CON-05 (data model, preliminary PR #40) | 6.1, section 3.2; closes `C05-OQ-21`, input to `C05-OQ-26` | Pair D |
| CON-02 (messaging) | P-07 must stay true, or CON-02 changes (6.1 option C) | Pair B |
| SPEC-F08 (history) | Title (3.2), P-10, 6.4; closes SPEC-F08 OQ-04 and the `problems` source of `problem_title` | Pair A |
| TST-01 (reference challenge fixtures) | 3.2, 4.2, 6.1–6.3, C04-OQ-12 (algorithmic and HTTP-handler fixtures) | To be confirmed by the mentor |

---

## 8. Open Questions

| # | Question | Affects | Owner | Blocks final approval | Proposed default |
|---|----------|---------|-------|-----------------------|------------------|
| C04-OQ-01 | Maximum length of a title. | 3.2, F02, F08 | Pair A, Pair D | No | 120 characters |
| C04-OQ-02 | Allowed difficulty values, and whether topic tags are in the MVP. Options: A. `easy`, `medium`, `hard`. B. Division style, `div4`, `div3`, `div2`, `div1`, with readable labels for learners. C. A numeric rating (for example 800–3500). | 3.2, F02, F11 | Mentor | No | A (`easy`, `medium`, `hard`), stored as a string with an allowed-values check so that B can be adopted later without a schema change; tags not in the MVP |
| C04-OQ-03 | Statement format (plain text or Markdown), how it is rendered safely, and whether the input/output contract and constraints are separate fields. | 3.2, F03 | Pair A, F03 owner | Yes | Markdown rendered without raw HTML; constraints as part of the statement |
| C04-OQ-04 | Are the public examples on the problem page the same as the sample tests used by "Run samples"? | 3.2, F03, F05 | Pair A, F05 owner | Yes | Yes, one set |
| C04-OQ-05 | Are time and memory limits set per problem, or global? What are the maximum values? | 3.2, F07-02a | Mentor, F07 owner | Yes | Per problem, capped by platform maximums set in F07-02a |
| C04-OQ-06 | Is a reference solution required, and must F09 check before publishing that it passes all tests and that the starter template compiles (T-07)? | 3.2, 4.2, F09-02 | Mentor | No | Yes to both |
| C04-OQ-07 | Are tests versioned when an admin replaces them? Does a submission record which test version judged it? | P-09, CON-05 | Pair A, Pair D | Yes | Yes, a version number on the test set, stored with the result |
| C04-OQ-08 | Can a problem with submissions be deleted, or only unpublished? | P-10, SPEC-F08, CON-05 | Mentor, Pair D | No | Only unpublished (`RESTRICT` on `submissions.problem_id`) |
| C04-OQ-09 | How does the judge Job keep hidden test assets unreadable by the learner's code (V-06)? | V-06, F07-04, DEC-08 | Pair C (F07-04), Pair B (SPK-05 / DEC-08 input) | Yes | Input for DEC-08 and SPK-05; no default here |
| C04-OQ-10 | Can compiler messages contain hidden test code, and if so, how are they filtered before they are shown (V-08)? | V-08, F07-05 | Pair C | Yes | Depends on 6.2 |
| C04-OQ-11 | Maximum size of a test asset or package, and allowed file types in a package (6.1 B and C). Note: SPK-05 SC-B proposes delivering code and tests into a `network=none` Pod via read-only volumes (e.g., ConfigMap, which imposes a ~1 MiB Kubernetes/etcd limit). | 6.1, F09-02, SPK-05 | Pair A, F09 owner | No | Decided after DEC-09 and DEC-08 |
| C04-OQ-12 | Go module (`go.mod`), Go version pin (DEC-14) and allowed package dependencies in a problem package under offline (`no network`) judge execution (SPK-01, SPK-05). | 3.2, 6.2, F07-04, F09-02, TST-01 | Pair A, F07 owner, Mentor | Yes | Standard library only, pinned `go.mod` per DEC-14 (no external module downloads) |

---

## 9. What Blocks Final Approval

The contract cannot be marked Approved until all of the following are true:

1. DEC-09 is decided by the mentor (MA-14), covering 6.1, 6.2, 6.3 and 6.4.
2. The open questions marked "Blocks final approval: Yes" in section 8 are resolved: C04-OQ-03, 04, 05, 07, 09, 10 and 12.
3. CON-05 (preliminary, PR #40) defines the problem and test-asset entities that follow section 3.2 and DEC-09, and closes `C05-OQ-21`.
4. CON-02 is confirmed to match the chosen storage (P-07).
5. Every dependent pair in section 7 has recorded its agreement in section 10.

---

## 10. Agreement Record

| Dependent pair | Agreed? | Recorded on | Date |
|----------------|---------|-------------|------|
| Pair B | No (draft) | — | — |
| Pair C | No (draft) | — | — |
| Pair D | No (draft) | — | — |

*Note: final sign-off rows will be confirmed once Sprint 3 owners are assigned for SPEC-F03, F03-01, F05, SPEC-F09, F09-02 and TST-01 (section 7).*

---

## 11. Change Log

| Date | Change | Author |
|------|--------|--------|
| 2026-11-22 | First draft for Sprint 2 (issue #31) | Pair A |
| 2026-11-24 | Added Go dependency/offline build question (C04-OQ-12), SPK-05 ConfigMap size constraint (C04-OQ-11), missing backlog dependents (SPEC-F09, TST-01), synced merged CON-02/F06/F08 refs vs CON-05 PR #40, and resolved T-07/V-05/F04 cross-ref gaps | Pair A |