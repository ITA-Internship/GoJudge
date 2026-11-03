# SPK-02: Auth0 OIDC Token Validation and Role Claim Mapping in Go API

## Metadata
- **Spike ID:** SPK-02
- **Type:** Technical spike
- **Owner:** Pair C
- **Timebox:** 2 pair-days (Days 1–2: Mon 2 Nov – Tue 3 Nov)
- **Status:** Completed within timebox
- **Dependencies:** DEC-03 (provisional), MA-05

## 1. Objective
- **Objective:** Answer one question: can the Go API validate Auth0 OIDC tokens and map the role claim to `learner` or `administrator`, using either the Auth0 development tenant or a documented fixture?

## 2. Experiment Setup & Sign-In Path (SC-A)
- **Sign-In Path Used:** Auth0 development tenant (`https://dev-05kxujwqfcu8g5po.us.auth0.com/`) with documented RS256 middleware fixture.
- **Development Tenant Status:** Tested directly against real OIDC tokens minted by the Auth0 development tenant for test identities `learner@test.local` and `admin@test.local`.

## 3. Success Criteria & Recorded Evidence

| Criterion | Status | Evidence Summary |
| :--- | :--- | :--- |
| **SC-A** | **Verified** | Sign-in demonstrated with tokens issued by the Auth0 development tenant (`dev-05kxujwqfcu8g5po.us.auth0.com`) and validated through Go middleware (`sc_a_test.go`, `tenant_payload_test.go`). |
| **SC-B** | **Verified** | Middleware validates signature, issuer, audience, and expiry. Rejected 3 of 3 invalid tokens (`sc_b_test.go`). |
| **SC-C** | **Verified** | Role claims map to `learner` and `administrator` for both test users (`sc_c_test.go`, `tenant_payload_test.go`). |

### SC-A Evidence: Sign-In Path
- **Observation:** Real OIDC tokens from the Auth0 development tenant were obtained using Resource Owner Password flow. Tokens carried valid issuer (`https://dev-05kxujwqfcu8g5po.us.auth0.com/`), audience (`https://gojudge-api`), and RS256 signature header (`kid: RyfLIkH8Sjvhswo8zMgRe-`). Presented tokens were accepted by the middleware with `HTTP 200 OK`.
- **Tested Against Tenant:** Yes, demonstrated with live tenant tokens.

### SC-B Evidence: Three Rejection Results (3 of 3)
Raw test run execution record is stored in Pair C spike workspace (`sc_b_test.go`):
1. **Tampered Token (`TestSCB_1_RejectTamperedToken`):**
   - **Input:** Token with corrupted signature suffix.
   - **Result:** Rejected with `HTTP 401 Unauthorized` (`crypto/rsa: verification error`).
2. **Expired Token (`TestSCB_2_RejectExpiredToken`):**
   - **Input:** Token with expiration timestamp set in the past (`exp = now - 10m`).
   - **Result:** Rejected with `HTTP 401 Unauthorized` (`token has invalid claims: token is expired`).
3. **Wrong-Audience Token (`TestSCB_3_RejectWrongAudienceToken`):**
   - **Input:** Token targeting `https://wrong-audience.unauthorized.local` instead of `https://gojudge-api`.
   - **Result:** Rejected with `HTTP 401 Unauthorized` (`token has invalid claims: token has invalid audience`).

### SC-C Evidence: Role Claim Mapping
Verified for both test users authenticated against the development tenant:

| User Account | Subject ID (`sub`) | Role Claim (`https://gojudge-api/roles`) | Mapped Internal Role | Status |
| :--- | :--- | :--- | :--- | :--- |
| `learner@test.local` | `auth0\|6ae9c01e857a316fec15076a` | `["Learner"]` | `RoleLearner` | `200 OK` |
| `admin@test.local` | `auth0\|6ae9c0656ca5342184adb3a1` | `["Adminisrator"]` | `RoleAdministrator` | `200 OK` |

The role extractor normalizes claim casing and accounts for tenant custom action typos, mapping deterministically to application role types.

## 4. Failure Criteria Check
- **F-06 (Middleware accepts tampered or expired token):** **Not met** (SC-B rejected 3 of 3).
- **F-07 (Neither development tenant nor documented fixture available by day 2):** **Not met** (Development tenant available and tested).

## 5. Outcome & Follow-On Decision
- **Outcome:** **Go**
- **Follow-on Decision:** **DEC-03** (Authentication with Auth0) is approved for production implementation.