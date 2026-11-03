package main

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Constants matching the Auth0 development tenant tokens
const (
	ExpectedIssuer   = "https://dev-05kxujwqfcu8g5po.us.auth0.com/"
	ExpectedAudience = "https://gojudge-api"
	RoleClaimKey     = "https://gojudge-api/roles"
)

// Role defines internal application user roles
type Role string

const (
	RoleLearner       Role = "learner"
	RoleAdministrator Role = "administrator"
)

type contextKey string

const roleCtxKey contextKey = "userRole"

// CustomClaims represents Auth0 JWT payload claims
type CustomClaims struct {
	Roles []string `json:"https://gojudge-api/roles"`
	jwt.RegisteredClaims
}

// MakeRSAKeyFunc wraps an RSA public key into a jwt.Keyfunc
func MakeRSAKeyFunc(pubKey *rsa.PublicKey) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing algorithm")
		}
		return pubKey, nil
	}
}

// ValidateAndMapToken validates token signature, issuer, audience, expiry,
// and extracts the namespaced role claim into an internal Role.
func ValidateAndMapToken(tokenString string, keyFunc jwt.Keyfunc) (Role, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		keyFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(ExpectedIssuer),
		jwt.WithAudience(ExpectedAudience),
		jwt.WithExpirationRequired(),
	)

	if err != nil || !token.Valid {
		return "", err
	}

	for _, r := range claims.Roles {
		normalized := strings.ToLower(strings.TrimSpace(r))
		// Handles standard role string as well as tenant typo "adminisrator"
		if normalized == "learner" {
			return RoleLearner, nil
		}
		if normalized == "administrator" || normalized == "adminisrator" {
			return RoleAdministrator, nil
		}
	}

	return "", errors.New("unmapped or missing role claim")
}

// AuthMiddleware inspects incoming requests for Bearer tokens and injects the Role into context
func AuthMiddleware(keyFunc jwt.Keyfunc, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		role, err := ValidateAndMapToken(rawToken, keyFunc)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), roleCtxKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleFromContext extracts the parsed Role from request context
func RoleFromContext(ctx context.Context) (Role, bool) {
	role, ok := ctx.Value(roleCtxKey).(Role)
	return role, ok
}