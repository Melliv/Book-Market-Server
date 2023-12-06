package main

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Book struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	Title   string `json:"title" bson:"title"`
	Author  string `json:"author" bson:"author"`
	OwnerId string `json:"ownerId" bson:"ownerId"`
}

type BookPagination struct {
	Books      []Book
	TotalCount int    `json:"totalCount" bson:"totalCount"`
}
