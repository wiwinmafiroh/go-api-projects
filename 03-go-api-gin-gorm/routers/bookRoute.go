package routers

import (
	"03-go-api-gin-gorm/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	route := gin.Default()

	route.GET("/books", controllers.GetAllBooks)
	route.POST("/books", controllers.CreateBook)
	route.GET("/books/:bookID", controllers.GetBookByID)
	route.PUT("/books/:bookID", controllers.UpdateBookByID)
	route.DELETE("/books/:bookID", controllers.DeleteBookByID)

	return route
}
