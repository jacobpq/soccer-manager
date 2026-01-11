package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jacobpq/soccer-manager/internal/api"
	"github.com/jacobpq/soccer-manager/internal/config"
	"github.com/jacobpq/soccer-manager/internal/domain/models"
	"github.com/jacobpq/soccer-manager/internal/locales"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				api.WriteError(w, http.StatusUnauthorized, locales.T(ctx, "unauthorized"))
				return
			}

			tokenString := authHeader
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}

			token, err := jwt.ParseWithClaims(tokenString, &models.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				api.WriteError(w, http.StatusUnauthorized, locales.T(ctx, "invalid_token"))
				return
			}

			if claims, ok := token.Claims.(*models.AuthClaims); ok {
				if claims.TokenType != "access" {
					api.WriteError(w, http.StatusUnauthorized, locales.T(ctx, "invalid_token"))
					return
				}

				ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				api.WriteError(w, http.StatusUnauthorized, locales.T(ctx, "invalid_token"))
			}

		})
	}
}
