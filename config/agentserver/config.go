package agentconfig

import (
	"flag"
)

var (
	addr           = flag.String("a", "localhost:8080", "адрес HTTP-сервера")
	reportInterval = flag.Int("r", 10, "частота отправки метрик на сервер в секундах")
	pollInterval   = flag.Int("p", 2, "частота опроса метрик из пакета runtime в секундах")
)

func Init() {
	flag.Parse()
}

func ListenAgentConfig() string {
	return *addr
}

func ReportInterval() int {
	return *reportInterval
}

func PollInterval() int {
	return *pollInterval
}
