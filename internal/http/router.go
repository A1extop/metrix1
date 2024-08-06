package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(handler *Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/update/{type}/{name}/{value}", handler.Update).Methods(http.MethodPost)
	return router
}
