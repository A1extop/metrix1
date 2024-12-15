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
	CryptoKey      string
}

func NewParameters() *Parameters {
	return &Parameters{
		AddressHTTP:    "",
		ReportInterval: 0,
		PollInterval:   0,
		Key:            "",
		RateLimit:      0,
		CryptoKey:      "",
	}
}

func (p *Parameters) GetParameters() {
	addr := flag.String("a", "localhost:8080", "address HTTP")
	reportInterval := flag.Int("r", 10, "frequency of sending metrics to the server in seconds")
	pollInterval := flag.Int("p", 2, "frequency of polling metrics from the runtime package in seconds")
	key := flag.String("k", "", "hash key")
	rateLimit := flag.Int("l", 1, "number of goroutines")
	cryptoKey := flag.String("c", "", "hash key")

	flag.Parse()
	p.AddressHTTP = *addr
	p.PollInterval = *pollInterval
	p.ReportInterval = *reportInterval
	p.Key = *key
	p.RateLimit = *rateLimit
	p.CryptoKey = *cryptoKey
}
func (p *Parameters) GetParametersEnvironmentVariables() {
	addr := os.Getenv("ADDRESS")
	if addr != "" {
		p.AddressHTTP = addr
	}
	repIntervalStr := os.Getenv("REPORT_INTERVAL")
	repIntervalInt, err := strconv.Atoi(repIntervalStr)
	if err == nil {
		p.ReportInterval = 10
	}
	if repIntervalStr != "" {
		p.ReportInterval = repIntervalInt
	}
	pollIntervalStr := os.Getenv("POLL_INTERVAL")
	pollIntervalInt, err := strconv.Atoi(pollIntervalStr)
	if err == nil {
		p.PollInterval = 2
	}
	if pollIntervalStr != "" {
		p.PollInterval = pollIntervalInt
	}

	key := os.Getenv("KEY")
	if key != "" {
		p.Key = key
	}
	rateLimitStr := os.Getenv("RATE_LIMIT")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		p.RateLimit = 1
	}

	if pollIntervalStr != "" {
		p.RateLimit = rateLimit
	}
	cryptoKey := os.Getenv("CRYPTO_KEY")
	if cryptoKey != "" {
		p.CryptoKey = cryptoKey
	}
}
