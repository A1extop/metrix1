package main

import (
	"flag"
	"log"
	"net/http"

	http2 "github.com/A1extop/metrix1/internal/http"
	"github.com/A1extop/metrix1/internal/storage"
)

func main() {
	newStorage := storage.NewMemStorage()
	handler := http2.NewHandler(newStorage)
	router := http2.NewRouter(handler)
	lis := flag.String("a", "localhost:8080", "address HTTP")
	flag.Parse()
	log.Printf("Starting server on port %s", *lis)
	err := http.ListenAndServe(*lis, router)
	if err != nil {
		log.Fatal(err)
	}
}
