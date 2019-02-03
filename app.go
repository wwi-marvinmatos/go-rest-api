package main

// app.go

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Impl Data Structure to define our router and database
type Impl struct {
	DB      *gorm.DB
	Product Product
	Router  mux.Router
}

// This unexported struct will define the needed data for the product that we will use to interact with the db
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&Product{})
}

//sql  open from below ("postgres", "user=test password=test dbname=test sslmode=disable")
// we are using sslmode=disabled to prevent the test from failing(go test are failing because ssl is not configured)

// Initialize method to initialize the Database connection utilizing a connection string and the sql driver to connect
// this will also create the necessary new mux router to register the proper routes
func (i *Impl) Initialize(user, password, dbname, sslmode string) {
	var err error

	// connStr := "host=localhost port=5432 user=marvin dbname=grom sslmode=disable"
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)

	i.DB, err = gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	i.DB.LogMode(true)
}

// GET Product handler
func (i *Impl) getProduct(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	product := Product{}
	if i.DB.First(&product, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(&product)
}

// GET Products handler
func (i *Impl) getProducts(w rest.ResponseWriter, r *rest.Request) {
	products := []Product{}
	i.DB.Find(&products)
	w.WriteJson(&products)
}

// CREATE product handler
func (i *Impl) createProduct(w rest.ResponseWriter, r *rest.Request) {
	product := Product{}
	if err := r.DecodeJsonPayload(&product); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := i.DB.Save(&product).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&product)
}

// UPDATE product handler
func (i *Impl) updateProduct(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	product := Product{}
	if i.DB.First(&product, id).Error != nil {
		rest.NotFound(w, r)
		return
	}

	updated := Product{}
	if err := r.DecodeJsonPayload(&updated); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.Name = updated.Name

	if err := i.DB.Save(&product).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&product)
}

// DELETE product handler
func (i *Impl) deleteProduct(w rest.ResponseWriter, r *rest.Request) {

	id := r.PathParam("id")
	product := Product{}
	if i.DB.First(&product, id).Error != nil {
		rest.NotFound(w, r)
		return
	}
	if err := i.DB.Delete(&product).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	// vars := mux.Vars(r)
	// id, err := strconv.Atoi(vars["id"])
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
	// 	return
	// }

	// p := product{ID: id}
	// if err := p.deleteProduct(a.DB); err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Helper functions for the handlers
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
