package db

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func applyMigrations(db *sql.DB) error {
	migrations, err := os.ReadFile("internal/db/migrations/001_initial_schema.sql")
	if err != nil {
		return fmt.Errorf("Failed to read migrations: %v", err)
	}

	_, err = db.Exec(string(migrations))
	return err
}

func newDB(dbPath string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("Failed to create database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath) 
	if err != nil {
		return nil, fmt.Errorf("Failed to open database: %v", err)
	}

	// migrations
	if err := applyMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to apply migrations: %v", err)
	}

	return &DB{db}, nil
}