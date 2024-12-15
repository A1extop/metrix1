package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/A1extop/metrix1/config/serverconfig"

	_ "net/http/pprof"

	"github.com/A1extop/metrix1/internal/server/data"
	http2 "github.com/A1extop/metrix1/internal/server/http"
	"github.com/A1extop/metrix1/internal/server/storage"
	psql "github.com/A1extop/metrix1/internal/server/store/postgrestore"
)

var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {
	if BuildVersion == "" {
		BuildVersion = "N/A"
	}
	if BuildDate == "" {
		BuildDate = "N/A"
	}
	if BuildCommit == "" {
		BuildCommit = "N/A"
	}
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
	if db != nil {
		psql.CreateOrConnectTable(db)
		go psql.WritingToBD(repos, parameters.StoreInterval, parameters.AddrDB, newStorage)
	}

	router := http2.NewRouter(handler, repos, parameters)

	if parameters.Restore {
		data.ReadingFromDisk(parameters.FileStoragePath, newStorage)
	}
	go data.WritingToDisk(parameters.StoreInterval, parameters.FileStoragePath, newStorage)

	log.Printf("Starting server on port %s", parameters.AddressHTTP)
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
	err = http.ListenAndServe(parameters.AddressHTTP, router)
	if err != nil {
		log.Println(err)
		return
	}
}
