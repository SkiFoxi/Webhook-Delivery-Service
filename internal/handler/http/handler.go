package http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/service"
)

type Handler struct {
	webhookService *service.WebhookService
	eventService   *service.EventService
}

func NewHandler(webhookServ *service.WebhookService, eventService *service.EventService) *Handler {
	return &Handler{
		webhookService: webhookServ,
		eventService:   eventService,
	}
}

func (h *Handler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	var req CreateWebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	webhook, err := h.webhookService.Create(r.Context(), req.URL, req.Events)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidURL):
			http.Error(w, "url is required", http.StatusBadRequest)
		case errors.Is(err, domain.ErrNoEvents):
			http.Error(w, "at least one event required", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	response := WebhookResponse{
		ID:        webhook.ID,
		URL:       webhook.URL,
		Events:    webhook.Events,
		CreatedAt: webhook.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response) // Можно делать и через json.Marshal и w.Write()

}

func (h *Handler) TriggerEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}
	var req EventRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}

	if err := h.eventService.Trigger(r.Context(), req.Type, req.Payload); err != nil {
		slog.Error("trigger error", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})

}
