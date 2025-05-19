package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jhonVitor-rs/url-shortener/pkg/utils"
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
			utils.WriteJSON(w, http.StatusUnauthorized, utils.ErrorResponse{
				Message: "Missing or invalid token",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.ErrorResponse{
				Message: "Invalid token", Error: &err,
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.ErrorResponse{
				Message: "Invalid token claims",
			})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.ErrorResponse{
				Message: "Invalid user ID in tokne",
			})
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIdFromContex(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
