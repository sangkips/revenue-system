package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRoleKey contextKey = "user_role"
const UserCountyIDKey contextKey = "user_county_id"

func JWTAuth(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrNotSupported
				}
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				log.Error().Err(err).Msg("Invalid JWT")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "Invalid user_id in token", http.StatusUnauthorized)
				return
			}

			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "Invalid role in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)

			if countyIDFloat, exists := claims["county_id"]; exists {
				if countyID, ok := countyIDFloat.(float64); ok {
					ctx = context.WithValue(ctx, UserCountyIDKey, int32(countyID))
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
