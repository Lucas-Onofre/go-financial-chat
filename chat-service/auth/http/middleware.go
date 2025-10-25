package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Lucas-Onofre/financial-chat/chat-service/auth/jwt"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	ClaimsKey contextKey = "claims"
)

func AuthMiddleware(jwtService *jwt.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, http.StatusUnauthorized, "Authorization header missing")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				writeJSONError(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}

			tokenString := parts[1]
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}
