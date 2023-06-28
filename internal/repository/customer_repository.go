package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
	"github.com/IbnuFarhanS/Golang_MNC/utils"
)

type BlacklistToken struct {
	Token string `json:"token"`
}

// CustomerRepository interface provides methods to interact with customer data
type CustomerRepository interface {
	GetByUsername(username string) (*models.Customer, error)
	GetByID(customerID string) (*models.Customer, error)
	SaveCustomer(customer *models.Customer) error
	SaveToFile() error
	SaveToken(username, token string) error
	DeleteToken(tokenString string) error
	SaveTokenToBlacklist(token string) error
	IsTokenBlacklisted(token string) (bool, error)
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
	err = r.SaveToFile()
	if err != nil {
		return fmt.Errorf("failed to save customer data: %v", err)
	}

	return nil
}

// SaveToken saves the token to the customer data
func (r *InMemoryCustomerRepository) SaveToken(username, token string) error {
	for _, customer := range r.customers {
		if customer.Username == username {
			customer.Token = token
			break
		}
	}

	// Save the updated data to the file
	err := r.SaveToFile()
	if err != nil {
		return fmt.Errorf("failed to save customer data: %v", err)
	}

	return nil
}

// saveToFile saves the updated customer data to the file
func (r *InMemoryCustomerRepository) SaveToFile() error {
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

// SaveTokenToBlacklist saves the token to the blacklist
func (r *InMemoryCustomerRepository) SaveTokenToBlacklist(token string) error {
	// Remove the "Bearer " prefix from the token
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	// Read the existing blacklist tokens
	filePath := "json/blacklist_token.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read blacklist token data: %v", err)
	}

	var tokens []BlacklistToken
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return fmt.Errorf("failed to unmarshal blacklist token data: %v", err)
	}

	// Add the new token to the blacklist
	tokens = append(tokens, BlacklistToken{Token: token})

	// Marshal the updated data back to JSON
	updatedData, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal blacklist token data: %v", err)
	}

	// Write the updated data back to the file
	err = ioutil.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write blacklist token data to file: %v", err)
	}

	return nil
}

// DeleteToken deletes the token from the customer data and saves it to the blacklist
func (r *InMemoryCustomerRepository) DeleteToken(tokenString string) error {
	// Find the customer with the token to be deleted
	for _, customer := range r.customers {
		if customer.Token == tokenString {
			// Save the token to the blacklist
			err := r.SaveTokenToBlacklist(tokenString)
			if err != nil {
				return fmt.Errorf("failed to save token to blacklist: %v", err)
			}

			// Set the customer's token to an empty string
			customer.Token = ""
			break
		}
	}

	// Save the updated customer data to the file
	err := r.SaveToFile()
	if err != nil {
		return fmt.Errorf("failed to save customer data: %v", err)
	}

	return nil
}

// IsTokenBlacklisted checks if a token is blacklisted
func (r *InMemoryCustomerRepository) IsTokenBlacklisted(token string) (bool, error) {
	// Read the existing blacklist tokens
	filePath := "json/blacklist_token.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to read blacklist token data: %v", err)
	}

	var tokens []BlacklistToken
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal blacklist token data: %v", err)
	}

	// Check if the token is blacklisted
	for _, t := range tokens {
		if strings.HasPrefix(token, t.Token) {
			return true, nil
		}
	}

	return false, nil
}
