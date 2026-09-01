package events

import "time"

// CommandExecutedEvent is published when a command is executed on a machine
type CommandExecutedEvent struct {
	MachineID  string
	Command    string
	Output     string
	ExitCode   int
	ExecutedBy string
	Time       time.Time
}

func (e CommandExecutedEvent) Name() string {
	return "command.executed"
}

func (e CommandExecutedEvent) Timestamp() time.Time {
	return e.Time
}

// TerminalOpenedEvent is published when a terminal session is opened
type TerminalOpenedEvent struct {
	MachineID    string
	SessionID    string
	UserID       string
	Username     string
	TerminalType string
	Time         time.Time
}

func (e TerminalOpenedEvent) Name() string {
	return "terminal.opened"
}

func (e TerminalOpenedEvent) Timestamp() time.Time {
	return e.Time
}

// TerminalClosedEvent is published when a terminal session is closed
type TerminalClosedEvent struct {
	MachineID string
	SessionID string
	UserID    string
	Username  string
	Duration  int64 // seconds
	Time      time.Time
}

func (e TerminalClosedEvent) Name() string {
	return "terminal.closed"
}

func (e TerminalClosedEvent) Timestamp() time.Time {
	return e.Time
}

// FileUploadedEvent is published when a file is uploaded to a machine
type FileUploadedEvent struct {
	MachineID  string
	FilePath   string
	FileName   string
	FileSize   int64
	UploadedBy string
	Time       time.Time
}

func (e FileUploadedEvent) Name() string {
	return "file.uploaded"
}

func (e FileUploadedEvent) Timestamp() time.Time {
	return e.Time
}

// FileDeletedEvent is published when a file is deleted from a machine
type FileDeletedEvent struct {
	MachineID string
	FilePath  string
	FileName  string
	DeletedBy string
	Time      time.Time
}

func (e FileDeletedEvent) Name() string {
	return "file.deleted"
}

func (e FileDeletedEvent) Timestamp() time.Time {
	return e.Time
}

// DockerContainerStoppedEvent is published when a Docker container stops
type DockerContainerStoppedEvent struct {
	MachineID     string
	ContainerID   string
	ContainerName string
	Image         string
	ExitCode      int
	Time          time.Time
}

func (e DockerContainerStoppedEvent) Name() string {
	return "docker.container.stopped"
}

func (e DockerContainerStoppedEvent) Timestamp() time.Time {
	return e.Time
}

// KubernetesPodFailedEvent is published when a Kubernetes pod fails
type KubernetesPodFailedEvent struct {
	MachineID string
	PodName   string
	Namespace string
	Reason    string
	Message   string
	Time      time.Time
}

func (e KubernetesPodFailedEvent) Name() string {
	return "kubernetes.pod.failed"
}

func (e KubernetesPodFailedEvent) Timestamp() time.Time {
	return e.Time
}

// UserLoginEvent is published when a user logs in
type UserLoginEvent struct {
	UserID    string
	Username  string
	Email     string
	IPAddress string
	UserAgent string
	Time      time.Time
}

func (e UserLoginEvent) Name() string {
	return "user.login"
}

func (e UserLoginEvent) Timestamp() time.Time {
	return e.Time
}

// UserLogoutEvent is published when a user logs out
type UserLogoutEvent struct {
	UserID   string
	Username string
	Time     time.Time
}

func (e UserLogoutEvent) Name() string {
	return "user.logout"
}

func (e UserLogoutEvent) Timestamp() time.Time {
	return e.Time
}

// APIKeyRotatedEvent is published when a machine API key is rotated
type APIKeyRotatedEvent struct {
	MachineID string
	UserID    string
	Time      time.Time
}

func (e APIKeyRotatedEvent) Name() string {
	return "api_key.rotated"
}

func (e APIKeyRotatedEvent) Timestamp() time.Time {
	return e.Time
}
