package main

import (
	"log/slog"
	"net/http"
	"os"

	httphandler "github.com/SkiFoxi/Webhook-Delivery-Service/internal/handler/http"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/repository/memory"
)

func main() {
	webhookRepo := memory.NewWebhookRepo()
	handler := httphandler.NewHandler(webhookRepo)
	
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhooks", handler.CreateWebhook)

	slog.Info("server starting", "addr", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}