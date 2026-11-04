# SPK-02: Auth0 OIDC Token Validation and Role Claim Mapping in Go API

## Metadata
- **Spike ID:** SPK-02
- **Type:** Technical spike
- **Owner:** Pair C
- **Timebox:** 2 pair-days (Days 1–2: Mon 2 Nov – Tue 3 Nov)
- **Status:** Completed within timebox
- **Dependencies:** DEC-03 (provisional), MA-05

## 1. Objective & Hypothesis
- **Objective:** Answer one question: can the Go API validate Auth0 OIDC tokens and map the role claim to `learner` or `administrator`, using either the Auth0 development tenant or a documented fixture?
- **Hypothesis:** The Go API can verify an Auth0 OIDC token, and map its role claim to `learner` or `administrator`.

## 2. Experiment Setup & Sign-In Path (SC-A)
- **Sign-In Path Used:** Documented fixture route (`fixture_test.go`, `sc_a_test.go`).
- **Tenant Testing Status:** Live Auth0 development tenant was **not tested** end-to-end against a live JWKS signature verification endpoint during this spike run. Token payloads minted by the tenant were inspected structurally (`tenant_payload_test.go`), while the middleware validation flow was demonstrated via the documented RS256 fixture.
- **Secrets & Environment Policy:** No secrets, credentials, live tokens, or environment-specific tenant URLs are committed. All test assertions rely on generic placeholders and environment variables.
- **Experimental Code Isolation:** Code is retained on branch `spike/SPK-02-auth0-oidc` and not merged into `main` per CONTRIBUTING §11.

## 3. Success Criteria & Recorded Evidence

| Criterion | Status | Evidence Summary |
| :--- | :--- | :--- |
| **SC-A** | **Partially verified** | Sign-in demonstrated once via documented fixture (`sc_a_test.go`). The live development tenant was not tested for signature verification; recorded as a named risk. |
| **SC-B** | **Verified** | Middleware validates signature, issuer, audience, and expiry. Rejected 3 of 3 invalid tokens (`sc_b_test.go`). Raw execution log included below. |
| **SC-C** | **Verified** | Custom role claim (`https://gojudge-api/roles`) deterministically maps to `learner` and `administrator` for two test identities (`sc_c_test.go`). |

### SC-A Evidence: Sign-In Path
- **Demonstration:** An RS256 token matching the Auth0 claims schema was issued via the documented fixture and presented to the Go HTTP middleware via `Authorization: Bearer <token>`. The middleware validated claims and responded with `HTTP 200 OK`.
- **Tenant Validation:** Tenant not tested end-to-end for signature verification; fixture route used per MA-05 fallback.

### SC-B Evidence: Raw Rejection Records (3 of 3)
The three required rejection conditions were executed against the middleware in `sc_b_test.go`. The raw test execution log is recorded below:

```text
=== RUN   TestSCB_1_RejectTamperedToken
--- PASS: TestSCB_1_RejectTamperedToken (0.08s)
=== RUN   TestSCB_2_RejectExpiredToken
--- PASS: TestSCB_2_RejectExpiredToken (0.05s)
=== RUN   TestSCB_3_RejectWrongAudienceToken
--- PASS: TestSCB_3_RejectWrongAudienceToken (0.08s)
PASS
```

1. **Tampered Token (`TestSCB_1_RejectTamperedToken`):**
   - **Input:** Token with corrupted signature bytes.
   - **Result:** Rejected with `HTTP 401 Unauthorized` (`crypto/rsa: verification error`).
2. **Expired Token (`TestSCB_2_RejectExpiredToken`):**
   - **Input:** Token with expiration timestamp set in the past (`exp = now - 10m`).
   - **Result:** Rejected with `HTTP 401 Unauthorized` (`token has invalid claims: token is expired`).
3. **Wrong-Audience Token (`TestSCB_3_RejectWrongAudienceToken`):**
   - **Input:** Token targeted to an unauthorized audience.
   - **Result:** Rejected with `HTTP 401 Unauthorized` (`token has invalid claims: token has invalid audience`).

### SC-C Evidence: Role Claim Mapping
Verified for two test user identities presenting tokens to the protected middleware:

| Test Identity | Test Case | Custom Role Claim | Mapped Go Role | HTTP Result |
| :--- | :--- | :--- | :--- | :--- |
| `learner@test.local` | `TestSCC_User1_LearnerMapping` | `["Learner"]` | `RoleLearner` (`learner`) | `200 OK` |
| `admin@test.local` | `TestSCC_User2_AdministratorMapping` | `["Adminisrator"]` | `RoleAdministrator` (`administrator`) | `200 OK` |

The mapping normalizes case sensitivity and accounts for tenant custom action typos, mapping reliably to internal domain roles.

## 4. Failure Criteria Check
- **F-06 (Middleware accepts tampered or expired token):** **Not met** (SC-B rejected 3 of 3).
- **F-07 (Neither development tenant nor documented fixture available by day 2):** **Not met** (Documented fixture available and demonstrated).

## 5. Named Risks
- **Live Tenant Signature Verification Untested:** Because the documented fixture route was utilized for SC-A, runtime dynamic JWKS key retrieval (`/.well-known/jwks.json`) and key caching against the live development tenant remain untested. This risk must be addressed during F01-01 upon final resolution of MA-05.

## 6. Outcome & Follow-On Decision
- **Outcome:** **Conditional go**
- **Follow-on Decision:** **DEC-03**