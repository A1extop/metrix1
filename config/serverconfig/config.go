package serverconfig

import (
	"flag"
	"os"
)

type Parameters struct {
	AddressHTTP string
}

func NewParameters() *Parameters {
	return &Parameters{AddressHTTP: ""}
}

func (p *Parameters) GetParameters() {
	addr := flag.String("a", "localhost:8080", "address HTTP")
	flag.Parse()
	p.AddressHTTP = *addr
}
func (p *Parameters) GetParametersEnvironmentVariables() {
	addr := os.Getenv("ADDRESS")
	if addr != "" {
		p.AddressHTTP = addr
	}
}
