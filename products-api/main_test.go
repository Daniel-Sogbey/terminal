package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(host, user, password, dbname, port)
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestEmptytable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest(http.MethodGet, "/products", nil)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	t.Errorf("%v", response.Body)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest(http.MethodGet, "/product/11", nil)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string

	json.Unmarshal(response.Body.Bytes(), &m)

	if m["error"] != "Product not found" {
		t.Errorf("Expected 'error' key of response to be set to 'Product not found'. Got %s", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"name":"Test Product", "price":11.99}`)

	req, _ := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Test Product" {
		t.Errorf("Expected product name to be 'Test Product'. Got %v", m["name"])
	}

	if m["price"] != 11.99 {
		t.Errorf("Expected price to be '11.99'. Got %v", m["price"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected ID to be '1.0'. Got %v", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest(http.MethodGet, "/product/1", nil)

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest(http.MethodGet, "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price":11.99}`)

	req, _ = http.NewRequest(http.MethodPut, "/product/1", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected id to remain the same %v. Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] != originalProduct["name"] {
		t.Errorf("Expected name to change from %v to %v. Got %v", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] != originalProduct["price"] {
		t.Errorf("Expected price to change from %v to %v. Got %v", originalProduct["price"], m["price"], m["price"])
	}

}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest(http.MethodDelete, "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest(http.MethodDelete, "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest(http.MethodGet, "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, r)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
	id SERIAL,
	name TEXT NOT NULL,
	price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
	CONSTRAINT products_pkey PRIMARY KEY (id)
)`

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products(name,price) VALUE ($1,$2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}
