package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// product repo
type productRepo struct {
	conn *gorm.DB
}

func NewProductRepo(conn *gorm.DB) ProductRepositoryInterface {
	return productRepo{
		conn: conn,
	}
}
func (p productRepo) Create(data Product) (Product, error) {
	result := p.conn.Create(&data)
	return data, result.Error
}
func (p productRepo) Read(q string) ([]Product, error) {
	products := []Product{}
	result := p.conn.Where("product_name LIKE ?", "%"+q+"%").Find(&products)
	if len(products) == 0 && result.Error == nil {
		return products, ErrNotFound
	}
	return products, result.Error
}
func (p productRepo) Find(id int) (Product, error) {
	product := Product{}
	result := p.conn.Where("id = ?", id).First(&product)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return product, ErrNotFound
	}
	return product, result.Error
}
func (p productRepo) Update(id int, data Product) (Product, error) {
	result := p.conn.Where("id = ?", id).Updates(data)
	if result.RowsAffected == 0 {
		return data, ErrNotFound
	}
	return data, result.Error
}
func (p productRepo) Delete(id int) error {
	result := p.conn.Where("id = ?", id).Delete(&Product{})
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}
func (p productRepo) BatchSearch(ids []uint) ([]Product, error) {
	products := []Product{}
	result := p.conn.Where("(id) IN ?", ids).Find(&products)
	if len(products) == 0 && result.Error == nil {
		return products, ErrNotFound
	}
	return products, result.Error
}

// end of product repo

// wallet repo
type walletRepo struct {
	conn *gorm.DB
}

func NewWalletRepo(conn *gorm.DB) WalletRepositoryInterface {
	return walletRepo{
		conn: conn,
	}
}

func (wa walletRepo) getWalletHistory(wallet Wallet) (string, error) {
	id := wallet.ID
	histories := []WalletHistory{}
	wa.conn.Where("wallet_id = ?", id).Find(&histories)
	data, err := json.Marshal(histories)
	if err != nil {
		return "", err
	}
	cipherText, err := EncryptData(string(data))
	if err != nil {
		return "", err
	}
	return cipherText, nil
}

