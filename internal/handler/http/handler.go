package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
	"uuid"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/domain"
)

type Handler struct {
	WebhookRepo domain.WebhookRepository
}

func NewHandler(WebhookRepo domain.WebhookRepository) *Handler {
	return &Handler{
		WebhookRepo: WebhookRepo,
	}
}

func (h *Handler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	var req CreateWebhookRequest
	if err := json.Unmarshal(body, req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.URL == "" && len(req.Events) == 0 {
		http.Error(w, "invalid url ", http.StatusBadRequest)
		return
	}

	webhook := &domain.Webhook{
		ID:        uuid.New().String(),
		URL:       req.URL,
		Events:    req.Events,
		Secret:    uuid.New().String(),
		CreatedAt: time.Now(),
	}

	if err := h.WebhookRepo.Create(r.Context(), webhook); err != nil {
		http.Error(w, "cannot creat webhook", http.StatusInternalServerError)
		return
	}

	response := WebhookResponse{
		ID:       webhook.ID,
		URL:      webhook.URL,
		Events:   webhook.Events,
		CreateAt: webhook.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response) // Можно делать и через json.Marshal и w.Write()

}
