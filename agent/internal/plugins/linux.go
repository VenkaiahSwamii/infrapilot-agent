package plugins

import (
	"net"
	"runtime"
	"sync"
	"time"

	"infrapilot/agent/internal/metrics"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	gnet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// LinuxPlugin collects Linux/Unix system metrics
type LinuxPlugin struct {
	enabled bool
	mu      sync.RWMutex
	health  PluginHealth
}

// NewLinuxPlugin creates a new Linux plugin instance
func NewLinuxPlugin() *LinuxPlugin {
	return &LinuxPlugin{
		enabled: true,
		health:  NewPluginHealth("linux"),
	}
}

// Name returns the plugin name
func (p *LinuxPlugin) Name() string {
	return "linux"
}

// IsEnabled returns whether the plugin is enabled
func (p *LinuxPlugin) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// SetEnabled enables or disables the plugin
func (p *LinuxPlugin) SetEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = enabled
}

// HealthCheck returns the health status of the plugin
func (p *LinuxPlugin) HealthCheck() PluginHealth {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.health
}

// Collect gathers Linux system metrics
func (p *LinuxPlugin) Collect() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	startTime := time.Now()

	result := make(map[string]interface{})

	// Collect host information
	hostInfo, err := p.collectHostInfo()
	if err != nil {
		p.health.Status = "degraded"
		p.health.Error = err.Error()
	} else {
		p.health.Status = "healthy"
		p.health.Error = ""
		for k, v := range hostInfo {
			result[k] = v
		}
	}

	// Collect CPU metrics
	cpuMetrics, err := p.collectCPUMetrics()
	if err != nil {
		// Non-critical, continue
		result["cpu_error"] = err.Error()
	} else {
		for k, v := range cpuMetrics {
			result[k] = v
		}
	}

	// Collect memory metrics
	memMetrics, err := p.collectMemoryMetrics()
	if err != nil {
		result["memory_error"] = err.Error()
	} else {
		for k, v := range memMetrics {
			result[k] = v
		}
	}

	// Collect disk metrics
	diskMetrics, err := p.collectDiskMetrics()
	if err != nil {
		result["disk_error"] = err.Error()
	} else {
		result["disk"] = diskMetrics
	}

	// Collect network metrics
	netMetrics, err := p.collectNetworkMetrics()
	if err != nil {
		result["network_error"] = err.Error()
	} else {
		for k, v := range netMetrics {
			result[k] = v
		}
	}

	// Collect load averages
	loadMetrics, err := p.collectLoadAverages()
	if err != nil {
		result["load_error"] = err.Error()
	} else {
		for k, v := range loadMetrics {
			result[k] = v
		}
	}

	// Collect filesystem information
	filesystems, err := p.collectFilesystems()
	if err != nil {
		result["filesystems_error"] = err.Error()
	} else {
		result["filesystems"] = filesystems
	}

	// Collect process information
	processes, err := p.collectProcesses()
	if err != nil {
		result["processes_error"] = err.Error()
	} else {
		result["processes"] = processes
	}

	// Collect network interfaces
	interfaces, err := p.collectNetworkInterfaces()
	if err != nil {
		result["interfaces_error"] = err.Error()
	} else {
		result["network_interfaces"] = interfaces
	}

	p.health.LastRun = time.Now()
	p.health.Metadata["duration_ms"] = time.Since(startTime).String()

	return result, nil
}

// collectHostInfo gathers host system information
func (p *LinuxPlugin) collectHostInfo() (map[string]interface{}, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	ip := ""
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil {
			ip = ipnet.IP.String()
			break
		}
	}

	return map[string]interface{}{
		"hostname":     info.Hostname,
		"ip_address":   ip,
		"os":           info.OS,
		"platform":     info.Platform,
		"platform_ver": info.PlatformVersion,
		"uptime":       info.Uptime,
		"kernel":       info.KernelVersion,
		"architecture": runtime.GOARCH,
		"timezone":     time.Now().Location().String(),
		"boot_time":    info.BootTime,
	}, nil
}

// collectCPUMetrics gathers CPU-related metrics
func (p *LinuxPlugin) collectCPUMetrics() (map[string]interface{}, error) {
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	cpuPerCore, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}

	cpuInfos, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var cpuVal float64
	if len(cpuUsage) > 0 {
		cpuVal = cpuUsage[0]
	}

	cpuFrequency := 0.0
	if len(cpuInfos) > 0 {
		cpuFrequency = cpuInfos[0].Mhz
	}

	return map[string]interface{}{
		"cpu_usage":         cpuVal,
		"cpu_temperature":   p.getCpuTemperature(),
		"cpu_cores":         runtime.NumCPU(),
		"cpu_per_core":      cpuPerCore,
		"cpu_frequency_mhz": cpuFrequency,
	}, nil
}

// collectMemoryMetrics gathers memory-related metrics
func (p *LinuxPlugin) collectMemoryMetrics() (map[string]interface{}, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"memory_total":   memInfo.Total,
		"memory_used":    memInfo.Used,
		"memory_free":    memInfo.Available,
		"memory_cached":  memInfo.Cached,
		"memory_percent": memInfo.UsedPercent,
		"swap_usage":     swapInfo.UsedPercent,
	}, nil
}

