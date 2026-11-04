package main

import (
	"os"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

// TestTenantTokens_UnverifiedPayloadInspection dynamically reads tenant tokens
// and email identities from environment variables to avoid committing secrets into code.
func TestTenantTokens_UnverifiedPayloadInspection(t *testing.T) {
	learnerToken := strings.TrimSpace(os.Getenv("AUTH0_LEARNER_TOKEN"))
	learnerEmail := strings.TrimSpace(os.Getenv("AUTH0_LEARNER_EMAIL"))
	adminToken := strings.TrimSpace(os.Getenv("AUTH0_ADMIN_TOKEN"))
	adminEmail := strings.TrimSpace(os.Getenv("AUTH0_ADMIN_EMAIL"))

	if learnerEmail == "" {
		learnerEmail = "learner@test.local"
	}
	if adminEmail == "" {
		adminEmail = "admin@test.local"
	}

	if learnerToken == "" && adminToken == "" {
		t.Skip("Skipping live tenant token inspection: AUTH0_LEARNER_TOKEN and AUTH0_ADMIN_TOKEN are not set")
	}

	parser := jwt.NewParser()

	// 1. Verify and inspect learner token if provided
	if learnerToken != "" {
		t.Run("Inspect_Learner_Token", func(t *testing.T) {
			t.Logf("Verifying token claims for user: %s", learnerEmail)
			var claims fixtureClaims
			_, _, err := parser.ParseUnverified(learnerToken, &claims)
			if err != nil {
				t.Fatalf("failed to parse live learner token: %v", err)
			}
			if claims.Issuer != ExpectedIssuer {
				t.Errorf("expected issuer %q, got %q", ExpectedIssuer, claims.Issuer)
			}
			role, err := MapRolesToInternal(claims.Roles)
			if err != nil || role != RoleLearner {
				t.Errorf("expected role %q in claims %v, got %q (err: %v)", RoleLearner, claims.Roles, role, err)
			}
		})
	}

	// 2. Verify and inspect admin token if provided
	if adminToken != "" {
		t.Run("Inspect_Admin_Token", func(t *testing.T) {
			t.Logf("Verifying token claims for user: %s", adminEmail)
			var claims fixtureClaims
			_, _, err := parser.ParseUnverified(adminToken, &claims)
			if err != nil {
				t.Fatalf("failed to parse live admin token: %v", err)
			}
			if claims.Issuer != ExpectedIssuer {
				t.Errorf("expected issuer %q, got %q", ExpectedIssuer, claims.Issuer)
			}
			role, err := MapRolesToInternal(claims.Roles)
			if err != nil || role != RoleAdministrator {
				t.Errorf("expected role %q in claims %v, got %q (err: %v)", RoleAdministrator, claims.Roles, role, err)
			}
		})
	}
}