package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/locales"
	"github.com/jacobpq/soccer-manager/internal/repository"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(sessionRepo *repository.SessionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				api.WriteError(w, http.StatusUnauthorized, locales.T(ctx, "unauthorized"))
				return
			}

			token := authHeader
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}

			if token == "" {
				api.WriteError(w, http.StatusUnauthorized, locales.T(ctx, "invalid_token"))
				return
			}

			userID, err := sessionRepo.GetUserIDByAccessToken(ctx, token)
			if err != nil {
				api.WriteError(w, http.StatusUnauthorized, locales.T(ctx, "invalid_token"))
				return
			}

			ctx = context.WithValue(ctx, UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
