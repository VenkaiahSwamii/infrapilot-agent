package alerts

import (
	"fmt"
	"strings"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"
)

// CheckStorageHealth evaluates filesystems, disk temperatures, SMART status, LVM status, and RAID array health.
func CheckStorageHealth(machine models.Machine, metric models.Metric, input services.SaveMetricInput) []models.LinuxAlert {
	var alerts []models.LinuxAlert

	// 1. Filesystem Capacity Check (Per Mount Point)
	if len(input.Filesystems) > 0 {
		for _, fs := range input.Filesystems {
			pct := fs.UsedPct
			if pct > 95.0 {
				alerts = append(alerts, models.LinuxAlert{
					Title:              "Filesystem Full",
					Description:        fmt.Sprintf("Filesystem '%s' usage critical: %.1f%%", fs.MountPoint, pct),
					Category:           "Filesystem",
					Component:          fs.MountPoint,
					Source:             "StorageChecker",
					Type:               "disk_percent",
					Severity:           "Critical",
					Priority:           models.MapSeverityToPriority("Critical"),
					Message:            fmt.Sprintf("Filesystem %s Full: %.1f%% > 95.0%%", fs.MountPoint, pct),
					MetricValue:        pct,
					Threshold:          95.0,
					Status:             "OPEN",
					RecoverySuggestion: fmt.Sprintf("Clean log files or temp files on '%s' (`df -h %s`, `du -sh %s/*`). Expand disk volume if needed.", fs.MountPoint, fs.MountPoint, fs.MountPoint),
				})
			} else if pct > 90.0 {
				alerts = append(alerts, models.LinuxAlert{
					Title:              "Filesystem Usage High",
					Description:        fmt.Sprintf("Filesystem '%s' usage warning: %.1f%%", fs.MountPoint, pct),
					Category:           "Filesystem",
					Component:          fs.MountPoint,
					Source:             "StorageChecker",
					Type:               "disk_percent",
					Severity:           "Warning",
					Priority:           models.MapSeverityToPriority("Warning"),
					Message:            fmt.Sprintf("Filesystem %s High: %.1f%% > 90.0%%", fs.MountPoint, pct),
					MetricValue:        pct,
					Threshold:          90.0,
					Status:             "OPEN",
					RecoverySuggestion: fmt.Sprintf("Monitor disk growth on mount point '%s'. Clean temporary files or log rotate.", fs.MountPoint),
				})
			}
		}
	} else {
		// Fallback to overall DiskPercent metric check
		diskPct := input.DiskPercent
		if diskPct == 0 {
			diskPct = metric.DiskPercent
		}
		if diskPct == 0 {
			diskPct = metric.DiskUsage
		}

		if diskPct > 95.0 {
			alerts = append(alerts, models.LinuxAlert{
				Title:              "Disk Almost Full",
				Description:        fmt.Sprintf("Overall disk space usage critical: %.1f%%", diskPct),
				Category:           "Storage",
				Component:          "root_disk",
				Source:             "StorageChecker",
				Type:               "disk_percent",
				Severity:           "Critical",
				Priority:           models.MapSeverityToPriority("Critical"),
				Message:            fmt.Sprintf("Disk Usage Critical: %.1f%% > 95.0%%", diskPct),
				MetricValue:        diskPct,
				Threshold:          95.0,
				Status:             "OPEN",
				RecoverySuggestion: "Purge package manager caches (`apt clean` or `yum clean all`) and rotate `/var/log`.",
			})
		} else if diskPct > 90.0 {
			alerts = append(alerts, models.LinuxAlert{
				Title:              "Disk Usage High",
				Description:        fmt.Sprintf("Overall disk space usage warning: %.1f%%", diskPct),
				Category:           "Storage",
				Component:          "root_disk",
				Source:             "StorageChecker",
				Type:               "disk_percent",
				Severity:           "Warning",
				Priority:           models.MapSeverityToPriority("Warning"),
				Message:            fmt.Sprintf("Disk Usage Warning: %.1f%% > 90.0%%", diskPct),
				MetricValue:        diskPct,
				Threshold:          90.0,
				Status:             "OPEN",
				RecoverySuggestion: "Perform disk cleanup or schedule storage volume expansion.",
			})
		}
	}

	// 2. Disk Temperature Check
	if input.DiskTemperature > 70.0 {
		alerts = append(alerts, models.LinuxAlert{
			Title:              "Disk Temperature Critical",
			Description:        fmt.Sprintf("Storage drive temperature dangerous: %.1f°C", input.DiskTemperature),
			Category:           "Storage",
			Component:          "disk_hardware",
			Source:             "StorageChecker",
			Type:               "disk_temperature",
			Severity:           "Critical",
			Priority:           models.MapSeverityToPriority("Critical"),
			Message:            fmt.Sprintf("Disk Temperature Critical: %.1f°C > 70.0°C", input.DiskTemperature),
			MetricValue:        input.DiskTemperature,
			Threshold:          70.0,
			Status:             "OPEN",
			RecoverySuggestion: "Check drive enclosure ventilation and drive bay cooling to prevent drive failure.",
		})
	} else if input.DiskTemperature > 60.0 {
		alerts = append(alerts, models.LinuxAlert{
			Title:              "Disk Temperature High",
			Description:        fmt.Sprintf("Storage drive temperature elevated: %.1f°C", input.DiskTemperature),
			Category:           "Storage",
			Component:          "disk_hardware",
			Source:             "StorageChecker",
			Type:               "disk_temperature",
			Severity:           "Warning",
			Priority:           models.MapSeverityToPriority("Warning"),
			Message:            fmt.Sprintf("Disk Temperature Warning: %.1f°C > 60.0°C", input.DiskTemperature),
			MetricValue:        input.DiskTemperature,
			Threshold:          60.0,
			Status:             "OPEN",
			RecoverySuggestion: "Monitor drive bay airflow.",
		})
	}

	// 3. SMART Status Check
	smart := strings.ToLower(input.SmartStatus)
	if strings.Contains(smart, "failed") || strings.Contains(smart, "failing") || strings.Contains(smart, "bad") || strings.Contains(smart, "error") {
		alerts = append(alerts, models.LinuxAlert{
			Title:              "SMART Failed",
			Description:        fmt.Sprintf("Disk SMART self-diagnostic report indicates failure: '%s'", input.SmartStatus),
			Category:           "Storage",
			Component:          "smart_health",
			Source:             "StorageChecker",
			Type:               "smart_status",
			Severity:           "Critical",
			Priority:           models.MapSeverityToPriority("Critical"),
			Message:            fmt.Sprintf("SMART Disk Diagnostics FAILED: %s", input.SmartStatus),
			Status:             "OPEN",
			RecoverySuggestion: "BACKUP DATA IMMEDIATELY! Schedule immediate physical drive replacement.",
		})
	}

	// 4. LVM Degraded Check
	lvm := strings.ToLower(input.LVMStatus)
	if strings.Contains(lvm, "degraded") || strings.Contains(lvm, "failed") || strings.Contains(lvm, "error") {
		alerts = append(alerts, models.LinuxAlert{
			Title:              "LVM Degraded",
			Description:        fmt.Sprintf("Logical Volume Manager reports volume degradation: '%s'", input.LVMStatus),
			Category:           "Storage",
			Component:          "lvm_volume",
			Source:             "StorageChecker",
			Type:               "lvm_status",
			Severity:           "Critical",
			Priority:           models.MapSeverityToPriority("Critical"),
			Message:            fmt.Sprintf("LVM Volume Degraded: %s", input.LVMStatus),
			Status:             "OPEN",
			RecoverySuggestion: "Inspect LVM volume group health (`vgdisplay`, `lvdisplay`) and replace faulty physical volumes.",
		})
	}

	// 5. RAID Degraded Check
	raid := strings.ToLower(input.RAIDStatus)
	if strings.Contains(raid, "degraded") || strings.Contains(raid, "failed") || strings.Contains(raid, "sync_failed") {
		alerts = append(alerts, models.LinuxAlert{
			Title:              "RAID Degraded",
			Description:        fmt.Sprintf("RAID array reports degraded or failing status: '%s'", input.RAIDStatus),
			Category:           "Storage",
			Component:          "raid_array",
			Source:             "StorageChecker",
			Type:               "raid_status",
			Severity:           "Critical",
			Priority:           models.MapSeverityToPriority("Critical"),
			Message:            fmt.Sprintf("RAID Array Degraded: %s", input.RAIDStatus),
			Status:             "OPEN",
			RecoverySuggestion: "Check RAID array status (`cat /proc/mdstat` or `mdadm --detail /dev/md0`). Replace failed member disk.",
		})
	}

	return alerts
}
