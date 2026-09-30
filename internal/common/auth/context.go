package auth

import (
	"context"
	"errors"
)

type contextKey int

const claimsContextKey contextKey = iota

var ErrNoClaimsInContext = errors.New("no claims in context")

func contextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// ClaimsFromContext retrieves the authenticated caller's claims, put there
// by Middleware earlier in the request's handler chain.
func ClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	if !ok || claims == nil {
		return nil, ErrNoClaimsInContext
	}
	return claims, nil
}
