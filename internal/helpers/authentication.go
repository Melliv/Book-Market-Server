package helpers

import (
	"context"
	"net/http"
	"time"

	"github.com/Melliv/Book-Market-Server/internal/types"
	"github.com/dgrijalva/jwt-go"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		jwtSecretKey := GetEnvVar("JWT_SECRET_KEY")

		token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*types.JWTClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "username", claims.Username)
			ctx = context.WithValue(ctx, "userId", claims.Id)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func GetJwtExpirationTime() time.Time {
	env := GetEnvVar("ENV")
	if env == "DEV" {
		return time.Unix(1<<62-1, 0)
	}
	return time.Now().Add(5 * time.Minute)
}