func (wa walletRepo) Create(data Wallet) (Wallet, error) {
	result := wa.conn.Create(&data)
	return data, result.Error
}
func (wa walletRepo) Read(userID int) ([]Wallet, error) {
	wallets := []Wallet{}
	result := wa.conn.Find(&wallets, Wallet{
		UserID: uint(userID),
	})
	if len(wallets) == 0 && result.Error == nil {
		return wallets, ErrNotFound
	}
	return wallets, result.Error
}
func (wa walletRepo) Find(id int) (Wallet, error) {
	wallet := Wallet{}
	result := wa.conn.Where("id = ?", id).First(&wallet)
	if wallet.ID == 0 && result.Error == nil {
		return wallet, ErrNotFound
	}
	history, err := wa.getWalletHistory(wallet)
	if err != nil {
		return Wallet{}, err
	}
	wallet.Histories = history
	return wallet, result.Error
}
func (wa walletRepo) Deposit(wallet Wallet, amount int64) error {
	balanceCipher := wallet.Balance
	decryptBalance, err := DecryptData(balanceCipher)
	if err != nil {
		return err
	}
	balance, _ := strconv.ParseInt(decryptBalance, 10, 64)
	finalBalance := balance + amount
	finalBalanceStr := strconv.FormatInt(finalBalance, 10)
	balanceCipher, err = EncryptData(finalBalanceStr)
	if err != nil {
		return err
	}

	return wa.conn.Transaction(func(tx *gorm.DB) error {
		wallet.Balance = balanceCipher
		result := tx.Save(wallet)
		if result.Error != nil {
			return result.Error
		}
		history := WalletHistory{
			WalletID:    wallet.ID,
			Date:        time.Now(),
			Amount:      amount,
			Description: fmt.Sprintf("deposit %d", amount),
		}
		result = tx.Create(&history)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}
func (wa walletRepo) Withdraw(wallet Wallet, amount int64) error {
	balanceCipher := wallet.Balance
	decryptBalance, err := DecryptData(balanceCipher)
	if err != nil {
		return err
	}
	balance, _ := strconv.ParseInt(decryptBalance, 10, 64)
	finalBalance := balance - amount
	if finalBalance < 0 {
		return ErrInvalidWithdrawAmount
	}
	finalBalanceStr := strconv.FormatInt(finalBalance, 10)
	balanceCipher, err = EncryptData(finalBalanceStr)
	if err != nil {
		return err
	}

	return wa.conn.Transaction(func(tx *gorm.DB) error {
		wallet.Balance = balanceCipher
		result := tx.Save(wallet)
		if result.Error != nil {
			return result.Error
		}
		history := WalletHistory{
			WalletID:    wallet.ID,
			Date:        time.Now(),
			Amount:      -amount,
			Description: fmt.Sprintf("withdraw %d", amount),
		}
		result = tx.Create(&history)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
}

// end of wallet repo

// transaction repo
type transactionRepo struct {
	conn *gorm.DB
}

func NewTransactionRepo(conn *gorm.DB) TransactionRepositoryInterface {
	return transactionRepo{
		conn: conn,
	}
}
func (t transactionRepo) CreateTransaction(data Transaction) (Transaction, error) {
	err := t.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			return err
		}
		details := data.Details
		prods := []Product{}
		for _, detail := range details {
			prod := Product{}
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&prod, detail.ProductID).Error; err != nil {
				return err
			}
			currentStock := prod.Stock
			finalStock := int64(currentStock) - int64(detail.Qty)
			if finalStock < 0 {
				return ErrInvalidQty
			}
			prod.Stock = uint64(finalStock)
			prods = append(prods, prod)
		}
		tx.Save(prods)
		return nil
	})
	return data, err
}
func (t transactionRepo) CreateTransactionPayment(data TransactionPayment) error {
	err := t.conn.Transaction(func(tx *gorm.DB) error {
		tr := Transaction{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&tr, data.TransactionID).Error; err != nil {
			return err
		}
		if tr.WalletDestinationID == 0 {
			data.WalletDestinationID = 0
		}
		amount := tr.FinalAmount
		data.Amount = amount
		if err := tx.Create(&data).Error; err != nil {
			return err
		}
		tr.Status = "ok"

		wallet1 := Wallet{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&tr, data.WalletID).Error; err != nil {
			return err
		}
		cipherbalance := wallet1.Balance
		balanceStr, err := DecryptData(cipherbalance)
		if err != nil {
			return err
		}
		balance, err := strconv.ParseInt(balanceStr, 10, 64)
		if err != nil {
			return err
		}

		finalbalance := balance - amount

		if finalbalance < 0 {
			return ErrInvalidWithdrawAmount
		}
		finalBalanceStr := strconv.FormatInt(finalbalance, 10)
		balanceCipher, err := EncryptData(finalBalanceStr)
		if err != nil {
			return err
		}
		wallet1.Balance = balanceCipher
		tx.Save(wallet1)
		h := WalletHistory{
			WalletID:    wallet1.ID,
			Date:        time.Now(),
			Amount:      -amount,
			Description: "sending",
		}
		if err := tx.Create(&h).Error; err != nil {
			return err
		}
		if data.WalletDestinationID != 0 {
			wallet1 = Wallet{}
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&tr, data.WalletDestinationID).Error; err != nil {
				return err
			}
			cipherbalance = wallet1.Balance
			balanceStr, err = DecryptData(cipherbalance)
			if err != nil {
				return err
			}
			balance, err := strconv.ParseInt(balanceStr, 10, 64)
			if err != nil {
				return err
			}

			finalbalance := balance + amount
			finalBalanceStr := strconv.FormatInt(finalbalance, 10)
			balanceCipher, err := EncryptData(finalBalanceStr)
			if err != nil {
				return err
			}
			wallet1.Balance = balanceCipher
			tx.Save(wallet1)
			h := WalletHistory{
				WalletID:    wallet1.ID,
				Date:        time.Now(),
				Amount:      amount,
				Description: "receiving",
			}
			if err := tx.Create(&h).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
func (t transactionRepo) UpdateTransaction(id uint, data Transaction) error {
	result := t.conn.Where("id = ?", id).Updates(data)
	return result.Error
}

// end of transaction repo
