package main

import (
	"fmt"
	"net/http"

	"github.com/salamanderman234/vocagame-technical-test/src"
	"github.com/spf13/viper"
)

func Migrate() {
	models := []any{
		src.Product{},
		src.User{},
		src.Wallet{},
		src.WalletHistory{},
		src.Transaction{},
		src.TransactionDetail{},
		src.TransactionPayment{},
	}
	connection := src.ConnectDatabase()
	connection.AutoMigrate(models...)
}

func init() {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	Migrate()
	src.ENCRYPTION_KEY = viper.GetString("ENCRYPTION_KEY")
	src.JWT_KEY = viper.GetString("JWT_KEY")
}

func main() {
	mux := http.NewServeMux()
	middlewares := []src.Middleware{
		src.GetUserFromRequest,
	}

	conn := src.ConnectDatabase()
	src.InitCipherBlock()

	// product
	productRepo := src.NewProductRepo(conn)
	productService := src.NewProductService(productRepo)
	productController := src.NewProductController(productService)
	mux.HandleFunc("/api/products/", productController.Handle)

	// wallet
	walletRepo := src.NewWalletRepo(conn)
	walletService := src.NewWalletService(walletRepo)
	walletController := src.NewWalletController(walletService)
	mux.HandleFunc("/api/wallets/", walletController.Handle)
	mux.HandleFunc("/api/wallets/deposit/", walletController.Deposit)
	mux.HandleFunc("/api/wallets/withdraw/", walletController.Withdraw)

	// transaction
	transactionRepo := src.NewTransactionRepo(conn)
	transactionService := src.NewTransactionService(transactionRepo, productRepo)
	transactionController := src.NewTransactionController(transactionService)
	mux.HandleFunc("/api/transactions/", transactionController.Handle)
	mux.HandleFunc("/api/transactions/payments", transactionController.CreatePaymentTransaction)

	finalHandler := src.RegisterMiddlewares(mux, middlewares...)
	fmt.Println("Starting server at port 8000...")
	http.ListenAndServe(":8000", finalHandler)
}
