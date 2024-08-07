package agentconfig

import (
	"flag"
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
	reportInterval := flag.Int("r", 10, "частота отправки метрик на сервер в секундах")
	pollInterval := flag.Int("p", 2, "частота опроса метрик из пакета runtime в секундах")
	flag.Parse()
	p.AddressHTTP = *addr
	p.PollInterval = *pollInterval
	p.ReportInterval = *reportInterval
}
