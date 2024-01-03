package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Melliv/Book-Market-Server/internal/helpers"
	"github.com/Melliv/Book-Market-Server/internal/types"
	"github.com/Melliv/Book-Market-Server/internal/repositories"
	"github.com/dgrijalva/jwt-go"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userPointer := repositories.GetUserByUsernameAndPassword(user.Username, user.Password)

	if userPointer != nil {
		user = *userPointer
		expirationTime := helpers.GetJwtExpirationTime()
		claims := &types.JWTClaims{
			Username: user.Username,
			StandardClaims: jwt.StandardClaims{
				Id: user.ID,
				ExpiresAt: expirationTime.Unix(),
			},
		}


		jwtSecretKey := helpers.GetEnvVar("JWT_SECRET_KEY")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecretKey))

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": tokenString, "id": userPointer.ID, "username": userPointer.Username})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}
