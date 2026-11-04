package main

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type testRig struct {
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
	handler http.Handler
}

// newTestRig boots an in-memory fixture environment
func newTestRig(t *testing.T) *testRig {
	t.Helper()
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed generating test RSA keypair: %v", err)
	}

	pubKey := &privKey.PublicKey
	keyFunc := MakeRSAKeyFunc(pubKey)

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := RoleFromContext(r.Context())
		if !ok {
			http.Error(w, "role missing from context", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(role))
	})

	return &testRig{
		privKey: privKey,
		pubKey:  pubKey,
		handler: AuthMiddleware(keyFunc, protectedHandler),
	}
}

// issueFixtureToken mints signed RS256 OIDC tokens for testing
func (r *testRig) issueFixtureToken(t *testing.T, sub, iss, aud string, roles []string, exp time.Time) string {
	t.Helper()
	claims := CustomClaims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			Issuer:    iss,
			Audience:  jwt.ClaimStrings{aud},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(r.privKey)
	if err != nil {
		t.Fatalf("failed to sign fixture token: %v", err)
	}
	return signed
}
