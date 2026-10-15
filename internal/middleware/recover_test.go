package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/middleware"
)

func TestRecover(t *testing.T) {
    // Хендлер, который паникует
    panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        panic("test panic")
    })

    recoveredHandler := middleware.Recover(panicHandler)

    req := httptest.NewRequest(http.MethodGet, "/test", nil)
    rec := httptest.NewRecorder()

    // Не должно паниковать
    recoveredHandler.ServeHTTP(rec, req)

    if rec.Code != http.StatusInternalServerError {
        t.Errorf("status = %d; want 500", rec.Code)
    }
}