package agentsend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/A1extop/metrix1/config/agentconfig"
	"io"
	"net/http"

	"github.com/A1extop/metrix1/internal/agent/compress"
	"github.com/A1extop/metrix1/internal/agent/hash"
	js "github.com/A1extop/metrix1/internal/agent/json"
)

func send(client *http.Client, serverAddress, path string, data []byte, key string, encryptionKey string) error {
	compressedData, err := compress.CompressData(data)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s%s", serverAddress, path)
	hs, err := hash.SignRequestWithSHA256(compressedData, key)

	if err != nil {
		return fmt.Errorf("error creating hash: %w", err)
	}

	encryptionData := make([]byte, 0)
	if encryptionKey != "" {
		encryptionData, err = hash.Encrypt(compressedData, encryptionKey)
		if err != nil {
			return err
		}
	}
	if len(encryptionData) == 0 {
		encryptionData = compressedData
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(encryptionData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("HashSHA256", hs)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error response from server: %s, failed to read response body: %v", resp.Status, err)
		}
		return fmt.Errorf("error response from server: %s, response body: %s", resp.Status, string(body))
	}

	return nil
}
func SendMetric(client *http.Client, metric js.Metrics, parameters *agentconfig.Parameters) error {
	metricJs, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("error marshalling metrics to JSON: %w", err)
	}
	return send(client, parameters.AddressHTTP, "/update/", metricJs, parameters.Key, parameters.CryptoKey)
}
func SendMetrics(client *http.Client, metrics []js.Metrics, parameters *agentconfig.Parameters) error {
	metricsJs, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("error marshalling metrics to JSON: %w", err)
	}
	return send(client, parameters.AddressHTTP, "/updates/", metricsJs, parameters.Key, parameters.CryptoKey)
}
