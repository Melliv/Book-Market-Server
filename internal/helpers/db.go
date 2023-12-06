package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Melliv/Book-Market-Server/internal/enums"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database


func ConnectDatabase(env enums.Env) {
	// MongoDB connection string. Replace it with your actual MongoDB connection string.
	connectionString := GetEnvVar("DB_HOST")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	if env == enums.Prod {
		DB = client.Database("Some")
	} else if env == enums.Test {
		DB = client.Database("Test")
	} else {
		DB = client.Database("Some")
	}

	fmt.Println("Connected to MongoDB!")
}

func DropCollection(collectionName string) {
	err := DB.Collection(collectionName).Drop(context.Background())
	if err != nil {
		log.Fatal("Database collection drop error! ", err)
	}
}
