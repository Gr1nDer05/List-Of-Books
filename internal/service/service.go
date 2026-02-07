package service

import (
	"context"
	"listOfBooks/internal/models"
)

type BooksService interface {
	Create(ctx context.Context, req models.BookRequest) (*models.Book, error)
	Books(ctx context.Context) ([]models.Book, error)
	BookByID(ctx context.Context, id string) (*models.Book, error)
	BookByStatus(ctx context.Context, status string) ([]models.Book, error)
	Update(ctx context.Context, id string, req models.BookRequest) (*models.Book, error)
	Delete(ctx context.Context, id string) error
}
