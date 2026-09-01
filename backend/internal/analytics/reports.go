package analytics

import (
	"fmt"
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type ReportEngine struct {
	repo Repository
}

func NewReportEngine(repo Repository) *ReportEngine {
	return &ReportEngine{repo: repo}
}

func (r *ReportEngine) ListGenerated(orgID string) ([]models.GeneratedReport, error) {
	return r.repo.GetGeneratedReports(orgID)
}

func (r *ReportEngine) GenerateReport(orgID, name, reportType, format, username string) (*models.GeneratedReport, error) {
	filePath := fmt.Sprintf("/storage/reports/report_%s_%d.%s", orgID, time.Now().Unix(), format)
	report := &models.GeneratedReport{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		Type:           reportType,
		Format:         format,
		FilePath:       filePath,
		FileSize:       1024 * 512, // 512KB mock
		GeneratedBy:    username,
		Status:         "completed",
		CreatedAt:      time.Now(),
	}

	if err := r.repo.CreateGeneratedReport(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (r *ReportEngine) ScheduleReport(orgID, name, reportType, format, frequency, recipients string) (*models.ScheduledReport, error) {
	sReport := &models.ScheduledReport{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           name,
		ReportType:     reportType,
		Format:         format,
		Frequency:      frequency,
		Recipients:     recipients,
		Enabled:        true,
		NextRunAt:      time.Now().Add(24 * time.Hour),
		CreatedAt:      time.Now(),
	}

	if err := r.repo.CreateScheduledReport(sReport); err != nil {
		return nil, err
	}
	return sReport, nil
}
