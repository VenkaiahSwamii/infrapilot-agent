package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"infrapilot/agent/internal/config"
)

func PostWithRetry(url string, body []byte, contentType string) error {
	if HTTPClient == nil {
		_ = InitHTTPClient()
		if HTTPClient == nil {
			HTTPClient = &http.Client{Timeout: 10 * time.Second}
		}
	}

	backoff := []time.Duration{
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
		30 * time.Second,
	}

	for i := 0; i < len(backoff); i++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", contentType)

		// Set API Key headers
		apiKey := config.Get().APIKey
		if apiKey != "" {
			req.Header.Set("X-API-Key", apiKey)
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}

		resp, err := HTTPClient.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}

		fmt.Println("Retrying in", backoff[i])

		time.Sleep(backoff[i])
	}

	return fmt.Errorf("server unreachable")
}
