package customer

import (
	"context"
	"testing"

	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/pkg/database/sql"
	"github.com/stretchr/testify/assert"
)

var customersToStoreAndGet = []domain.Customer{
	{
		Id:        40000,
		FirstName: "Pepe",
		LastName:  "Argento",
		Situation: "Inactivo",
	},
}

var customersToStore = []domain.Customer{
	{
		Id:        1000,
		FirstName: "Pepe",
		LastName:  "Argento",
		Situation: "Inactivo",
	},
	{
		Id:        1001,
		FirstName: "Coki",
		LastName:  "Argento",
		Situation: "Activo",
	},
	{
		Id:        1002,
		FirstName: "Paola",
		LastName:  "Argento",
		Situation: "Bloqueado",
	},
}

var customersToStoreErrorPrepare = []domain.Customer{}

var customersToStoreErrorStore = []domain.Customer{
	{
		Id:        10000,
		FirstName: "Pepe",
		LastName:  "Argento",
		Situation: "Inactivo",
	},
	{
		Id:        10000,
		FirstName: "Coki",
		LastName:  "Argento",
		Situation: "Activo",
	},
}

func TestCustomerGet(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewCustomerRepository(db)

	// Act
	customerStored, _ := repository.StoreBulk(context.Background(), customersToStoreAndGet)
	result, err := repository.Get(context.Background(), customerStored[0].Id)

	// Assert
	assert.Equal(t, customerStored[0], result, "result should be equal customer stored")
	assert.Nil(t, err, "error should be nil")
}

func TestCustomerGetNotFound(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewCustomerRepository(db)

	// Act
	result, err := repository.Get(context.Background(), 99999)

	// Assert
	assert.Error(t, err, "error should exists")
	assert.Equal(t, domain.Customer{}, result, "result should be equal to expected result")
}

func TestCustomerStoreBulk(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewCustomerRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), customersToStore)

	// Assert
	assert.Nil(t, err, "error should be nil")
	assert.True(t, len(result) == 3, "len of result should be equal to 3")
}

func TestCustomerStoreBulkErrorPrepare(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewCustomerRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), customersToStoreErrorPrepare)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}

func TestCustomerStoreBulkError(t *testing.T) {
	// Arrange
	db, err := sql.InitTxSqlDb()
	assert.Nil(t, err, "error should be nil")
	repository := NewCustomerRepository(db)

	// Act
	result, err := repository.StoreBulk(context.Background(), customersToStoreErrorStore)

	// Assert
	assert.Error(t, err, "should exist an error")
	assert.Nil(t, result, "result should be nil")
}
