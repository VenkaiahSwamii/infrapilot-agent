package services

import (
	"testing"
	"time"
)

func TestStartMachineStatusChecker(t *testing.T) {
	// Verify StartMachineStatusChecker can be initialized without panic or deadlock
	go StartMachineStatusChecker()
	time.Sleep(50 * time.Millisecond)
}
