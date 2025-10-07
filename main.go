package main

import (
	"listOfBooks/database"
	"listOfBooks/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := "postgres://postgres:postgres@localhost:5432/books?sslmode=disable"
	if err := database.InitDB(dsn); err != nil {
		log.Fatal("error", err)
	}

	server := gin.Default()

	server.GET("/books", handlers.GetBooks)
	server.GET("/books/:id", handlers.GetBookByID)
	server.GET("/books/status", handlers.GetBookByStatus)
	server.POST("/books", handlers.AddBook)
	server.PUT("/books/update/:id", handlers.UpdateBook)
	server.DELETE("/books/delete/:id", handlers.DeleteBook)
	server.Run(":8080")
}
