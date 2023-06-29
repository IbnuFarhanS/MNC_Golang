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

// CustomerServic untuke menangani operasi terkait pelanggan
type CustomerService struct {
	repo repository.CustomerRepository
}

// NewCustomerService untuk membuat instance baru dari CustomerService
func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

// Login untuk menangani operasi login
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
		return false, fmt.Errorf("gagal menghasilkan token: %w", err)
	}

	// Simpan token ke data pelanggan
	err = s.repo.SaveToken(username, token)
	if err != nil {
		return false, fmt.Errorf("gagal menyimpan token: %w", err)
	}

	return true, nil
}

// Register untuk menangani operasi registrasi pelanggan
func (s *CustomerService) Register(name, username, password string, phone int) error {
	log.Println("Mendaftarkan pelanggan baru...")

	// Periksa apakah username sudah ada
	log.Println("Memeriksa ketersediaan username...")
	_, err := s.repo.GetByUsername(username)
	if err == nil {
		return errors.New("username sudah digunakan")
	}

	// Buat pelanggan baru
	log.Println("Membuat pelanggan baru...")
	customer := &models.Customer{
		Name:     name,
		Username: username,
		Password: password,
		Phone:    phone,
	}

	// Simpan pelanggan ke repositori
	log.Println("Menyimpan pelanggan...")
	err = s.repo.SaveCustomer(customer)
	if err != nil {
		return fmt.Errorf("gagal menyimpan pelanggan: %w", err)
	}

	log.Println("Registrasi pelanggan berhasil.")

	return nil
}

// Logout untuk menangani operasi logout pelanggan
func (s *CustomerService) Logout(tokenString string) error {
	log.Println("Menghapus token dari data pelanggan...")
	err := s.repo.DeleteToken(tokenString)
	if err != nil {
		return fmt.Errorf("gagal menghapus token: %w", err)
	}

	log.Println("Menyimpan token yang dihapus ke daftar hitam...")
	err = s.repo.SaveTokenToBlacklist(tokenString)
	if err != nil {
		return fmt.Errorf("gagal menyimpan token yang dihapus ke daftar hitam: %w", err)
	}

	return nil
}

// GetByID mengambil untuk pelanggan berdasarkan ID
func (s *CustomerService) GetByID(customerID string) (*models.Customer, error) {
	return s.repo.GetByID(customerID)
}

// GetByUsername untuk mengambil pelanggan berdasarkan username
func (s *CustomerService) GetByUsername(username string) (*models.Customer, error) {
	return s.repo.GetByUsername(username)
}

// SaveCustomer untuk menyimpan pelanggan baru
func (s *CustomerService) SaveCustomer(customer *models.Customer) error {
	return s.repo.SaveCustomer(customer)
}

// SaveToken untuk menyimpan token ke data pelanggan
func (s *CustomerService) SaveToken(username, token string) error {
	return s.repo.SaveToken(username, token)
}

// DeleteToken untuk menghapus token dari data pelanggan dan menyimpannya ke daftar hitam
func (s *CustomerService) DeleteToken(tokenString string) error {
	return s.repo.DeleteToken(tokenString)
}

// SaveTokenToBlacklist untuk menyimpan token ke daftar hitam
func (s *CustomerService) SaveTokenToBlacklist(token string) error {
	return s.repo.SaveTokenToBlacklist(token)
}

// IsTokenBlacklisted untuk memeriksa apakah token ada dalam daftar hitam
func (s *CustomerService) IsTokenBlacklisted(token string) (bool, error) {
	return s.repo.IsTokenBlacklisted(token)
}

// SaveToFile untuk menyimpan data pelanggan ke file
func (s *CustomerService) SaveToFile() error {
	return s.repo.SaveToFile()
}
