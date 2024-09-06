package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zinrai/alertmanager-notification-router/internal/domain"
	"github.com/zinrai/alertmanager-notification-router/pkg/logger"
)

type AlertRepository struct {
	logger     logger.Logger
	apiBaseURL string
	client     *http.Client
}

func NewAlertRepository(logger logger.Logger) *AlertRepository {
	return &AlertRepository{
		logger:     logger,
		apiBaseURL: "http://localhost:8000/api/alerts/",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *AlertRepository) SaveAlert(alert domain.Alert) error {
	jsonData, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("error marshaling alert: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.apiBaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending alert to API: %w", err)
	}
	defer func() {
		// Drain the body before closing
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	r.logger.Info(fmt.Sprintf("Alert saved successfully: %s", alert.Identifier))
	return nil
}
