package agentsend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	js "github.com/A1extop/metrix1/internal/agent/json"
)

func SendMetric(client *http.Client, serverAddress string, metrics *js.Metrics) error {
	metricsJs, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("error marshalling metrics to JSON: %w", err)
	}
	url := fmt.Sprintf("%s/update/", serverAddress)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(metricsJs))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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
