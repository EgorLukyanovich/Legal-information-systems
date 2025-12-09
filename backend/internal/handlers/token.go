package handlers

import (
	"context"
	"net/http"
	"os"
	"strings"

	json_resp "github.com/egor_lukyanovich/legal-information-systems/backend/pkg/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"

// ключ из .env если его нет - ставим дефолтный (для тестов)
func getJWTKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("test_secret_key_123")
	}
	return []byte(secret)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("token")

		tokenStr = strings.TrimSpace(tokenStr)

		if tokenStr == "" {
			json_resp.RespondError(w, 401, "UNAUTHORIZED", "missing token")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTKey(), nil
		})

		if err != nil || !token.Valid {
			json_resp.RespondError(w, 401, "UNAUTHORIZED", "invalid token")
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			json_resp.RespondError(w, 401, "UNAUTHORIZED", "invalid token claims")
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			json_resp.RespondError(w, 401, "UNAUTHORIZED", "invalid user id in token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
