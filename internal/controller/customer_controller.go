package controller

import (
	"encoding/json"
	"net/http"

	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
	"github.com/IbnuFarhanS/Golang_MNC/internal/service"
	"github.com/IbnuFarhanS/Golang_MNC/utils"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
}

type RegisterResponse struct {
	Success  bool   `json:"success"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
	Message  string `json:"message"`
}

type RegisterFailed struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginResponse struct {
	Success  bool   `json:"success"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Token    string `json:"token"`
}

type LoginFailed struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CustomerController menangani permintaan HTTP terkait pelanggan
type CustomerController struct {
	CustomerRepo repository.CustomerRepository
	service      *service.CustomerService
}

// NewCustomerController membuat instance baru dari CustomerController
func NewCustomerController(customerRepo repository.CustomerRepository) *CustomerController {
	service := service.NewCustomerService(customerRepo)
	return &CustomerController{
		CustomerRepo: customerRepo,
		service:      service,
	}
}

// LoginRequest mewakili payload permintaan login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login menangani permintaan HTTP login
func (h *CustomerController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Payload permintaan tidak valid", http.StatusBadRequest)
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
			http.Error(w, "Gagal menghasilkan token", http.StatusInternalServerError)
			return
		}

		resp := LoginResponse{
			Success:  true,
			Username: req.Username,
			Message:  "Login berhasil",
			Token:    token,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, "Gagal mengodekan respons JSON", http.StatusInternalServerError)
			return
		}
	} else {
		resp := LoginFailed{
			Success: false,
			Message: "Login gagal, username atau password salah",
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, "Gagal mengodekan respons JSON", http.StatusInternalServerError)
			return
		}
	}
}

// Register menangani permintaan HTTP registrasi pelanggan
func (h *CustomerController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Payload permintaan tidak valid", http.StatusBadRequest)
		return
	}

	err = h.service.Register(req.Name, req.Username, req.Password, req.Phone)
	if err != nil {
		resp := RegisterFailed{
			Success: false,
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, "Gagal mengodekan respons JSON", http.StatusInternalServerError)
			return
		}
		return
	}

	hashedPassword, err := utils.GenerateHash(req.Password)
	if err != nil {
		http.Error(w, "Gagal menghash password", http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		Success:  true,
		Name:     req.Name,
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Message:  "Registrasi berhasil",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		http.Error(w, "Gagal mengodekan respons JSON", http.StatusInternalServerError)
		return
	}
}

// Logout menangani permintaan HTTP logout pelanggan
func (h *CustomerController) Logout(w http.ResponseWriter, r *http.Request) {
	// Mengambil token dari header permintaan
	token := r.Header.Get("Authorization")

	err := h.service.Logout(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := LogoutResponse{
		Success: true,
		Message: "Logout berhasil",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		http.Error(w, "Gagal mengodekan respons JSON", http.StatusInternalServerError)
		return
	}
}
