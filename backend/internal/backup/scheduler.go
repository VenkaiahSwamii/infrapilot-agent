package backup

import (
	"fmt"
	"sync"
	"time"

	"infrapilot/backend/internal/logger"
)

// BackupScheduler manages scheduled backup operations
type BackupScheduler struct {
	backupManager *BackupManager
	logger        *logger.Logger
	tickers       map[string]*time.Ticker
	done          chan bool
	mu            sync.RWMutex
	running       bool
}

// NewBackupScheduler creates a new backup scheduler
func NewBackupScheduler(backupManager *BackupManager) *BackupScheduler {
	return &BackupScheduler{
		backupManager: backupManager,
		logger:        logger.Get(),
		tickers:       make(map[string]*time.Ticker),
		done:          make(chan bool),
	}
}

// Schedule defines a backup schedule
type Schedule struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        BackupType `json:"type"`
	Frequency   string     `json:"frequency"` // hourly, daily, weekly, monthly
	Retention   string     `json:"retention"` // retention period
	Enabled     bool       `json:"enabled"`
	CronExpr    string     `json:"cron_expr"`
	Compress    bool       `json:"compress"`
	Encrypt     bool       `json:"encrypt"`
	StorageType string     `json:"storage_type"`
	LastRun     *time.Time `json:"last_run,omitempty"`
	NextRun     *time.Time `json:"next_run,omitempty"`
}

// Start begins the scheduler
func (s *BackupScheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	s.logger.Info("Starting backup scheduler")

	// Start default schedules
	go s.runScheduleLoop()
}

// Stop halts the scheduler
func (s *BackupScheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.done)

	for id, ticker := range s.tickers {
		ticker.Stop()
		delete(s.tickers, id)
	}

	s.logger.Info("Backup scheduler stopped")
}

// runScheduleLoop manages periodic schedule checks
func (s *BackupScheduler) runScheduleLoop() {
	// Check schedules every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Run initial checks on startup
	s.checkSchedules()

	for {
		select {
		case <-ticker.C:
			s.checkSchedules()
		case <-s.done:
			return
		}
	}
}

// checkSchedules evaluates all configured schedules and triggers backups
func (s *BackupScheduler) checkSchedules() {
	schedules := s.getDefaultSchedules()
	now := time.Now()

	for _, schedule := range schedules {
		if !schedule.Enabled {
			continue
		}

		shouldRun := false

		if schedule.NextRun != nil && now.After(schedule.NextRun.Add(-1*time.Minute)) {
			shouldRun = true
		}

		if shouldRun {
			s.logger.Info("Triggering scheduled backup",
				"id", schedule.ID,
				"type", schedule.Type,
				"frequency", schedule.Frequency)

			go s.executeScheduledBackup(schedule)
		}
	}
}

// executeScheduledBackup runs a backup for the given schedule
func (s *BackupScheduler) executeScheduledBackup(schedule Schedule) {
	var backupFile string
	var err error

	switch schedule.Type {
	case BackupTypePostgreSQL:
		backupFile, err = s.backupManager.BackupDatabase()
	case BackupTypeQdrant:
		backupFile, err = s.backupManager.BackupQdrant()
	case BackupTypeKubernetes:
		backupFile, err = s.backupManager.BackupKubernetes()
	case BackupTypeConfig:
		backupFile, err = s.backupManager.BackupConfig()
	case BackupTypeGrafana:
		backupFile, err = s.backupManager.BackupGrafana()
	case BackupTypeFull:
		backupFile, err = s.backupManager.BackupAll()
	default:
		err = fmt.Errorf("unknown backup type: %s", schedule.Type)
	}

	if err != nil {
		s.logger.Error("Scheduled backup failed",
			"id", schedule.ID,
			"type", schedule.Type,
			"error", err)
		return
	}

	s.logger.Info("Scheduled backup completed",
		"id", schedule.ID,
		"type", schedule.Type,
		"file", backupFile)
}

// getDefaultSchedules returns the default backup schedules
func (s *BackupScheduler) getDefaultSchedules() []Schedule {
	now := time.Now()

	// Daily at 02:00
	dailyNext := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
	dailyLast := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())

	// Weekly on Sunday at 02:00
	weeklyNext := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if now.Weekday() != time.Sunday {
		daysUntilSunday := (7 - int(now.Weekday())) % 7
		if daysUntilSunday == 0 {
			daysUntilSunday = 7
		}
		weeklyNext = time.Date(now.Year(), now.Month(), now.Day()+daysUntilSunday, 2, 0, 0, 0, now.Location())
	}

	// Monthly on 1st at 02:00
	monthlyNext := time.Date(now.Year(), now.Month()+1, 1, 2, 0, 0, 0, now.Location())

	// Hourly at minute 0
	hourlyNext := now.Truncate(time.Hour).Add(1 * time.Hour)

	return []Schedule{
		{
			ID:        "hourly-full",
			Name:      "Hourly Database Backup",
			Type:      BackupTypeFull,
			Frequency: "hourly",
			Retention: "24h",
			Enabled:   true,
			LastRun:   &dailyLast,
			NextRun:   &hourlyNext,
		},
		{
			ID:        "daily-full",
			Name:      "Daily Full Backup",
			Type:      BackupTypeFull,
			Frequency: "daily",
			Retention: "30d",
			Enabled:   true,
			LastRun:   &dailyLast,
			NextRun:   &dailyNext,
		},
		{
			ID:        "weekly-full",
			Name:      "Weekly Full Backup",
			Type:      BackupTypeFull,
			Frequency: "weekly",
			Retention: "12w",
			Enabled:   true,
			LastRun:   &dailyLast,
			NextRun:   &weeklyNext,
		},
		{
			ID:        "monthly-full",
			Name:      "Monthly Full Backup",
			Type:      BackupTypeFull,
			Frequency: "monthly",
			Retention: "12m",
			Enabled:   true,
			LastRun:   &dailyLast,
			NextRun:   &monthlyNext,
		},
	}
}

// AddSchedule adds a new backup schedule
func (s *BackupScheduler) AddSchedule(schedule Schedule) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if schedule.NextRun == nil {
		nextRun := s.calculateNextRun(schedule)
		schedule.NextRun = &nextRun
	}

	s.logger.Info("Backup schedule added",
		"id", schedule.ID,
		"type", schedule.Type,
		"frequency", schedule.Frequency)
}

// RemoveSchedule removes a backup schedule
func (s *BackupScheduler) RemoveSchedule(scheduleID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ticker, exists := s.tickers[scheduleID]; exists {
		ticker.Stop()
		delete(s.tickers, scheduleID)
	}

	s.logger.Info("Backup schedule removed", "id", scheduleID)
}

// calculateNextRun calculates the next run time for a schedule
func (s *BackupScheduler) calculateNextRun(schedule Schedule) time.Time {
	now := time.Now()

	switch schedule.Frequency {
	case "hourly":
		return now.Truncate(time.Hour).Add(1 * time.Hour)
	case "daily":
		return time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
	case "weekly":
		daysUntilSunday := (7 - int(now.Weekday())) % 7
		if daysUntilSunday == 0 {
			daysUntilSunday = 7
		}
		return time.Date(now.Year(), now.Month(), now.Day()+daysUntilSunday, 2, 0, 0, 0, now.Location())
	case "monthly":
		return time.Date(now.Year(), now.Month()+1, 1, 2, 0, 0, 0, now.Location())
	default:
		return now.Add(24 * time.Hour)
	}
}

// GetSchedules returns all configured schedules
func (s *BackupScheduler) GetSchedules() []Schedule {
	return s.getDefaultSchedules()
}
