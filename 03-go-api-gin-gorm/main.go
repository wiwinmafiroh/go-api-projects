package main

import (
	"03-go-api-gin-gorm/database"
	"03-go-api-gin-gorm/routers"
)

func main() {
	var PORT = ":8080"

	database.StartDB()
	routers.StartServer().Run(PORT)
}
