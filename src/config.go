package src

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ENCRYPTION_KEY string
	JWT_KEY        string
)

var CiperBlock cipher.Block

func ConnectDatabase() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASS"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"),
	)
	c, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError:                           true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	return c
}

func InitCipherBlock() {
	key := []byte(ENCRYPTION_KEY)
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	CiperBlock = c
}

// for testing
func GenerateKey() []byte {
	key := make([]byte, 128/8) // Convert bits to bytes
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}
