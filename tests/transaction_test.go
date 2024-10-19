package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/salamanderman234/vocagame-technical-test/src"
)

func TestInitiateTransaction(t *testing.T) {
	client := &http.Client{}
	// test case membuat transaksi dengan data yang tidak lengkap
	data := `{}`
	t.Logf("Menggunakan data yang tidak lengkap : %s", data)
	request, err := http.NewRequest("POST", "http://localhost:8000/api/transactions/", bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	// end test case
	// menggunakan tipe transaksi yang salah
	data = `{
		"type" : "test"
	}`
	t.Logf("Menggunakan tipe transaksi yang salah : %s", data)
	request, err = http.NewRequest("POST", "http://localhost:8000/api/transactions/", bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	// end test case
	// tidak menyertakan user
	data = `{
		"payment_method" : "transfer"
	}`
	request, err = http.NewRequest("POST", "http://localhost:8000/api/wallets/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
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
	resp.Body.Close()
	walletIDStr := strconv.Itoa(int(wallet.ID))
	data = `{
		"type" : "wallet",
		"final_amount" : 1000,
		"wallet_destination_id" : ` + walletIDStr + `
	}`
	t.Logf("Tidak menyertakan user : %s", data)
	request, err = http.NewRequest("POST", "http://localhost:8000/api/transactions/", bytes.NewBuffer([]byte(data)))
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusUnauthorized, resp.StatusCode)
	}
	resp.Body.Close()
	// end test case
	// test case transaksi wallet
	t.Logf("Membuat transaksi wallet : %s", data)
	request, err = http.NewRequest("POST", "http://localhost:8000/api/transactions/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusCreated, resp.StatusCode)
	}
	resp.Body.Close()
	// end test case
	// test case transaksi produk
	prodStr := strconv.Itoa(int(prod.ID))
	data = `{
		"type" : "product",
		"products" : [
			{
				"id" : ` + prodStr + `,
				"stock" : 5
			}
		]
	}`
	t.Logf("Membuat transaksi produk : %s", data)
	request, err = http.NewRequest("POST", "http://localhost:8000/api/transactions/", bytes.NewBuffer([]byte(data)))
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusCreated, resp.StatusCode)
	}
	resp.Body.Close()
	// end test case
}
