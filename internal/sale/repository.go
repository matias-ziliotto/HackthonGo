package sale

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
	GetSaleQuery       = "SELECT id, invoice_id, product_id, quantity FROM sales WHERE id = ?"
	StoreSaleStatement = "INSERT INTO sales(invoice_id, product_id, quantity) VALUES(?, ?, ?)"

	// Errors
	ErrorSaleNotFound              = errors.New("sale not found")
	ErrorSalePrepareStoreStatement = errors.New("can not prepare store statement")
	ErrorSaleExecStoreStatement    = errors.New("error executing store statement")
)

type SaleRepository interface {
	Get(ctx context.Context, id int) (domain.Sale, error)
	StoreBulk(ctx context.Context, sales []domain.Sale) ([]domain.Sale, error)
}

func NewSaleRepository(db *sql.DB) SaleRepository {
	return &saleRepository{
		db: db,
	}
}

type saleRepository struct {
	db *sql.DB
}

func (r *saleRepository) Get(ctx context.Context, id int) (domain.Sale, error) {
	var sale domain.Sale
	err := r.db.QueryRowContext(ctx, GetSaleQuery, id).Scan(&sale.Id, &sale.Invoice_id, &sale.Product_id, &sale.Quantity)

	if err != nil {
		return domain.Sale{}, ErrorSaleNotFound
	}

	return sale, nil
}

func (r *saleRepository) StoreBulk(ctx context.Context, sales []domain.Sale) ([]domain.Sale, error) {
	valueStrings := make([]string, 0, len(sales))
	valueArgs := make([]interface{}, 0, len(sales)*4)

	for _, sale := range sales {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, sale.Id)
		valueArgs = append(valueArgs, sale.Invoice_id)
		valueArgs = append(valueArgs, sale.Product_id)
		valueArgs = append(valueArgs, sale.Quantity)
	}

	stmtString := fmt.Sprintf("INSERT INTO sales (id, invoice_id, product_id, quantity) VALUES %s", strings.Join(valueStrings, ","))
	stmt, err := r.db.PrepareContext(ctx, stmtString)
	if err != nil {
		return nil, ErrorSalePrepareStoreStatement
	}

	defer stmt.Close()

	result, err := r.db.Exec(stmtString, valueArgs...)

	if err != nil {
		return nil, ErrorSaleExecStoreStatement
	}

	_, err = result.RowsAffected()
	if err != nil {
		return nil, err
	}

	return sales, nil
}
