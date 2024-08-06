package sendagent

import (
	"fmt"
	"net/http"
)

func SendMetric(client *http.Client, serverAddress, metricType, metricName string, value interface{}) error {
	url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, metricType, metricName, value)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

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
