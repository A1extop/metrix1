package storage

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	send "github.com/A1extop/metrix1/internal/agent/agentsend"
	js "github.com/A1extop/metrix1/internal/agent/json"
)

type MetricUpdater interface {
	updateRuntimeMetrics()
	updateCustomMetrics()
	UpdateMetrics()
	ReportMetrics(semaphore chan struct{}, client *http.Client, serverAddress string, key string)
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m *MemStorage) updateRuntimeMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.gauges["Alloc"] = float64(memStats.Alloc)
	m.gauges["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.gauges["Frees"] = float64(memStats.Frees)
	m.gauges["GCCPUFraction"] = memStats.GCCPUFraction
	m.gauges["GCSys"] = float64(memStats.GCSys)
	m.gauges["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.gauges["HeapIdle"] = float64(memStats.HeapIdle)
	m.gauges["HeapInuse"] = float64(memStats.HeapInuse)
	m.gauges["HeapObjects"] = float64(memStats.HeapObjects)
	m.gauges["HeapReleased"] = float64(memStats.HeapReleased)
	m.gauges["HeapSys"] = float64(memStats.HeapSys)
	m.gauges["LastGC"] = float64(memStats.LastGC)
	m.gauges["Lookups"] = float64(memStats.Lookups)
	m.gauges["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.gauges["MCacheSys"] = float64(memStats.MCacheSys)
	m.gauges["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.gauges["MSpanSys"] = float64(memStats.MSpanSys)
	m.gauges["Mallocs"] = float64(memStats.Mallocs)
	m.gauges["NextGC"] = float64(memStats.NextGC)
	m.gauges["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.gauges["NumGC"] = float64(memStats.NumGC)
	m.gauges["OtherSys"] = float64(memStats.OtherSys)
	m.gauges["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.gauges["StackInuse"] = float64(memStats.StackInuse)
	m.gauges["StackSys"] = float64(memStats.StackSys)
	m.gauges["Sys"] = float64(memStats.Sys)
	m.gauges["TotalAlloc"] = float64(memStats.TotalAlloc)
}

func (m *MemStorage) updateCustomMetrics() {
	m.counters["PollCount"]++
	m.gauges["RandomValue"] = rand.Float64()
}

func (m *MemStorage) UpdateMetrics() {
	m.updateRuntimeMetrics()
	m.updateCustomMetrics()
}

func (m *MemStorage) ReportMetrics(semaphore chan struct{}, client *http.Client, serverAddress string, key string) {
	defer func() {
		<-semaphore
	}()
	var metrics []js.Metrics

	for name, value := range m.gauges {
		metric := js.NewMetrics()
		metric.ID = name
		metric.Value = &value
		metric.MType = "gauge"
		metrics = append(metrics, *metric)
	}

	for name, value := range m.counters {
		metric := js.NewMetrics()
		metric.ID = name
		metric.MType = "counter"
		metric.Delta = &value
		metrics = append(metrics, *metric)
	}

	if len(metrics) > 0 {
		TimesDuration := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
		targetError := errors.New("error sending request")
		for _, times := range TimesDuration {

			err := send.SendMetrics(client, serverAddress, metrics, key)
			if err == nil {
				break
			}
			if errors.Is(err, targetError) {
				log.Printf("Error encountered: %s. Retrying in %v...", err.Error(), times)
				time.Sleep(times)
			} else {
				log.Printf("Failed to send metrics: %s", err.Error())
				break
			}
		}
	}

	m.counters["PollCount"] = 0
}
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}
