package services

import (
	"testing"
)

func TestOverviewService_GetDashboardStats(t *testing.T) {
	stats, err := GetDashboardStats()
	if err != nil {
		t.Fatalf("Expected nil error for GetDashboardStats, got %v", err)
	}

	if stats.TotalMachines <= 0 {
		t.Errorf("Expected totalMachines > 0, got %d", stats.TotalMachines)
	}

	if stats.MachineHealthScore < 0 || stats.MachineHealthScore > 100 {
		t.Errorf("Expected machineHealthScore between 0 and 100, got %.2f", stats.MachineHealthScore)
	}
}
