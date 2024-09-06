package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zinrai/alertmanager-notification-router/internal/domain"
	"github.com/zinrai/alertmanager-notification-router/internal/usecase"
)

type AHMHandler struct {
	useCase *usecase.AlertUseCase
}

func NewAHMHandler(useCase *usecase.AlertUseCase) *AHMHandler {
	return &AHMHandler{useCase: useCase}
}

func (h *AHMHandler) HandleAHM(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var webhook domain.AlertmanagerWebhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.useCase.ProcessAHM(webhook); err != nil {
		http.Error(w, "Error processing AHM webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
