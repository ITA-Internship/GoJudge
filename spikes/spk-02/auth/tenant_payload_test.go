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
	// Retrieve sensitive values from environment variables
	learnerToken := strings.TrimSpace(os.Getenv("AUTH0_LEARNER_TOKEN"))
	learnerEmail := strings.TrimSpace(os.Getenv("AUTH0_LEARNER_EMAIL"))
	adminToken := strings.TrimSpace(os.Getenv("AUTH0_ADMIN_TOKEN"))
	adminEmail := strings.TrimSpace(os.Getenv("AUTH0_ADMIN_EMAIL"))

	// Default expected fallback identities if emails are not explicitly supplied
	if learnerEmail == "" {
		learnerEmail = "learner@test.local"
	}
	if adminEmail == "" {
		adminEmail = "admin@test.local"
	}

	// If no tokens are supplied via environment variables, safely skip live token parsing
	if learnerToken == "" && adminToken == "" {
		t.Skip("Skipping live tenant token inspection: AUTH0_LEARNER_TOKEN and AUTH0_ADMIN_TOKEN are not set")
	}

	parser := jwt.NewParser()

	// 1. Verify and inspect learner token if provided
	if learnerToken != "" {
		t.Run("Inspect_Learner_Token", func(t *testing.T) {
			t.Logf("Verifying token claims for user: %s", learnerEmail)

			var claims CustomClaims
			_, _, err := parser.ParseUnverified(learnerToken, &claims)
			if err != nil {
				t.Fatalf("failed to parse live learner token: %v", err)
			}

			if claims.Issuer != ExpectedIssuer {
				t.Errorf("expected issuer %q, got %q", ExpectedIssuer, claims.Issuer)
			}

			roleFound := false
			for _, r := range claims.Roles {
				if strings.ToLower(strings.TrimSpace(r)) == "learner" {
					roleFound = true
					break
				}
			}
			if !roleFound {
				t.Errorf("expected role 'learner' in claims %v", claims.Roles)
			}
		})
	}

	// 2. Verify and inspect admin token if provided
	if adminToken != "" {
		t.Run("Inspect_Admin_Token", func(t *testing.T) {
			t.Logf("Verifying token claims for user: %s", adminEmail)

			var claims CustomClaims
			_, _, err := parser.ParseUnverified(adminToken, &claims)
			if err != nil {
				t.Fatalf("failed to parse live admin token: %v", err)
			}

			if claims.Issuer != ExpectedIssuer {
				t.Errorf("expected issuer %q, got %q", ExpectedIssuer, claims.Issuer)
			}

			roleFound := false
			for _, r := range claims.Roles {
				normalized := strings.ToLower(strings.TrimSpace(r))
				if normalized == "administrator" || normalized == "adminisrator" {
					roleFound = true
					break
				}
			}
			if !roleFound {
				t.Errorf("expected role 'administrator' in claims %v", claims.Roles)
			}
		})
	}
}