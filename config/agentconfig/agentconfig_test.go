package agentconfig

import (
	"os"
	"testing"
)

func TestGetParametersEnvironmentVariables(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("ADDRESS", "localhost:9090")
	os.Setenv("KEY", "key")
	os.Setenv("POLL_INTERVAL", "400")
	os.Setenv("REPORT_INTERVAL", "200")
	os.Setenv("RATE_LIMIT", "300")

	parameters := NewParameters()
	parameters.GetParametersEnvironmentVariables()

	if parameters.AddressHTTP != "localhost:9090" {
		t.Errorf("Expected AddressHTTP to be 'localhost:9090', got '%s'", parameters.AddressHTTP)
	}
	if parameters.Key != "key" {
		t.Errorf("Expected Key to be 'key', got '%s'", parameters.Key)
	}
	if parameters.ReportInterval != 200 {
		t.Errorf("Expected ReportInterval to be 300, got '%d'", parameters.ReportInterval)
	}
	if parameters.PollInterval != 400 {
		t.Errorf("Expected PollInterval to be 400, got '%d'", parameters.PollInterval)
	}
	if parameters.RateLimit != 300 {
		t.Errorf("Expected RateLimit to be 300, got '%d'", parameters.RateLimit)
	}
}

func TestGetParameters(t *testing.T) {

	os.Args = []string{
		"program",
		"-a", "localhost:9090",
		"-r", "20",
		"-p", "30",
		"-k", "key",
		"-l", "2",
	}

	parameters := NewParameters()
	parameters.GetParameters()

	if parameters.AddressHTTP != "localhost:9090" {
		t.Errorf("Expected AddressHTTP to be 'localhost:9090', got '%s'", parameters.AddressHTTP)
	}
	if parameters.ReportInterval != 20 {
		t.Errorf("Expected StoreInterval to be 500, got '%d'", parameters.ReportInterval)
	}
	if parameters.PollInterval != 30 {
		t.Errorf("Expected FileStoragePath to be '/tmp/storage', got '%d'", parameters.PollInterval)
	}
	if parameters.Key != "key" {
		t.Errorf("Expected AddrDB to be 'storage', got '%s'", parameters.Key)
	}
	if parameters.RateLimit != 2 {
		t.Errorf("Expected Key to be 'key', got '%d'", parameters.RateLimit)
	}
}
