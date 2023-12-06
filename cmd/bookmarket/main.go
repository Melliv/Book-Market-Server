package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Melliv/Book-Market-Server/internal/controllers"
	"github.com/Melliv/Book-Market-Server/internal/helpers"
	"github.com/Melliv/Book-Market-Server/internal/init"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	env := helpers.GetEnv()
	port, _ := strconv.Atoi(helpers.GetEnvVar("PORT"))
	fmt.Println("Env:", env)

	helpers.ConnectDatabase(env)
	initi.InitData()
	r.Use(helpers.BaseMiddleware)

	r.HandleFunc("/login", controllers.LoginHandler).Methods("POST")

	r.HandleFunc("/books", controllers.GetBooks).Methods("GET")
	r.Handle("/user-books", helpers.Authenticate(http.HandlerFunc(controllers.GetUserBooks))).Methods("GET")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", helpers.Authenticate(http.HandlerFunc(controllers.GetBook))).Methods("GET")
	r.Handle("/book", helpers.Authenticate(http.HandlerFunc(controllers.CreateBook))).Methods("POST")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", helpers.Authenticate(http.HandlerFunc(controllers.UpdateBook))).Methods("PUT")
	r.Handle("/book/{id:[0-9a-fA-F]{24}}", helpers.Authenticate(http.HandlerFunc(controllers.DeleteBook))).Methods("DELETE")

	fmt.Printf("Server is running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
