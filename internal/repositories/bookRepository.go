package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/Melliv/Book-Market-Server/internal/helpers"
	"github.com/Melliv/Book-Market-Server/internal/types"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetBook(bookId string) types.Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book types.Book

	err := helpers.DB.Collection("Book").FindOne(ctx, bson.M{"_id": helpers.GetPrimitiveObjectID(bookId)}).Decode(&book)
	if err != nil {
		log.Fatal(err)
	}
	return book
}

func GetBooks(limit int, offset int) []types.Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := helpers.DB.Collection("Book").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var books []types.Book
	helpers.MapCursorToTarget(ctx, cursor, &books)

	return books
}

func GetBooksCount() int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	totalCount, err := helpers.DB.Collection("Book").CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	return int(totalCount)
}

func GetBooksByOwnerId(ownerId string) []types.Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"ownerId": ownerId}
	cursor, err := helpers.DB.Collection("Book").Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var books []types.Book
	helpers.MapCursorToTarget(ctx, cursor, &books)
	return books
}

func CreateBook(book types.Book) types.Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := helpers.DB.Collection("Book").InsertOne(ctx, book)
	if err != nil {
		log.Fatal(err)
	}
	book.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return book
}

func CreateBooks(books []types.Book) []types.Book {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bookInterfaces []interface{}
	for _, book := range books {
		bookInterfaces = append(bookInterfaces, book)
	}

	result, err := helpers.DB.Collection("Book").InsertMany(ctx, bookInterfaces)
	if err != nil {
		log.Fatal(err)
	}

	for i, insertedID := range result.InsertedIDs {
		books[i].ID = insertedID.(primitive.ObjectID).Hex()
	}

	return books
}

func UpdateBook(book types.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": helpers.GetPrimitiveObjectID(book.ID)}
	update := bson.M{
		"$set": bson.M{
			"title":   book.Title,
			"author":  book.Author,
			"ownerId": book.OwnerId,
		},
	}
	result, err := helpers.DB.Collection("Book").UpdateOne(ctx, filter, update)

	if result != nil && result.MatchedCount != 1 {
		return fmt.Errorf("MatchedCount not 1")
	}
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func DeleteBook(bookId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": helpers.GetPrimitiveObjectID(bookId)}
	result, err := helpers.DB.Collection("Book").DeleteOne(ctx, filter)

	if result != nil && result.DeletedCount != 1 {
		return fmt.Errorf("MatchedCount not 1")
	}
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
