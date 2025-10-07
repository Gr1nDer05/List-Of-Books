package models

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Status string `json:"status"` //planned, reading, completed
}

type BookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type StatusRequest struct {
	Status string `json:"status" binding:"required"`
}
