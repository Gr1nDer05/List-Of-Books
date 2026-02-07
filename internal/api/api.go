package api

import "github.com/gin-gonic/gin"

type BooksApi interface {
	Create(c *gin.Context)
	Books(c *gin.Context)
	BookByID(c *gin.Context)
	BookByStatus(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
