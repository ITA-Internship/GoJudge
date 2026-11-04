package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/auth0/go-jwt-middleware/v3/validator"
)

// TestSCA_DemonstrateSignIn verifies that sign-in works and produces a token accepted by middleware (Demonstrated once)
func TestSCA_DemonstrateSignIn(t *testing.T) {
	rig := newTestRig(t)

	token := rig.issueFixtureToken(
		t,
		"auth0|6ae9c01e857a316fec15076a",
		ExpectedIssuer,
		ExpectedAudience,
		[]string{"Learner"},
		time.Now().Add(15*time.Minute),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	rig.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("SC-A failed: expected HTTP 200 OK on valid sign-in token, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestSCA_LiveTenantToken_ThroughMiddleware presents real Auth0 development tenant tokens to the
// official Auth0 middleware, verified against the tenant's published signing keys via jwks.NewCachingProvider.
// Tokens come from AUTH0_LEARNER_TOKEN and AUTH0_ADMIN_TOKEN and are never logged. Skipped when neither is set.
func TestSCA_LiveTenantToken_ThroughMiddleware(t *testing.T) {
	cases := []struct {
		name string
		env  string
		want Role
	}{
		{"learner", "AUTH0_LEARNER_TOKEN", RoleLearner},
		{"administrator", "AUTH0_ADMIN_TOKEN", RoleAdministrator},
	}

	var jwtValidator *validator.Validator
	ran := false

	for _, c := range cases {
		token := strings.TrimSpace(os.Getenv(c.env))
		if token == "" {
			continue
		}
		ran = true

		if jwtValidator == nil {
			var err error
			jwtValidator, err = NewJWKSValidator(ExpectedIssuer)
			if err != nil {
				t.Fatalf("failed to initialize JWKS validator: %v", err)
			}
		}

		t.Run(c.name, func(t *testing.T) {
			handler := AuthMiddleware(jwtValidator, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				role, _ := RoleFromContext(r.Context())
				_, _ = w.Write([]byte(role))
			}))

			req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("live %s token rejected by middleware: HTTP %d: %s", c.name, rec.Code, rec.Body.String())
			}
			if got := Role(rec.Body.String()); got != c.want {
				t.Fatalf("live %s token mapped to %q, want %q", c.name, got, c.want)
			}
			t.Logf("live %s token verified against tenant JWKS: HTTP 200, role=%s", c.name, c.want)
		})
	}

	if !ran {
		t.Skip("AUTH0_LEARNER_TOKEN and AUTH0_ADMIN_TOKEN are not set")
	}
}