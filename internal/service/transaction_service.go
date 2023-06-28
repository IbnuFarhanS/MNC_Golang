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

// TransactionService handles transaction-related operations
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
	log.Println("Processing transaction...")

	// Validate the customer ID
	log.Println("Validating customer ID...")
	_, err := s.customerRepository.GetByID(customerID)
	if err != nil {
		return errors.New("invalid customer ID")
	}

	// Validate the merchant ID
	_, err = s.merchantRepository.GetByID(merchantID)
	if err != nil {
		return errors.New("invalid merchant ID")
	}

	// Validate the amount
	log.Println("Validating amount...")
	if amount <= 0 {
		return errors.New("amount cannot be less than or equal to zero")
	}

	// Create a new transaction
	log.Println("Creating new transaction...")
	transaction := &models.Transaction{
		ID:         generateTransactionID(),
		CustomerID: customerID,
		MerchantID: merchantID,
		Amount:     amount,
	}

	// Save the transaction to the repository
	log.Println("Saving transaction...")
	err = s.transactionRepository.SaveTransaction(transaction)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	log.Println("Transaction processed successfully.")

	return nil
}

func (s *TransactionService) GetMerchantNameByID(merchantID string) (string, error) {
	merchant, err := s.merchantRepository.GetByID(merchantID)
	if err != nil {
		return "", fmt.Errorf("failed to get merchant: %w", err)
	}

	return merchant.Name, nil
}

// Helper function to generate a unique transaction ID
func generateTransactionID() string {
	transactionCounter := time.Now().UnixNano() // Menggunakan timestamp saat ini sebagai basis ID transaksi
	return strconv.FormatInt(transactionCounter, 10)
}
