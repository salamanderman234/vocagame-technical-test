package tests

import (
	"github.com/salamanderman234/vocagame-technical-test/src"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var conn *gorm.DB

var user = src.User{
	Role:  "user",
	Email: "test@gmail.com",
}
var user1 = src.User{
	Role:  "user",
	Email: "test2@gmail.com",
}

var prod = src.Product{
	ProductName: "TEST",
	Price:       1000,
	Stock:       10,
	Visible:     true,
	Description: "TEST",
}

var token = ""
var token1 = ""
var wallet = src.Wallet{}
var wallet1 = src.Wallet{}

func InitFunc() {
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()
	src.ENCRYPTION_KEY = viper.GetString("ENCRYPTION_KEY")
	src.JWT_KEY = viper.GetString("JWT_KEY")
	src.InitCipherBlock()
	conn = src.ConnectDatabase()
	conn.Create(&user)
	conn.Create(&user1)
	conn.Create(&prod)

	tkn, err := src.CreateToken(user)
	if err != nil {
		panic(err)
	}
	token = tkn
	tkn, err = src.CreateToken(user1)
	if err != nil {
		panic(err)
	}
	token1 = tkn
}

func ResetDatabase() {
	tables := []any{
		src.Product{},
		src.User{},
		src.Wallet{},
		src.WalletHistory{},
		src.Transaction{},
		src.TransactionDetail{},
		src.TransactionPayment{},
	}
	for _, table := range tables {
		conn.Unscoped().Where("id != ?", 0).Delete(&table)
	}
}
