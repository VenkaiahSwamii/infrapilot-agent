package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func LoadClientCert(certPath, keyPath string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("load client cert/key: %w", err)
	}
	return &cert, nil
}

func LoadRootCAs(caPath string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read CA cert: %w", err)
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("append CA cert failed")
	}

	return pool, nil
}

func NewTLSConfig(cert *tls.Certificate, rootCAs *x509.CertPool) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		RootCAs:      rootCAs,
		MinVersion:   tls.VersionTLS13,
	}
}
