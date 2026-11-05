# Feature Spec: F01 – Accounts and access

> **How to use this template**
> 1. Copy this file to `specs/features/F00-short-name.md` and replace the placeholders.
> 2. Mark sections **Required**, **Conditional** or **Optional** as shown. Delete any section marked *Conditional* or *Optional* that does not apply, or write "Not applicable" with a reason.
> 3. Do not mark the feature **Approved** while acceptance criteria are ambiguous, a required contract is undefined, or a critical question is open. See [Definition of Ready](../../CONTRIBUTING.md#10-definition-of-ready).
> 4. Delete this blockquote before review.

---

## 1. Feature Metadata *(Required)*

| Field | Value |
|-------|-------|
| Feature ID | F01 |
| Title | Accounts and access |
| Status | Draft |
| Responsible pair | Pair D |
| Reviewer pair | Pair A, Pair C |
| GitHub Issue | [issue](https://github.com/ITA-Internship/GoJudge/issues/10) |
| Last updated | 2026-11-05 |

---

## 2. Overview and User Story *(Required)*

**Purpose:** It is important to know who is calling the API and what that person is allowed to do. F01 lets people register, sign in and out. Also it separates learners from administrators, so that every other feature can rely on a verified identity and an enforced role.

**User story:**

> As a **learner**,
> I want to register, sign in and sign out,
> so that my submissions and progress are tied to me and are not visible to others.

> As an **administrator**,
> I want to sign in with an account that carries the administrator role,
> so that I can reach problem and test administration that learners cannot.

**Expected value:** The foundation for ownership of submissions (F06, F08), progress (F10) and role-protected administration (F09). Delivering it early lets dependent pairs build against a stable identity and role contract.

---

## 3. Scope *(Required)*

**In scope:**

- Registration of a new learner accounts.
- Sign in and sign out from the web application.
- A local user record in PostgreSQL, created on first authenticated request, that gives other features a stable user ID.
- Two roles, `learner` and `administrator`, with `learner` as the default.
- Validation of the access token on every protected API request.
- Role enforcement, denying by default, as reusable API middleware.
- An endpoint returning the caller's own profile and role.
- A documented way to designate an administrator without any public self-service path.

**Out of scope:**

- Storing or handling passwords. Credentials are held only by the identity provider.
- Custom-built password reset, email verification or multi-factor flows.
- Social login or enterprise single sign-on.
- Profile editing, account deletion and data export.
- A user-management UI for administrators, including promoting or demoting users through the application.
- Roles beyond learner and administrator, and per-problem permissions.
- API rate limiting and abuse protection beyond what the identity provider offers.
- Server-side revocation of already-issued access tokens.
- Any behaviour of submissions, problems or the judge.

---

## 4. Functional Requirements *(Required)*

| ID | Requirement |
|----|-------------|
| FR-01 | The system shall allow a visitor to register a new account with an email address and password through the identity provider's registration flow. |
| FR-02 | The system shall reject registration with an email address that is already registered. |
| FR-03 | The system shall allow a registered user to sign in with valid credentials and shall give the web application an access token for calling the API. |
| FR-04 | The system shall not issue an access token when sign-in credentials are invalid. |
| FR-05 | The system shall allow a signed-in user to sign out, which ends the session in the web application and at the identity provider. |
| FR-06 | The system shall support exactly two roles: `learner` and `administrator`. |
| FR-07 | The system shall assign the `learner` role to every newly registered user. |
| FR-08 | The system shall create a local user record, keyed by the identity provider's subject identifier, the first time an authenticated request is received for a subject, and shall not create more than one record per subject. |
| FR-09 | The system shall validate the access token on every protected API request |
| FR-10 | The system shall respond `401 Unauthorized` to a protected request with a missing or invalid token. |
| FR-11 | The system shall respond `403 Forbidden` when an authenticated user calls an action that requires a role they do not hold. |
| FR-12 | The system shall deny by default: an endpoint is accessible only to callers it explicitly allows. |
| FR-13 | The system shall provide `GET /api/v1/me` returning the caller's own profile. |
| FR-14 | The system shall determine a caller's role only from server-controlled data and shall ignore any role value supplied in a request body, query string or client-controlled token claim. |
| FR-15 | The system shall provide a documented way to designate at least one administrator that is not available through any public endpoint. |
| FR-16 | The system shall not store, log or return passwords or access tokens. |
| FR-17 | The system shall log authentication and authorisation failures with a reason code and, where known, the user ID, using structured logging. |
| FR-18 | The system shall expose role checks as reusable middleware so that F02–F10 apply them without re-implementing authentication. |


---

## 5. Acceptance Criteria *(Required)*

### AC-01 – Register a new learner

- **Given** a visitor with an email address that is not registered
- **When** they complete the registration flow with a valid email and password
- **Then** an account is created at the identity provider and the visitor can sign in with those credentials

### AC-02 – Duplicate registration rejected

- **Given** an email address that is already registered
- **When** a visitor attempts to register with the same email address
- **Then** registration is rejected, no second account exists, and the visitor is told the address cannot be used

### AC-03 – Sign in with valid credentials

- **Given** a registered user
- **When** they sign in with valid credentials
- **Then** the web application holds an access token, and `GET /api/v1/me` with that token returns `200` with the user's own profile

### AC-04 – First request provisions a learner record

- **Given** a valid access token for a subject that has no local user record
- **When** the API receives its first protected request with that token
- **Then** exactly one user record is created with role `learner`, and `GET /api/v1/me` returns that user's ID and `"role": "learner"`

### AC-05 – Provisioning is idempotent

- **Given** a valid access token for a subject with no local user record
- **When** several protected requests with that token arrive concurrently, or the user signs in again later
- **Then** exactly one user record exists for that subject and every response returns the same user ID

### AC-06 – Invalid credentials

- **Given** a registered user
- **When** they attempt to sign in with a wrong password
- **Then** no access token is issued, the user remains unauthenticated and the protected pages stay inaccessible

### AC-07 – Sign out

- **Given** a signed-in user
- **When** they sign out
- **Then** the web application discards its session state, the identity provider session ends, and opening a protected page requires signing in again

### AC-08 – Missing token

- **Given** no `Authorization` header
- **When** a protected endpoint is called
- **Then** the response is `401` and contains no user data

### AC-09 – Invalid token

- **Given** a token that is expired, has a bad signature, has the wrong issuer or audience, or is malformed
- **When** a protected endpoint is called with it
- **Then** the response is `401`, no local user record is created, and the failure is logged with a reason code but without the token

### AC-10 – Learner denied administrator action

- **Given** a signed-in user with role `learner`
- **When** they call an administrator-only endpoint
- **Then** the response is `403` and no administrator action takes place

### AC-11 – Administrator allowed

- **Given** a signed-in user with role `administrator`
- **When** they call an administrator-only endpoint
- **Then** the request passes the role check

### AC-12 – Role cannot be self-assigned

- **Given** a signed-in `learner`
- **When** they send a request containing a role value in the body or query string, or present a token carrying a role claim they could have set themselves
- **Then** their role remains `learner`, and `GET /api/v1/me` still returns `"role": "learner"`

### AC-13 – Profile is own data only

- **Given** two signed-in users
- **When** each calls `GET /api/v1/me`
- **Then** each receives only their own profile, with no token or credential material

### AC-14 – Sensitive data kept out of logs

- **Given** any sign-in, failed authentication or denied authorisation
- **When** the structured log output is inspected
- **Then** no access token, password or Authorization header value appears in it

### AC-15 – Identity provider unavailable

- **Given** the identity provider's key endpoint cannot be reached
- **When** a protected request arrives with a valid token and the signing keys are cached
- **Then** the request is processed normally; and when no signing key is available, the response is `503`

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

- [x] Invalid or missing input
- [x] Timeouts
- [x] Retries and their limits
- [x] Duplicate processing
- [x] Unavailable dependencies

| Scenario | Expected behaviour |
|----------|--------------------|
| No `Authorization` header, or a scheme other than Bearer | `401`, no user data (AC-08). |
| Token malformed, expired, wrong issuer or audience, bad signature | `401`, no user record created, reason code logged (AC-09). |
| Token with a signing key identifier not yet known | Refresh the key set once, within a bounded time, then validate. If still unknown, `401`. |
| Key endpoint slow or unreachable | Use cached keys if present (AC-15). Key fetches have a timeout and a small bounded number of retries, with no unbounded retry loop. With no usable key, `503`. |
| Two first requests for the same new subject arrive concurrently | One record is created, relying on the unique `auth_subject` constraint. The losing request reads the existing record (AC-05). |
| Token has no email claim | User is created with `email` null. Sign-in is not blocked. |
| Email or display name changes at the identity provider | Identity is unaffected because it is keyed by `auth_subject`. Profile fields may be refreshed on a later request. |
| Display name or email claim exceeds the stored length | Value is truncated to the limit or replaced by the fallback. The request is not rejected on this ground. |
| Client sends a `role` field or a forged role claim | Ignored (AC-12). |
| Administrator role is stored for a user whose token is later invalid | The token is rejected first. The role has no effect without a valid token. |
| User is blocked or deleted at the identity provider | No new tokens are issued. Already-issued tokens work until they expire. |
| Registration with an already registered email | Rejected at the identity provider (AC-02). |
| Identity provider unavailable during sign-in or registration | The application shows an error and the user stays signed out. Valid existing tokens continue to work. |
| Database unavailable during provisioning | `503`. No partial user state is returned. |

---

## 9. Security Considerations *(Required)*

- [x] **Authentication:** *How is the caller identified?*
- [x] **Authorisation:** Roles are `learner` and `administrator`. The role is read from server-controlled data and never from client input (FR-14, AC-12). Administrator-only actions, including everything that exposes hidden tests (F09), return `403` to learners (AC-10). The API is the single authority. Denied attempts are logged with the user ID.
- [x] **Input validation:** The Authorization header is parsed strictly. Token claims used for the local profile are length-bounded and treated as untrusted text. No role is accepted from any request field. Registration input (email format, password strength) is validated by the identity provider. `GET /api/v1/me` takes no input.
- [x] **Sensitive data:** Application never receives or stores passwords. Access tokens are not stored server-side and are never logged or echoed in responses (AC-14). The only personal data stored is email and display name. The email is returned only to its owner and kept out of logs.
- [] **Execution isolation:** Not applicable. F01 does not run learner-submitted code.

---

## 10. Verification Plan *(Required)*

| AC | Verification type | Test or scenario | Automated? |
|----|-------------------|------------------|------------|
| AC-01 | End-to-end | Register a new user through the hosted registration flow, then sign in. | No |
| AC-02 | End-to-end | Attempt to register an existing email and confirm rejection and a single account. | No |
| AC-03 | Integration | Valid test token on `GET /api/v1/me` returns `200` with the expected profile. | Yes |
| AC-04 | Integration | Valid token for an unknown subject creates one `learner` row. Check the DB and the `/me` response. | Yes |
| AC-05 | Integration | N concurrent first requests and a repeated sign-in leave exactly one row and one ID. | Yes |
| AC-06 | End-to-end | Wrong password at the hosted login: no token, protected pages stay inaccessible. | No |
| AC-07 | End-to-end | Sign out, confirm provider session ended and protected page demands sign-in. | No |
| AC-08 | Integration | Request with no header and with a non-Bearer scheme returns `401`. | Yes |
| AC-09 | Unit and integration | Table test of expired, wrong issuer, wrong audience, bad signature, malformed tokens: all `401`, no row created, no token in the log. | Yes |
| AC-10 | Integration | Learner token against a role-protected test route returns `403`. | Yes |
| AC-11 | Integration | Administrator token against the same route passes the check. | Yes |
| AC-12 | Integration | Learner sends a role in body and query, and a token with a forged role claim: stored role and `/me` stay `learner`. | Yes |
| AC-13 | Integration | Two users each call `/me` and receive only their own fields. | Yes |
| AC-14 | Integration | Run AC-03, AC-09 and AC-10 with log capture and assert no token, password or header value appears. | Yes |
| AC-15 | Integration | Fake provider stopped: cached keys still validate. With an empty cache, `503`. | Yes |

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
| Author | Pair D | Date | Submitted |
| Reviewer | Pair A | Date | Pending |
| Reviewer | Pair C | Date | Pending |

**Approval checklist:**

- [ ] Acceptance criteria are unambiguous and numbered.
- [ ] Every required contract is defined and agreed.
- [ ] Every acceptance criterion has a verification entry.
- [ ] Security considerations are documented.
- [ ] No critical open questions remain.
- [ ] Reviewer is from a different pair.
