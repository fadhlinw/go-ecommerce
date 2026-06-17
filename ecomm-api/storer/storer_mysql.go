package storer

import (
	"context"
	"time"

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

// CreateProduct inserts a new product into the database and returns the newly generated ID.
func (s *Storer) CreateProduct(ctx context.Context, p *Product) (int, error) {
	query := `
		INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	if p.CreatedAt == nil {
		p.CreatedAt = &now
	}
	p.UpdatedAt = &now

	result, err := s.db.ExecContext(ctx, query,
		p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	p.ID = int(id)
	return p.ID, nil
}

// GetProductByID retrieves a single product by its ID.
func (s *Storer) GetProductByID(ctx context.Context, id int) (*Product, error) {
	query := `SELECT * FROM products WHERE id = ?`
	var p Product
	err := s.db.GetContext(ctx, &p, query, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListProducts retrieves a list of products.
func (s *Storer) ListProducts(ctx context.Context) ([]Product, error) {
	query := `SELECT * FROM products ORDER BY id DESC`
	var products []Product
	err := s.db.SelectContext(ctx, &products, query)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct updates an existing product's details.
func (s *Storer) UpdateProduct(ctx context.Context, p *Product) error {
	query := `
		UPDATE products
		SET name = ?, image = ?, category = ?, description = ?, rating = ?, num_reviews = ?, price = ?, count_in_stock = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	p.UpdatedAt = &now

	_, err := s.db.ExecContext(ctx, query,
		p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.UpdatedAt, p.ID,
	)
	return err
}

// DeleteProduct deletes a product by its ID.
func (s *Storer) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
