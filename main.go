package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
)

func setEnvs() {
	os.Setenv("APP_DB_USERNAME", "marvin")
	os.Setenv("APP_DB_PASSWORD", "test")
	os.Setenv("APP_DB_NAME", "api_server")
	os.Setenv("APP_DB_SSLMODE", "disable")
}

func main() {
	setEnvs()

	i := Impl{}
	i.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_DB_SSLMODE"))

	i.InitSchema()
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/products", i.getProducts),
		rest.Post("/product", i.createProduct),
		rest.Get("/product/{id:[0-9]+}", i.getProduct),
		rest.Put("/product/{id:[0-9]+}", i.updateProduct),
		rest.Delete("/product/{id:[0-9]+}", i.updateProduct),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8000", api.MakeHandler()))
}
