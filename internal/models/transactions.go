package models

// Transaction represents a transaction
type Transaction struct {
	ID          string  `json:"id"`
	CustomerID  string  `json:"customer_id"`
	MerchantID  string  `json:"merchant_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}
