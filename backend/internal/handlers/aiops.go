package handlers

import (
	"fmt"
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
)

type AIOpsInsight struct {
	Type           string    `json:"type"` // anomaly, cost, prediction, capacity, self_healing
	MachineID      string    `json:"machine_id"`
	Hostname       string    `json:"hostname"`
	Severity       string    `json:"severity"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Recommendation string    `json:"recommendation"`
	Timestamp      time.Time `json:"timestamp"`
}

// GetAIOpsDashboard analyzes system baselines, disk growth trends, cost opportunities, and returns AIOps insights
func GetAIOpsDashboard(c *gin.Context) {
	var machines []models.Machine
	if err := database.DB.Find(&machines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query machines"})
		return
	}

	var insights []AIOpsInsight = []AIOpsInsight{}

	for _, m := range machines {
		// 1. Fetch historical CPU metrics
		var cpuHistory []models.HistoricalMetric
		database.DB.Where("machine_id = ? AND metric = ?", m.ID, "cpu").Order("sampled_at asc").Find(&cpuHistory)

		// 2. Fetch historical Memory metrics
		var memHistory []models.HistoricalMetric
		database.DB.Where("machine_id = ? AND metric = ?", m.ID, "memory").Order("sampled_at asc").Find(&memHistory)

		// 3. Fetch historical Disk metrics
		var diskHistory []models.HistoricalMetric
		database.DB.Where("machine_id = ? AND metric = ?", m.ID, "disk").Order("sampled_at asc").Find(&diskHistory)

		var avgCPU float64
		var avgMem float64
		var currentCPU float64
		var currentMem float64
		var currentDisk float64

		var latest models.Metric
		errLatest := database.DB.Where("machine_id = ?", m.ID).Order("created_at desc").First(&latest).Error
		if errLatest == nil {
			currentCPU = latest.CPUUsage
			currentMem = latest.MemoryUsage
			currentDisk = latest.DiskUsage
		}

		// Calculate CPU average
		nCPU := float64(len(cpuHistory))
		if nCPU > 0 {
			var totalCPU float64
			for _, h := range cpuHistory {
				totalCPU += h.Value
			}
			avgCPU = totalCPU / nCPU
		} else {
			avgCPU = currentCPU
		}

		// Calculate Memory average
		nMem := float64(len(memHistory))
		if nMem > 0 {
			var totalMem float64
			for _, h := range memHistory {
				totalMem += h.Value
			}
			avgMem = totalMem / nMem
		} else {
			avgMem = currentMem
		}

		// 1. Anomaly Detection: Baseline deviation spike
		if currentCPU > 50.0 && (currentCPU-avgCPU) > 20.0 {
			insights = append(insights, AIOpsInsight{
				Type:           "anomaly",
				MachineID:      m.ID.String(),
				Hostname:       m.Hostname,
				Severity:       "Warning",
				Title:          "CPU Anomaly Detected",
				Description:    fmt.Sprintf("Current CPU (%.1f%%) is significantly higher than historical baseline average (%.1f%%).", currentCPU, avgCPU),
				Recommendation: "Analyze active processes for execution spikes or thread lock leaks.",
				Timestamp:      time.Now(),
			})
		}

		// 2. Cost Optimization: Idle virtual machines detection
		if nCPU >= 5 && avgCPU < 8.0 && avgMem < 25.0 {
			insights = append(insights, AIOpsInsight{
				Type:           "cost",
				MachineID:      m.ID.String(),
				Hostname:       m.Hostname,
				Severity:       "Info",
				Title:          "Underutilized Server Capacity",
				Description:    fmt.Sprintf("Machine has average CPU of %.1f%% and Memory of %.1f%% over past historical cycles.", avgCPU, avgMem),
				Recommendation: "Reduce VM size allocation (vCPU/RAM allocation) to optimize monthly cost by ~40%.",
				Timestamp:      time.Now(),
			})
		}

		// 3. Predictive Maintenance: Disk full forecast
		if len(diskHistory) >= 2 {
			firstDisk := diskHistory[0].Value
			lastDisk := diskHistory[len(diskHistory)-1].Value
			diskChange := lastDisk - firstDisk
			if diskChange > 0.1 {
				daysSpan := diskHistory[len(diskHistory)-1].SampledAt.Sub(diskHistory[0].SampledAt).Hours() / 24.0
				if daysSpan > 0 {
					growthRatePerDay := diskChange / daysSpan
					remainingDisk := 100.0 - currentDisk
					daysToFull := remainingDisk / growthRatePerDay
					if daysToFull < 14 {
						severity := "Warning"
						if daysToFull < 5 {
							severity = "Critical"
						}
						insights = append(insights, AIOpsInsight{
							Type:           "prediction",
							MachineID:      m.ID.String(),
							Hostname:       m.Hostname,
							Severity:       severity,
							Title:          "Disk Depletion Forecast",
							Description:    fmt.Sprintf("Disk usage is at %.1f%% and growing at %.2f%% per day. Volume will reach capacity in %.1f days.", currentDisk, growthRatePerDay, daysToFull),
							Recommendation: "Initiate temp folder cleanup, compress database indices, or allocate larger storage blocks.",
							Timestamp:      time.Now(),
						})
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"insights_count": len(insights),
		"insights":       insights,
	})
}
