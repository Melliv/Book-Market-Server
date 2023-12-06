package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BookController interface {
	getBook(w http.ResponseWriter, r *http.Request)
	getBooks(w http.ResponseWriter, r *http.Request)
	getUserBooks(w http.ResponseWriter, r *http.Request)
	createBook(w http.ResponseWriter, r *http.Request)
	updateBook(w http.ResponseWriter, r *http.Request)
	deleteBook(w http.ResponseWriter, r *http.Request)
}

type DefaultBookController struct {
	bookService BookService
}

func NewDefaultBookController() *DefaultBookController {
	return &DefaultBookController{bookService: NewDefaultBookService()}
}


func (c *DefaultBookController) getBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	books := c.bookService.getBook(bookId)
	json.NewEncoder(w).Encode(books)
}

func (c *DefaultBookController) getBooks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	var limit int
	var offset int

	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	booksPagination, err := c.bookService.getBooks(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(*booksPagination)
}

func (c *DefaultBookController) getUserBooks(w http.ResponseWriter, r *http.Request) {
	ownerId := r.Context().Value("userId").(string)
	books := c.bookService.getBooksByOwnerId(ownerId)
	json.NewEncoder(w).Encode(books)
}

func (c *DefaultBookController) createBook(w http.ResponseWriter, r *http.Request) {
	ownerId := r.Context().Value("userId").(string)
	var book Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	book.OwnerId = ownerId
	book = c.bookService.createBook(book)
	json.NewEncoder(w).Encode(book)
}

func (c *DefaultBookController) updateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]
	var book Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil || bookId != book.ID {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = c.bookService.updateBook(book)
	if err != nil {
		log.Print(err)
		http.Error(w, "ID mismatching", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w)
}

func (c *DefaultBookController) deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	err := c.bookService.deleteBook(bookId)
	if err != nil {
		log.Print(err)
		http.Error(w, "ID mismatching", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w)
}

