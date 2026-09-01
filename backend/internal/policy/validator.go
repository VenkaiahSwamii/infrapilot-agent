package policy

import (
	"fmt"

	"infrapilot/backend/internal/logger"
)

type EndpointTrustInfo struct {
	AgentVersion      string
	APIKey            string
	OrganizationID    string
	DeviceFingerprint string
	mTLSVerified      bool
	HealthStatus      string
}

// ValidateEndpointTrust validates agent authenticity before accepting enrollment
func ValidateEndpointTrust(info EndpointTrustInfo) (bool, error) {
	if info.APIKey == "" {
		return false, fmt.Errorf("missing registration API key")
	}

	if info.AgentVersion == "" {
		return false, fmt.Errorf("missing agent version header")
	}

	if info.HealthStatus == "unhealthy" || info.HealthStatus == "compromised" {
		return false, fmt.Errorf("endpoint health status '%s' fails trust policy", info.HealthStatus)
	}

	logger.Info("Endpoint trust validated", "fingerprint", info.DeviceFingerprint, "org", info.OrganizationID)
	return true, nil
}
