package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Melliv/Book-Market-Server/internal/helpers"
	"github.com/Melliv/Book-Market-Server/internal/services"
	"github.com/Melliv/Book-Market-Server/internal/types"
	"github.com/gorilla/mux"
)

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	books := services.GetBook(bookId)
	json.NewEncoder(w).Encode(books)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
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

	booksPagination, err := services.GetBooks(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(*booksPagination)
}

func GetUserBooks(w http.ResponseWriter, r *http.Request) {
	ownerId := r.Context().Value("userId").(string)
	books := services.GetBooksByOwnerId(ownerId)
	json.NewEncoder(w).Encode(books)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	ownerId := r.Context().Value("userId").(string)
	var book types.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request body!", http.StatusBadRequest)
		return
	}

	err = helpers.ValidateBody(book)
	if err != nil {
		http.Error(w, "Invalid request body!\n" + err.Error(), http.StatusBadRequest)
		return
	}

	book.OwnerId = ownerId
	book = services.CreateBook(book)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]
	var book types.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil || bookId != book.ID {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = services.UpdateBook(book)
	if err != nil {
		log.Print(err)
		http.Error(w, "ID mismatching", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["id"]

	err := services.DeleteBook(bookId)
	if err != nil {
		log.Print(err)
		http.Error(w, "ID mismatching", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w)
}

