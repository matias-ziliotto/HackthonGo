package invoice

import (
	"context"
	"strconv"

	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/pkg/file"
)

var (
	InvoiceTxtPath = "../../datos/invoices.txt"
)

type InvoiceService interface {
	Get(ctx context.Context, id int) (domain.Invoice, error)
	StoreBulk(ctx context.Context) ([]domain.Invoice, error)
	UpdateTotal(ctx context.Context) ([]domain.InvoiceTotalDTO, error)
}

func NewInvoiceService(pr InvoiceRepository) InvoiceService {
	return &invoiceService{
		repository: pr,
	}
}

type invoiceService struct {
	repository InvoiceRepository
}

func (s *invoiceService) Get(ctx context.Context, id int) (domain.Invoice, error) {
	invoice, err := s.repository.Get(ctx, id)

	if err != nil {
		return domain.Invoice{}, err
	}

	return invoice, nil
}

func (s *invoiceService) StoreBulk(ctx context.Context) ([]domain.Invoice, error) {
	data, err := file.ReadFile(InvoiceTxtPath)

	if err != nil {
		return nil, err
	}

	var invoices []domain.Invoice

	for _, line := range data {
		id, errId := strconv.Atoi(line[0])
		datetime := line[1]
		customerId, errCustomerId := strconv.Atoi(line[2])

		if errId == nil && errCustomerId == nil {
			emptyInvoice := domain.Invoice{}
			invoiceAux := domain.Invoice{
				Id:          id,
				Customer_id: customerId,
				Datetime:    datetime,
				Total:       0, // TODO: ver
			}

			// Check if invoice already exists
			invoiceExists, err := s.repository.Get(ctx, invoiceAux.Id)
			if err != nil && invoiceExists == emptyInvoice {
				invoices = append(invoices, invoiceAux)
			}
		}
	}

	if len(invoices) == 0 {
		return invoices, nil
	}

	invoicesSaved, err := s.repository.StoreBulk(ctx, invoices)

	if err != nil {
		return nil, err
	}

	return invoicesSaved, nil
}

func (s *invoiceService) UpdateTotal(ctx context.Context) ([]domain.InvoiceTotalDTO, error) {
	invoicesIds, err := s.repository.GetAllTotalEmpty(ctx)
	if err != nil {
		return nil, err
	}

	if len(invoicesIds) == 0 {
		return nil, nil
	}

	invoicesTotals, err := s.repository.CalculateTotal(ctx, invoicesIds)
	if err != nil {
		return nil, err
	}

	for _, invoice := range invoicesTotals {
		var invoiceAux domain.Invoice
		invoiceAux.Id = invoice.Id
		invoiceAux.Total = invoice.Total

		_, err = s.repository.UpdateTotal(ctx, invoiceAux)
		if err != nil {
			return nil, err
		}
	}

	return invoicesTotals, nil
}
