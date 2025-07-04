package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON urls (alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil

}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"
	stmt, err := s.db.Prepare(`
		INSERT INTO urls (alias, url) VALUES (?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"
	stmt, err := s.db.Prepare(`
		SELECT url FROM urls WHERE alias = ?
	`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resUrl string

	err = stmt.QueryRow(alias).Scan(&resUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.sqlite.DeleteUrl"
	stmt, err := s.db.Prepare(`
		DELETE FROM urls WHERE alias = ?
	`)
	_, err = stmt.Exec(alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
