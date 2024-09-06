package usecase

import (
	"fmt"
	"strings"

	"github.com/zinrai/alertmanager-notification-router/internal/domain"
)

type AlertRepository interface {
	SaveAlert(alert domain.Alert) error
}

type AlertUseCase struct {
	repo AlertRepository
}

func NewAlertUseCase(repo AlertRepository) *AlertUseCase {
	return &AlertUseCase{repo: repo}
}

func (uc *AlertUseCase) ProcessAHM(webhook domain.AlertmanagerWebhook) error {
	for _, amAlert := range webhook.Alerts {
		alert := domain.Alert{
			Subject:    fmt.Sprintf("%s: %s", amAlert.Labels["alertname"], amAlert.Annotations["summary"]),
			Body:       buildAlertBody(amAlert),
			Identifier: amAlert.Fingerprint,
			Urgency:    determineUrgency(amAlert),
		}

		if err := uc.repo.SaveAlert(alert); err != nil {
			return err
		}
	}
	return nil
}

func buildAlertBody(alert domain.AlertmanagerAlert) string {
	var body strings.Builder

	body.WriteString(fmt.Sprintf("Status: %s\n", alert.Status))
	body.WriteString(fmt.Sprintf("Description: %s\n", alert.Annotations["description"]))
	body.WriteString(fmt.Sprintf("Starts At: %s\n", alert.StartsAt))
	body.WriteString("Labels:\n")
	for k, v := range alert.Labels {
		body.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	}

	return body.String()
}

func determineUrgency(alert domain.AlertmanagerAlert) string {
	severity, exists := alert.Labels["severity"]
	if !exists {
		return "MEDIUM"
	}

	switch strings.ToUpper(severity) {
	case "CRITICAL":
		return "HIGH"
	case "WARNING":
		return "MEDIUM"
	default:
		return "LOW"
	}
}
