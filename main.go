package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"strconv"
)

// Book struct
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice book struct
var books []Book

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get a single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //Get params
	// loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		} else {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) //Mock id
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r) // Get params

    for index, item := range books {
        if item.ID == params["id"] {
            // Decode the partial update from the request body
            var updates map[string]interface{}
            if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
                http.Error(w, "Invalid request payload", http.StatusBadRequest)
                return
            }

            // Create a copy of the existing book for updates
            updatedBook := item

            // Update fields in the copy of the book
            if title, ok := updates["title"].(string); ok {
                updatedBook.Title = title
            }
            if isbn, ok := updates["isbn"].(string); ok {
                updatedBook.Isbn = isbn
            }
            if author, ok := updates["author"].(map[string]interface{}); ok {
                if firstname, ok := author["firstname"].(string); ok {
                    updatedBook.Author.Firstname = firstname
                }
                if lastname, ok := author["lastname"].(string); ok {
                    updatedBook.Author.Lastname = lastname
                }
            }

            // Replace the book in the slice with the updated book
            books[index] = updatedBook

            // Return the updated book
            json.NewEncoder(w).Encode(updatedBook)
            return
        }
    }

    // If no book was found with the given ID, return a 404 error
    http.Error(w, "Book not found", http.StatusNotFound)
}


// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //Get params
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// Init router
	r := mux.NewRouter()

	books = append(books, Book{ID: "1", Isbn: "44454", Title: "Book One", Author: &Author{Firstname: "Adewole", Lastname: "Fidelis"}})
	books = append(books, Book{ID: "2", Isbn: "43455", Title: "Book Two", Author: &Author{Firstname: "Adedayo", Lastname: "Adewoye"}})

	// Route Handlers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Start server
	http.ListenAndServe(":8000", r)
}
