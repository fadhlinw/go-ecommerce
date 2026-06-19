package storer

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*Storer, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "mysql")
	storer := NewStorer(sqlxDB)

	return storer, mock
}

func TestCreateProduct_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	desc := "A nice phone"
	p := &Product{
		Name:         "Phone",
		Image:        "phone.png",
		Category:     "Electronics",
		Description:  &desc,
		Rating:       4,
		NumReviews:   12,
		Price:        499.99,
		CountInStock: 20,
	}

	mock.ExpectExec("INSERT INTO products").
		WithArgs(p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := s.CreateProduct(context.Background(), p)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.Equal(t, 1, p.ID)
	assert.NotNil(t, p.CreatedAt)
	assert.NotNil(t, p.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProduct_Error(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	p := &Product{Name: "Phone"}
	mock.ExpectExec("INSERT INTO products").WillReturnError(errors.New("db error"))

	id, err := s.CreateProduct(context.Background(), p)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Equal(t, "db error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductByID_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
		AddRow(1, "Phone", "img.png", "Cat", nil, 5, 10, 99.99, 50, now, now)

	mock.ExpectQuery("SELECT \\* FROM products WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	p, err := s.GetProductByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, 1, p.ID)
	assert.Equal(t, "Phone", p.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductByID_NotFound(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	mock.ExpectQuery("SELECT \\* FROM products WHERE id = \\?").
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	p, err := s.GetProductByID(context.Background(), 99)

	assert.Error(t, err)
	assert.Nil(t, p)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListProducts_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
		AddRow(2, "Laptop", "laptop.png", "Electronics", nil, 4, 5, 999.99, 10, now, now).
		AddRow(1, "Phone", "phone.png", "Electronics", nil, 5, 10, 499.99, 20, now, now)

	mock.ExpectQuery("SELECT \\* FROM products ORDER BY id DESC").
		WillReturnRows(rows)

	products, err := s.ListProducts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, 2, products[0].ID)
	assert.Equal(t, 1, products[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListProducts_Error(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	mock.ExpectQuery("SELECT \\* FROM products ORDER BY id DESC").
		WillReturnError(errors.New("db error"))

	products, err := s.ListProducts(context.Background())

	assert.Error(t, err)
	assert.Nil(t, products)
	assert.Equal(t, "db error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	p := &Product{
		ID:           1,
		Name:         "Updated Phone",
		Image:        "img.png",
		Category:     "Electronics",
		Rating:       5,
		NumReviews:   15,
		Price:        450.00,
		CountInStock: 25,
	}

	mock.ExpectExec("UPDATE products SET").
		WithArgs(p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, sqlmock.AnyArg(), p.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.UpdateProduct(context.Background(), p)

	assert.NoError(t, err)
	assert.NotNil(t, p.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_Error(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	p := &Product{ID: 1}

	mock.ExpectExec("UPDATE products SET").
		WillReturnError(errors.New("db error"))

	err := s.UpdateProduct(context.Background(), p)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProduct_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	mock.ExpectExec("DELETE FROM products WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.DeleteProduct(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProduct_Error(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	mock.ExpectExec("DELETE FROM products WHERE id = \\?").
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err := s.DeleteProduct(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrder_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	o := &Order{
		PaymentMethod: "Credit Card",
		TaxPrice:      10.0,
		ShippingPrice: 5.0,
		TotalPrice:    115.0,
	}
	items := []OrderItem{
		{ProductID: 1, Name: "Phone", Quantity: 1, Price: 100},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO orders").
		WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO order_items").
		WithArgs(1, items[0].ProductID, items[0].Name, items[0].Quantity, items[0].Image, items[0].Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	id, err := s.CreateOrder(context.Background(), o, items)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOrder_TxRollback(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	o := &Order{}
	items := []OrderItem{{ProductID: 1}}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO orders").
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	mock.ExpectExec("INSERT INTO order_items").
		WillReturnError(errors.New("insert item failed"))

	mock.ExpectRollback()

	id, err := s.CreateOrder(context.Background(), o, items)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "insert item failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderByID_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
		AddRow(1, "Credit Card", 10.0, 5.0, 115.0, now, now)

	mock.ExpectQuery("SELECT \\* FROM orders WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	o, err := s.GetOrderByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, 1, o.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderItemsByOrderID_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	rows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "name", "quantity", "image", "price"}).
		AddRow(1, 1, 1, "Phone", 1, "img.png", 100)

	mock.ExpectQuery("SELECT \\* FROM order_items WHERE order_id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	items, err := s.GetOrderItemsByOrderID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Phone", items[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListOrders_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
		AddRow(2, "Cash", 0.0, 0.0, 50.0, now, now).
		AddRow(1, "Credit Card", 10.0, 5.0, 115.0, now, now)

	mock.ExpectQuery("SELECT \\* FROM orders ORDER BY id DESC").
		WillReturnRows(rows)

	orders, err := s.ListOrders(context.Background())

	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, 2, orders[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOrder_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	o := &Order{
		ID:            1,
		PaymentMethod: "Transfer",
		TaxPrice:      0,
		ShippingPrice: 0,
		TotalPrice:    100,
	}

	mock.ExpectExec("UPDATE orders SET").
		WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, sqlmock.AnyArg(), o.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.UpdateOrder(context.Background(), o)

	assert.NoError(t, err)
	assert.NotNil(t, o.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteOrder_Success(t *testing.T) {
	s, mock := setupMockDB(t)
	defer s.db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM order_items WHERE order_id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM orders WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := s.DeleteOrder(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
