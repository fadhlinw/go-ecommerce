package storer

import "time"

// Product represents a product in the database.
type Product struct {
	ID           int        `db:"id" json:"id"`
	Name         string     `db:"name" json:"name"`
	Image        string     `db:"image" json:"image"`
	Category     string     `db:"category" json:"category"`
	Description  *string    `db:"description" json:"description"`
	Rating       int        `db:"rating" json:"rating"`
	NumReviews   int        `db:"num_reviews" json:"num_reviews"`
	Price        float64    `db:"price" json:"price"`
	CountInStock int        `db:"count_in_stock" json:"count_in_stock"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at"`
}

// Order represents an order in the database.
type Order struct {
	ID            int        `db:"id" json:"id"`
	PaymentMethod string     `db:"payment_method" json:"payment_method"`
	TaxPrice      float64    `db:"tax_price" json:"tax_price"`
	ShippingPrice float64    `db:"shipping_price" json:"shipping_price"`
	TotalPrice    float64    `db:"total_price" json:"total_price"`
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
}

// OrderItem represents an order item in the database.
type OrderItem struct {
	ID        int    `db:"id" json:"id"`
	OrderID   int    `db:"order_id" json:"order_id"`
	ProductID int    `db:"product_id" json:"product_id"`
	Name      string `db:"name" json:"name"`
	Quantity  int    `db:"quantity" json:"quantity"`
	Image     string `db:"image" json:"image"`
	Price     int    `db:"price" json:"price"`
}
