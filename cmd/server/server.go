package main

import (
	"log"
	"net/http"

	config "github.com/A1extop/metrix1/config/serverconfig"
	http2 "github.com/A1extop/metrix1/internal/http"
	"github.com/A1extop/metrix1/internal/storage"
)

func main() {
	newStorage := storage.NewMemStorage()
	handler := http2.NewHandler(newStorage)
	router := http2.NewRouter(handler)
	parameters := config.NewParameters()
	parameters.GetParameters()
	parameters.GetParametersEnvironmentVariables()
	log.Printf("Starting server on port %s", parameters.AddressHTTP)
	err := http.ListenAndServe(parameters.AddressHTTP, router)
	if err != nil {
		log.Fatal(err)
	}
}
