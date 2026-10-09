package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()
	//Механизм Shutdown на 10 секунд
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan // Блокируется здесь до получения системного сигнала
	slog.Info("signal received", "signal", sig)
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("http shutdown error", "err", err)
	}
	slog.Info("http server stopped")

	slog.Info("stopping worker pool...")
	workerPool.Stop()
	slog.Info("worker pool stopped")
	cancelWorkers() //Страховка
}
