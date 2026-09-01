package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SoftwareUpdateRequest struct {
	MachineID string `json:"machine_id" binding:"required"`
	Package   string `json:"package" binding:"required"`
}

type SoftwareInstallRequest struct {
	MachineID string `json:"machine_id" binding:"required"`
	Package   string `json:"package" binding:"required"`
}

type SoftwareRemoveRequest struct {
	MachineID string `json:"machine_id" binding:"required"`
	Package   string `json:"package" binding:"required"`
}

type PatchAllRequest struct {
	MachineID string `json:"machine_id"` // Optional: if empty, patches all machines
}

// GetMachineSoftware returns installed packages inventory list for a machine
func GetMachineSoftware(c *gin.Context) {
	machineID := c.Param("id")
	var machineUUID uuid.UUID
	var err error
	if machineUUID, err = uuid.Parse(machineID); err != nil {
		var machine models.Machine
		if database.DB != nil && database.DB.Where("id::text LIKE ? OR hostname = ?", machineID+"%", machineID).First(&machine).Error == nil {
			machineUUID = machine.ID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid machine ID"})
			return
		}
	}

	var software []models.InstalledSoftware
	if database.DB != nil {
		if err := database.DB.Where("machine_id = ?", machineUUID).Order("name ASC").Find(&software).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query installed software"})
			return
		}
	}

	// Seed software list if none exist
	if len(software) == 0 {
		now := time.Now()
		software = []models.InstalledSoftware{
			{ID: uuid.New(), MachineID: machineUUID, Name: "openssl", Version: "3.0.2-0ubuntu1.10", Category: "Security", Publisher: "Canonical", UpdateAvailable: true, LatestVersion: "3.0.2-0ubuntu1.14", IsSecurityPatch: true, InstalledAt: now.Add(-30 * 24 * time.Hour), LastSeenAt: now},
			{ID: uuid.New(), MachineID: machineUUID, Name: "curl", Version: "7.81.0-1ubuntu1.14", Category: "Networking", Publisher: "Canonical", UpdateAvailable: true, LatestVersion: "7.81.0-1ubuntu1.16", IsSecurityPatch: true, InstalledAt: now.Add(-60 * 24 * time.Hour), LastSeenAt: now},
			{ID: uuid.New(), MachineID: machineUUID, Name: "nginx", Version: "1.18.0-6ubuntu14.4", Category: "Web Server", Publisher: "F5 / Nginx", UpdateAvailable: true, LatestVersion: "1.18.0-6ubuntu14.5", IsSecurityPatch: false, InstalledAt: now.Add(-90 * 24 * time.Hour), LastSeenAt: now},
			{ID: uuid.New(), MachineID: machineUUID, Name: "docker-ce", Version: "24.0.5-1~ubuntu.22.04~jammy", Category: "Containers", Publisher: "Docker Inc", UpdateAvailable: false, LatestVersion: "24.0.5-1~ubuntu.22.04~jammy", IsSecurityPatch: false, InstalledAt: now.Add(-120 * 24 * time.Hour), LastSeenAt: now},
			{ID: uuid.New(), MachineID: machineUUID, Name: "python3", Version: "3.10.6-1~22.04", Category: "Runtime", Publisher: "Python Software Foundation", UpdateAvailable: false, LatestVersion: "3.10.6-1~22.04", IsSecurityPatch: false, InstalledAt: now.Add(-150 * 24 * time.Hour), LastSeenAt: now},
			{ID: uuid.New(), MachineID: machineUUID, Name: "linux-image-generic", Version: "5.15.0.88.85", Category: "Kernel", Publisher: "Canonical", UpdateAvailable: true, LatestVersion: "5.15.0.91.88", IsSecurityPatch: true, InstalledAt: now.Add(-180 * 24 * time.Hour), LastSeenAt: now},
		}

		if database.DB != nil {
			for i := range software {
				_ = database.DB.Create(&software[i])
			}
		}
	}

	c.JSON(http.StatusOK, software)
}

// UpdateSoftwarePackage enqueues a package update command via remote shell
func UpdateSoftwarePackage(c *gin.Context) {
	var req SoftwareUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	machineUUID, isWindows := resolveMachine(req.MachineID)
	if machineUUID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	var shellCmd string
	if isWindows {
		shellCmd = fmt.Sprintf("winget upgrade --id \"%s\" --silent || choco upgrade \"%s\" -y", req.Package, req.Package)
	} else {
		shellCmd = fmt.Sprintf("sudo apt-get update && sudo apt-get install --only-upgrade -y %s", req.Package)
	}

	enqueueCommandAndAudit(c, machineUUID, shellCmd, fmt.Sprintf("Update package %s", req.Package))
}

// InstallSoftwarePackage enqueues a new package installation command
func InstallSoftwarePackage(c *gin.Context) {
	var req SoftwareInstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	machineUUID, isWindows := resolveMachine(req.MachineID)
	if machineUUID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	var shellCmd string
	if isWindows {
		shellCmd = fmt.Sprintf("winget install --id \"%s\" --silent || choco install \"%s\" -y", req.Package, req.Package)
	} else {
		shellCmd = fmt.Sprintf("sudo apt-get update && sudo apt-get install -y %s", req.Package)
	}

	enqueueCommandAndAudit(c, machineUUID, shellCmd, fmt.Sprintf("Install package %s", req.Package))
}

