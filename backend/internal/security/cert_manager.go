package security

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type CertManager struct {
	logger *logger.Logger
}

func NewCertManager() *CertManager {
	return &CertManager{
		logger: logger.Get(),
	}
}

// ListCertificates returns TLS and Internal CA certificate lifecycles
func (cm *CertManager) ListCertificates(orgID string) ([]models.CertificateRecord, error) {
	var certs []models.CertificateRecord

	if database.DB != nil {
		database.DB.Where("organization_id = ? OR organization_id = ''", orgID).Find(&certs)
	}

	if len(certs) == 0 {
		now := time.Now()
		certs = []models.CertificateRecord{
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
				CommonName:     "*.infrapilot.company.com",
				Issuer:         "Let's Encrypt Authority X3",
				SerialNumber:   "04:A1:8B:C9:2D",
				ValidFrom:      now.AddDate(0, -2, 0),
				ValidTo:        now.AddDate(0, 1, 0),
				AutoRenew:      true,
				Status:         "active",
				CreatedAt:      now,
			},
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
				CommonName:     "internal-ca.infrapilot.local",
				Issuer:         "InfraPilot Root CA",
				SerialNumber:   "01:00:99:FF:E2",
				ValidFrom:      now.AddDate(-1, 0, 0),
				ValidTo:        now.AddDate(4, 0, 0),
				AutoRenew:      true,
				Status:         "active",
				CreatedAt:      now,
			},
		}
	}

	return certs, nil
}
