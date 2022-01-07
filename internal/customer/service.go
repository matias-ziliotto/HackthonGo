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
	GetTotalByCondition(ctx context.Context) ([]domain.CustomerTotalByConditionDTO, error)
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
			emptyCustomer := domain.Customer{}
			customerAux := domain.Customer{
				Id:        id,
				FirstName: firstName,
				LastName:  lastName,
				Situation: situation,
			}

			// Check if customer already exists
			customerExists, err := s.repository.Get(ctx, customerAux.Id)
			if err != nil && customerExists == emptyCustomer {
				customers = append(customers, customerAux)
			}
		}
	}

	if len(customers) == 0 {
		return customers, nil
	}

	customersSaved, err := s.repository.StoreBulk(ctx, customers)

	if err != nil {
		return nil, err
	}

	return customersSaved, nil
}

func (s *customerService) GetTotalByCondition(ctx context.Context) ([]domain.CustomerTotalByConditionDTO, error) {
	customersTotalByContidion, err := s.repository.GetTotalByCondition(ctx)

	if err != nil {
		return nil, err
	}

	return customersTotalByContidion, nil
}
