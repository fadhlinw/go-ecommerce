package db

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Database wraps the sqlx.DB connection.
// This allows for easy dependency injection and mocking in tests.
type Database struct {
	db *sqlx.DB
}

// NewDatabase establishes and returns a new database connection.
// It reads credentials from the environment variables.
func NewDatabase() (*Database, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format the Data Source Name (DSN) with our credentials.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)
	
	// Open connection to the database. This validates the DSN parameters.
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Return a Database instance containing the connection pointer.
	return &Database{db: db}, nil
}

// Close safely terminates the database connection pool.
func (d *Database) Close() error {
	return d.db.Close()
}

// GetDB returns the underlying sqlx.DB pointer.
// This is used by repository layers to execute queries.
func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
