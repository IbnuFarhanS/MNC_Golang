package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
)

// TransactionRepository menangani penyimpanan dan pengambilan transaksi
type TransactionRepository struct {
	filePath string
}

// NewTransactionRepository membuat instance baru dari TransactionRepository
func NewTransactionRepository(filePath string) *TransactionRepository {
	return &TransactionRepository{
		filePath: filePath,
	}
}

// SaveTransaction menyimpan transaksi ke file JSON
func (r *TransactionRepository) SaveTransaction(transaction *models.Transaction) error {
	// Buka file JSON
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Baca transaksi yang sudah ada dari file
	transactions, err := r.getTransactionsFromFile()
	if err != nil {
		return err
	}

	// Tambahkan transaksi baru
	transactions = append(transactions, *transaction)

	// Encode transaksi yang diperbarui menjadi JSON
	transactionJSON, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	// Tulis transaksi yang diperbarui ke file
	err = ioutil.WriteFile(r.filePath, transactionJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetTransactionsByCustomerID mengambil semua transaksi yang terkait dengan ID pelanggan
func (r *TransactionRepository) GetTransactionsByCustomerID(customerID string) ([]models.Transaction, error) {
	// Baca file JSON
	file, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction

	// Decode transaksi dari JSON
	err = json.Unmarshal(file, &transactions)
	if err != nil {
		return nil, err
	}

	filteredTransactions := make([]models.Transaction, 0)

	// Filter transaksi berdasarkan ID pelanggan
	for _, transaction := range transactions {
		if transaction.CustomerID == customerID {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	return filteredTransactions, nil
}

// Fungsi bantu untuk mendapatkan transaksi dari file
func (r *TransactionRepository) getTransactionsFromFile() ([]models.Transaction, error) {
	// Baca file JSON
	file, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction

	// Decode transaksi dari JSON
	err = json.Unmarshal(file, &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
