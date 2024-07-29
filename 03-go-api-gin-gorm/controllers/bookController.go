package controllers

import (
	"03-go-api-gin-gorm/database"
	"03-go-api-gin-gorm/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookRequest struct {
	NameBook string `json:"name_book"`
	Author   string `json:"author"`
}

func GetAllBooks(ctx *gin.Context) {
	var books []models.Book

	err := database.DB.Order("id asc").Find(&books).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve books from the database",
		})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func CreateBook(ctx *gin.Context) {
	if ctx.GetHeader("Content-Type") != "application/json" {
		ctx.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
			"error": "Unsupported content-type",
		})
		return
	}

	var bookRequest BookRequest

	if err := ctx.ShouldBindJSON(&bookRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	if bookRequest.NameBook == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Book name cannot be empty",
		})
		return
	}

	newBook := models.Book{
		NameBook: bookRequest.NameBook,
		Author:   bookRequest.Author,
	}

	err := database.DB.Create(&newBook).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create a book",
		})
		return
	}

	ctx.JSON(http.StatusCreated, newBook)
}

func GetBookByID(ctx *gin.Context) {
	bookID := ctx.Param("bookID")

	parsedBookID, err := strconv.Atoi(bookID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID in path",
		})
		return
	}

	var book models.Book

	err = database.DB.First(&book, parsedBookID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg := fmt.Sprintf("Book with ID %d doesn't exist", parsedBookID)

			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": msg,
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve book details",
		})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func UpdateBookByID(ctx *gin.Context) {
	bookID := ctx.Param("bookID")

	parsedBookID, err := strconv.Atoi(bookID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID in path",
		})
		return
	}

	if ctx.GetHeader("Content-Type") != "application/json" {
		ctx.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
			"error": "Unsupported content-type",
		})
		return
	}

	var existingBook models.Book

	err = database.DB.First(&existingBook, parsedBookID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg := fmt.Sprintf("Book with ID %d doesn't exist", parsedBookID)

			ctx.JSON(http.StatusNotFound, gin.H{
				"error": msg,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve book details",
		})
		return
	}

	var bookRequest BookRequest

	if err = ctx.ShouldBindJSON(&bookRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	if bookRequest.NameBook == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Book name cannot be empty",
		})
		return
	}

	updateBook := models.Book{
		NameBook: bookRequest.NameBook,
		Author:   bookRequest.Author,
	}

	err = database.DB.Model(&existingBook).Updates(updateBook).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update a book",
		})
		return
	}

	ctx.JSON(http.StatusOK, existingBook)
}

func DeleteBookByID(ctx *gin.Context) {
	bookID := ctx.Param("bookID")

	parsedBookID, err := strconv.Atoi(bookID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID in path",
		})
		return
	}

	var existingBook models.Book

	err = database.DB.First(&existingBook, parsedBookID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			msg := fmt.Sprintf("Book with ID %d doesn't exist", parsedBookID)

			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": msg,
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve book details",
		})
		return
	}

	err = database.DB.Delete(&existingBook).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete a book",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})
}
