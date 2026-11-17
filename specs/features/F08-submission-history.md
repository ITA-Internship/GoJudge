# Feature Spec: F00 – Feature Title

> **How to use this template**
> 1. Copy this file to `specs/features/F00-short-name.md` and replace the placeholders.
> 2. Mark sections **Required**, **Conditional** or **Optional** as shown. Delete any section marked *Conditional* or *Optional* that does not apply, or write "Not applicable" with a reason.
> 3. Do not mark the feature **Approved** while acceptance criteria are ambiguous, a required contract is undefined, or a critical question is open. See [Definition of Ready](../../CONTRIBUTING.md#10-definition-of-ready).
> 4. Delete this blockquote before review.

---

## 1. Feature Metadata *(Required)*

| Field | Value |
|-------|-------|
| Feature ID | F00 |
| Title | *Short, descriptive title* |
| Status | Draft / In Review / Approved |
| Responsible pair | *Pair names* |
| Reviewer pair | *Pair names, must differ from the responsible pair* |
| GitHub Issue | *Link to the issue* |
| Last updated | *YYYY-MM-DD* |

---

## 2. Overview and User Story *(Required)*

**Purpose:** *What problem does this feature solve, in one or two sentences?*

**User story:**

> As a *[user role, e.g. learner or administrator]*,
> I want *[capability]*,
> so that *[expected value]*.

**Expected value:** *How will the project or its users benefit once this is delivered?*

---

## 3. Scope *(Required)*

**In scope:**

- *Item*

**Out of scope:**

- *Item*

---

## 4. Functional Requirements *(Required)*

Each requirement must be numbered and observable, so that a tester can confirm it.

| ID | Requirement |
|----|-------------|
| FR-01 | *The system shall …* |
| FR-02 | *The system shall …* |

---

## 5. Acceptance Criteria *(Required)*

Write each criterion as Given / When / Then. Number them AC-01, AC-02, and so on. Every criterion must be testable.

### AC-01 – *Short name*

- **Given** *[precondition]*
- **When** *[action]*
- **Then** *[observable outcome]*

### AC-02 – *Short name*

- **Given** *[precondition]*
- **When** *[action]*
- **Then** *[observable outcome]*

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

- *F00 – Title, and the nature of the dependency*

**Open questions:**

| # | Question | Owner | Status |
|---|----------|-------|--------|
| 1 | *Question* | *Pair or person* | Open / Resolved |

**Blockers:**

- *Blocker, or "None"*

---

## 12. Review and Approval *(Required)*

| Role | Name or pair | Date | Outcome |
|------|--------------|------|---------|
| Author | *Pair* | *YYYY-MM-DD* | *Submitted* |
| Reviewer | *Pair* | *YYYY-MM-DD* | *Approved / Changes requested* |

**Approval checklist:**

- [ ] Acceptance criteria are unambiguous and numbered.
- [ ] Every required contract is defined and agreed.
- [ ] Every acceptance criterion has a verification entry.
- [ ] Security considerations are documented.
- [ ] No critical open questions remain.
- [ ] Reviewer is from a different pair.
