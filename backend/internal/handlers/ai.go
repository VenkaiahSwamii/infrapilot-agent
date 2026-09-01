package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatRequest struct {
	Question string `json:"question" binding:"required"`
}

type ChatResponse struct {
	Answer string `json:"answer"`
}

// AIChat handles incoming question queries using a hybrid NLP/LLM engine
func AIChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question := strings.TrimSpace(req.Question)
	lowerQuest := strings.ToLower(question)

	// Fetch system metrics context for the AI prompt or rules matcher
	var machines []models.Machine
	database.DB.Find(&machines)

	var metrics []models.Metric
	// Fetch latest metric sample for each machine
	for _, m := range machines {
		var latest models.Metric
		if database.DB.Where("machine_id = ?", m.ID).Order("created_at desc").First(&latest).Error == nil {
			metrics = append(metrics, latest)
		}
	}

	var alerts []models.LinuxAlert
	database.DB.Where("status = 'OPEN'").Find(&alerts)

	// Context formatted for LLM and rule-matching fallback
	var ctxBuf bytes.Buffer
	ctxBuf.WriteString("System Status:\n")
	for _, m := range machines {
		status := "Offline"
		if strings.ToUpper(m.Status) == "ONLINE" && time.Since(m.LastSeen) <= 90*time.Second {
			status = "Online"
		}

		cpu := 0.0
		mem := 0.0
		disk := 0.0
		for _, met := range metrics {
			if met.MachineID == m.ID {
				cpu = met.CPUUsage
				mem = met.MemoryUsage
				disk = met.DiskUsage
				break
			}
		}

		ctxBuf.WriteString(fmt.Sprintf("- Machine: %s, CPU: %.1f%%, RAM: %.1f%%, Disk: %.1f%%, Status: %s, LastSeen: %s\n",
			m.Hostname, cpu, mem, disk, status, m.LastSeen.Format(time.RFC3339)))
	}

	ctxBuf.WriteString("\nActive Alerts:\n")
	if len(alerts) == 0 {
		ctxBuf.WriteString("- No active alerts.\n")
	} else {
		for _, al := range alerts {
			ctxBuf.WriteString(fmt.Sprintf("- [%s] %s: %s (Severity: %s)\n", al.Type, al.Message, al.Status, al.Severity))
		}
	}

	systemContext := ctxBuf.String()

	// 1. Hybrid Rule Matcher (guarantees instantaneous correct answers for preset queries)
	var answer string
	matched := true

	switch {
	case strings.Contains(lowerQuest, "offline") || strings.Contains(lowerQuest, "down"):
		var offline []string
		for _, m := range machines {
			if strings.ToUpper(m.Status) != "ONLINE" || time.Since(m.LastSeen) > 90*time.Second {
				offline = append(offline, m.Hostname)
			}
		}
		if len(offline) == 0 {
			answer = "All machines are currently Online 🟢. There are no offline machines."
		} else {
			answer = fmt.Sprintf("The following machines are currently Offline 🔴:\n- %s\n\nI recommend verifying the agent process or network connectivity on these hosts.", strings.Join(offline, "\n- "))
		}

	case strings.Contains(lowerQuest, "highest cpu") || strings.Contains(lowerQuest, "max cpu"):
		var highestHost string
		highestCPU := -1.0
		for _, m := range machines {
			for _, met := range metrics {
				if met.MachineID == m.ID && met.CPUUsage > highestCPU {
					highestCPU = met.CPUUsage
					highestHost = m.Hostname
				}
			}
		}
		if highestHost == "" {
			answer = "No metrics gathered yet to determine CPU utilization."
		} else {
			answer = fmt.Sprintf("Machine **%s** currently has the highest CPU usage at **%.1f%%**.", highestHost, highestCPU)
			if highestCPU > 90 {
				answer += "\n\n⚠️ CPU usage exceeds the critical threshold of 90%. Consider reviewing running processes."
			}
		}

	case strings.Contains(lowerQuest, "highest memory") || strings.Contains(lowerQuest, "max memory") || strings.Contains(lowerQuest, "highest ram"):
		var highestHost string
		highestMem := -1.0
		for _, m := range machines {
			for _, met := range metrics {
				if met.MachineID == m.ID && met.MemoryUsage > highestMem {
					highestMem = met.MemoryUsage
					highestHost = m.Hostname
				}
			}
		}
		if highestHost == "" {
			answer = "No metrics gathered yet to determine memory usage."
		} else {
			answer = fmt.Sprintf("Machine **%s** currently has the highest memory utilization at **%.1f%%**.", highestHost, highestMem)
		}

	case strings.Contains(lowerQuest, "alert") || strings.Contains(lowerQuest, "incident"):
		if len(alerts) == 0 {
			answer = "There are no open or active alerts on the system 🟢."
		} else {
			var alertLines []string
			for _, al := range alerts {
				alertLines = append(alertLines, fmt.Sprintf("- **%s**: %s (%s)", al.Severity, al.Message, al.Type))
			}
			answer = fmt.Sprintf("Found %d active alert(s):\n%s", len(alerts), strings.Join(alertLines, "\n"))
		}

	case strings.Contains(lowerQuest, "need attention") || strings.Contains(lowerQuest, "trouble") || strings.Contains(lowerQuest, "slow"):
		var attentionLines []string
		for _, m := range machines {
			isOffline := strings.ToUpper(m.Status) != "ONLINE" || time.Since(m.LastSeen) > 90*time.Second
			cpu := 0.0
			for _, met := range metrics {
				if met.MachineID == m.ID {
					cpu = met.CPUUsage
					break
				}
			}
			if isOffline {
				attentionLines = append(attentionLines, fmt.Sprintf("- **%s**: Machine is OFFLINE", m.Hostname))
			} else if cpu > 90 {
				attentionLines = append(attentionLines, fmt.Sprintf("- **%s**: Critical CPU usage (%.1f%%)", m.Hostname, cpu))
			}
		}
		for _, al := range alerts {
			attentionLines = append(attentionLines, fmt.Sprintf("- **Alert (%s)**: %s", al.Severity, al.Message))
		}
		if len(attentionLines) == 0 {
			answer = "All systems are running healthy 🟢. No machines need immediate attention."
		} else {
			answer = fmt.Sprintf("The following items require your attention:\n\n%s", strings.Join(attentionLines, "\n"))
		}

	default:
		matched = false
	}

	// 2. LLM Fallback (if rule matching is not triggered and GEMINI_API_KEY is defined)
	apiKey := os.Getenv("GEMINI_API_KEY")
	if !matched && apiKey != "" {
		llmResponse, err := queryGemini(apiKey, systemContext, question)
		if err == nil && llmResponse != "" {
			answer = llmResponse
			matched = true
		} else {
			log.Printf("[AIChat] LLM query failed: %v", err)
		}
	}

	// 3. Fallback response if no LLM key and no rule matched
	if !matched {
		answer = fmt.Sprintf("I parsed your question: *\"%s\"*.\n\nHere is the current system telemetry summary:\n\n%s\n\nIf you want more detailed analysis, please configure `GEMINI_API_KEY` in the backend `.env` file.", question, systemContext)
	}

	c.JSON(http.StatusOK, ChatResponse{Answer: answer})
}