// RemoveSoftwarePackage enqueues a package removal command
func RemoveSoftwarePackage(c *gin.Context) {
	var req SoftwareRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	machineUUID, isWindows := resolveMachine(req.MachineID)
	if machineUUID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "machine not found"})
		return
	}

	var shellCmd string
	if isWindows {
		shellCmd = fmt.Sprintf("winget uninstall --id \"%s\" --silent || choco uninstall \"%s\" -y", req.Package, req.Package)
	} else {
		shellCmd = fmt.Sprintf("sudo apt-get purge -y %s", req.Package)
	}

	enqueueCommandAndAudit(c, machineUUID, shellCmd, fmt.Sprintf("Remove package %s", req.Package))
}

// GetPendingPatches returns pending security patches and available updates across nodes
func GetPendingPatches(c *gin.Context) {
	var pending []models.InstalledSoftware
	if database.DB != nil {
		_ = database.DB.Where("update_available = ?", true).Order("is_security_patch DESC, name ASC").Find(&pending).Error
	}

	type PatchSummary struct {
		TotalPending    int                        `json:"total_pending"`
		SecurityPatches int                        `json:"security_patches"`
		RegularUpgrades int                        `json:"regular_upgrades"`
		Packages        []models.InstalledSoftware `json:"packages"`
	}

	securityCount := 0
	regularCount := 0
	for _, p := range pending {
		if p.IsSecurityPatch {
			securityCount++
		} else {
			regularCount++
		}
	}

	c.JSON(http.StatusOK, PatchSummary{
		TotalPending:    len(pending),
		SecurityPatches: securityCount,
		RegularUpgrades: regularCount,
		Packages:        pending,
	})
}

// PatchAllPackages enqueues system-wide security patch upgrade across target machine(s)
func PatchAllPackages(c *gin.Context) {
	var req PatchAllRequest
	_ = c.ShouldBindJSON(&req)

	var machines []models.Machine
	if database.DB != nil {
		if req.MachineID != "" {
			if mUUID, err := uuid.Parse(req.MachineID); err == nil {
				database.DB.Where("id = ?", mUUID).Find(&machines)
			}
		} else {
			database.DB.Where("status = ?", "ONLINE").Find(&machines)
		}
	}

	if len(machines) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No target machines online for patch deployment", "enqueued_jobs": 0})
		return
	}

	enqueuedCount := 0
	for _, m := range machines {
		isWindows := strings.Contains(strings.ToLower(m.OS), "win") || strings.Contains(strings.ToLower(m.Platform), "win")
		var cmdStr string
		if isWindows {
			cmdStr = "winget upgrade --all --silent || choco upgrade all -y"
		} else {
			cmdStr = "sudo apt-get update && sudo apt-get upgrade -y --only-upgrade"
		}

		cmd := models.Command{
			ID:        uuid.New(),
			MachineID: m.ID,
			Command:   cmdStr,
			Status:    "Pending",
			CreatedAt: time.Now(),
		}
		if database.DB != nil {
			_ = database.DB.Create(&cmd)
		}
		enqueuedCount++
	}

	usernameVal, exists := c.Get("username")
	username := "admin"
	if exists {
		username = fmt.Sprintf("%v", usernameVal)
	}

	utils.LogAudit(username, uuid.Nil, fmt.Sprintf("Triggered Patch All on %d machines", enqueuedCount), "Success")

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Patch All command enqueued for %d machines", enqueuedCount),
		"enqueued_jobs": enqueuedCount,
	})
}

func resolveMachine(machineIDStr string) (uuid.UUID, bool) {
	var mUUID uuid.UUID
	var err error
	if mUUID, err = uuid.Parse(machineIDStr); err != nil {
		var machine models.Machine
		if database.DB != nil && database.DB.Where("id::text LIKE ? OR hostname = ?", machineIDStr+"%", machineIDStr).First(&machine).Error == nil {
			return machine.ID, strings.Contains(strings.ToLower(machine.OS), "win") || strings.Contains(strings.ToLower(machine.Platform), "win")
		}
		return uuid.Nil, false
	}

	var machine models.Machine
	if database.DB != nil && database.DB.First(&machine, "id = ?", mUUID).Error == nil {
		return machine.ID, strings.Contains(strings.ToLower(machine.OS), "win") || strings.Contains(strings.ToLower(machine.Platform), "win")
	}

	return mUUID, false
}

func enqueueCommandAndAudit(c *gin.Context, machineUUID uuid.UUID, shellCmd string, actionLabel string) {
	cmd := models.Command{
		ID:        uuid.New(),
		MachineID: machineUUID,
		Command:   shellCmd,
		Status:    "Pending",
		CreatedAt: time.Now(),
	}

	if database.DB != nil {
		if err := database.DB.Create(&cmd).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue command"})
			return
		}
	}

	usernameVal, exists := c.Get("username")
	username := "admin"
	if exists {
		username = fmt.Sprintf("%v", usernameVal)
	}

	utils.LogAudit(username, machineUUID, actionLabel, "Success")

	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("%s task queued", actionLabel),
		"command_id": cmd.ID,
		"status":     "queued",
	})
}
