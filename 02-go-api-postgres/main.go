package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

const PORT = ":4000"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"desc"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"desc"`
}

func init() {
	err = godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err = sql.Open(os.Getenv("DB_DIALECT"), psqlInfo)
	if err != nil {
		log.Panicln("Failed to validate the database configuration:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panicln("Failed to connect to the database:", err)
	}

	log.Println("Successfully connected to the database")

	createBooksTableQuery := `
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			author VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			created_at timestamptz DEFAULT now(),
			updated_at timestamptz DEFAULT now()
		);
	`

	_, err = db.Exec(createBooksTableQuery)
	if err != nil {
		log.Panicln("Failed to create 'books' table:", err)
	}
}

func main() {
	defer db.Close()

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllBooks(w, r)

		case "POST":
			addBook(w, r)

		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Listening on port", PORT)
	http.ListenAndServe(PORT, nil)
}

func getAllBooks(w http.ResponseWriter, _ *http.Request) {
	var books = []Book{}

	getBooksQuery := `
		SELECT * FROM books
		ORDER BY id;
	`

	rows, err := db.Query(getBooksQuery)
	if err != nil {
		http.Error(w, "Failed to execute the database query", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var book Book

		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to scan rows for book data", http.StatusInternalServerError)
			return
		}

		books = append(books, book)
	}

	booksJSON, err := json.Marshal(books)
	if err != nil {
		http.Error(w, "Failed to encode books data to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(booksJSON)
}

func getBookByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	segments := strings.Split(path, "/")

	bookID, err := strconv.Atoi(segments[2])
	if err != nil {
		http.Error(w, "Invalid parameter in path", http.StatusBadRequest)
		return
	}

	getBookByIDQuery := `
		SELECT * FROM books
		WHERE id = $1;
	`

	var book Book

	row := db.QueryRow(getBookByIDQuery, bookID)

	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			message := fmt.Sprintf("Book with ID %d doesn't exist", bookID)

			http.Error(w, message, http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to retrieve book details", http.StatusInternalServerError)
		return
	}

	bookJSON, err := json.Marshal(&book)
	if err != nil {
		http.Error(w, "Failed to encode book data to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bookJSON)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported content-type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the request body", http.StatusBadRequest)
		return
	}

	var bookRequest BookRequest
	err = json.Unmarshal(body, &bookRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if bookRequest.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	addBookQuery := `
		INSERT INTO books (title, author, description)
		VALUES ($1, $2, $3);
	`

	res, err := db.Exec(addBookQuery, bookRequest.Title, bookRequest.Author, bookRequest.Description)
	if err != nil {
		http.Error(w, "Failed to execute the database query", http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to get the number of affected rows", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "Failed to create a book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created"))
}

func updateBookByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	segments := strings.Split(path, "/")

	bookID, err := strconv.Atoi(segments[2])
	if err != nil {
		http.Error(w, "Invalid parameter in path", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported content-type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read the request body", http.StatusBadRequest)
		return
	}

	var bookRequest BookRequest
	err = json.Unmarshal(body, &bookRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if bookRequest.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	updateBookQuery := `
		UPDATE books
		SET title = $2, author = $3, description = $4, updated_at = $5
		WHERE id = $1;
	`

	res, err := db.Exec(updateBookQuery, bookID, bookRequest.Title, bookRequest.Author, bookRequest.Description, time.Now())
	if err != nil {
		http.Error(w, "Failed to execute the database query", http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to get the number of affected rows", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		message := fmt.Sprintf("Book with ID %d doesn't exist", bookID)

		http.Error(w, message, http.StatusNotFound)
		return
	}

	w.Write([]byte("Updated"))
}

func deleteBookByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	segments := strings.Split(path, "/")

	bookID, err := strconv.Atoi(segments[2])
	if err != nil {
		http.Error(w, "Invalid parameter in path", http.StatusBadRequest)
		return
	}

	deleteBookQuery := `
		DELETE FROM books
		WHERE id = $1;
	`

	res, err := db.Exec(deleteBookQuery, bookID)
	if err != nil {
		http.Error(w, "Failed to execute the database query", http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to get the number of affected rows", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		message := fmt.Sprintf("Book with ID %d doesn't exist", bookID)

		http.Error(w, message, http.StatusNotFound)
		return
	}

	w.Write([]byte("Deleted"))
}
