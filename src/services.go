package src

import (
	"context"
	"slices"
	"time"
)

// product service
type productService struct {
	repo ProductRepositoryInterface
}

func NewProductService(repo ProductRepositoryInterface) ProductServiceInterface {
	return productService{
		repo: repo,
	}
}
func (p productService) Create(ctx context.Context, data Product) (Product, error) {
	err := Authorize(ctx, data, ProductPolicy)
	if err != nil {
		return Product{}, err
	}
	data.ProductName = Sanitize(data.ProductName)
	data.Description = Sanitize(data.Description)
	return p.repo.Create(data)
}
func (p productService) Read(ctx context.Context, q string) ([]Product, error) {
	return p.repo.Read(q)
}
func (p productService) Find(ctx context.Context, id int) (Product, error) {
	return p.repo.Find(id)
}
func (p productService) Update(ctx context.Context, id int, data Product) (Product, error) {
	err := Authorize(ctx, data, ProductPolicy)
	if err != nil {
		return Product{}, err
	}
	data.ProductName = Sanitize(data.ProductName)
	data.Description = Sanitize(data.Description)
	return p.repo.Update(id, data)
}
func (p productService) Delete(ctx context.Context, id int) error {
	err := Authorize(ctx, nil, ProductPolicy)
	if err != nil {
		return err
	}
	return p.repo.Delete(id)
}

// end of product service

// wallet service
type walletService struct {
	repo WalletRepositoryInterface
}

func NewWalletService(repo WalletRepositoryInterface) WalletServiceInterface {
	return walletService{
		repo: repo,
	}
}
func (wa walletService) CreateNewWallet(ctx context.Context, paymentMethod string) (Wallet, error) {
	err := Authorize(ctx, nil, WalletCreatePolicy)
	if err != nil {
		return Wallet{}, err
	}
	user, err := GetUserContext(ctx)
	if err != nil {
		return Wallet{}, err
	}

	// encrypt payment method
	cipherPaymentMethod, err := EncryptData(paymentMethod)
	if err != nil {
		return Wallet{}, err
	}
	// encrypt balance
	cipherBalance, err := EncryptData("0")
	if err != nil {
		return Wallet{}, err
	}

	data := Wallet{
		UserID:        user.ID,
		Balance:       cipherBalance,
		PaymentMethod: cipherPaymentMethod,
	}
	return wa.repo.Create(data)
}
func (wa walletService) GetWallet(ctx context.Context) ([]Wallet, error) {
	err := Authorize(ctx, nil, WalletReadPolicy)
	if err != nil {
		return []Wallet{}, err
	}
	user, err := GetUserContext(ctx)
	if err != nil {
		return []Wallet{}, err
	}
	return wa.repo.Read(int(user.ID))
}
func (wa walletService) WalletDetail(ctx context.Context, id int) (Wallet, error) {
	wallet, errRepo := wa.repo.Find(id)
	err := Authorize(ctx, wallet, WalletDetailPolicy)
	if err != nil {
		return Wallet{}, err
	}
	if errRepo != nil {
		return Wallet{}, errRepo
	}
	return wallet, nil
}
func (wa walletService) Deposit(ctx context.Context, walletID int, amount int64) error {
	wallet, errRepo := wa.repo.Find(walletID)
	err := Authorize(ctx, wallet, WalletDetailPolicy)
	if err != nil {
		return err
	}
	if errRepo != nil {
		return errRepo
	}
	return wa.repo.Deposit(wallet, amount)
}
func (wa walletService) Withdraw(ctx context.Context, walletID int, amount int64) error {
	wallet, errRepo := wa.repo.Find(walletID)
	err := Authorize(ctx, wallet, WalletDetailPolicy)
	if err != nil {
		return err
	}
	if errRepo != nil {
		return errRepo
	}
	return wa.repo.Withdraw(wallet, amount)
}

// end of wallet service

// transaction service
type transactionService struct {
	repo        TransactionRepositoryInterface
	productRepo ProductRepositoryInterface
}

func NewTransactionService(repo TransactionRepositoryInterface, pr ProductRepositoryInterface) TransactionServiceInterface {
	return transactionService{
		repo:        repo,
		productRepo: pr,
	}
}
func (t transactionService) CreateProductTransaction(ctx context.Context, data TransactionForm) (Transaction, error) {
	err := Authorize(ctx, data, CreateTransactionPolicy)
	if err != nil {
		return Transaction{}, err
	}
	user, _ := GetUserContext(ctx)
	transaction := Transaction{
		UserID:              user.ID,
		FinalAmount:         0,
		WalletDestinationID: 0,
		Status:              "wait",
		Date:                time.Now(),
		Type:                "product",
		Note:                Sanitize(data.Note),
		Details:             []TransactionDetail{},
	}
	ids := []uint{}
	for _, product := range data.Products {
		ids = append(ids, product.ID)
	}
	products, err := t.productRepo.BatchSearch(ids)
	if err != nil {
		return Transaction{}, err
	}
	for _, product := range products {
		qty := data.Products[slices.IndexFunc(data.Products, func(c TransactionProductForm) bool { return c.ID == product.ID })].Qty
		price := product.Price
		totalPrice := int64(qty) * price
		detail := TransactionDetail{
			ProductID:  product.ID,
			Product:    product,
			Qty:        qty,
			Price:      price,
			TotalPrice: totalPrice,
		}
		transaction.Details = append(transaction.Details, detail)
	}
	return t.repo.CreateTransaction(transaction)
}
func (t transactionService) CreateWalletTransaction(ctx context.Context, data TransactionForm) (Transaction, error) {
	err := Authorize(ctx, data, CreateTransactionPolicy)
	if err != nil {
		return Transaction{}, err
	}
	user, _ := GetUserContext(ctx)
	if data.WalletDestinationID == 0 {
		return Transaction{}, ErrBadRequest
	}
	transaction := Transaction{
		UserID:              user.ID,
		FinalAmount:         data.FinalAmount,
		WalletDestinationID: data.WalletDestinationID,
		Status:              "wait",
		Date:                time.Now(),
		Type:                "wallet",
		Note:                Sanitize(data.Note),
	}
	return t.repo.CreateTransaction(transaction)
}
func (t transactionService) CreatePaymentTransaction(ctx context.Context, data TransactionPaymentForm) error {
	err := Authorize(ctx, data, UpdateTransactionPolicy)
	if err != nil {
		return err
	}
	transactionPayment := TransactionPayment{
		TransactionID:       data.TransactionID,
		WalletID:            data.WalletID,
		WalletDestinationID: data.WalletDestinationID,
		Status:              "ok",
		Note:                Sanitize(data.Note),
	}
	return t.repo.CreateTransactionPayment(transactionPayment)
}
func (t transactionService) UpdateTransaction(ctx context.Context, id uint, data Transaction) error {
	err := Authorize(ctx, data, UpdateTransactionPolicy)
	if err != nil {
		return err
	}
	transaction := Transaction{
		WalletDestinationID: data.WalletDestinationID,
	}
	return t.repo.UpdateTransaction(id, transaction)
}

// end of transaction service
