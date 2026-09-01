package azure

import (
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type AzureProvider struct{}

func NewAzureProvider() *AzureProvider {
	return &AzureProvider{}
}

func (a *AzureProvider) DiscoverResources(accountID uuid.UUID, orgID string) ([]models.AzureResource, error) {
	now := time.Now()
	return []models.AzureResource{
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "/subscriptions/sub-az-01/resourceGroups/rg-prod/providers/Microsoft.Compute/virtualMachines/az-vm-web-01",
			ResourceType:   "vm",
			Name:           "az-vm-web-01",
			ResourceGroup:  "rg-prod",
			Region:         "eastus",
			State:          "running",
			DetailsJSON:    `{"vm_size":"Standard_D4s_v5","os":"Ubuntu 22.04"}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "/subscriptions/sub-az-01/resourceGroups/rg-prod/providers/Microsoft.ContainerService/managedClusters/az-aks-cluster",
			ResourceType:   "aks",
			Name:           "az-aks-cluster",
			ResourceGroup:  "rg-prod",
			Region:         "eastus",
			State:          "Succeeded",
			DetailsJSON:    `{"node_pools":3,"k8s_version":"1.27.7"}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "/subscriptions/sub-az-01/resourceGroups/rg-prod/providers/Microsoft.Sql/servers/sql-primary/databases/db-analytics",
			ResourceType:   "sql",
			Name:           "az-sql-analytics",
			ResourceGroup:  "rg-prod",
			Region:         "westeurope",
			State:          "Online",
			DetailsJSON:    `{"edition":"GeneralPurpose","max_size_gb":250}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}, nil
}
