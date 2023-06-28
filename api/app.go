package api

import (
	"log"
	"net/http"

	"github.com/IbnuFarhanS/Golang_MNC/internal/controller"
	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
	"github.com/IbnuFarhanS/Golang_MNC/internal/router"
	"github.com/IbnuFarhanS/Golang_MNC/internal/service"
)

// App represents the API application
type App struct {
	router *router.Router
}

// NewApp creates a new instance of the App
func NewApp() *App {
	return &App{
		router: router.NewRouter(),
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
	merchantRepo, err := repository.NewInMemoryMerchantRepository("json/merchants.json")
	if err != nil {
		log.Fatal(err)
	}
	transactionService := service.NewTransactionService(transactionRepo, customerRepo, merchantRepo)
	transactionController := controller.NewTransactionController(customerRepo, transactionService)

	// Register customer routes
	log.Println("Registering customer routes...")
	a.router.RegisterCustomerRoutes(customerController)
	log.Println("Customer routes registered.")

	// Register transaction routes
	log.Println("Registering transaction routes...")
	a.router.RegisterTransactionRoutes(transactionController)
	log.Println("Transaction routes registered.")

	log.Println("Application initialized.")
}

// Run starts the application
func (a *App) Run(port string) {
	log.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, a.router.GetHandler()))
}
