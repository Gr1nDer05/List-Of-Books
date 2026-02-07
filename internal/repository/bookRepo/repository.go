package bookrepo

import (
	"context"
	"database/sql"
	"errors"
	"listOfBooks/database"
	"listOfBooks/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) Books(ctx context.Context) ([]models.Book, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, author, year, status FROM books`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Status)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *repository) BookByID(ctx context.Context, id string) (*models.Book, error) {
	var book models.Book
	err := r.db.QueryRow(ctx, `SELECT id, title, author, year, status FROM books WHERE id = $1`, id).Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Book{}, nil
		}
		return &models.Book{}, err
	}
	return &book, nil
}

func (r *repository) BookByStatus(ctx context.Context, status string) ([]models.Book, error) {
	var books []models.Book
	rows, err := r.db.Query(ctx, `SELECT id, title, author, year, status FROM books WHERE status = $1`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year, &book.Status)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *repository) Create(ctx context.Context, req models.BookRequest) (*models.Book, error) {
	var id string
	err := r.db.QueryRow(ctx, `INSERT INTO books (title, author, year, status) VALUES ($1, $2, $3, $4) RETURNING id`, req.Title, req.Author, req.Year, req.Status).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &models.Book{
		ID:     id,
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
		Status: req.Status,
	}, nil
}

func (r *repository) Update(ctx context.Context, id string, req models.BookRequest) (*models.Book, error) {
	_, err := database.DB.Exec(ctx, `UPDATE books SET title = $1, author = $2, year = $3, status = $4 WHERE id = $5`, req.Title, req.Author, req.Year, req.Status, id)
	if err != nil {
		return nil, err
	}
	return &models.Book{
		ID:     id,
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
		Status: req.Status,
	}, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	var exist bool
	if err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM books WHERE id = $1)`, id).Scan(&exist); err != nil {
		return err
	}
	_, err := database.DB.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
