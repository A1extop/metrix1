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
}

func NewParameters() *Parameters {
	return &Parameters{
		AddressHTTP:    "",
		ReportInterval: 0,
		PollInterval:   0,
	}
}

func (p *Parameters) GetParameters() {
	addr := flag.String("a", "localhost:8080", "address HTTP")
	reportInterval := flag.Int("r", 10, "frequency of sending metrics to the server in seconds")
	pollInterval := flag.Int("p", 2, "frequency of polling metrics from the runtime package in seconds")
	flag.Parse()
	p.AddressHTTP = *addr
	p.PollInterval = *pollInterval
	p.ReportInterval = *reportInterval
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

}
