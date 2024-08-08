package main

import (
	"net/http"
	"time"

	config "github.com/A1extop/metrix1/config/agentconfig"
	"github.com/A1extop/metrix1/internal/agent/storage"
)

func main() {
	var storage storage.MetricUpdater = storage.NewMemStorage()
	parameters := config.NewParameters()
	parameters.GetParameters()
	parameters.GetParametersEnvironmentVariables()
	pollTicker := time.NewTicker(time.Duration(parameters.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(parameters.ReportInterval) * time.Second)
	client := &http.Client{}
	for {
		select {
		case <-pollTicker.C:
			storage.UpdateMetrics()
		case <-reportTicker.C:
			storage.ReportMetrics(client, "http://"+parameters.AddressHTTP)
		}
	}
}
