# Feature Spec: F04 – Browser Code Editor

---

## 1. Feature Metadata *(Required)*

| Field | Value |
|-------|-------|
| Feature ID | F04 |
| Title | Browser code editor |
| Status | In Review |
| Responsible pair | Pair D |
| Reviewer pair | Pair A |
| GitHub Issue | [issue](https://github.com/ITA-Internship/GoJudge/issues/24) |
| Last updated | 2026-11-11 |

---

## 2. Overview and User Story *(Required)*

**Purpose:** Let a learner write a Go solution inside the platform. The editor starts from a starter template supplied with the problem and keeps the learner's current work while their session is active.

**User story:**

> As a learner,
> I want to edit Go source code in the browser, starting from a starter template, and keep my current work in the active session,
> so that I can write my solution without leaving the platform.

**Expected value:** Learners can write, leave and return to a solution without losing work. Run samples (F05) and Submit (F06) get one well-defined source of the learner's code.

---

## 3. Scope *(Required)*

**In scope:**

- Editing Go source text in the browser. Go is the only language.
- Loading a starter template from problem data.
- Preserving the learner's current work during the active session. The mechanism is DEC-13 and is recorded in this spec as Open.
- The edge cases in section 8: blank content, unusable or missing template, reload or navigation away, session expiry or sign-out, two tabs, very large source, and non-Go content.
- What the editor stores, and where session data may be kept (section 9).

**Out of scope:**

- Implementation of the editor (F04-01, F04-02).
- Deciding DEC-13. It stays Open until the mechanism is approved.
- Defining the CON-04 starter template format. It is an input to this spec only.
- Running, compiling, testing or submitting code (F05, F06, F07). The editor only provides the content.
- Saving drafts beyond the active session, or across devices. The MVP entry says "active session".
- Multiple files, other languages, linting, formatting, autocomplete and language-server features. None is part of the F04 entry in the MVP Proposal.
- A "reset to starter template" control. This is not in the MVP entry and is recorded as OQ-08.

---

## 4. Functional Requirements *(Required)*

| ID | Requirement |
|----|-------------|
| FR-01 | The system shall present one editable text area holding Go source for the current problem, with Go syntax highlighting. Go is the only supported language, and the editor shall not offer a language selector. |
| FR-02 | The editor shall support typing, deleting, selecting, cut, copy, paste, undo and redo. |
| FR-03 | The editor shall treat its content only as text. It shall not execute, compile, evaluate, or render the content as HTML, and editing shall not cause any request that runs the code. |
| FR-04 | The editor shall allow blank content. It shall tell dependent features (F05, F06) whether the current content is blank. Whether blank content may be run or submitted is decided in SPEC-F05 and SPEC-F06 (OQ-10). |
| FR-05 | Pasted content shall be inserted as plain text only. Formatting is dropped. Files or images dropped onto the editor shall not be inserted. |
| FR-06 | The editor shall not reject, change or flag text only because it is not valid Go. Compilation problems are reported through the judge path, not by the editor. |
| FR-07 | The editor shall not accept the NUL character (U+0000). A paste containing it shall be rejected as a whole, with a message, and the existing content left unchanged. |
| FR-08 | The editor shall enforce a max source size. Input that would take the content over it, by typing or by paste, shall be rejected as a whole. The existing content shall stay unchanged and the learner shall be told the limit. The editor shall never truncate content silently. |
| FR-09 | For content up to the max source size, the editor shall keep accepting input and shall not freeze. A measurable responsiveness target shall be set together with max source size (OQ-09). |
| FR-10 | When a learner opens a problem and has no preserved work for it in the active session, the editor shall load its initial content from the starter template in the problem data. |
| FR-11 | The editor shall not contain template text of its own. The template comes only from problem data. The data format is an input from CON-04 (OQ-06). |
| FR-12 | If the starter template is missing, empty, whitespace-only, contains the NUL character (U+0000), or is longer than the max source size, the editor shall open with empty content and show a visible notice that no usable starter template is available. The learner shall still be able to edit. |
| FR-13 | If the problem data cannot be retrieved (error or timeout, OQ-13) and there is no preserved work, the editor shall show an error state with a manual retry. It shall not open as if the template were simply absent. The editor shall not retry automatically. If preserved work exists, the editor shall open with it. |
| FR-14 | The starter template shall be treated as untrusted text and shown as plain text only. |
| FR-15 | While the session is active, the editor shall preserve the learner's current content for each problem. The content shall be restored after (a) moving to another page of the platform and returning, (b) a page reload, and (c) leaving the platform and returning within the active session. Whether closing the tab or browser counts as part of (c) is OQ-02. |
| FR-16 | Preserved work shall take precedence over the starter template, including preserved blank content. "No preserved work" and "preserved empty content" are different states. Later changes to a problem's template shall not alter preserved work. |
| FR-17 | Preserved work shall be tied to the signed-in user and the problem. It shall never be shown to a different user or on a different problem. |
| FR-18 | When the learner signs out, their preserved work shall be removed, because the active session has ended. It shall not reappear at the next sign-in (OQ-02). |
| FR-19 | If the session expires while work is open, the content on screen shall not be cleared. The learner shall be told that the session expired and that they must sign in again. The editor shall not claim the work is saved if it is not. Whether the work is recoverable after signing in again is OQ-04. |
| FR-20 | If preservation fails (for example, storage is unavailable or full), editing shall continue. The learner shall be told, at the time of failure, that their work is not being preserved. |
| FR-21 | Preserved data shall contain only the source text and the minimum metadata needed to restore it (user, problem, last-modified time). It shall contain no credentials, tokens, test data or other secrets. |
| FR-22 | When two browser tabs are open on the same problem, behaviour shall be deterministic and documented, and no tab's work shall be discarded without notice to the learner. The exact policy is OQ-03. |
| FR-23 | Typed content shall be preserved within a bounded delay, and the editor shall make a best-effort final preservation when the page is unloaded (OQ-05). |

---

## 5. Acceptance Criteria *(Required)*

### AC-01 – Edit Go source

- **Given** a signed-in learner has opened a problem page and the editor is ready
- **When** the learner types, deletes, selects, pastes, cuts, copies, undoes and redoes text
- **Then** the content changes accordingly, the text is highlighted as Go, and no language selector is offered.

### AC-02 – Starter template on first open

- **Given** the problem data contains a non-empty starter template, and the learner has no preserved work for that problem in the active session
- **When** the learner opens the problem
- **Then** the editor content equals the template text exactly, and the learner can edit it.

### AC-03 – Template missing, empty or unusable

- **Given** the problem data has a missing, empty or whitespace-only template, or one that contains NUL or is longer than max source size, and the learner has no preserved work
- **When** the learner opens the problem
- **Then** the editor opens with empty content, a visible notice says that no usable starter template is available, editing works, and no error page is shown.

### AC-04 – Problem data cannot be retrieved

- **Given** the request for problem data fails or times out (OQ-13), and the learner has no preserved work
- **When** the learner opens the problem
- **Then** an error state with a retry action is shown, no empty editor is presented as if the template were absent, no automatic retry happens, and after a successful manual retry the editor behaves as in AC-02 or AC-03.

### AC-05 – Preserved blank content is not replaced by the template

- **Given** the learner has cleared the editor, and the problem has a non-empty starter template
- **When** the learner leaves the problem and returns within the active session
- **Then** the editor is still blank and the template is not loaded.

### AC-06 – Navigate away and back

- **Given** the learner has changed the editor content and has not submitted it
- **When** the learner moves to another page of the platform and returns to the same problem within the active session
- **Then** the editor shows exactly the content the learner left.

### AC-07 – Page reload

- **Given** the learner has changed the editor content and has not submitted it
- **When** the learner reloads the page within the active session
- **Then** the editor shows all changes made more than X ms before the reload (OQ-05). The mechanism that provides this is DEC-13.

### AC-08 – Preserved work beats a changed template

- **Given** the learner has preserved work for a problem, and the problem's starter template has since changed
- **When** the learner opens the problem in the same active session
- **Then** the editor shows the preserved work and not the new template.

### AC-09 – Sign-out removes preserved work

- **Given** the learner has preserved work for a problem
- **When** the learner signs out and signs in again as the same user
- **Then** the editor shows the starter template, not the earlier work.

### AC-10 – Another user does not see the work

- **Given** a learner has preserved work on a browser, and a different user then signs in on the same browser
- **When** the second user opens the same problem
- **Then** the second user sees the starter template, never the first user's work.

### AC-11 – Session expires while work is open

- **Given** the learner has unsaved content in the editor
- **When** the session expires
- **Then** the content stays visible and is not cleared, the learner is told the session expired and that they must sign in again, and the editor does not claim the work is saved unless it is. Recovery after signing in again is OQ-04.

### AC-12 – Two tabs on the same problem

- **Given** the same problem is open in two tabs of the same browser, and the learner edits in one tab and then in the other
- **When** the learner continues, reloads either tab, or returns to the problem
- **Then** no tab's visible content is changed by the other tab without notice, and the content restored after a reload follows the documented policy.

### AC-13 – Very large source

- **Given** the editor holds content close to max source size bytes (OQ-09)
- **When** the learner pastes text that would take the content over max source size, and separately types past it
- **Then** each over-limit paste or keystroke is rejected as a whole, the existing content is unchanged, a message states the limit, and the editor stays responsive within the target set in OQ-09.

### AC-14 – Content that is not Go

- **Given** the editor is ready
- **When** the learner pastes plain text that is not Go, pastes HTML-formatted content, drops a file onto the editor, or pastes text containing NUL
- **Then** the plain text is inserted without rejection or warning, formatting is dropped, the dropped file is not inserted, the text containing NUL is rejected as a whole, and nothing pasted is rendered as HTML or run.

### AC-15 – Preservation fails

- **Given** the preservation mechanism is unavailable or full
- **When** the learner edits
- **Then** editing continues, and the learner is told that their work is not being preserved.

### AC-16 – No execution from the editor

- **Given** the learner is editing
- **When** the learner types or pastes, with no run or submit action taken
- **Then** no request that runs, compiles or evaluates the code is sent, and the content is not interpreted as HTML or script.

### AC-17 – Blank state is available to Run and Submit

- **Given** the editor is ready
- **When** the learner clears all content, and then types any content again
- **Then** the editor reports the content as blank to F05 and F06 in the first case and as not blank in the second, using the definition of blank from OQ-11, and the editor itself does not block or allow run or submit.

### AC-18 – Template comes only from problem data

- **Given** two problems have different non-empty starter templates, and a third problem has none
- **When** the learner opens each of them for the first time in the active session
- **Then** the first two editors contain exactly their own template and the third opens empty, with no built-in or default text in the editor content of any of them.

### AC-19 – Template is shown as plain text

- **Given** the problem's starter template contains HTML tags and script text
- **When** the learner opens the problem
- **Then** the template appears literally as text, nothing is rendered as HTML, and no script runs.

### AC-20 – Problem data fails but preserved work exists

- **Given** the learner has preserved work for a problem, and the request for problem data fails or times out (OQ-13)
- **When** the learner opens the problem
- **Then** the editor opens with the preserved work, no error state replaces it, and no automatic retry happens.

### AC-21 – Leave the platform and return

- **Given** the learner has changed the editor content and has not submitted it
- **When** the learner opens another site, returns to the platform within the active session, and opens the same problem
- **Then** the editor shows exactly the content the learner left. Whether closing the tab or browser also keeps the work is OQ-02.

### AC-22 – Preservation delay and unload

- **Given** the learner types in the editor
- **When** the page is closed or navigated away less than X ms after the last keystroke (OQ-05)
- **Then** the editor makes a best-effort final preservation, and everything typed more than X ms before the unload is restored on return.

### AC-23 – Preserved work is per problem

- **Given** a learner has preserved work on problem A and none on problem B
- **When** the learner opens problem B
- **Then** problem B shows its own starter template, never the work from problem A.

### AC-24 – Preserved data contains no secrets

- **Given** the learner has preserved work
- **When** the stored preserved data is inspected
- **Then** it contains only the source text, the user, the problem and the last-modified time, and no credentials, tokens, test data or other secrets.

---

## 6. Technical Design *(Required)*

**Component responsibilities:**

| Component | Responsibility |
|-----------|----------------|
| *e.g. Go API* | *What it does for this feature* |

**Dependencies:**

- *Internal or external component, library or service, with reason*

**Architecture decisions:**

- *Link to any decision record, or state "None". Do not introduce a new architectural decision without raising it as an issue first.*

---

## 7. Contracts *(Conditional)*

Complete this section if the feature exposes or consumes an API, stores data, publishes or consumes messages, or changes state. Contracts must be agreed with dependent pairs before dependent development begins.

### 7.1 API Endpoints *(if applicable)*

| Method | Path | Purpose | Auth required |
|--------|------|---------|---------------|
| *GET* | */path* | *Purpose* | *Yes / No* |

**Request schema:** *Link to, or describe, the request body and parameters.*

**Response schema:** *Link to, or describe, success and error responses.*

### 7.2 Data Models *(if applicable)*

*Describe each entity, its fields, types and constraints.*

### 7.3 Messaging Contracts *(if applicable)*

*Describe the queue or topic name, message fields, and expected consumer behaviour.*

### 7.4 State Transitions *(if applicable)*

*Describe the allowed states and transitions, for example a submission moving from queued to running to completed.*

---

## 8. Edge Cases and Failure Handling *(Conditional)*

Complete this section where relevant. Consider the following.

- [ ] Invalid or missing input
- [ ] Timeouts
- [ ] Retries and their limits
- [ ] Duplicate processing
- [ ] Unavailable dependencies

| Scenario | Expected behaviour |
|----------|--------------------|
| *Scenario* | *Behaviour* |

---

## 9. Security Considerations *(Required)*

Confirm each point, or state "Not applicable" with a reason.

- [ ] **Authentication:** *How is the caller identified?*
- [ ] **Authorisation:** *Which roles may perform each action?*
- [ ] **Input validation:** *What input is validated, and where?*
- [ ] **Sensitive data:** *What sensitive data is handled, and how is it protected?*
- [ ] **Execution isolation:** *Required if the feature runs learner-submitted code. Describe the isolation approach and how it is validated.*

---

## 10. Verification Plan *(Required)*

Map every acceptance criterion to at least one test or verification scenario.

| AC | Verification type | Test or scenario | Automated? |
|----|-------------------|------------------|------------|
| AC-01 | *Unit / integration / end-to-end / manual* | *Description* | *Yes / No* |
| AC-02 | *Unit / integration / end-to-end / manual* | *Description* | *Yes / No* |

---

## 11. Dependencies and Open Questions *(Required)*

**Related features:**

- **F01 – Accounts and access:** defines sign-in, sign-out, expiry and user identity, which the editor relies on (FR-17 to FR-19).
- **F03 – Problem details:** supplies the problem data that carries the starter template, and hosts the editor.
- **F05 – Run sample tests:** consumes the editor content and the blank state (FR-04).
- **F06 – Submit solution:** consumes the editor content and the blank state; its request size must be aligned with max source size (OQ-09).
- **F07 – Automated judge:** executes the code. It owns the isolation requirements, and its limits must be aligned with max source size.
- **Contracts and decisions:** CON-04, CON-03, DEC-13, DEC-09.
- **Gated items:** F04-01 (editor integration with Go mode and starter template) and F04-02 (active-session preservation)

**Open questions:**

| # | Question | Owner | Status |
|---|----------|-------|--------|
| OQ-01 | **DEC-13.** Which mechanism preserves current work in the active session? | Mentor | Open |
| OQ-02 | What does "active session" mean: the signed-in session, or the browser tab or browser session? Does closing the tab or browser end it? This spec proposes the signed-in session. | Pair D | Open |
| OQ-03 | What should happen when two tabs are open on the same problem: independent copies, a warning, or synchronisation? | Pair D | Open |
| OQ-04 | After session expiry, should unsaved work be recoverable once the learner signs in again, or discarded? | Mentor | Open |
| OQ-05 | How soon after a keystroke must content be preserved, and is a best-effort save on page unload enough? | Pair D | Open |
| OQ-06 | What is the starter template format in problem data? (CON-04) | Pair C | Open |
| OQ-07 | When no usable template exists, should the editor open empty, or with a generic skeleton? Spec assumes empty. | Mentor | Open |
| OQ-08 | Should the learner have a "reset to starter template" control? | Mentor | Open |
| OQ-09 | What is the max source size, and what responsiveness target applies? | Pair D | Open |
| OQ-10 | May blank content be run or submitted, or should it be blocked? | SPEC-F05 and SPEC-F06 authors | Open |
| OQ-11 | What counts as blank content: zero characters only, or also whitespace-only content? | Pair D | Open |
| OQ-12 | How long may preserved work stay in the browser after the session expires without sign-out, or after the tab or browser closes? FR-18 only removes it on sign-out. | Mentor | Open |
| OQ-13 | What is the timeout for retrieving problem data? | Pair D | Open |
| OQ-14 | Are line endings (CRLF or LF) and tabs kept as typed or normalised? Is the source size counted before or after normalisation? | Pair D | Open |
| OQ-15 | If the learner signs out in one tab, what happens to the editor content still open in other tabs? | Pair D | Open |

**Blockers:**

- **DEC-13 is Open.** The spec records it and does not claim a decision. FR-15, FR-16, FR-23, AC-06, AC-07, AC-08, AC-09, AC-10, AC-15, AC-21, AC-22, AC-23 and AC-24 depend on it. F04-02 cannot start until DEC-13 is approved and this spec is Approved.
- **The CON-04 starter template format is an open input** (OQ-06). CON-04 is in draft. F04-01 cannot start until it is agreed.

---

## 12. Review and Approval *(Required)*

| Role | Name or pair | Date | Outcome |
|------|--------------|------|---------|
| Author | Pair D | 2026-11-11 | Submitted |
| Reviewer | Pair A | — | Pending |

Approval is not claimed. The status is **In Review**.

**Approval checklist:**

- [ ] Acceptance criteria are unambiguous and numbered.
- [ ] Every required contract is defined and agreed.
- [ ] Every acceptance criterion has a verification entry.
- [ ] Security considerations are documented.
- [ ] No critical open questions remain.
- [ ] Reviewer is from a different pair.