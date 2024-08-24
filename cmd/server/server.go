package main

import (
	"log"
	"net/http"

	config "github.com/A1extop/metrix1/config/serverconfig"
	"github.com/A1extop/metrix1/internal/server/data"
	http2 "github.com/A1extop/metrix1/internal/server/http"
	"github.com/A1extop/metrix1/internal/server/storage"
)

func main() {
	newStorage := storage.NewMemStorage()
	handler := http2.NewHandler(newStorage)
	router := http2.NewRouter(handler)
	parameters := config.NewParameters()
	parameters.GetParameters()
	parameters.GetParametersEnvironmentVariables()
	if parameters.Restore == true {
		data.ReadingFromDisk(parameters.FileStoragePath, newStorage)
	}
	go data.WritingToDisk(parameters.StoreInterval, parameters.FileStoragePath, newStorage)

	log.Printf("Starting server on port %s", parameters.AddressHTTP)
	err := http.ListenAndServe(parameters.AddressHTTP, router)
	if err != nil {
		log.Fatal(err)
	}
}
