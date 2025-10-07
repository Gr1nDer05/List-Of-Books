package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(connectionString string) error {
	var err error

	config, err := pgxpool.ParseConfig(connectionString)

	if err != nil {
		return err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return err
	}

	DB = pool

	if err := createTables(); err != nil {
		return err
	}

	log.Println("Database succefully initialized")
	return nil
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		year INT NOT NULL,
		STATUS TEXT NOT NULL
	);
	`
	_, err := DB.Exec(context.Background(), query)
	return err
}
