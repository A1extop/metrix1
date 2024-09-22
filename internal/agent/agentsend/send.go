package agentsend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/A1extop/metrix1/internal/agent/compress"
	"github.com/A1extop/metrix1/internal/agent/hash"
	js "github.com/A1extop/metrix1/internal/agent/json"
)

func SendMetrics(client *http.Client, serverAddress string, metrics []js.Metrics, key string) error {
	metricsJs, err := json.Marshal(metrics)
	fmt.Println(string(metricsJs))
	if err != nil {
		return fmt.Errorf("error marshalling metrics to JSON: %w", err)
	}

	compressedData, err := compress.CompressData(metricsJs)
	if err != nil {
		return err
	}
	hs, err := hash.SignRequestWithSHA256(compressedData, key)
	if err != nil {
		return fmt.Errorf("error creating hash: %w", err)
	}

	url := fmt.Sprintf("%s/updates/", serverAddress)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(compressedData))
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
		return fmt.Errorf("error response from server: %s", resp.Status)
	}

	return nil
}
