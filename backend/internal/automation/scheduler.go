package automation

import (
	"fmt"
)

type AutomationScheduler struct{}

func NewAutomationScheduler() *AutomationScheduler {
	return &AutomationScheduler{}
}

func (s *AutomationScheduler) ScheduleJob(name, cronExpression, action string) error {
	fmt.Printf("[AUTOMATION SCHEDULER] Job '%s' registered with cron '%s' (action: %s)\n", name, cronExpression, action)
	return nil
}
