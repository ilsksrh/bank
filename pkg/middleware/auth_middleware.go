package middleware

import (
	"context"
	"jusan_demo/pkg/auth"
	"net/http"
	"strings"
)

type contextKey string

const RoleKey contextKey = "userRole"

func AuthMiddleware(authService *auth.AuthService, requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authService.VerifyAccessToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			allowed := false
			for _, rr := range requiredRoles {
				if rr == claims.Role {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), RoleKey, claims.Role)
			ctx = context.WithValue(ctx, "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
