package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		jwtSecretKey := getEnvVar("JWT_SECRET_KEY")

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "username", claims.Username)
			ctx = context.WithValue(ctx, "userId", claims.Id)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userRepository := NewUserRepository()
	userPointer := userRepository.getUserByUsernameAndPassword(user.Username, user.Password)

	if userPointer != nil {
		user = *userPointer
		expirationTime := getJwtExpirationTime()
		claims := &JWTClaims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				Id: user.ID,
				ExpiresAt: expirationTime.Unix(),
			},
		}


		jwtSecretKey := getEnvVar("JWT_SECRET_KEY")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecretKey))

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func getJwtExpirationTime() time.Time {
	env := getEnvVar("ENV")
	if env == "DEV" {
		return time.Unix(1<<62-1, 0)
	}
	return time.Now().Add(5 * time.Minute)
}
