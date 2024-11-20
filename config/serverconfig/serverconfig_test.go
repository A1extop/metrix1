package serverconfig

import (
	"os"
	"testing"
)

func TestGetParametersEnvironmentVariables(t *testing.T) {

	os.Setenv("ADDRESS", "localhost:9090")
	os.Setenv("KEY", "key")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/storage")
	os.Setenv("RESTORE", "true")
	os.Setenv("DATABASE_DSN", "storage")
	os.Setenv("STORE_INTERVAL", "300")

	parameters := NewParameters()
	parameters.GetParametersEnvironmentVariables()

	if parameters.AddressHTTP != "localhost:9090" {
		t.Errorf("Expected AddressHTTP to be 'localhost:9090', got '%s'", parameters.AddressHTTP)
	}
	if parameters.Key != "key" {
		t.Errorf("Expected Key to be 'key', got '%s'", parameters.Key)
	}
	if parameters.StoreInterval != 300 {
		t.Errorf("Expected StoreInterval to be 300, got '%d'", parameters.StoreInterval)
	}
	if parameters.FileStoragePath != "/tmp/storage" {
		t.Errorf("Expected FileStoragePath to be '/tmp/storage', got '%s'", parameters.FileStoragePath)
	}
	if parameters.Restore != true {
		t.Errorf("Expected Restore to be true, got '%v'", parameters.Restore)
	}
	if parameters.AddrDB != "storage" {
		t.Errorf("Expected AddrDB to be 'storage', got '%s'", parameters.AddrDB)
	}
}

func TestGetParameters(t *testing.T) {

	os.Args = []string{
		"program",
		"-a", "localhost:9090",
		"-i", "500",
		"-f", "/tmp/storage",
		"-d", "storage",
		"-k", "key",
	}

	parameters := NewParameters()
	parameters.GetParameters()

	if parameters.AddressHTTP != "localhost:9090" {
		t.Errorf("Expected AddressHTTP to be 'localhost:9090', got '%s'", parameters.AddressHTTP)
	}
	if parameters.StoreInterval != 500 {
		t.Errorf("Expected StoreInterval to be 500, got '%d'", parameters.StoreInterval)
	}
	if parameters.FileStoragePath != "/tmp/storage" {
		t.Errorf("Expected FileStoragePath to be '/tmp/storage', got '%s'", parameters.FileStoragePath)
	}
	if parameters.AddrDB != "storage" {
		t.Errorf("Expected AddrDB to be 'storage', got '%s'", parameters.AddrDB)
	}
	if parameters.Key != "key" {
		t.Errorf("Expected Key to be 'key', got '%s'", parameters.Key)
	}
}
