package commands

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"infrapilot/agent/internal/config"
)

// KeyRotationResponse matches the backend rotation response
type KeyRotationResponse struct {
	Rotate    bool   `json:"rotate"`
	APIKey    string `json:"api_key"`
	Version   int    `json:"version"`
	MachineID string `json:"machine_id"`
}

// PollKeyRotation polls the backend for API key rotation commands
func PollKeyRotation(backendURL, machineID, apiKey string, store *config.ConfigStore) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		checkKeyRotation(backendURL, machineID, apiKey, store)
	}
}

func checkKeyRotation(backendURL, machineID, apiKey string, store *config.ConfigStore) {
	client := &http.Client{Timeout: 10 * time.Second}
	reqURL := backendURL + "/api/v1/agent/key"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Printf("[KeyRotation] Error creating request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[KeyRotation] Backend unreachable: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[KeyRotation] Backend returned status: %d", resp.StatusCode)
		return
	}

	var rotationResp KeyRotationResponse
	if err := json.NewDecoder(resp.Body).Decode(&rotationResp); err != nil {
		log.Printf("[KeyRotation] Decode error: %v", err)
		return
	}

	if rotationResp.Rotate && rotationResp.APIKey != "" && rotationResp.APIKey != apiKey {
		log.Println("[KeyRotation] Backend requested key rotation. Updating API key...")
		apiKey = rotationResp.APIKey

		cfg, err := store.Load()
		if err != nil {
			log.Printf("[KeyRotation] Failed to load config for key update: %v", err)
			return
		}
		cfg.APIKey = apiKey
		cfg.KeyVersion = rotationResp.Version
		if err := store.Save(cfg); err != nil {
			log.Printf("[KeyRotation] Failed to save rotated key to config.json: %v", err)
			return
		}
		log.Println("[KeyRotation] API key rotated and saved successfully.")
	}
}

// RequestKeyRotation requests a key rotation from the backend
func RequestKeyRotation(backendURL, machineID, apiKey string, store *config.ConfigStore) bool {
	client := &http.Client{Timeout: 10 * time.Second}
	reqURL := backendURL + "/api/v1/agent/key-rotation"
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		log.Printf("[KeyRotation] Error creating request: %v", err)
		return false
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[KeyRotation] Backend unreachable during rotation: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[KeyRotation] Backend returned status: %d during rotation", resp.StatusCode)
		return false
	}

	var rotationResp KeyRotationResponse
	if err := json.NewDecoder(resp.Body).Decode(&rotationResp); err != nil {
		log.Printf("[KeyRotation] Decode error: %v", err)
		return false
	}

	if rotationResp.APIKey != "" && rotationResp.APIKey != apiKey {
		log.Println("[KeyRotation] Received new API key from backend. Updating config...")
		cfg, err := store.Load()
		if err != nil {
			log.Printf("[KeyRotation] Failed to load config for key update: %v", err)
			return false
		}
		cfg.APIKey = rotationResp.APIKey
		cfg.KeyVersion = rotationResp.Version
		if err := store.Save(cfg); err != nil {
			log.Printf("[KeyRotation] Failed to save rotated key to config.json: %v", err)
			return false
		}
		log.Println("[KeyRotation] API key rotated and saved successfully.")
		return true
	}

	return false
}
