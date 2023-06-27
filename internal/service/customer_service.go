package service

import (
	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// CustomerService handles customer-related operations
type CustomerService struct {
	repo repository.CustomerRepository
}

// NewCustomerService creates a new instance of CustomerService
func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

// Login handles the login operation
func (s *CustomerService) Login(username, password string) (bool, error) {
	customer, err := s.repo.GetByUsername(username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}
