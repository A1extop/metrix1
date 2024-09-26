package updatereportmetrics

import (
	"context"
	"net/http"
	"sync"
	"time"

	config "github.com/A1extop/metrix1/config/agentconfig"
	"github.com/A1extop/metrix1/internal/agent/storage"
)

func NewAction(updater storage.MetricUpdater, metricsCh chan struct{}) *Updater {
	return &Updater{
		updater:   updater,
		metricsCh: metricsCh,
	}
}

type Updater struct {
	metricsCh chan struct{}
	updater   storage.MetricUpdater
}

func (u *Updater) Action(ctx context.Context, parameters *config.Parameters) {
	pollTicker := time.NewTicker(time.Duration(parameters.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(parameters.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	client := &http.Client{}
	jobs := make(chan struct{})
	rateLimit := make(chan struct{}, parameters.RateLimit)
	var wg sync.WaitGroup

	for i := 0; i < parameters.RateLimit; i++ {
		wg.Add(1)
		go worker(ctx, &wg, jobs, client, rateLimit, u, "http://"+parameters.AddressHTTP, parameters.Key)
	}

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

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			case <-reportTicker.C:
				jobs <- struct{}{}
			}
		}
	}()

	wg.Wait()
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan struct{}, client *http.Client, rateLimit chan struct{}, u *Updater, url string, key string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-jobs:
			rateLimit <- struct{}{}
			go func() {
				defer func() { <-rateLimit }()
				u.updater.ReportMetrics(client, url, key)
			}()
		}
	}
}
