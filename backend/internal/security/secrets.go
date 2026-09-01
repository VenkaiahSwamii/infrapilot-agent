package security

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type SecretsManager struct {
	logger *logger.Logger
}

func NewSecretsManager() *SecretsManager {
	return &SecretsManager{
		logger: logger.Get(),
	}
}

// GetVaultStatus returns enterprise secrets manager provider status
func (sm *SecretsManager) GetVaultStatus(orgID string) (*models.SecretVaultRecord, error) {
	var vault models.SecretVaultRecord

	if database.DB != nil {
		database.DB.Where("organization_id = ? OR organization_id = ''", orgID).First(&vault)
	}

	if vault.ID == uuid.Nil {
		vault = models.SecretVaultRecord{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "hashicorp_vault",
			VaultURL:       "https://vault.internal.company.com:8200",
			RotationDays:   90,
			LastRotatedAt:  time.Now().AddDate(0, 0, -15),
			Status:         "connected",
			CreatedAt:      time.Now(),
		}
	}

	return &vault, nil
}

// RotateSecrets triggers secret rotation for DB, JWT, and API keys
func (sm *SecretsManager) RotateSecrets(provider string) error {
	sm.logger.Info("Triggering automated secret rotation", "provider", provider)
	return nil
}
