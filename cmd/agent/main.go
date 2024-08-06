package main

import (
	"net/http"
	"time"

	config "github.com/A1extop/metrix1/config/agentserver"
	"github.com/A1extop/metrix1/internal/storage"
)

func main() {
	var storage storage.MetricUpdater = storage.NewMemStorage()
	config.Init()
	addr := config.ListenAgentConfig()
	reportInterval := config.ReportInterval()
	pollInterval := config.PollInterval()

	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	client := &http.Client{}
	for {
		select {
		case <-pollTicker.C:
			storage.UpdateMetrics()
		case <-reportTicker.C:
			storage.ReportMetrics(client, addr)
		}
	}
}
