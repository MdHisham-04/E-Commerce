package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MdHisham-04/E-Commerce/internal/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware checks if user has required role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(UserContextKey).(*auth.Claims)
			if claims.Role != role {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user claims from request context
func GetUserFromContext(r *http.Request) *auth.Claims {
	claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}
