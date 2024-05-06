package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, err
    }

    password := os.Getenv("POSTGRES_DB_PASSWORD")
    if password == "" {
        return nil, fmt.Errorf("POSTGRES_DB_PASSWORD is not set")
    }

    connStr := fmt.Sprintf("user=postgres password=%s dbname=postgres host=localhost port=5432 sslmode=disable", password)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    err =  db.Ping()
    if err != nil {
        return nil, err
    }
    return db, nil
}
