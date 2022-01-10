package invoice

import (
	"context"
	"testing"

	"github.com/matias-ziliotto/HackthonGo/internal/customer"
	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/internal/product"
	"github.com/matias-ziliotto/HackthonGo/internal/sale"
	"github.com/matias-ziliotto/HackthonGo/pkg/database/sql"
	"github.com/stretchr/testify/assert"
)

var invoicesToStoreAndGet = []domain.Invoice{
	{
		Id:          40000,
		Customer_id: 1000,
		Datetime:    "2022-01-06 11:11:11",
		Total:       0,
	},
}

var invoiceToUpdate = domain.Invoice{
	Id:    1000,
	Total: 1,
}

var invoicesToStore = []domain.Invoice{
	{
		Id:          1000,
		Customer_id: 1000,
		Datetime:    "2022-01-06 11:11:11",
		Total:       200.5,
	},
	{
		Id:          1001,
		Customer_id: 1000,
		Datetime:    "2022-01-06 11:11:12",
		Total:       300.5,
	},
	{
		Id:          1002,
		Customer_id: 1000,
		Datetime:    "2022-01-06 11:11:13",
		Total:       400.5,
	},
}

var invoicesToStoreErrorPrepare = []domain.Invoice{}

var invoicesToStoreErrorStore = []domain.Invoice{
	{
		Id:          10000,
		Customer_id: 1000,
		Datetime:    "2022-01-06 11:11:12",
		Total:       300.5,
	},
	{
		Id:          10000,
		Customer_id: 1000,
		Datetime:    "2022-01-06 11:11:12",
		Total:       300.5,
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

var products = []domain.Product{
	{
		Id:          50000,
		Description: "Description de producto",
		Price:       100.5,
	},
}

var sales = []domain.Sale{
	{
		Id:         50000,
		Invoice_id: 40000,
		Product_id: 50000,
		Quantity:   1,
	},
}

func TestInvoiceGet(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewInvoiceRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")

	// Act
	invoiceStored, _ := repository.StoreBulk(context.Background(), invoicesToStoreAndGet)
	result, err := repository.Get(context.Background(), invoiceStored[0].Id)

	// Assert
	assert.Equal(t, invoiceStored[0], result, "result should be equal invoice stored")
	assert.Nil(t, err, "error should be nil")
}

func TestInvoiceGetNotFound(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewInvoiceRepository(db)

	// Act
	result, err := repository.Get(context.Background(), 99999)

	// Assert
	assert.Error(t, err, "error should exists")
	assert.Equal(t, domain.Invoice{}, result, "result should be equal to expected result")
}

func TestInvoiceStoreBulk(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewInvoiceRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")

	// Act
	result, err := repository.StoreBulk(context.Background(), invoicesToStore)

	// Assert
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(result) == 3, "len of result should be equal to 3")
}

func TestInvoiceStoreBulkErrorPrepare(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewInvoiceRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), invoicesToStoreErrorPrepare)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}

func TestInvoiceStoreBulkError(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewInvoiceRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")

	// Act
	result, err := repository.StoreBulk(context.Background(), invoicesToStoreErrorStore)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}

func TestInvoiceUpdateTotal(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewInvoiceRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")

	// Act
	_, err = repository.StoreBulk(context.Background(), invoicesToStore)
	assert.Nil(t, err, "error should be nil")

	invoiceUpdated, _ := repository.UpdateTotal(context.Background(), invoiceToUpdate)

	// Assert
	assert.Equal(t, invoiceToUpdate, invoiceUpdated, "invoice updated should be equal invoice to update")
	assert.Nil(t, err, "error should be nil")
}

func TestInvoiceGetAllTotalEmpty(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewInvoiceRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")

	// Act
	_, _ = repository.StoreBulk(context.Background(), invoicesToStoreAndGet)
	result, err := repository.GetAllTotalEmpty(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestInvoiceCalculateTotal(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	defer db.Close()
	repository := NewInvoiceRepository(db)

	repositoryCustomer := customer.NewCustomerRepository(db)
	_, err = repositoryCustomer.StoreBulk(context.Background(), customers) // insert dummy customer
	assert.Nil(t, err, "error should be nil")

	repositoryProduct := product.NewProductRepository(db)
	_, err = repositoryProduct.StoreBulk(context.Background(), products) // insert dummy product
	assert.Nil(t, err, "error should be nil")

	// Act
	_, _ = repository.StoreBulk(context.Background(), invoicesToStoreAndGet)

	repositorySale := sale.NewSaleRepository(db)
	_, err = repositorySale.StoreBulk(context.Background(), sales) // insert dummy customer
	assert.Nil(t, err, "error should be nil")

	ids, _ := repository.GetAllTotalEmpty(context.Background())
	result, err := repository.CalculateTotal(context.Background(), ids)

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}
