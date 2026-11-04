package main

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/jwks"
	"github.com/auth0/go-jwt-middleware/v3/validator"
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

// CustomClaims represents Auth0 JWT custom payload claims and implements validator.CustomClaims.
type CustomClaims struct {
	Roles      []string `json:"https://gojudge-api/roles"`
	MappedRole Role     `json:"-"`
}

// Validate is called automatically by the official Auth0 validator after parsing the JWT.
func (c *CustomClaims) Validate(ctx context.Context) error {
	role, err := MapRolesToInternal(c.Roles)
	if err != nil {
		return err
	}
	c.MappedRole = role
	return nil
}

// MapRolesToInternal extracts the namespaced role claim into an internal Role.
func MapRolesToInternal(roles []string) (Role, error) {
	for _, r := range roles {
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

// MakeStaticKeyFunc wraps an RSA public key into a keyFunc compatible with validator.WithKeyFunc.
func MakeStaticKeyFunc(pubKey *rsa.PublicKey) func(context.Context) (any, error) {
	return func(ctx context.Context) (any, error) {
		return pubKey, nil
	}
}

// NewStaticValidator creates an official Auth0 JWT validator backed by an in-memory RSA public key (for fixtures).
func NewStaticValidator(pubKey *rsa.PublicKey) (*validator.Validator, error) {
	return newValidatorWithKeyFunc(MakeStaticKeyFunc(pubKey))
}

// NewJWKSValidator creates an official Auth0 JWT validator using the official Caching JWKS Provider.
func NewJWKSValidator(issuer string) (*validator.Validator, error) {
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issuer URL: %w", err)
	}

	provider, err := jwks.NewCachingProvider(
		jwks.WithIssuerURL(issuerURL),
		jwks.WithCacheTTL(5*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS provider: %w", err)
	}

	return newValidatorWithKeyFunc(provider.KeyFunc)
}

func newValidatorWithKeyFunc(keyFunc func(context.Context) (any, error)) (*validator.Validator, error) {
	return validator.New(
		validator.WithKeyFunc(keyFunc),
		validator.WithAlgorithm(validator.RS256),
		validator.WithIssuer(ExpectedIssuer),
		validator.WithAudience(ExpectedAudience),
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &CustomClaims{}
		}),
		validator.WithAllowedClockSkew(30*time.Second),
	)
}

// AuthMiddleware wraps the official Auth0 go-jwt-middleware/v3 CheckJWT handler.
func AuthMiddleware(jwtValidator *validator.Validator, next http.Handler) http.Handler {
	mw, err := jwtmiddleware.New(
		jwtmiddleware.WithValidator(jwtValidator),
		jwtmiddleware.WithValidateOnOptions(false),
		jwtmiddleware.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		}),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize official Auth0 jwtmiddleware: %v", err))
	}

	return mw.CheckJWT(next)
}

// RoleFromContext extracts the validated CustomClaims from the official Auth0 middleware context.
func RoleFromContext(ctx context.Context) (Role, bool) {
	claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](ctx)
	if err != nil || claims == nil {
		return "", false
	}

	custom, ok := claims.CustomClaims.(*CustomClaims)
	if !ok || custom == nil || custom.MappedRole == "" {
		return "", false
	}

	return custom.MappedRole, true
}