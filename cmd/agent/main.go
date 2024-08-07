package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/A1extop/metrix1/internal/storage"
)

func main() {
	var storage storage.MetricUpdater = storage.NewMemStorage()
	addr := flag.String("a", "localhost:8080", "адрес HTTP-сервера")
	reportInterval := flag.Int("r", 10, "частота отправки метрик на сервер в секундах")
	pollInterval := flag.Int("p", 2, "частота опроса метрик из пакета runtime в секундах")
	flag.Parse()

	pollTicker := time.NewTicker(time.Duration(*pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(*reportInterval) * time.Second)
	client := &http.Client{}
	for {
		select {
		case <-pollTicker.C:
			storage.UpdateMetrics()
		case <-reportTicker.C:
			storage.ReportMetrics(client, "http://"+*addr)
		}
	}
}
