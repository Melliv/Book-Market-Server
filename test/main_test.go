package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Melliv/Book-Market-Server/internal/controllers"
	"github.com/Melliv/Book-Market-Server/internal/enums"
	"github.com/Melliv/Book-Market-Server/internal/helpers"
	"github.com/Melliv/Book-Market-Server/internal/types"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	helpers.ConnectDatabase(enums.Test)
	helpers.DropCollection("Book")

	token, err := logIn()
	assert.NoError(t, err)

	bookId, err := createBook(token)
	assert.NoError(t, err)

	err = getBook(token, bookId)
	assert.NoError(t, err)

	err = updateBook(token, bookId)
	assert.NoError(t, err)

	err = updateBook(token, bookId)
	assert.NoError(t, err)

	err = getBooks()
	assert.NoError(t, err)
}

func logIn() (string, error) {
	validCredentials := types.User{
		Username: "admin1",
		Password: "admin1",
	}
	validCredentialsJSON, err := json.Marshal(validCredentials)
	if err != nil {
		return "", err
	}

	reqValid := httptest.NewRequest("POST", "/login", bytes.NewBuffer(validCredentialsJSON))
	reqValid.Header.Set("Content-Type", "application/json")


	// Create a new recorder for the response
	w := httptest.NewRecorder()

	// Call the loginHandler function with the valid credentials
	controllers.LoginHandler(w, reqValid)

	// Check the status code and the presence of the token in the response body
	if w.Code != http.StatusOK {
		return "", fmt.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return "", err
	}

	token, ok := response["token"]
	if !ok {
		return "", fmt.Errorf("Token not found in response body")
	}

	return token, nil
}

func createBook(token string) (string, error) {
	reqBook := types.Book{
		Title: "Test title 1",
		Author: "Test author 1",
	}
	requestByte, _ := json.Marshal(reqBook)
	requestReader := bytes.NewReader(requestByte)

	req := httptest.NewRequest("POST", "/book", requestReader)
	req.Header.Set("Authorization", token)


	w := httptest.NewRecorder()

	authMiddleware := helpers.BaseMiddleware(helpers.Authenticate(http.HandlerFunc(controllers.CreateBook)))

	authMiddleware.ServeHTTP(w, req)


	if w.Code != http.StatusOK {
		return "", fmt.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resBook types.Book
	err := json.NewDecoder(w.Body).Decode(&resBook)
	if err != nil {
		return "", fmt.Errorf("Invalid request body. %s", err)
	}
	if resBook.ID == "" {
		return "", fmt.Errorf("BookId missing")
	}

	if reqBook.Title != resBook.Title {
		return "", fmt.Errorf("Expected book title %s, got %s", reqBook.Title, resBook.Title)
	}

	return resBook.ID, nil
}

func updateBook(token string, bookId string) (error) {
	book := types.Book{
		ID: bookId,
		Title: "Test title 1 new",
		Author: "Test author 1 new",
	}
	requestByte, _ := json.Marshal(book)
	requestReader := bytes.NewReader(requestByte)

    vars := map[string]string{
        "id": bookId,
    }

	req := httptest.NewRequest("PUT", "/book/" + bookId, requestReader)
	req.Header.Set("Authorization", token)
	req = mux.SetURLVars(req, vars)


	w := httptest.NewRecorder()

	authMiddleware := helpers.BaseMiddleware(helpers.Authenticate(http.HandlerFunc(controllers.UpdateBook)))

	authMiddleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body)
	}

	return nil
}

func getBook(token string, bookId string) (error) {
	book := types.Book{
		ID: bookId,
		Title: "Test title 1",
		Author: "Test author 1",
	}

    vars := map[string]string{
        "id": bookId,
    }

	req := httptest.NewRequest("GET", "/book/" + bookId, nil)
	req.Header.Set("Authorization", token)
	req = mux.SetURLVars(req, vars)


	w := httptest.NewRecorder()

	authMiddleware := helpers.BaseMiddleware(helpers.Authenticate(http.HandlerFunc(controllers.GetBook)))

	authMiddleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body)
	}

	var resBook types.Book
	err := json.NewDecoder(w.Body).Decode(&resBook)
	if err != nil {
		return fmt.Errorf("Invalid request body. %s", err)
	}
	if book.Title != resBook.Title || book.Author != resBook.Author {
		return fmt.Errorf("Book values are not right.\n Title: Expected %s, was %s \n Author: Expected %s, was %s", book.Title, resBook.Title, book.Author, resBook.Author)
	}

	return nil
}

func getBooks() (error) {
	req := httptest.NewRequest("GET", "/books?limit=1&offset=0", nil)

	w := httptest.NewRecorder()

	middleware := helpers.BaseMiddleware(http.HandlerFunc(controllers.GetBooks))

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body)
	}

	var resBookPagination types.BookPagination

	err := json.NewDecoder(w.Body).Decode(&resBookPagination)
	if err != nil {
		return fmt.Errorf("Invalid request body. %s", err)
	}
	if resBookPagination.TotalCount != 1 {
		return fmt.Errorf("Expected books list lenght %d, but was %d", 1, resBookPagination.TotalCount)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/books?limit=-1", nil)
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		return fmt.Errorf("Expected status code %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/books?offset=-1", nil)
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		return fmt.Errorf("Expected status code %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body)
	}


	return nil
}
