package auth

import (
	"context"
	"net/http"
	"strings"
)

// Define a key type for context to avoid collisions
type contextKey string

const userKey contextKey = "username"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		claims, err := VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		// store username in context
		ctx := context.WithValue(r.Context(), userKey, username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
