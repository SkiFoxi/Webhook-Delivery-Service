package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/middleware"
)

func TestLogging(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	loggedHandler := middleware.Logging(handler)

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d; want 201", rec.Code)
	}
}
