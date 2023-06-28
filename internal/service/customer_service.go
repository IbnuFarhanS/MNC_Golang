package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
	"github.com/IbnuFarhanS/Golang_MNC/utils"
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

	// Generate token
	token, err := utils.GenerateToken(username, "user")
	if err != nil {
		return false, fmt.Errorf("failed to generate token: %w", err)
	}

	// Save token to customer's data
	err = s.repo.SaveToken(username, token)
	if err != nil {
		return false, fmt.Errorf("failed to save token: %w", err)
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

// Logout handles the customer logout operation
func (s *CustomerService) Logout(tokenString string) error {
	log.Println("Deleting token from customer data...")
	err := s.repo.DeleteToken(tokenString)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	log.Println("Saving deleted token to blacklist...")
	err = s.repo.SaveTokenToBlacklist(tokenString)
	if err != nil {
		return fmt.Errorf("failed to save deleted token to blacklist: %w", err)
	}

	return nil
}

// GetByID retrieves a customer by ID
func (s *CustomerService) GetByID(customerID string) (*models.Customer, error) {
	return s.repo.GetByID(customerID)
}

// GetByUsername retrieves a customer by username
func (s *CustomerService) GetByUsername(username string) (*models.Customer, error) {
	return s.repo.GetByUsername(username)
}

// SaveCustomer saves a new customer
func (s *CustomerService) SaveCustomer(customer *models.Customer) error {
	return s.repo.SaveCustomer(customer)
}

// SaveToken saves the token to the customer data
func (s *CustomerService) SaveToken(username, token string) error {
	return s.repo.SaveToken(username, token)
}

// DeleteToken deletes the token from the customer data and saves it to the blacklist
func (s *CustomerService) DeleteToken(tokenString string) error {
	return s.repo.DeleteToken(tokenString)
}

// SaveTokenToBlacklist saves the token to the blacklist
func (s *CustomerService) SaveTokenToBlacklist(token string) error {
	return s.repo.SaveTokenToBlacklist(token)
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *CustomerService) IsTokenBlacklisted(token string) (bool, error) {
	return s.repo.IsTokenBlacklisted(token)
}

// SaveToFile saves the customer data to the file
func (s *CustomerService) SaveToFile() error {
	return s.repo.SaveToFile()
}
