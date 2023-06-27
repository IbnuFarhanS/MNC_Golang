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
}

func NewTransactionService(transactionRepository *repository.TransactionRepository, customerRepository repository.CustomerRepository) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		customerRepository:    customerRepository,
	}
}

func (s *TransactionService) ProcessTransaction(customerID string, amount float64, description string) error {
	log.Println("Processing transaction...")

	// Validate the customer ID
	log.Println("Validating customer ID...")
	_, err := s.customerRepository.GetByID(customerID)
	if err != nil {
		return errors.New("invalid customer ID")
	}

	// Validate the amount
	log.Println("Validating amount...")
	if amount <= 0 {
		return errors.New("minimum amount is 20000")
	}

	// Create a new transaction
	log.Println("Creating new transaction...")
	transaction := &models.Transaction{
		ID:          generateTransactionID(),
		CustomerID:  customerID,
		Amount:      amount,
		Description: description,
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

// Helper function to generate a unique transaction ID
func generateTransactionID() string {
	transactionCounter := time.Now().UnixNano() // Menggunakan timestamp saat ini sebagai basis ID transaksi
	return strconv.FormatInt(transactionCounter, 10)
}
