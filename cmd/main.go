package main

import (
	"context"
	"fmt"
	"listOfBooks/database"
	bookapi "listOfBooks/internal/api/bookApi"
	bookrepo "listOfBooks/internal/repository/bookRepo"
	bookservice "listOfBooks/internal/service/bookService"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func loadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "Titanit"),
		ServerPort: getEnv("SERVER_PORT", "8081"),
	}

	// Проверяем обязательные переменные
	if config.DBPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func main() {
	config := loadConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	log.Println("Connecting to database...")
	if err := database.InitDB(dsn); err != nil {
		log.Fatal("Database connection error:", err)
	}
	log.Println("Database connected successfully!")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	bookRepo := bookrepo.NewRepository(database.DB)
	bookService := bookservice.NewService(bookRepo, logger)
	booksHandler := bookapi.NewHandler(bookService, logger)

	server := gin.Default()

	server.GET("/test-db", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var result int
		err := database.DB.QueryRow(ctx, "SELECT 1").Scan(&result)
		if err != nil {
			c.JSON(500, gin.H{"error": "Database test failed: " + err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Database connection OK", "result": result})
	})

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "test",
		})
	})

	server.GET("/books", booksHandler.Books)
	server.GET("/books/:id", booksHandler.BookByID)
	server.GET("/books/status", booksHandler.BookByStatus)

	server.POST("/addbook", booksHandler.Create)

	server.PUT("/books/update/:id", booksHandler.Update)

	server.DELETE("/books/delete/:id", booksHandler.Delete)

	server.Run(":8080")
}
