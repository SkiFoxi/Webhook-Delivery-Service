package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/middleware"
)

func TestAuth(t *testing.T) {
    const apiKey = "secret-key"

    // Хендлер, который ставит 200
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Оборачиваем в Auth
    protectedHandler := middleware.Auth(apiKey)(nextHandler)

    tests := []struct {
        name       string
        apiKey     string   // ключ, который пришлёт клиент
        wantStatus int
    }{
        {"no key", "", http.StatusUnauthorized},
        {"wrong key", "wrong", http.StatusUnauthorized},
        {"correct key", apiKey, http.StatusOK},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, "/test", nil)
            if tt.apiKey != "" {
                req.Header.Set("X-API-Key", tt.apiKey)
            }
            rec := httptest.NewRecorder()

            protectedHandler.ServeHTTP(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
            }
        })
    }
}