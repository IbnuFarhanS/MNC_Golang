package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
	"github.com/IbnuFarhanS/Golang_MNC/utils"
)

// CustomerRepository interface provides methods to interact with customer data
type CustomerRepository interface {
	GetByUsername(username string) (*models.Customer, error)
	GetByID(customerID string) (*models.Customer, error)
}

// InMemoryCustomerRepository implements CustomerRepository interface using in-memory data
type InMemoryCustomerRepository struct {
	customers []*models.Customer
}

// NewInMemoryCustomerRepository creates a new instance of InMemoryCustomerRepository
func NewInMemoryCustomerRepository(filePath string) (*InMemoryCustomerRepository, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read customer data: %v", err)
	}

	var customers []*models.Customer
	err = json.Unmarshal(data, &customers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer data: %v", err)
	}

	// Hash password for each customer
	for _, customer := range customers {
		hashedPassword, err := utils.GenerateHash(customer.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		customer.Password = hashedPassword
	}

	return &InMemoryCustomerRepository{
		customers: customers,
	}, nil
}

// GetByUsername retrieves a customer by username
func (r *InMemoryCustomerRepository) GetByUsername(username string) (*models.Customer, error) {
	for _, c := range r.customers {
		if c.Username == username {
			return c, nil
		}
	}
	return nil, fmt.Errorf("customer not found")
}

// GetByID retrieves a customer by ID
func (r *InMemoryCustomerRepository) GetByID(customerID string) (*models.Customer, error) {
	for _, c := range r.customers {
		if c.ID == customerID {
			return c, nil
		}
	}
	return nil, fmt.Errorf("customer not found")
}
