package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Melliv/Book-Market-Server/internal/controllers"
	"github.com/Melliv/Book-Market-Server/internal/helpers"
	"github.com/Melliv/Book-Market-Server/internal/init"
)

func main() {
	r := mux.NewRouter()

	helpers.ConnectDatabase(helpers.GetEnv())
	initi.InitData()
	r.Use(helpers.BaseMiddleware)

	r.HandleFunc("/login", controllers.LoginHandler).Methods("POST")

	r.HandleFunc("/books", controllers.GetBook).Methods("GET")
	r.Handle("/user-books", helpers.Authenticate(http.HandlerFunc(controllers.GetUserBooks))).Methods("GET")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", helpers.Authenticate(http.HandlerFunc(controllers.GetBook))).Methods("GET")
	r.Handle("/book", helpers.Authenticate(http.HandlerFunc(controllers.CreateBook))).Methods("POST")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", helpers.Authenticate(http.HandlerFunc(controllers.UpdateBook))).Methods("PUT")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", helpers.Authenticate(http.HandlerFunc(controllers.DeleteBook))).Methods("DELETE")

	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
