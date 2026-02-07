package bookapi

import (
	"listOfBooks/internal/models"
	"listOfBooks/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	StatusPlanned   = "planned"
	StatusReading   = "reading"
	StatusCompleted = "completed"
)

type handler struct {
	service service.BooksService
	log     *slog.Logger
}

func NewHandler(service service.BooksService, log *slog.Logger) *handler {
	return &handler{service: service, log: log}
}

func (h *handler) Books(c *gin.Context) {
	books, err := h.service.Books(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (h *handler) BookByID(c *gin.Context) {
	bookID := c.Param("id")
	book, err := h.service.BookByID(c, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (h *handler) BookByStatus(c *gin.Context) {
	status := c.Query("status")
	book, err := h.service.BookByStatus(c, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (h *handler) Create(c *gin.Context) {
	var req models.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	book, err := h.service.Create(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "book created successfully",
		"book":    book,
	})
}

func (h *handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is requiered"})
		return
	}
	var req models.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	book, err := h.service.Update(c, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "book succefully updated",
		"book":    book,
	})
}

func (h *handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is requiered"})
		return
	}
	if err := h.service.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}
}
