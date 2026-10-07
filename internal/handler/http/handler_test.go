package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httphandler "github.com/SkiFoxi/Webhook-Delivery-Service/internal/handler/http"
	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/repository/memory"
)

func TestCreateWebhook(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{{
		name:       "success",
		method:     http.MethodPost,
		body:       `{"url": "https://example.com/hooks", "events": ["user.created"]}`,
		wantStatus: http.StatusCreated,
	},
	{
		name:       "empty url",
		method:     http.MethodPost,
		body:       `{"url": "", "events": ["user.created"]}`,
		wantStatus: http.StatusBadRequest,
	},
	{
		name:       "no events",
		method:     http.MethodPost,
		body:       `{"url": "https://example.com", "events": []}`,
		wantStatus: http.StatusBadRequest,
	},
	{
		name:       "invalid json",
		method:     http.MethodPost,
		body:       `{not json`,
		wantStatus: http.StatusBadRequest,
	},
	{
		name:       "wrong method",
		method:     http.MethodGet,
		body:       `{}`,
		wantStatus: http.StatusMethodNotAllowed,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webhookRepo := memory.NewWebhookRepo()
			h := httphandler.NewHandler(webhookRepo)

			req := httptest.NewRequest(tt.method, "/webhooks", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.CreateWebhook(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusCreated {
				var resp httphandler.WebhookResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if resp.ID == "" {
					t.Error("ID is empty")
				}
			}
		})
	}
}
