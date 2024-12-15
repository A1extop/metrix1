package serverconfig

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Parameters struct {
	AddressHTTP     string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	AddrDB          string
	Key             string
	CryptoKey       string
}

func NewParameters() *Parameters {
	return &Parameters{
		AddressHTTP:     "",
		StoreInterval:   0,
		FileStoragePath: "",
		Restore:         true,
		Key:             "",
		CryptoKey:       "",
	}
}
func (p *Parameters) LoadFromJSONFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(p); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	return nil
}
func (p *Parameters) GetParameters() {
	configPath := flag.String("config", "", "path to JSON config file")
	addr := flag.String("a", "localhost:8080", "address HTTP")
	storeInterval := flag.Int("i", 300, "the time interval in seconds after which the current server readings are saved to disk")
	fileStoragePath := flag.String("f", "", "the path to the file where the current values are saved")
	restore := flag.Bool("r", true, "whether or not to load previously saved values from the specified file when the server starts")
	addrDB := flag.String("d", "", "String with database connection address")
	key := flag.String("k", "", "hash key")
	cryptoKey := flag.String("c", "", "encryption key")

	flag.Parse()
	if *configPath == "" {
		*configPath = os.Getenv("CONFIG")
	}
	if *configPath != "" {
		if err := p.LoadFromJSONFile(*configPath); err != nil {
			log.Printf("Failed to load config from file: %v\n", err)
		}
	}
	p.AddressHTTP = *addr
	p.StoreInterval = *storeInterval
	p.FileStoragePath = *fileStoragePath
	p.Restore = *restore
	p.AddrDB = *addrDB
	p.Key = *key
	p.CryptoKey = *cryptoKey
}
func (p *Parameters) GetParametersEnvironmentVariables() {
	addr := os.Getenv("ADDRESS")
	if addr != "" {
		p.AddressHTTP = addr
	}
	storeIntervalStr := os.Getenv("STORE_INTERVAL")
	if storeIntervalStr != "" {
		storeInterval, err := strconv.Atoi(storeIntervalStr)
		if err != nil {
			log.Printf("Invalid StoreInterval: %v\n", err)
		} else {
			p.StoreInterval = storeInterval
		}
	}
	fileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePath != "" {
		p.FileStoragePath = fileStoragePath
	}
	restoreStr := os.Getenv("RESTORE")
	if restoreStr != "" {
		restore, err := strconv.ParseBool(restoreStr)
		if err != nil {
			log.Printf("Invalid Restore: %v\n", err)
		} else {
			p.Restore = restore
		}
	}
	addrDB := os.Getenv("DATABASE_DSN")
	if addrDB != "" {
		p.AddrDB = addrDB
	}

	key := os.Getenv("KEY")
	if key != "" {
		p.Key = key
	}
	cryptoKey := os.Getenv("CRYPTO_KEY")
	if cryptoKey != "" {
		p.CryptoKey = cryptoKey
	}

}
