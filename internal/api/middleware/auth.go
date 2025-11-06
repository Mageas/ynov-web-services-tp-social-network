package middleware

import (
	"context"
	"net/http"
	"strings"

	"ynov-social-api/internal/api/response"
	"ynov-social-api/internal/pkg/apperrors"
	"ynov-social-api/internal/service/auth"
)

type contextKey string

const userEmailKey contextKey = "userEmail"

// Auth middleware verifies JWT token and adds user email to context
func Auth(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, apperrors.ErrMissingAuth)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, apperrors.ErrMissingAuth)
				return
			}

			email, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				response.Error(w, apperrors.ErrInvalidToken)
				return
			}

			ctx := context.WithValue(r.Context(), userEmailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserEmail extracts the user email from the request context
func GetUserEmail(r *http.Request) string {
	email, _ := r.Context().Value(userEmailKey).(string)
	return email
}
