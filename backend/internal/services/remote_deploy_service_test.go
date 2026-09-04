package services

import (
	"context"
	"testing"
	"time"
)

func TestExtractKey(t *testing.T) {
	sampleOutput := `HOSTNAME=web-node-01
OS=Linux
DISTRO="Ubuntu 22.04 LTS"
ARCH=x86_64
HAS_SYSTEMD=yes`

	if got := extractKey(sampleOutput, "HOSTNAME", "fallback"); got != "web-node-01" {
		t.Errorf("expected web-node-01, got %s", got)
	}
	if got := extractKey(sampleOutput, "DISTRO", "fallback"); got != "Ubuntu 22.04 LTS" {
		t.Errorf("expected 'Ubuntu 22.04 LTS', got %s", got)
	}
	if got := extractKey(sampleOutput, "NON_EXISTENT", "default_val"); got != "default_val" {
		t.Errorf("expected default_val, got %s", got)
	}
}

func TestRemoteDeployService_BuildSSHConfig(t *testing.T) {
	svc := NewRemoteDeployService(nil)

	// 1. Password auth
	target := RemoteDeployTarget{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "ubuntu",
		Password: "SecretPassword123",
	}

	cfg, err := svc.buildSSHConfig(target)
	if err != nil {
		t.Fatalf("unexpected error building ssh config: %v", err)
	}
	if cfg.User != "ubuntu" {
		t.Errorf("expected user 'ubuntu', got %s", cfg.User)
	}
	if len(cfg.Auth) == 0 {
		t.Errorf("expected at least 1 auth method")
	}

	// 2. Missing credentials
	targetNoCreds := RemoteDeployTarget{
		Host:     "192.168.1.100",
		Username: "ubuntu",
	}
	_, err = svc.buildSSHConfig(targetNoCreds)
	if err == nil {
		t.Errorf("expected error when no password or key is provided")
	}
}

func TestRemoteDeployService_History(t *testing.T) {
	svc := NewRemoteDeployService(nil)

	res := &RemoteDeployResult{
		DeploymentID: "dep_test_123",
		Success:      true,
		Host:         "10.0.0.5",
		Hostname:     "db-primary",
		CreatedAt:    time.Now(),
	}

	svc.saveHistory(res)

	retrieved := svc.GetDeployment("dep_test_123")
	if retrieved == nil {
		t.Fatalf("expected to retrieve saved deployment")
	}
	if retrieved.Hostname != "db-primary" {
		t.Errorf("expected hostname 'db-primary', got %s", retrieved.Hostname)
	}

	history := svc.GetDeploymentHistory()
	if len(history) != 1 {
		t.Errorf("expected 1 history entry, got %d", len(history))
	}
}

func TestRemoteDeployService_UnreachableHostTest(t *testing.T) {
	svc := NewRemoteDeployService(nil)

	// Port 59999 on 127.0.0.1 should fail cleanly without crashing
	target := RemoteDeployTarget{
		Host:     "127.0.0.1",
		Port:     59999,
		Username: "root",
		Password: "password",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := svc.TestConnection(ctx, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Success {
		t.Errorf("expected connection to fail on unreachable port 59999")
	}
	if res.Error == "" {
		t.Errorf("expected non-empty error message for unreachable host")
	}
}
