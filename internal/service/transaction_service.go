package service

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
)

// TransactionService menangani operasi terkait transaksi
type TransactionService struct {
	transactionRepository *repository.TransactionRepository
	customerRepository    repository.CustomerRepository
	merchantRepository    repository.MerchantRepository
}

func NewTransactionService(transactionRepository *repository.TransactionRepository, customerRepository repository.CustomerRepository, merchantRepository repository.MerchantRepository) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		customerRepository:    customerRepository,
		merchantRepository:    merchantRepository,
	}
}

func (s *TransactionService) ProcessTransaction(customerID string, merchantID string, amount float64) error {
	log.Println("Memproses transaksi...")

	// Validasi customer ID
	log.Println("Memvalidasi customer ID...")
	_, err := s.customerRepository.GetByID(customerID)
	if err != nil {
		return errors.New("ID customer tidak valid")
	}

	// Validasi merchant ID
	_, err = s.merchantRepository.GetByID(merchantID)
	if err != nil {
		return errors.New("ID merchant tidak valid")
	}

	// Validasi jumlah transaksi
	log.Println("Memvalidasi jumlah transaksi...")
	if amount <= 0 {
		return errors.New("jumlah transaksi tidak boleh kurang dari atau sama dengan nol")
	}

	// Membuat transaksi baru
	log.Println("Membuat transaksi baru...")
	transaction := &models.Transaction{
		ID:         generateTransactionID(),
		CustomerID: customerID,
		MerchantID: merchantID,
		Amount:     amount,
	}

	// Menyimpan transaksi ke repository
	log.Println("Menyimpan transaksi...")
	err = s.transactionRepository.SaveTransaction(transaction)
	if err != nil {
		return fmt.Errorf("gagal menyimpan transaksi: %w", err)
	}

	log.Println("Transaksi berhasil diproses.")

	return nil
}

func (s *TransactionService) GetMerchantNameByID(merchantID string) (string, error) {
	merchant, err := s.merchantRepository.GetByID(merchantID)
	if err != nil {
		return "", fmt.Errorf("gagal mendapatkan merchant: %w", err)
	}

	return merchant.Name, nil
}

// Fungsi bantu untuk menghasilkan ID transaksi yang unik
func generateTransactionID() string {
	transactionCounter := time.Now().UnixNano() // Menggunakan timestamp saat ini sebagai basis ID transaksi
	return strconv.FormatInt(transactionCounter, 10)
}
