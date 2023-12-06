package main

import (
	"fmt"
	"log"
	"net/http"
	"controller"

	"github.com/gorilla/mux"
)

var bookController BookController

func main() {
	r := mux.NewRouter()

	connectDatabase(getEnv())
	initData()
	r.Use(baseMiddleware)

	bookController = NewDefaultBookController()

	r.HandleFunc("/login", loginHandler).Methods("POST")

	r.HandleFunc("/books", bookController.getBooks).Methods("GET")
	r.Handle("/user-books", authenticate(http.HandlerFunc(bookController.getUserBooks))).Methods("GET")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", authenticate(http.HandlerFunc(bookController.getBook))).Methods("GET")
	r.Handle("/book", authenticate(http.HandlerFunc(bookController.createBook))).Methods("POST")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", authenticate(http.HandlerFunc(bookController.updateBook))).Methods("PUT")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", authenticate(http.HandlerFunc(bookController.deleteBook))).Methods("DELETE")

	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
