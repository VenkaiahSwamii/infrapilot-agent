package client

import (
	"net/http"
	"time"

	"infrapilot/agent/internal/config"
	"infrapilot/agent/internal/security"
)

// HTTPClient is the shared HTTP client for all agent requests
var HTTPClient *http.Client

func InitHTTPClient() error {
	cfg := config.Get()

	if !cfg.TLSEnabled {
		HTTPClient = &http.Client{Timeout: 10 * time.Second}
		return nil
	}

	certPath := cfg.CertPath
	keyPath := cfg.KeyPath
	caPath := cfg.CAPath

	clientCert, err := security.LoadClientCert(certPath, keyPath)
	if err != nil {
		return err
	}

	rootCAs, err := security.LoadRootCAs(caPath)
	if err != nil {
		return err
	}

	tlsConfig := security.NewTLSConfig(clientCert, rootCAs)

	HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return nil
}
