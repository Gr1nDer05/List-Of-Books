package bookservice

import (
	"context"
	"fmt"
	"listOfBooks/internal/models"
		"listOfBooks/internal/repository"
	"log/slog"
)

const (
	StatusPlanned   = "planned"
	StatusReading   = "reading"
	StatusCompleted = "completed"
)

type service struct {
	repository repository.BooksRepository
	log        *slog.Logger
}

func NewService(repository repository.BooksRepository, log *slog.Logger) *service {
	return &service{repository: repository, log: log}
}

func (s *service) Books(ctx context.Context) ([]models.Book, error) {
	const op = "service.getBooks"
	s.log.With(
		slog.String("op", op),
	)
	s.log.Info("getting books")

	books, err := s.repository.Books(ctx)
	if err != nil {
		s.log.Error("failed to get books", "error", err.Error())
		return nil, err
	}
	s.log.Info("books getted")
	return books, nil
}

func (s *service) BookByID(ctx context.Context, id string) (*models.Book, error) {
	const op = "service.getBooksByID"
	s.log.With(
		slog.String("op", op),
	)
	s.log.Info("getting book by id")
	books, err := s.repository.BookByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get book", "error", err.Error())
		return nil, err
	}
	return books, nil
}

func (s *service) BookByStatus(ctx context.Context, status string) ([]models.Book, error) {
	const op = "service.getBooksByStatus"
	s.log.With(
		slog.String("op", op),
	)
	s.log.Info("getting book by status")
	book, err := s.repository.BookByStatus(ctx, status)
	if err != nil {
		s.log.Error("failed to get book", "error", err.Error())
		return nil, err
	}
	return book, nil
}

func (s *service) Create(ctx context.Context, req models.BookRequest) (*models.Book, error) {
	const op = "service.CreateBook"
	s.log.With(
		slog.String("op", op),
	)
	s.log.Info("creating book")
	reqBook := models.BookRequest{
		Title:  req.Title,
		Author: req.Author,
		Status: req.Status,
		Year:   req.Year,
	}

	if req.Year > 2025 {
		return nil, fmt.Errorf("year cannot be bigger than 2025")
	}
	if req.Status != StatusPlanned && req.Status != StatusReading && req.Status != StatusCompleted {
		return nil, fmt.Errorf("status must be planned, reading, completed")
	}

	book, err := s.repository.Create(ctx, reqBook)
	if err != nil {
		s.log.Error("failed to create book", "error", err.Error())
		return nil, err
	}
	return book, nil
}

func (s *service) Update(ctx context.Context, id string, req models.BookRequest) (*models.Book, error) {
	const op = "service.UpdateBook"
	s.log.With(
		slog.String("op", op),
	)
	s.log.Info("updating book")
	book, err := s.repository.Update(ctx, id, req)
	if err != nil {
		s.log.Error("failed to create book", "error", err.Error())
		return nil, err
	}
	return book, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	const op = "service.DeleteBook"
	s.log.With(
		slog.String("op", op),
	)
	s.log.Info("deleting book")
	if err := s.repository.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete book", "error", err.Error())
		return err
	}
	return nil
}
