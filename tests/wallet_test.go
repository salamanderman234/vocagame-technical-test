package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/salamanderman234/vocagame-technical-test/src"
)

func init() {
	InitFunc()
}

func TestCreateNewWallet(t *testing.T) {
	client := &http.Client{}
	// test case tidak menyertakan authorization token
	data := `{
		"payment_method" : "transfer"
	}`
	request, err := http.NewRequest("POST", "http://localhost:8000/api/wallets/", bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}
	t.Logf("Tidak menyertakan authorization token : %s", data)
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusUnauthorized, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case berhasil
	data = `{
		"payment_method" : "transfer"
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	t.Logf("Menyertakan data yang lengkap : %s", data)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	respBody := &struct {
		Message string     `json:"message"`
		Details src.Wallet `json:"details"`
	}{}
	derr := json.NewDecoder(resp.Body).Decode(respBody)
	wallet = respBody.Details
	if derr != nil {
		panic(derr)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusCreated, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case berhasil
	data = `{
		"payment_method" : "transfer"
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token1)
	if err != nil {
		panic(err)
	}
	t.Logf("Menyertakan data yang lengkap : %s", data)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	respBody = &struct {
		Message string     `json:"message"`
		Details src.Wallet `json:"details"`
	}{}
	derr = json.NewDecoder(resp.Body).Decode(respBody)
	wallet1 = respBody.Details
	if derr != nil {
		panic(derr)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusCreated, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
}

func TestDepositWallet(t *testing.T) {
	client := &http.Client{}
	// test case data kurang lengkap
	walletID := strconv.Itoa(int(wallet.ID))
	data := `{
		"wallet_id" : ` + walletID + `
	}`
	request, err := http.NewRequest("POST", "http://localhost:8000/api/wallets/deposit/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	t.Logf("Menyertakan data yang kurang lengkap : %s", data)
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case tidak dengan user
	walletID = strconv.Itoa(int(wallet.ID))
	data = `{
		"wallet_id" : ` + walletID + `,
		"amount" : 10000
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/deposit/", bytes.NewBuffer([]byte(data)))
	// request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	t.Logf("Tidak menyertakan user : %s", data)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusUnauthorized, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case dengan menggunakan wallet yang bukan milik user
	walletID = strconv.Itoa(int(wallet.ID))
	data = `{
		"wallet_id" : ` + walletID + `,
		"amount" : 10000
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/deposit/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token1)
	if err != nil {
		panic(err)
	}
	t.Logf("Menggunakan wallet yang bukan milik user : %s", data)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusUnauthorized, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case dengan data yang lengkap
	data = `{
		"wallet_id" : ` + walletID + `,
		"amount" : 10000
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/deposit/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	t.Logf("Menyertakan data yang lengkap : %s", data)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusOK, resp.StatusCode)
	}
	t.Logf("History total balance harus sesuai dengan total balance terakhir")
	wall := src.Wallet{}
	conn.Where("id = ?", wallet.ID).First(&wall)
	cipherBalance := wall.Balance
	balance, err := src.DecryptData(cipherBalance)
	if err != nil {
		panic(err)
	}
	if balance != "10000" {
		t.Errorf("Test gagal, ekspetasi balance adalah 10000 tapi hasilnya %s", balance)
	}
	totalBalance := int64(0)
	histories := []src.WalletHistory{}
	conn.Where("wallet_id = ?", wallet.ID).Find(&histories)
	for _, history := range histories {
		totalBalance += history.Amount
	}
	totalBalanceStr := strconv.Itoa(int(totalBalance))
	if totalBalanceStr != balance {
		t.Errorf("Test gagal, ekspetasi total balance dari history adalah %s tapi hasilnya %s", balance, totalBalanceStr)
	}
	resp.Body.Close()
	// end of test case
}

func TestWithDraw(t *testing.T) {
	client := &http.Client{}
	// test case data kurang lengkap
	walletID := strconv.Itoa(int(wallet.ID))
	data := `{
		"wallet_id" : ` + walletID + `
	}`
	request, err := http.NewRequest("POST", "http://localhost:8000/api/wallets/withdraw/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	t.Logf("Menyertakan data yang kurang lengkap : %s", data)
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case tidak dengan user
	walletID = strconv.Itoa(int(wallet.ID))
	data = `{
		"wallet_id" : ` + walletID + `,
		"amount" : 10000
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/withdraw/", bytes.NewBuffer([]byte(data)))
	// request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	t.Logf("Tidak menyertakan user : %s", data)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusUnauthorized, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case dengan menggunakan wallet yang bukan milik user
	walletID = strconv.Itoa(int(wallet.ID))
	data = `{
		"wallet_id" : ` + walletID + `,
		"amount" : 10000
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/withdraw/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token1)
	if err != nil {
		panic(err)
	}
	t.Logf("Menggunakan wallet yang bukan milik user : %s", data)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusUnauthorized, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case dengan data yang lengkap
	data = `{
		"wallet_id" : ` + walletID + `,
		"amount" : 10000
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/withdraw/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	t.Logf("Menyertakan data yang lengkap : %s", data)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusOK, resp.StatusCode)
	}
	t.Logf("History total balance harus sesuai dengan total balance terakhir")
	wall := src.Wallet{}
	conn.Where("id = ?", wallet.ID).First(&wall)
	cipherBalance := wall.Balance
	balance, err := src.DecryptData(cipherBalance)
	if err != nil {
		panic(err)
	}
	if balance != "0" {
		t.Errorf("Test gagal, ekspetasi balance adalah 0 tapi hasilnya %s", balance)
	}
	totalBalance := int64(0)
	histories := []src.WalletHistory{}
	conn.Where("wallet_id = ?", wallet.ID).Find(&histories)
	for _, history := range histories {
		totalBalance += history.Amount
	}
	totalBalanceStr := strconv.Itoa(int(totalBalance))
	if totalBalanceStr != balance {
		t.Errorf("Test gagal, ekspetasi total balance dari history adalah %s tapi hasilnya %s", balance, totalBalanceStr)
	}
	resp.Body.Close()
	// end of test case
}
