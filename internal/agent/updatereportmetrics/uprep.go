package updatereportmetrics

import (
	"net/http"
	"time"

	config "github.com/A1extop/metrix1/config/agentconfig"
	"github.com/A1extop/metrix1/internal/agent/storage"
)

func NewAction(updater storage.MetricUpdater) *Updater {
	return &Updater{updater: updater}
}

type Updater struct {
	updater storage.MetricUpdater
}

func (u *Updater) Action(parameters *config.Parameters) {
	semaphore := make(chan struct{}, parameters.RateLimit)
	pollTicker := time.NewTicker(time.Duration(parameters.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(parameters.ReportInterval) * time.Second)
	client := &http.Client{}
	go func() {
		for {
			<-pollTicker.C
			go u.updater.UpdateMetrics()
		}
	}()

	go func() {
		for {
			semaphore <- struct{}{}
			<-reportTicker.C
			u.updater.ReportMetrics(semaphore, client, "http://"+parameters.AddressHTTP)
		}
	}()
}
