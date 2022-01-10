package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/stretchr/testify/assert"
)

var expectedResultGet = domain.Customer{
	Id:        1000,
	FirstName: "Pepe",
	LastName:  "Argento",
	Situation: "Inactivo",
}

var expectedResultGetNotFound = domain.Customer{}

func TestServiceCustomerGet(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	customerRepository := NewCustomerRepository(db)
	customerService := NewCustomerService(customerRepository)

	rows := mock.NewRows([]string{"id", "first_name", "last_name", "situation"})
	rows.AddRow(1000, "Pepe", "Argento", "Inactivo")
	mock.ExpectQuery(GetCustomerQuery).WithArgs(1000).WillReturnRows(rows)

	// Act
	result, err := customerService.Get(context.Background(), 1000)

	// Assert
	assert.Equal(t, expectedResultGet, result, "result should be equal to expected result")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceCustomerGetNotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	customerRepository := NewCustomerRepository(db)
	customerService := NewCustomerService(customerRepository)

	mock.ExpectQuery(GetCustomerQuery).WithArgs(1000).WillReturnError(ErrorCustomerNotFound)

	// Act
	result, err := customerService.Get(context.Background(), 1000)

	// Assert
	assert.Equal(t, expectedResultGetNotFound, result, "result should be equal to expected result")
	assert.Error(t, err, "should exists an error")
}

func TestServiceCustomerStoreBulk(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	customerRepository := NewCustomerRepository(db)
	customerService := NewCustomerService(customerRepository)

	for i := 1; i <= 50; i++ {
		mock.ExpectQuery(GetCustomerQuery).WithArgs(i).WillReturnError(ErrorCustomerNotFound)
	}

	mock.ExpectPrepare("INSERT INTO customers")
	for i := 1; i <= 50; i++ {
		mock.ExpectExec("INSERT INTO customers").WillReturnResult(sqlmock.NewResult(int64(i), 1))
	}

	// Act
	result, err := customerService.StoreBulk(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceCustomerGetTotalByCondition(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	customerRepository := NewCustomerRepository(db)
	customerService := NewCustomerService(customerRepository)

	rows := mock.NewRows([]string{"situation", "total"})
	rows.AddRow("Inactivo", 100.0)
	rows.AddRow("Bloqueado", 200.0)
	rows.AddRow("Activo", 300.0)
	mock.ExpectQuery("SELECT customers.situation").WillReturnRows(rows)

	// Act
	result, err := customerService.GetTotalByCondition(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceCustomerGetTotalByConditionError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	customerRepository := NewCustomerRepository(db)
	customerService := NewCustomerService(customerRepository)

	mock.ExpectQuery("SELECT customers.situation").WillReturnError(errors.New("error"))

	// Act
	result, err := customerService.GetTotalByCondition(context.Background())

	// Assert
	assert.Error(t, err, "error should exists")
	assert.Nil(t, result, "result should be nil")
}

func TestServiceCustomerGetCheaperProduct(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	customerRepository := NewCustomerRepository(db)
	customerService := NewCustomerService(customerRepository)

	rows := mock.NewRows([]string{"last_name", "first_name", "price"})
	rows.AddRow("Vend", "testo", 1)
	rows.AddRow("COmp", "testo", 1)
	rows.AddRow("Pep", "testo", 1)
	mock.ExpectQuery("SELECT DISTINCT").WillReturnRows(rows)

	// Act
	result, err := customerService.GetCustomerCheaperProducts(context.Background())

	// Assert
	assert.True(t, len(result) > 0, "result should has more than 0 results")
	assert.Nil(t, err, "error should be nil")
}

func TestServiceCustomerGetCheaperProductError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.Nil(t, err, "error should be nil")
	customerRepository := NewCustomerRepository(db)
	customerService := NewCustomerService(customerRepository)

	mock.ExpectQuery("SELECT DISTINCT").WillReturnError(errors.New("error"))

	// Act
	result, err := customerService.GetCustomerCheaperProducts(context.Background())

	// Assert
	assert.Error(t, err, "error should exists")
	assert.Nil(t, result, "result should be nil")
}