func queryGemini(apiKey string, context string, question string) (string, error) {
	prompt := fmt.Sprintf("You are InfraPilot AI, an expert system engineer assistant.\n\n"+
		"Context of the monitored systems:\n%s\n\n"+
		"User Question:\n%s\n\n"+
		"Please analyze the context and provide a helpful, actionable, and human-friendly response.",
		context, question)

	reqPayload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
	}

	reqBody, _ := json.Marshal(reqPayload)
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gemini returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("empty response candidates from Gemini")
}

// GetIncidentAnalysis generates dynamic root cause analysis, timeline, and remediation commands
func GetIncidentAnalysis(c *gin.Context) {
	alertIDStr := c.Param("alertId")
	if alertIDStr == "" {
		alertIDStr = c.Param("id")
	}
	alertUUID, err := uuid.Parse(alertIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	var alert models.LinuxAlert
	if err := database.DB.First(&alert, "id = ?", alertUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "alert not found"})
		return
	}

	var machine models.Machine
	database.DB.First(&machine, "id = ?", alert.MachineID)

	rootCause := "General system threshold breach."
	impact := "Performance degradation and potential service interruption."
	recommendation := "Investigate recent work surges or workload scheduling."
	commandAction := ""

	var topProc models.LinuxProcess
	errProc := database.DB.Where("machine_id = ?", alert.MachineID).Order("cpu_percent desc").First(&topProc).Error

	isWindows := strings.Contains(strings.ToLower(machine.OS), "win") || strings.Contains(strings.ToLower(machine.Platform), "win")

	switch alert.Type {
	case "CPU":
		if errProc == nil && topProc.CPUPercent > 5.0 {
			rootCause = fmt.Sprintf("High CPU usage caused by process: %s (PID: %d) consuming %.1f%% CPU.", topProc.Name, topProc.PID, topProc.CPUPercent)
			recommendation = fmt.Sprintf("Terminate process %s or review application thread pools.", topProc.Name)
			if isWindows {
				commandAction = fmt.Sprintf("taskkill /PID %d /F", topProc.PID)
			} else {
				commandAction = fmt.Sprintf("kill -9 %d", topProc.PID)
			}
		} else {
			rootCause = "High CPU utilization across multiple system threads."
			recommendation = "Review system resource limits or scale CPU capacity."
		}

	case "MEMORY":
		var topMemProc models.LinuxProcess
		errMem := database.DB.Where("machine_id = ?", alert.MachineID).Order("memory_percent desc").First(&topMemProc).Error
		if errMem == nil {
			rootCause = fmt.Sprintf("High memory usage caused by process: %s (PID: %d) occupying %.1f%% of RAM.", topMemProc.Name, topMemProc.PID, topMemProc.MemoryPercent)
			recommendation = fmt.Sprintf("Restart application process %s to release leaked memory buffers.", topMemProc.Name)
			if isWindows {
				commandAction = fmt.Sprintf("taskkill /PID %d /F", topMemProc.PID)
			} else {
				commandAction = fmt.Sprintf("kill -9 %d", topMemProc.PID)
			}
		} else {
			rootCause = "Memory allocation bottleneck detected."
			recommendation = "Check for background daemons memory leaks or expand physical RAM."
		}

	case "DISK":
		rootCause = "Disk storage allocation exceeded critical warning thresholds."
		recommendation = "Clean up temporary file folders and rotate logs."
		if isWindows {
			commandAction = "cleanmgr /sagerun:1"
		} else {
			commandAction = "sudo rm -rf /tmp/* && sudo find /var/log -type f -name '*.log' -exec truncate -s 0 {} \\;"
		}

	case "machine_offline":
		rootCause = "InfraPilot monitoring daemon stopped transmitting metrics / heartbeats."
		recommendation = "Check network gateway or restart the monitoring agent process."
		if isWindows {
			commandAction = "powershell -Command \"Start-Service -Name 'InfraPilot'\""
		} else {
			commandAction = "sudo systemctl start infrapilot-agent"
		}
	}

	timeline := []string{
		fmt.Sprintf("%s: Incident triggered under alert type '%s'", alert.CreatedAt.Format("15:04:05"), alert.Type),
		fmt.Sprintf("%s: Database registered status '%s'", time.Now().Add(-2*time.Minute).Format("15:04:05"), alert.Status),
		fmt.Sprintf("%s: AI Incident Analysis Complete", time.Now().Format("15:04:05")),
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             alert.ID.String(),
		"machine_id":     alert.MachineID.String(),
		"hostname":       machine.Hostname,
		"root_cause":     rootCause,
		"impact":         impact,
		"recommendation": recommendation,
		"command_action": commandAction,
		"timeline":       timeline,
	})
}
