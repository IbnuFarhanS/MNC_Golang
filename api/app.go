package api

import (
	"log"
	"net/http"

	"github.com/IbnuFarhanS/Golang_MNC/internal/controller"
	"github.com/IbnuFarhanS/Golang_MNC/internal/repository"
	"github.com/IbnuFarhanS/Golang_MNC/internal/router"
	"github.com/IbnuFarhanS/Golang_MNC/internal/service"
)

// App mewakili aplikasi API
type App struct {
	router *router.Router
}

// NewApp membuat instance baru dari App
func NewApp() *App {
	return &App{
		router: router.NewRouter(),
	}
}

// Initialize menginisialisasi aplikasi
func (a *App) Initialize() {
	log.Println("Menginisialisasi aplikasi...")

	customerRepo, err := repository.NewInMemoryCustomerRepository("json/customers.json")
	if err != nil {
		// Log fatal jika gagal membuat repository pelanggan dalam memori
		log.Fatal(err)
	}
	// Membuat layanan pelanggan baru dengan repository yang sudah dibuat
	customerService := service.NewCustomerService(customerRepo)
	// Membuat kontroler pelanggan baru dengan layanan pelanggan
	customerController := controller.NewCustomerController(customerService)

	// Membuat repository transaksi baru
	transactionRepo := repository.NewTransactionRepository("json/transactions.json")
	merchantRepo, err := repository.NewInMemoryMerchantRepository("json/merchants.json")
	if err != nil {
		// Log fatal jika gagal membuat repository merchant dalam memori
		log.Fatal(err)
	}
	// Membuat layanan transaksi baru dengan repository yang sudah dibuat
	transactionService := service.NewTransactionService(transactionRepo, customerRepo, merchantRepo)
	// Membuat kontroler transaksi baru dengan layanan transaksi
	transactionController := controller.NewTransactionController(customerRepo, transactionService)

	// Mendaftarkan rute pelanggan
	log.Println("Mendaftarkan rute pelanggan...")
	a.router.RegisterCustomerRoutes(customerController)
	log.Println("Rute pelanggan terdaftar.")

	// Mendaftarkan rute transaksi
	log.Println("Mendaftarkan rute transaksi...")
	a.router.RegisterTransactionRoutes(transactionController)
	log.Println("Rute transaksi terdaftar.")

	log.Println("Aplikasi diinisialisasi.")
}

// Run menjalankan aplikasi
func (a *App) Run(port string) {
	log.Printf("Server berjalan pada port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, a.router.GetHandler()))
}
