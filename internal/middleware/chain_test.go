package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/SkiFoxi/Webhook-Delivery-Service/internal/middleware"
)

func TestChain(t *testing.T) {
	var order []string

	makeMW := func(name string) middleware.Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+"-before")
				next.ServeHTTP(w, r)
				order = append(order, name+"-after")
			})
		}
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	chained := middleware.Chain(handler, makeMW("A"), makeMW("B"))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	chained.ServeHTTP(rec, req)

	want := []string{"A-before", "B-before", "handler", "B-after", "A-after"}
	if !slices.Equal(order, want) {
		t.Errorf("order = %v; want %v", order, want)
	}
}
