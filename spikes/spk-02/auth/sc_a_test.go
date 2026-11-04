package main

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

// jwksKeyFunc resolves the RSA public key for a token's kid from the tenant's JWKS endpoint.
func jwksKeyFunc(t *testing.T, issuer string) jwt.Keyfunc {
	t.Helper()
	resp, err := http.Get(issuer + ".well-known/jwks.json")
	if err != nil {
		t.Fatalf("failed to fetch tenant JWKS: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("tenant JWKS returned HTTP %d", resp.StatusCode)
	}

	var set struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		t.Fatalf("failed to decode tenant JWKS: %v", err)
	}

	keys := map[string]*rsa.PublicKey{}
	for _, k := range set.Keys {
		if k.Kty != "RSA" {
			continue
		}
		n, errN := base64.RawURLEncoding.DecodeString(k.N)
		e, errE := base64.RawURLEncoding.DecodeString(k.E)
		if errN != nil || errE != nil {
			continue
		}
		keys[k.Kid] = &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}
	}

	return func(token *jwt.Token) (interface{}, error) {
		kid, _ := token.Header["kid"].(string)
		pub, ok := keys[kid]
		if !ok {
			return nil, errors.New("no JWKS key for kid")
		}
		return pub, nil
	}
}

// TestSCA_LiveTenantToken_ThroughMiddleware presents real Auth0 development tenant tokens to the
// middleware, verified against the tenant's published signing keys (JWKS). Tokens come from
// AUTH0_LEARNER_TOKEN and AUTH0_ADMIN_TOKEN and are never logged. Skipped when neither is set.
func TestSCA_LiveTenantToken_ThroughMiddleware(t *testing.T) {
	cases := []struct {
		name string
		env  string
		want Role
	}{
		{"learner", "AUTH0_LEARNER_TOKEN", RoleLearner},
		{"administrator", "AUTH0_ADMIN_TOKEN", RoleAdministrator},
	}

	keyFunc := (jwt.Keyfunc)(nil)
	ran := false
	for _, c := range cases {
		token := strings.TrimSpace(os.Getenv(c.env))
		if token == "" {
			continue
		}
		ran = true
		if keyFunc == nil {
			keyFunc = jwksKeyFunc(t, ExpectedIssuer)
		}
		t.Run(c.name, func(t *testing.T) {
			handler := AuthMiddleware(keyFunc, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
