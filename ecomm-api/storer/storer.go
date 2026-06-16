package storer

import (
	"github.com/jmoiron/sqlx"
)

// Storer wraps the raw sqlx DB connection.
// This is where you will define methods to interact with your tables (e.g., GetProduct, CreateOrder).
type Storer struct {
	db *sqlx.DB
}

// NewStorer is the constructor that creates a new Storer instance.
// We inject the raw database connection here, enforcing dependency injection.
func NewStorer(db *sqlx.DB) *Storer {
	return &Storer{db: db}
}
