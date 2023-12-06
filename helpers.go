package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	value := os.Getenv(key)

	return value
}

func getEnv() Env {
	value := getEnvVar("ENV")
	switch value {
	case "PROD":
		return Prod
	case "DEV":
		return Dev
	case "TEST":
		return Test
	default:
		log.Fatal("Env not exist: ", value)
	}
	return Dev
}

func baseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func mapCursorToTarget(ctx context.Context, cursor *mongo.Cursor, target interface{}) error {
	if err := cursor.All(ctx, target); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func getPrimitiveObjectID(value string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		panic(err)
	}
	return objID
}

func initData() {
	err := db.Collection("Book").Drop(context.Background())
	if err != nil {
		log.Fatal("Database collection drop error! ", err)
	}

	var user User
	err = db.Collection("User").FindOne(context.Background(), bson.D{}).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	bookRepository := NewDefaultBookRepository()

	var books []Book
	for i := 0; i < 20; i++ {
		book := Book{
			Title:   "Book " + strconv.Itoa(i),
			Author:  "Author " + strconv.Itoa(i),
			OwnerId: user.ID,
		}
		books = append(books, book)
	}
	books = bookRepository.createBooks(books)

	fmt.Println("Initial data initialised!", len(books), books[0].ID)
}
