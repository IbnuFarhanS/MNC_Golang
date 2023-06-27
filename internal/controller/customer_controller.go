package controller

import (
	"encoding/json"
	"net/http"

	"github.com/IbnuFarhanS/Golang_MNC/internal/service"
	"github.com/IbnuFarhanS/Golang_MNC/utils"
)

type LoginResponse struct {
	Success  bool   `json:"success"`
	Username string `json:"username"`
	Password string `json:"password"`
	Message  string `json:"message"`
	Token    string `json:"token"`
}

type LoginFailed struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CustomerController handles customer-related HTTP requests
type CustomerController struct {
	service *service.CustomerService
}

// NewCustomerController creates a new instance of CustomerController
func NewCustomerController(service *service.CustomerService) *CustomerController {
	return &CustomerController{
		service: service,
	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login handles the login HTTP request
func (h *CustomerController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	success, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if success {
		token, err := utils.GenerateToken(req.Username, "user")
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		hashedPassword, err := utils.GenerateHash(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		resp := LoginResponse{
			Success:  true,
			Username: req.Username,
			Password: hashedPassword,
			Message:  "Login successful",
			Token:    token,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
			return
		}
	} else {
		resp := LoginFailed{
			Success: false,
			Message: "Login failed, username or password incorrect",
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
			return
		}
	}
}
