package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"infrapilot/backend/internal/config"
)

func NewTLSServer(handler http.Handler) (*http.Server, error) {
	cfg := config.Get()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	// Enable MTLS if configured
	if cfg.MTLSEnabled {
		caCert, err := os.ReadFile(cfg.TLSCAPath)
		if err != nil {
			return nil, fmt.Errorf("read CA certificate: %w", err)
		}

		caPool := x509.NewCertPool()
		if ok := caPool.AppendCertsFromPEM(caCert); !ok {
			return nil, fmt.Errorf("append CA certificate to pool failed")
		}

		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientCAs = caPool
	}

	// Load server certificate
	cert, err := tls.LoadX509KeyPair(cfg.TLSCertPath, cfg.TLSKeyPath)
	if err != nil {
		return nil, fmt.Errorf("load server certificate: %w", err)
	}

	tlsConfig.Certificates = []tls.Certificate{cert}

	return &http.Server{
		Addr:      ":" + cfg.TLSPort,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}, nil
}
