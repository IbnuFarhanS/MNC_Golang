package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
)

// TransactionRepository handles storage and retrieval of transactions
type TransactionRepository struct {
	filePath string
}

// NewTransactionRepository creates a new instance of TransactionRepository
func NewTransactionRepository(filePath string) *TransactionRepository {
	return &TransactionRepository{
		filePath: filePath,
	}
}

// SaveTransaction saves a transaction to the JSON file
func (r *TransactionRepository) SaveTransaction(transaction *models.Transaction) error {
	// Open the JSON file
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read existing transactions from the file
	transactions, err := r.getTransactionsFromFile()
	if err != nil {
		return err
	}

	// Append the new transaction
	transactions = append(transactions, *transaction)

	// Encode the updated transactions as JSON
	transactionJSON, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	// Write the updated transactions to the file
	err = ioutil.WriteFile(r.filePath, transactionJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetTransactionsByCustomerID retrieves all transactions associated with a customer ID
func (r *TransactionRepository) GetTransactionsByCustomerID(customerID string) ([]models.Transaction, error) {
	// Read the JSON file
	file, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction

	// Decode transactions from JSON
	err = json.Unmarshal(file, &transactions)
	if err != nil {
		return nil, err
	}

	filteredTransactions := make([]models.Transaction, 0)

	// Filter transactions by customer ID
	for _, transaction := range transactions {
		if transaction.CustomerID == customerID {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	return filteredTransactions, nil
}

// Helper function to get transactions from the file
func (r *TransactionRepository) getTransactionsFromFile() ([]models.Transaction, error) {
	// Read the JSON file
	file, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction

	// Decode transactions from JSON
	err = json.Unmarshal(file, &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
