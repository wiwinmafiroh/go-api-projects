package handler

import (
	"06-go-api-with-unittest/database"
	"06-go-api-with-unittest/repository/product_repository/product_postgres"
	"06-go-api-with-unittest/repository/user_repository/user_postgres"
	"06-go-api-with-unittest/service"

	"github.com/gin-gonic/gin"
)

func StartApp() {
	const PORT = ":3000"

	db := database.GetDatabaseInstance()
	defer db.Close()

	userRepository := user_postgres.NewUserPostgres(db)
	userService := service.NewUserService(userRepository)
	userHandler := NewUserHandler(userService)

	productRepository := product_postgres.NewProductPostgres(db)
	productService := service.NewProductService(productRepository)
	productHandler := NewProductHandler(productService)

	authService := service.NewAuthService(userRepository, productRepository)

	route := gin.Default()

	userRoute := route.Group("/users")
	{
		userRoute.POST("/register", userHandler.UserRegister)
		userRoute.POST("/login", userHandler.UserLogin)
	}

	productRoute := route.Group("/products")
	{
		productRoute.Use(authService.Authentication())
		productRoute.Use(authService.AuthorizationRole())
		productRoute.POST("/", productHandler.CreateProduct)
		productRoute.GET("/", productHandler.GetProducts)
		productRoute.GET("/:productId", authService.AuthorizationProduct(), productHandler.GetProductById)
		productRoute.PUT("/:productId", productHandler.UpdateProductById)
		productRoute.DELETE("/:productId", productHandler.DeleteProductById)
	}

	route.Run(PORT)
}
