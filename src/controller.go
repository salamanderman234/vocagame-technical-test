package src

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// product
type productController struct {
	service ProductServiceInterface
}

func NewProductController(service ProductServiceInterface) productController {
	return productController{
		service: service,
	}
}
func (p productController) Handle(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodGet:
		p.read(w, r)
	case http.MethodPost:
		p.create(w, r)
	case http.MethodPatch:
		p.update(w, r)
	case http.MethodDelete:
		p.delete(w, r)
	default:
		SendJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
}

func (p productController) create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	product := Product{}
	err := decoder.Decode(&product)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	err = ValidateStruct(product)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	result, err := p.service.Create(GetRequestContext(r), product)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, "created", result)
}
func (p productController) read(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	q := r.URL.Query().Get("q")

	if id != "" {
		idInt, _ := strconv.Atoi(id)
		data, err := p.service.Find(GetRequestContext(r), idInt)
		if err != nil {
			ErrorHandler(w, err)
			return
		}
		SendJSON(w, http.StatusOK, "ok", data)
	}

	data, err := p.service.Read(GetRequestContext(r), q)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusOK, "ok", data)
}
func (p productController) update(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		SendJSON(w, http.StatusBadRequest, "missing id paramter", nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	product := Product{}
	err := decoder.Decode(&product)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	err = ValidateStruct(product)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	idInt, _ := strconv.Atoi(id)
	result, err := p.service.Update(GetRequestContext(r), idInt, product)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusOK, "ok", result)
}
func (p productController) delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		SendJSON(w, http.StatusBadRequest, "missing id paramter", nil)
		return
	}
	idInt, _ := strconv.Atoi(id)
	err := p.service.Delete(GetRequestContext(r), idInt)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusOK, "ok", nil)
}

// end of product

// wallet
type walletController struct {
	service WalletServiceInterface
}

func NewWalletController(service WalletServiceInterface) WalletControllerInterface {
	return &walletController{
		service: service,
	}
}
func (wa walletController) Handle(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodGet:
		wa.Read(w, r)
	case http.MethodPost:
		wa.Create(w, r)
	default:
		SendJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
}
func (wa walletController) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := struct {
		PaymentMethod string `json:"amount" valid:"required"`
	}{}
	err := decoder.Decode(&data)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	err = ValidateStruct(data)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	result, err := wa.service.CreateNewWallet(GetRequestContext(r), data.PaymentMethod)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, "created", result)
}
func (wa walletController) Read(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id != "" {
		idInt, _ := strconv.Atoi(id)
		data, err := wa.service.WalletDetail(GetRequestContext(r), idInt)
		if err != nil {
			ErrorHandler(w, err)
			return
		}
		SendJSON(w, http.StatusOK, "ok", data)
	}

	data, err := wa.service.GetWallet(GetRequestContext(r))
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusOK, "ok", data)
}
func (wa walletController) Withdraw(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != http.MethodPost {
		SendJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	data := struct {
		WalletID int   `json:"wallet_id" valid:"required"`
		Amount   int64 `json:"amount" valid:"required"`
	}{}
	err := decoder.Decode(&data)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	err = ValidateStruct(data)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	err = wa.service.Withdraw(GetRequestContext(r), data.WalletID, data.Amount)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusOK, "ok", nil)
}
func (wa walletController) Deposit(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != http.MethodPost {
		SendJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	data := struct {
		WalletID int   `json:"wallet_id" valid:"required"`
		Amount   int64 `json:"amount" valid:"required"`
	}{}
	err := decoder.Decode(&data)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	err = ValidateStruct(data)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	err = wa.service.Deposit(GetRequestContext(r), data.WalletID, data.Amount)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusOK, "ok", nil)
}

// end of wallet
// transaction
type transactionController struct {
	service TransactionServiceInterface
}

func NewTransactionController(service TransactionServiceInterface) TransactionControllerInterface {
	return transactionController{
		service: service,
	}
}
func (t transactionController) Handle(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodPost:
		t.CreateTransaction(w, r)
	case http.MethodPatch:
		t.UpdateTransaction(w, r)
	default:
		SendJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
}
func (t transactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := TransactionForm{}
	err := decoder.Decode(&data)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	err = ValidateStruct(data)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	var result Transaction
	if data.Type == "product" {
		result, err = t.service.CreateProductTransaction(GetRequestContext(r), data)
	} else if data.Type == "wallet" {
		result, err = t.service.CreateWalletTransaction(GetRequestContext(r), data)
	} else {
		SendJSON(w, http.StatusBadRequest, "invalid transaction type", nil)
	}

	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, "created", result)
}
func (t transactionController) CreatePaymentTransaction(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != http.MethodPost {
		SendJSON(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
	decoder := json.NewDecoder(r.Body)
	data := TransactionPaymentForm{}
	err := decoder.Decode(&data)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	err = ValidateStruct(data)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	err = t.service.CreatePaymentTransaction(GetRequestContext(r), data)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, "created", nil)
}
func (t transactionController) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		SendJSON(w, http.StatusBadRequest, "missing id paramter", nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	data := Transaction{}
	err := decoder.Decode(&data)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	err = ValidateStruct(data)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	id, _ := strconv.Atoi(idStr)
	err = t.service.UpdateTransaction(GetRequestContext(r), uint(id), data)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, "created", nil)
}

// end of transaction
