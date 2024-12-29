package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/A1extop/metrix1/config/serverconfig"

	"context"
	"github.com/A1extop/metrix1/internal/server/data"
	http2 "github.com/A1extop/metrix1/internal/server/http"
	"github.com/A1extop/metrix1/internal/server/storage"
	psql "github.com/A1extop/metrix1/internal/server/store/postgrestore"
	_ "net/http/pprof"
)

var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	defer stop()
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
		log.Println("failed to connect to database at startup:", err)
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

	srv := &http.Server{
		Addr:    parameters.AddressHTTP,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen on %s: %v\n", parameters.AddressHTTP, err)
		}
	}()
	<-ctx.Done()
	stop()
	log.Println("Shutting down graceful")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
