package agentconfig

import (
	"flag"
	"os"
	"strconv"
)

type Parameters struct {
	AddressHTTP    string
	ReportInterval int
	PollInterval   int
	Key            string
	RateLimit      int
}

func NewParameters() *Parameters {
	return &Parameters{
		AddressHTTP:    "",
		ReportInterval: 0,
		PollInterval:   0,
		Key:            "",
		RateLimit:      0,
	}
}

func (p *Parameters) GetParameters() {
	addr := flag.String("a", "localhost:8080", "address HTTP")
	reportInterval := flag.Int("r", 10, "frequency of sending metrics to the server in seconds")
	pollInterval := flag.Int("p", 2, "frequency of polling metrics from the runtime package in seconds")
	key := flag.String("k", "", "hash key")
	rateLimit := flag.Int("l", 1, "number of goroutines")

	flag.Parse()
	p.AddressHTTP = *addr
	p.PollInterval = *pollInterval
	p.ReportInterval = *reportInterval
	p.Key = *key
	p.RateLimit = *rateLimit
}
func (p *Parameters) GetParametersEnvironmentVariables() {
	addr := os.Getenv("ADDRESS")
	if addr != "" {
		p.AddressHTTP = addr
	}
	repIntervalStr := os.Getenv("REPORT_INTERVAL")
	repIntervalInt, _ := strconv.Atoi(repIntervalStr)
	if repIntervalStr != "" {
		p.ReportInterval = repIntervalInt
	}
	pollIntervalStr := os.Getenv("POLL_INTERVAL")
	pollIntervalInt, _ := strconv.Atoi(pollIntervalStr)
	if pollIntervalStr != "" {
		p.PollInterval = pollIntervalInt
	}

	key := os.Getenv("KEY")
	if key != "" {
		p.Key = key
	}
	rateLimitStr := os.Getenv("RATE_LIMIT")
	rateLimit, _ := strconv.Atoi(rateLimitStr)
	if pollIntervalStr != "" {
		p.RateLimit = rateLimit
	}
}
