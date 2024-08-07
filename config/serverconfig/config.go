package serverconfig

import (
	"flag"
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
