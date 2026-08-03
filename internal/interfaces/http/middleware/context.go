package middleware

import (
	"context"

	"github.com/google/uuid"

	"github.com/GitAlex9/go-order-service/internal/pkg/jwt"
)

type contextKey string

const claimsContextKey contextKey = "claims"

func WithClaims(ctx context.Context, claims *jwt.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*jwt.Claims)
	return claims, ok
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return uuid.Nil, false
	}
	return claims.UserID, true
}
