package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type MetricStorage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}
func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.gauges[name] = value
}

func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.counters[name] += value
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	value, exists := m.gauges[name]
	return value, exists
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	value, exists := m.counters[name]
	return value, exists
}

func validate(w http.ResponseWriter, req *http.Request, expectedContentType, metricName string) bool {
	if req.Header.Get("Content-Type") != expectedContentType {
		http.Error(w, "Неверный Content-Type", http.StatusUnsupportedMediaType)
		return false
	}

	if metricName == "" {
		http.Error(w, "Пустое metricName", http.StatusNotFound)
		return false
	}
	return true
}

var storage MetricStorage = NewMemStorage()

func action(metricType, metricValue, metricName string) error {
	var resErr error
	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			resErr = fmt.Errorf("invalid metricValue")
		} else {
			storage.UpdateGauge(metricName, value)
		}
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			resErr = fmt.Errorf("invalid metricValue")
		} else {
			storage.UpdateCounter(metricName, value)
		}
	default:
		resErr = fmt.Errorf("invalid metricType")
	}
	return resErr
}

func Update(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	metricType := vars["type"]

	metricName := vars["name"]
	if !validate(w, req, "text/plain", metricName) {
		return
	}
	metricValue := vars["value"]
	err := action(metricType, metricValue, metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	currentTime := time.Now().Format(time.RFC1123)
	w.Header().Set("Date", currentTime)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

type starter struct{}

func New() *starter {
	return &starter{}
}

func (s *starter) start() {
	mu := mux.NewRouter()
	mu.HandleFunc("/update/{type}/{name}/{value}", Update).Methods(http.MethodPost)

	http.ListenAndServe(":8080", mu)
}
func main() {
	starter := New()
	starter.start()
}
