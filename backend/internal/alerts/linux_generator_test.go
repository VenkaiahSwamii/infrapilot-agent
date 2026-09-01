package alerts

import (
	"testing"

	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/services"

	"github.com/google/uuid"
)

func TestGenerateLinuxAlerts_CPUAndMemory(t *testing.T) {
	machine := models.Machine{
		ID:       uuid.New(),
		Hostname: "ubuntu-server-01",
	}
	metric := models.Metric{
		ID:        uuid.New(),
		MachineID: machine.ID,
	}

	input := services.SaveMetricInput{
		CPUUsage:       97.5,
		CPUTemperature: 98.0,
		CPUCores:       4,
		Load1:          8.5,
		MemoryPercent:  96.0,
		SwapUsage:      85.0,
	}

	alerts := GenerateLinuxAlerts(machine, metric, input)
	if len(alerts) == 0 {
		t.Fatalf("Expected generated alerts for critical CPU & Memory payload, got 0")
	}

	// Phase 13 Incident Correlation should combine simultaneous critical CPU + Memory + Load into "Infrastructure Overloaded"
	hasCorrelatedIncident := false
	for _, alert := range alerts {
		if alert.IsCorrelated && alert.Title == "Infrastructure Overloaded" {
			hasCorrelatedIncident = true
			if alert.Priority != "P1" {
				t.Errorf("Expected Priority P1 for Infrastructure Overloaded incident, got %s", alert.Priority)
			}
			if alert.RecoverySuggestion == "" {
				t.Errorf("Expected non-empty recovery suggestion for correlated incident")
			}
		}
	}

	if !hasCorrelatedIncident {
		t.Errorf("Expected incident correlation engine to produce 'Infrastructure Overloaded' alert")
	}
}

func TestGenerateLinuxAlerts_StorageAndSMART(t *testing.T) {
	machine := models.Machine{
		ID:       uuid.New(),
		Hostname: "storage-node-01",
	}
	metric := models.Metric{
		ID:        uuid.New(),
		MachineID: machine.ID,
	}

	input := services.SaveMetricInput{
		DiskPercent: 96.5,
		Filesystems: []services.FilesystemInput{
			{MountPoint: "/", UsedPct: 97.2},
			{MountPoint: "/var", UsedPct: 99.0},
		},
		SmartStatus: "FAILED",
		RAIDStatus:  "Degraded",
		LVMStatus:   "Degraded",
	}

	alerts := GenerateLinuxAlerts(machine, metric, input)
	if len(alerts) < 3 {
		t.Fatalf("Expected multiple storage alerts (filesystem, SMART, RAID, LVM), got %d", len(alerts))
	}

	foundSMART := false
	foundRAID := false
	for _, alert := range alerts {
		if alert.Title == "SMART Failed" && alert.Category == "Storage" {
			foundSMART = true
		}
		if alert.Title == "RAID Degraded" && alert.Category == "Storage" {
			foundRAID = true
		}
	}

	if !foundSMART {
		t.Errorf("Expected 'SMART Failed' alert")
	}
	if !foundRAID {
		t.Errorf("Expected 'RAID Degraded' alert")
	}
}

func TestGenerateLinuxAlerts_ServicesAndDocker(t *testing.T) {
	machine := models.Machine{
		ID:       uuid.New(),
		Hostname: "docker-host-01",
	}
	metric := models.Metric{
		ID:        uuid.New(),
		MachineID: machine.ID,
	}

	input := services.SaveMetricInput{
		Services: []services.ServiceInput{
			{Name: "docker", Status: "stopped"},
			{Name: "nginx", Status: "failed"},
		},
		DockerContainers: []services.DockerContainerInput{
			{Name: "mysql-db", Status: "Exited (137)", State: "exited", RestartCount: 5},
		},
	}

	alerts := GenerateLinuxAlerts(machine, metric, input)

	foundDockerSvc := false
	foundContainerExited := false
	foundRestartLoop := false

	for _, alert := range alerts {
		if alert.Title == "Service Down" && alert.Component == "docker" {
			foundDockerSvc = true
		}
		if alert.Title == "Docker Container Exited" {
			foundContainerExited = true
		}
		if alert.Title == "Docker Container Restart Loop" {
			foundRestartLoop = true
		}
	}

	if !foundDockerSvc {
		t.Errorf("Expected 'Service Down' alert for docker service")
	}
	if !foundContainerExited {
		t.Errorf("Expected 'Docker Container Exited' alert")
	}
	if !foundRestartLoop {
		t.Errorf("Expected 'Docker Container Restart Loop' alert")
	}
}

func TestGenerateLinuxAlerts_KubernetesAndSecurity(t *testing.T) {
	machine := models.Machine{
		ID:       uuid.New(),
		Hostname: "k8s-node-01",
	}
	metric := models.Metric{
		ID:        uuid.New(),
		MachineID: machine.ID,
	}

	k8sPodsJSON := `[{"name":"api-server","namespace":"prod","status":"CrashLoopBackOff","reason":"OOMKilled","restarts":12}]`
	k8sNodesJSON := `[{"name":"k8s-node-01","status":"NotReady","ready":false}]`

	input := services.SaveMetricInput{
		K8sPodsJSON:  k8sPodsJSON,
		K8sNodesJSON: k8sNodesJSON,
		OpenPorts:    []string{"23", "21"},
		Services: []services.ServiceInput{
			{Name: "sshd", Status: "stopped"},
			{Name: "ufw", Status: "inactive"},
		},
	}

	alerts := GenerateLinuxAlerts(machine, metric, input)

	foundPodCrash := false
	foundNodeNotReady := false
	foundSSHDown := false
	foundInsecurePort := false

	for _, alert := range alerts {
		if alert.Title == "Kubernetes Pod CrashLoopBackOff" && alert.Category == "Kubernetes" {
			foundPodCrash = true
		}
		if alert.Title == "Kubernetes Node NotReady" && alert.Category == "Kubernetes" {
			foundNodeNotReady = true
		}
		if alert.Title == "SSH Service Down" && alert.Category == "Security" {
			foundSSHDown = true
		}
		if alert.Title == "Insecure Open Port Detected" && alert.Category == "Security" {
			foundInsecurePort = true
		}
	}

	if !foundPodCrash {
		t.Errorf("Expected Kubernetes Pod CrashLoopBackOff alert")
	}
	if !foundNodeNotReady {
		t.Errorf("Expected Kubernetes Node NotReady alert")
	}
	if !foundSSHDown {
		t.Errorf("Expected Security SSH Service Down alert")
	}
	if !foundInsecurePort {
		t.Errorf("Expected Security Insecure Open Port alert")
	}
}
