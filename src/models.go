package src

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Email string
	Role  string
}

type User struct {
	gorm.Model
	Role  string `json:"role"`
	Email string `json:"email"`
}

type Product struct {
	gorm.Model
	ProductName string `json:"product_name" valid:"required"`
	Price       int64  `json:"price" valid:"required"`
	Stock       uint64 `json:"stock" valid:"required"`
	Visible     bool   `json:"visible" valid:"required"`
	Description string `json:"description" valid:"required"`
}

type Wallet struct {
	gorm.Model
	UserID        uint   `json:"user_id"`
	Balance       string `json:"balance"`
	PaymentMethod string `json:"payment_method"`
	Histories     string `json:"histories"`
}

type WalletHistory struct {
	gorm.Model
	WalletID    uint      `json:"wallet_id"`
	Date        time.Time `json:"date"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
}

type TransactionProductForm struct {
	ID  uint `json:"id"`
	Qty uint `json:"qty"`
}

type TransactionForm struct {
	FinalAmount         int64                    `json:"final_amount"`
	WalletDestinationID uint                     `json:"wallet_destination_id"`
	Date                time.Time                `json:"date"`
	Note                string                   `json:"note"`
	Type                string                   `json:"type"`
	Products            []TransactionProductForm `json:"products"`
}

type Transaction struct {
	gorm.Model
	UserID              uint                 `json:"user_id"`
	User                User                 `json:"user"`
	FinalAmount         int64                `json:"final_amount"`
	WalletDestinationID uint                 `json:"wallet_destination_id"`
	Status              string               `json:"status"`
	Date                time.Time            `json:"date"`
	Note                string               `json:"note"`
	Type                string               `json:"type"`
	Details             []TransactionDetail  `json:"details"`
	Payments            []TransactionPayment `json:"payments"`
}

type TransactionDetail struct {
	gorm.Model
	TransactionID uint        `json:"transaction_id"`
	Transaction   Transaction `json:"transaction"`
	ProductID     uint        `json:"product_id"`
	Product       Product     `json:"product"`
	Qty           uint        `json:"qty"`
	Price         int64       `json:"price"`
	TotalPrice    int64       `json:"total_price"`
}
type TransactionPaymentForm struct {
	TransactionID       uint   `json:"transaction_id"`
	WalletID            uint   `json:"wallet_id"`
	WalletDestinationID uint   `json:"wallet_destination_id"`
	Note                string `json:"note"`
}
type TransactionPayment struct {
	gorm.Model
	TransactionID       uint   `json:"transaction_id"`
	WalletID            uint   `json:"wallet_id"`
	WalletDestinationID uint   `json:"wallet_destination_id"`
	Amount              int64  `json:"amount"`
	Note                string `json:"note"`
	Status              string `json:"status"`
}
