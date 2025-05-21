package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jhonVitor-rs/url-shortener/internal/api/hooks"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

var jwtSecret = []byte(os.Getenv("MY_SECRET_KEY"))

type contextKey string

const userIDKey contextKey = "user_id"

func GenerateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			hooks.SendResponse(w, http.StatusUnauthorized, nil, wraperrors.UnauthorizedErr("Missing or invalid token"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			hooks.SendResponse(w, http.StatusUnauthorized, nil, wraperrors.UnauthorizedErr("Invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			hooks.SendResponse(w, http.StatusUnauthorized, nil, wraperrors.UnauthorizedErr("Invalid token claims"))
			return
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			hooks.SendResponse(w, http.StatusUnauthorized, nil, wraperrors.UnauthorizedErr("Invalid user ID in tokne"))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIdFromContext(ctx context.Context) (string, error) {
	id, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", wraperrors.UnauthorizedErr("Invalid user ID in token")
	}
	return id, nil
}
