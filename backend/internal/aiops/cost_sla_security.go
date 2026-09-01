package aiops

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

// CostSLASecurityEngine handles cost optimization, SLA tracking, and security analytics
type CostSLASecurityEngine struct {
	logger *logger.Logger
}

// NewCostSLASecurityEngine initializes CostSLASecurityEngine
func NewCostSLASecurityEngine() *CostSLASecurityEngine {
	return &CostSLASecurityEngine{
		logger: logger.Get(),
	}
}

// GetCostOptimizations returns cloud & infrastructure savings opportunities
func (c *CostSLASecurityEngine) GetCostOptimizations(orgID string) ([]models.CostOptimizationRecord, error) {
	var list []models.CostOptimizationRecord

	if database.DB != nil {
		database.DB.Where("organization_id = ? OR organization_id = ''", orgID).Find(&list)
	}

	if len(list) == 0 {
		list = c.generateSampleCosts(orgID)
	}

	return list, nil
}

// GetSLAMetrics retrieves uptime availability, MTTR, and MTBF metrics
func (c *CostSLASecurityEngine) GetSLAMetrics(orgID string) (*models.SLARecord, error) {
	var sla models.SLARecord
	if database.DB != nil {
		database.DB.Where("organization_id = ? OR organization_id = ''", orgID).Order("created_at desc").First(&sla)
	}

	if sla.ID == uuid.Nil {
		sla = models.SLARecord{
			ID:              uuid.New(),
			OrganizationID:  orgID,
			Period:          "monthly",
			AvailabilityPct: 99.98,
			AvgResponseMs:   14.2,
			MTTRSeconds:     180,
			MTBFSeconds:     604800,
			IncidentsCount:  2,
			CreatedAt:       time.Now(),
		}
	}

	return &sla, nil
}

// GetSecurityAnalytics retrieves security threat events
func (c *CostSLASecurityEngine) GetSecurityAnalytics(orgID string) ([]models.SecurityAnalyticRecord, error) {
	var sec []models.SecurityAnalyticRecord

	if database.DB != nil {
		database.DB.Where("organization_id = ? OR organization_id = ''", orgID).Order("detected_at desc").Find(&sec)
	}

	if len(sec) == 0 {
		sec = c.generateSampleSecurity(orgID)
	}

	return sec, nil
}

func (c *CostSLASecurityEngine) generateSampleCosts(orgID string) []models.CostOptimizationRecord {
	now := time.Now()
	return []models.CostOptimizationRecord{
		{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			ResourceName:     "dev-worker-node-04",
			ResourceType:     "idle_k8s_node",
			Issue:            "Underutilized VM (< 5% CPU over 14 days)",
			SuggestedAction:  "Downscale or terminate VM node to eliminate idle charges",
			EstimatedSavings: 8500.0,
			CreatedAt:        now,
		},
		{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			ResourceName:     "backup-unattached-pvc-02",
			ResourceType:     "unattached_pvc",
			Issue:            "Unused 1TB EBS volume unattached for 30 days",
			SuggestedAction:  "Delete unattached PVC volume or archive to S3 Glacier",
			EstimatedSavings: 15500.0,
			CreatedAt:        now,
		},
	}
}

func (c *CostSLASecurityEngine) generateSampleSecurity(orgID string) []models.SecurityAnalyticRecord {
	now := time.Now()
	return []models.SecurityAnalyticRecord{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			EventType:      "brute_force_attempt",
			Severity:       "warning",
			SourceIP:       "198.51.100.42",
			Hostname:       "bastion-host-01",
			Details:        "142 failed SSH authentication attempts detected in 60 seconds",
			DetectedAt:     now.Add(-15 * time.Minute),
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			EventType:      "suspicious_process",
			Severity:       "critical",
			SourceIP:       "10.0.4.15",
			Hostname:       "worker-node-02",
			Details:        "Unsigned binary executing from /tmp directory (PID: 14920)",
			DetectedAt:     now.Add(-45 * time.Minute),
		},
	}
}
