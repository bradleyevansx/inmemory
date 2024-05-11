package stor

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db        *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) QueryRow(query string, args ...interface{}) *sql.Row {
	return r.db.QueryRow(query, args...)
}

func (r *Repository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return rows, nil
}

func (r *Repository) Exec(query string, args ...interface{}) (*sql.Result, error) {
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return &res, nil
}

