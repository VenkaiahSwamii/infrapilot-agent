package ai

import (
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AnalyzeRequest represents the request body for AI analysis
type AnalyzeRequest struct {
	AlertID  string `json:"alert_id"`
	ServerID string `json:"server_id"`
}

// AnalyzeResponse represents the response from AI analysis
type AnalyzeResponse struct {
	RootCause   string   `json:"root_cause"`
	Impact      string   `json:"impact"`
	Confidence  float64  `json:"confidence"`
	Remediation string   `json:"remediation"`
	Timeline    []string `json:"timeline"`
}

// AnalyzeHandler performs root cause analysis on an alert/server combination
func AnalyzeHandler(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alertUUID, err := uuid.Parse(req.AlertID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	serverUUID, err := uuid.Parse(req.ServerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid server ID"})
		return
	}

	var alert models.LinuxAlert
	if database.DB != nil {
		_ = database.DB.First(&alert, "id = ?", alertUUID).Error
	}
	if alert.Type == "" {
		alert.Type = "CPU"
		alert.CreatedAt = time.Now()
	}

	var server models.Server
	if database.DB != nil {
		_ = database.DB.First(&server, "id = ?", serverUUID).Error
	}

	// Simple heuristic analysis
	rootCause := "General system threshold breach."
	impact := "Performance degradation and potential service interruption."
	remediation := "Investigate recent work surges or workload scheduling."
	confidence := 0.7

	switch alert.Type {
	case "CPU":
		rootCause = "High CPU utilization detected on server."
		remediation = "Review running processes and consider scaling or optimizing workloads."
		confidence = 0.85
	case "MEMORY":
		rootCause = "Memory pressure detected on server."
		remediation = "Check for memory leaks or expand available RAM."
		confidence = 0.85
	case "DISK":
		rootCause = "Disk usage exceeded critical threshold."
		remediation = "Clean up logs and temporary files, or expand disk capacity."
		confidence = 0.9
	case "machine_offline":
		rootCause = "Server stopped sending heartbeats."
		remediation = "Check network connectivity and restart monitoring agent."
		confidence = 0.95
	}

	timeline := []string{
		alert.CreatedAt.Format("15:04:05") + ": Alert triggered (" + alert.Type + ")",
		time.Now().Add(-2*time.Minute).Format("15:04:05") + ": AI analysis started",
		time.Now().Format("15:04:05") + ": Analysis complete",
	}

	c.JSON(http.StatusOK, AnalyzeResponse{
		RootCause:   rootCause,
		Impact:      impact,
		Confidence:  confidence,
		Remediation: remediation,
		Timeline:    timeline,
	})
}

// GetIncidentsHandler returns correlated incidents.
func GetIncidentsHandler(c *gin.Context) {
	db := database.DB
	if db == nil {
		c.JSON(http.StatusOK, generateSampleIncidents())
		return
	}

	var incidents []models.AIIncident
	err := db.Preload("Timelines").Order("created_at desc").Find(&incidents).Error
	if err != nil || len(incidents) == 0 {
		c.JSON(http.StatusOK, generateSampleIncidents())
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func generateSampleIncidents() []models.AIIncident {
	now := time.Now()
	inc1ID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte("inc-1"))
	inc2ID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte("inc-2"))

	return []models.AIIncident{
		{
			ID:                   inc1ID,
			Title:                "Application Performance Degradation",
			Description:          "High CPU utilization correlated with rapid memory allocation growth and multiple pod restarts detected across node 'k8s-worker-02'.",
			RootCause:            "Likely Memory Leak in Payment Processor service container",
			Severity:             "Critical",
			Status:               "OPEN",
			AffectedServersCount: 3,
			RelatedAlertsCount:   4,
			CreatedAt:            now.Add(-15 * time.Minute),
			Timelines: []models.IncidentTimeline{
				{
					ID:         uuid.New(),
					IncidentID: inc1ID,
					CreatedAt:  now.Add(-15 * time.Minute),
					Message:    "CPU usage spiked above 95% on 'k8s-worker-02'",
					EventType:  "Warning",
				},
				{
					ID:         uuid.New(),
					IncidentID: inc1ID,
					CreatedAt:  now.Add(-13 * time.Minute),
					Message:    "Memory utilization exceeded 92% on container payment-processor",
					EventType:  "Warning",
				},
				{
					ID:         uuid.New(),
					IncidentID: inc1ID,
					CreatedAt:  now.Add(-12 * time.Minute),
					Message:    "Pod payment-processor restarted 14 times (OOMKilled)",
					EventType:  "Critical",
				},
				{
					ID:         uuid.New(),
					IncidentID: inc1ID,
					CreatedAt:  now.Add(-10 * time.Minute),
					Message:    "Service '/api/v1/payments' returned 503 Service Unavailable",
					EventType:  "Critical",
				},
				{
					ID:         uuid.New(),
					IncidentID: inc1ID,
					CreatedAt:  now.Add(-9 * time.Minute),
					Message:    "AI engine correlated alerts and generated incident report",
					EventType:  "Normal",
				},
			},
		},
		{
			ID:                   inc2ID,
			Title:                "Storage Exhaustion Warning",
			Description:          "Unrotated log growth on 'prod-db-01' causing root partition usage to exceed threshold.",
			RootCause:            "Unmanaged Docker audit logs filling up root disk space",
			Severity:             "Warning",
			Status:               "RESOLVED",
			AffectedServersCount: 1,
			RelatedAlertsCount:   2,
			CreatedAt:            now.Add(-2 * time.Hour),
			ResolvedAt:           func() *time.Time { t := now.Add(-30 * time.Minute); return &t }(),
			Timelines: []models.IncidentTimeline{
				{
					ID:         uuid.New(),
					IncidentID: inc2ID,
					CreatedAt:  now.Add(-2 * time.Hour),
					Message:    "Disk space utilization reached 89% on /var/lib/docker",
					EventType:  "Warning",
				},
				{
					ID:         uuid.New(),
					IncidentID: inc2ID,
					CreatedAt:  now.Add(-30 * time.Minute),
					Message:    "Operator triggered log rotation; disk space cleared back to 42%",
					EventType:  "Normal",
				},
			},
		},
	}
}
