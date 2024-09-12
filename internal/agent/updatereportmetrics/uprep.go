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

	pollTicker := time.NewTicker(time.Duration(parameters.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(parameters.ReportInterval) * time.Second)
	client := &http.Client{}
	for {
		select {
		case <-pollTicker.C:
			u.updater.UpdateMetrics()
		case <-reportTicker.C:

			u.updater.ReportMetrics(client, "http://"+parameters.AddressHTTP)
		}
	}
}
