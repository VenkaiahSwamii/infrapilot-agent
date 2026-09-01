package cloud

import (
	"time"

	"infrapilot/backend/internal/cloud/aws"
	"infrapilot/backend/internal/cloud/azure"
	"infrapilot/backend/internal/cloud/gcp"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Service struct {
	repo            Repository
	awsProvider     *aws.AWSProvider
	azureProvider   *azure.AzureProvider
	gcpProvider     *gcp.GCPProvider
	inventoryEngine *InventoryEngine
	securityEngine  *SecurityEngine
	costEngine      *CostEngine
	alertEngine     *AlertEngine
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:            repo,
		awsProvider:     aws.NewAWSProvider(),
		azureProvider:   azure.NewAzureProvider(),
		gcpProvider:     gcp.NewGCPProvider(),
		inventoryEngine: NewInventoryEngine(repo),
		securityEngine:  NewSecurityEngine(repo),
		costEngine:      NewCostEngine(repo),
		alertEngine:     NewAlertEngine(repo),
	}
}

func (s *Service) ConnectAWS(orgID, accountID, accountName, accessKey, secretKey, roleARN, regions string) (*models.CloudAccount, error) {
	acc := &models.CloudAccount{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Provider:       "aws",
		AccountID:      accountID,
		AccountName:    accountName,
		Regions:        regions,
		Status:         "connected",
		TotalVMs:       24,
		TotalClusters:  2,
		MonthlyCost:    842.0,
		LastSyncedAt:   time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateAccount(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *Service) ConnectAzure(orgID, subscriptionID, name, tenantID, clientID, clientSecret string) (*models.CloudAccount, error) {
	acc := &models.CloudAccount{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Provider:       "azure",
		AccountID:      subscriptionID,
		AccountName:    name,
		Regions:        "eastus, westeurope",
		Status:         "connected",
		TotalVMs:       16,
		TotalClusters:  2,
		MonthlyCost:    511.0,
		LastSyncedAt:   time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateAccount(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *Service) ConnectGCP(orgID, projectID, name, serviceAccountJSON string) (*models.CloudAccount, error) {
	acc := &models.CloudAccount{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Provider:       "gcp",
		AccountID:      projectID,
		AccountName:    name,
		Regions:        "us-central1, europe-west1",
		Status:         "connected",
		TotalVMs:       12,
		TotalClusters:  1,
		MonthlyCost:    218.0,
		LastSyncedAt:   time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateAccount(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *Service) GetAccounts(orgID string) ([]models.CloudAccount, error) {
	return s.repo.GetAccounts(orgID)
}

func (s *Service) GetHybridInventory(orgID string) (HybridResourceSummary, error) {
	return s.inventoryEngine.GetHybridSummary(orgID)
}

func (s *Service) GetCosts(orgID string) (CostSummary, error) {
	return s.costEngine.GetCostSummary(orgID)
}

func (s *Service) GetSecurity(orgID string) ([]models.CloudSecurityFinding, error) {
	return s.securityEngine.GetFindings(orgID)
}

func (s *Service) GetAlerts(orgID string) ([]models.CloudAlert, error) {
	return s.alertEngine.GetAlerts(orgID)
}

func (s *Service) GetAWSResources(orgID string) ([]models.AWSResource, error) {
	accounts, _ := s.repo.GetAccounts(orgID)
	accID := uuid.Nil
	if len(accounts) > 0 {
		accID = accounts[0].ID
	}
	return s.awsProvider.DiscoverResources(accID, orgID)
}

func (s *Service) GetAzureResources(orgID string) ([]models.AzureResource, error) {
	accounts, _ := s.repo.GetAccounts(orgID)
	accID := uuid.Nil
	if len(accounts) > 0 {
		accID = accounts[0].ID
	}
	return s.azureProvider.DiscoverResources(accID, orgID)
}

func (s *Service) GetGCPResources(orgID string) ([]models.GCPResource, error) {
	accounts, _ := s.repo.GetAccounts(orgID)
	accID := uuid.Nil
	if len(accounts) > 0 {
		accID = accounts[0].ID
	}
	return s.gcpProvider.DiscoverResources(accID, orgID)
}
