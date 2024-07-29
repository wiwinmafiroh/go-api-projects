package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func init() {
	handleDatabaseConnection()
	handleRequiredTables()
}

func handleDatabaseConnection() {
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
}

func handleRequiredTables() {
	usersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			created_at timestamptz DEFAULT now(),
			updated_at timestamptz DEFAULT now()
		);
	`

	productsTable := `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(150) NOT NULL,
			description TEXT,
			price NUMERIC NOT NULL,
			image_url TEXT NOT NULL,
			user_id INT NOT NULL,
			created_at timestamptz DEFAULT now(),
			updated_at timestamptz DEFAULT now(),
			CONSTRAINT fk_products_user_id
				FOREIGN KEY(user_id)
				REFERENCES users(id)
				ON UPDATE CASCADE
				ON DELETE CASCADE
		);
	`

	createTableQueries := fmt.Sprintf("%s %s", usersTable, productsTable)

	_, err = db.Exec(createTableQueries)
	if err != nil {
		log.Panicln("Failed to create required tables:", err)
	}
}

func GetDatabaseInstance() *sql.DB {
	return db
}
