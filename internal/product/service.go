package product

import (
	"context"
	"strconv"

	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/pkg/file"
)

var (
	ProductTxtPath = "../../datos/products.txt"
)

type ProductService interface {
	Get(ctx context.Context, id int) (domain.Product, error)
	StoreBulk(ctx context.Context) ([]domain.Product, error)
	GetProductsMostSelled(ctx context.Context) ([]domain.ProductMostSelledDTO, error)
}

func NewProductService(pr ProductRepository) ProductService {
	return &productService{
		repository: pr,
	}
}

type productService struct {
	repository ProductRepository
}

func (s *productService) Get(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.repository.Get(ctx, id)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

func (s *productService) StoreBulk(ctx context.Context) ([]domain.Product, error) {
	data, err := file.ReadFile(ProductTxtPath)

	if err != nil {
		return nil, err
	}

	var products []domain.Product

	for _, line := range data {
		id, errId := strconv.Atoi(line[0])
		description := line[1]
		price, errPrice := strconv.ParseFloat(line[2], 64)

		if errId == nil && errPrice == nil {
			emptyProduct := domain.Product{}
			productAux := domain.Product{
				Id:          id,
				Description: description,
				Price:       price,
			}

			// Check if product already exists
			productExists, err := s.repository.Get(ctx, productAux.Id)
			if err != nil && productExists == emptyProduct {
				products = append(products, productAux)
			}
		}
	}

	if len(products) == 0 {
		return products, nil
	}

	productsSaved, err := s.repository.StoreBulk(ctx, products)

	if err != nil {
		return nil, err
	}

	return productsSaved, nil
}

func (s *productService) GetProductsMostSelled(ctx context.Context) ([]domain.ProductMostSelledDTO, error) {
	productsMostSelled, err := s.repository.ProductsMostSelled(ctx)

	if err != nil {
		return nil, err
	}

	return productsMostSelled, nil
}
