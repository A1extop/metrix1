package http

import (
	"github.com/A1extop/metrix1/internal/storage"
	"github.com/A1extop/metrix1/internal/usecase"
	"github.com/A1extop/metrix1/pkg/validator"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Handler struct {
	storage storage.MetricStorage
}

func NewHandler(storage storage.MetricStorage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) Update(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	metricType := vars["type"]
	metricName := vars["name"]
	metricValue := vars["value"]

	if !validator.ValidateRequest(w, req, "text/plain", metricName) {
		return
	}

	err := usecase.UpdateMetric(h.storage, metricType, metricValue, metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentTime := time.Now().Format(time.RFC1123)
	w.Header().Set("Date", currentTime)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
