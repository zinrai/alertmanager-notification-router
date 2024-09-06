package main

import (
	"log"
	"net/http"

	"github.com/zinrai/alertmanager-notification-router/internal/infrastructure/repository"
	"github.com/zinrai/alertmanager-notification-router/internal/interface/handler"
	"github.com/zinrai/alertmanager-notification-router/internal/usecase"
	"github.com/zinrai/alertmanager-notification-router/pkg/logger"
)

func main() {
	logger := logger.NewLogger()
	repo := repository.NewAlertRepository(logger)
	useCase := usecase.NewAlertUseCase(repo)
	ahmHandler := handler.NewAHMHandler(useCase)

	http.HandleFunc("/ahm", ahmHandler.HandleAHM)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
