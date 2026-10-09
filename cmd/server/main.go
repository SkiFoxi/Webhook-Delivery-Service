package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	httphandler "github.com/SkiFoxi/Webhook-Delivery-Service/internal/handler/http"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/repository/memory"
	httpsender "github.com/SkiFoxi/Webhook-Delivery-Service/internal/sender"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/service"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/worker"
)

func main() {
	webhookRepo := memory.NewWebhookRepo()
	deliveryRepo := memory.NewDeliveryRepo()

	httpClient := &http.Client{Timeout: 10 * time.Second}
	deliverySender := httpsender.NewHTTPSender(httpClient, 10*time.Second)

	//Запуск пула воркеров
	workerPool := worker.NewPool(5, 1000, deliverySender, deliveryRepo)
	workersCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	workerPool.Start(workersCtx)

	//Сервисы
	webhookService := service.NewWebhookService(webhookRepo)
	eventService := service.NewEventService(webhookRepo, deliveryRepo, workerPool)

	handler := httphandler.NewHandler(webhookService, eventService)
	
	//Роутинг
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhooks", handler.CreateWebhook)
	mux.HandleFunc("POST /events", handler.TriggerEvent)

	slog.Info("server starting", "addr", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}