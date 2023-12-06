package main

import "fmt"

type BookService interface {
	getBook(bookId string) Book
	getBooks(limit int, offset int) (*BookPagination, error)
	getBooksByOwnerId(ownerId string) []Book
	createBook(book Book) Book
	updateBook(book Book) error
	deleteBook(bookId string) error
}

type DefaultBookService struct {
	bookRepository BookRepository
}

func NewDefaultBookService() *DefaultBookService {
	return &DefaultBookService{bookRepository: NewDefaultBookRepository()}
}

func (s *DefaultBookService) getBook(bookId string) Book {
	return s.bookRepository.getBook(bookId)
}

func (s *DefaultBookService) getBooks(limit int, offset int) (*BookPagination, error) {
	if limit < 0 {
		return nil, fmt.Errorf("Limit must be >= 0")
	}
	if offset < 0 {
		return nil, fmt.Errorf("Offset must be >= 0")
	}
	bookPagination := BookPagination{
		Books:      s.bookRepository.getBooks(limit, offset),
		TotalCount: s.bookRepository.getBooksCount(),
	}

	return &bookPagination, nil
}

func (s *DefaultBookService) getBooksByOwnerId(ownerId string) []Book {
	return s.bookRepository.getBooksByOwnerId(ownerId)
}

func (s *DefaultBookService) createBook(book Book) Book {
	return s.bookRepository.createBook(book)
}

func (s *DefaultBookService) updateBook(book Book) error {
	return s.bookRepository.updateBook(book)
}

func (s *DefaultBookService) deleteBook(bookId string) error {
	return s.bookRepository.deleteBook(bookId)
}
