package database

import (
	"03-go-api-gin-gorm/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func StartDB() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln("Failed to connect to the database:", err)
	}

	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		log.Panicln("Failed to perform auto migration:", err)
	}

	DB = db
	fmt.Println("Successfully connected to the database")
}
