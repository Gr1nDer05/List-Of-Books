package handlers

import (
	"database/sql"
	"errors"
	"listOfBooks/database"
	"listOfBooks/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	StatusPlanned   = "planned"
	StatusReading   = "reading"
	StatusCompleted = "completed"
)

func GetBooks(c *gin.Context) {
	rows, err := database.DB.Query(c, `SELECT id, title, author, year, status FROM books`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при открытии бд",
		})
		return
	}

	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Не удалось найти книги",
			})
			return
		}

		books = append(books, book)
	}
	if len(books) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Книг нет",
			"книги":   books,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"количество книг": len(books),
		"книги":           books,
	})
}

func GetBookByID(c *gin.Context) {
	strid := c.Param("id")

	id, err := strconv.Atoi(strid)
	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неправильный id",
		})
		return
	}

	var book models.Book

	err = database.DB.QueryRow(c, `SELECT id, title, author, year, status FROM books WHERE id = $1`, id).Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Книга не найдена",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"книга по id": book,
	})
}

func GetBookByStatus(c *gin.Context) {
	var statusRequest models.StatusRequest

	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неправильный формат данных",
		})
		return
	}
	if statusRequest.Status != StatusPlanned && statusRequest.Status != StatusReading && statusRequest.Status != StatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неправильный статус",
			"message": "Статус должен быть planned, reading, completed",
		})
		return
	}

	var books []models.Book

	rows, err := database.DB.Query(c, `SELECT id, title, author, year, status FROM books WHERE status = $1`, statusRequest.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка при открытии бд",
		})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Не удалось найти книги",
			})
			return
		}
		books = append(books, book)
	}
	c.JSON(http.StatusOK, gin.H{
		"книги с выбранным статусом": books,
	})
}

func AddBook(c *gin.Context) {
	var request models.BookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неправильный формат данных",
		})
		return
	}

	if request.Year > 2025 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Год не может быть больше 2025",
		})
	}

	if request.Status != StatusPlanned && request.Status != StatusReading && request.Status != StatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неправильный статус",
			"message": "Статус должен быть planned, reading, completed",
		})
		return
	}

	var id int64
	err := database.DB.QueryRow(c, `INSERT INTO books (title, author, year, status) VALUES ($1, $2, $3, $4) RETURNING id`, request.Title, request.Author, request.Year, request.Status).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при добавлении книги",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Книга успешно добавлена",
		"книга":   request,
		"ID":      id,
	})
}

func UpdateBook(c *gin.Context) {
	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неправильный id",
		})
		return
	}

	var request models.BookRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неправильный формат данных",
		})
		return
	}
	_, err = database.DB.Exec(c, `UPDATE books SET title = $1, author = $2, year = $3, status = $4 WHERE id = $5`, request.Title, request.Author, request.Year, request.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при обновлении книги",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Книга успешно обновлена",
		"книга":   request,
	})
}

func DeleteBook(c *gin.Context) {
	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неправильный id",
		})
		return
	}
	var exist bool
	if err = database.DB.QueryRow(c, `SELECT EXISTS (SELECT 1 FROM books WHERE id = $1)`, id).Scan(&exist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при удалении книги",
		})
		return
	}
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Книга не найдена",
		})
		return
	}

	_, err = database.DB.Exec(c, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при удалении книги",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Книга успешно удалена",
	})
}
