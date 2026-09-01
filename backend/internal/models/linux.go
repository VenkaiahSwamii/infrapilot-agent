package models

import (
	"time"

	"github.com/google/uuid"
)

type LinuxServer struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"machine_id"`
	Hostname     string    `json:"hostname"`
	OS           string    `json:"os"`
	Kernel       string    `json:"kernel"`
	Architecture string    `json:"architecture"`
	AgentVersion string    `json:"agent_version"`
	IPAddress    string    `json:"ip_address"`
	MACAddress   string    `json:"mac_address"`
	Timezone     string    `json:"timezone"`
	BootTime     uint64    `json:"boot_time"`
	LastSeenAt   time.Time `json:"last_seen_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LinuxMetric struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID       uuid.UUID `gorm:"type:uuid;index;not null" json:"machine_id"`
	SampledAt       time.Time `gorm:"index" json:"sampled_at"`
	CPUUsage        float64   `json:"cpu_usage"`
	CPUPerCoreJSON  string    `gorm:"type:jsonb" json:"cpu_per_core"`
	CPUFrequencyMHz float64   `json:"cpu_frequency_mhz"`
	CPUTemperature  float64   `json:"cpu_temperature"`
	Load1           float64   `json:"load_1"`
	Load5           float64   `json:"load_5"`
	Load15          float64   `json:"load_15"`
	MemoryUsed      uint64    `json:"memory_used"`
	MemoryFree      uint64    `json:"memory_free"`
	MemoryCached    uint64    `json:"memory_cached"`
	MemoryPercent   float64   `json:"memory_percent"`
	SwapUsage       float64   `json:"swap_usage"`
	DiskUsage       float64   `json:"disk_usage"`
	DiskReadBps     float64   `json:"disk_read_bps"`
	DiskWriteBps    float64   `json:"disk_write_bps"`
	DiskIOPS        float64   `json:"disk_iops"`
	UploadMbps      float64   `json:"upload_mbps"`
	DownloadMbps    float64   `json:"download_mbps"`
	LatencyMs       float64   `json:"latency_ms"`
	PacketLoss      float64   `json:"packet_loss"`
	CreatedAt       time.Time `json:"created_at"`
}

type LinuxProcess struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID     uuid.UUID `gorm:"type:uuid;index;not null" json:"machine_id"`
	SampledAt     time.Time `gorm:"index" json:"sampled_at"`
	PID           int32     `json:"pid"`
	Name          string    `json:"name"`
	User          string    `json:"user"`
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryPercent float32   `json:"memory_percent"`
	Status        string    `json:"status"`
	Command       string    `json:"command"`
}

type LinuxService struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_linux_service_machine_name;not null" json:"machine_id"`
	Name         string    `gorm:"uniqueIndex:idx_linux_service_machine_name" json:"name"`
	Status       string    `json:"status"`
	RestartCount int       `json:"restart_count"`
	LastSeenAt   time.Time `json:"last_seen_at"`
}

type LinuxNetwork struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID      uuid.UUID `gorm:"type:uuid;index;not null" json:"machine_id"`
	SampledAt      time.Time `gorm:"index" json:"sampled_at"`
	InterfacesJSON string    `gorm:"type:jsonb" json:"interfaces"`
	Connections    int       `json:"connections"`
	OpenPortsJSON  string    `gorm:"type:jsonb" json:"open_ports"`
	UploadMbps     float64   `json:"upload_mbps"`
	DownloadMbps   float64   `json:"download_mbps"`
	PacketLoss     float64   `json:"packet_loss"`
	LatencyMs      float64   `json:"latency_ms"`
}

type LinuxStorage struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID       uuid.UUID `gorm:"type:uuid;index;not null" json:"machine_id"`
	SampledAt       time.Time `gorm:"index" json:"sampled_at"`
	FilesystemsJSON string    `gorm:"type:jsonb" json:"filesystems"`
	DiskReadBps     float64   `json:"disk_read_bps"`
	DiskWriteBps    float64   `json:"disk_write_bps"`
	DiskIOPS        float64   `json:"disk_iops"`
	SmartStatus     string    `json:"smart_status"`
	RAIDStatus      string    `json:"raid_status"`
	LVMStatus       string    `json:"lvm_status"`
	DiskTemperature float64   `json:"disk_temperature"`
}

type HistoricalMetric struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID uuid.UUID `gorm:"type:uuid;index;not null" json:"machine_id"`
	Metric    string    `gorm:"index" json:"metric"`
	Value     float64   `json:"value"`
	SampledAt time.Time `gorm:"index" json:"sampled_at"`
}

type LinuxDocker struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID      uuid.UUID `gorm:"type:uuid;index;not null" json:"machine_id"`
	SampledAt      time.Time `gorm:"index" json:"sampled_at"`
	ContainersJSON string    `gorm:"type:jsonb" json:"containers"`
	ImagesJSON     string    `gorm:"type:jsonb" json:"images"`
}

type LinuxKubernetes struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID uuid.UUID `gorm:"type:uuid;index;not null" json:"machine_id"`
	SampledAt time.Time `gorm:"index" json:"sampled_at"`
	NodesJSON string    `gorm:"type:jsonb" json:"nodes"`
	PodsJSON  string    `gorm:"type:jsonb" json:"pods"`
}

type LinuxLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MachineID uuid.UUID `gorm:"type:uuid;index;not null" json:"machine_id"`
	Timestamp time.Time `gorm:"index" json:"timestamp"`
	Source    string    `json:"source"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}
