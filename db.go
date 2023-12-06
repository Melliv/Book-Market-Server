package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database


func connectDatabase(env Env) {
	// MongoDB connection string. Replace it with your actual MongoDB connection string.
	connectionString := getEnvVar("DB_HOST")

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

	if env == Prod {
		db = client.Database("Some")
	} else if env == Test {
		db = client.Database("Test")
	} else {
		db = client.Database("Some")
	}

	fmt.Println("Connected to MongoDB!")
}

func dropCollection(collectionName string) {
	err := db.Collection(collectionName).Drop(context.Background())
	if err != nil {
		log.Fatal("Database collection drop error! ", err)
	}
}