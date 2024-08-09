package main

import (
	config "github.com/A1extop/metrix1/config/agentconfig"
	"github.com/A1extop/metrix1/internal/agent/storage"
	uprep "github.com/A1extop/metrix1/internal/agent/updatereportmetrics"
)

func main() {
	storage := storage.NewMemStorage()
	action := uprep.NewAction(storage)
	parameters := config.NewParameters()
	parameters.GetParameters()
	parameters.GetParametersEnvironmentVariables()
	action.Action(parameters)
}
