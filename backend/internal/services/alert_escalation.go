package services

type AlertEscalation struct {
	notification *NotificationService
}

func NewAlertEscalation() *AlertEscalation {
	return &AlertEscalation{
		notification: NewNotificationService(),
	}
}

func (a *AlertEscalation) Escalate(alert string) {
	a.notification.SendEmailMsg(
		"Critical Infrastructure Alert",
		alert,
	)

	a.notification.SendSlackMsg(alert)

	a.notification.SendTeamsMsg(alert)
}
