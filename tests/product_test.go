package tests

import (
	"bytes"
	"net/http"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	client := &http.Client{}
	// test case data yang tidak lengkap
	product1 := `{
		"product_name": "Test Name",
		"price":        1000
	}`
	request1, err := http.NewRequest("POST", "http://localhost:8000/api/products/", bytes.NewBuffer([]byte(product1)))
	if err != nil {
		panic(err)
	}
	request1.Header.Add("Content-Type", "application/json")
	t.Logf("Menggunakan data yang tidak lengkap : %s", product1)
	resp, err := client.Do(request1)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case data yang lengkap
	product2 := `{
		"product_name": "Test Name",
		"price":        1000,
		"stock" : 1,
		"visible" : true,
		"description" : "test"
	}`
	request2, err := http.NewRequest("POST", "http://localhost:8000/api/products/", bytes.NewBuffer([]byte(product2)))
	if err != nil {
		panic(err)
	}
	request2.Header.Add("Content-Type", "application/json")
	t.Logf("Menggunakan data yang lengkap : %s", product2)
	resp2, err := client.Do(request2)
	if err != nil {
		panic(err)
	}
	if resp2.StatusCode != http.StatusCreated {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusCreated, resp2.StatusCode)
	}
	resp2.Body.Close()
	// end of test case
}

func TestReadProduct(t *testing.T) {
	client := &http.Client{}
	// test case berhasil
	request, err := http.NewRequest("GET", "http://localhost:8000/api/products/", nil)
	if err != nil {
		panic(err)
	}
	t.Logf("Tidak menggunakan query dan id pada database yang sudah terisi")
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusOK, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case dengan kata kunci yang tidak ada
	request, err = http.NewRequest("GET", "http://localhost:8000/api/products/?q=fasdfsadfasdfasd", nil)
	if err != nil {
		panic(err)
	}
	t.Logf("Menggunakan query yang tidak ada : fasdfsadfasdfasd")
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusNotFound, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case dengan kata kunci yang tidak ada
	request, err = http.NewRequest("GET", "http://localhost:8000/api/products/?id=9980913", nil)
	if err != nil {
		panic(err)
	}
	t.Logf("Menggunakan id yang tidak ada : 9980913")
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusNotFound, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
}

func TestUpdateProduct(t *testing.T) {
	client := &http.Client{}
	// test case tidak menyertakan id
	body := `{}`
	request, err := http.NewRequest("PATCH", "http://localhost:8000/api/products/", bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(err)
	}
	t.Logf("Tidak menyertakan id : %s", body)
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case menggunakan id yang tidak ada
	body = `{
		"product_name" : "Test",
		"price" : 1000,
		"stock" : 1,
		"visible" : true,
		"description" : "test"
	}`
	request, err = http.NewRequest("PATCH", "http://localhost:8000/api/products/?id=100093", bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(err)
	}
	t.Logf("Menggunakan id yang tidak ada : %s", body)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusNotFound, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case data tidak lengkap
	body = `{}`
	request, err = http.NewRequest("PATCH", "http://localhost:8000/api/products/?id=1", bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(err)
	}
	t.Logf("Menggunakan data yang tidak lengkap : %s", body)
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
}

func TestDeleteProduct(t *testing.T) {
	client := &http.Client{}
	// test case tidak menyertakan id
	request, err := http.NewRequest("DELETE", "http://localhost:8000/api/products/", nil)
	if err != nil {
		panic(err)
	}
	t.Logf("Tidak menyertakan id")
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// test case menggunakan id yang tidak ada
	request, err = http.NewRequest("DELETE", "http://localhost:8000/api/products/?id=1038947", nil)
	if err != nil {
		panic(err)
	}
	t.Logf("Menggunakan id yang tidak ada")
	resp, err = client.Do(request)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Test gagal, ekspetasi kode status adalah %d tapi hasilnya %d", http.StatusNotFound, resp.StatusCode)
	}
	resp.Body.Close()
	// end of test case
	// ResetDatabase()
}
