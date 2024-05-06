package main

import (
	"database/sql"
	"fmt"
)

type IRepository[T IEntity] interface {
	Get() ([]T, error)
}

type BaseRepository[T IEntity] struct {
	db        *sql.DB
	tableName string
	rowMapper func(scan func(dest ...any) error) (*T, error)
}

func NewRepository[T IEntity](db *sql.DB, tableName string, rowMapper  func(scan func(dest ...any) error) (*T, error)) *BaseRepository[T] {
	return &BaseRepository[T]{db: db, tableName: tableName, rowMapper: rowMapper}
}

func (e *BaseRepository[T]) Get() ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s", e.tableName)

	rows, err := e.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var results []T

	for rows.Next() {
		res, err := e.rowMapper(rows.Scan)
		if err != nil{
			return nil, fmt.Errorf("error parsing row: %w", err)
		}

		results = append(results, *res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return results, nil
}
