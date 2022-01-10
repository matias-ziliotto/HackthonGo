package invoice

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

var expectedResultGet = domain.Invoice{
	Id:          1000,
	Customer_id: 1000,
	Datetime:    "2022-01-10 14:16:05",
	Total:       200.5,
}

var expectedResultGetNotFound = domain.Invoice{}

func TestServiceInvoiceGet(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	rows := mock.NewRows([]string{"id", "invoice_id", "datetime", "total"})
	rows.AddRow(1000, 1000, "2022-01-10 14:16:05", 200.5)
	mock.ExpectQuery(GetInvoiceQuery).WithArgs(1000).WillReturnRows(rows)

	// Act
	result, err := invoiceService.Get(context.Background(), 1000)

	// Assert
	assert.Equal(t, expectedResultGet, result, "result should be equal to expected result")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceInvoiceGetNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	mock.ExpectQuery(GetInvoiceQuery).WithArgs(1000).WillReturnError(errors.New("error"))

	// Act
	result, err := invoiceService.Get(context.Background(), 1000)

	// Assert
	assert.Equal(t, expectedResultGetNotFound, result, "result should be equal to expected result")
	assert.Error(t, err, "should exists an error")
}

func TestServiceInvoiceStoreBulk(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	for i := 1; i <= 50; i++ {
		mock.ExpectQuery(GetInvoiceQuery).WithArgs(i).WillReturnError(ErrorInvoiceNotFound)
	}

	mock.ExpectPrepare("INSERT INTO invoices")
	mock.ExpectExec("INSERT INTO invoices").WillReturnResult(sqlmock.NewResult(100, 100))

	// Act
	result, err := invoiceService.StoreBulk(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceInvoiceStoreBulkError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	for i := 1; i <= 50; i++ {
		mock.ExpectQuery(GetInvoiceQuery).WithArgs(i).WillReturnError(ErrorInvoiceNotFound)
	}

	mock.ExpectPrepare("INSERT INTO customers")
	mock.ExpectExec("INSERT INTO customers").WillReturnError(errors.New("error"))

	// Act
	result, err := invoiceService.StoreBulk(context.Background())

	// Assert
	assert.Error(t, err, "should exists an error")
	assert.Nil(t, result, "result should be nil")
}

func TestServiceInvoiceUpdateTotal(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	rowsGetAllTotalEmpty := mock.NewRows([]string{"id"})
	rowsGetAllTotalEmpty.AddRow(1)
	rowsGetAllTotalEmpty.AddRow(2)
	rowsGetAllTotalEmpty.AddRow(3)
	mock.ExpectQuery(GetAllTotalEmptyInvoiceQuery).WillReturnRows(rowsGetAllTotalEmpty)

	rowsCalculateTotal := mock.NewRows([]string{"id", "total"})
	rowsCalculateTotal.AddRow(1, 100.0)
	rowsCalculateTotal.AddRow(2, 200.0)
	rowsCalculateTotal.AddRow(3, 300.0)
	mock.ExpectQuery("SELECT DISTINCT").WillReturnRows(rowsCalculateTotal)

	for i := 1; i <= 3; i++ {
		mock.ExpectPrepare("UPDATE invoices SET total")
		mock.ExpectExec("UPDATE invoices SET total").WithArgs(float64(i*100), i).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	// Act
	result, err := invoiceService.UpdateTotal(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceInvoiceUpdateTotalGetAllTotalEmptyError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	mock.ExpectQuery(GetAllTotalEmptyInvoiceQuery).WillReturnError(errors.New("error"))

	// Act
	result, err := invoiceService.UpdateTotal(context.Background())

	// Assert
	assert.Error(t, err, "should exists an error")
	assert.Nil(t, result, "result should be nil")
}

func TestServiceInvoiceUpdateTotalGetAllTotalEmpty(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	rows := mock.NewRows([]string{"id"})
	mock.ExpectQuery(GetAllTotalEmptyInvoiceQuery).WillReturnRows(rows)

	// Act
	result, err := invoiceService.UpdateTotal(context.Background())

	// Assert
	assert.Nil(t, err, "error should be nil")
	assert.Nil(t, result, "result should be nil")
}

func TestServiceInvoiceUpdateTotalCalculateTotal(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	rowsGetAllTotalEmpty := mock.NewRows([]string{"id"})
	rowsGetAllTotalEmpty.AddRow(1)
	mock.ExpectQuery(GetAllTotalEmptyInvoiceQuery).WillReturnRows(rowsGetAllTotalEmpty)

	mock.ExpectQuery("SELECT DISTINCT").WillReturnError(errors.New("error"))

	// Act
	result, err := invoiceService.UpdateTotal(context.Background())

	// Assert
	assert.Error(t, err, "should exists an error")
	assert.Nil(t, result, "result should be nil")
}

func TestServiceInvoiceUpdateTotalError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	invoiceRepository := NewInvoiceRepository(db)
	invoiceService := NewInvoiceService(invoiceRepository)

	rowsGetAllTotalEmpty := mock.NewRows([]string{"id"})
	rowsGetAllTotalEmpty.AddRow(1)
	rowsGetAllTotalEmpty.AddRow(2)
	rowsGetAllTotalEmpty.AddRow(3)
	mock.ExpectQuery(GetAllTotalEmptyInvoiceQuery).WillReturnRows(rowsGetAllTotalEmpty)

	rowsCalculateTotal := mock.NewRows([]string{"id", "total"})
	rowsCalculateTotal.AddRow(1, 100.0)
	rowsCalculateTotal.AddRow(2, 200.0)
	rowsCalculateTotal.AddRow(3, 300.0)
	mock.ExpectQuery("SELECT DISTINCT").WillReturnRows(rowsCalculateTotal)

	for i := 1; i <= 3; i++ {
		mock.ExpectPrepare("UPDATE invoices SET total")
		mock.ExpectExec("UPDATE invoices SET total").WithArgs(float64(i*100), i).WillReturnError(errors.New("error"))
	}

	// Act
	result, err := invoiceService.UpdateTotal(context.Background())

	// Assert
	assert.Error(t, err, "should exists an error")
	assert.Nil(t, result, "result should be nil")
}
