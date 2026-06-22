package storer

import (
	"context"
	"fmt"
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

// execTx executes a function within a database transaction.
func (s *Storer) execTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// CreateOrder inserts a new order and its items into the database.
func (s *Storer) CreateOrder(ctx context.Context, o *Order, items []OrderItem) (int, error) {
	var orderID int
	err := s.execTx(ctx, func(tx *sqlx.Tx) error {
		now := time.Now()
		if o.CreatedAt == nil {
			o.CreatedAt = &now
		}
		o.UpdatedAt = &now

		queryOrder := `
			INSERT INTO orders (payment_method, tax_price, shipping_price, total_price, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		res, err := tx.ExecContext(ctx, queryOrder, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		orderID = int(id)
		o.ID = orderID

		if len(items) > 0 {
			queryItems := `
				INSERT INTO order_items (order_id, product_id, name, quantity, image, price)
				VALUES (?, ?, ?, ?, ?, ?)
			`
			for _, item := range items {
				_, err := tx.ExecContext(ctx, queryItems, orderID, item.ProductID, item.Name, item.Quantity, item.Image, item.Price)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return orderID, nil
}

// GetOrderByID retrieves an order by its ID.
func (s *Storer) GetOrderByID(ctx context.Context, id int) (*Order, error) {
	query := `SELECT * FROM orders WHERE id = ?`
	var o Order
	err := s.db.GetContext(ctx, &o, query, id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetOrderItemsByOrderID retrieves order items by order ID.
func (s *Storer) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]OrderItem, error) {
	query := `SELECT * FROM order_items WHERE order_id = ?`
	var items []OrderItem
	err := s.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// ListOrders retrieves a list of orders.
func (s *Storer) ListOrders(ctx context.Context) ([]Order, error) {
	query := `SELECT * FROM orders ORDER BY id DESC`
	var orders []Order
	err := s.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateOrder updates an existing order's details.
func (s *Storer) UpdateOrder(ctx context.Context, o *Order) error {
	query := `
		UPDATE orders
		SET payment_method = ?, tax_price = ?, shipping_price = ?, total_price = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	o.UpdatedAt = &now

	_, err := s.db.ExecContext(ctx, query,
		o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.UpdatedAt, o.ID,
	)
	return err
}

// DeleteOrder deletes an order and its items by order ID.
func (s *Storer) DeleteOrder(ctx context.Context, id int) error {
	return s.execTx(ctx, func(tx *sqlx.Tx) error {
		// Delete order items first due to foreign key constraints
		_, err := tx.ExecContext(ctx, `DELETE FROM order_items WHERE order_id = ?`, id)
		if err != nil {
			return err
		}

		// Delete the order
		_, err = tx.ExecContext(ctx, `DELETE FROM orders WHERE id = ?`, id)
		return err
	})
}
