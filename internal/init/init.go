package initi

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Melliv/Book-Market-Server/internal/helpers"
	"github.com/Melliv/Book-Market-Server/internal/repositories"
	"github.com/Melliv/Book-Market-Server/internal/types"
	"go.mongodb.org/mongo-driver/bson"
)

func InitData() {
	err := helpers.DB.Collection("Book").Drop(context.Background())
	if err != nil {
		log.Fatal("Database collection drop error! ", err)
	}

	var user types.User
	err = helpers.DB.Collection("User").FindOne(context.Background(), bson.D{}).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	var books []types.Book
	for i := 0; i < 20; i++ {
		book := types.Book{
			Title:   "Book " + strconv.Itoa(i),
			Author:  "Author " + strconv.Itoa(i),
			OwnerId: user.ID,
		}
		books = append(books, book)
	}
	books = repositories.CreateBooks(books)

	fmt.Println("Initial data initialised!", len(books), books[0].ID)
}
