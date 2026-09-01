package backup

import (
	"time"

	"infrapilot/backend/internal/logger"
)

// RetentionPolicy defines backup retention rules
type RetentionPolicy struct {
	HourlyRetention  time.Duration
	DailyRetention   time.Duration
	WeeklyRetention  time.Duration
	MonthlyRetention time.Duration
}

// NewRetentionPolicy creates a new retention policy with defaults
func NewRetentionPolicy() *RetentionPolicy {
	return &RetentionPolicy{
		HourlyRetention:  24 * time.Hour,
		DailyRetention:   30 * 24 * time.Hour,
		WeeklyRetention:  12 * 7 * 24 * time.Hour,
		MonthlyRetention: 12 * 30 * 24 * time.Hour,
	}
}

// CalculateExpiration calculates expiration time for a given retention type
func CalculateExpiration(retentionType string, createdAt time.Time) time.Time {
	policy := NewRetentionPolicy()

	switch retentionType {
	case "hourly":
		return createdAt.Add(policy.HourlyRetention)
	case "daily":
		return createdAt.Add(policy.DailyRetention)
	case "weekly":
		return createdAt.Add(policy.WeeklyRetention)
	case "monthly":
		return createdAt.Add(policy.MonthlyRetention)
	default:
		return createdAt.Add(policy.DailyRetention)
	}
}

// EnforceRetention removes expired backups
func (r *RetentionPolicy) EnforceRetention(backupManager *BackupManager) {
	log := logger.Get()
	log.Info("Enforcing backup retention policies")

	// This would integrate with the database to find and remove expired backups
	// For now, we'll log the retention policy details
	log.Info("Retention policy configured",
		"hourly", r.HourlyRetention.String(),
		"daily", r.DailyRetention.String(),
		"weekly", r.WeeklyRetention.String(),
		"monthly", r.MonthlyRetention.String())

}
