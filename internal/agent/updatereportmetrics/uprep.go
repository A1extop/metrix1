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

	go func() {
		defer pollTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-pollTicker.C:
				go u.updater.UpdateMetrics()
			}
		}
	}()

	go func() {
		workerPool(u, ctx, parameters)
	}()
}

func workerPool(u *Updater, ctx context.Context, parameters *config.Parameters) {
	reportTicker := time.NewTicker(time.Duration(parameters.ReportInterval) * time.Second)
	client := &http.Client{}
	reportChan := make(chan struct{}, parameters.RateLimit)

	for i := 0; i < parameters.RateLimit; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-reportChan:
					u.updater.ReportMetrics(client, "http://"+parameters.AddressHTTP, parameters.Key)
				}
			}
		}()
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-reportTicker.C:
				reportChan <- struct{}{}
			}
		}
	}()
}
