package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
	"github.com/IbnuFarhanS/Golang_MNC/utils"
)

// CustomerRepository interface provides methods to interact with customer data
type CustomerRepository interface {
	GetByUsername(username string) (*models.Customer, error)
	GetByID(customerID string) (*models.Customer, error)
	SaveCustomer(customer *models.Customer) error
	DeleteToken(tokenString string) error
}

// InMemoryCustomerRepository implements CustomerRepository interface using in-memory data
type InMemoryCustomerRepository struct {
	customers       []*models.Customer
	customerCounter int
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

// SaveCustomer saves a new customer
func (r *InMemoryCustomerRepository) SaveCustomer(customer *models.Customer) error {
	// Increment customer counter
	r.customerCounter++

	// Set customer ID
	customer.ID = strconv.Itoa(r.customerCounter)
	// Hash password before saving
	hashedPassword, err := utils.GenerateHash(customer.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	customer.Password = hashedPassword

	r.customers = append(r.customers, customer)

	// Save the updated data to the file
	err = r.saveToFile()
	if err != nil {
		return fmt.Errorf("failed to save customer data: %v", err)
	}

	return nil
}

// saveToFile saves the updated customer data to the file
func (r *InMemoryCustomerRepository) saveToFile() error {
	data, err := json.Marshal(r.customers)
	if err != nil {
		return fmt.Errorf("failed to marshal customer data: %v", err)
	}

	err = ioutil.WriteFile("json/customers.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write customer data to file: %v", err)
	}

	return nil
}

// DeleteToken menghapus token dari data pelanggan
func (r *InMemoryCustomerRepository) DeleteToken(tokenString string) error {
	// Baca file .json ke dalam memori
	filePath := "json/customers.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read customer data: %v", err)
	}

	// Unmarshal data ke dalam slice pelanggan yang sesuai
	var customers []*models.Customer
	err = json.Unmarshal(data, &customers)
	if err != nil {
		return fmt.Errorf("failed to unmarshal customer data: %v", err)
	}

	// Cari pelanggan yang memiliki token yang ingin dihapus
	for _, customer := range customers {
		// Jika token ditemukan pada pelanggan
		if customer.Token == tokenString {
			// Hapus token
			customer.Token = ""
			break
		}
	}

	// Marshal kembali data pelanggan ke dalam format .json
	updatedData, err := json.Marshal(customers)
	if err != nil {
		return fmt.Errorf("failed to marshal customer data: %v", err)
	}

	// Tulis kembali file .json dengan data yang diperbarui
	err = ioutil.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write customer data to file: %v", err)
	}

	return nil
}
