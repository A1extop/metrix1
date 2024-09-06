package main

import (
	"log"
	"net/http"

	config "github.com/A1extop/metrix1/config/serverconfig"
	"github.com/A1extop/metrix1/internal/server/data"
	http2 "github.com/A1extop/metrix1/internal/server/http"
	"github.com/A1extop/metrix1/internal/server/storage"
	psql "github.com/A1extop/metrix1/internal/server/store/postgrestore"
)

func main() {
	newStorage := storage.NewMemStorage()
	handler := http2.NewHandler(newStorage)

	parameters := config.NewParameters()
	parameters.GetParameters()
	parameters.GetParametersEnvironmentVariables()

	db, err := psql.ConnectDB(parameters.AddrDB)
	if err != nil {
		log.Println("Failed to connect to database at startup:", err)
	}

	store := psql.NewStore(db)
	repos := psql.NewRepository(store)

	router := http2.NewRouter(handler, repos)

	if parameters.Restore {
		data.ReadingFromDisk(parameters.FileStoragePath, newStorage)
	}
	go data.WritingToDisk(parameters.StoreInterval, parameters.FileStoragePath, newStorage)

	log.Printf("Starting server on port %s", parameters.AddressHTTP)
	err = http.ListenAndServe(parameters.AddressHTTP, router)
	if err != nil {
		log.Fatal(err)
	}
}
