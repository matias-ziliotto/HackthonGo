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
	GetAllTotalEmptyInvoiceQuery = "SELECT id FROM invoices WHERE total = 0"
	GetInvoiceQuery              = "SELECT id, customer_id, datetime, total FROM invoices WHERE id = ?"
	CalculateTotalInvoiceQuery   = "SELECT DISTINCT(invoices.id), SUM(calc.total) as total FROM invoices INNER JOIN ( SELECT sales.invoice_id,  SUM(products.price) * sales.quantity AS total FROM sales INNER JOIN products ON products.id = sales.product_id WHERE sales.invoice_id IN replace_with_invoices_ids GROUP BY sales.id ) calc ON calc.invoice_id = invoices.id GROUP BY calc.invoice_id;"
	StoreInvoiceStatement        = "INSERT INTO invoices(customer_id, datetime, total) VALUES(?, ?, ?)"
	UpdateInvoiceStatement       = "UPDATE invoices SET total = ? WHERE id = ?"

	// Errors
	ErrorInvoiceNotFound               = errors.New("invoice not found")
	ErrorInvoicePrepareStoreStatement  = errors.New("can not prepare store statement")
	ErrorInvoiceExecStoreStatement     = errors.New("error executing store statement")
	ErrorInvoicePrepareUpdateStatement = errors.New("can not prepare update statement")
	ErrorInvoiceExecUpdateStatement    = errors.New("error executing update statement")
)

type InvoiceRepository interface {
	GetAllTotalEmpty(ctx context.Context) ([]int, error)
	Get(ctx context.Context, id int) (domain.Invoice, error)
	StoreBulk(ctx context.Context, invoices []domain.Invoice) ([]domain.Invoice, error)
	UpdateTotal(ctx context.Context, invoice domain.Invoice) (domain.Invoice, error)
	CalculateTotal(ctx context.Context, ids []int) ([]domain.InvoiceTotalDTO, error)
}

func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &invoiceRepository{
		db: db,
	}
}

type invoiceRepository struct {
	db *sql.DB
}

func (r *invoiceRepository) GetAllTotalEmpty(ctx context.Context) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, GetAllTotalEmptyInvoiceQuery)

	if err != nil {
		return nil, err
	}

	var invoicesIds []int

	for rows.Next() {
		var invoiceId int
		err = rows.Scan(&invoiceId)
		if err != nil {
			return nil, err
		}

		invoicesIds = append(invoicesIds, invoiceId)
	}

	return invoicesIds, nil
}

func (r *invoiceRepository) Get(ctx context.Context, id int) (domain.Invoice, error) {
	var invoice domain.Invoice
	err := r.db.QueryRowContext(ctx, GetInvoiceQuery, id).Scan(&invoice.Id, &invoice.Customer_id, &invoice.Datetime, &invoice.Total)

	if err != nil {
		return domain.Invoice{}, ErrorInvoiceNotFound
	}

	return invoice, nil
}

func (r *invoiceRepository) StoreBulk(ctx context.Context, invoices []domain.Invoice) ([]domain.Invoice, error) {
	valueStrings := make([]string, 0, len(invoices))
	valueArgs := make([]interface{}, 0, len(invoices)*4)

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

func (r *invoiceRepository) UpdateTotal(ctx context.Context, invoice domain.Invoice) (domain.Invoice, error) {
	stmt, err := r.db.PrepareContext(ctx, UpdateInvoiceStatement)

	if err != nil {
		return domain.Invoice{}, ErrorInvoicePrepareUpdateStatement
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, invoice.Total, invoice.Id)

	if err != nil {
		return domain.Invoice{}, ErrorInvoiceExecUpdateStatement
	}

	_, err = result.RowsAffected()

	if err != nil {
		return domain.Invoice{}, err
	}

	return invoice, nil
}

func (r *invoiceRepository) CalculateTotal(ctx context.Context, ids []int) ([]domain.InvoiceTotalDTO, error) {
	// replace replace_with_invoices_ids with (id1, id2, id3) in query
	query := CalculateTotalInvoiceQuery
	query = strings.ReplaceAll(query, "replace_with_invoices_ids", "("+strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]")+")")

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	var invoiceTotals []domain.InvoiceTotalDTO
	for rows.Next() {
		var invoiceAux domain.InvoiceTotalDTO

		err = rows.Scan(&invoiceAux.Id, &invoiceAux.Total)
		if err != nil {
			return nil, err
		}

		invoiceTotals = append(invoiceTotals, invoiceAux)
	}

	return invoiceTotals, nil
}
