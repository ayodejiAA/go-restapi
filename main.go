package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author of book
type Author struct {
	Name string `json:"name"`
}

// Collection of books
var books []Book

func getValue(initial, new string) string {
	if len(new) > 0 {
		return new
	}
	return initial
}

func getBooks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(books)
}

func getBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := mux.Vars(req)["id"]

	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(res).Encode(book)
			return
		}
	}

	json.NewEncoder(res).Encode(struct {
		Error string `json:"error"`
		Book  Book   `json:"book"`
	}{"book not found", Book{}})
}

func createBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(req.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000))
	books = append(books, book)
	json.NewEncoder(res).Encode(book)
}

func updateBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := mux.Vars(req)["id"]

	for index, book := range books {
		if book.ID == id {
			var body Book
			_ = json.NewDecoder(req.Body).Decode(&body)
			body.ID = id
			books[index] = body
			json.NewEncoder(res).Encode(body)
			return
		}
	}
}

func patchBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := mux.Vars(req)["id"]

	for index, book := range books {
		if book.ID == id {
			var body Book
			_ = json.NewDecoder(req.Body).Decode(&body)

			book.Title = getValue(book.Title, body.Title)
			if body.Author != nil {
				book.Author.Name = getValue(book.Author.Name, body.Author.Name)
			}

			books[index] = book

			json.NewEncoder(res).Encode(book)
			return
		}
	}
}

func deleteBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := mux.Vars(req)["id"]

	for index, book := range books {
		if book.ID == id {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}

	json.NewEncoder(res).Encode(books)
}

func main() {
	// Init router
	router := mux.NewRouter()

	// Mock Data todo: Implement DB
	books = append(books, Book{ID: "1", Title: "Book One", Author: &Author{Name: "Mike"}})
	books = append(books, Book{ID: "2", Title: "Book Two", Author: &Author{Name: "Steve"}})
	books = append(books, Book{ID: "3", Title: "Book Three", Author: &Author{Name: "Ali"}})

	// Route Handlers
	router.HandleFunc("/api/v1/books", getBooks).Methods("GET")
	router.HandleFunc("/api/v1/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/v1/books", createBook).Methods("POST")
	router.HandleFunc("/api/v1/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/v1/books/{id}", patchBook).Methods("PATCH")
	router.HandleFunc("/api/v1/books/{id}", deleteBook).Methods("DELETE")

	// Http listener
	fmt.Println("Listening on PORT 7000")
	log.Fatal(http.ListenAndServe(":7000", router))
}
