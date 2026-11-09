# Feature Spec: F01 – Accounts and access

## 1. Feature Metadata *(Required)*

| Field | Value |
|-------|-------|
| Feature ID | F01 |
| Title | Accounts and access |
| Status | In Review |
| Responsible pair | Pair D |
| Reviewer pair | Pair A, Pair C |
| GitHub Issue | [issue](https://github.com/ITA-Internship/GoJudge/issues/10) |
| Last updated | 2026-11-09 |

---

## 2. Overview and User Story *(Required)*

**Purpose:** It is important to know who is calling the API and what that person is allowed to do. F01 lets people register, sign in and out. It also separates learners from administrators, so that every other feature can rely on a verified identity and an enforced role.

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

- Registration of new learner accounts.
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
- Any behavior of submissions, problems or the judge, including the submission lifecycle and verdicts.

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
| FR-07 | The system shall assign the learner role to every newly registered user, subject to the OQ-1 decision. The default is learner unless the administrator designation procedure (FR-15) has applied. |
| FR-08 | The system shall create a local user record, keyed by the identity provider's subject identifier, the first time an authenticated request is received for a subject, and shall not create more than one record per subject. |
| FR-09 | The system shall validate the access token on every protected API request: RS256 signature, issuer, audience, expiry, and presence of a subject. |
| FR-10 | The system shall respond `401 Unauthorized` to a protected request with a missing or invalid token. |
| FR-11 | The system shall respond `403 Forbidden` when an authenticated user calls an action that requires a role they do not hold. |
| FR-12 | The system shall deny by default: an endpoint is accessible only to callers it explicitly allows. |
| FR-13 | The system shall provide `GET /api/v1/me` returning the caller's own profile. |
| FR-14 | The system shall determine a caller's role only from server-controlled data and shall ignore any role value supplied in a request body, query string or any token claim not designated as trusted by the OQ-1 decision. |
| FR-15 | The system shall provide a documented way to designate at least one administrator that is not available through any public endpoint. |
| FR-16 | The system shall not store, log or return passwords or access tokens. |
| FR-17 | The system shall log authentication and authorization failures with a reason code and, where known, the user ID, using structured logging. |
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

- **Given** no `Authorization` header or an `Authorization` header with a scheme other than Bearer
- **When** a protected endpoint is called
- **Then** the response is `401` and contains no user data

### AC-09 – Invalid token

- **Given** an access token that is expired, has an invalid signature, has the wrong issuer or audience, has no subject, or is malformed
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
- **When** they send a request containing a role value in the body or query string, or present a token whose role claim is not the one designated as trusted by the OQ-1 decision
- **Then** their role remains `learner`, and `GET /api/v1/me` still returns `"role": "learner"`

### AC-13 – Profile is own data only

- **Given** two signed-in users
- **When** each calls `GET /api/v1/me`
- **Then** each receives only their own profile, with no token or credential material

### AC-14 – Sensitive data kept out of logs

- **Given** any sign-in, failed authentication or denied authorization
- **When** the structured log output is inspected
- **Then** no access token, password or Authorization header value appears in it

### AC-15 – Identity provider unavailable

- **Given** the identity provider's key endpoint cannot be reached
- **When** a protected request arrives with a valid token and the signing keys are cached
- **Then** the request is processed normally; and when no signing key is available, the response is `503`

### AC-16 – Database unavailable during provisioning

- **Given** a valid access token for a subject with no local user record
- **When** the database cannot be reached during provisioning
- **Then** the response is `503`, no partial user data is returned, and no user record exists afterwards

### AC-17 – Deny by default

- **Given** the registered API routes
- **When** a route is registered without an explicit access policy (public, authenticated, or a named role)
- **Then** the route is not reachable, and the automated route check fails

### AC-18 – Administrator designation

- **Given** a subject with a valid token and no administrator role
- **When** the documented designation procedure is carried out
- **Then** `GET /api/v1/me` returns `"role": "administrator"`, and no public endpoint can perform this change

### AC-19 – Denied authorization is logged

- **Given** a signed-in user with role `learner`
- **When** they call an administrator-only endpoint
- **Then** the response is `403`, and a structured log entry contains a reason code and the user ID but no token, password or Authorization header value

---

## 6. Technical Design *(Required)*

**Component responsibilities:**

| Component | Responsibility |
|-----------|----------------|
| Identity provider | Hosts registration, sign-in and sign-out. Stores credentials. Rejects duplicate emails. Issues RS256-signed access tokens. Publishes signing keys (JWKS). |
| Web application | Starts the hosted login, registration and logout flows. Holds the access token in memory for the session. Sends it as `Authorization: Bearer`. Never sends or sets a role. |
| Go API – authentication middleware | Extracts the Bearer token and validates signature (RS256 only), issuer, audience and expiry with the official `github.com/auth0/go-jwt-middleware/v3` and its `validator` and `jwks` sub-packages. Signing keys are cached by `jwks.NewCachingProvider`. Maps failures to `401` or `503` and to log reason codes. |
| Go API – provisioning step | Creates the local user record for a new subject with an insert that relies on the unique `auth_subject` constraint, then reads the record. Refreshes bounded profile fields. |
| Go API – role resolution and role middleware | Resolves the caller's role into the request context and exposes a single accessor to handlers. Provides reusable "require role" middleware. Denies by default and rejects routes without an explicit policy. |
| Go API – `GET /api/v1/me` | Returns the caller's own profile. |
| PostgreSQL | Stores the `users` table. |
| Structured logger | Records authentication and authorization failures with reason codes and user ID, never token, password or header values. |

**Dependencies:**

- Auth0 tenant (identity provider). A development tenant was verified in SPK-02. Tenant access was confirmed with the mentor (MA-05). No secrets or tokens are stored in the repository.
- `github.com/auth0/go-jwt-middleware/v3`. SPK-02 recommends it for production because it greatly shortens parsing code and provides JWKS caching, `RS256` algorithm pinning and standard clock-skew handling.
- PostgreSQL, for the local `users` table.
- DEC-03 (Authentication with Auth0): provisional. Formal approval and ADR recording remain with the mentor (OQ-2).
- SPK-02 (Auth0 OIDC token validation and role claim mapping): Go outcome. Its code is on branch `spike/SPK-02-auth0-oidc`. F01-01 may reuse its test approach.

**Architecture decisions:**

- DEC-03 (Authentication with Auth0): provisional, confirmed viable by SPK-02, not yet approved.

---

## 7. Contracts *(Conditional)*

Complete this section if the feature exposes or consumes an API, stores data, publishes or consumes messages, or changes state. Contracts must be agreed with dependent pairs before dependent development begins.

### 7.1 API Endpoints *(if applicable)*

| Method | Path | Purpose | Auth required |
|--------|------|---------|---------------|
| GET | `/api/v1/me` | Return the caller's own profile and role | Yes |

Registration, sign-in and sign-out are performed at the identity provider and are not API endpoints. F01 adds no public API endpoint. Public routes, if any, are explicitly allow-listed (FR-12).

**Request schema:** All protected requests carry `Authorization: Bearer <access token>`. `GET /api/v1/me` takes no body, no query parameters and no role input.

**Response schema:**
Success, `200`:

```json
{
  "id": "6f1c0a52-3c8e-4d7b-9a51-1f2d3e4a5b6c",
  "email": "learner@example.com",
  "display_name": "Learner",
  "role": "learner"
}
```

- `id` is the local user ID (UUID). `email` and `display_name` may be `null`. The identity provider's `auth_subject` is internal and is not returned. Tokens and credential material are never returned.

Errors (same envelope for all F01 failures; the generic message never contains parser or claim details):

| Status | `code` | When |
|--------|--------|------|
| `401` | `unauthorized` | Missing or invalid token (AC-08, AC-09) |
| `403` | `forbidden` | Authenticated, but the required role is not held (AC-10) |
| `503` | `service_unavailable` | No signing key available, or database unavailable during provisioning (AC-15, AC-16) |

```json
{ "error": { "code": "unauthorized", "message": "Authentication required." } }
```

### 7.2 Data Models

All models are provisional and to be confirmed.

Table `users`:

| Field | Type | Constraints |
|-------|------|-------------|
| `id` | UUID | Primary key, generated by the API or database. Stable ID used by F06, F08, F09 and F10. |
| `auth_subject` | text | Not null, unique. Identity provider subject (`sub`), at most 255 characters. Never changes. |
| `email` | text | Nullable, at most 320 characters. Returned only to its owner. Kept out of logs. |
| `display_name` | text | Nullable, at most 100 characters. Truncated to the limit. |
| `role` | text | Not null, default `learner`, check constraint: `learner` or `administrator`. |
| `created_at` | timestamptz | Not null, default now. |
| `updated_at` | timestamptz | Not null, default now. |


### 7.3 Messaging Contracts

Not applicable. F01 neither publishes nor consumes queue or topic messages.

### 7.4 State Transitions

Not applicable. A user record has no lifecycle states.

---

## 8. Edge Cases and Failure Handling *(Conditional)*

- [x] Invalid or missing input
- [x] Timeouts
- [x] Retries and their limits
- [x] Duplicate processing
- [x] Unavailable dependencies

| Scenario | Expected behavior |
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

- [x] **Authentication:** The caller is identified by an RS256-signed OIDC access token issued by the identity provider (Auth0, DEC-03 provisional). The API validates signature, issuer, audience and expiry on every protected request and requires a subject (FR-09). The algorithm is pinned to RS256. Issuer and audience come from configuration. Signing keys come from the provider's JWKS endpoint and are cached with a bounded TTL. This was verified by SPK-02 with the official middleware. The production sign-in flow for the web application is hosted login with Authorization Code and PKCE (provisional, OQ-6). The validation error text is logged as a reason code only and is never returned.
- [x] **Authorization:** Roles are `learner` and `administrator`. The role is read from server-controlled data and never from client input (FR-14, AC-12). Administrator-only actions, including everything that exposes hidden tests (F09), return `403` to learners (AC-10). The API is the single authority. Denied attempts are logged with the user ID.
- [x] **Input validation:** The Authorization header is parsed strictly. Token claims used for the local profile are length-bounded and treated as untrusted text. No role is accepted from any request field. Registration input (email format, password strength) is validated by the identity provider. `GET /api/v1/me` takes no input.
- [x] **Sensitive data:** The application never receives or stores passwords. Access tokens are not stored server-side and are never logged or echoed in responses (AC-14). The only personal data stored is email and display name. The email is returned only to its owner and kept out of logs.
- [ ] **Execution isolation:** Not applicable. F01 does not run learner-submitted code.

**Open questions from SPK-02:**

| Gap identified in SPK-02 | Security impact | Recorded as |
|----------------|-----------------|-------------|
| The spike read the role from the claim `https://gojudge-api/roles` in the token. It did not decide whether that claim or the local database is the source of truth. | Decides who can grant `administrator`. | OQ-1 |
| DEC-03 is confirmed viable but not formally approved. | The identity provider choice could change. | OQ-2 |
| Stale-cache behavior after the 5-minute TTL, and refresh on an unknown key ID, were not exercised. | AC-15 and the unknown-key edge case depend on it. | OQ-3 |
| The administrator test user's claim contained a misspelling ("Adminisrator") and the learner's was capitalized ("Learner"). The spike mapper lower-cases and trims, and the evidence notes tenant typo handling. | Approximate matching must never grant `administrator`. | OQ-4 |
| Which claims (`email`, name) appear in the access token was inspected but not recorded. | Affects profile fields only. | OQ-5 |
| Sign-in in the spike used the Resource Owner Password flow. The hosted Authorization Code flow with PKCE was not exercised. | The product flow is unverified. | OQ-6 |

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
| AC-09 | Integration | Table test of expired, wrong issuer, wrong audience, bad signature, malformed tokens: all `401`, no row created, no token in the log. | Yes |
| AC-10 | Integration | Learner token against a role-protected test route returns `403`. | Yes |
| AC-11 | Integration | Administrator token against the same route passes the check. | Yes |
| AC-12 | Integration | Learner sends a role in body and query, and a token with a forged role claim: stored role and `/me` stay `learner`. | Yes |
| AC-13 | Integration | Two users each call `/me` and receive only their own fields. | Yes |
| AC-14 | Integration | Run AC-03, AC-09 and AC-10 with log capture and assert no token, password or header value appears. | Yes |
| AC-15 | Integration | Fake provider stopped: cached keys still validate. With an empty cache, `503`. | Yes |
| AC-16 | Integration | Fake DB down on first request: `503`, no row, no partial body. | Yes |
| AC-17 | Unit | Register a route without a policy: it is unreachable and the route-check test fails. | Yes |
| AC-18 | Integration | Apply the documented designation procedure to a test subject, then `/me` returns `administrator`. Confirm no public route can do it. | Yes |
| AC-19 | Integration | Learner calls an admin route, which returns `403`. Assert the log has a reason code and user ID and no token. | Yes |

---

## 11. Dependencies and Open Questions *(Required)*

**Related features:**

- **F02–F10 – Consumers of identity and roles.** They apply the role middleware and do not re-implement authentication.
- **F06, F08 – Submissions.** Submission records reference `users.id` as owner. Visibility of other users' submissions is decided by those features, using the identity F01 provides.
- **F09 – Problem and test administration.** Requires the `administrator` role (FR-11, AC-10). Hidden test definitions are administrator-only. Detailed test feedback to learners follows each problem's visibility rules.
- **F10 – Progress.** Keyed by `users.id`.
- **CON-01 (S1-08) – Submission lifecycle and verdict model.** F01 defines no submission states or verdicts and supplies only `users.id` as owner.
- **DEC-03 – Authentication with Auth0.** Provisional.

**Open questions:**

| # | Question | Owner | Status | Critical |
|---|----------|-------|--------|----------|
| OQ-1 | What is the source of truth for the role: the local `users.role` column or the validated Auth0 claim, which SPK-02 showed can be mapped? Or a defined combination of both? | Pair C, Pair D | Open | Yes |
| OQ-2 | Formal approval and ADR for DEC-03 (Auth0). The spec is written against it as provisional. | Mentor | Open | Yes |
| OQ-3 | How does the JWKS caching provider behave after the 5-minute TTL when the key endpoint is down, and on an unknown key ID? Not exercised in SPK-02. AC-15 and §8 depend on it. | Pair C | Open | No |
| OQ-4 | If the claim is used (OQ-1) fix tenant role values to exact `learner` / `administrator`, and drop any approximate matching from the production mapper. | Pair C | Open | No |
| OQ-5 | Which claims (`email`, display name) are present in the access token, and under what names? The spec allows both to be null. | Pair C | Open | No |
| OQ-6 | Is the production web sign-in flow hosted login with Authorization Code and PKCE, and what is the logout URL and allowed redirect configuration? | Pair D | Open | No |
| OQ-7 | Agree the error envelope, JSON field naming (`display_name`), user ID type (UUID) and log reason codes. | Pair D | Open | No |
| OQ-8 | Access-token lifetime and clock-skew leeway. These bound how long a blocked user's token still works. | Pair C | Open | No |
| OQ-9 | Agree the (`users`) table model. | Pair D | Open | No |

**Blockers:**

- OQ-1 and OQ-2 must be closed, with decision references, before the specification can be approved.

---

## 12. Review and Approval *(Required)*

| Role | Name or pair | Date | Outcome |
|------|--------------|------|---------|
| Author | Pair D | 2026-11-09 | Submitted |
| Reviewer | Pair A | - | Pending |
| Reviewer | Pair C | - | Pending |

**Approval checklist:**

- [ ] Acceptance criteria are unambiguous and numbered.
- [ ] Every required contract is defined and agreed.
- [ ] Every acceptance criterion has a verification entry.
- [ ] Security considerations are documented.
- [ ] No critical open questions remain.
- [ ] Reviewer is from a different pair.
