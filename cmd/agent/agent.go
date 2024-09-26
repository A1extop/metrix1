package main

import (
	"context"

	config "github.com/A1extop/metrix1/config/agentconfig"
	"github.com/A1extop/metrix1/internal/agent/storage"
	uprep "github.com/A1extop/metrix1/internal/agent/updatereportmetrics"
)

func main() {
	storage := storage.NewMemStorage()
	ch := make(chan struct{})
	action := uprep.NewAction(storage, ch)
	parameters := config.NewParameters()
	parameters.GetParameters()
	parameters.GetParametersEnvironmentVariables()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	action.Action(ctx, parameters)
	select {}
}
