package storage

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"

	send "github.com/A1extop/metrix1/internal/agentsend"
)

type MetricUpdater interface {
	updateRuntimeMetrics()
	updateCustomMetrics()
	UpdateMetrics()
	ReportMetrics(client *http.Client, serverAddress string)
}

func (m *MemStorage) updateRuntimeMetrics() { ////////
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

func (m *MemStorage) updateCustomMetrics() { ///////
	m.counters["PollCount"]++
	m.gauges["RandomValue"] = rand.Float64()
}

func (m *MemStorage) UpdateMetrics() { /////////
	m.updateRuntimeMetrics()
	m.updateCustomMetrics()
}

func (m *MemStorage) ReportMetrics(client *http.Client, serverAddress string) {
	for name, value := range m.gauges {
		err := send.SendMetric(client, serverAddress, "gauge", name, value)
		if err != nil {
			fmt.Println(err) //оставил, чтобы отслеживать, что происходит
		}
	}
	for name, value := range m.counters {
		err := send.SendMetric(client, serverAddress, "counter", name, value)
		if err != nil {
			fmt.Println(err) //оставил, чтобы отслеживать, что происходит
		}
	}
}
