package requestValidator

import (
	"context"
	"fmt"
	"mailinglist-backend-go/services/configReader"
	"mailinglist-backend-go/services/jwtValidator"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name     string
	LastName string
	Email    string
	Admin    bool
}

// normalizePublicKey takes the env value and returns a PEM-formatted public key string.
// Supports three formats in KEYCLOAK_PUBLIC_KEY:
// 1) Full PEM including BEGIN/END lines (possibly with \n escaped) -> used as-is (after unescaping \n)
// 2) Single-line base64 body, no headers -> headers/footers will be added
// 3) Multi-line body without headers (rare) -> headers/footers will be added
func normalizePublicKey(raw string) string {
	if raw == "" {
		return raw
	}

	// Remove surrounding quotes if present (some .env examples quote the value)
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "\"")
	trimmed = strings.TrimSuffix(trimmed, "\"")

	// Allow storing with literal \n sequences in .env
	trimmed = strings.ReplaceAll(trimmed, "\\n", "\n")

	// If already contains PEM header, assume it's complete
	if strings.Contains(trimmed, "BEGIN PUBLIC KEY") && strings.Contains(trimmed, "END PUBLIC KEY") {
		return trimmed
	}

	// Otherwise, wrap with PEM headers/footers (backward compatibility)
	return "-----BEGIN PUBLIC KEY-----\n" + trimmed + "\n-----END PUBLIC KEY-----"
}

func ValidateRequest(r *http.Request) (jwt.MapClaims, error) {
	publicKey := configReader.Value("KEYCLOAK_PUBLIC_KEY")
	publicKeyComplete := normalizePublicKey(publicKey)
	bearerToken := r.Header.Get("Authorization")
	token := strings.Split(bearerToken, "Bearer ")

	if len(token) < 2 {
		return nil, fmt.Errorf("No token found in header")
	}

	return jwtValidator.ValidateToken(token[1], publicKeyComplete)
}

// Context helpers for attaching and retrieving claims set by auth middleware

type ctxKey string

var claimsCtxKey ctxKey = "jwtClaims"

// WithClaims returns a new context with JWT claims stored
func WithClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, claimsCtxKey, claims)
}

// ClaimsFromContext extracts JWT claims from context
func ClaimsFromContext(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(claimsCtxKey).(jwt.MapClaims)
	return claims, ok
}

// ClaimsFromRequest is a helper to extract claims from the request context
func ClaimsFromRequest(r *http.Request) (jwt.MapClaims, error) {
	if claims, ok := ClaimsFromContext(r.Context()); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("no claims in context")
}

func isAdmin(claims jwt.MapClaims) bool {
	var groups []interface{}
	groups = claims["groups"].([]interface{})

	for _, group := range groups {
		if group == "Admin" {
			return true
		}
	}
	return false
}

func CurrentUser(claims jwt.MapClaims) User {
	return User{
		Name:     claims["given_name"].(string),
		LastName: claims["family_name"].(string),
		Email:    claims["email"].(string),
		Admin:    isAdmin(claims),
	}
}
