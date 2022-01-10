package product

import (
	"context"
	"testing"

	"github.com/matias-ziliotto/HackthonGo/internal/customer"
	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/internal/invoice"
	"github.com/matias-ziliotto/HackthonGo/internal/sale"
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

var customers = []domain.Customer{
	{
		Id:        1000,
		FirstName: "Coki",
		LastName:  "Argento",
		Situation: "Activo",
	},
}

var invoices = []domain.Invoice{
	{
		Id:          40000,
		Customer_id: 1000,
		Datetime:    "2022-01-06 11:11:11",
		Total:       0,
	},
}

var sales = []domain.Sale{
	{
		Id:         50000,
		Invoice_id: 40000,
		Product_id: 40000,
		Quantity:   1,
	},
}

func TestProductGet(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
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
	defer db.Close()
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
	defer db.Close()
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
	defer db.Close()
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
	defer db.Close()
	repository := NewProductRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), productsToStoreErrorStore)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}

func TestProductProductsMostSelled(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewProductRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")

	repositoryInvoice := invoice.NewInvoiceRepository(db)
	_, err = repositoryInvoice.StoreBulk(context.Background(), invoices) // insert dummy invoice
	assert.Nil(t, err, "error should be nil")

	// Act
	_, _ = repository.StoreBulk(context.Background(), productsToStoreAndGet)

	repositorySale := sale.NewSaleRepository(db)
	_, err = repositorySale.StoreBulk(context.Background(), sales) // insert dummy sale
	assert.Nil(t, err, "error should be nil")

	result, err := repository.ProductsMostSelled(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}
