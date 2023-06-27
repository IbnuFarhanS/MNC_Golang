package api

import (
	"log"
	"net/http"

	"github.com/IbnuFarhanS/Golang_MNC/internal/controller"
	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
	"github.com/IbnuFarhanS/Golang_MNC/internal/service"
	"github.com/IbnuFarhanS/Golang_MNC/middleware"
	"github.com/gorilla/mux"
)

// App represents the API application
type App struct {
	router *mux.Router
}

// NewApp creates a new instance of the App
func NewApp() *App {
	return &App{
		router: mux.NewRouter(),
	}
}

// Initialize initializes the application
func (a *App) Initialize() {
	log.Println("Initializing the application...")

	customerRepo, err := repository.NewInMemoryCustomerRepository("json/customers.json")
	if err != nil {
		log.Fatal(err)
	}
	customerService := service.NewCustomerService(customerRepo)
	customerController := controller.NewCustomerController(customerService)

	transactionRepo := repository.NewTransactionRepository("json/transactions.json")
	transactionService := service.NewTransactionService(transactionRepo, customerRepo)
	transactionController := controller.NewTransactionController(transactionService)

	// Register customer routes
	log.Println("Registering customer routes...")
	a.router.HandleFunc("/login", customerController.Login).Methods(http.MethodPost)
	log.Println("Customer routes registered.")

	// Register transaction route with middleware
	log.Println("Registering transaction routes...")
	transactionHandler := middleware.AuthMiddleware(http.HandlerFunc(transactionController.ProcessTransaction))
	a.router.Handle("/transaction", transactionHandler).Methods(http.MethodPost)
	log.Println("Transaction routes registered.")

	log.Println("Application initialized.")
}

// Run starts the application
func (a *App) Run(port string) {
	log.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, a.router))
}
