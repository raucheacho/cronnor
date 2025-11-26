package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Repository handles database operations
type Repository struct {
	db *sql.DB
}

// New creates a new repository instance
func New(dbPath string) (*Repository, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Set pragmas for better performance
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to set journal mode: %w", err)
	}

	return &Repository{db: db}, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}

// RunMigrations applies database migrations
func (r *Repository) RunMigrations(migrationsPath string) error {
	schemaSQL, err := os.ReadFile(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	if _, err := r.db.Exec(string(schemaSQL)); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// DB returns the underlying database connection
func (r *Repository) DB() *sql.DB {
	return r.db
}
