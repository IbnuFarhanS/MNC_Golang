package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
	"github.com/IbnuFarhanS/Golang_MNC/internal/service"
	"github.com/IbnuFarhanS/Golang_MNC/middleware"
)

type TransactionController struct {
	CustomerRepo       repository.CustomerRepository
	transactionService *service.TransactionService
}

func NewTransactionController(customerRepo repository.CustomerRepository, transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		CustomerRepo:       customerRepo,
		transactionService: transactionService,
	}
}

type TransactionRequest struct {
	CustomerID string  `json:"customer_id"`
	MerchantID string  `json:"merchant_id"`
	Amount     float64 `json:"amount"`
}

type TransactionResponse struct {
	Success      bool    `json:"success"`
	CustomerID   string  `json:"customer_id"`
	MerchantName string  `json:"merchant_name"`
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
	Message      string  `json:"message"`
}

func (h *TransactionController) ProcessTransaction(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing transaction...")

	// Extract user ID from request context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		log.Println("Failed to extract user ID from context")
		http.Error(w, "Failed to extract user ID from context", http.StatusInternalServerError)
		return
	}

	log.Println("User ID:", userID)

	var req TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Invalid request payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println("Received transaction request:", req)

	err = h.transactionService.ProcessTransaction(req.CustomerID, req.MerchantID, req.Amount)
	if err != nil {
		log.Println("Failed to process transaction:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the merchant name
	merchantName, err := h.transactionService.GetMerchantNameByID(req.MerchantID)
	if err != nil {
		log.Println("Failed to get merchant name:", err)
		http.Error(w, "Failed to get merchant name", http.StatusInternalServerError)
		return
	}

	resp := TransactionResponse{
		Success:      true,
		CustomerID:   req.CustomerID,
		MerchantName: merchantName,
		Amount:       req.Amount,
		Description:  fmt.Sprintf("payment for %s with amount %.2f success", merchantName, req.Amount),
		Message:      "Transaction success",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Println("Failed to encode JSON response:", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	log.Println("Transaction processed successfully.")
}
