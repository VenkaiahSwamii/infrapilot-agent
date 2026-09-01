package services

import (
	"testing"
)

func TestDashboardService_GetDashboardStats(t *testing.T) {
	stats, err := GetDashboardStats()
	if err != nil {
		t.Fatalf("Unexpected error getting dashboard stats: %v", err)
	}

	if stats == nil {
		t.Fatalf("Expected non-nil DashboardStats")
	}
}
