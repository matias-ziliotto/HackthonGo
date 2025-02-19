package customer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/matias-ziliotto/HackthonGo/internal/domain"
)

var (
	// Db queries & statements
	GetCustomerQuery                  = "SELECT id, first_name, last_name, situation FROM customers WHERE id = ?"
	GetCustomersTotalByConditionQuery = "SELECT customers.situation, ROUND(SUM(invoices.total), 2) FROM customers INNER JOIN invoices ON invoices.customer_id = customers.id GROUP BY customers.situation;"
	GetCustomersCheaperProductsQuery  = "SELECT DISTINCT(customers.last_name), customers.first_name, products.price FROM customers INNER JOIN invoices ON invoices.customer_id = customers.id INNER JOIN sales ON sales.invoice_id = invoices.id INNER JOIN products ON sales.product_id = products.id ORDER BY products.price ASC, customers.last_name ASC LIMIT 5;"
	StoreCustomerStatement            = "INSERT INTO customers(first_name, last_name, situation) VALUES(?, ?, ?)"

	// Errors
	ErrorCustomerNotFound              = errors.New("customer not found")
	ErrorCustomerPrepareStoreStatement = errors.New("can not prepare store statement")
	ErrorCustomerExecStoreStatement    = errors.New("error executing store statement")
)

type CustomerRepository interface {
	Get(ctx context.Context, id int) (domain.Customer, error)
	GetTotalByCondition(ctx context.Context) ([]domain.CustomerTotalByConditionDTO, error)
	GetCustomerCheaperProducts(ctx context.Context) ([]domain.CustomerCheaperProductDTO, error)
	StoreBulk(ctx context.Context, customers []domain.Customer) ([]domain.Customer, error)
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

type customerRepository struct {
	db *sql.DB
}

func (r *customerRepository) Get(ctx context.Context, id int) (domain.Customer, error) {
	var customer domain.Customer
	err := r.db.QueryRowContext(ctx, GetCustomerQuery, id).Scan(&customer.Id, &customer.FirstName, &customer.LastName, &customer.Situation)

	if err != nil {
		return domain.Customer{}, ErrorCustomerNotFound
	}

	return customer, nil
}

func (r *customerRepository) StoreBulk(ctx context.Context, customers []domain.Customer) ([]domain.Customer, error) {
	valueStrings := make([]string, 0, len(customers))
	valueArgs := make([]interface{}, 0, len(customers)*4)

	for _, customer := range customers {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, customer.Id)
		valueArgs = append(valueArgs, customer.FirstName)
		valueArgs = append(valueArgs, customer.LastName)
		valueArgs = append(valueArgs, customer.Situation)
	}

	stmtString := fmt.Sprintf("INSERT INTO customers (id, first_name, last_name, situation) VALUES %s", strings.Join(valueStrings, ","))
	stmt, err := r.db.PrepareContext(ctx, stmtString)
	if err != nil {
		return nil, ErrorCustomerPrepareStoreStatement
	}

	defer stmt.Close()

	result, err := r.db.Exec(stmtString, valueArgs...)

	if err != nil {
		return nil, ErrorCustomerExecStoreStatement
	}

	_, err = result.RowsAffected()
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *customerRepository) GetTotalByCondition(ctx context.Context) ([]domain.CustomerTotalByConditionDTO, error) {
	rows, err := r.db.QueryContext(ctx, GetCustomersTotalByConditionQuery)

	if err != nil {
		return nil, err
	}

	var customerTotalByConditions []domain.CustomerTotalByConditionDTO

	for rows.Next() {
		var customerTotalByCondition domain.CustomerTotalByConditionDTO
		err = rows.Scan(&customerTotalByCondition.Situation, &customerTotalByCondition.Total)
		if err != nil {
			return nil, err
		}

		customerTotalByConditions = append(customerTotalByConditions, customerTotalByCondition)
	}

	return customerTotalByConditions, nil
}

func (r *customerRepository) GetCustomerCheaperProducts(ctx context.Context) ([]domain.CustomerCheaperProductDTO, error) {
	rows, err := r.db.QueryContext(ctx, GetCustomersCheaperProductsQuery)

	if err != nil {
		return nil, err
	}

	var customerCheaperProducts []domain.CustomerCheaperProductDTO

	for rows.Next() {
		var customerCheaperProduct domain.CustomerCheaperProductDTO
		var price float64
		err = rows.Scan(&customerCheaperProduct.LastName, &customerCheaperProduct.FirstName, &price)
		if err != nil {
			return nil, err
		}

		customerCheaperProducts = append(customerCheaperProducts, customerCheaperProduct)
	}

	return customerCheaperProducts, nil
}
