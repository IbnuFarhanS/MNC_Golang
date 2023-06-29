package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
)

// Mendefinisikan interface MerchantRepository yang menyediakan method-method
type MerchantRepository interface {
	GetByID(merchantID string) (*models.Merchant, error)
	GetMerchantNameByID(merchantID string) (string, error)
}

// Data merchant disimpan dalam slice of Merchant
type InMemoryMerchantRepository struct {
	merchants []*models.Merchant
}

func GetMerchants(filePath string) ([]*models.Merchant, error) {
	// Membaca file yang berisi data merchant
	data, err := ioutil.ReadFile(filePath)

	// Mengembalikan error jika gagal membaca file
	if err != nil {
		return nil, fmt.Errorf("failed to read merchant data: %v", err)
	}

	// Mendekode data JSON menjadi slice of Merchant
	var merchants []*models.Merchant
	err = json.Unmarshal(data, &merchants)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal merchant data: %v", err)
	}

	// Mengembalikan slice of Merchant yang berhasil didapatkan
	return merchants, nil
}

func NewInMemoryMerchantRepository(filePath string) (*InMemoryMerchantRepository, error) {
	// Mendapatkan data merchant menggunakan fungsi GetMerchants
	merchants, err := GetMerchants(filePath)

	// Mengembalikan error jika gagal mendapatkan data merchant
	if err != nil {
		return nil, err
	}

	// Mengembalikan instance InMemoryMerchantRepository yang berisi data merchant
	return &InMemoryMerchantRepository{
		merchants: merchants,
	}, nil
}

func (r *InMemoryMerchantRepository) GetByID(merchantID string) (*models.Merchant, error) {
	for _, merchant := range r.merchants {
		// Mencari merchant berdasarkan ID
		if merchant.ID == merchantID {
			// Mengembalikan merchant yang ditemukan
			return merchant, nil
		}
	}

	// Mengembalikan error jika merchant tidak ditemukan
	return nil, fmt.Errorf("merchant not found")
}

func (r *InMemoryMerchantRepository) GetMerchantNameByID(merchantID string) (string, error) {
	for _, merchant := range r.merchants {
		// Mencari merchant berdasarkan ID
		if merchant.ID == merchantID {
			// Mengembalikan nama merchant yang ditemukan
			return merchant.Name, nil
		}
	}

	// Mengembalikan error jika merchant tidak ditemukan
	return "", fmt.Errorf("merchant not found")
}
