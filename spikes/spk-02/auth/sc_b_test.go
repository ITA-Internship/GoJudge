package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestSCB_1_RejectTamperedToken verifies signature verification failure on tampered tokens
func TestSCB_1_RejectTamperedToken(t *testing.T) {
	rig := newTestRig(t)

	validToken := rig.issueFixtureToken(
		t,
		"auth0|learner_test_user",
		ExpectedIssuer,
		ExpectedAudience,
		[]string{"Learner"},
		time.Now().Add(15*time.Minute),
	)
	tamperedToken := validToken[:len(validToken)-6] + "TAMPER"

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tamperedToken)
	rec := httptest.NewRecorder()

	rig.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("SC-B 1/3 failed: expected HTTP 401 for tampered token, got %d", rec.Code)
	}
	t.Logf("rejected: HTTP %d: %s", rec.Code, strings.TrimSpace(rec.Body.String()))
}

// TestSCB_2_RejectExpiredToken verifies rejection of expired tokens
func TestSCB_2_RejectExpiredToken(t *testing.T) {
	rig := newTestRig(t)

	expiredToken := rig.issueFixtureToken(
		t,
		"auth0|learner_test_user",
		ExpectedIssuer,
		ExpectedAudience,
		[]string{"Learner"},
		time.Now().Add(-10*time.Minute),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	rec := httptest.NewRecorder()

	rig.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("SC-B 2/3 failed: expected HTTP 401 for expired token, got %d", rec.Code)
	}
	t.Logf("rejected: HTTP %d: %s", rec.Code, strings.TrimSpace(rec.Body.String()))
}

// TestSCB_3_RejectWrongAudienceToken verifies rejection of mismatched audience tokens
func TestSCB_3_RejectWrongAudienceToken(t *testing.T) {
	rig := newTestRig(t)

	wrongAudToken := rig.issueFixtureToken(
		t,
		"auth0|learner_test_user",
		ExpectedIssuer,
		"https://wrong-audience.unauthorized.local",
		[]string{"Learner"},
		time.Now().Add(15*time.Minute),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+wrongAudToken)
	rec := httptest.NewRecorder()

	rig.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("SC-B 3/3 failed: expected HTTP 401 for wrong-audience token, got %d", rec.Code)
	}
	t.Logf("rejected: HTTP %d: %s", rec.Code, strings.TrimSpace(rec.Body.String()))
}