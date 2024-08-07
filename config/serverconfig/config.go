package serverconfig

import (
	"flag"
)

var lis *string

func Init() {
	lis = flag.String("a", "localhost:8080", "address HTTP")
	flag.Parse()
}

func ListenServerConfig() string {
	return *lis
}
