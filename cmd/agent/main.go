package main

import (
	"net/http"
	"time"

	"github.com/A1extop/metrix1/internal/storage"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://localhost:8080"
)

func main() {
	var storage storage.MetricUpdater = storage.NewMemStorage()

	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)
	client := &http.Client{}
	for {
		select {
		case <-pollTicker.C:
			storage.UpdateMetrics()
		case <-reportTicker.C:
			storage.ReportMetrics(client, serverAddress)
		}
	}
}
