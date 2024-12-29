package updatereportmetrics

import (
	"context"
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

func (u *Updater) Action(ctx context.Context, parameters *config.Parameters) {
	pollTicker := time.NewTicker(time.Duration(parameters.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(parameters.ReportInterval) * time.Second)
	client := &http.Client{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
				go u.updater.UpdateMetrics()
			}
		}
	}()
	if parameters.RateLimit == 0 {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-reportTicker.C:

					parameters.AddressHTTP = "http://" + parameters.AddressHTTP
					u.updater.ReportMetrics(client, parameters)
				}
			}
		}()
	} else if parameters.RateLimit > 0 {
		u.updater.Report(ctx, client, parameters, reportTicker) //parameters.AddressHTTP, parameters.Key, parameters.RateLimit
	}
}
