package customer

import (
	"context"
	"strconv"

	"github.com/matias-ziliotto/HackthonGo/internal/domain"
	"github.com/matias-ziliotto/HackthonGo/pkg/file"
)

var (
	CustomerTxtPath = "../../datos/customers.txt"
)

type CustomerService interface {
	Get(ctx context.Context, id int) (domain.Customer, error)
	StoreBulk(ctx context.Context) ([]domain.Customer, error)
}

func NewCustomerService(pr CustomerRepository) CustomerService {
	return &customerService{
		repository: pr,
	}
}

type customerService struct {
	repository CustomerRepository
}

func (s *customerService) Get(ctx context.Context, id int) (domain.Customer, error) {
	customer, err := s.repository.Get(ctx, id)

	if err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}

func (s *customerService) StoreBulk(ctx context.Context) ([]domain.Customer, error) {
	data, err := file.ReadFile(CustomerTxtPath)

	if err != nil {
		return nil, err
	}

	var customers []domain.Customer

	for _, line := range data {
		id, errId := strconv.Atoi(line[0])
		lastName := line[1]
		firstName := line[2]
		situation := line[3]

		if errId == nil {
			customerAux := domain.Customer{
				Id:        id,
				FirstName: firstName,
				LastName:  lastName,
				Situation: situation,
			}

			customers = append(customers, customerAux)
		}
	}

	customersSaved, err := s.repository.StoreBulk(ctx, customers)

	if err != nil {
		return nil, err
	}

	return customersSaved, nil
}
