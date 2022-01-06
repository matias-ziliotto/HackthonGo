package invoice

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
	GetInvoiceQuery       = "SELECT id, customer_id, datetime, total FROM invoices WHERE id = ?"
	StoreInvoiceStatement = "INSERT INTO invoices(customer_id, datetime, total) VALUES(?, ?, ?)"

	// Errors
	ErrorInvoiceNotFound              = errors.New("invoice not found")
	ErrorInvoicePrepareStoreStatement = errors.New("can not prepare store statement")
	ErrorInvoiceExecStoreStatement    = errors.New("error executing store statement")
)

type InvoiceRepository interface {
	Get(ctx context.Context, id int) (domain.Invoice, error)
	Store(ctx context.Context, invoice domain.Invoice) (domain.Invoice, error)
	StoreBulk(ctx context.Context, invoices []domain.Invoice) ([]domain.Invoice, error)
}

func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &invoiceRepository{
		db: db,
	}
}

type invoiceRepository struct {
	db *sql.DB
}

func (r *invoiceRepository) Get(ctx context.Context, id int) (domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.QueryRowContext(ctx, GetInvoiceQuery, id).Scan(&invoice)

	if err != nil {
		return domain.Invoice{}, ErrorInvoiceNotFound
	}

	return invoice, nil
}

func (r *invoiceRepository) Store(ctx context.Context, invoice domain.Invoice) (domain.Invoice, error) {
	stmt, err := r.db.PrepareContext(ctx, StoreInvoiceStatement)

	if err != nil {
		return domain.Invoice{}, ErrorInvoicePrepareStoreStatement
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, invoice.Customer_id, invoice.Datetime, invoice.Total)

	if err != nil {
		return domain.Invoice{}, ErrorInvoiceExecStoreStatement
	}

	id, err := result.LastInsertId()

	if err != nil {
		return domain.Invoice{}, err
	}

	invoice.Id = int(id)

	return invoice, nil
}

func (r *invoiceRepository) StoreBulk(ctx context.Context, invoices []domain.Invoice) ([]domain.Invoice, error) {
	valueStrings := make([]string, 0, len(invoices))
	valueArgs := make([]interface{}, 0, len(invoices)*3)

	for _, invoice := range invoices {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, invoice.Id)
		valueArgs = append(valueArgs, invoice.Customer_id)
		valueArgs = append(valueArgs, invoice.Datetime)
		valueArgs = append(valueArgs, invoice.Total)
	}

	stmtString := fmt.Sprintf("INSERT INTO invoices (id, customer_id, datetime, total) VALUES %s", strings.Join(valueStrings, ","))
	stmt, err := r.db.PrepareContext(ctx, stmtString)
	if err != nil {
		return nil, ErrorInvoicePrepareStoreStatement
	}

	defer stmt.Close()

	result, err := r.db.Exec(stmtString, valueArgs...)

	if err != nil {
		return nil, ErrorInvoiceExecStoreStatement
	}

	_, err = result.RowsAffected()
	if err != nil {
		return nil, err
	}

	return invoices, nil
}
