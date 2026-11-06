# SPK-02: Auth0 OIDC Token Validation and Role Claim Mapping in Go API

## Metadata
- **Spike ID:** SPK-02
- **Type:** Technical spike
- **Owner:** Pair C
- **Timebox:** 2 pair-days (Days 1–2: Mon 2 Nov – Tue 3 Nov)
- **Status:** Completed within timebox
- **Dependencies:** DEC-03 (provisional), MA-05
- **Experiment code:** branch `spike/SPK-02-auth0-oidc`, directory `spikes/spk-02/` (not merged to `main`)
- **SDK Reference:** Official Auth0 Go Quickstart (`github.com/auth0/go-jwt-middleware/v3`)

## 1. Objective
- **Objective:** Answer one question: can the Go API validate Auth0 OIDC tokens and map the role claim to `learner` or `administrator`, using either the Auth0 development tenant or a documented fixture?

## 2. Experiment Setup & Architecture
- **Official Middleware Adoption:** The spike branch uses the official Auth0 middleware suite (`github.com/auth0/go-jwt-middleware/v3`) together with its sub-packages `validator` and `jwks`.
  - **AuthN (JWT Validation):** Handled via `jwtmiddleware.New(jwtmiddleware.WithValidator(v))` and `validator.New(...)`. Key discovery and rotation are managed via `jwks.NewCachingProvider` with a 5-minute cache TTL.
  - **AuthZ & Claim Mapping:** Implemented via a domain-specific `CustomClaims` struct satisfying `validator.CustomClaims` with `Validate(ctx context.Context) error`. Validated claims are retrieved in a type-safe manner from the request context using `jwtmiddleware.GetClaims[*validator.ValidatedClaims](ctx)`.
- **Sign-In Path Used:** Auth0 development tenant (`dev-05kxujwqfcu8g5po.us.auth0.com`). Tokens were issued via Resource Owner Password Flow for two identities (`learner@test.local` and `admin@test.local`). A static key validator fixture (`NewStaticValidator` in `fixture_test.go`) was used for isolated unit verification and negative test cases.
- **Development Tenant Status:** Tested against the live tenant. `TestSCA_LiveTenantToken_ThroughMiddleware` (`sc_a_test.go`) initializes `NewJWKSValidator(ExpectedIssuer)` backed by `jwks.NewCachingProvider` and verifies real tokens supplied via environment variables (`AUTH0_LEARNER_TOKEN` and `AUTH0_ADMIN_TOKEN`).
- **Tenant access (MA-05):** Confirmed with the mentor. No secrets or tokens are stored in the repository.

## 3. Success Criteria & Recorded Evidence

| Criterion | Status | Evidence Summary |
| :--- | :--- | :--- |
| **SC-A** | **Verified** | Live tokens from the Auth0 development tenant were accepted by the official Auth0 middleware (`HTTP 200`), with signatures verified against the tenant's JWKS via `jwks.NewCachingProvider`. Raw run record in §3, SC-A Evidence. |
| **SC-B** | **Verified** | Official middleware validates signature, issuer, audience, and expiry. Rejected 3 of 3 invalid tokens with `HTTP 401`. Raw run record in §3, SC-B Evidence. |
| **SC-C** | **Verified** | Role claims map to `learner` and `administrator` for both test users via `CustomClaims.Validate` (`sc_c_test.go`, `tenant_payload_test.go`). |

### SC-A Evidence: Sign-In Path
- **Tenant:** Tested against the live tenant: **yes**. Raw run record: [`sc_a_live_run.txt`](https://github.com/ITA-Internship/GoJudge/blob/spike/SPK-02-auth0-oidc/spikes/spk-02/evidence/sc_a_live_run.txt) (branch `spike/SPK-02-auth0-oidc`). Command: `go test ./auth -run TestSCA_LiveTenantToken_ThroughMiddleware -v -count=1`. Both subtests (`learner` and `administrator`) passed with `HTTP 200`, and roles were successfully resolved from the verified context.
- **Fixture:** `TestSCA_DemonstrateSignIn` confirms offline operation using `NewStaticValidator` with an ephemeral RSA keypair.
- **Payload Inspection:** `tenant_payload_test.go` verifies the token claims structure without committing sensitive data.

### SC-B Evidence: Three Rejection Results (3 of 3)
Raw test run record: [`sc_b_run.txt`](https://github.com/ITA-Internship/GoJudge/blob/spike/SPK-02-auth0-oidc/spikes/spk-02/evidence/sc_b_run.txt) (branch `spike/SPK-02-auth0-oidc`). Command: `go test ./... -run TestSCB -v -count=1`.
1. **Tampered Token (`TestSCB_1_RejectTamperedToken`):** Corrupted signature suffix. Rejected by official middleware with `HTTP 401` (`unauthorized: jwt invalid: failed to parse and validate token: ... verification error`).
2. **Expired Token (`TestSCB_2_RejectExpiredToken`):** `exp = now - 10m`. Rejected with `HTTP 401` (`unauthorized: jwt invalid: failed to parse and validate token: ... "exp" not satisfied: token is expired`).
3. **Wrong-Audience Token (`TestSCB_3_RejectWrongAudienceToken`):** Audience mismatched. Rejected with `HTTP 401` (`unauthorized: jwt invalid: failed to parse and validate token: ... "aud" not satisfied: token has invalid audience`).

### SC-C Evidence: Role Claim Mapping
Mapping verified for both test users. Claims are extracted from `https://gojudge-api/roles` inside `CustomClaims.Validate()` during token processing.

| User Account | Role Claim (`https://gojudge-api/roles`) | Mapped Internal Role |
| :--- | :--- | :--- |
| `learner@test.local` | `["Learner"]` | `learner` |
| `admin@test.local` | `["Adminisrator"]` *(includes tenant typo handling)* | `administrator` |

The role extractor normalizes (lower-cases and trims) the claim values and maps them directly into `CustomClaims.MappedRole`, which is retrieved via `RoleFromContext(ctx)`.

## 4. Failure Criteria Check
- **F-06 (Middleware accepts tampered or expired token):** **Not met** (Official Auth0 validator rejected all 3 invalid token scenarios in SC-B).
- **F-07 (Neither development tenant nor documented fixture available by day 2):** **Not met** (Both live tenant verification and offline fixtures are fully operational).

## 5. Outcome & Recommendations
- **Outcome:** **Go**. All success criteria (SC-A, SC-B, SC-C) are Verified, and failure criteria were avoided.
- **Recommendation for F01-01:** Adopt `github.com/auth0/go-jwt-middleware/v3` as the foundation for authentication in production. It replaces ~150 lines of bespoke parsing code while offering robust JWKS caching, automatic algorithm pinning (`RS256`), and standard clock-skew handling out of the box.
- **Follow-On Decision:** **DEC-03** (Authentication with Auth0) is confirmed viable. Formal approval and ADR recording remain with the mentor.