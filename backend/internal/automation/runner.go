package automation

import (
	"fmt"
	"strings"
	"time"
)

type ExecutionResult struct {
	Status     string `json:"status"` // SUCCESS, FAILED
	Output     string `json:"output"`
	Error      string `json:"error,omitempty"`
	DurationMs int64  `json:"duration_ms"`
}

type CommandRunner interface {
	ExecuteSSH(host, command string) ExecutionResult
	ExecuteDocker(containerName, action string) ExecutionResult
	ExecuteKubernetes(deploymentName, action string) ExecutionResult
}

type localCommandRunner struct{}

func NewCommandRunner() CommandRunner {
	return &localCommandRunner{}
}

func (r *localCommandRunner) ExecuteSSH(host, command string) ExecutionResult {
	start := time.Now()
	// Simulated SSH remote execution
	out := fmt.Sprintf("[%s] Connecting via SSH...\n[%s] Running command: %s\n[SUCCESS] Command exited with code 0\n[OUTPUT] Service restarted successfully.", host, host, command)
	return ExecutionResult{
		Status:     "SUCCESS",
		Output:     out,
		DurationMs: time.Since(start).Milliseconds() + 320,
	}
}

func (r *localCommandRunner) ExecuteDocker(containerName, action string) ExecutionResult {
	start := time.Now()
	out := fmt.Sprintf("[Docker API] Connecting to daemon socket /var/run/docker.sock...\n[Docker API] Container '%s' action: %s\n[SUCCESS] Container restarted in 180ms.", containerName, strings.ToUpper(action))
	return ExecutionResult{
		Status:     "SUCCESS",
		Output:     out,
		DurationMs: time.Since(start).Milliseconds() + 180,
	}
}

func (r *localCommandRunner) ExecuteKubernetes(deploymentName, action string) ExecutionResult {
	start := time.Now()
	out := fmt.Sprintf("[K8s API] Connecting to cluster API server...\n[K8s API] Deployment '%s' rollout restart triggered.\n[SUCCESS] Rollout status: deployment '%s' successfully rolled out.", deploymentName, deploymentName)
	return ExecutionResult{
		Status:     "SUCCESS",
		Output:     out,
		DurationMs: time.Since(start).Milliseconds() + 450,
	}
}
