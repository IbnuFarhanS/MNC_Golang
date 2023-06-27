package router

import (
	"log"
	"net/http"

	"github.com/IbnuFarhanS/Golang_MNC/internal/controller"
	"github.com/IbnuFarhanS/Golang_MNC/middleware"
	"github.com/gorilla/mux"
)

// Router represents the HTTP router
type Router struct {
	router *mux.Router
}

// NewRouter creates a new instance of Router
func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

// RegisterCustomerRoutes registers customer-related routes
func (r *Router) RegisterCustomerRoutes(customerController *controller.CustomerController) {
	log.Println("Registering customer routes...")
	r.router.HandleFunc("/login", customerController.Login).Methods(http.MethodPost)
	log.Println("Customer routes registered.")
}

// RegisterTransactionRoutes registers transaction-related routes
func (r *Router) RegisterTransactionRoutes(transactionController *controller.TransactionController) {
	log.Println("Registering transaction routes...")
	// Create a new subrouter for transaction routes
	subrouter := r.router.PathPrefix("/transaction").Subrouter()

	// Apply the AuthMiddleware to the subrouter
	subrouter.Use(middleware.AuthMiddleware)

	// Register the transaction route
	subrouter.HandleFunc("", transactionController.ProcessTransaction).Methods(http.MethodPost)
	log.Println("Transaction routes registered.")
}

// GetHandler returns the HTTP handler
func (r *Router) GetHandler() http.Handler {
	return r.router
}
