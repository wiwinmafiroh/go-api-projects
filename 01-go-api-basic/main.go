package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"desc"`
}

var Books = []Book{}

var PORT = ":4000"

func main() {
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllBooks(w, r)

		case "POST":
			addBook(w, r)

		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getBookByID(w, r)

		case "PUT":
			updateBookByID(w, r)

		case "DELETE":
			deleteBookByID(w, r)

		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Application is Listening on PORT", PORT)
	http.ListenAndServe(PORT, nil)
}

func getAllBooks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	booksJSON, err := json.Marshal(Books)
	if err != nil {
		http.Error(w, "Failed to Encode Book to JSON", http.StatusInternalServerError)
		return
	}

	w.Write(booksJSON)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to Read Request Body", http.StatusBadRequest)
		return
	}

	var requestBookData Book
	err = json.Unmarshal(requestBody, &requestBookData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if requestBookData.Title == "" {
		http.Error(w, "Title Cannot be Empty", http.StatusBadRequest)
		return
	}

	newBookID := 1
	if len(Books) > 0 {
		newBookID = Books[len(Books)-1].ID + 1
	}

	newBook := Book{
		ID:          newBookID,
		Title:       requestBookData.Title,
		Author:      requestBookData.Author,
		Description: requestBookData.Description,
	}

	Books = append(Books, newBook)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created"))
}

func getBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path

	segments := strings.Split(path, "/")

	bookID, err := strconv.Atoi(segments[2])
	if err != nil {
		http.Error(w, "Invalid Parameter in Path", http.StatusBadRequest)
		return
	}

	for _, book := range Books {
		if book.ID == bookID {
			bookJSON, err := json.Marshal(book)
			if err != nil {
				http.Error(w, "Failed to Encode Book to JSON", http.StatusInternalServerError)
				return
			}

			w.Write(bookJSON)
			return
		}
	}

	errMessage := fmt.Sprintf("Book with ID %d Doesn't Exist", bookID)
	http.Error(w, errMessage, http.StatusNotFound)
}

func updateBookByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	segments := strings.Split(path, "/")

	bookID, err := strconv.Atoi(segments[2])
	if err != nil {
		http.Error(w, "Invalid Parameter in Path", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to Read Request Body", http.StatusBadRequest)
		return
	}

	var requestBookData Book
	err = json.Unmarshal(requestBody, &requestBookData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for index, book := range Books {
		if bookID == book.ID {
			if requestBookData.Title != "" {
				Books[index].Title = requestBookData.Title
			}

			if requestBookData.Author != "" {
				Books[index].Author = requestBookData.Author
			}

			if requestBookData.Description != "" {
				Books[index].Description = requestBookData.Description
			}

			Books[index].ID = bookID

			w.Write([]byte("Updated"))
			return
		}
	}

	errMessage := fmt.Sprintf("Book with ID %d Doesn't Exist", bookID)
	http.Error(w, errMessage, http.StatusNotFound)
}

func deleteBookByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	segments := strings.Split(path, "/")

	bookID, err := strconv.Atoi(segments[2])
	if err != nil {
		http.Error(w, "Invalid Parameter in Path", http.StatusBadRequest)
		return
	}

	for index, book := range Books {
		if bookID == book.ID {
			copy(Books[index:], Books[index+1:])
			Books = Books[:len(Books)-1]

			w.Write([]byte("Deleted"))
			return
		}
	}

	errMessage := fmt.Sprintf("Book with ID %d Doesn't Exist", bookID)
	http.Error(w, errMessage, http.StatusNotFound)
}
