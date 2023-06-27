package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
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

// Register handles the customer registration operation
func (s *CustomerService) Register(name, username, password string, phone int) error {
	log.Println("Registering new customer...")

	// Check if username already exists
	log.Println("Checking username availability...")
	_, err := s.repo.GetByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	// Create a new customer
	log.Println("Creating new customer...")
	customer := &models.Customer{
		Name:     name,
		Username: username,
		Password: password,
		Phone:    phone,
	}

	// Save the customer to the repository
	log.Println("Saving customer...")
	err = s.repo.SaveCustomer(customer)
	if err != nil {
		return fmt.Errorf("failed to save customer: %w", err)
	}

	log.Println("Customer registration successful.")

	return nil
}
