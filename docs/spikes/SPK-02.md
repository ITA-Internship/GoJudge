# SPK-02: Auth0 OIDC Token Validation and Role Claim Mapping in Go API

## Metadata
- **Spike ID:** SPK-02
- **Type:** Technical spike
- **Owner:** Pair C
- **Timebox:** 2 pair-days (Days 1–2: Mon 2 Nov – Tue 3 Nov)
- **Status:** Completed within timebox
- **Dependencies:** DEC-03 (provisional), MA-05
- **Experiment code:** branch `spike/SPK-02-auth0-oidc`, directory `spikes/spk-02/` (not merged to `main`)

## 1. Objective & Hypothesis
- **Objective:** Answer one question: can the Go API validate Auth0 OIDC tokens and map the role claim to `learner` or `administrator`, using either the Auth0 development tenant or a documented fixture?
- **Hypothesis:** The Go API can verify an Auth0 OIDC token, and map its role claim to `learner` or `administrator`.

## 2. Experiment Setup & Sign-In Path (SC-A)
- **Sign-In Path Used:** Auth0 development tenant (`<dev-tenant>.auth0.com`). Real tokens were minted with the Resource Owner Password flow for two test identities. A locally generated RSA fixture (`fixture_test.go`) was used for the negative cases in SC-B.
- **Development Tenant Status:** Tested against the tenant. `TestSCA_LiveTenantToken_ThroughMiddleware` (`sc_a_test.go`) fetches the tenant's JWKS and passes the real learner and administrator tokens (from `AUTH0_LEARNER_TOKEN` and `AUTH0_ADMIN_TOKEN`, not stored) through the middleware, which verifies signature, issuer, audience and expiry.
- **Tenant access (MA-05):** agreed verbally with the mentor at the daily on 3 Nov 2026. No written record exists yet (see §5).
- **Secrets & Environment Policy:** No secrets, credentials, live tokens or environment-specific tenant URLs are committed. Tokens are supplied through environment variables and are not logged.

## 3. Success Criteria & Recorded Evidence

| Criterion | Status | Evidence Summary |
| :--- | :--- | :--- |
| **SC-A** | **Verified** | Real tokens from the Auth0 development tenant were accepted by the middleware (`HTTP 200`), with the signature verified against the tenant JWKS. Raw run record in §3, SC-A Evidence. |
| **SC-B** | **Verified** | Middleware validates signature, issuer, audience, and expiry. Rejected 3 of 3 invalid tokens. Raw run record in §3, SC-B Evidence. |
| **SC-C** | **Verified** | Role claims map to `learner` and `administrator` for both test users (`sc_c_test.go`, `tenant_payload_test.go`). |

### SC-A Evidence: Sign-In Path
- **Tenant:** Tested against the tenant: **yes**. Raw run record (secret-free, tokens are not logged): [`sc_a_live_run.txt`](https://github.com/ITA-Internship/GoJudge/blob/8b85f4d4c735f0897e7d61e2c6e55bba355a7c3b/spikes/spk-02/evidence/sc_a_live_run.txt) (branch `spike/SPK-02-auth0-oidc`, commit `8b85f4d`). Command: `go test ./auth -run SCA_Live -v -count=1`. Result: `PASS` for both subtests. The learner and administrator tokens each returned `HTTP 200` with the signature verified against the tenant JWKS.
- **Fixture:** `TestSCA_DemonstrateSignIn` repeats the check offline with a locally signed token. It is a regression check, not the evidence for SC-A.
- **Earlier check:** `tenant_payload_test.go` only parses the payload without verifying the signature (`ParseUnverified`). It is not used as evidence for SC-A.

### SC-B Evidence: Three Rejection Results (3 of 3)
Raw test run record (secret-free): [`sc_b_run.txt`](https://github.com/ITA-Internship/GoJudge/blob/8b85f4d4c735f0897e7d61e2c6e55bba355a7c3b/spikes/spk-02/evidence/sc_b_run.txt) (branch `spike/SPK-02-auth0-oidc`, commit `8b85f4d`). It was produced with `go test ./... -run TestSCB -v -count=1` on 4 Nov 2026 (Go 1.27.3).
1. **Tampered Token (`TestSCB_1_RejectTamperedToken`):** token with corrupted signature suffix. Rejected with `HTTP 401` (`token signature is invalid: crypto/rsa: verification error`).
2. **Expired Token (`TestSCB_2_RejectExpiredToken`):** `exp = now - 10m`. Rejected with `HTTP 401` (`token has invalid claims: token is expired`).
3. **Wrong-Audience Token (`TestSCB_3_RejectWrongAudienceToken`):** audience `https://wrong-audience.unauthorized.local` instead of the expected one. Rejected with `HTTP 401` (`token has invalid claims: token has invalid audience`).

### SC-C Evidence: Role Claim Mapping
Mapping checked for two test users. Claims come from the development tenant. The live run in [`sc_a_live_run.txt`](https://github.com/ITA-Internship/GoJudge/blob/8b85f4d4c735f0897e7d61e2c6e55bba355a7c3b/spikes/spk-02/evidence/sc_a_live_run.txt) also shows the mapped roles (`role=learner`, `role=administrator`) after full verification. Mapping code: `middleware.go`.

| User Account | Role Claim (`https://gojudge-api/roles`) | Mapped Internal Role |
| :--- | :--- | :--- |
| `learner@test.local` | `["Learner"]` | `learner` |
| `admin@test.local` | `["Administrator"]` | `administrator` |

The role extractor lower-cases and trims the claim value before mapping it.

## 4. Failure Criteria Check
- **F-06 (Middleware accepts tampered or expired token):** **Not met** (SC-B rejected 3 of 3).
- **F-07 (Neither development tenant nor documented fixture available by day 2):** **Not met** (development tenant available and used; fixture also available).

## 5. Outcome & Follow-On Decision
- **Outcome:** **Go**. SC-A, SC-B and SC-C are all Verified, and neither failure criterion was met.
- **Follow-on decision:** **DEC-03** (Authentication with Auth0). Recommendation: approve Auth0 for F01-01. MA-05 (tenant access) was agreed verbally with the mentor on 3 Nov 2026. DEC-03 stays **provisional** in this report. Its approval is the mentor's decision and is to be recorded outside this report (decision record or ADR, as the mentor decides).