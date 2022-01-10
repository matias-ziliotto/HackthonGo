package product

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

var expectedResultGet = domain.Product{
	Id:          1000,
	Description: "Mate",
	Price:       1250.5,
}

var expectedResultGetNotFound = domain.Product{}

func TestServiceProductGet(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	productRepository := NewProductRepository(db)
	productService := NewProductService(productRepository)

	rows := mock.NewRows([]string{"id", "description", "price"})
	rows.AddRow(1000, "Mate", 1250.5)
	mock.ExpectQuery(GetProductQuery).WithArgs(1000).WillReturnRows(rows)

	// Act
	result, err := productService.Get(context.Background(), 1000)

	// Assert
	assert.Equal(t, expectedResultGet, result, "result should be equal to expected result")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceProductGetNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	productRepository := NewProductRepository(db)
	productService := NewProductService(productRepository)

	mock.ExpectQuery(GetProductQuery).WithArgs(1000).WillReturnError(ErrorProductNotFound)

	// Act
	result, err := productService.Get(context.Background(), 1000)

	// Assert
	assert.Equal(t, expectedResultGetNotFound, result, "result should be equal to expected result")
	assert.Error(t, err, "should exists an error")
}

func TestServiceProductStoreBulk(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	productRepository := NewProductRepository(db)
	productService := NewProductService(productRepository)

	for i := 1; i <= 100; i++ {
		mock.ExpectQuery(GetProductQuery).WithArgs(i).WillReturnError(ErrorProductNotFound)
	}

	mock.ExpectPrepare("INSERT INTO products")
	mock.ExpectExec("INSERT INTO products").WillReturnResult(sqlmock.NewResult(100, 100))

	// Act
	result, err := productService.StoreBulk(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceProductStoreBulkError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	productRepository := NewProductRepository(db)
	productService := NewProductService(productRepository)

	for i := 1; i <= 100; i++ {
		mock.ExpectQuery(GetProductQuery).WithArgs(i).WillReturnError(ErrorProductNotFound)
	}

	mock.ExpectPrepare("INSERT INTO products")
	mock.ExpectExec("INSERT INTO products").WillReturnError(errors.New("error"))

	// Act
	result, err := productService.StoreBulk(context.Background())

	// Assert
	assert.Error(t, err, "should exists an error")
	assert.Nil(t, result, "result should be nil")
}

func TestServiceProductGetMostSelled(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	productRepository := NewProductRepository(db)
	productService := NewProductService(productRepository)

	rows := mock.NewRows([]string{"id", "description", "price"})
	rows.AddRow(1, "Mate", 1250.5)
	mock.ExpectQuery("SELECT COUNT").WillReturnRows(rows)

	// Act
	result, err := productService.GetProductsMostSelled(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceProductGetMostSelledError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	productRepository := NewProductRepository(db)
	productService := NewProductService(productRepository)

	mock.ExpectQuery("SELECT COUNT").WillReturnError(errors.New("error"))

	// Act
	result, err := productService.GetProductsMostSelled(context.Background())

	// Assert
	assert.Error(t, err, "should exists an error")
	assert.Nil(t, result, "result should be nil")
}
