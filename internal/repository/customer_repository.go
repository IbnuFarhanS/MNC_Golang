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

// Mendefinisikan tipe data BlacklistToken yang akan digunakan untuk menyimpan token yang akan dimasukkan ke dalam daftar hitam (blacklist)
type BlacklistToken struct {
	Token string `json:"token"`
}

// Mendefinisikan interface CustomerRepository yang menyediakan method-method
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

// Mendefinisikan tipe data InMemoryCustomerRepository yang merupakan implementasi dari interface CustomerRepository
type InMemoryCustomerRepository struct {
	customers       []*models.Customer
	customerCounter int
}

// Mendefinisikan fungsi NewInMemoryCustomerRepository yang digunakan untuk membuat instance baru dari InMemoryCustomerRepository
func NewInMemoryCustomerRepository(filePath string) (*InMemoryCustomerRepository, error) {
	// Membaca file yang berisi data pelanggan
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read customer data: %v", err)
	}

	// Mendekode data JSON menjadi slice of Customer
	var customers []*models.Customer
	err = json.Unmarshal(data, &customers)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer data: %v", err)
	}

	// Melakukan hashing pada password pelanggan
	for _, customer := range customers {
		hashedPassword, err := utils.GenerateHash(customer.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		// Mengganti password dengan password yang sudah di-hash
		customer.Password = hashedPassword
	}

	// Menginisialisasi slice customers pada InMemoryCustomerRepository
	return &InMemoryCustomerRepository{
		customers: customers,
	}, nil
}

// Implementasi method GetByUsername yang mengambil data pelanggan berdasarkan username
func (r *InMemoryCustomerRepository) GetByUsername(username string) (*models.Customer, error) {
	for _, c := range r.customers {
		if c.Username == username {
			return c, nil
		}
	}
	return nil, fmt.Errorf("customer not found")
}

// Implementasi method GetByID yang mengambil data pelanggan berdasarkan ID
func (r *InMemoryCustomerRepository) GetByID(customerID string) (*models.Customer, error) {
	for _, c := range r.customers {
		if c.ID == customerID {
			return c, nil
		}
	}
	return nil, fmt.Errorf("customer not found")
}

// Implementasi method SaveCustomer untuk menyimpan data pelanggan baru
func (r *InMemoryCustomerRepository) SaveCustomer(customer *models.Customer) error {
	r.customerCounter++ // Increment customer counter

	// Mengatur ID pelanggan
	customer.ID = strconv.Itoa(r.customerCounter)
	// Melakukan hashing pada password pelanggan
	hashedPassword, err := utils.GenerateHash(customer.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	// Mengganti password dengan password yang sudah di-hash
	customer.Password = hashedPassword

	// Menambahkan pelanggan baru ke dalam slice customers
	r.customers = append(r.customers, customer)

	// Menyimpan data yang sudah diupdate ke dalam file
	err = r.SaveToFile()
	if err != nil {
		return fmt.Errorf("failed to save customer data: %v", err)
	}

	return nil
}

// Implementasi method SaveToken untuk menyimpan token ke data pelanggan
func (r *InMemoryCustomerRepository) SaveToken(username, token string) error {
	// Menyimpan token ke dalam data pelanggan
	for _, customer := range r.customers {
		if customer.Username == username {
			customer.Token = token
			break
		}
	}

	// Menyimpan data yang sudah diupdate ke dalam file
	err := r.SaveToFile()
	if err != nil {
		return fmt.Errorf("failed to save customer data: %v", err)
	}

	return nil
}

// Implementasi method SaveToFile untuk menyimpan data pelanggan ke file
func (r *InMemoryCustomerRepository) SaveToFile() error {
	// Melakukan encoding data pelanggan menjadi JSON
	data, err := json.Marshal(r.customers)
	if err != nil {
		return fmt.Errorf("failed to marshal customer data: %v", err)
	}

	// Menulis data JSON ke dalam file
	err = ioutil.WriteFile("json/customers.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write customer data to file: %v", err)
	}

	return nil
}

// Implementasi method SaveTokenToBlacklist untuk menyimpan token ke dalam daftar hitam (blacklist)
func (r *InMemoryCustomerRepository) SaveTokenToBlacklist(token string) error {
	// Menghapus prefix "Bearer " dari token
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	// Membaca file yang berisi data daftar hitam (blacklist) token
	filePath := "json/blacklist_token.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read blacklist token data: %v", err)
	}

	// Mendekode data JSON menjadi slice of BlacklistToken
	var tokens []BlacklistToken
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return fmt.Errorf("failed to unmarshal blacklist token data: %v", err)
	}

	// Menambahkan token baru ke dalam daftar hitam (blacklist)
	tokens = append(tokens, BlacklistToken{Token: token})

	// Melakukan encoding data daftar hitam (blacklist) menjadi JSON
	updatedData, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to marshal blacklist token data: %v", err)
	}

	// Menulis data JSON yang sudah diupdate ke dalam file
	err = ioutil.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write blacklist token data to file: %v", err)
	}

	return nil
}

// Implementasi method DeleteToken untuk menghapus token dari data pelanggan dan menyimpannya ke dalam daftar hitam (blacklist)
func (r *InMemoryCustomerRepository) DeleteToken(tokenString string) error {
	for _, customer := range r.customers {
		if customer.Token == tokenString {
			// Menyimpan token ke dalam daftar hitam (blacklist)
			err := r.SaveTokenToBlacklist(tokenString)
			if err != nil {
				return fmt.Errorf("failed to save token to blacklist: %v", err)
			}

			// Menghapus token dari data pelanggan
			customer.Token = ""
			break
		}
	}

	// Menyimpan data pelanggan yang sudah diupdate ke dalam file
	err := r.SaveToFile()
	if err != nil {
		return fmt.Errorf("failed to save customer data: %v", err)
	}

	return nil
}

// Implementasi method IsTokenBlacklisted untuk memeriksa apakah token berada dalam daftar hitam (blacklist)
func (r *InMemoryCustomerRepository) IsTokenBlacklisted(token string) (bool, error) {
	// Membaca file yang berisi data daftar hitam (blacklist) token
	filePath := "json/blacklist_token.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to read blacklist token data: %v", err)
	}

	// Mendekode data JSON menjadi slice of BlacklistToken
	var tokens []BlacklistToken
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal blacklist token data: %v", err)
	}

	// Memeriksa apakah token berada dalam daftar hitam (blacklist)
	for _, t := range tokens {
		if strings.HasPrefix(token, t.Token) {
			return true, nil
		}
	}

	return false, nil
}
