package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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