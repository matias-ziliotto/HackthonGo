package sale

import (
	"context"
	"strconv"

	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/pkg/file"
)

var (
	SaleTxtPath = "../../datos/sales.txt"
)

type SaleService interface {
	Get(ctx context.Context, id int) (domain.Sale, error)
	StoreBulk(ctx context.Context) ([]domain.Sale, error)
}

func NewSaleService(pr SaleRepository) SaleService {
	return &saleService{
		repository: pr,
	}
}

type saleService struct {
	repository SaleRepository
}

func (s *saleService) Get(ctx context.Context, id int) (domain.Sale, error) {
	sale, err := s.repository.Get(ctx, id)

	if err != nil {
		return domain.Sale{}, err
	}

	return sale, nil
}

func (s *saleService) StoreBulk(ctx context.Context) ([]domain.Sale, error) {
	data, err := file.ReadFile(SaleTxtPath)

	if err != nil {
		return nil, err
	}

	var sales []domain.Sale

	for _, line := range data {
		id, errId := strconv.Atoi(line[0])
		productId, errProductId := strconv.Atoi(line[1])
		invoiceId, errInvoiceId := strconv.Atoi(line[2])
		quantity, errQuantity := strconv.ParseFloat(line[3], 64)

		if errId == nil && errInvoiceId == nil && errProductId == nil && errQuantity == nil {
			emptySale := domain.Sale{}
			saleAux := domain.Sale{
				Id:         id,
				Invoice_id: invoiceId,
				Product_id: productId,
				Quantity:   quantity,
			}

			// Check if customer already exists
			saleExists, err := s.repository.Get(ctx, saleAux.Id)
			if err == nil && saleExists == emptySale {
				sales = append(sales, saleAux)
			}
		}
	}

	salesSaved, err := s.repository.StoreBulk(ctx, sales)

	if err != nil {
		return nil, err
	}

	return salesSaved, nil
}
