package aws

import (
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type AWSProvider struct{}

func NewAWSProvider() *AWSProvider {
	return &AWSProvider{}
}

func (a *AWSProvider) DiscoverResources(accountID uuid.UUID, orgID string) ([]models.AWSResource, error) {
	now := time.Now()
	return []models.AWSResource{
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "i-0a1b2c3d4e5f6g7h8",
			ResourceType:   "ec2",
			Name:           "prod-web-ec2-01",
			Region:         "us-east-1",
			State:          "running",
			DetailsJSON:    `{"instance_type":"t3.large","vcpu":2,"memory_gb":8}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "i-0b2c3d4e5f6g7h8i9",
			ResourceType:   "ec2",
			Name:           "prod-api-ec2-02",
			Region:         "us-east-1",
			State:          "running",
			DetailsJSON:    `{"instance_type":"c5.xlarge","vcpu":4,"memory_gb":8}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "arn:aws:eks:us-east-1:123456789012:cluster/prod-eks-cluster",
			ResourceType:   "eks",
			Name:           "prod-eks-cluster",
			Region:         "us-east-1",
			State:          "ACTIVE",
			DetailsJSON:    `{"nodes":12,"k8s_version":"1.28"}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "db-aws-rds-postgres-primary",
			ResourceType:   "rds",
			Name:           "prod-rds-postgres",
			Region:         "us-west-2",
			State:          "available",
			DetailsJSON:    `{"engine":"postgres","instance_class":"db.r6g.xlarge","allocated_storage_gb":500}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "arn:aws:lambda:us-east-1:123456789012:function:image-processor",
			ResourceType:   "lambda",
			Name:           "image-processor",
			Region:         "us-east-1",
			State:          "active",
			DetailsJSON:    `{"runtime":"go1.x","timeout_sec":30}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			AccountID:      accountID,
			OrganizationID: orgID,
			ResourceID:     "s3-prod-assets-bucket-acme",
			ResourceType:   "s3",
			Name:           "s3-prod-assets-bucket-acme",
			Region:         "us-east-1",
			State:          "active",
			DetailsJSON:    `{"public":false,"encrypted":true}`,
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}, nil
}
