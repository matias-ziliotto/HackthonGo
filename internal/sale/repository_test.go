package sale

import (
	"context"
	"testing"

	"github.com/matias-ziliotto/HackthonGo/internal/customer"
	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/internal/invoice"
	"github.com/matias-ziliotto/HackthonGo/internal/product"
	"github.com/matias-ziliotto/HackthonGo/pkg/database/sql"
	"github.com/stretchr/testify/assert"
)

var salesToStoreAndGet = []domain.Sale{
	{
		Id:         400000,
		Invoice_id: 1000,
		Product_id: 2000,
		Quantity:   1,
	},
}

var salesToStore = []domain.Sale{
	{
		Id:         10000,
		Invoice_id: 1000,
		Product_id: 2000,
		Quantity:   1,
	},
	{
		Id:         10001,
		Invoice_id: 1000,
		Product_id: 2000,
		Quantity:   2,
	},
	{
		Id:         10002,
		Invoice_id: 1000,
		Product_id: 2000,
		Quantity:   3,
	},
}

var salesToStoreErrorPrepare = []domain.Sale{}

var salesToStoreErrorStore = []domain.Sale{
	{
		Id:         100000,
		Invoice_id: 1000,
		Product_id: 2000,
		Quantity:   1,
	},
	{
		Id:         100000,
		Invoice_id: 1000,
		Product_id: 2000,
		Quantity:   1,
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
		Id:          1000,
		Customer_id: 1000,
		Datetime:    "2022-01-06 11:11:11",
		Total:       200.5,
	},
}

var products = []domain.Product{
	{
		Id:          2000,
		Description: "Descripcion x",
		Price:       50.0,
	},
}

func TestSaleGet(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewSaleRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	repositoryInvoice := invoice.NewInvoiceRepository(db)
	repositoryProduct := product.NewProductRepository(db)

	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")
	_, err = repositoryInvoice.StoreBulk(context.Background(), invoices) // insert dummy invoice
	assert.Nil(t, err, "error should be nil")
	_, err = repositoryProduct.StoreBulk(context.Background(), products) // insert dummy invoice
	assert.Nil(t, err, "error should be nil")

	// Act
	saleStored, _ := repository.StoreBulk(context.Background(), salesToStoreAndGet)
	result, err := repository.Get(context.Background(), saleStored[0].Id)

	// Assert
	assert.Equal(t, saleStored[0], result, "result should be equal sale stored")
	assert.Nil(t, err, "error should be nil")
}

func TestSaleGetNotFound(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewSaleRepository(db)

	// Act
	result, err := repository.Get(context.Background(), 99999)

	// Assert
	assert.Error(t, err, "error should exists")
	assert.Equal(t, domain.Sale{}, result, "result should be equal to expected result")
}

func TestSaleStoreBulk(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewSaleRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	repositoryInvoice := invoice.NewInvoiceRepository(db)
	repositoryProduct := product.NewProductRepository(db)

	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")
	_, err = repositoryInvoice.StoreBulk(context.Background(), invoices) // insert dummy invoice
	assert.Nil(t, err, "error should be nil")
	_, err = repositoryProduct.StoreBulk(context.Background(), products) // insert dummy invoice
	assert.Nil(t, err, "error should be nil")

	// Act
	result, err := repository.StoreBulk(context.Background(), salesToStore)

	// Assert
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(result) == 3, "len of result should be equal to 3")
}

func TestSaleStoreBulkErrorPrepare(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewSaleRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), salesToStoreErrorPrepare)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}

func TestSaleStoreBulkError(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewSaleRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	repositoryInvoice := invoice.NewInvoiceRepository(db)
	repositoryProduct := product.NewProductRepository(db)

	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")
	_, err = repositoryInvoice.StoreBulk(context.Background(), invoices) // insert dummy invoice
	assert.Nil(t, err, "error should be nil")
	_, err = repositoryProduct.StoreBulk(context.Background(), products) // insert dummy invoice
	assert.Nil(t, err, "error should be nil")

	// Act
	result, err := repository.StoreBulk(context.Background(), salesToStoreErrorStore)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}
