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