// collectDiskMetrics gathers disk-related metrics
func (p *LinuxPlugin) collectDiskMetrics() (map[string]interface{}, error) {
	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "C:"
	}

	diskInfo, err := disk.Usage(diskPath)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"disk_total":     diskInfo.Total,
		"disk_used":      diskInfo.Used,
		"disk_percent":   diskInfo.UsedPercent,
		"disk_read_bps":  0.0,
		"disk_write_bps": 0.0,
		"disk_iops":      0.0,
	}, nil
}

// collectNetworkMetrics gathers network-related metrics
func (p *LinuxPlugin) collectNetworkMetrics() (map[string]interface{}, error) {
	netInfo, err := gnet.IOCounters(false)
	if err != nil {
		return nil, err
	}

	var bytesSent, bytesRecv uint64
	if len(netInfo) > 0 {
		bytesSent = netInfo[0].BytesSent
		bytesRecv = netInfo[0].BytesRecv
	}

	upload, download, _ := GetNetworkSpeed()

	return map[string]interface{}{
		"bytes_sent":     bytesSent,
		"bytes_received": bytesRecv,
		"upload_mbps":    upload,
		"download_mbps":  download,
		"latency_ms":     0.0,
		"packet_loss":    0.0,
	}, nil
}

// collectLoadAverages gathers system load averages
func (p *LinuxPlugin) collectLoadAverages() (map[string]interface{}, error) {
	loadInfo, err := load.Avg()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"load_1":  loadInfo.Load1,
		"load_5":  loadInfo.Load5,
		"load_15": loadInfo.Load15,
	}, nil
}

// collectFilesystems gathers filesystem information
func (p *LinuxPlugin) collectFilesystems() ([]metrics.FilesystemMetric, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	result := make([]metrics.FilesystemMetric, 0, len(partitions))
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}
		result = append(result, metrics.FilesystemMetric{
			MountPoint: partition.Mountpoint,
			FSType:     partition.Fstype,
			Total:      usage.Total,
			Used:       usage.Used,
			Free:       usage.Free,
			UsedPct:    usage.UsedPercent,
		})
	}

	return result, nil
}

// collectProcesses gathers process information
func (p *LinuxPlugin) collectProcesses() ([]metrics.ProcessMetric, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	result := make([]metrics.ProcessMetric, 0, 15)
	for _, proc := range processes {
		name, _ := proc.Name()
		user, _ := proc.Username()
		cpuPct, _ := proc.CPUPercent()
		memPct, _ := proc.MemoryPercent()
		cmd, _ := proc.Cmdline()
		statusValues, _ := proc.Status()
		status := ""
		if len(statusValues) > 0 {
			status = statusValues[0]
		}
		result = append(result, metrics.ProcessMetric{
			PID:           proc.Pid,
			Name:          name,
			User:          user,
			CPUPercent:    cpuPct,
			MemoryPercent: memPct,
			Status:        status,
			Command:       cmd,
		})
	}

	// Sort by CPU usage
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].CPUPercent > result[i].CPUPercent {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if len(result) > 15 {
		return result[:15], nil
	}
	return result, nil
}

// collectNetworkInterfaces gathers network interface information
func (p *LinuxPlugin) collectNetworkInterfaces() ([]metrics.NetworkInterfaceMetric, error) {
	interfaces, err := gnet.Interfaces()
	if err != nil {
		return nil, err
	}

	result := make([]metrics.NetworkInterfaceMetric, 0, len(interfaces))
	for _, item := range interfaces {
		addresses := make([]string, 0, len(item.Addrs))
		for _, addr := range item.Addrs {
			addresses = append(addresses, addr.Addr)
		}
		result = append(result, metrics.NetworkInterfaceMetric{
			Name:      item.Name,
			MAC:       item.HardwareAddr,
			Addresses: addresses,
			SpeedMbps: 0.0,
		})
	}

	return result, nil
}

// getCpuTemperature returns zero when no hardware sensor integration exists.
func (p *LinuxPlugin) getCpuTemperature() float64 {
	return 0
}

var previousSent uint64
var previousRecv uint64
var previousTime time.Time

// GetNetworkSpeed calculates network upload/download speeds
func GetNetworkSpeed() (float64, float64, error) {
	stats, err := gnet.IOCounters(false)
	if err != nil {
		return 0, 0, err
	}

	now := time.Now()

	if previousTime.IsZero() {
		previousSent = stats[0].BytesSent
		previousRecv = stats[0].BytesRecv
		previousTime = now
		return 0, 0, nil
	}

	elapsed := now.Sub(previousTime).Seconds()
	if elapsed <= 0 {
		return 0, 0, nil
	}

	upload := float64(stats[0].BytesSent-previousSent) * 8 / elapsed / 1000000
	download := float64(stats[0].BytesRecv-previousRecv) * 8 / elapsed / 1000000

	previousSent = stats[0].BytesSent
	previousRecv = stats[0].BytesRecv
	previousTime = now

	return upload, download, nil
}
