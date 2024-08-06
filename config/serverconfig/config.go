package serverconfig

import (
	"flag"
)

func ListenServerConfig() string {
	lis := flag.String("a", "localhost:8080", "address HTTP")
	flag.Parse()
	return *lis
}
