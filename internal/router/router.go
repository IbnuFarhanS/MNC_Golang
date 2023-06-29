package router

import (
	"log"
	"net/http"

	"github.com/IbnuFarhanS/Golang_MNC/internal/controller"
	"github.com/IbnuFarhanS/Golang_MNC/middleware"
	"github.com/gorilla/mux"
)

// Router mewakili router HTTP
type Router struct {
	router *mux.Router
}

// NewRouter membuat instance baru dari Router
func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

// RegisterCustomerRoutes mendaftarkan rute terkait pelanggan
func (r *Router) RegisterCustomerRoutes(customerController *controller.CustomerController) {
	log.Println("Mendaftarkan rute pelanggan...")
	r.router.HandleFunc("/register", customerController.Register).Methods(http.MethodPost)
	r.router.HandleFunc("/login", customerController.Login).Methods(http.MethodPost)

	// Membuat subrouter baru untuk rute terkait pelanggan
	customerSubrouter := r.router.PathPrefix("/customer").Subrouter()

	// Menerapkan AuthMiddleware ke subrouter pelanggan
	customerSubrouter.Use(middleware.AuthMiddleware(customerController.CustomerRepo))

	// Mendaftarkan rute logout
	customerSubrouter.HandleFunc("/logout", customerController.Logout).Methods(http.MethodPost)
	log.Println("Rute pelanggan terdaftar.")
}

// RegisterTransactionRoutes mendaftarkan rute terkait transaksi
func (r *Router) RegisterTransactionRoutes(transactionController *controller.TransactionController) {
	log.Println("Mendaftarkan rute transaksi...")
	// Membuat subrouter baru untuk rute transaksi
	subrouter := r.router.PathPrefix("/transaction").Subrouter()

	// Menerapkan AuthMiddleware ke subrouter transaksi
	subrouter.Use(middleware.AuthMiddleware(transactionController.CustomerRepo))

	// Mendaftarkan rute transaksi
	subrouter.HandleFunc("", transactionController.ProcessTransaction).Methods(http.MethodPost)
	log.Println("Rute transaksi terdaftar.")
}

// GetHandler mengembalikan handler HTTP
func (r *Router) GetHandler() http.Handler {
	return r.router
}
