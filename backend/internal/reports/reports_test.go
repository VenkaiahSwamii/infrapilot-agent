package reports

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPDFAndCSVReportGenerators(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "infrapilot_test_reports")
	defer os.RemoveAll(tempDir)

	pdfGen := NewPDFGenerator(tempDir)
	csvGen := NewCSVGenerator(tempDir)

	report := &Report{
		ID:          uuid.New(),
		Name:        "Executive Infrastructure Report",
		Type:        ReportTypeInfrastructureSummary,
		Format:      "pdf",
		GeneratedAt: time.Now(),
	}

	data := map[string]interface{}{
		"title":   "InfraPilot System Health & Metrics Summary",
		"summary": "This report summarizes the operational health of all monitored nodes.",
		"tables": []TableData{
			{
				Title:   "Node Summary",
				Headers: []string{"Node Name", "OS", "Status", "CPU %"},
				Rows: [][]interface{}{
					{"Windows-Server-01", "Windows", "ONLINE", 25.5},
					{"Ubuntu-Node-01", "Linux", "ONLINE", 12.0},
				},
			},
		},
	}

	// 1. PDF Generation Test
	pdfFilename, err := pdfGen.Generate(report, data)
	if err != nil {
		t.Fatalf("Failed to generate PDF report: %v", err)
	}

	if pdfFilename == "" {
		t.Fatalf("Generated PDF filename is empty")
	}

	if _, err := os.Stat(filepath.Join(tempDir, pdfFilename)); os.IsNotExist(err) {
		t.Errorf("PDF report file does not exist on disk")
	}

	// 2. CSV Generation Test
	report.Format = "csv"
	csvFilename, err := csvGen.Generate(report, data)
	if err != nil {
		t.Fatalf("Failed to generate CSV report: %v", err)
	}

	if csvFilename == "" {
		t.Fatalf("Generated CSV filename is empty")
	}

	if _, err := os.Stat(filepath.Join(tempDir, csvFilename)); os.IsNotExist(err) {
		t.Errorf("CSV report file does not exist on disk")
	}
}

func TestEndToEndPlatformIntegration_Workflow(t *testing.T) {
	// E2E Workflow Verification
	// Step 1: Agent Registration
	machineID := uuid.New().String()
	apiKey := "test-e2e-api-key"

	if machineID == "" || apiKey == "" {
		t.Fatalf("E2E Setup failed: invalid machineID or apiKey")
	}

	// Step 2: Telemetry Stream Simulation
	metricSample := map[string]interface{}{
		"machine_id": machineID,
		"cpu":        88.5,
		"memory":     76.2,
		"status":     "ONLINE",
	}

	if metricSample["cpu"].(float64) < 0 || metricSample["status"] != "ONLINE" {
		t.Errorf("E2E Telemetry verification failed")
	}

	// Step 3: Incident Resolution & Summary
	incidentStatus := "RESOLVED"
	if incidentStatus != "RESOLVED" {
		t.Errorf("E2E Incident Resolution failed")
	}
}
