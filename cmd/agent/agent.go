package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	config "github.com/A1extop/metrix1/config/agentconfig"
	"github.com/A1extop/metrix1/internal/agent/storage"
	uprep "github.com/A1extop/metrix1/internal/agent/updatereportmetrics"
)

var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
)

func main() {
	if BuildVersion == "" {
		BuildVersion = "N/A"
	}
	if BuildDate == "" {
		BuildDate = "N/A"
	}
	if BuildCommit == "" {
		BuildCommit = "N/A"
	}
	storage := storage.NewMemStorage()
	action := uprep.NewAction(storage)
	parameters := config.NewParameters()
	parameters.GetParameters()
	parameters.GetParametersEnvironmentVariables()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	log.Printf("Starting server on port %s", parameters.AddressHTTP)
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
	go func() {
		action.Action(ctx, parameters)
	}()
	<-ctx.Done()
}
