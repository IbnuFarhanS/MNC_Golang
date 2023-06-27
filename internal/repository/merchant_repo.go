package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/IbnuFarhanS/Golang_MNC/internal/models"
)

type MerchantRepository interface {
	GetByID(merchantID string) (*models.Merchant, error)
	GetMerchantNameByID(merchantID string) (string, error)
}

func GetMerchants(filePath string) ([]*models.Merchant, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read merchant data: %v", err)
	}

	var merchants []*models.Merchant
	err = json.Unmarshal(data, &merchants)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal merchant data: %v", err)
	}

	return merchants, nil
}

type InMemoryMerchantRepository struct {
	merchants []*models.Merchant
}

func NewInMemoryMerchantRepository(filePath string) (*InMemoryMerchantRepository, error) {
	merchants, err := GetMerchants(filePath)
	if err != nil {
		return nil, err
	}

	return &InMemoryMerchantRepository{
		merchants: merchants,
	}, nil
}

func (r *InMemoryMerchantRepository) GetByID(merchantID string) (*models.Merchant, error) {
	for _, merchant := range r.merchants {
		if merchant.ID == merchantID {
			return merchant, nil
		}
	}

	return nil, fmt.Errorf("merchant not found")
}

func (r *InMemoryMerchantRepository) GetMerchantNameByID(merchantID string) (string, error) {
	for _, merchant := range r.merchants {
		if merchant.ID == merchantID {
			return merchant.Name, nil
		}
	}

	return "", fmt.Errorf("merchant not found")
}
