package cloud

import (
	"infrapilot/backend/internal/models"
)

type CostSummary struct {
	TotalMonthly float64            `json:"total_monthly"`
	AWSMonthly   float64            `json:"aws_monthly"`
	AzureMonthly float64            `json:"azure_monthly"`
	GCPMonthly   float64            `json:"gcp_monthly"`
	Breakdown    []models.CloudCost `json:"breakdown"`
}

type CostEngine struct {
	repo Repository
}

func NewCostEngine(repo Repository) *CostEngine {
	return &CostEngine{repo: repo}
}

func (c *CostEngine) GetCostSummary(orgID string) (CostSummary, error) {
	costs, err := c.repo.GetCosts(orgID)
	if err != nil {
		return CostSummary{}, err
	}

	awsSum, azSum, gcpSum := 0.0, 0.0, 0.0
	for _, item := range costs {
		switch item.Provider {
		case "aws":
			awsSum += item.MonthlySpend
		case "azure":
			azSum += item.MonthlySpend
		case "gcp":
			gcpSum += item.MonthlySpend
		}
	}

	return CostSummary{
		TotalMonthly: awsSum + azSum + gcpSum,
		AWSMonthly:   awsSum,
		AzureMonthly: azSum,
		GCPMonthly:   gcpSum,
		Breakdown:    costs,
	}, nil
}
