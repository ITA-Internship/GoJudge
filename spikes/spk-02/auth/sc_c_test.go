package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestSCC_User1_LearnerMapping verifies claim mapping for test user 1 (learner)
func TestSCC_User1_LearnerMapping(t *testing.T) {
	rig := newTestRig(t)

	token := rig.issueFixtureToken(
		t,
		"auth0|6ae9c01e857a316fec15076a",
		ExpectedIssuer,
		ExpectedAudience,
		[]string{"Learner"},
		time.Now().Add(15*time.Minute),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	rig.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != string(RoleLearner) {
		t.Fatalf("SC-C user 1 failed: expected HTTP 200 with role %q, got %d: %s", RoleLearner, rec.Code, rec.Body.String())
	}
}

// TestSCC_User2_AdministratorMapping verifies claim mapping for test user 2 (administrator)
func TestSCC_User2_AdministratorMapping(t *testing.T) {
	rig := newTestRig(t)

	token := rig.issueFixtureToken(
		t,
		"auth0|6ae9c0656ca5342184adb3a1",
		ExpectedIssuer,
		ExpectedAudience,
		[]string{"Adminisrator"},
		time.Now().Add(15*time.Minute),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	rig.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != string(RoleAdministrator) {
		t.Fatalf("SC-C user 2 failed: expected HTTP 200 with role %q, got %d: %s", RoleAdministrator, rec.Code, rec.Body.String())
	}
}