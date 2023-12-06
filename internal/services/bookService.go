package services

import (
	"fmt"

	"github.com/Melliv/Book-Market-Server/internal/types"
	"github.com/Melliv/Book-Market-Server/internal/repositories"
)

func GetBook(bookId string) types.Book {
	return repositories.GetBook(bookId)
}

func GetBooks(limit int, offset int) (*types.BookPagination, error) {
	if limit < 0 {
		return nil, fmt.Errorf("Limit must be >= 0")
	}
	if offset < 0 {
		return nil, fmt.Errorf("Offset must be >= 0")
	}
	bookPagination := types.BookPagination{
		Books:      repositories.GetBooks(limit, offset),
		TotalCount: repositories.GetBooksCount(),
	}

	return &bookPagination, nil
}

func GetBooksByOwnerId(ownerId string) []types.Book {
	return repositories.GetBooksByOwnerId(ownerId)
}

func CreateBook(book types.Book) types.Book {
	return repositories.CreateBook(book)
}

func UpdateBook(book types.Book) error {
	return repositories.UpdateBook(book)
}

func DeleteBook(bookId string) error {
	return repositories.DeleteBook(bookId)
}
