package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookRepository interface {
	getBook(bookId string) Book
	getBooks(limit int, offset int) []Book
	getBooksCount() int
	getBooksByOwnerId(ownerId string) []Book
	createBook(book Book) Book
	createBooks(books []Book) []Book
	updateBook(book Book) error
	deleteBook(bookId string) error
}

type DefaultBookRepository struct {
}

func NewDefaultBookRepository() *DefaultBookRepository {
	return &DefaultBookRepository{}
}

func (r *DefaultBookRepository) getBook(bookId string) Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book Book

	err := db.Collection("Book").FindOne(ctx, bson.M{"_id": getPrimitiveObjectID(bookId)}).Decode(&book)
	if err != nil {
		log.Fatal(err)
	}
	return book
}

func (r *DefaultBookRepository) getBooks(limit int, offset int) []Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := db.Collection("Book").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var books []Book
	mapCursorToTarget(ctx, cursor, &books)

	return books
}

func (r *DefaultBookRepository) getBooksCount() int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	totalCount, err := db.Collection("Book").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	return int(totalCount)
}

func (r *DefaultBookRepository) getBooksByOwnerId(ownerId string) []Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"ownerId": ownerId}
	cursor, err := db.Collection("Book").Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var books []Book
	mapCursorToTarget(ctx, cursor, &books)
	return books
}

func (r *DefaultBookRepository) createBook(book Book) Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := db.Collection("Book").InsertOne(ctx, book)
	if err != nil {
		log.Fatal(err)
	}
	book.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return book
}

func (r *DefaultBookRepository) createBooks(books []Book) []Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bookInterfaces []interface{}
	for _, book := range books {
		bookInterfaces = append(bookInterfaces, book)
	}

	result, err := db.Collection("Book").InsertMany(ctx, bookInterfaces)
	if err != nil {
		log.Fatal(err)
	}

	for i, insertedID := range result.InsertedIDs {
		books[i].ID = insertedID.(primitive.ObjectID).Hex()
	}

	return books
}

func (r *DefaultBookRepository) updateBook(book Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": getPrimitiveObjectID(book.ID)}
	update := bson.M{
		"$set": bson.M{
			"title":   book.Title,
			"author":  book.Author,
			"ownerId": book.OwnerId,
		},
	}
	result, err := db.Collection("Book").UpdateOne(ctx, filter, update)

	if result != nil && result.MatchedCount != 1 {
		return fmt.Errorf("MatchedCount not 1")
	}
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (r *DefaultBookRepository) deleteBook(bookId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": getPrimitiveObjectID(bookId)}
	result, err := db.Collection("Book").DeleteOne(ctx, filter)

	if result != nil && result.DeletedCount != 1 {
		return fmt.Errorf("MatchedCount not 1")
	}
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
