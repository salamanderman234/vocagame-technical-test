package src

import (
	"context"
	"net/http"
)

// product domain

type ProductRepositoryInterface interface {
	Create(data Product) (Product, error)
	Read(q string) ([]Product, error)
	Find(id int) (Product, error)
	BatchSearch(ids []uint) ([]Product, error)
	Update(id int, data Product) (Product, error)
	Delete(id int) error
}

type ProductServiceInterface interface {
	Create(ctx context.Context, data Product) (Product, error)
	Read(ctx context.Context, q string) ([]Product, error)
	Find(ctx context.Context, id int) (Product, error)
	Update(ctx context.Context, id int, data Product) (Product, error)
	Delete(ctx context.Context, id int) error
}

type ProductControllerInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// end of product domain

// wallet domain
type WalletRepositoryInterface interface {
	Create(data Wallet) (Wallet, error)
	Read(userID int) ([]Wallet, error)
	Find(id int) (Wallet, error)
	Deposit(wallet Wallet, amount int64) error
	Withdraw(wallet Wallet, amount int64) error
}
type WalletServiceInterface interface {
	CreateNewWallet(ctx context.Context, paymentMethod string) (Wallet, error)
	GetWallet(ctx context.Context) ([]Wallet, error)
	WalletDetail(ctx context.Context, id int) (Wallet, error)
	Deposit(ctx context.Context, walletID int, amount int64) error
	Withdraw(ctx context.Context, walletID int, amount int64) error
}
type WalletControllerInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Read(w http.ResponseWriter, r *http.Request)
	Withdraw(w http.ResponseWriter, r *http.Request)
	Deposit(w http.ResponseWriter, r *http.Request)
}

// end of wallet domain
// transaction
type TransactionRepositoryInterface interface {
	CreateTransaction(data Transaction) (Transaction, error)
	CreateTransactionPayment(data TransactionPayment) error
	UpdateTransaction(id uint, data Transaction) error
}
type TransactionServiceInterface interface {
	CreateProductTransaction(ctx context.Context, data TransactionForm) (Transaction, error)
	CreateWalletTransaction(ctx context.Context, data TransactionForm) (Transaction, error)
	CreatePaymentTransaction(ctx context.Context, data TransactionPaymentForm) error
	UpdateTransaction(ctx context.Context, id uint, data Transaction) error
}
type TransactionControllerInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
	CreateTransaction(w http.ResponseWriter, r *http.Request)
	CreatePaymentTransaction(w http.ResponseWriter, r *http.Request)
	UpdateTransaction(w http.ResponseWriter, r *http.Request)
}

// end of transaction
