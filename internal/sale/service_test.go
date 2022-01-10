package sale

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

var expectedResultGet = domain.Sale{
	Id:         1000,
	Invoice_id: 1000,
	Product_id: 1000,
	Quantity:   1,
}

var expectedResultGetNotFound = domain.Sale{}

func TestServiceSaleGet(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	saleRepository := NewSaleRepository(db)
	saleService := NewSaleService(saleRepository)

	rows := mock.NewRows([]string{"id", "invoice_id", "product_id", "quantity"})
	rows.AddRow(1000, 1000, 1000, 1)
	mock.ExpectQuery(GetSaleQuery).WithArgs(1000).WillReturnRows(rows)

	// Act
	result, err := saleService.Get(context.Background(), 1000)

	// Assert
	assert.Equal(t, expectedResultGet, result, "result should be equal to expected result")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceSaleGetNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	saleRepository := NewSaleRepository(db)
	saleService := NewSaleService(saleRepository)

	mock.ExpectQuery(GetSaleQuery).WithArgs(1000).WillReturnError(ErrorSaleNotFound)

	// Act
	result, err := saleService.Get(context.Background(), 1000)

	// Assert
	assert.Equal(t, expectedResultGetNotFound, result, "result should be equal to expected result")
	assert.Error(t, err, "should exists an error")
}

func TestServiceSaleStoreBulk(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	saleRepository := NewSaleRepository(db)
	saleService := NewSaleService(saleRepository)

	for i := 1; i <= 1000; i++ {
		mock.ExpectQuery(GetSaleQuery).WithArgs(i).WillReturnError(ErrorSaleNotFound)
	}

	mock.ExpectPrepare("INSERT INTO sales")
	mock.ExpectExec("INSERT INTO sales").WillReturnResult(sqlmock.NewResult(1000, 1000))

	// Act
	result, err := saleService.StoreBulk(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceSaleStoreBulkError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	saleRepository := NewSaleRepository(db)
	saleService := NewSaleService(saleRepository)

	for i := 1; i <= 1000; i++ {
		mock.ExpectQuery(GetSaleQuery).WithArgs(i).WillReturnError(ErrorSaleNotFound)
	}

	mock.ExpectPrepare("INSERT INTO sales")
	mock.ExpectExec("INSERT INTO sales").WillReturnError(errors.New("error"))

	// Act
	result, err := saleService.StoreBulk(context.Background())

	// Assert
	assert.Error(t, err, "should exists an error")
	assert.Nil(t, result, "result should be nil")
}
