package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a Impl

var j *sql.DB

func TestMain(m *testing.M) {
	setEnvs()

	a := Impl{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_DB_SSLMODE"))

	// Checking to ensure our table is already existing
	// ensureTableExists()
	a.DB.HasTable(&Product{})

	code := m.Run()

	a.DB.Model(&Product{}).Delete(&Product{})
	// clearTable()
	// a.DB.Exec("DELETE FROM products")
	// a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")

	os.Exit(code)
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS
(
id SERIAL,
name TEXT NOT NULL,
price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
CONSTRAINT products_pkey PRIMARY KEY (id)
);`

func ensureTableExists() {
	j.Exec(tableCreationQuery)
}

func clearTable() {
	j.Exec("DELETE FROM products")
	j.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
	// clearTable()
	// a.DB.Model(&Product{}).Delete(&Product{})

	a.DB.Debug().Exec("DELETE FROM products")
	a.DB.Debug().Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	// clearTable()

	a.DB.Delete(&Product{})

	// a.DB.Raw("DELETE FROM products", nil)
	// a.DB.Raw("ALTER SEQUENCE products_id_seq RESTART WITH 1")

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {
	// clearTable()

	a.DB.Delete(&Product{})

	// a.DB.Raw("DELETE FROM products")
	// a.DB.Raw("ALTER SEQUENCE products_id_seq RESTART WITH 1")

	payload := []byte(`{"name":"test product","price":11.22}`)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	// clearTable()
	a.DB.Delete(&Product{})

	// a.DB.Raw("DELETE FROM products")
	// a.DB.Raw("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	addProducts()

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
func addProducts() {
	a.DB.Create(Product{Name: "L1212", Price: 1000})
}

func TestUpdateProduct(t *testing.T) {
	// clearTable()
	a.DB.Delete(&Product{})

	// a.DB.Raw("DELETE FROM products")
	// a.DB.Raw("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	addProducts()

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	payload := []byte(`{"name":"test product - updated name","price":11.22}`)

	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	// clearTable()
	a.DB.Delete(&Product{})

	// a.DB.Raw("DELETE FROM products")
	// a.DB.Raw("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	addProducts()

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
