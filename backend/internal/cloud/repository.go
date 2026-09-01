package cloud

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetAccounts(orgID string) ([]models.CloudAccount, error)
	CreateAccount(acc *models.CloudAccount) error
	GetAWSResources(orgID string) ([]models.AWSResource, error)
	GetAzureResources(orgID string) ([]models.AzureResource, error)
	GetGCPResources(orgID string) ([]models.GCPResource, error)
	GetCosts(orgID string) ([]models.CloudCost, error)
	GetSecurityFindings(orgID string) ([]models.CloudSecurityFinding, error)
	GetAlerts(orgID string) ([]models.CloudAlert, error)
}

type postgresRepository struct{}

func NewRepository() Repository {
	return &postgresRepository{}
}

func (r *postgresRepository) GetAccounts(orgID string) ([]models.CloudAccount, error) {
	if database.DB == nil {
		return generateMockAccounts(orgID), nil
	}
	var accounts []models.CloudAccount
	db := database.DB.Model(&models.CloudAccount{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&accounts).Error
	if err != nil || len(accounts) == 0 {
		return generateMockAccounts(orgID), nil
	}
	return accounts, nil
}

func (r *postgresRepository) CreateAccount(acc *models.CloudAccount) error {
	if database.DB == nil {
		return nil
	}
	if acc.ID == uuid.Nil {
		acc.ID = uuid.New()
	}
	return database.DB.Create(acc).Error
}

func (r *postgresRepository) GetAWSResources(orgID string) ([]models.AWSResource, error) {
	if database.DB == nil {
		return []models.AWSResource{}, nil
	}
	var res []models.AWSResource
	db := database.DB.Model(&models.AWSResource{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Find(&res).Error
	return res, err
}

func (r *postgresRepository) GetAzureResources(orgID string) ([]models.AzureResource, error) {
	if database.DB == nil {
		return []models.AzureResource{}, nil
	}
	var res []models.AzureResource
	db := database.DB.Model(&models.AzureResource{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Find(&res).Error
	return res, err
}

func (r *postgresRepository) GetGCPResources(orgID string) ([]models.GCPResource, error) {
	if database.DB == nil {
		return []models.GCPResource{}, nil
	}
	var res []models.GCPResource
	db := database.DB.Model(&models.GCPResource{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Find(&res).Error
	return res, err
}

func (r *postgresRepository) GetCosts(orgID string) ([]models.CloudCost, error) {
	if database.DB == nil {
		return generateMockCosts(orgID), nil
	}
	var costs []models.CloudCost
	db := database.DB.Model(&models.CloudCost{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Find(&costs).Error
	if err != nil || len(costs) == 0 {
		return generateMockCosts(orgID), nil
	}
	return costs, nil
}

func (r *postgresRepository) GetSecurityFindings(orgID string) ([]models.CloudSecurityFinding, error) {
	if database.DB == nil {
		return generateMockSecurityFindings(orgID), nil
	}
	var findings []models.CloudSecurityFinding
	db := database.DB.Model(&models.CloudSecurityFinding{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&findings).Error
	if err != nil || len(findings) == 0 {
		return generateMockSecurityFindings(orgID), nil
	}
	return findings, nil
}

func (r *postgresRepository) GetAlerts(orgID string) ([]models.CloudAlert, error) {
	if database.DB == nil {
		return generateMockCloudAlerts(orgID), nil
	}
	var alerts []models.CloudAlert
	db := database.DB.Model(&models.CloudAlert{})
	if orgID != "" && orgID != "default" {
		db = db.Where("organization_id = ?", orgID)
	}
	err := db.Order("created_at desc").Find(&alerts).Error
	if err != nil || len(alerts) == 0 {
		return generateMockCloudAlerts(orgID), nil
	}
	return alerts, nil
}

// Fallback Mock Data Generators
func generateMockAccounts(orgID string) []models.CloudAccount {
	now := time.Now()
	return []models.CloudAccount{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "aws",
			AccountID:      "123456789012",
			AccountName:    "AWS Production Account",
			Regions:        "us-east-1, us-west-2, eu-central-1",
			Status:         "connected",
			TotalVMs:       24,
			TotalClusters:  2,
			MonthlyCost:    842.0,
			LastSyncedAt:   now,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "aws",
			AccountID:      "987654321098",
			AccountName:    "AWS Staging Account",
			Regions:        "us-east-1",
			Status:         "connected",
			TotalVMs:       8,
			TotalClusters:  1,
			MonthlyCost:    290.0,
			LastSyncedAt:   now,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "azure",
			AccountID:      "sub-az-production-01",
			AccountName:    "Azure Enterprise Subscription",
			Regions:        "eastus, westeurope",
			Status:         "connected",
			TotalVMs:       16,
			TotalClusters:  2,
			MonthlyCost:    511.0,
			LastSyncedAt:   now,
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "gcp",
			AccountID:      "gcp-prod-infra-4092",
			AccountName:    "GCP Infrastructure Project",
			Regions:        "us-central1, europe-west1",
			Status:         "connected",
			TotalVMs:       12,
			TotalClusters:  1,
			MonthlyCost:    218.0,
			LastSyncedAt:   now,
			CreatedAt:      now,
		},
	}
}

func generateMockCosts(orgID string) []models.CloudCost {
	return []models.CloudCost{
		{OrganizationID: orgID, Provider: "aws", Category: "compute", MonthlySpend: 420.0, ForecastedSpend: 460.0, PeriodMonth: "2026-07"},
		{OrganizationID: orgID, Provider: "aws", Category: "database", MonthlySpend: 210.0, ForecastedSpend: 220.0, PeriodMonth: "2026-07"},
		{OrganizationID: orgID, Provider: "aws", Category: "storage", MonthlySpend: 112.0, ForecastedSpend: 125.0, PeriodMonth: "2026-07"},
		{OrganizationID: orgID, Provider: "aws", Category: "network", MonthlySpend: 100.0, ForecastedSpend: 110.0, PeriodMonth: "2026-07"},

		{OrganizationID: orgID, Provider: "azure", Category: "compute", MonthlySpend: 280.0, ForecastedSpend: 300.0, PeriodMonth: "2026-07"},
		{OrganizationID: orgID, Provider: "azure", Category: "database", MonthlySpend: 140.0, ForecastedSpend: 150.0, PeriodMonth: "2026-07"},
		{OrganizationID: orgID, Provider: "azure", Category: "storage", MonthlySpend: 91.0, ForecastedSpend: 95.0, PeriodMonth: "2026-07"},

		{OrganizationID: orgID, Provider: "gcp", Category: "compute", MonthlySpend: 120.0, ForecastedSpend: 130.0, PeriodMonth: "2026-07"},
		{OrganizationID: orgID, Provider: "gcp", Category: "database", MonthlySpend: 68.0, ForecastedSpend: 75.0, PeriodMonth: "2026-07"},
		{OrganizationID: orgID, Provider: "gcp", Category: "storage", MonthlySpend: 30.0, ForecastedSpend: 32.0, PeriodMonth: "2026-07"},
	}
}

func generateMockSecurityFindings(orgID string) []models.CloudSecurityFinding {
	now := time.Now()
	return []models.CloudSecurityFinding{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "aws",
			ResourceID:     "s3-acme-public-logs",
			FindingType:    "public_s3",
			Severity:       "CRITICAL",
			Title:          "Public S3 Bucket Detected",
			Description:    "Bucket 's3-acme-public-logs' has global read ACL permissions enabled.",
			Remediation:    "Enable S3 Block Public Access at bucket policy level.",
			Status:         "OPEN",
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "aws",
			ResourceID:     "sg-0a1b2c3d4e5f",
			FindingType:    "open_security_group",
			Severity:       "HIGH",
			Title:          "Open SSH Port (22) to Internet",
			Description:    "Security group ingress rule allows 0.0.0.0/0 on port 22.",
			Remediation:    "Restrict ingress to corporate VPN IP cidr.",
			Status:         "OPEN",
			CreatedAt:      now,
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "azure",
			ResourceID:     "az-disk-os-01",
			FindingType:    "unencrypted_disk",
			Severity:       "MEDIUM",
			Title:          "Unencrypted Managed Disk",
			Description:    "Azure VM OS disk does not have customer-managed encryption (CMEK) enabled.",
			Remediation:    "Enable Disk Encryption Set using Azure Key Vault.",
			Status:         "OPEN",
			CreatedAt:      now,
		},
	}
}

func generateMockCloudAlerts(orgID string) []models.CloudAlert {
	now := time.Now()
	return []models.CloudAlert{
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "aws",
			ResourceName:   "i-0a1b2c3d4e5f6g7h8 (prod-web-ec2-01)",
			AlertType:      "ec2_down",
			Severity:       "CRITICAL",
			Message:        "EC2 instance status check failed (System & Instance unreachable).",
			Resolved:       false,
			CreatedAt:      now.Add(-15 * time.Minute),
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "aws",
			ResourceName:   "db-aws-rds-postgres-primary",
			AlertType:      "rds_high_cpu",
			Severity:       "HIGH",
			Message:        "RDS CPU utilization exceeded 92% threshold for >10 consecutive minutes.",
			Resolved:       false,
			CreatedAt:      now.Add(-45 * time.Minute),
		},
		{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Provider:       "azure",
			ResourceName:   "az-aks-cluster",
			AlertType:      "k8s_node_fail",
			Severity:       "HIGH",
			Message:        "AKS Node Pool node 'aks-agentpool-1928' transitioned to NotReady.",
			Resolved:       true,
			CreatedAt:      now.Add(-120 * time.Minute),
		},
	}
}
