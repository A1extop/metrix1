package main

import (
	http2 "github.com/A1extop/metrix1/internal/http"
	"github.com/A1extop/metrix1/internal/storage"
	"log"
	"net/http"
)

func main() {
	newStorage := storage.NewMemStorage()
	handler := http2.NewHandler(newStorage)
	router := http2.NewRouter(handler)

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
