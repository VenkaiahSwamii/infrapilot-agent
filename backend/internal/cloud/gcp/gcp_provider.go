package gcp

import (
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type GCPProvider struct{}

func NewGCPProvider() *GCPProvider {
	return &GCPProvider{}
}

func (g *GCPProvider) DiscoverResources(accountID uuid.UUID, orgID string) ([]models.GCPResource, error) {
	now := time.Now()
	return []models.GCPResource{
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "projects/gcp-prod-infra/zones/us-central1-a/instances/gcp-instance-worker-01",
			ResourceType:   "compute",
			Name:           "gcp-instance-worker-01",
			ProjectID:      "gcp-prod-infra",
			Zone:           "us-central1-a",
			State:          "RUNNING",
			DetailsJSON:    `{"machine_type":"n2-standard-4","cpus":4,"ram_gb":16}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "projects/gcp-prod-infra/locations/us-central1/clusters/gcp-gke-prod",
			ResourceType:   "gke",
			Name:           "gcp-gke-prod",
			ProjectID:      "gcp-prod-infra",
			Zone:           "us-central1",
			State:          "RUNNING",
			DetailsJSON:    `{"num_nodes":8,"node_config":"n2-standard-4"}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "projects/gcp-prod-infra/instances/cloudsql-replica-01",
			ResourceType:   "sql",
			Name:           "cloudsql-replica-01",
			ProjectID:      "gcp-prod-infra",
			Zone:           "us-central1-b",
			State:          "RUNNABLE",
			DetailsJSON:    `{"database_version":"POSTGRES_15","tier":"db-custom-4-16384"}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}, nil
}
