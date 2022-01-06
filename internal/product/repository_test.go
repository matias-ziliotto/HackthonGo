package product

import (
	"context"
	"testing"

	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/pkg/database/sql"
	"github.com/stretchr/testify/assert"
)

var productsToStoreAndGet = []domain.Product{
	{
		Id:          40000,
		Description: "Descripcion x",
		Price:       50.0,
	},
}

var productsToStore = []domain.Product{
	{
		Id:          1000,
		Description: "Descripcion 1000",
		Price:       1000.0,
	},
	{
		Id:          1001,
		Description: "Descripcion 1001",
		Price:       1000.1,
	},
	{
		Id:          1002,
		Description: "Descripcion 1002",
		Price:       1000.2,
	},
}

var productsToStoreErrorPrepare = []domain.Product{}

var productsToStoreErrorStore = []domain.Product{
	{
		Id:          2000,
		Description: "Descripcion 1000",
		Price:       1000.0,
	},
	{
		Id:          2000,
		Description: "Descripcion 1000",
		Price:       1000.0,
	},
}

func TestProductGet(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewProductRepository(db)

	// Act
	productStored, _ := repository.StoreBulk(context.Background(), productsToStoreAndGet)
	result, err := repository.Get(context.Background(), productStored[0].Id)

	// Assert
	assert.Equal(t, productStored[0], result, "result should be equal product stored")
	assert.Nil(t, err, "error should be nil")
}

func TestProductGetNotFound(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewProductRepository(db)

	// Act
	result, err := repository.Get(context.Background(), 99999)

	// Assert
	assert.Error(t, err, "error should exists")
	assert.Equal(t, domain.Product{}, result, "result should be equal to expected result")
}

func TestProductStoreBulk(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewProductRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), productsToStore)

	// Assert
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(result) == 3, "len of result should be equal to 3")
}

func TestProductStoreBulkErrorPrepare(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewProductRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), productsToStoreErrorPrepare)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}

func TestProductStoreBulkError(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewProductRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), productsToStoreErrorStore)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}
